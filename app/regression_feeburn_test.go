package app

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"

	feeburnante "github.com/ChihuahuaChain/chihuahua/x/feeburn/ante"
	feeburntypes "github.com/ChihuahuaChain/chihuahua/x/feeburn/types"
)

// Regression tests for the x/feeburn fee deduction. They pin the current
// behaviour so the SDK v0.54 port can be checked against it.

func setupFeeburn(t *testing.T, burnPercent string) (*App, sdk.Context, sdk.AccAddress) {
	t.Helper()
	app := Setup(t)
	ctx := app.BaseApp.NewContext(false)
	require.NoError(t, app.FeeburnKeeper.SetParams(ctx, feeburntypes.NewParams(burnPercent)))

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
		require.NoError(t, feeburntypes.NewParams(v).Validate(), v)
	}
	for _, v := range []string{"-1", "101", "", "abc", "0.5"} {
		require.Error(t, feeburntypes.NewParams(v).Validate(), v)
	}
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
