package keeper

import (
	"fmt"

	"cosmossdk.io/log/v2"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ChihuahuaChain/chihuahua/x/feeburn/types"
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storetypes.StoreKey
		// the address capable of executing a MsgUpdateParams message. Typically, this should be the x/gov module account.
		authority sdk.AccAddress
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	authority sdk.AccAddress,
) *Keeper {
	// ensure gov module account is set and is not nil
	if err := sdk.VerifyAddressFormat(authority); err != nil {
		panic(err)
	}

	return &Keeper{
		cdc:       cdc,
		storeKey:  storeKey,
		authority: authority,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
