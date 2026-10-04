package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	cmtcfg "github.com/cometbft/cometbft/config"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"cosmossdk.io/log/v2"
	confixcmd "cosmossdk.io/tools/confix/cmd"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/pruning"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/snapshot"
	"github.com/cosmos/cosmos-sdk/server"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmcli "github.com/CosmWasm/wasmd/x/wasm/client/cli"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	"github.com/ChihuahuaChain/chihuahua/app"
	"github.com/ChihuahuaChain/chihuahua/app/params"
)

// NewRootCmd creates a new root command for chihuahuad. It is called once in the
// main function.
func NewRootCmd() (*cobra.Command, params.EncodingConfig) {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	cfg.SetBech32PrefixForValidator(app.Bech32PrefixValAddr, app.Bech32PrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(app.Bech32PrefixConsAddr, app.Bech32PrefixConsPub)
	cfg.SetAddressVerifier(wasmtypes.VerifyAddressLen())
	cfg.Seal()

	// we "pre"-instantiate the application for getting the injected/configured encoding configuration
	temp := tempDir()
	// cleanup temp dir after we are done with the tempApp, so we don't leave behind a
	// new temporary directory for every invocation
	defer os.RemoveAll(temp)
	tempApp := app.New(log.NewNopLogger(), dbm.NewMemDB(), nil, false, map[int64]bool{}, temp, 0,
		simtestutil.NewAppOptionsWithFlagHome(temp), []wasmkeeper.Option{})
	encodingConfig := params.EncodingConfig{
		InterfaceRegistry: tempApp.InterfaceRegistry(),
		Codec:             tempApp.AppCodec(),
		TxConfig:          tempApp.TxConfig(),
		Amino:             tempApp.LegacyAmino(),
	}

	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Codec).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithHomeDir(app.DefaultNodeHome).
		WithViper("")

	rootCmd := &cobra.Command{
		Use:           version.AppName,
		Short:         "The Most Loyal Cosmos Dog Coin",
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx = initClientCtx.WithCmdContext(cmd.Context())
			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := initAppConfig()
			customCMTConfig := initCometBFTConfig()

			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, customCMTConfig)
		},
	}

	initRootCmd(rootCmd, encodingConfig.TxConfig, tempApp.BasicModuleManager)

	// add keyring to autocli opts
	autoCliOpts := tempApp.AutoCliOpts()
	initClientCtx, _ = config.ReadFromClientConfig(initClientCtx)
	autoCliOpts.ClientCtx = initClientCtx

	if err := autoCliOpts.EnhanceRootCommand(rootCmd); err != nil {
		panic(err)
	}

	return rootCmd, encodingConfig
}

// initCometBFTConfig helps to override default CometBFT Config values.
// return cmtcfg.DefaultConfig if no custom configuration is required for the application.
func initCometBFTConfig() *cmtcfg.Config {
	return cmtcfg.DefaultConfig()
}

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, any) {
	type CustomAppConfig struct {
		serverconfig.Config

		Wasm wasmtypes.NodeConfig `mapstructure:"wasm"`
	}

	srvCfg := serverconfig.DefaultConfig()
	srvCfg.StateSync.SnapshotInterval = 1500
	srvCfg.StateSync.SnapshotKeepRecent = 2
	srvCfg.BaseConfig.IAVLDisableFastNode = false

	customAppConfig := CustomAppConfig{
		Config: *srvCfg,
		Wasm:   wasmtypes.DefaultNodeConfig(),
	}

	customAppTemplate := serverconfig.DefaultConfigTemplate +
		wasmtypes.DefaultConfigTemplate()

	return customAppTemplate, customAppConfig
}

func initRootCmd(rootCmd *cobra.Command, txConfig client.TxConfig, basicManager module.BasicManager) {
	rootCmd.AddCommand(
		genutilcli.InitCmd(basicManager, app.DefaultNodeHome),
		debug.Cmd(),
		confixcmd.ConfigCommand(),
		pruning.Cmd(newApp, app.DefaultNodeHome),
		snapshot.Cmd(newApp),
	)

	server.AddCommands(rootCmd, app.DefaultNodeHome, newApp, appExport, addModuleInitFlags)
	wasmcli.ExtendUnsafeResetAllCmd(rootCmd)

	// Warn (without blocking) when the node starts with a local minimum-gas-prices
	// below the protocol floor enforced by x/feeburn.
	if startCmd, _, err := rootCmd.Find([]string{"start"}); err == nil && startCmd != nil {
		orig := startCmd.PreRunE
		startCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
			if orig != nil {
				if err := orig(cmd, args); err != nil {
					return err
				}
			}
			warnIfMinGasPricesBelowFloor(cmd)
			return nil
		}
	}

	// add keybase, auxiliary RPC, query, genesis, and tx child commands
	rootCmd.AddCommand(
		server.StatusCommand(),
		genutilcli.Commands(txConfig, basicManager, app.DefaultNodeHome),
		queryCommand(),
		txCommand(),
		keys.Commands(),
		barkCommand(),
	)
}

// barkCommand says hello, the chihuahua way.
func barkCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "bark",
		Aliases: []string{"woof"},
		Short:   "Woof woof",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.Println("woof woof")
			return nil
		},
	}
}

func addModuleInitFlags(startCmd *cobra.Command) {
	wasm.AddModuleInitFlags(startCmd)
}

// warnIfMinGasPricesBelowFloor prints a non-blocking warning when the node's
// local minimum-gas-prices is empty or below the protocol floor
// (app.RecommendedMinGasPrices). The floor is enforced in consensus by
// x/feeburn regardless of this local setting; a lower local value only makes
// the mempool accept txs that are then rejected at block execution. The warning
// pauses a few seconds so operators notice it, then startup continues.
func warnIfMinGasPricesBelowFloor(cmd *cobra.Command) {
	serverCtx := server.GetServerContextFromCmd(cmd)

	floor, err := sdk.ParseDecCoins(app.RecommendedMinGasPrices)
	if err != nil {
		return // RecommendedMinGasPrices is a constant; this cannot happen.
	}

	local, err := sdk.ParseDecCoins(serverCtx.Viper.GetString("minimum-gas-prices"))
	below := err != nil
	if !below {
		for _, f := range floor {
			if local.AmountOf(f.Denom).LT(f.Amount) {
				below = true
				break
			}
		}
	}
	if !below {
		return
	}

	out := cmd.ErrOrStderr()
	fmt.Fprintf(out, "\n⚠️  minimum-gas-prices is below the chain floor (%s).\n", app.RecommendedMinGasPrices)
	fmt.Fprintf(out, "    The protocol enforces this minimum in consensus: transactions under it are\n")
	fmt.Fprintf(out, "    rejected when the block executes, on every node. A lower local value only\n")
	fmt.Fprintf(out, "    makes this node's mempool accept txs that will then fail. Set\n")
	fmt.Fprintf(out, "    minimum-gas-prices to at least %s in app.toml (or --minimum-gas-prices).\n", app.RecommendedMinGasPrices)
	fmt.Fprintf(out, "    Starting anyway in 5s...\n\n")
	time.Sleep(5 * time.Second)
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		rpc.QueryEventForTxCmd(),
		server.QueryBlockCmd(),
		authcmd.QueryTxsByEventsCmd(),
		server.QueryBlocksCmd(),
		authcmd.QueryTxCmd(),
		server.QueryBlockResultsCmd(),
	)

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetMultiSignBatchCmd(),
		authcmd.GetValidateSignaturesCommand(),
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
		authcmd.GetSimulateCmd(),
	)

	return cmd
}

// newApp creates the application
func newApp(logger log.Logger, db dbm.DB, appOpts servertypes.AppOptions) servertypes.Application {
	baseappOptions := server.DefaultBaseappOptions(appOpts)

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	return app.New(
		logger, db, nil, true, skipUpgradeHeights,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		0,
		appOpts,
		app.GetWasmOpts(appOpts),
		baseappOptions...,
	)
}

// appExport creates a new app (optionally at a given height) and exports state.
func appExport(
	logger log.Logger,
	db dbm.DB,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	// this check is necessary as we use the flag in x/upgrade.
	// we can exit more gracefully by checking the flag here.
	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home is not set")
	}

	if _, ok := appOpts.(*viper.Viper); !ok {
		return servertypes.ExportedApp{}, errors.New("appOpts is not viper.Viper")
	}

	chihuahuaApp := app.New(logger, db, nil, height == -1, map[int64]bool{}, homePath, 0, appOpts, app.GetWasmOpts(appOpts))

	if height != -1 {
		if err := chihuahuaApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	}

	return chihuahuaApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}

var tempDir = func() string {
	dir, err := os.MkdirTemp("", "chihuahuad")
	if err != nil {
		panic("failed to create temp dir: " + err.Error())
	}

	return dir
}
