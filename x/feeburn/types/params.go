package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyTxFeeBurnPercent = []byte("TxFeeBurnPercent")
	// TODO: Determine the default value
	DefaultTxFeeBurnPercent string = "tx_fee_burn_percent"
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	txFeeBurnPercent string,
) Params {
	return Params{
		TxFeeBurnPercent: txFeeBurnPercent,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultTxFeeBurnPercent,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyTxFeeBurnPercent, &p.TxFeeBurnPercent, validateTxFeeBurnPercent),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateTxFeeBurnPercent(p.TxFeeBurnPercent); err != nil {
		return err
	}

	return nil
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

	// TODO implement validation
	_ = txFeeBurnPercent

	return nil
}
