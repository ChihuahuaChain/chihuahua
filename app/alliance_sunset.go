package app

import (
	"bytes"
	"math"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	alliancetypes "github.com/terra-money/alliance/x/alliance/types"
)

// SunsetAlliance shuts x/alliance down so the module can be removed in a later
// upgrade (governance proposal 99). Everything the module holds is burned:
//
//  1. the virtual stake minted by the module is unbonded and burned on every
//     validator, bonded or not, together with its last distribution rewards
//  2. every balance of the alliance and alliance_rewards module accounts is
//     burned: the staked and unbonding ampGASH, the unclaimed rewards
//     (HUAHUA, ampGASH and factory tokens) and any dust
//  3. the alliance store is cleared, keeping only the module params
//
// The alliance undelegation code is not used: it fails on the mainnet state
// because of rounding inconsistencies in the alliance share accounting.
func (app *App) SunsetAlliance(ctx sdk.Context) error {
	logger := ctx.Logger().With("upgrade", "alliance-sunset")

	burned, err := app.unbondAllianceVirtualStake(ctx)
	if err != nil {
		return err
	}
	logger.Info("burned alliance virtual stake", "amount", burned)

	allianceAddr := app.AccountKeeper.GetModuleAddress(alliancetypes.ModuleName)
	rewardsAddr := app.AccountKeeper.GetModuleAddress(alliancetypes.RewardsPoolName)
	// the rewards pool account cannot burn: gather everything in the alliance account first
	if rewards := app.BankKeeper.GetAllBalances(ctx, rewardsAddr); !rewards.IsZero() {
		if err := app.BankKeeper.SendCoinsFromModuleToModule(ctx, alliancetypes.RewardsPoolName, alliancetypes.ModuleName, rewards); err != nil {
			return err
		}
	}
	if balance := app.BankKeeper.GetAllBalances(ctx, allianceAddr); !balance.IsZero() {
		if err := app.BankKeeper.BurnCoins(ctx, alliancetypes.ModuleName, balance); err != nil {
			return err
		}
		logger.Info("burned alliance module balances", "amount", balance)
	}

	deleted, err := app.clearAllianceStore(ctx)
	if err != nil {
		return err
	}
	logger.Info("cleared alliance store", "keys", deleted)
	return nil
}

// unbondAllianceVirtualStake unbonds and burns every native delegation held by
// the alliance module account, returning the burned amount. The pending
// distribution rewards are withdrawn to the module account first.
func (app *App) unbondAllianceVirtualStake(ctx sdk.Context) (sdkmath.Int, error) {
	moduleAddr := app.AccountKeeper.GetModuleAddress(alliancetypes.ModuleName)
	bondDenom, err := app.StakingKeeper.BondDenom(ctx)
	if err != nil {
		return sdkmath.Int{}, err
	}

	delegations, err := app.StakingKeeper.GetDelegatorDelegations(ctx, moduleAddr, math.MaxUint16)
	if err != nil {
		return sdkmath.Int{}, err
	}

	burned := sdkmath.ZeroInt()
	for _, d := range delegations {
		valAddr, err := sdk.ValAddressFromBech32(d.ValidatorAddress)
		if err != nil {
			return sdkmath.Int{}, err
		}
		if _, err := app.DistrKeeper.WithdrawDelegationRewards(ctx, moduleAddr, valAddr); err != nil {
			return sdkmath.Int{}, err
		}
		validator, err := app.StakingKeeper.GetValidator(ctx, valAddr)
		if err != nil {
			return sdkmath.Int{}, err
		}
		amount, err := app.StakingKeeper.Unbond(ctx, moduleAddr, valAddr, d.Shares)
		if err != nil {
			return sdkmath.Int{}, err
		}
		if amount.IsZero() {
			continue
		}
		pool := stakingtypes.NotBondedPoolName
		if validator.IsBonded() {
			pool = stakingtypes.BondedPoolName
		}
		if err := app.BankKeeper.BurnCoins(ctx, pool, sdk.NewCoins(sdk.NewCoin(bondDenom, amount))); err != nil {
			return sdkmath.Int{}, err
		}
		burned = burned.Add(amount)
	}
	return burned, nil
}

// clearAllianceStore deletes every alliance record (assets, validator infos,
// delegations, redelegations, undelegation queues, snapshots) except the
// module params, returning the number of deleted keys.
func (app *App) clearAllianceStore(ctx sdk.Context) (int, error) {
	store := runtime.KVStoreAdapter(app.AllianceKeeper.StoreService().OpenKVStore(ctx))
	iter := store.Iterator(nil, nil)
	var keys [][]byte
	for ; iter.Valid(); iter.Next() {
		if bytes.HasPrefix(iter.Key(), alliancetypes.ParamsKey) || bytes.HasPrefix(iter.Key(), alliancetypes.ModuleAccKey) {
			continue
		}
		keys = append(keys, iter.Key())
	}
	if err := iter.Close(); err != nil {
		return 0, err
	}
	for _, key := range keys {
		store.Delete(key)
	}
	return len(keys), nil
}
