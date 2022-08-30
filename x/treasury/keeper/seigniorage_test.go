package keeper

import (
	"math/rand"
	"testing"

	core "github.com/ChihuahuaChain/chihuahua/types"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestSettle(t *testing.T) {
	input := CreateTestInput(t)

	faucetBalance := input.BankKeeper.GetBalance(input.Ctx, input.AccountKeeper.GetModuleAddress(faucetAccountName), core.MicroHuahuaDenom)
	burnAmt := sdk.NewInt(rand.Int63()%faucetBalance.Amount.Int64() + 1)
	initialHuahuaSupply := input.BankKeeper.GetSupply(input.Ctx, core.MicroHuahuaDenom)
	// input.TreasuryKeeper.RecordEpochInitialIssuance(input.Ctx)

	input.Ctx = input.Ctx.WithBlockHeight(int64(core.BlocksPerWeek))
	err := input.BankKeeper.BurnCoins(input.Ctx, faucetAccountName, sdk.NewCoins(sdk.NewCoin(core.MicroHuahuaDenom, burnAmt)))
	require.NoError(t, err)

	// check seigniorage update
	require.Equal(t, burnAmt, input.TreasuryKeeper.PeekEpochSeigniorage(input.Ctx))

	input.TreasuryKeeper.SettleSeigniorage(input.Ctx)
	lunaSupply := input.BankKeeper.GetSupply(input.Ctx, core.MicroHuahuaDenom)
	feePool := input.DistrKeeper.GetFeePool(input.Ctx)

	// Reward weight portion of seigniorage burned
	rewardWeight := input.TreasuryKeeper.GetRewardWeight(input.Ctx)
	communityPoolAmt := burnAmt.Sub(rewardWeight.MulInt(burnAmt).TruncateInt())

	require.Equal(t, lunaSupply.Amount, initialHuahuaSupply.Amount.Sub(burnAmt).Add(communityPoolAmt))
	require.Equal(t, communityPoolAmt, feePool.CommunityPool.AmountOf(core.MicroHuahuaDenom).TruncateInt())
}

func TestOneRewardWeightSettle(t *testing.T) {
	input := CreateTestInput(t)

	// set zero reward weight
	input.TreasuryKeeper.SetRewardWeight(input.Ctx, sdk.OneDec())

	faucetBalance := input.BankKeeper.GetBalance(input.Ctx, input.AccountKeeper.GetModuleAddress(faucetAccountName), core.MicroHuahuaDenom)
	burnAmt := sdk.NewInt(rand.Int63()%faucetBalance.Amount.Int64() + 1)
	initialHuahuaSupply := input.BankKeeper.GetSupply(input.Ctx, core.MicroHuahuaDenom)
	// input.TreasuryKeeper.RecordEpochInitialIssuance(input.Ctx)

	input.Ctx = input.Ctx.WithBlockHeight(int64(core.BlocksPerWeek))
	err := input.BankKeeper.BurnCoins(input.Ctx, faucetAccountName, sdk.NewCoins(sdk.NewCoin(core.MicroHuahuaDenom, burnAmt)))
	require.NoError(t, err)

	// check seigniorage update
	require.Equal(t, burnAmt, input.TreasuryKeeper.PeekEpochSeigniorage(input.Ctx))

	input.TreasuryKeeper.SettleSeigniorage(input.Ctx)
	lunaSupply := input.BankKeeper.GetSupply(input.Ctx, core.MicroHuahuaDenom)
	feePool := input.DistrKeeper.GetFeePool(input.Ctx)

	// Reward weight portion of seigniorage burned
	require.Equal(t, lunaSupply.Amount, initialHuahuaSupply.Amount.Sub(burnAmt))
	require.Equal(t, sdk.ZeroInt(), feePool.CommunityPool.AmountOf(core.MicroHuahuaDenom).TruncateInt())
}
