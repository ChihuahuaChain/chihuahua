package types

import (
	context "context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	// RemoveAccount and GetModuleAccount let EndBlock materialize the burn
	// address as a module account before burning, see module.go.
	RemoveAccount(ctx context.Context, acc sdk.AccountI)
	GetModuleAccount(ctx context.Context, moduleName string) sdk.ModuleAccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances
// and burn the balance accumulated by the burn address.
type BankKeeper interface {
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	// GetBalance and IterateAccountBalances let EndBlock find the burnable
	// balance of the burn address without a full GetAllBalances scan over every
	// stuck non-burnable denom, see module.go.
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	IterateAccountBalances(ctx context.Context, addr sdk.AccAddress, cb func(sdk.Coin) bool)
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	// Methods imported from bank should be defined here
}
