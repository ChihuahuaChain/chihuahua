package v310

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type CosMints struct {
	Address       string `json:"address"`
	AmountUhuahua string `json:"amount"`
}

var (
	cosValidatorAddress = "chihuahuavaloper12kupfxepn4lu2kjmed84rvhvjddy3wz6s92zvq"
	cosConsensusAddress = "chihuahuavalcons1d7pu4xvswxx4es7ntawfpqklp0tlv02dv8f6sf"
)

func mintLostTokens(
	ctx sdk.Context,
	bankKeeper bankkeeper.BaseKeeper,
	stakingKeeper stakingkeeper.Keeper,
	mintKeeper mintkeeper.Keeper,
) {
	var cosMints []CosMints
	err := json.Unmarshal([]byte(recordsJSONString), &cosMints)
	if err != nil {
		panic(fmt.Sprintf("error reading COS JSON: %+v", err))
	}

	cosValAddress, err := sdk.ValAddressFromBech32(cosValidatorAddress)
	if err != nil {
		panic(fmt.Sprintf("validator address is not valid bech32: %s", cosValAddress))
	}

	cosValidator, found := stakingKeeper.GetValidator(ctx, cosValAddress)
	if !found {
		panic(fmt.Sprintf("cos validator '%s' not found", cosValAddress))
	}

	for _, mintRecord := range cosMints {
		coinAmount, mintOk := sdk.NewIntFromString(mintRecord.AmountUhuahua)
		if !mintOk {
			panic(fmt.Sprintf("error parsing mint of %suhuahua to %s", mintRecord.AmountUhuahua, mintRecord.Address))
		}

		coin := sdk.NewCoin("uhuahua", coinAmount)
		coins := sdk.NewCoins(coin)

		err = mintKeeper.MintCoins(ctx, coins)
		if err != nil {
			panic(fmt.Sprintf("error minting %suhuahua to %s: %+v", mintRecord.AmountUhuahua, mintRecord.Address, err))
		}

		delegatorAccount, err := sdk.AccAddressFromBech32(mintRecord.Address)
		if err != nil {
			panic(fmt.Sprintf("error converting human address %s to sdk.AccAddress: %+v", mintRecord.Address, err))
		}

		err = bankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delegatorAccount, coins)
		if err != nil {
			panic(fmt.Sprintf("error sending minted %suhuahua to %s: %+v", mintRecord.AmountUhuahua, mintRecord.Address, err))
		}

		sdkAddress, err := sdk.AccAddressFromBech32(mintRecord.Address)
		if err != nil {
			panic(fmt.Sprintf("account address is not valid bech32: %s", mintRecord.Address))
		}

		_, err = stakingKeeper.Delegate(ctx, sdkAddress, coin.Amount, stakingtypes.Unbonded, cosValidator, true)
		if err != nil {
			panic(fmt.Sprintf("error delegating minted %suhuahua from %s to %s: %+v", mintRecord.AmountUhuahua, mintRecord.Address, cosValidatorAddress, err))
		}
	}
}

func revertTombstone(ctx sdk.Context, slashingKeeper slashingkeeper.Keeper) error {
	cosValAddress, err := sdk.ValAddressFromBech32(cosValidatorAddress)
	if err != nil {
		panic(fmt.Sprintf("validator address is not valid bech32: %s", cosValAddress))
	}

	cosConsAddress, err := sdk.ConsAddressFromBech32(cosConsensusAddress)
	if err != nil {
		panic(fmt.Sprintf("consensus address is not valid bech32: %s", cosValAddress))
	}

	// Revert Tombstone info
	slashingKeeper.RevertTombstone(ctx, cosConsAddress)

	// Set jail until=now, the validator then must unjail manually
	slashingKeeper.JailUntil(ctx, cosConsAddress, ctx.BlockTime())

	return nil
}

func RevertCosTombstoning(
	ctx sdk.Context,
	slashingKeeper slashingkeeper.Keeper,
	mintKeeper mintkeeper.Keeper,
	bankKeeper bankkeeper.BaseKeeper,
	stakingKeeper stakingkeeper.Keeper,
) error {
	err := revertTombstone(ctx, slashingKeeper)
	if err != nil {
		return err
	}

	mintLostTokens(ctx, bankKeeper, stakingKeeper, mintKeeper)

	return nil
}
