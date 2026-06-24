package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/Victor118/liquidity/x/liquidity/types"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.LegacyParamChange {
	return []simtypes.LegacyParamChange{
		simulation.NewSimLegacyParamChange(types.ModuleName, string(types.KeyMinInitDepositAmount),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMinInitDepositAmount(r).Int64())
			},
		),
		simulation.NewSimLegacyParamChange(types.ModuleName, string(types.KeyInitPoolCoinMintAmount),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenInitPoolCoinMintAmount(r).Int64())
			},
		),
		simulation.NewSimLegacyParamChange(types.ModuleName, string(types.KeyMaxReserveCoinAmount),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMaxReserveCoinAmount(r).Int64())
			},
		),
		simulation.NewSimLegacyParamChange(types.ModuleName, string(types.KeySwapFeeRate),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenSwapFeeRate(r))
			},
		),
		simulation.NewSimLegacyParamChange(types.ModuleName, string(types.KeyWithdrawFeeRate),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenWithdrawFeeRate(r))
			},
		),
		simulation.NewSimLegacyParamChange(types.ModuleName, string(types.KeyMaxOrderAmountRatio),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenMaxOrderAmountRatio(r))
			},
		),
		simulation.NewSimLegacyParamChange(types.ModuleName, string(types.KeyUnitBatchHeight),
			func(r *rand.Rand) string {
				return fmt.Sprintf("%d", GenUnitBatchHeight(r))
			},
		),
	}
}
