package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/collections/indexes"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/ChihuahuaChain/chihuahua/x/tokenfactory/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeKey     storetypes.StoreKey
		storeService store.KVStoreService

		accountKeeper       types.AccountKeeper
		bankKeeper          types.BankKeeper
		communityPoolKeeper types.CommunityPoolKeeper

		enabledCapabilities []string

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority         string
		StakedropSequence collections.Sequence
		ActiveStakedrop   *collections.IndexedMap[collections.Pair[uint64, uint64], types.Stakedrop, StakedropIndexes]
		FeeCollectorName  string
	}
)

type StakedropIndexes struct {
	StakedropByDenom *indexes.Multi[string, collections.Pair[uint64, uint64], types.Stakedrop]
}

func (stkdIndexes StakedropIndexes) IndexesList() []collections.Index[collections.Pair[uint64, uint64], types.Stakedrop] {
	return []collections.Index[collections.Pair[uint64, uint64], types.Stakedrop]{stkdIndexes.StakedropByDenom}
}

// NewKeeper returns a new instance of the x/tokenfactory keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	storeService store.KVStoreService,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	communityPoolKeeper types.CommunityPoolKeeper,
	enabledCapabilities []string,
	feeCollectorName string,
	authority string,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:                 cdc,
		storeKey:            storeKey,
		storeService:        storeService,
		accountKeeper:       accountKeeper,
		bankKeeper:          bankKeeper,
		communityPoolKeeper: communityPoolKeeper,

		enabledCapabilities: enabledCapabilities,

		authority:         authority,
		FeeCollectorName:  feeCollectorName,
		StakedropSequence: collections.NewSequence(sb, types.AirdropSequenceKey, "airdrop_sequence"),
		ActiveStakedrop: collections.NewIndexedMap(sb, types.ActiveStakedropPrefix, "active_airdrop", collections.PairKeyCodec(collections.Uint64Key, collections.Uint64Key), codec.CollValue[types.Stakedrop](cdc),
			StakedropIndexes{
				StakedropByDenom: indexes.NewMulti(sb, types.StakedropIndexKey, "stakedrop_by_denom", collections.StringKey, collections.PairKeyCodec(collections.Uint64Key, collections.Uint64Key), func(pk collections.Pair[uint64, uint64], value types.Stakedrop) (string, error) {
					return value.Amount.Denom, nil
				}),
			}),
	}
	_, err := sb.Build()
	if err != nil {
		panic(err)
	}

	return k
}

func (k Keeper) getNextStakedropSequence(ctx sdk.Context) (uint64, error) {
	seq, err := k.StakedropSequence.Next(ctx)
	return seq, err
}

// GetAuthority returns the x/mint module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a logger for the x/tokenfactory module
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetDenomPrefixStore returns the substore for a specific denom
func (k Keeper) GetDenomPrefixStore(ctx sdk.Context, denom string) storetypes.KVStore {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.GetDenomPrefixStore(denom))
}

// GetCreatorPrefixStore returns the substore for a specific creator address
func (k Keeper) GetCreatorPrefixStore(ctx sdk.Context, creator string) storetypes.KVStore {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.GetCreatorPrefix(creator))
}

// GetCreatorsPrefixStore returns the substore that contains a list of creators
func (k Keeper) GetCreatorsPrefixStore(ctx sdk.Context) storetypes.KVStore {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.GetCreatorsPrefix())
}

// CreateModuleAccount creates a module account with minting and burning capabilities
// This account isn't intended to store any coins,
// it purely mints and burns them on behalf of the admin of respective denoms,
// and sends to the relevant address.
func (k Keeper) CreateModuleAccount(ctx sdk.Context) {
	moduleAcc := authtypes.NewEmptyModuleAccount(types.ModuleName, authtypes.Minter, authtypes.Burner)
	k.accountKeeper.SetModuleAccount(ctx, moduleAcc)
}
