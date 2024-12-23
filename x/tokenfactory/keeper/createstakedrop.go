package keeper

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	"github.com/ChihuahuaChain/chihuahua/app/params"
	"github.com/ChihuahuaChain/chihuahua/x/tokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) CreateStakedropByDenom(ctx sdk.Context, creatorAddr string, amount sdk.Coin, startBlock uint64, endBlock uint64) error {
	err := k.validateCreateStakedrop(ctx, creatorAddr, amount, startBlock, endBlock)
	if err != nil {
		return err
	}

	err = k.chargeForCreateStakedrop(ctx, creatorAddr, startBlock, endBlock)
	if err != nil {
		return err
	}

	err = k.createStakedropAfterValidation(ctx, amount, creatorAddr, startBlock, endBlock)
	return err
}

func (k Keeper) createStakedropAfterValidation(ctx sdk.Context, amount sdk.Coin, creatorAddr string, startBlock uint64, endBlock uint64) error {

	seq, err := k.getNextStakedropSequence(ctx)
	if err != nil {
		return err
	}
	if amount.Denom == params.BondDenom {
		//native token, we can mint it so sender need to send it to module account
		sender := sdk.MustAccAddressFromBech32(creatorAddr)
		k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(amount))
	}
	amountPerBlock := amount.Amount.Quo(math.NewInt(int64(endBlock - startBlock)))
	newStakedrop := types.Stakedrop{
		Amount:         amount,
		AmountPerBlock: sdk.NewCoin(amount.Denom, amountPerBlock),
		StartBlock:     int64(startBlock),
		EndBlock:       int64(endBlock),
	}
	key := collections.Join(startBlock, seq)
	return k.ActiveStakedrop.Set(ctx, key, newStakedrop)

}

func (k Keeper) validateCreateStakedrop(ctx sdk.Context, creatorAddr string, amount sdk.Coin, startBlock uint64, endBlock uint64) error {
	moduleParams := k.GetParams(ctx)
	if !isWhitelistedAddress(creatorAddr, moduleParams) {
		return types.ErrUnauthorized
	}

	if params.BondDenom == amount.Denom {
		return nil
	}
	//verify sender has created a subdenom (amount.Denom)
	creator, _, err := types.DeconstructDenom(amount.Denom)
	if err != nil {
		return err
	}
	if creator != creatorAddr {
		return types.ErrInvalidCreator
	}
	_, found := k.bankKeeper.GetDenomMetaData(ctx, amount.Denom)
	if !found {
		return types.ErrDenomDoesNotExist
	}

	if ctx.BlockHeight() > int64(startBlock) || startBlock >= endBlock {
		return types.ErrBadBlockParameters
	}

	return nil
}

func (k Keeper) chargeForCreateStakedrop(ctx sdk.Context, creatorAddr string, startBlock uint64, endBlock uint64) (err error) {
	params := k.GetParams(ctx)

	// if StakedropChargePerBlock is non-zero, transfer the tokens from the creator
	// account to community pool
	if (params.StakedropChargePerBlock != sdk.Coin{} && params.StakedropChargePerBlock.Amount.GT(math.NewInt(0))) {
		accAddr, err := sdk.AccAddressFromBech32(creatorAddr)
		if err != nil {
			return err
		}
		charge := params.StakedropChargePerBlock.Amount.Mul(math.NewIntFromUint64(endBlock - startBlock))

		if err := k.communityPoolKeeper.FundCommunityPool(ctx, sdk.NewCoins(sdk.NewCoin(params.StakedropChargePerBlock.Denom, charge)), accAddr); err != nil {
			return err
		}
	}

	// if DenomCreationGasConsume is non-zero, consume the gas
	if params.DenomCreationGasConsume != 0 {
		ctx.GasMeter().ConsumeGas(params.DenomCreationGasConsume, "consume denom creation gas")
	}

	return nil
}
