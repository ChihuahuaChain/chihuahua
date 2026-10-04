package types

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"gopkg.in/yaml.v2"
)

var (
	KeyTxFeeBurnPercent = []byte("TxFeeBurnPercent")
	// TODO: Determine the default value
	DefaultTxFeeBurnPercent = "0"

	// DefaultMinGasPrices is nil on purpose: the code default imposes no floor,
	// so fresh genesis, gentx processing, simulations and existing tests keep
	// working. nil (not an empty slice) matches what a stored, unset param
	// decodes to. The production floor is set on mainnet by the v10.0.1 upgrade
	// and thereafter tuned by governance via MsgUpdateParams.
	DefaultMinGasPrices sdk.DecCoins = nil
)

// NewParams creates a new Params instance
func NewParams(
	txFeeBurnPercent string,
	minGasPrices sdk.DecCoins,
) Params {
	return Params{
		TxFeeBurnPercent: txFeeBurnPercent,
		MinGasPrices:     minGasPrices,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultTxFeeBurnPercent,
		DefaultMinGasPrices,
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateTxFeeBurnPercent(p.TxFeeBurnPercent); err != nil {
		return err
	}
	return validateMinGasPrices(p.MinGasPrices)
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// validateTxFeeBurnPercent validates the TxFeeBurnPercent param
func validateTxFeeBurnPercent(v interface{}) error {
	txFeeBurnPercent, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	txFeeBurnPercentInt, err := strconv.Atoi(txFeeBurnPercent)
	if err != nil {
		return err
	}
	if txFeeBurnPercentInt < 0 || txFeeBurnPercentInt > 100 {
		return fmt.Errorf("fee must be between 0 and 100")
	}

	return nil
}

// validateMinGasPrices validates the MinGasPrices param. An empty/nil set is
// valid and means "no floor". sdk.DecCoins.Validate covers valid denoms,
// non-negative amounts, sorted order and no duplicates. On top of that, the
// protocol floor may only be denominated in the chain's native denom
// (BurnDenom, "uhuahua"): governance must not set a floor in any other denom.
func validateMinGasPrices(v interface{}) error {
	coins, ok := v.(sdk.DecCoins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}
	if err := coins.Validate(); err != nil {
		return fmt.Errorf("invalid min_gas_prices: %w", err)
	}
	for _, c := range coins {
		if c.Denom != BurnDenom {
			return fmt.Errorf("invalid min_gas_prices: only %s is allowed, got %s", BurnDenom, c.Denom)
		}
	}
	return nil
}
