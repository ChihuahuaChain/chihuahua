package app

import (
	"testing"
	"time"

	"cosmossdk.io/collections"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/store/v2/prefix"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

// TestFeegrantExpiryHaltAndFix is the end-to-end proof that the conversion
// prevents the real chain halt: it drives the actual feegrant EndBlocker
// (RemoveExpiredAllowances) over a queue entry written the SDK v0.50 way (empty
// value). Before the conversion the EndBlocker fails to decode it and returns an
// error (which the module manager turns into a FinalizeBlock failure, i.e. the
// halt); after fixFeegrantQueue runs it decodes and prunes the allowance.
func TestFeegrantExpiryHaltAndFix(t *testing.T) {
	app := Setup(t)
	now := time.Now().UTC()
	ctx := app.BaseApp.NewContext(false).WithBlockHeight(10).WithBlockTime(now)

	accs := simtestutil.CreateIncrementalAccounts(2)
	granter, grantee := accs[0], accs[1]
	// make both accounts exist on chain, as they would in production
	initAccountWithCoins(app, ctx, granter, sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 1_000_000)))
	initAccountWithCoins(app, ctx, grantee, sdk.NewCoins(sdk.NewInt64Coin("uhuahua", 1_000_000)))
	exp := now.Add(time.Minute)
	require.NoError(t, app.FeeGrantKeeper.GrantAllowance(ctx, granter, grantee, &feegrant.BasicAllowance{Expiration: &exp}))

	// Rewrite the expiry-queue entry the SDK v0.50 way: an empty value.
	queue := prefix.NewStore(ctx.KVStore(app.GetKey(feegrant.StoreKey)), feegrant.FeeAllowanceQueueKeyPrefix.Bytes())
	var keys [][]byte
	it := queue.Iterator(nil, nil)
	for ; it.Valid(); it.Next() {
		keys = append(keys, append([]byte(nil), it.Key()...))
	}
	it.Close()
	require.NotEmpty(t, keys, "granting with an expiration must enqueue a queue entry")
	for _, k := range keys {
		queue.Set(k, []byte{})
	}

	// Stand after the expiration so the entry is in the EndBlocker's range.
	future := ctx.WithBlockTime(exp.Add(time.Minute))

	// Without the conversion: the EndBlocker cannot decode the empty value and
	// errors. On a live chain this is the consensus failure that halts it.
	require.Error(t, app.FeeGrantKeeper.RemoveExpiredAllowances(future, 100),
		"an empty v0.50 queue value must make the feegrant EndBlocker fail (the chain halt)")

	// Apply the conversion (clear the marker first: Setup already ran a block,
	// so PreBlocker set it), then the EndBlocker succeeds and prunes the grant.
	future.KVStore(app.GetKey(feeburntypes.StoreKey)).Delete(feegrantQueueFixedKey)
	app.fixFeegrantQueue(future)
	require.NoError(t, app.FeeGrantKeeper.RemoveExpiredAllowances(future, 100),
		"after conversion the feegrant EndBlocker must not fail")
	_, err := app.FeeGrantKeeper.GetAllowance(future, granter, grantee)
	require.Error(t, err, "the expired allowance must have been removed")
}
