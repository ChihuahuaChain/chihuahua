package types

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types"
)

type WasmMsg struct {
	/// Contracts can create denoms, namespaced under the contract's address.
	/// A contract may create any number of independent sub-denoms.
	CreateDenom *CreateDenom `json:"create_denom,omitempty"`
	/// Contracts can change the admin of a denom that they are the admin of.
	ChangeAdmin *ChangeAdmin `json:"change_admin,omitempty"`
	/// Contracts can mint native tokens for an existing factory denom
	/// that they are the admin of.
	MintTokens *MintTokens `json:"mint_tokens,omitempty"`
	/// Contracts can burn native tokens for an existing factory denom
	/// that they are the admin of.
	/// Currently, the burn from address must be the admin contract.
	BurnTokens *BurnTokens `json:"burn_tokens,omitempty"`
	/// Sets the metadata on a denom which the contract controls
	SetMetadata *SetMetadata `json:"set_metadata,omitempty"`
	/// Forces a transfer of tokens from one address to another.
	ForceTransfer *ForceTransfer `json:"force_transfer,omitempty"`

	CreateStakedrop *CreateStakedrop `json:"create_stakedrop,omitempty"`

	CreatePool *CreatePool `json:"create_pool,omitempty"`

	DirectSwap *DirectSwap `json:"direct_swap,omitempty"`
}

type CreateStakedrop struct {
	Denom      string   `json:"denom"`
	Amount     math.Int `json:"amount"`
	StartBlock int64    `json:"start_block"`
	EndBlock   int64    `json:"end_block"`
}

// CreateDenom creates a new factory denom, of denomination:
// factory/{creating contract address}/{Subdenom}
// Subdenom can be of length at most 44 characters, in [0-9a-zA-Z./]
// The (creating contract address, subdenom) pair must be unique.
// The created denom's admin is the creating contract address,
// but this admin can be changed using the ChangeAdmin binding.
type CreateDenom struct {
	Subdenom string    `json:"subdenom"`
	Metadata *Metadata `json:"metadata,omitempty"`
}

// ChangeAdmin changes the admin for a factory denom.
// If the NewAdminAddress is empty, the denom has no admin.
type ChangeAdmin struct {
	Denom           string `json:"denom"`
	NewAdminAddress string `json:"new_admin_address"`
}

type MintTokens struct {
	Denom         string   `json:"denom"`
	Amount        math.Int `json:"amount"`
	MintToAddress string   `json:"mint_to_address"`
}

type BurnTokens struct {
	Denom           string   `json:"denom"`
	Amount          math.Int `json:"amount"`
	BurnFromAddress string   `json:"burn_from_address"`
}

type SetMetadata struct {
	Denom    string   `json:"denom"`
	Metadata Metadata `json:"metadata"`
}

type ForceTransfer struct {
	Denom       string   `json:"denom"`
	Amount      math.Int `json:"amount"`
	FromAddress string   `json:"from_address"`
	ToAddress   string   `json:"to_address"`
}

// Message for liquidity module
// TODO binding folder should be extracted from tokenfactory module
type CreatePool struct {
	PoolCreatorAddress string   `json:"denom"`
	PoolTypeId         uint32   `json:"pool_type_id"`
	Amount1            math.Int `json:"amount1"`
	Denom1             string   `json:"denom1"`
	Amount2            math.Int `json:"amount2"`
	Denom2             string   `json:"denom2"`
}

type DirectSwap struct {
	SwapRequesterAddress string         `json:"swap_requester_address"`
	PoolId               uint64         `json:"pool_id"`
	SwapTypeId           uint32         `json:"swap_type_id"`
	OfferCoin            types.Coin     `json:"offer_coin"`
	DemandCoinDenom      string         `json:"demand_coin_denom"`
	OfferCoinFee         types.Coin     `json:"offer_coin_fee"`
	OrderPrice           math.LegacyDec `json:"order_price"`
}
