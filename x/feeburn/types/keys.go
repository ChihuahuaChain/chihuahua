package types

const (
	// ModuleName defines the module name
	ModuleName = "feeburn"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_feeburn"
)

var (
	ParamsKey            = []byte{0x00} // Prefix for params key
	TotalBurnedKeyPrefix = []byte{0x01} // Prefix for the burned fees, one key per denom
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
