package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// EnsureBurnModuleAccount makes the burn address a module account with the
// module's permissions. If a plain account squats the address (someone sent to
// it while it was an ordinary, unblocked address), its record is dropped first:
// GetModuleAccount panics on a non-module account, and x/bank keeps the balance
// by address, so the balance survives the record being replaced.
//
// It is the single source of truth for materializing the burn account, called
// both by the v10.0.1 upgrade handler (proactively, so the account exists from
// the first post-upgrade block) and by EndBlock (self-healing, so a chain that
// never ran the handler cannot halt on the first burn).
func EnsureBurnModuleAccount(ctx sdk.Context, ak AccountKeeper) {
	addr := authtypes.NewModuleAddress(ModuleName)
	if acc := ak.GetAccount(ctx, addr); acc != nil {
		if _, ok := acc.(sdk.ModuleAccountI); ok {
			return // already a module account, the common case: nothing to do
		}
		ak.RemoveAccount(ctx, acc)
	}
	ak.GetModuleAccount(ctx, ModuleName)
}
