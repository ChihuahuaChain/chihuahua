package types

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AlphabeticalDenomPair returns denom pairs that are alphabetically sorted.
func AlphabeticalDenomPair(denom1, denom2 string) (resDenom1, resDenom2 string) {
	if denom1 > denom2 {
		return denom2, denom1
	}
	return denom1, denom2
}

// SortDenoms sorts denoms in alphabetical order.
func SortDenoms(denoms []string) []string {
	sort.Strings(denoms)
	return denoms
}

// GetPoolReserveAcc returns the module address of the pool's reserve account.
func GetPoolReserveAcc(poolName string) sdk.AccAddress {

	poolCoinDenom := GetPoolCoinDenom(poolName)
	poolCoinDenom = strings.TrimPrefix(poolCoinDenom, PoolCoinDenomPrefix)
	return sdk.AccAddress(address.Module(ModuleName, []byte(poolCoinDenom)))

}

// CreateModuleAccount creates a module account at the provided address.
// It overrides an account if it exists at that address, with a non-zero sequence number & pubkey
// Contract: addr is derived from `address.Module(ModuleName, key)`
func CreateModuleAccount(ctx sdk.Context, ak AccountKeeper, addr sdk.AccAddress) error {

	acc := ak.NewAccount(
		ctx,
		authtypes.NewModuleAccount(
			authtypes.NewBaseAccountWithAddress(addr),
			addr.String(),
		),
	)
	ak.SetAccount(ctx, acc)
	return nil
}

// GetPoolCoinDenom returns the denomination of the pool coin.
func GetPoolCoinDenom(poolName string) string {
	// Originally pool coin denom has prefix with / splitter, but removed prefix for pass validation of ibc-transfer
	return fmt.Sprintf("%s%X", PoolCoinDenomPrefix, sha256.Sum256([]byte(poolName)))
}

// GetReserveAcc extracts and returns reserve account from pool coin denom.
func GetReserveAcc(poolCoinDenom string) (sdk.AccAddress, error) {

	if err := sdk.ValidateDenom(poolCoinDenom); err != nil {
		return nil, err
	}
	if !strings.HasPrefix(poolCoinDenom, PoolCoinDenomPrefix) {
		return nil, ErrInvalidDenom
	}
	poolCoinDenom = strings.TrimPrefix(poolCoinDenom, PoolCoinDenomPrefix)
	if len(poolCoinDenom) != 64 {
		return nil, ErrInvalidDenom
	}

	return sdk.AccAddress(address.Module(ModuleName, []byte(poolCoinDenom))), nil

}

// GetCoinsTotalAmount returns total amount of all coins in sdk.Coins.
func GetCoinsTotalAmount(coins sdk.Coins) math.Int {
	totalAmount := math.ZeroInt()
	for _, coin := range coins {
		totalAmount = totalAmount.Add(coin.Amount)
	}
	return totalAmount
}

// ValidateReserveCoinLimit checks if total amounts of depositCoins exceed maxReserveCoinAmount.
func ValidateReserveCoinLimit(maxReserveCoinAmount math.Int, depositCoins sdk.Coins) error {
	totalAmount := GetCoinsTotalAmount(depositCoins)
	if maxReserveCoinAmount.IsZero() {
		return nil
	} else if totalAmount.GT(maxReserveCoinAmount) {
		return ErrExceededReserveCoinLimit
	} else {
		return nil
	}
}

func GetOfferCoinFee(offerCoin sdk.Coin, swapFeeRate math.LegacyDec) sdk.Coin {
	if swapFeeRate.IsZero() {
		return sdk.NewCoin(offerCoin.Denom, math.ZeroInt())
	}
	// apply half-ratio swap fee rate and ceiling
	// see https://github.com/tendermint/liquidity/issues/41 for details
	return sdk.NewCoin(offerCoin.Denom, offerCoin.Amount.ToLegacyDec().Mul(swapFeeRate.QuoInt64(2)).Ceil().TruncateInt()) // Ceil(offerCoin.Amount * (swapFeeRate/2))
}

func MustParseCoinsNormalized(coinStr string) sdk.Coins {
	coins, err := sdk.ParseCoinsNormalized(coinStr)
	if err != nil {
		panic(err)
	}
	return coins
}

func CheckOverflow(a, b math.Int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrOverflowAmount
		}
	}()
	a.Mul(b)
	a.Quo(b)
	b.Quo(a)
	return nil
}

func CheckOverflowWithDec(a, b math.LegacyDec) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrOverflowAmount
		}
	}()
	a.Mul(b)
	a.Quo(b)
	b.Quo(a)
	return nil
}
