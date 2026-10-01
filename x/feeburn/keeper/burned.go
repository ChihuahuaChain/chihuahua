package keeper

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/v2/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ChihuahuaChain/chihuahua/x/feeburn/types"
)

func (k Keeper) burnedStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), types.TotalBurnedKeyPrefix)
}

// GetTotalBurned returns the transaction fees burned so far.
func (k Keeper) GetTotalBurned(ctx sdk.Context) sdk.Coins {
	iter := k.burnedStore(ctx).Iterator(nil, nil)
	defer iter.Close()

	total := sdk.NewCoins()
	for ; iter.Valid(); iter.Next() {
		var amount math.Int
		if err := amount.Unmarshal(iter.Value()); err != nil {
			panic(err)
		}
		total = total.Add(sdk.NewCoin(string(iter.Key()), amount))
	}
	return total
}

// SetTotalBurned replaces the burned fees total.
func (k Keeper) SetTotalBurned(ctx sdk.Context, total sdk.Coins) error {
	if err := total.Validate(); err != nil {
		return err
	}

	store := k.burnedStore(ctx)
	iter := store.Iterator(nil, nil)
	var denoms [][]byte
	for ; iter.Valid(); iter.Next() {
		denoms = append(denoms, iter.Key())
	}
	iter.Close()
	for _, denom := range denoms {
		store.Delete(denom)
	}

	for _, coin := range total {
		k.setBurned(store, coin.Denom, coin.Amount)
	}
	return nil
}

// AddBurned adds freshly burned fees to the total. It does not consume gas, so
// recording a burn doesn't change what a transaction costs.
func (k Keeper) AddBurned(ctx sdk.Context, burned sdk.Coins) {
	store := k.burnedStore(ctx.WithGasMeter(storetypes.NewInfiniteGasMeter()))
	for _, coin := range burned {
		if !coin.IsPositive() {
			continue
		}
		total := coin.Amount
		if bz := store.Get([]byte(coin.Denom)); bz != nil {
			var prev math.Int
			if err := prev.Unmarshal(bz); err != nil {
				panic(err)
			}
			total = total.Add(prev)
		}
		k.setBurned(store, coin.Denom, total)
	}
}

func (k Keeper) setBurned(store prefix.Store, denom string, amount math.Int) {
	bz, err := amount.Marshal()
	if err != nil {
		panic(err)
	}
	store.Set([]byte(denom), bz)
}
