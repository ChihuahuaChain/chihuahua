package tokenfactory

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	"github.com/ChihuahuaChain/chihuahua/app/params"
	"github.com/ChihuahuaChain/chihuahua/x/tokenfactory/keeper"
	"github.com/ChihuahuaChain/chihuahua/x/tokenfactory/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// BeginBlocker
func BeginBlocker(ctx sdk.Context, k keeper.Keeper, bankKeeper types.BankKeeper) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyEndBlocker)

	rng := collections.NewPrefixUntilPairRange[uint64, uint64](uint64(ctx.BlockHeight()))
	iter, err := k.ActiveStakedrop.Iterate(ctx, rng)

	if err != nil {

		return err
	}
	defer iter.Close()

	coinsToSend := sdk.NewCoins()
	nativeCoinToSend := sdk.NewCoin(params.BondDenom, math.NewInt(0))

	for ; iter.Valid(); iter.Next() {
		stakeDrop, err := iter.KeyValue()
		if err != nil {

			return err
		}

		if ctx.BlockHeight() >= stakeDrop.Value.StartBlock && ctx.BlockHeight() < (stakeDrop.Value.EndBlock) {
			if stakeDrop.Value.Amount.Denom == params.BondDenom {
				nativeCoinToSend = nativeCoinToSend.Add(stakeDrop.Value.AmountPerBlock)
			} else {
				coinsToSend = coinsToSend.Add(stakeDrop.Value.AmountPerBlock)
			}

		} else if ctx.BlockHeight() == stakeDrop.Value.EndBlock {
			restAmount := stakeDrop.Value.Amount.Amount.Sub(stakeDrop.Value.AmountPerBlock.Amount.Mul(math.NewInt(stakeDrop.Value.EndBlock - stakeDrop.Value.StartBlock)))
			if stakeDrop.Value.Amount.Denom == params.BondDenom {

			} else {
				coinsToSend = coinsToSend.Add(sdk.NewCoin(stakeDrop.Value.Amount.Denom, restAmount))
			}

		} else if ctx.BlockHeight() > stakeDrop.Value.EndBlock {
			k.ActiveStakedrop.Remove(ctx, stakeDrop.Key)
		}
	}
	if !coinsToSend.Empty() {
		bankKeeper.MintCoins(ctx, types.ModuleName, coinsToSend)
		err = bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.FeeCollectorName, sdk.NewCoins(coinsToSend...))
		if err != nil {
			return err
		}
	}
	if nativeCoinToSend.Amount.IsPositive() {
		spendableAmount := bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(types.ModuleName), params.BondDenom)
		if spendableAmount.IsLT(nativeCoinToSend) {
			nativeCoinToSend = spendableAmount
		}
		if nativeCoinToSend.IsPositive() {
			err = bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.FeeCollectorName, sdk.NewCoins(nativeCoinToSend))
			if err != nil {
				return err
			}
		}

	}

	return nil
}
