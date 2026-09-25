package app

import (
	"crypto/sha256"
	"fmt"
	"testing"

	transferkeeper "github.com/cosmos/ibc-go/v11/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v11/modules/apps/transfer/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protowire"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// the only denom trace of chihuahua-1 that ibc-go cannot migrate
var ambiguousTrace = [2]string{"transfer/channel-208", "IRO/heilelonmusk_653667-1"}

func setDenomTrace(t *testing.T, app *App, ctx sdk.Context, path, base string) []byte {
	t.Helper()
	var bz []byte
	bz = protowire.AppendTag(bz, 1, protowire.BytesType)
	bz = protowire.AppendString(bz, path)
	bz = protowire.AppendTag(bz, 2, protowire.BytesType)
	bz = protowire.AppendString(bz, base)
	hash := sha256.Sum256([]byte(path + "/" + base))
	store := runtime.KVStoreAdapter(runtime.NewKVStoreService(app.GetKey(ibctransfertypes.StoreKey)).OpenKVStore(ctx))
	store.Set(append(append([]byte{}, ibctransfertypes.DenomTraceKey...), hash[:]...), bz)
	return hash[:]
}

func TestTransferMigrationPanicsOnAmbiguousTrace(t *testing.T) {
	app := Setup(t)
	ctx := app.BaseApp.NewContext(false)
	setDenomTrace(t, app, ctx, ambiguousTrace[0], ambiguousTrace[1])

	require.Panics(t, func() {
		_ = transferkeeper.NewMigrator(*app.TransferKeeper).MigrateDenomTraceToDenom(ctx)
	})
}

func TestTakeAmbiguousDenomTraces(t *testing.T) {
	app := Setup(t)
	ctx := app.BaseApp.NewContext(false)
	normalHash := setDenomTrace(t, app, ctx, "transfer/channel-1", "uatom")
	multiHopHash := setDenomTrace(t, app, ctx, "transfer/channel-7/transfer/channel-143", "erc20/tether/usdt")
	ambiguousHash := setDenomTrace(t, app, ctx, ambiguousTrace[0], ambiguousTrace[1])

	denoms, err := app.takeAmbiguousDenomTraces(ctx)
	require.NoError(t, err)
	require.Len(t, denoms, 1)

	require.NoError(t, transferkeeper.NewMigrator(*app.TransferKeeper).MigrateDenomTraceToDenom(ctx))
	for _, d := range denoms {
		app.TransferKeeper.SetDenom(ctx, d)
	}

	for _, tc := range []struct {
		hash []byte
		path string
		base string
	}{
		{normalHash, "transfer/channel-1/uatom", "uatom"},
		{multiHopHash, "transfer/channel-7/transfer/channel-143/erc20/tether/usdt", "erc20/tether/usdt"},
		{ambiguousHash, "transfer/channel-208/IRO/heilelonmusk_653667-1", "IRO/heilelonmusk_653667-1"},
	} {
		denom, found := app.TransferKeeper.GetDenom(ctx, tc.hash)
		require.True(t, found, tc.path)
		require.Equal(t, tc.path, denom.Path())
		require.Equal(t, tc.base, denom.Base)
		require.Equal(t, fmt.Sprintf("ibc/%X", tc.hash), denom.IBCDenom())
	}
}
