package keeper

import (
	"testing"

	"cosmossdk.io/log/v2"
	"github.com/ChihuahuaChain/chihuahua/x/feeburn/keeper"
	"github.com/ChihuahuaChain/chihuahua/x/feeburn/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	db2 "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/v2"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"
)

func FeeburnKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	db := db2.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		authtypes.NewModuleAddress(govtypes.ModuleName),
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	_ = k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}
