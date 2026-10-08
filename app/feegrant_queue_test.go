package app

import (
	"testing"

	"cosmossdk.io/collections"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/store/v2/prefix"
	"github.com/cosmos/cosmos-sdk/x/feegrant"

	feeburntypes "github.com/ChihuahuaChain/chihuahua/x/feeburn/types"
)

// TestFixFeegrantQueue reproduces the chain halt the v10.0.1 upgrade fixes:
// SDK v0.50 (v9.x) wrote the fee allowance expiry queue with empty values,
// while v0.54 decodes those values with collections.BoolValue, which rejects
// an empty value. The first pre-v10 allowance to expire then made the feegrant
// EndBlocker fail and halted the chain (woofnet-6 at height 432248).
// fixFeegrantQueue converts the empty values to the v0.54 one-byte bool once,
// before anything reads the queue.
func TestFixFeegrantQueue(t *testing.T) {
	app := Setup(t)
	ctx := app.BaseApp.NewContext(false)

	queue := prefix.NewStore(ctx.KVStore(app.GetKey(feegrant.StoreKey)), feegrant.FeeAllowanceQueueKeyPrefix.Bytes())
	legacyKey := []byte("legacy-v9-entry")  // written by SDK v0.50: empty value
	v054Key := []byte("already-v054-entry") // written by SDK v0.54: one-byte bool
	queue.Set(legacyKey, []byte{})
	queue.Set(v054Key, []byte{1})

	// the empty value is exactly what BoolValue cannot decode: this is the halt.
	_, err := collections.BoolValue.Decode([]byte{})
	require.Error(t, err)

	// Setup() already ran a block, so PreBlocker set the marker; clear it to
	// stand at the pre-upgrade state, with a legacy empty-valued queue entry.
	ctx.KVStore(app.GetKey(feeburntypes.StoreKey)).Delete(feegrantQueueFixedKey)
	require.False(t, ctx.KVStore(app.GetKey(feeburntypes.StoreKey)).Has(feegrantQueueFixedKey), "marker must be unset before the fix runs")

	app.fixFeegrantQueue(ctx)

	// the legacy entry is now a decodable one-byte bool; the v0.54 entry is left as is.
	require.Equal(t, []byte{1}, queue.Get(legacyKey))
	require.Equal(t, []byte{1}, queue.Get(v054Key))
	got, err := collections.BoolValue.Decode(queue.Get(legacyKey))
	require.NoError(t, err)
	require.True(t, got)

	require.True(t, ctx.KVStore(app.GetKey(feeburntypes.StoreKey)).Has(feegrantQueueFixedKey), "marker must be set after the fix runs")

	// idempotent: once the marker is set, a later empty value is NOT reconverted,
	// so post-upgrade blocks only pay for one store read.
	queue.Set(legacyKey, []byte{})
	app.fixFeegrantQueue(ctx)
	require.Equal(t, []byte{}, queue.Get(legacyKey), "the fix must be a no-op after the marker is set")
}
