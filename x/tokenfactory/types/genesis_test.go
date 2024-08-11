package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/ChihuahuaChain/chihuahua/x/tokenfactory/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/cosmos1t7egva48prqmzl59x5ngv4zx0dtrwewcdqdjr8/bitcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "cosmos1t7egva48prqmzl59x5ngv4zx0dtrwewcdqdjr8",
						},
					},
				},
				Params: types.Params{
					BuildersCommission: math.LegacyNewDecWithPrec(1, 2),
				},
			},
			valid: true,
		},
		{
			desc: "different admin from creator",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/cosmos1t7egva48prqmzl59x5ngv4zx0dtrwewcdqdjr8/bitcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "cosmos1ft6e5esdtdegnvcr3djd3ftk4kwpcr6jta8eyh",
						},
					},
				},
				Params: types.Params{
					BuildersCommission: math.LegacyNewDecWithPrec(1, 2),
				},
			},
			valid: true,
		},
		{
			desc: "empty admin",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/cosmos1t7egva48prqmzl59x5ngv4zx0dtrwewcdqdjr8/bitcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "",
						},
					},
				},
				Params: types.Params{
					BuildersCommission: math.LegacyNewDecWithPrec(1, 2),
				},
			},
			valid: true,
		},
		{
			desc: "no admin",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/cosmos1t7egva48prqmzl59x5ngv4zx0dtrwewcdqdjr8/bitcoin",
					},
				},
				Params: types.Params{
					BuildersCommission: math.LegacyNewDecWithPrec(1, 2),
				},
			},
			valid: true,
		},
		{
			desc: "invalid admin",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/cosmos1t7egva48prqmzl59x5ngv4zx0dtrwewcdqdjr8/bitcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "moose",
						},
					},
				},
				Params: types.Params{
					BuildersCommission: math.LegacyNewDecWithPrec(1, 2),
				},
			},
			valid: false,
		},
		{
			desc: "multiple denoms",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/cosmos1t7egva48prqmzl59x5ngv4zx0dtrwewcdqdjr8/bitcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "",
						},
					},
					{
						Denom: "factory/cosmos1t7egva48prqmzl59x5ngv4zx0dtrwewcdqdjr8/litecoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "",
						},
					},
				},
				Params: types.Params{
					BuildersCommission: math.LegacyNewDecWithPrec(1, 2),
				},
			},
			valid: true,
		},
		{
			desc: "duplicate denoms",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/cosmos1t7egva48prqmzl59x5ngv4zx0dtrwewcdqdjr8/bitcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "",
						},
					},
					{
						Denom: "factory/cosmos1t7egva48prqmzl59x5ngv4zx0dtrwewcdqdjr8/bitcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "",
						},
					},
				},
				Params: types.Params{
					BuildersCommission: math.LegacyNewDecWithPrec(1, 2),
				},
			},
			valid: false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
