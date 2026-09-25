package app

import (
	chihuahuaapp "github.com/ChihuahuaChain/chihuahua/app"
)

// LiquidityApp is the chihuahua app: the liquidity tests run against the
// actual chain wiring.
type LiquidityApp = chihuahuaapp.App

// GenesisState is the chihuahua app genesis state.
type GenesisState = chihuahuaapp.GenesisState

func init() {
	// set the chihuahua bech32 prefixes before any test converts an address:
	// the sdk caches address strings, so a later change would not apply to them
	chihuahuaapp.SetTestAddressPrefixes()
}
