package app

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"

	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"

	"github.com/ChihuahuaChain/chihuahua/x/feeburn"
	feeburnante "github.com/ChihuahuaChain/chihuahua/x/feeburn/ante"
	feeburntypes "github.com/ChihuahuaChain/chihuahua/x/feeburn/types"
)

// Regression tests for the x/feeburn fee deduction. They pin the current
// behaviour so the SDK v0.54 port can be checked against it.

func setupFeeburn(t *testing.T, burnPercent string) (*App, sdk.Context, sdk.AccAddress) {
	t.Helper()
	app := Setup(t)
	ctx := app.BaseApp.NewContext(false)
	require.NoError(t, app.FeeburnKeeper.SetParams(ctx, feeburntypes.NewParams(burnPercent, nil)))

	addr := simtestutil.CreateIncrementalAccounts(1)[0]
	initAccountWithCoins(app, ctx, addr, sdk.NewCoins(
		sdk.NewInt64Coin("uhuahua", 1_000_000),
		sdk.NewInt64Coin("uother", 1_000_000),
	))
	return app, ctx, addr
}

func TestDeductFeesBurnsPercentage(t *testing.T) {
	tests := []struct {
		name        string
		burnPercent string
		fee         sdk.Coins
		wantBurned  sdk.Coins
	}{
		{"no burn", "0", sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 1000)), sdk.NewCoins()},
		{"half burn", "50", sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 1000)), sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 500))},
		{"full burn", "100", sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 1000)), sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 1000))},
		{"burn rounds down", "50", sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 3)), sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 1))},
		{
			"every fee denom is burned",
			"50",
			sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 1000), sdk.NewInt64Coin("uother", 200)),
			sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 500), sdk.NewInt64Coin("uother", 100)),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app, ctx, addr := setupFeeburn(t, tc.burnPercent)
			feeCollector := app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)

			supplyBefore := supplyOf(app, ctx, tc.fee)
			payerBefore := app.BankKeeper.GetAllBalances(ctx, addr)
			collectorBefore := app.BankKeeper.GetAllBalances(ctx, feeCollector)

			bp, ok := sdkmath.NewIntFromString(tc.burnPercent)
			require.True(t, ok)
			acc := app.AccountKeeper.GetAccount(ctx, addr)
			require.NoError(t, feeburnante.DeductFees(app.BankKeeper, ctx, acc, tc.fee, bp))

			require.Equal(t, payerBefore.Sub(tc.fee...), app.BankKeeper.GetAllBalances(ctx, addr))
			require.Equal(t, collectorBefore.Add(tc.fee...).Sub(tc.wantBurned...), app.BankKeeper.GetAllBalances(ctx, feeCollector))
			require.Equal(t, supplyBefore.Sub(tc.wantBurned...), supplyOf(app, ctx, tc.fee))
		})
	}
}

func TestDeductFeesInsufficientFunds(t *testing.T) {
	app, ctx, addr := setupFeeburn(t, "50")
	acc := app.AccountKeeper.GetAccount(ctx, addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 2_000_000))
	require.Error(t, feeburnante.DeductFees(app.BankKeeper, ctx, acc, fee, sdkmath.NewInt(50)))
	require.Equal(t, sdkmath.NewInt(1_000_000), app.BankKeeper.GetBalance(ctx, addr, "uhuahua").Amount)
}

func TestDeductFeeDecoratorBurnsTxFee(t *testing.T) {
	app, ctx, addr := setupFeeburn(t, "50")
	feeCollector := app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	fee := sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 10_000))

	txBuilder := app.TxConfig().NewTxBuilder()
	require.NoError(t, txBuilder.SetMsgs(banktypes.NewMsgSend(addr, addr, sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 1)))))
	txBuilder.SetFeeAmount(fee)
	txBuilder.SetGasLimit(200_000)

	supplyBefore := app.BankKeeper.GetSupply(ctx, "uhuahua")
	collectorBefore := app.BankKeeper.GetBalance(ctx, feeCollector, "uhuahua")

	dfd := feeburnante.NewDeductFeeDecorator(app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, nil, app.FeeburnKeeper)
	_, err := dfd.AnteHandle(ctx, txBuilder.GetTx(), false, func(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) {
		return ctx, nil
	})
	require.NoError(t, err)

	require.Equal(t, sdkmath.NewInt(1_000_000-10_000), app.BankKeeper.GetBalance(ctx, addr, "uhuahua").Amount)
	require.Equal(t, collectorBefore.Amount.AddRaw(5_000), app.BankKeeper.GetBalance(ctx, feeCollector, "uhuahua").Amount)
	require.Equal(t, supplyBefore.Amount.SubRaw(5_000), app.BankKeeper.GetSupply(ctx, "uhuahua").Amount)
}

func TestFeeburnParamsValidation(t *testing.T) {
	for _, v := range []string{"0", "1", "50", "100"} {
		require.NoError(t, feeburntypes.NewParams(v, nil).Validate(), v)
	}
	for _, v := range []string{"-1", "101", "", "abc", "0.5"} {
		require.Error(t, feeburntypes.NewParams(v, nil).Validate(), v)
	}
}

// TestBurnAddressBurnsUhuahuaInEndBlock pins the public burn address behaviour:
// uhuahua and token factory denoms sent to the feeburn module account are
// burned out of the supply in EndBlock and recorded in the burned total, while
// IBC vouchers and any other denom are left alone.
func TestBurnAddressBurnsUhuahuaInEndBlock(t *testing.T) {
	const (
		factoryDenom = "factory/chihuahua1kjey0s32mpsfq5sseazjshulxlc54fg7uspjac/cryptobank"
		ibcDenom     = "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2"
	)

	app := Setup(t)
	ctx := app.BaseApp.NewContext(false)

	// the upgrade handler materializes the module account; mirror that here so
	// the burn address is a module account before it receives funds
	app.AccountKeeper.GetModuleAccount(ctx, feeburntypes.ModuleName)

	burnAddr := app.AccountKeeper.GetModuleAddress(feeburntypes.ModuleName)
	initAccountWithCoins(app, ctx, burnAddr, sdk.NewCoins(
		sdk.NewInt64Coin("uhuahua", 1_000),
		sdk.NewInt64Coin(factoryDenom, 700),
		sdk.NewInt64Coin(ibcDenom, 300),
		sdk.NewInt64Coin("uother", 500),
	))

	huahuaSupplyBefore := app.BankKeeper.GetSupply(ctx, "uhuahua")
	factorySupplyBefore := app.BankKeeper.GetSupply(ctx, factoryDenom)
	ibcSupplyBefore := app.BankKeeper.GetSupply(ctx, ibcDenom)
	require.True(t, app.FeeburnKeeper.GetTotalBurned(ctx).IsZero())

	// Drive EndBlock through the module manager, exactly as a live block does.
	// A direct module.EndBlock(ctx) call would pass even when feeburn is not
	// wired as an end blocker: that is how the SDK v0.54 manager silently
	// skipped it (its EndBlock signature matched no end-block interface) and
	// nothing was ever burned on chain. Going through the manager pins the
	// wiring too.
	_, err := app.mm.EndBlock(ctx)
	require.NoError(t, err)

	// uhuahua and the factory denom at the burn address are gone; the IBC
	// voucher and uother are untouched
	require.True(t, app.BankKeeper.GetBalance(ctx, burnAddr, "uhuahua").IsZero())
	require.True(t, app.BankKeeper.GetBalance(ctx, burnAddr, factoryDenom).IsZero())
	require.Equal(t, int64(300), app.BankKeeper.GetBalance(ctx, burnAddr, ibcDenom).Amount.Int64())
	require.Equal(t, int64(500), app.BankKeeper.GetBalance(ctx, burnAddr, "uother").Amount.Int64())

	// only the burnable denoms left the supply, and exactly those were recorded
	require.Equal(t, huahuaSupplyBefore.Amount.SubRaw(1_000).Int64(), app.BankKeeper.GetSupply(ctx, "uhuahua").Amount.Int64())
	require.Equal(t, factorySupplyBefore.Amount.SubRaw(700).Int64(), app.BankKeeper.GetSupply(ctx, factoryDenom).Amount.Int64())
	require.Equal(t, ibcSupplyBefore.Amount.Int64(), app.BankKeeper.GetSupply(ctx, ibcDenom).Amount.Int64())
	require.Equal(t, sdk.NewCoins(
		sdk.NewInt64Coin("uhuahua", 1_000),
		sdk.NewInt64Coin(factoryDenom, 700),
	), app.FeeburnKeeper.GetTotalBurned(ctx))

	// a second EndBlock with no burnable balance left is a no-op
	_, err = app.mm.EndBlock(ctx)
	require.NoError(t, err)
	require.Equal(t, sdk.NewCoins(
		sdk.NewInt64Coin("uhuahua", 1_000),
		sdk.NewInt64Coin(factoryDenom, 700),
	), app.FeeburnKeeper.GetTotalBurned(ctx))
}

// TestBurnAddressIsNotBlocked ensures the public burn address can receive funds,
// so a community-pool spend (or anyone) can send uhuahua to it.
func TestBurnAddressIsNotBlocked(t *testing.T) {
	burnAddr := authtypes.NewModuleAddress(feeburntypes.ModuleName).String()
	require.False(t, BlockedAddresses()[burnAddr], "burn address must be able to receive funds")
}

func supplyOf(app *App, ctx sdk.Context, denoms sdk.Coins) sdk.Coins {
	supply := sdk.NewCoins()
	for _, c := range denoms {
		supply = supply.Add(app.BankKeeper.GetSupply(ctx, c.Denom))
	}
	return supply
}

func TestDeductFeeDecoratorFeePayerEvent(t *testing.T) {
	app, ctx, addr := setupFeeburn(t, "50")
	ctx = ctx.WithEventManager(sdk.NewEventManager())

	txBuilder := app.TxConfig().NewTxBuilder()
	require.NoError(t, txBuilder.SetMsgs(banktypes.NewMsgSend(addr, addr, sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 1)))))
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 100)))
	txBuilder.SetGasLimit(200_000)

	dfd := feeburnante.NewDeductFeeDecorator(app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, nil, app.FeeburnKeeper)
	_, err := dfd.AnteHandle(ctx, txBuilder.GetTx(), false, func(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) {
		return ctx, nil
	})
	require.NoError(t, err)

	for _, ev := range ctx.EventManager().Events() {
		if v, ok := ev.GetAttribute(sdk.AttributeKeyFeePayer); ok {
			require.Equal(t, addr.String(), v.Value)
			return
		}
	}
	t.Fatal("fee_payer attribute not emitted")
}

func TestDeductFeeDecoratorRecordsBurnedFees(t *testing.T) {
	app, ctx, addr := setupFeeburn(t, "50")

	txBuilder := app.TxConfig().NewTxBuilder()
	require.NoError(t, txBuilder.SetMsgs(banktypes.NewMsgSend(addr, addr, sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 1)))))
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 10_000), sdk.NewInt64Coin("uother", 300)))
	txBuilder.SetGasLimit(200_000)

	dfd := feeburnante.NewDeductFeeDecorator(app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, nil, app.FeeburnKeeper)
	next := func(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) { return ctx, nil }

	require.True(t, app.FeeburnKeeper.GetTotalBurned(ctx).IsZero())
	for i := 0; i < 2; i++ {
		_, err := dfd.AnteHandle(ctx, txBuilder.GetTx(), false, next)
		require.NoError(t, err)
	}
	want := sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 10_000), sdk.NewInt64Coin("uother", 300))
	require.Equal(t, want, app.FeeburnKeeper.GetTotalBurned(ctx))

	// the recorded total matches what left the supply
	res, err := app.FeeburnKeeper.TotalBurned(ctx, &feeburntypes.QueryTotalBurnedRequest{})
	require.NoError(t, err)
	require.Equal(t, want, res.TotalBurned)
}

func TestDeductFeeDecoratorNoBurnRecordsNothing(t *testing.T) {
	app, ctx, addr := setupFeeburn(t, "0")

	txBuilder := app.TxConfig().NewTxBuilder()
	require.NoError(t, txBuilder.SetMsgs(banktypes.NewMsgSend(addr, addr, sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 1)))))
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 10_000)))
	txBuilder.SetGasLimit(200_000)

	dfd := feeburnante.NewDeductFeeDecorator(app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, nil, app.FeeburnKeeper)
	_, err := dfd.AnteHandle(ctx, txBuilder.GetTx(), false, func(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) {
		return ctx, nil
	})
	require.NoError(t, err)
	require.True(t, app.FeeburnKeeper.GetTotalBurned(ctx).IsZero())
}

func TestRecordingBurnedFeesUsesNoGas(t *testing.T) {
	app, ctx, _ := setupFeeburn(t, "50")
	ctx = ctx.WithGasMeter(storetypes.NewGasMeter(1_000_000))

	app.FeeburnKeeper.AddBurned(ctx, sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 7)))
	app.FeeburnKeeper.AddBurned(ctx, sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 5)))
	require.Zero(t, ctx.GasMeter().GasConsumed())
	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 12)), app.FeeburnKeeper.GetTotalBurned(ctx))
}

// buildFeeTx builds a MsgSend tx with the given fee and gas limit.
func buildFeeTx(t *testing.T, app *App, from sdk.AccAddress, fee sdk.Coins, gas uint64) sdk.Tx {
	t.Helper()
	txBuilder := app.TxConfig().NewTxBuilder()
	require.NoError(t, txBuilder.SetMsgs(banktypes.NewMsgSend(from, from, sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 1)))))
	txBuilder.SetFeeAmount(fee)
	txBuilder.SetGasLimit(gas)
	return txBuilder.GetTx()
}

// TestDeductFeeDecoratorEnforcesMinGasPricesInDeliverTx pins the consensus
// floor: in DeliverTx (not a CheckTx), a fee below the on-chain min_gas_prices
// is rejected regardless of the node's local config.
func TestDeductFeeDecoratorEnforcesMinGasPricesInDeliverTx(t *testing.T) {
	app := Setup(t)
	ctx := app.BaseApp.NewContext(false).WithBlockHeight(1) // DeliverTx, height > 0
	require.False(t, ctx.IsCheckTx())

	floor := sdk.NewDecCoins(sdk.NewDecCoin("uhuahua", sdkmath.NewInt(500)))
	require.NoError(t, app.FeeburnKeeper.SetParams(ctx, feeburntypes.NewParams("50", floor)))

	addr := simtestutil.CreateIncrementalAccounts(1)[0]
	initAccountWithCoins(app, ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 1_000_000_000)))

	dfd := feeburnante.NewDeductFeeDecorator(app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, nil, app.FeeburnKeeper)
	next := func(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) { return ctx, nil }

	// required = ceil(500 * 200000) = 100_000_000 uhuahua
	_, err := dfd.AnteHandle(ctx, buildFeeTx(t, app, addr, sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 99_999_999)), 200_000), false, next)
	require.ErrorIs(t, err, sdkerrors.ErrInsufficientFee)

	_, err = dfd.AnteHandle(ctx, buildFeeTx(t, app, addr, sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 100_000_000)), 200_000), false, next)
	require.NoError(t, err)
}

// TestDeductFeeDecoratorMinGasPricesCheckTxUsesMax pins that in CheckTx the
// effective floor is the per-denom max of the global param and the validator's
// local minimum-gas-prices.
func TestDeductFeeDecoratorMinGasPricesCheckTxUsesMax(t *testing.T) {
	app := Setup(t)
	ctx := app.BaseApp.NewContext(true).WithBlockHeight(1). // CheckTx
								WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoin("uhuahua", sdkmath.NewInt(500))))
	require.True(t, ctx.IsCheckTx())

	// global floor (100) is lower than the local one (500): 500 must win.
	floor := sdk.NewDecCoins(sdk.NewDecCoin("uhuahua", sdkmath.NewInt(100)))
	require.NoError(t, app.FeeburnKeeper.SetParams(ctx, feeburntypes.NewParams("50", floor)))

	addr := simtestutil.CreateIncrementalAccounts(1)[0]
	initAccountWithCoins(app, ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 1_000_000_000)))

	dfd := feeburnante.NewDeductFeeDecorator(app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, nil, app.FeeburnKeeper)
	next := func(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) { return ctx, nil }

	// above the global floor (20_000_000) but below the local one (100_000_000): rejected.
	_, err := dfd.AnteHandle(ctx, buildFeeTx(t, app, addr, sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 20_000_001)), 200_000), false, next)
	require.ErrorIs(t, err, sdkerrors.ErrInsufficientFee)

	_, err = dfd.AnteHandle(ctx, buildFeeTx(t, app, addr, sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 100_000_000)), 200_000), false, next)
	require.NoError(t, err)
}

// TestDeductFeeDecoratorNoFloorWhenEmpty pins that an empty min_gas_prices means
// no floor, so a tiny fee passes (current mainnet behaviour before the upgrade).
func TestDeductFeeDecoratorNoFloorWhenEmpty(t *testing.T) {
	app, ctx, addr := setupFeeburn(t, "50") // NewParams(..., nil) → empty floor
	ctx = ctx.WithBlockHeight(1)

	dfd := feeburnante.NewDeductFeeDecorator(app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, nil, app.FeeburnKeeper)
	next := func(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) { return ctx, nil }

	_, err := dfd.AnteHandle(ctx, buildFeeTx(t, app, addr, sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 1)), 200_000), false, next)
	require.NoError(t, err)
}

// TestFeeburnMinGasPricesValidation covers the new min_gas_prices param validation.
func TestFeeburnMinGasPricesValidation(t *testing.T) {
	valid := []sdk.DecCoins{
		nil,
		sdk.NewDecCoins(sdk.NewDecCoin("uhuahua", sdkmath.NewInt(500))),
	}
	for _, v := range valid {
		require.NoError(t, feeburntypes.NewParams("50", v).Validate())
	}

	bad := []sdk.DecCoins{
		{sdk.DecCoin{Denom: "uhuahua", Amount: sdkmath.LegacyNewDec(-1)}},                                            // negative amount
		{sdk.DecCoin{Denom: "1bad", Amount: sdkmath.LegacyNewDec(1)}},                                                // invalid denom
		{{Denom: "btc", Amount: sdkmath.LegacyNewDec(1)}, {Denom: "atom", Amount: sdkmath.LegacyNewDec(1)}},          // unsorted
		sdk.NewDecCoins(sdk.NewDecCoin("uother", sdkmath.NewInt(1))),                                                 // non-native denom
		sdk.NewDecCoins(sdk.NewDecCoin("uhuahua", sdkmath.NewInt(500)), sdk.NewDecCoin("uother", sdkmath.NewInt(1))), // native + non-native
	}
	for _, v := range bad {
		require.Error(t, feeburntypes.NewParams("50", v).Validate())
	}
}

// TestFeeburnMinGasPricesMigration mirrors the v10.0.1 upgrade step: it sets the
// floor while preserving the existing TxFeeBurnPercent.
func TestFeeburnMinGasPricesMigration(t *testing.T) {
	app := Setup(t)
	ctx := app.BaseApp.NewContext(false)

	// pre-upgrade params: only the burn percent, no floor
	require.NoError(t, app.FeeburnKeeper.SetParams(ctx, feeburntypes.NewParams("50", nil)))

	minGasPrices, err := sdk.ParseDecCoins(RecommendedMinGasPrices)
	require.NoError(t, err)
	params := app.FeeburnKeeper.GetParams(ctx)
	params.MinGasPrices = minGasPrices
	require.NoError(t, app.FeeburnKeeper.SetParams(ctx, params))

	got := app.FeeburnKeeper.GetParams(ctx)
	require.Equal(t, "50", got.TxFeeBurnPercent, "burn percent must be preserved")
	require.Equal(t, sdk.NewDecCoins(sdk.NewDecCoin("uhuahua", sdkmath.NewInt(500))), got.MinGasPrices)
}

func TestFeeburnGenesisKeepsTotalBurned(t *testing.T) {
	app, ctx, _ := setupFeeburn(t, "50")
	total := sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 123_456_789), sdk.NewInt64Coin("uother", 5))
	require.NoError(t, app.FeeburnKeeper.SetTotalBurned(ctx, total))

	gs := feeburn.ExportGenesis(ctx, app.FeeburnKeeper)
	require.NoError(t, gs.Validate())
	require.Equal(t, total, gs.TotalBurned)

	app2 := Setup(t)
	ctx2 := app2.BaseApp.NewContext(false)
	feeburn.InitGenesis(ctx2, app2.FeeburnKeeper, *gs)
	require.Equal(t, total, app2.FeeburnKeeper.GetTotalBurned(ctx2))

	// SetTotalBurned replaces the total, dropping denoms that are no longer in it
	require.NoError(t, app2.FeeburnKeeper.SetTotalBurned(ctx2, sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 1))))
	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 1)), app2.FeeburnKeeper.GetTotalBurned(ctx2))
}
