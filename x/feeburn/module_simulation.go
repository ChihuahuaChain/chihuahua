package feeburn

import (
	"github.com/ChihuahuaChain/chihuahua/testutil/sample"

	// simappparams "cosmossdk.io/simapp/params"
	feeburnsimulation "github.com/ChihuahuaChain/chihuahua/x/feeburn/simulation"
	"github.com/ChihuahuaChain/chihuahua/x/feeburn/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = feeburnsimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	feeburnGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&feeburnGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized  param changes for the simulator
// func (am AppModule) RandomizedParams(_ *rand.Rand) []simtypes.ParamChange {
// 	feeburnParams := types.DefaultParams()
// 	return []simtypes.ParamChange{
// 		simulation.NewSimParamChange(types.ModuleName, string(types.KeyTxFeeBurnPercent), func(r *rand.Rand) string {
// 			return string(types.Amino.MustMarshalJSON(feeburnParams.TxFeeBurnPercent))
// 		}),
// 	}
// }

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
