package testutil

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govcli "github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	paramscli "github.com/cosmos/cosmos-sdk/x/params/client/cli"

	"github.com/Victor118/liquidity/app"
	liquidityapp "github.com/Victor118/liquidity/app"
	"github.com/Victor118/liquidity/app/params"
	liquiditycli "github.com/Victor118/liquidity/x/liquidity/client/cli"

	"cosmossdk.io/math"
	pruningtypes "cosmossdk.io/store/pruning/types"
	dbm "github.com/cosmos/cosmos-db"
)

// NewConfig returns config that defines the necessary testing requirements
// used to bootstrap and start an in-process local testing network.
func NewConfig(dbm *dbm.MemDB) network.Config {
	encCfg := params.MakeTestEncodingConfig()

	cfg := network.DefaultConfig(app.NewTestNetworkFixture)
	cfg.AppConstructor = NewAppConstructor(encCfg, dbm)                    // the ABCI application constructor
	cfg.GenesisState = liquidityapp.ModuleBasics.DefaultGenesis(cfg.Codec) // liquidity genesis state to provide
	return cfg
}

// NewAppConstructor returns a new network AppConstructor.
func NewAppConstructor(encodingCfg params.EncodingConfig, db *dbm.MemDB) network.AppConstructor {
	return func(val network.ValidatorI) servertypes.Application {
		return app.NewLiquidityApp(
			val.GetCtx().Logger, dbm.NewMemDB(), nil, true,
			simtestutil.NewAppOptionsWithFlagHome(val.GetCtx().Config.RootDir),
			baseapp.SetPruning(pruningtypes.NewPruningOptionsFromString(val.GetAppConfig().Pruning)),
			baseapp.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
			baseapp.SetChainID(val.GetCtx().Viper.GetString(flags.FlagChainID)),
		)
	}
}

var commonArgs = []string{
	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
	fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(10))).String()),
}

// MsgCreatePoolExec creates a transaction for creating liquidity pool.
func MsgCreatePoolExec(clientCtx client.Context, from, poolID, depositCoins string,
	extraArgs ...string) (testutil.BufferWriter, error) {

	args := append([]string{
		poolID,
		depositCoins,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}, commonArgs...)

	args = append(args, commonArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, liquiditycli.NewCreatePoolCmd(), args)
}

// MsgDepositWithinBatchExec creates a transaction to deposit new amounts to the pool.
func MsgDepositWithinBatchExec(clientCtx client.Context, from, poolID, depositCoins string,
	extraArgs ...string) (testutil.BufferWriter, error) {

	args := append([]string{
		poolID,
		depositCoins,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}, commonArgs...)

	args = append(args, commonArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, liquiditycli.NewDepositWithinBatchCmd(), args)
}

// MsgWithdrawWithinBatchExec creates a transaction to withraw pool coin amount from the pool.
func MsgWithdrawWithinBatchExec(clientCtx client.Context, from, poolID, poolCoin string,
	extraArgs ...string) (testutil.BufferWriter, error) {

	args := append([]string{
		poolID,
		poolCoin,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}, commonArgs...)

	args = append(args, commonArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, liquiditycli.NewWithdrawWithinBatchCmd(), args)
}

// MsgSwapWithinBatchExec creates a transaction to swap coins in the pool.
func MsgSwapWithinBatchExec(clientCtx client.Context, from, poolID, swapTypeID,
	offerCoin, demandCoinDenom, orderPrice, swapFeeRate string, extraArgs ...string) (testutil.BufferWriter, error) {

	args := append([]string{
		poolID,
		swapTypeID,
		offerCoin,
		demandCoinDenom,
		orderPrice,
		swapFeeRate,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}, commonArgs...)

	args = append(args, commonArgs...)
	buf, error := clitestutil.ExecTestCLICmd(clientCtx, liquiditycli.NewSwapWithinBatchCmd(), args)
	return buf, error
}

// MsgParamChangeProposalExec creates a transaction for submitting param change proposal
func MsgParamChangeProposalExec(clientCtx client.Context, from string, file string) (testutil.BufferWriter, error) {

	args := append([]string{
		file,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}, commonArgs...)

	paramChangeCmd := paramscli.NewSubmitParamChangeProposalTxCmd()
	flags.AddTxFlagsToCmd(paramChangeCmd)

	return clitestutil.ExecTestCLICmd(clientCtx, paramChangeCmd, args)
}

// MsgVote votes for a proposal
func MsgVote(clientCtx client.Context, from, id, vote string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := append([]string{
		id,
		vote,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}, commonArgs...)

	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, govcli.NewCmdWeightedVote(), args)
}
