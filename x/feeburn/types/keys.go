package types

import "strings"

const (
	// ModuleName defines the module name
	ModuleName = "feeburn"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_feeburn"

	// BurnDenom is the chain's native denom, always destroyed by the burn
	// address in EndBlock.
	BurnDenom = "uhuahua"

	// FactoryDenomPrefix is the prefix of every x/tokenfactory denom (mirrors
	// tokenfactory's ModuleDenomPrefix + "/"). The burn address also destroys
	// these: they are native, locally minted tokens, so burning them reduces
	// real supply. IBC vouchers ("ibc/...") and any other denom are left
	// untouched, because burning an IBC voucher would strand the escrowed
	// asset on the source chain.
	FactoryDenomPrefix = "factory/"
)

// IsBurnable reports whether the burn address should destroy denom: the native
// denom and any token factory denom, never IBC vouchers or anything else.
func IsBurnable(denom string) bool {
	return denom == BurnDenom || strings.HasPrefix(denom, FactoryDenomPrefix)
}

var (
	ParamsKey            = []byte{0x00} // Prefix for params key
	TotalBurnedKeyPrefix = []byte{0x01} // Prefix for the burned fees, one key per denom
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
