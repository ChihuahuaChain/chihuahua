package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Victor118/liquidity/x/liquidity/types"
)

// Keeper of the liquidity store
type Keeper struct {
	cdc           codec.BinaryCodec
	storeKey      storetypes.StoreKey
	bankKeeper    types.BankKeeper
	accountKeeper types.AccountKeeper
	distrKeeper   types.DistributionKeeper
	// the address capable of executing a MsgUpdateParams message. Typically, this should be the x/gov module account.
	authority sdk.AccAddress
}

// NewKeeper returns a liquidity keeper. It handles:
// - creating new ModuleAccounts for each pool ReserveAccount
// - sending to and from ModuleAccounts
// - minting, burning PoolCoins
func NewKeeper(cdc codec.BinaryCodec, key storetypes.StoreKey, bankKeeper types.BankKeeper, accountKeeper types.AccountKeeper, distrKeeper types.DistributionKeeper, authority sdk.AccAddress) Keeper {
	// ensure liquidity module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// ensure gov module account is set and is not nil
	if err := sdk.VerifyAddressFormat(authority); err != nil {
		panic(err)
	}

	return Keeper{
		storeKey:      key,
		bankKeeper:    bankKeeper,
		accountKeeper: accountKeeper,
		distrKeeper:   distrKeeper,
		cdc:           cdc,
		authority:     authority,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

// GetParams gets the parameters for the liquidity module.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams sets the parameters for the liquidity module.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}
	store.Set(types.ParamsKey, bz)

	return nil
}

// GetCircuitBreakerEnabled returns circuit breaker enabled param .
func (k Keeper) GetCircuitBreakerEnabled(ctx sdk.Context) (enabled bool) {
	enabled = k.GetParams(ctx).CircuitBreakerEnabled
	return
}
