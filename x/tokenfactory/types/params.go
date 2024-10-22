package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewParams(denomCreationFee sdk.Coins) Params {
	return Params{
		DenomCreationFee: denomCreationFee,
	}
}

// default tokenfactory module parameters.
func DefaultParams() Params {
	return Params{
		DenomCreationFee:           sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10_000_000)),
		DenomCreationGasConsume:    2_000_000,
		BuildersCommission:         math.LegacyNewDecWithPrec(5, 3), // "0.005" if there is builders addresses, this commission rate from minted amount is redirected to builders
		BuildersAddresses:          []WeightedAddress(nil),
		FreeMintWhitelistAddresses: []string(nil),
		StakedropChargePerBlock:    sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(0)),
	}
}

// validate params.
func (p Params) Validate() error {
	err := validateDenomCreationFee(p.DenomCreationFee)
	if err != nil {
		return err
	}
	err = validateDenomCreationFeeGasConsume(p.DenomCreationGasConsume)
	if err != nil {
		return err
	}
	err = validateBuildersCommission(p.BuildersCommission)
	if err != nil {
		return err
	}
	err = validateBuildersAddresses(p.BuildersAddresses)
	if err != nil {
		return err
	}
	err = validateFreeMintWhitelistAddresses(p.FreeMintWhitelistAddresses)
	if err != nil {
		return err
	}
	err = validateStakedropChargePerBlock(p.StakedropChargePerBlock)
	if err != nil {
		return err
	}
	return nil

}

func validateStakedropChargePerBlock(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.Validate() != nil {
		return fmt.Errorf("invalid denom stakedrop creation fee: %+v", i)
	}

	return nil
}

func validateFreeMintWhitelistAddresses(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// whitelist can be empty
	if len(v) == 0 {
		return nil
	}
	for i, addr := range v {
		_, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			return fmt.Errorf("invalid address at %dth", i)
		}

	}
	return nil
}

func validateBuildersAddresses(i interface{}) error {
	v, ok := i.([]WeightedAddress)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// if builders addresses is empty, there's no commission
	if len(v) == 0 {
		return nil
	}

	weightSum := math.LegacyZeroDec()
	for i, w := range v {
		_, err := sdk.AccAddressFromBech32(w.Address)
		if err != nil {
			return fmt.Errorf("invalid address at %dth", i)
		}

		if !w.Weight.IsPositive() {
			return fmt.Errorf("non-positive weight at %dth", i)
		}
		if w.Weight.GT(math.LegacyNewDec(1)) {
			return fmt.Errorf("more than 1 weight at %dth", i)
		}
		weightSum = weightSum.Add(w.Weight)
	}

	if !weightSum.Equal(math.LegacyOneDec()) {
		return fmt.Errorf("invalid weight sum: %s", weightSum.String())
	}

	return nil
}

func validateDenomCreationFee(i interface{}) error {
	v, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Validate() != nil {
		return fmt.Errorf("invalid denom creation fee: %+v", i)
	}

	return nil
}

func validateDenomCreationFeeGasConsume(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateBuildersCommission(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("builders commission must not be nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("builders commission must not be negative: %s", v)
	}

	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("builders commission too large: %s", v)
	}

	return nil
}
