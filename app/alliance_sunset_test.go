package app

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/terra-money/alliance/x/alliance"
	alliancekeeper "github.com/terra-money/alliance/x/alliance/keeper"
	alliancetypes "github.com/terra-money/alliance/x/alliance/types"
)

const sunsetDenom = "ibc/AMPGASH"

type sunsetFixture struct {
	app          *App
	ctx          sdk.Context
	bondDenom    string
	valAddr      sdk.ValAddress
	users        []sdk.AccAddress
	allianceAddr sdk.AccAddress
}

// setupSunset creates the ampGASH-like alliance with three stakers, one of them
// unbonding, lets the rebalance mint the virtual stake and leaves more virtual
// stake on a validator that is not bonded.
func setupSunset(t *testing.T) sunsetFixture {
	t.Helper()
	app := Setup(t)
	ctx := app.BaseApp.NewContext(false).WithBlockHeight(10).WithBlockTime(time.Now().UTC())
	bondDenom, err := app.StakingKeeper.BondDenom(ctx)
	require.NoError(t, err)
	allianceMsgs := alliancekeeper.NewMsgServerImpl(app.AllianceKeeper)

	vals, err := app.StakingKeeper.GetBondedValidatorsByPower(ctx)
	require.NoError(t, err)
	require.Len(t, vals, 1)
	valAddr, err := sdk.ValAddressFromBech32(vals[0].OperatorAddress)
	require.NoError(t, err)

	require.NoError(t, app.AllianceKeeper.CreateAlliance(ctx, &alliancetypes.MsgCreateAllianceProposal{
		Denom:             sunsetDenom,
		RewardWeight:      sdkmath.LegacyMustNewDecFromStr("0.03"),
		RewardWeightRange: alliancetypes.RewardWeightRange{Min: sdkmath.LegacyMustNewDecFromStr("0.03"), Max: sdkmath.LegacyMustNewDecFromStr("0.03")},
		// a take rate makes the shares/tokens ratio non trivial, as on mainnet
		TakeRate:             sdkmath.LegacyMustNewDecFromStr("0.0000003"),
		RewardChangeRate:     sdkmath.LegacyOneDec(),
		RewardChangeInterval: time.Hour,
	}))

	users := simtestutil.CreateIncrementalAccounts(3)
	staked := []int64{17_193_737_763_233, 2_500_000_000_007, 700_000_000_001}
	for i, u := range users {
		initAccountWithCoins(app, ctx, u, sdk.NewCoins(sdk.NewInt64Coin(sunsetDenom, staked[i])))
		_, err := allianceMsgs.Delegate(ctx, alliancetypes.NewMsgDelegate(u.String(), valAddr.String(), sdk.NewInt64Coin(sunsetDenom, staked[i])))
		require.NoError(t, err)
	}
	// users[2] starts unbonding part of its stake: the tokens stay in the module until completion
	_, err = allianceMsgs.Undelegate(ctx, alliancetypes.NewMsgUndelegate(users[2].String(), valAddr.String(), sdk.NewInt64Coin(sunsetDenom, 123_456_789_000)))
	require.NoError(t, err)

	// let the rewards start and rebalance, which mints and delegates the virtual stake
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(8 * 24 * time.Hour)).WithBlockHeight(20)
	require.NoError(t, app.AllianceKeeper.QueueAssetRebalanceEvent(ctx))
	require.NoError(t, alliance.EndBlocker(ctx, app.AllianceKeeper))
	allianceAddr := app.AccountKeeper.GetModuleAddress(alliancetypes.ModuleName)
	virtual, err := app.AllianceKeeper.GetAllianceBondedAmount(ctx, allianceAddr)
	require.NoError(t, err)
	require.True(t, virtual.IsPositive(), "the rebalance should have minted virtual stake")

	// virtual stake left on a validator that is not bonded, which the rebalance ignores
	unbondedVal := createUnbondedValidator(t, app, ctx, bondDenom)
	leftover := sdkmath.NewInt(12_345)
	require.NoError(t, app.BankKeeper.MintCoins(ctx, alliancetypes.ModuleName, sdk.NewCoins(sdk.NewCoin(bondDenom, leftover))))
	_, err = app.StakingKeeper.Delegate(ctx, allianceAddr, leftover, stakingtypes.Unbonded, unbondedVal, true)
	require.NoError(t, err)

	return sunsetFixture{app: app, ctx: ctx.WithBlockHeight(30), bondDenom: bondDenom, valAddr: valAddr, users: users, allianceAddr: allianceAddr}
}

// allocateValidatorRewards gives distribution rewards to the validator, so the
// virtual stake of the alliance has pending rewards.
func (f sunsetFixture) allocateValidatorRewards(t *testing.T, amount int64) {
	t.Helper()
	coins := sdk.NewCoins(sdk.NewInt64Coin(f.bondDenom, amount))
	require.NoError(t, f.app.BankKeeper.MintCoins(f.ctx, minttypes.ModuleName, coins))
	require.NoError(t, f.app.BankKeeper.SendCoinsFromModuleToModule(f.ctx, minttypes.ModuleName, distrtypes.ModuleName, coins))
	val, err := f.app.StakingKeeper.GetValidator(f.ctx, f.valAddr)
	require.NoError(t, err)
	require.NoError(t, f.app.DistrKeeper.AllocateTokensToValidator(f.ctx, val, sdk.NewDecCoinsFromCoins(coins...)))
}

type sunsetBefore struct {
	userBalances []sdkmath.Int
	moduleAsset  sdkmath.Int
	assetSupply  sdkmath.Int
	bondSupply   sdkmath.Int
	virtual      sdkmath.Int
	fundBalance  sdkmath.Int
	params       alliancetypes.Params
}

func (f sunsetFixture) before() sunsetBefore {
	b := sunsetBefore{
		moduleAsset: f.app.BankKeeper.GetBalance(f.ctx, f.allianceAddr, sunsetDenom).Amount.
			Add(f.app.BankKeeper.GetBalance(f.ctx, f.app.AccountKeeper.GetModuleAddress(alliancetypes.RewardsPoolName), sunsetDenom).Amount),
		assetSupply: f.app.BankKeeper.GetSupply(f.ctx, sunsetDenom).Amount,
		bondSupply:  f.app.BankKeeper.GetSupply(f.ctx, f.bondDenom).Amount,
		params:      f.app.AllianceKeeper.GetParams(f.ctx),
		virtual:     sdkmath.ZeroInt(),
		fundBalance: f.app.BankKeeper.GetBalance(f.ctx, sdk.MustAccAddressFromBech32(EcosystemFundAddress), f.bondDenom).Amount,
	}
	delegations, err := f.app.StakingKeeper.GetDelegatorDelegations(f.ctx, f.allianceAddr, 100)
	if err != nil {
		panic(err)
	}
	for _, d := range delegations {
		valAddr, err := sdk.ValAddressFromBech32(d.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		val, err := f.app.StakingKeeper.GetValidator(f.ctx, valAddr)
		if err != nil {
			panic(err)
		}
		b.virtual = b.virtual.Add(val.TokensFromShares(d.Shares).TruncateInt())
	}
	for _, u := range f.users {
		b.userBalances = append(b.userBalances, f.app.BankKeeper.GetBalance(f.ctx, u, sunsetDenom).Amount)
	}
	return b
}

func (f sunsetFixture) requireSunsetDone(t *testing.T, b sunsetBefore) {
	t.Helper()
	app, ctx := f.app, f.ctx

	// nothing is returned: the staked and unbonding assets are burned
	for i, u := range f.users {
		require.Equal(t, b.userBalances[i].String(), app.BankKeeper.GetBalance(ctx, u, sunsetDenom).Amount.String(), "user %d", i)
	}
	require.True(t, b.moduleAsset.IsPositive())
	require.Equal(t, b.assetSupply.Sub(b.moduleAsset).String(), app.BankKeeper.GetSupply(ctx, sunsetDenom).Amount.String())

	// the virtual stake went to the ecosystem fund, bonded and not bonded alike
	delegations, err := app.StakingKeeper.GetDelegatorDelegations(ctx, f.allianceAddr, 100)
	require.NoError(t, err)
	require.Empty(t, delegations)
	require.True(t, b.virtual.IsPositive())
	fund := sdk.MustAccAddressFromBech32(EcosystemFundAddress)
	require.Equal(t, b.fundBalance.Add(b.virtual).String(), app.BankKeeper.GetBalance(ctx, fund, f.bondDenom).Amount.String())

	// the HUAHUA rewards were burned, the alliance accounts are empty
	require.True(t, app.BankKeeper.GetSupply(ctx, f.bondDenom).Amount.LT(b.bondSupply))
	require.True(t, app.BankKeeper.GetSupply(ctx, f.bondDenom).Amount.GT(b.bondSupply.Sub(b.virtual)))
	require.True(t, app.BankKeeper.GetAllBalances(ctx, f.allianceAddr).IsZero())
	require.True(t, app.BankKeeper.GetAllBalances(ctx, app.AccountKeeper.GetModuleAddress(alliancetypes.RewardsPoolName)).IsZero())

	// the store holds nothing but the params
	require.Empty(t, app.AllianceKeeper.GetAllAssets(ctx))
	var remaining int
	require.NoError(t, app.AllianceKeeper.IterateDelegations(ctx, func(alliancetypes.Delegation) bool {
		remaining++
		return false
	}))
	require.Zero(t, remaining)
	for _, u := range f.users {
		unbondings, err := app.AllianceKeeper.GetUnbondingsByDelegator(ctx, u)
		require.NoError(t, err)
		require.Empty(t, unbondings)
	}
	infos, err := app.AllianceKeeper.GetAllAllianceValidatorInfo(ctx)
	require.NoError(t, err)
	require.Empty(t, infos)
	require.Equal(t, b.params, app.AllianceKeeper.GetParams(ctx))

	// the module keeps running without assets, including its staking hooks
	next := ctx.WithBlockTime(ctx.BlockTime().Add(time.Hour)).WithBlockHeight(ctx.BlockHeight() + 1)
	require.NoError(t, alliance.EndBlocker(next, app.AllianceKeeper))
	require.NoError(t, app.AllianceKeeper.StakingHooks().BeforeValidatorSlashed(next, f.valAddr, sdkmath.LegacyNewDecWithPrec(1, 2)))
	require.NoError(t, alliance.EndBlocker(next.WithBlockHeight(next.BlockHeight()+1), app.AllianceKeeper))
}

func TestSunsetAlliance(t *testing.T) {
	f := setupSunset(t)
	f.allocateValidatorRewards(t, 1_000_000)
	b := f.before()

	require.NoError(t, f.app.SunsetAlliance(f.ctx))

	f.requireSunsetDone(t, b)
}

// Proposal 99 sets the reward weight and the take rate to 0.
func TestSunsetAllianceWithZeroRewardWeight(t *testing.T) {
	f := setupSunset(t)
	asset, found := f.app.AllianceKeeper.GetAssetByDenom(f.ctx, sunsetDenom)
	require.True(t, found)
	asset.RewardWeight = sdkmath.LegacyZeroDec()
	asset.RewardWeightRange = alliancetypes.RewardWeightRange{Min: sdkmath.LegacyZeroDec(), Max: sdkmath.LegacyZeroDec()}
	asset.TakeRate = sdkmath.LegacyZeroDec()
	require.NoError(t, f.app.AllianceKeeper.SetAsset(f.ctx, asset))
	f.allocateValidatorRewards(t, 1_000_000)

	b := f.before()
	require.NoError(t, f.app.SunsetAlliance(f.ctx))
	f.requireSunsetDone(t, b)
}

func createUnbondedValidator(t *testing.T, app *App, ctx sdk.Context, bondDenom string) stakingtypes.Validator {
	t.Helper()
	operator := simtestutil.CreateIncrementalAccounts(10)[9]
	selfBond := sdk.NewCoin(bondDenom, sdkmath.NewInt(1_000_000))
	require.NoError(t, app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(selfBond)))
	require.NoError(t, app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, operator, sdk.NewCoins(selfBond)))

	valAddr := sdk.ValAddress(operator)
	msg, err := stakingtypes.NewMsgCreateValidator(
		valAddr.String(), ed25519.GenPrivKey().PubKey(), selfBond,
		stakingtypes.NewDescription("unbonded", "", "", "", ""),
		stakingtypes.NewCommissionRates(sdkmath.LegacyNewDecWithPrec(1, 1), sdkmath.LegacyOneDec(), sdkmath.LegacyZeroDec()),
		sdkmath.OneInt(),
	)
	require.NoError(t, err)
	_, err = stakingkeeper.NewMsgServerImpl(app.StakingKeeper).CreateValidator(ctx, msg)
	require.NoError(t, err)

	// the staking EndBlocker does not run in this test, so the validator stays unbonded
	val, err := app.StakingKeeper.GetValidator(ctx, valAddr)
	require.NoError(t, err)
	require.False(t, val.IsBonded())
	return val
}
