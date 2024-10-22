package keeper

import (
	"cosmossdk.io/math"
	"github.com/ChihuahuaChain/chihuahua/x/tokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) CreateStakedrop(ctx sdk.Context, creatorAddr string, amount sdk.Coin, startBlock uint64, endBlock uint64) error {
	denom, err := k.validateCreateStakedrop(ctx, creatorAddr, amount, startBlock, endBlock)
	if err != nil {
		return err
	}

	err = k.chargeForCreateStakedrop(ctx, creatorAddr, startBlock, endBlock)
	if err != nil {
		return err
	}

	err = k.createStakedropAfterValidation(ctx, creatorAddr, denom)
	return err
}

func (k Keeper) createStakedropAfterValidation(ctx sdk.Context, amount sdk.Coin, startBlock uint64, endBlock uint64) error {

}

func (k Keeper) validateCreateStakedrop(ctx sdk.Context, creatorAddr string, amount sdk.Coin, startBlock uint64, endBlock uint64) (string, error) {
	//verify sender has created a subdenom (amount.Denom)
	denom, err := types.GetTokenDenom(creatorAddr, amount.Denom)
	if err != nil {
		return "", err
	}
	_, found := k.bankKeeper.GetDenomMetaData(ctx, denom)
	if !found {
		return "", types.ErrDenomDoesNotExist
	}

	if ctx.BlockHeight() > int64(startBlock) || startBlock >= endBlock {
		return "", types.ErrBadBlockParameters
	}

	return denom, nil
}

func (k Keeper) chargeForCreateStakedrop(ctx sdk.Context, creatorAddr string, startBlock uint64, endBlock uint64) (err error) {
	params := k.GetParams(ctx)

	// if StakedropChargePerBlock is non-zero, transfer the tokens from the creator
	// account to community pool
	if params.StakedropChargePerBlock.Amount.GT(math.NewInt(0)) {
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
