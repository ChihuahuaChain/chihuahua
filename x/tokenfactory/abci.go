package tokenfactory

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	"github.com/ChihuahuaChain/chihuahua/x/tokenfactory/keeper"
	"github.com/ChihuahuaChain/chihuahua/x/tokenfactory/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	for ; iter.Valid(); iter.Next() {
		stakeDrop, err := iter.KeyValue()
		if err != nil {

			return err
		}

		if ctx.BlockHeight() >= stakeDrop.Value.StartBlock && ctx.BlockHeight() < (stakeDrop.Value.EndBlock) {
			ctx.Logger().Error(fmt.Sprintf("Stakedrop %v %s in Block %d", stakeDrop.Value.AmountPerBlock.Amount, stakeDrop.Value.Amount.Denom, ctx.BlockHeight()))
			coinsToSend = coinsToSend.Add(stakeDrop.Value.AmountPerBlock)

		} else if ctx.BlockHeight() == stakeDrop.Value.EndBlock {
			restAmount := stakeDrop.Value.Amount.Amount.Sub(stakeDrop.Value.AmountPerBlock.Amount.Mul(math.NewInt(stakeDrop.Value.EndBlock - stakeDrop.Value.StartBlock)))
			ctx.Logger().Error(fmt.Sprintf("Last Stakedrop %v in Block %d", restAmount, ctx.BlockHeight()))
			coinsToSend = coinsToSend.Add(sdk.NewCoin(stakeDrop.Value.Amount.Denom, restAmount))

		} else if ctx.BlockHeight() > stakeDrop.Value.EndBlock {
			ctx.Logger().Error(fmt.Sprintf("Remove Stakedrop in Block %d", ctx.BlockHeight()))
			k.ActiveStakedrop.Remove(ctx, stakeDrop.Key)
		}
	}

	bankKeeper.MintCoins(ctx, types.ModuleName, coinsToSend)
	return bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.FeeCollectorName, sdk.NewCoins(coinsToSend...))

}
