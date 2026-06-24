package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/Victor118/liquidity/app"
	"github.com/Victor118/liquidity/x/liquidity"
	"github.com/Victor118/liquidity/x/liquidity/keeper"
	"github.com/Victor118/liquidity/x/liquidity/types"
)

func TestWithdrawRatioInvariant(t *testing.T) {
	require.NotPanics(t, func() {
		keeper.WithdrawAmountInvariant(math.NewInt(1), math.NewInt(1), math.NewInt(2), math.NewInt(3), math.NewInt(1), math.NewInt(2), types.DefaultParams().WithdrawFeeRate)
	})
	require.Panics(t, func() {
		keeper.WithdrawAmountInvariant(math.NewInt(1), math.NewInt(1), math.NewInt(2), math.NewInt(5), math.NewInt(1), math.NewInt(2), types.DefaultParams().WithdrawFeeRate)
	})
}

func TestMintingPoolCoinsInvariant(t *testing.T) {
	for _, tc := range []struct {
		poolCoinSupply  int64
		mintingPoolCoin int64
		reserveA        int64
		depositA        int64
		refundedA       int64
		reserveB        int64
		depositB        int64
		refundedB       int64
		expectPanic     bool
	}{
		{
			10000, 1000,
			100000, 10000, 0,
			100000, 10000, 0,
			false,
		},
		{
			10000, 1000,
			100000, 10000, 5000,
			100000, 10000, 300,
			true,
		},
		{
			3000000, 100,
			100000000, 4000, 667,
			200000000, 8000, 1334,
			false,
		},
		{
			3000000, 100,
			100000000, 4000, 0,
			200000000, 8000, 1334,
			true,
		},
	} {
		f := require.NotPanics
		if tc.expectPanic {
			f = require.Panics
		}
		f(t, func() {
			keeper.MintingPoolCoinsInvariant(
				math.NewInt(tc.poolCoinSupply),
				math.NewInt(tc.mintingPoolCoin),
				math.NewInt(tc.depositA),
				math.NewInt(tc.depositB),
				math.NewInt(tc.reserveA),
				math.NewInt(tc.reserveB),
				math.NewInt(tc.refundedA),
				math.NewInt(tc.refundedB),
			)
		})
	}
}

func TestLiquidityPoolsEscrowAmountInvariant(t *testing.T) {
	simapp, ctx := app.CreateTestInput(t)

	// define test denom X, Y for Liquidity Pool
	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)

	X := math.NewInt(1000000000)
	Y := math.NewInt(1000000000)

	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, math.NewInt(10000))
	poolID := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])

	// begin block, init
	app.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolID, true)
	app.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolID, true)

	invariant := keeper.AllInvariants(simapp.LiquidityKeeper)
	_, broken := invariant(ctx)
	require.False(t, broken)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	_, broken = invariant(ctx)
	require.False(t, broken)

	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
	_, broken = invariant(ctx)
	require.False(t, broken)

	price, _ := math.LegacyNewDecFromStr("1.1")
	priceY, _ := math.LegacyNewDecFromStr("1.2")
	xOfferCoins := []sdk.Coin{sdk.NewCoin(denomX, math.NewInt(10000))}
	yOfferCoins := []sdk.Coin{sdk.NewCoin(denomY, math.NewInt(5000))}
	xOrderPrices := []math.LegacyDec{price}
	yOrderPrices := []math.LegacyDec{priceY}
	xOrderAddrs := addrs[1:2]
	yOrderAddrs := addrs[2:3]
	app.TestSwapPool(t, simapp, ctx, xOfferCoins, xOrderPrices, xOrderAddrs, poolID, false)
	app.TestSwapPool(t, simapp, ctx, xOfferCoins, xOrderPrices, xOrderAddrs, poolID, false)
	app.TestSwapPool(t, simapp, ctx, xOfferCoins, xOrderPrices, xOrderAddrs, poolID, false)
	app.TestSwapPool(t, simapp, ctx, yOfferCoins, yOrderPrices, yOrderAddrs, poolID, false)

	_, broken = invariant(ctx)
	require.False(t, broken)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	_, broken = invariant(ctx)
	require.False(t, broken)

	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
	_, broken = invariant(ctx)
	require.False(t, broken)

	batchEscrowAcc := simapp.AccountKeeper.GetModuleAddress(types.ModuleName)
	escrowAmt := simapp.BankKeeper.GetAllBalances(ctx, batchEscrowAcc)
	require.NotEmpty(t, escrowAmt)
	err := simapp.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addrs[0],
		sdk.NewCoins(sdk.NewCoin(xOfferCoins[0].Denom, xOfferCoins[0].Amount.QuoRaw(2))))
	require.NoError(t, err)
	escrowAmt = simapp.BankKeeper.GetAllBalances(ctx, batchEscrowAcc)

	msg, broken := invariant(ctx)
	require.True(t, broken)
	require.Equal(t, "liquidity: batch escrow amount invariant broken invariant\nbatch escrow amount LT batch remaining amount\n", msg)
}
