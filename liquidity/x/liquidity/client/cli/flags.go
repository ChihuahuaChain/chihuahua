package cli

// DONTCOVER

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagPoolCoinDenom = "pool-coin-denom"
	FlagReserveAcc    = "reserve-acc"
	FlagDenoms        = "denoms"
)

func flagSetPool() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagPoolCoinDenom, "", "The denomination of the pool coin")
	fs.String(FlagReserveAcc, "", "The Bech32 address of the reserve account")
	fs.String(FlagDenoms, "", "Coin denoms of the pool")

	return fs
}
