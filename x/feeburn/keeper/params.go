package keeper

import (
	"github.com/ChihuahuaChain/chihuahua/x/feeburn/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.TxFeeBurnPercent(ctx),
	)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// TxFeeBurnPercent returns the TxFeeBurnPercent param
func (k Keeper) TxFeeBurnPercent(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyTxFeeBurnPercent, &res)
	return
}
