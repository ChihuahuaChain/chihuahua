package v505

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

type CosMints struct {
	Address       string `json:"address"`
	AmountUhuahua string `json:"amount"`
}

func mintLostTokens(
	ctx sdk.Context,
	bankKeeper bankkeeper.Keeper,
	mintKeeper mintkeeper.Keeper,
) {
	var cosMints []CosMints
	err := json.Unmarshal([]byte(recordsJSONString), &cosMints)
	if err != nil {
		panic(fmt.Sprintf("error reading COS JSON: %+v", err))
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
	}	
}

func restoreDelegatorTokens(
	ctx sdk.Context,
	bankKeeper bankkeeper.Keeper,
	mintKeeper mintkeeper.Keeper,
) error {
	mintLostTokens(ctx, bankKeeper, mintKeeper)

	return nil
}