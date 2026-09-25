package app

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"

	appparams "github.com/ChihuahuaChain/chihuahua/app/params"
	"github.com/ChihuahuaChain/chihuahua/x/tokenfactory"
	tfkeeper "github.com/ChihuahuaChain/chihuahua/x/tokenfactory/keeper"
	tftypes "github.com/ChihuahuaChain/chihuahua/x/tokenfactory/types"
)

// Regression tests for the chihuahua-specific parts of x/tokenfactory
// (stakedrops, builders commission, free mint whitelist). They pin the current
// behaviour so the SDK v0.54 port can be checked against it.

type tfFixture struct {
	app         *App
	ctx         sdk.Context
	msgServer   tftypes.MsgServer
	whitelisted sdk.AccAddress
	user        sdk.AccAddress
	builders    []sdk.AccAddress
}

func setupTokenfactory(t *testing.T) tfFixture {
	t.Helper()
	app := Setup(t)
	ctx := app.BaseApp.NewContext(false).WithBlockHeight(100)

	addrs := simtestutil.CreateIncrementalAccounts(4)
	f := tfFixture{
		app:         app,
		ctx:         ctx,
		msgServer:   tfkeeper.NewMsgServerImpl(app.TokenFactoryKeeper),
		whitelisted: addrs[0],
		user:        addrs[1],
		builders:    addrs[2:],
	}

	params := tftypes.DefaultParams()
	params.DenomCreationFee = nil
	params.DenomCreationGasConsume = 0
	params.BuildersCommission = sdkmath.LegacyNewDecWithPrec(1, 2) // 1%
	params.BuildersAddresses = []tftypes.WeightedAddress{
		{Address: f.builders[0].String(), Weight: sdkmath.LegacyNewDecWithPrec(10, 2)},
		{Address: f.builders[1].String(), Weight: sdkmath.LegacyNewDecWithPrec(90, 2)},
	}
	params.FreeMintWhitelistAddresses = []string{f.whitelisted.String()}
	params.StakedropChargePerBlock = sdk.NewInt64Coin(appparams.BondDenom, 0)
	require.NoError(t, app.TokenFactoryKeeper.SetParams(ctx, params))

	for _, a := range addrs {
		initAccountWithCoins(app, ctx, a, sdk.NewCoins(sdk.NewInt64Coin(appparams.BondDenom, 1_000_000)))
	}
	return f
}

func (f tfFixture) createDenom(t *testing.T, creator sdk.AccAddress, subdenom string) string {
	t.Helper()
	res, err := f.msgServer.CreateDenom(f.ctx, tftypes.NewMsgCreateDenom(creator.String(), subdenom))
	require.NoError(t, err)
	return res.NewTokenDenom
}

func (f tfFixture) mint(t *testing.T, admin sdk.AccAddress, coin sdk.Coin) {
	t.Helper()
	_, err := f.msgServer.Mint(f.ctx, tftypes.NewMsgMint(admin.String(), coin))
	require.NoError(t, err)
}

func (f tfFixture) balance(addr sdk.AccAddress, denom string) sdkmath.Int {
	return f.app.BankKeeper.GetBalance(f.ctx, addr, denom).Amount
}

func (f tfFixture) feeCollector() sdk.AccAddress {
	return f.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
}

func (f tfFixture) moduleAddr() sdk.AccAddress {
	return f.app.AccountKeeper.GetModuleAddress(tftypes.ModuleName)
}

// runBeginBlocker runs the tokenfactory BeginBlocker at the given height and
// returns how much of denom the fee collector received.
func (f tfFixture) runBeginBlocker(t *testing.T, height int64, denom string) sdkmath.Int {
	t.Helper()
	ctx := f.ctx.WithBlockHeight(height)
	before := f.app.BankKeeper.GetBalance(ctx, f.feeCollector(), denom).Amount
	require.NoError(t, tokenfactory.BeginBlocker(ctx, f.app.TokenFactoryKeeper, f.app.BankKeeper))
	return f.app.BankKeeper.GetBalance(ctx, f.feeCollector(), denom).Amount.Sub(before)
}

func (f tfFixture) activeStakedrops(t *testing.T) int {
	t.Helper()
	iter, err := f.app.TokenFactoryKeeper.ActiveStakedrop.Iterate(f.ctx, nil)
	require.NoError(t, err)
	defer iter.Close()
	keys, err := iter.Keys()
	require.NoError(t, err)
	return len(keys)
}

func TestMintBuildersCommission(t *testing.T) {
	f := setupTokenfactory(t)
	denom := f.createDenom(t, f.user, "dog")

	f.mint(t, f.user, sdk.NewInt64Coin(denom, 10_000))

	// 1% commission = 100, split 10% / 90% between builders
	require.Equal(t, sdkmath.NewInt(9_900), f.balance(f.user, denom))
	require.Equal(t, sdkmath.NewInt(10), f.balance(f.builders[0], denom))
	require.Equal(t, sdkmath.NewInt(90), f.balance(f.builders[1], denom))
	require.Equal(t, sdkmath.NewInt(10_000), f.app.BankKeeper.GetSupply(f.ctx, denom).Amount)
	require.True(t, f.balance(f.moduleAddr(), denom).IsZero())
}

func TestMintWhitelistedHasNoCommission(t *testing.T) {
	f := setupTokenfactory(t)
	denom := f.createDenom(t, f.whitelisted, "cat")

	f.mint(t, f.whitelisted, sdk.NewInt64Coin(denom, 10_000))

	require.Equal(t, sdkmath.NewInt(10_000), f.balance(f.whitelisted, denom))
	require.True(t, f.balance(f.builders[0], denom).IsZero())
	require.True(t, f.balance(f.builders[1], denom).IsZero())
}

func TestMintOnlyAdmin(t *testing.T) {
	f := setupTokenfactory(t)
	denom := f.createDenom(t, f.user, "dog")

	_, err := f.msgServer.Mint(f.ctx, tftypes.NewMsgMint(f.whitelisted.String(), sdk.NewInt64Coin(denom, 1)))
	require.ErrorIs(t, err, tftypes.ErrUnauthorized)
}

func TestStakedropTokenfactoryDenom(t *testing.T) {
	f := setupTokenfactory(t)
	denom := f.createDenom(t, f.whitelisted, "drop")
	f.mint(t, f.whitelisted, sdk.NewInt64Coin(denom, 1_005))

	start, end := f.ctx.BlockHeight()+1, f.ctx.BlockHeight()+11
	_, err := f.msgServer.CreateStakeDrop(f.ctx, tftypes.NewMsgCreateStakeDrop(f.whitelisted.String(), sdk.NewInt64Coin(denom, 1_005), start, end))
	require.NoError(t, err)
	require.True(t, f.balance(f.whitelisted, denom).IsZero())
	require.Equal(t, sdkmath.NewInt(1_005), f.balance(f.moduleAddr(), denom))

	// nothing before start
	require.True(t, f.runBeginBlocker(t, start-1, denom).IsZero())

	// 100 per block from start to end-1
	for h := start; h < end; h++ {
		require.Equal(t, sdkmath.NewInt(100), f.runBeginBlocker(t, h, denom), "height %d", h)
	}
	// the remainder at end
	require.Equal(t, sdkmath.NewInt(5), f.runBeginBlocker(t, end, denom))
	require.Equal(t, 1, f.activeStakedrops(t))

	// removed the block after end, nothing more is paid
	require.True(t, f.runBeginBlocker(t, end+1, denom).IsZero())
	require.Equal(t, 0, f.activeStakedrops(t))
	require.True(t, f.balance(f.moduleAddr(), denom).IsZero())
}

func TestStakedropNativeDenom(t *testing.T) {
	f := setupTokenfactory(t)
	bond := appparams.BondDenom

	start, end := f.ctx.BlockHeight()+1, f.ctx.BlockHeight()+5
	_, err := f.msgServer.CreateStakeDrop(f.ctx, tftypes.NewMsgCreateStakeDrop(f.whitelisted.String(), sdk.NewInt64Coin(bond, 1_002), start, end))
	require.NoError(t, err)
	require.Equal(t, sdkmath.NewInt(1_000_000-1_002), f.balance(f.whitelisted, bond))

	total := sdkmath.ZeroInt()
	for h := start; h < end; h++ {
		got := f.runBeginBlocker(t, h, bond)
		require.Equal(t, sdkmath.NewInt(250), got, "height %d", h)
		total = total.Add(got)
	}
	total = total.Add(f.runBeginBlocker(t, end, bond))
	require.Equal(t, sdkmath.NewInt(1_002), total)

	require.True(t, f.runBeginBlocker(t, end+1, bond).IsZero())
	require.Equal(t, 0, f.activeStakedrops(t))
}

func TestStakedropRejections(t *testing.T) {
	f := setupTokenfactory(t)
	h := f.ctx.BlockHeight()
	userDenom := f.createDenom(t, f.user, "dog")
	f.mint(t, f.user, sdk.NewInt64Coin(userDenom, 1_000))
	wlDenom := f.createDenom(t, f.whitelisted, "drop")
	f.mint(t, f.whitelisted, sdk.NewInt64Coin(wlDenom, 1_000))

	tests := []struct {
		name string
		msg  *tftypes.MsgCreateStakeDrop
	}{
		{"sender not whitelisted", tftypes.NewMsgCreateStakeDrop(f.user.String(), sdk.NewInt64Coin(userDenom, 100), h+1, h+10)},
		{"not admin of denom", tftypes.NewMsgCreateStakeDrop(f.whitelisted.String(), sdk.NewInt64Coin(userDenom, 100), h+1, h+10)},
		{"start in the past", tftypes.NewMsgCreateStakeDrop(f.whitelisted.String(), sdk.NewInt64Coin(wlDenom, 100), h-1, h+10)},
		{"start equals end", tftypes.NewMsgCreateStakeDrop(f.whitelisted.String(), sdk.NewInt64Coin(wlDenom, 100), h+5, h+5)},
		{"insufficient funds", tftypes.NewMsgCreateStakeDrop(f.whitelisted.String(), sdk.NewInt64Coin(wlDenom, 10_000), h+1, h+10)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := f.msgServer.CreateStakeDrop(f.ctx, tc.msg)
			require.Error(t, err)
		})
	}
	require.Equal(t, 0, f.activeStakedrops(t))
}

func TestStakedropChargePerBlock(t *testing.T) {
	f := setupTokenfactory(t)
	params := f.app.TokenFactoryKeeper.GetParams(f.ctx)
	params.StakedropChargePerBlock = sdk.NewInt64Coin(appparams.BondDenom, 10)
	require.NoError(t, f.app.TokenFactoryKeeper.SetParams(f.ctx, params))

	denom := f.createDenom(t, f.whitelisted, "drop")
	f.mint(t, f.whitelisted, sdk.NewInt64Coin(denom, 1_000))
	poolBefore, err := f.app.DistrKeeper.FeePool.Get(f.ctx)
	require.NoError(t, err)

	h := f.ctx.BlockHeight()
	_, err = f.msgServer.CreateStakeDrop(f.ctx, tftypes.NewMsgCreateStakeDrop(f.whitelisted.String(), sdk.NewInt64Coin(denom, 1_000), h+1, h+11))
	require.NoError(t, err)

	// 10 blocks * 10 uhuahua to the community pool
	require.Equal(t, sdkmath.NewInt(1_000_000-100), f.balance(f.whitelisted, appparams.BondDenom))
	poolAfter, err := f.app.DistrKeeper.FeePool.Get(f.ctx)
	require.NoError(t, err)
	require.Equal(t, sdkmath.LegacyNewDec(100), poolAfter.CommunityPool.AmountOf(appparams.BondDenom).Sub(poolBefore.CommunityPool.AmountOf(appparams.BondDenom)))
}

func TestStakedropNativeDenomValidatesBlocks(t *testing.T) {
	f := setupTokenfactory(t)
	h := f.ctx.BlockHeight()
	coin := sdk.NewInt64Coin(appparams.BondDenom, 1_000)

	for name, msg := range map[string]*tftypes.MsgCreateStakeDrop{
		"start in the past": tftypes.NewMsgCreateStakeDrop(f.whitelisted.String(), coin, h-1, h+10),
		"start equals end":  tftypes.NewMsgCreateStakeDrop(f.whitelisted.String(), coin, h+5, h+5),
	} {
		_, err := f.msgServer.CreateStakeDrop(f.ctx, msg)
		require.ErrorIs(t, err, tftypes.ErrBadBlockParameters, name)
	}
}

// A stakedrop starting at the current height would miss its first block, as
// the BeginBlocker for that height already ran before the tx was executed.
func TestStakedropCannotStartAtCurrentHeight(t *testing.T) {
	f := setupTokenfactory(t)
	h := f.ctx.BlockHeight()
	denom := f.createDenom(t, f.whitelisted, "drop")
	f.mint(t, f.whitelisted, sdk.NewInt64Coin(denom, 1_000))

	for _, coin := range []sdk.Coin{sdk.NewInt64Coin(denom, 100), sdk.NewInt64Coin(appparams.BondDenom, 100)} {
		_, err := f.msgServer.CreateStakeDrop(f.ctx, tftypes.NewMsgCreateStakeDrop(f.whitelisted.String(), coin, h, h+10))
		require.ErrorIs(t, err, tftypes.ErrBadBlockParameters, coin.Denom)
	}
}

func TestMintAndBurnEventAddresses(t *testing.T) {
	f := setupTokenfactory(t)
	denom := f.createDenom(t, f.whitelisted, "dog")

	ctx := f.ctx.WithEventManager(sdk.NewEventManager())
	_, err := f.msgServer.Mint(ctx, &tftypes.MsgMint{Sender: f.whitelisted.String(), Amount: sdk.NewInt64Coin(denom, 100), MintToAddress: f.user.String()})
	require.NoError(t, err)
	requireEventAttr(t, ctx, tftypes.TypeMsgMint, tftypes.AttributeMintToAddress, f.user.String())

	ctx = f.ctx.WithEventManager(sdk.NewEventManager())
	_, err = f.msgServer.Burn(ctx, tftypes.NewMsgBurn(f.whitelisted.String(), sdk.NewInt64Coin(denom, 0)))
	require.NoError(t, err)
	requireEventAttr(t, ctx, tftypes.TypeMsgBurn, tftypes.AttributeBurnFromAddress, f.whitelisted.String())
}

func requireEventAttr(t *testing.T, ctx sdk.Context, eventType, key, want string) {
	t.Helper()
	for _, ev := range ctx.EventManager().Events() {
		if ev.Type != eventType {
			continue
		}
		v, ok := ev.GetAttribute(key)
		require.True(t, ok, "attribute %s missing on %s", key, eventType)
		require.Equal(t, want, v.Value)
		return
	}
	t.Fatalf("event %s not emitted", eventType)
}
