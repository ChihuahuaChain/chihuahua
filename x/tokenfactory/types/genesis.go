package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// this line is used by starport scaffolding # genesis/types/import

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:           DefaultParams(),
		FactoryDenoms:    []GenesisDenom{},
		ActiveStakedrops: []GenesisStakedrop{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	err := gs.Params.Validate()
	if err != nil {
		return err
	}

	seenDenoms := map[string]bool{}

	for _, denom := range gs.GetFactoryDenoms() {
		if seenDenoms[denom.GetDenom()] {
			return errorsmod.Wrapf(ErrInvalidGenesis, "duplicate denom: %s", denom.GetDenom())
		}
		seenDenoms[denom.GetDenom()] = true

		_, _, err := DeconstructDenom(denom.GetDenom())
		if err != nil {
			return err
		}

		if denom.AuthorityMetadata.Admin != "" {
			_, err = sdk.AccAddressFromBech32(denom.AuthorityMetadata.Admin)
			if err != nil {
				return errorsmod.Wrapf(ErrInvalidAuthorityMetadata, "Invalid admin address (%s)", err)
			}
		}
	}

	seenStakedrops := map[string]bool{}
	for _, genDrop := range gs.GetActiveStakedrops() {
		drop := genDrop.Stakedrop
		key := fmt.Sprintf("%d/%d", drop.StartBlock, genDrop.Sequence)
		if seenStakedrops[key] {
			return errorsmod.Wrapf(ErrInvalidGenesis, "duplicate stakedrop: start block %d, sequence %d", drop.StartBlock, genDrop.Sequence)
		}
		seenStakedrops[key] = true

		if drop.StartBlock <= 0 || drop.EndBlock <= drop.StartBlock {
			return errorsmod.Wrapf(ErrInvalidGenesis, "stakedrop %s: invalid blocks %d..%d", key, drop.StartBlock, drop.EndBlock)
		}
		if err := drop.Amount.Validate(); err != nil {
			return errorsmod.Wrapf(ErrInvalidGenesis, "stakedrop %s: %s", key, err)
		}
		if err := drop.AmountPerBlock.Validate(); err != nil {
			return errorsmod.Wrapf(ErrInvalidGenesis, "stakedrop %s: %s", key, err)
		}
		if drop.AmountPerBlock.Denom != drop.Amount.Denom {
			return errorsmod.Wrapf(ErrInvalidGenesis, "stakedrop %s: denom mismatch", key)
		}
		if genDrop.Sequence >= gs.StakedropSequence {
			return errorsmod.Wrapf(ErrInvalidGenesis, "stakedrop %s: sequence must be below stakedrop_sequence %d", key, gs.StakedropSequence)
		}
	}

	return nil
}
