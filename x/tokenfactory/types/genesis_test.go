package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/ChihuahuaChain/chihuahua/x/tokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func TestGenesisState_ValidateStakedrops(t *testing.T) {
	drop := func(start, end int64) types.Stakedrop {
		return types.Stakedrop{
			Amount:         sdk.NewInt64Coin("uhuahua", 1_000),
			AmountPerBlock: sdk.NewInt64Coin("uhuahua", 100),
			StartBlock:     start,
			EndBlock:       end,
		}
	}
	for _, tc := range []struct {
		desc  string
		drops []types.GenesisStakedrop
		seq   uint64
		valid bool
	}{
		{"valid", []types.GenesisStakedrop{{Sequence: 0, Stakedrop: drop(10, 20)}, {Sequence: 1, Stakedrop: drop(10, 30)}}, 2, true},
		{"duplicate key", []types.GenesisStakedrop{{Sequence: 0, Stakedrop: drop(10, 20)}, {Sequence: 0, Stakedrop: drop(10, 30)}}, 2, false},
		{"sequence not below next", []types.GenesisStakedrop{{Sequence: 2, Stakedrop: drop(10, 20)}}, 2, false},
		{"negative end", []types.GenesisStakedrop{{Sequence: 0, Stakedrop: drop(10, -1)}}, 1, false},
		{"end before start", []types.GenesisStakedrop{{Sequence: 0, Stakedrop: drop(20, 10)}}, 1, false},
		{"denom mismatch", []types.GenesisStakedrop{{Sequence: 0, Stakedrop: types.Stakedrop{
			Amount: sdk.NewInt64Coin("uhuahua", 1_000), AmountPerBlock: sdk.NewInt64Coin("other", 100), StartBlock: 10, EndBlock: 20,
		}}}, 1, false},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			gs := types.DefaultGenesis()
			gs.ActiveStakedrops = tc.drops
			gs.StakedropSequence = tc.seq
			err := gs.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
