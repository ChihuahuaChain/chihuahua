package main

import (
	"github.com/ChihuahuaChain/chihuahua/app"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/cosmoscmd"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmcli "github.com/CosmWasm/wasmd/x/wasm/client/cli"
)

func AddGenesisWasmMsgCmd(defaultNodeHome string) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        "add-wasm-genesis-message",
		Short:                      "Wasm genesis subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	genesisIO := wasmcli.NewDefaultGenesisIO()
	txCmd.AddCommand(
		wasmcli.GenesisStoreCodeCmd(defaultNodeHome, genesisIO),
		wasmcli.GenesisInstantiateContractCmd(defaultNodeHome, genesisIO),
		wasmcli.GenesisExecuteContractCmd(defaultNodeHome, genesisIO),
		wasmcli.GenesisListContractsCmd(defaultNodeHome, genesisIO),
		wasmcli.GenesisListCodesCmd(defaultNodeHome, genesisIO),
	)

	return txCmd
}

func GetWasmCmdOptions() []cosmoscmd.Option {
	var options []cosmoscmd.Option

	options = append(options,
		cosmoscmd.CustomizeStartCmd(func(startCmd *cobra.Command) {
			wasm.AddModuleInitFlags(startCmd)
		}),
		cosmoscmd.AddSubCmd(AddGenesisWasmMsgCmd(app.DefaultNodeHome)),
	)

	return options
}
