package app

import (
	"github.com/cosmos/cosmos-sdk/x/feegrant"

	"github.com/cosmos/cosmos-sdk/store/v2/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	feeburnmoduletypes "github.com/ChihuahuaChain/chihuahua/x/feeburn/types"
)

// feegrantQueueFixedKey marks, in the feeburn store, that the feegrant expiry queue has been converted.
var feegrantQueueFixedKey = []byte{0x7f, 'f', 'g', 'q'}

// fixFeegrantQueue converts the fee allowance expiry queue written by SDK v0.50 (v9.x), whose values are empty, to
// the v0.54 encoding, a one-byte bool. v0.54 decodes those values with collections.BoolValue, which refuses an empty
// value: the first allowance created before v10 that expires makes the feegrant EndBlocker fail and halts the chain
// (woofnet-6 at height 432248). The conversion runs once, before anything reads the queue, and leaves a marker so
// later blocks only pay for one store read; queue entries written by v0.54 are already encoded.
func (app *App) fixFeegrantQueue(ctx sdk.Context) {
	marker := ctx.KVStore(app.keys[feeburnmoduletypes.StoreKey])
	if marker.Has(feegrantQueueFixedKey) {
		return
	}
	queue := prefix.NewStore(ctx.KVStore(app.keys[feegrant.StoreKey]), feegrant.FeeAllowanceQueueKeyPrefix.Bytes())
	var empty [][]byte
	it := queue.Iterator(nil, nil)
	for ; it.Valid(); it.Next() {
		if len(it.Value()) == 0 {
			empty = append(empty, append([]byte(nil), it.Key()...))
		}
	}
	it.Close()
	for _, k := range empty {
		queue.Set(k, []byte{1})
	}
	marker.Set(feegrantQueueFixedKey, []byte{1})
	if len(empty) > 0 {
		app.Logger().Info("feegrant: converted the expiry queue to the v0.54 encoding", "entries", len(empty))
	}
}
