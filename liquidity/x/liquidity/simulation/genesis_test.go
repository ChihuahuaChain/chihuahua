package simulation_test

import (
	"encoding/json"
	"math/rand"
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/stretchr/testify/require"

	"github.com/Victor118/liquidity/x/liquidity/simulation"
	"github.com/Victor118/liquidity/x/liquidity/types"
)

// TestRandomizedGenState tests the normal scenario of applying RandomizedGenState.
// Abnormal scenarios are not tested here.
func TestRandomizedGenState(t *testing.T) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)
	s := rand.NewSource(1)
	r := rand.New(s)

	simState := module.SimulationState{
		AppParams:    make(simtypes.AppParams),
		Cdc:          cdc,
		Rand:         r,
		NumBonded:    3,
		Accounts:     simtypes.RandomAccounts(r, 3),
		InitialStake: math.NewInt(1000),
		GenState:     make(map[string]json.RawMessage),
	}

	simulation.RandomizedGenState(&simState)

	var liquidityGenesis types.GenesisState
	simState.Cdc.MustUnmarshalJSON(simState.GenState[types.ModuleName], &liquidityGenesis)

	dec1, _ := math.NewIntFromString("4122540")
	dec2, _ := math.NewIntFromString("11240456")
	dec3, _ := math.NewIntFromString("2040480279449")
	dec4, _ := math.LegacyNewDecFromStr("0.448590000000000000")
	dec5, _ := math.LegacyNewDecFromStr("0.732160000000000000")
	dec6, _ := math.LegacyNewDecFromStr("0.237840000000000000")

	require.Equal(t, dec1, liquidityGenesis.Params.MinInitDepositAmount)
	require.Equal(t, dec2, liquidityGenesis.Params.InitPoolCoinMintAmount)
	require.Equal(t, dec3, liquidityGenesis.Params.MaxReserveCoinAmount)
	require.Equal(t, dec4, liquidityGenesis.Params.SwapFeeRate)
	require.Equal(t, dec5, liquidityGenesis.Params.WithdrawFeeRate)
	require.Equal(t, dec6, liquidityGenesis.Params.MaxOrderAmountRatio)
	require.Equal(t, uint32(6), liquidityGenesis.Params.UnitBatchHeight)
}

// TestRandomizedGenState tests abnormal scenarios of applying RandomizedGenState.
func TestRandomizedGenState1(t *testing.T) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	s := rand.NewSource(1)
	r := rand.New(s)

	// all these tests will panic
	tests := []struct {
		simState module.SimulationState
		panicMsg string
	}{
		{ // panic => reason: incomplete initialization of the simState
			module.SimulationState{}, "invalid memory address or nil pointer dereference"},
		{ // panic => reason: incomplete initialization of the simState
			module.SimulationState{
				AppParams: make(simtypes.AppParams),
				Cdc:       cdc,
				Rand:      r,
			}, "assignment to entry in nil map"},
	}

	for _, tt := range tests {
		require.Panicsf(t, func() { simulation.RandomizedGenState(&tt.simState) }, tt.panicMsg)
	}
}
