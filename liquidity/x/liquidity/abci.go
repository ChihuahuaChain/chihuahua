package liquidity

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"

	"github.com/Victor118/liquidity/x/liquidity/keeper"
	"github.com/Victor118/liquidity/x/liquidity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// In the Begin blocker of the liquidity module,
// Reinitialize batch messages that were not executed in the previous batch and delete batch messages that were executed or ready to delete.
func BeginBlocker(ctx context.Context, k keeper.Keeper) {

	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	context := sdk.UnwrapSDKContext(ctx)
	k.DeleteAndInitPoolBatches(context)
}

// In case of deposit, withdraw, and swap msgs, unlike other normal tx msgs,
// collect them in the liquidity pool batch and perform an execution once at the endblock to calculate and use the universal price.
func EndBlocker(ctx context.Context, k keeper.Keeper) {
	context := sdk.UnwrapSDKContext(ctx)
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)
	k.ExecutePoolBatches(context)
}
