package cli_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govcli "github.com/cosmos/cosmos-sdk/x/gov/client/cli"

	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"

	liquidityapp "github.com/ChihuahuaChain/chihuahua/x/liquidity/testutil/app"
	"github.com/cosmos/cosmos-sdk/crypto/hd"

	"cosmossdk.io/math"
	"github.com/gogo/protobuf/proto"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	genutiltest "github.com/cosmos/cosmos-sdk/x/genutil/client/testutil"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/ChihuahuaChain/chihuahua/x/liquidity"
	"github.com/ChihuahuaChain/chihuahua/x/liquidity/client/cli"
	liquiditytestutil "github.com/ChihuahuaChain/chihuahua/x/liquidity/client/testutil"
	liquiditytypes "github.com/ChihuahuaChain/chihuahua/x/liquidity/types"

	"cosmossdk.io/log/v2"
	tmdb "github.com/cometbft/cometbft-db"
	tmcli "github.com/cometbft/cometbft/libs/cli"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network

	db *tmdb.MemDB // in-memory database is needed for exporting genesis cli integration test
}

// SetupTest creates a new network for _each_ integration test. We create a new
// network for each test because there are some state modifications that are
// needed to be made in order to make useful queries. However, we don't want
// these state changes to be present in other tests.
const testValidatorMnemonic = "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art"

func (s *IntegrationTestSuite) SetupTest() {

	s.T().Helper()

	//privVal := mock.NewPV()
	//pubKey, err := privVal.GetPubKey()
	//s.Require().NoError(err)
	// create validator set with single validator
	//validator := tmtypes.NewValidator(pubKey, 1)
	//valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account
	//senderPrivKey := secp256k1.GenPrivKey()
	//acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	// balance := banktypes.Balance{
	// 	Address: acc.GetAddress().String(),
	// 	Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100000000000000))),
	// }

	//addr, _ := sdk.AccAddressFromHexUnsafe(validator.Address.String())
	//s.T().Logf("Validator address : " + addr.String())
	// balanceVal := banktypes.Balance{
	// 	Address: addr.String(),
	// 	Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100_000_000_000)), sdk.NewCoin("node0token", math.NewInt(100_000_000_000))),
	// }

	db := tmdb.NewMemDB()

	//cfg := liquiditytestutil.NewConfig(db)
	//cfg.NumValidators = 1
	// run on the chihuahua app, with a validator derived from a fixed mnemonic
	// so that it can be allowed to create pools in the genesis
	cfg := network.DefaultConfig(liquidityapp.NewTestNetworkFixture)
	cfg.NumValidators = 1
	cfg.Mnemonics = []string{testValidatorMnemonic}
	cfg.AccountTokens = math.NewInt(1_000_000_000_000)
	cfg.StakingTokens = math.NewInt(1_000_000_000_000)
	valKey, err := hd.Secp256k1.Derive()(testValidatorMnemonic, "", sdk.GetConfig().GetFullBIP44Path())
	s.Require().NoError(err)
	valAddr := sdk.AccAddress(hd.Secp256k1.Generate()(valKey).PubKey().Address())
	var liquidityGenesis liquiditytypes.GenesisState
	s.Require().NoError(cfg.Codec.UnmarshalJSON(cfg.GenesisState[liquiditytypes.ModuleName], &liquidityGenesis))
	liquidityGenesis.Params.PoolPermissionedCreatorAddresses = []string{valAddr.String()}
	cfg.GenesisState[liquiditytypes.ModuleName] = cfg.Codec.MustMarshalJSON(&liquidityGenesis)
	//genesisState, err := simtestutil.GenesisStateWithValSet(cfg.Codec, cfg.GenesisState, valSet, []authtypes.GenesisAccount{acc}, balance, balanceVal)
	//s.Require().NoError(err)
	//cfg.GenesisState = genesisState
	// var liquidityGenState liquiditytypes.GenesisState
	// err = cfg.Codec.UnmarshalJSON(cfg.GenesisState[liquiditytypes.ModuleName], &liquidityGenState)
	// s.Require().NoError(err)
	//probably put theses coins and put it in simtestutil.GenesisStateWithValSet in balance parameters ?
	//cfg.AccountTokens = math.NewInt(100_000_000_000) // node0token denom
	//cfg.StakingTokens = math.NewInt(100_000_000_000) // stake denom

	genesisStateGov := v1.DefaultGenesisState()
	maxDepPeriod := time.Duration(15) * time.Second
	votingPeriod := time.Duration(5) * time.Second
	genesisStateGov.Params.MaxDepositPeriod = &maxDepPeriod
	genesisStateGov.Params.VotingPeriod = &votingPeriod
	var depositPeriod time.Duration = time.Duration(15) * time.Second
	genesisStateGov.Params.MaxDepositPeriod = &depositPeriod
	genesisStateGov.Params.Quorum = "0.01"
	bz, err := cfg.Codec.MarshalJSON(genesisStateGov)
	s.Require().NoError(err)
	cfg.GenesisState["gov"] = bz

	s.cfg = cfg
	s.network, err = network.New(s.T(), s.T().TempDir(), s.cfg)
	if err != nil {
		s.T().Log(err)
	}
	s.Require().NoError(err)
	s.db = db

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

// TearDownTest cleans up the curret test network after each test in the suite.
func (s *IntegrationTestSuite) TearDownTest() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
	time.Sleep(100 * time.Millisecond)
}

// TestIntegrationTestSuite every integration test suite.
func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestNewCreatePoolCmd() {
	val := s.network.Validators[0]
	// use two different tokens that are minted to the test account
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)

	testCases := []struct {
		name         string
		args         []string
		expectErr    bool
		respType     proto.Message
		expectedCode uint32
	}{
		{
			"invalid pool type id",
			[]string{
				"invalidpooltypeid",
				sdk.NewCoins(sdk.NewCoin(denomX, math.NewInt(100_000_000)), sdk.NewCoin(denomY, math.NewInt(100_000_000))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastAsync),
				fmt.Sprintf("--%s=%d", flags.FlagGas, 1_000_000),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.NewInt(10))).String()),
			},
			true, nil, 0,
		},
		{
			"pool type id is not supported",
			[]string{
				fmt.Sprintf("%d", uint32(2)),
				sdk.NewCoins(sdk.NewCoin(denomX, math.NewInt(100_000_000)), sdk.NewCoin(denomY, math.NewInt(100_000_000)), sdk.NewCoin("denomZ", math.NewInt(100_000_000))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastAsync),
				fmt.Sprintf("--%s=%d", flags.FlagGas, 1_000_000),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.NewInt(10))).String()),
			},
			true, nil, 0,
		},
		{
			"invalid number of denoms",
			[]string{
				fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeID),
				sdk.NewCoins(sdk.NewCoin(denomX, math.NewInt(100_000_000)), sdk.NewCoin(denomY, math.NewInt(100_000_000)), sdk.NewCoin("denomZ", math.NewInt(100_000_000))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastAsync),
				fmt.Sprintf("--%s=%d", flags.FlagGas, 1_000_000),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.NewInt(10))).String()),
			},
			true, nil, 0,
		},
		{
			"deposit coin less than minimum deposit amount",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				sdk.NewCoins(sdk.NewCoin(denomX, math.NewInt(1_000)), sdk.NewCoin(denomY, math.NewInt(10_000_000))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastAsync),
				fmt.Sprintf("--%s=%d", flags.FlagGas, 1_000_000),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.NewInt(10))).String()),
			},
			false, &sdk.TxResponse{}, 9,
		},
		{
			"valid transaction",
			[]string{
				fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeID),
				sdk.NewCoins(sdk.NewCoin(denomX, math.NewInt(100_000_000)), sdk.NewCoin(denomY, math.NewInt(100_000_000))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastAsync),
				fmt.Sprintf("--%s=%d", flags.FlagGas, 1_000_000),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.NewInt(10))).String()),
			},
			false, &sdk.TxResponse{}, 0,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewCreatePoolCmd()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err, out.String())
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())

				txResp := s.txResult(clientCtx, out.Bytes())
				s.Require().Equal(tc.expectedCode, txResp.Code, txResp.RawLog)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestNewDepositWithinBatchCmd() {
	val := s.network.Validators[0]

	// use two different tokens that are minted to the test account
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)

	// liquidity pool should be created prior to test this integration test
	_, err := liquiditytestutil.MsgCreatePoolExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeID),
		sdk.NewCoins(sdk.NewCoin(denomX, math.NewInt(100_000_000)), sdk.NewCoin(denomY, math.NewInt(100_000_000))).String(),
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	testCases := []struct {
		name         string
		args         []string
		expectErr    bool
		respType     proto.Message
		expectedCode uint32
	}{
		{
			"invalid pool id",
			[]string{
				"invalidpoolid",
				sdk.NewCoins(sdk.NewCoin(denomX, math.NewInt(1_000_000)), sdk.NewCoin(denomY, math.NewInt(1_000_000))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastAsync),
				fmt.Sprintf("--%s=%d", flags.FlagGas, 1_000_000),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.NewInt(10))).String()),
			},
			true, nil, 0,
		},
		{
			"invalid number of denoms",
			[]string{
				fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeID),
				sdk.NewCoins(sdk.NewCoin(denomX, math.NewInt(1_000_000)), sdk.NewCoin(denomY, math.NewInt(1_000_000)), sdk.NewCoin("denomZ", math.NewInt(1_000_000))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastAsync),
				fmt.Sprintf("--%s=%d", flags.FlagGas, 1_000_000),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.NewInt(10))).String()),
			},
			true, nil, 0,
		},
		{
			"valid transaction",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				sdk.NewCoins(sdk.NewCoin(denomX, math.NewInt(10_000_000)), sdk.NewCoin(denomY, math.NewInt(10_000_000))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastAsync),
				fmt.Sprintf("--%s=%d", flags.FlagGas, 1_000_000),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.NewInt(10))).String()),
			},
			false, &sdk.TxResponse{}, 0,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewDepositWithinBatchCmd()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err, out.String())
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())

				txResp := s.txResult(clientCtx, out.Bytes())
				s.Require().Equal(tc.expectedCode, txResp.Code, txResp.RawLog)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestNewWithdrawWithinBatchCmd() {
	val := s.network.Validators[0]

	// use two different tokens that are minted to the test account
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)

	// liquidity pool should be created prior to test this integration test
	_, err := liquiditytestutil.MsgCreatePoolExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeID),
		sdk.NewCoins(sdk.NewCoin(denomX, math.NewInt(100_000_000)), sdk.NewCoin(denomY, math.NewInt(100_000_000))).String(),
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	testCases := []struct {
		name         string
		args         []string
		expectErr    bool
		respType     proto.Message
		expectedCode uint32
	}{
		{
			"invalid pool id",
			[]string{
				"invalidpoolid",
				sdk.NewCoins(sdk.NewCoin("poolC33A77E752C183913636A37FE1388ACA22FE7BED792BEB2E72EF2DA857703D8D", math.NewInt(10_000))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastAsync),
				fmt.Sprintf("--%s=%d", flags.FlagGas, 1_000_000),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.NewInt(10))).String()),
			},
			true, nil, 0,
		},
		{
			"bad pool coin",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				sdk.NewCoins(sdk.NewCoin("badpoolcoindenom", math.NewInt(10_000))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastAsync),
				fmt.Sprintf("--%s=%d", flags.FlagGas, 1_000_000),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.NewInt(10))).String()),
			},
			false, &sdk.TxResponse{}, 29,
		},
		{
			"valid transaction",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				sdk.NewCoins(sdk.NewCoin("poolC33A77E752C183913636A37FE1388ACA22FE7BED792BEB2E72EF2DA857703D8D", math.NewInt(10_000))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastAsync),
				fmt.Sprintf("--%s=%d", flags.FlagGas, 1_000_000),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.NewInt(10))).String()),
			},
			false, &sdk.TxResponse{}, 0,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewWithdrawWithinBatchCmd()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err, out.String())
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())

				txResp := s.txResult(clientCtx, out.Bytes())
				s.Require().Equal(tc.expectedCode, txResp.Code, txResp.RawLog)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestNewSwapWithinBatchCmd() {

	val := s.network.Validators[0]

	// use two different tokens that are minted to the test account
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)

	// liquidity pool should be created prior to test this integration test
	_, err := liquiditytestutil.MsgCreatePoolExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeID),
		sdk.NewCoins(sdk.NewCoin(denomX, math.NewInt(100_000_000)), sdk.NewCoin(denomY, math.NewInt(100_000_000))).String(),
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	testCases := []struct {
		name         string
		args         []string
		expectErr    bool
		respType     proto.Message
		expectedCode uint32
	}{
		{
			"invalid pool id",
			[]string{
				"invalidpoolid",
				fmt.Sprintf("%d", uint32(1)),
				sdk.NewCoins(sdk.NewCoin(denomX, math.NewInt(10_000))).String(),
				denomY,
				fmt.Sprintf("%.2f", 0.02),
				fmt.Sprintf("%.3f", 0.003),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastAsync),
				fmt.Sprintf("--%s=%d", flags.FlagGas, 1_000_000),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.NewInt(10))).String()),
			},
			true, nil, 0,
		},
		{
			"swap type id not supported",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("%d", uint32(2)),
				sdk.NewCoins(sdk.NewCoin(denomX, math.NewInt(10_000))).String(),
				denomY,
				fmt.Sprintf("%.2f", 0.02),
				fmt.Sprintf("%.2f", 0.03),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastAsync),
				fmt.Sprintf("--%s=%d", flags.FlagGas, 1_000_000),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.NewInt(10))).String()),
			},
			true, nil, 0,
		},
		{
			"bad offer coin fee",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("%d", uint32(1)),
				sdk.NewCoins(sdk.NewCoin(denomX, math.NewInt(10_000))).String(),
				denomY,
				fmt.Sprintf("%.2f", 0.02),
				fmt.Sprintf("%.2f", 0.01),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastAsync),
				fmt.Sprintf("--%s=%d", flags.FlagGas, 1_000_000),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.NewInt(10))).String()),
			},
			false, &sdk.TxResponse{}, 35,
		},
		{
			"valid transaction",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("%d", liquiditytypes.DefaultSwapTypeID),
				sdk.NewCoins(sdk.NewCoin(denomX, math.NewInt(10_000))).String(),
				denomY,
				fmt.Sprintf("%.2f", 1.0),
				fmt.Sprintf("%.3f", 0.003),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastAsync),
				fmt.Sprintf("--%s=%d", flags.FlagGas, 1_000_000),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.NewInt(10))).String()),
			},
			false, &sdk.TxResponse{}, 0,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewSwapWithinBatchCmd()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err, out.String())
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())

				txResp := s.txResult(clientCtx, out.Bytes())
				s.Require().Equal(tc.expectedCode, txResp.Code, txResp.RawLog)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetCmdQueryParams() {

	val := s.network.Validators[0]

	testCases := []struct {
		name           string
		args           []string
		expectedOutput string
	}{
		{
			"json output",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			`{"pool_types":[{"id":1,"name":"StandardLiquidityPool","min_reserve_coin_num":2,"max_reserve_coin_num":2,"description":"Standard liquidity pool with pool price function X/Y, ESPM constraint, and two kinds of reserve coins"}],"min_init_deposit_amount":"1000000","init_pool_coin_mint_amount":"1000000","max_reserve_coin_amount":"0","pool_creation_fee":[{"denom":"stake","amount":"40000000"}],"swap_fee_rate":"0.003000000000000000","withdraw_fee_rate":"0.000000000000000000","max_order_amount_ratio":"0.100000000000000000","unit_batch_height":1,"circuit_breaker_enabled":false,"builders_addresses":[],"builders_commission":"0.200000000000000000","pool_permissioned_creator_addresses":["VALADDR"]}`,
		},
		{
			"text output",
			[]string{fmt.Sprintf("--%s=text", tmcli.OutputFlag)},
			`builders_addresses: []
builders_commission: "0.200000000000000000"
circuit_breaker_enabled: false
init_pool_coin_mint_amount: "1000000"
max_order_amount_ratio: "0.100000000000000000"
max_reserve_coin_amount: "0"
min_init_deposit_amount: "1000000"
pool_creation_fee:
- amount: "40000000"
  denom: stake
pool_permissioned_creator_addresses:
- VALADDR
pool_types:
- description: Standard liquidity pool with pool price function X/Y, ESPM constraint,
    and two kinds of reserve coins
  id: 1
  max_reserve_coin_num: 2
  min_reserve_coin_num: 2
  name: StandardLiquidityPool
swap_fee_rate: "0.003000000000000000"
unit_batch_height: 1
withdraw_fee_rate: "0.000000000000000000"`,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryParams()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			s.Require().NoError(err)
			s.Require().Equal(strings.ReplaceAll(tc.expectedOutput, "VALADDR", val.Address.String()), strings.TrimSpace(out.String()))
		})
	}
}

func (s *IntegrationTestSuite) TestGetCmdQueryLiquidityPool() {

	val := s.network.Validators[0]

	// use two different tokens that are minted to the test account
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)

	// liquidity pool should be created prior to test this integration test
	_, err := liquiditytestutil.MsgCreatePoolExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeID),
		sdk.NewCoins(sdk.NewCoin(denomX, math.NewInt(100_000_000)), sdk.NewCoin(denomY, math.NewInt(100_000_000))).String(),
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			"with invalid pool id",
			[]string{
				"invalidpoolid",
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"with not supported pool id",
			[]string{
				fmt.Sprintf("%d", uint32(2)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"valid case with pool id",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			false,
		},
		{
			"with invalid pool coin denom",
			[]string{
				fmt.Sprintf("--%s=%s", cli.FlagPoolCoinDenom, "invalid_value"),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"with empty pool coin denom",
			[]string{
				fmt.Sprintf("--%s", cli.FlagPoolCoinDenom),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"valid case with pool coin denom",
			[]string{
				fmt.Sprintf("--%s=%s", cli.FlagPoolCoinDenom, "poolC33A77E752C183913636A37FE1388ACA22FE7BED792BEB2E72EF2DA857703D8D"),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			false,
		},
		{
			"with invalid reserve acc",
			[]string{
				fmt.Sprintf("--%s=%s", cli.FlagReserveAcc, "invalid_value"),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"with empty reserve acc",
			[]string{
				fmt.Sprintf("--%s", cli.FlagReserveAcc),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"valid case with reserve acc",
			[]string{
				fmt.Sprintf("--%s=%s", cli.FlagReserveAcc, liquiditytypes.GetPoolReserveAcc(liquiditytypes.PoolName([]string{denomX, denomY}, liquiditytypes.DefaultPoolTypeID)).String()),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryLiquidityPool()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				var resp liquiditytypes.QueryLiquidityPoolResponse
				err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &resp)
				s.Require().NoError(err)
				s.Require().Equal(uint64(1), resp.GetPool().Id)
				s.Require().Equal(uint32(1), resp.GetPool().TypeId)
				s.Require().Len(resp.GetPool().ReserveCoinDenoms, 2)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetCmdQueryLiquidityPools() {

	val := s.network.Validators[0]

	// use two different tokens that are minted to the test account
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)

	// liquidity pool should be created prior to test this integration test
	_, err := liquiditytestutil.MsgCreatePoolExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeID),
		sdk.NewCoins(sdk.NewCoin(denomX, math.NewInt(100_000_000)), sdk.NewCoin(denomY, math.NewInt(100_000_000))).String(),
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			"valid case",
			[]string{
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryLiquidityPools()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				var resps liquiditytypes.QueryLiquidityPoolsResponse
				err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &resps)
				s.Require().NoError(err)

				for _, pool := range resps.GetPools() {
					s.Require().Equal(uint64(1), pool.Id)
					s.Require().Equal(uint32(1), pool.TypeId)
					s.Require().Len(pool.ReserveCoinDenoms, 2)
				}
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetCmdQueryLiquidityPoolBatch() {

	val := s.network.Validators[0]

	// use two different tokens that are minted to the test account
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)

	// liquidity pool should be created prior to test this integration test
	createOut, err := liquiditytestutil.MsgCreatePoolExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeID),
		sdk.NewCoins(sdk.NewCoin(denomX, math.NewInt(100_000_000)), sdk.NewCoin(denomY, math.NewInt(100_000_000))).String(),
	)
	s.Require().NoError(err)
	s.Require().Zero(s.txResult(val.ClientCtx, createOut.Bytes()).Code)

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			"with invalid pool id",
			[]string{
				"invalidpoolid",
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"with not supported pool id",
			[]string{
				fmt.Sprintf("%d", uint32(2)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"valid case",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryLiquidityPoolBatch()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				var resp liquiditytypes.QueryLiquidityPoolBatchResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &resp)
				s.Require().NoError(err)
				s.Require().Equal(uint64(1), resp.GetBatch().PoolId)
				s.Require().Equal(false, resp.GetBatch().Executed)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchDepositMsg() {

	val := s.network.Validators[0]

	// use two different tokens that are minted to the test account
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)

	// liquidity pool should be created prior to test this integration test
	_, err := liquiditytestutil.MsgCreatePoolExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeID),
		sdk.NewCoins(sdk.NewCoin(denomX, math.NewInt(100_000_000)), sdk.NewCoin(denomY, math.NewInt(100_000_000))).String(),
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	// create new deposit
	batchOut, err := liquiditytestutil.MsgDepositWithinBatchExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeID),
		sdk.NewCoins(sdk.NewCoin(denomX, math.NewInt(10_000_000)), sdk.NewCoin(denomY, math.NewInt(10_000_000))).String(),
	)
	s.Require().NoError(err)
	batchHeight := s.txResult(val.ClientCtx, batchOut.Bytes()).Height

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			"with invalid pool id",
			[]string{
				"invalidpoolid",
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"with not supported pool id",
			[]string{
				fmt.Sprintf("%d", uint32(2)),
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"valid case",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryPoolBatchDepositMsg()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, append(tc.args, fmt.Sprintf("--%s=%d", flags.FlagHeight, batchHeight)))

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				var resp liquiditytypes.QueryPoolBatchDepositMsgResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &resp)
				s.Require().NoError(err)
				s.Require().Equal(val.Address.String(), resp.GetDeposit().Msg.DepositorAddress)
				s.Require().Equal(true, resp.GetDeposit().Executed)
				s.Require().Equal(true, resp.GetDeposit().Succeeded)
				s.Require().Equal(true, resp.GetDeposit().ToBeDeleted)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchDepositMsgs() {

	val := s.network.Validators[0]

	// use two different tokens that are minted to the test account
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)

	// liquidity pool should be created prior to test this integration test
	_, err := liquiditytestutil.MsgCreatePoolExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeID),
		sdk.NewCoins(sdk.NewCoin(denomX, math.NewInt(100_000_000)), sdk.NewCoin(denomY, math.NewInt(100_000_000))).String(),
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	// create new deposit
	batchOut, err := liquiditytestutil.MsgDepositWithinBatchExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeID),
		sdk.NewCoins(sdk.NewCoin(denomX, math.NewInt(10_000_000)), sdk.NewCoin(denomY, math.NewInt(10_000_000))).String(),
	)
	s.Require().NoError(err)
	batchHeight := s.txResult(val.ClientCtx, batchOut.Bytes()).Height

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			"with invalid pool id",
			[]string{
				"invalidpoolid",
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"with not supported pool id",
			[]string{
				fmt.Sprintf("%d", uint32(2)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"valid case",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryPoolBatchDepositMsgs()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, append(tc.args, fmt.Sprintf("--%s=%d", flags.FlagHeight, batchHeight)))

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				var resps liquiditytypes.QueryPoolBatchDepositMsgsResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &resps)
				s.Require().NoError(err)

				for _, deposit := range resps.GetDeposits() {
					s.Require().Equal(true, deposit.Executed)
					s.Require().Equal(true, deposit.Succeeded)
					s.Require().Equal(true, deposit.ToBeDeleted)
				}
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchWithdrawMsg() {

	val := s.network.Validators[0]

	// use two different tokens that are minted to the test account
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)

	// liquidity pool should be created prior to test this integration test
	_, err := liquiditytestutil.MsgCreatePoolExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeID),
		sdk.NewCoins(sdk.NewCoin(denomX, math.NewInt(100_000_000)), sdk.NewCoin(denomY, math.NewInt(100_000_000))).String(),
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	// withdraw pool coin from the pool
	poolCoinDenom := "poolC33A77E752C183913636A37FE1388ACA22FE7BED792BEB2E72EF2DA857703D8D"
	batchOut, err := liquiditytestutil.MsgWithdrawWithinBatchExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", uint32(1)),
		sdk.NewCoins(sdk.NewCoin(poolCoinDenom, math.NewInt(10_000))).String(),
	)
	s.Require().NoError(err)
	batchHeight := s.txResult(val.ClientCtx, batchOut.Bytes()).Height

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			"with invalid pool id",
			[]string{
				"invalidpoolid",
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"with not suppoorted pool id",
			[]string{
				fmt.Sprintf("%d", uint32(2)),
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"valid case",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryPoolBatchWithdrawMsg()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, append(tc.args, fmt.Sprintf("--%s=%d", flags.FlagHeight, batchHeight)))

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				var resp liquiditytypes.QueryPoolBatchWithdrawMsgResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &resp)
				s.Require().NoError(err)
				s.Require().Equal(val.Address.String(), resp.GetWithdraw().Msg.WithdrawerAddress)
				s.Require().Equal(poolCoinDenom, resp.GetWithdraw().Msg.PoolCoin.Denom)
				s.Require().Equal(true, resp.GetWithdraw().Executed)
				s.Require().Equal(true, resp.GetWithdraw().Succeeded)
				s.Require().Equal(true, resp.GetWithdraw().ToBeDeleted)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchWithdrawMsgs() {

	val := s.network.Validators[0]

	// use two different tokens that are minted to the test account
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)

	// liquidity pool should be created prior to test this integration test
	_, err := liquiditytestutil.MsgCreatePoolExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeID),
		sdk.NewCoins(sdk.NewCoin(denomX, math.NewInt(100_000_000)), sdk.NewCoin(denomY, math.NewInt(100_000_000))).String(),
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	// withdraw pool coin from the pool
	batchOut, err := liquiditytestutil.MsgWithdrawWithinBatchExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", uint32(1)),
		sdk.NewCoins(sdk.NewCoin("poolC33A77E752C183913636A37FE1388ACA22FE7BED792BEB2E72EF2DA857703D8D", math.NewInt(10_000))).String(),
	)
	s.Require().NoError(err)
	batchHeight := s.txResult(val.ClientCtx, batchOut.Bytes()).Height

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			"with invalid pool id",
			[]string{
				"invalidpoolid",
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"with not supported pool id",
			[]string{
				fmt.Sprintf("%d", uint32(2)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"valid case",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryPoolBatchWithdrawMsgs()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, append(tc.args, fmt.Sprintf("--%s=%d", flags.FlagHeight, batchHeight)))

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				var resps liquiditytypes.QueryPoolBatchWithdrawMsgsResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &resps)
				s.Require().NoError(err)

				for _, withdraw := range resps.GetWithdraws() {
					s.Require().Equal(val.Address.String(), withdraw.Msg.WithdrawerAddress)
					s.Require().Equal(true, withdraw.Executed)
					s.Require().Equal(true, withdraw.Succeeded)
					s.Require().Equal(true, withdraw.ToBeDeleted)
				}
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchSwapMsg() {

	val := s.network.Validators[0]

	// use two different tokens that are minted to the test account
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)
	X := sdk.NewCoin(denomX, math.NewInt(1_000_000_000))
	Y := sdk.NewCoin(denomY, math.NewInt(5_000_000_000))

	// liquidity pool should be created prior to test this integration test
	_, err := liquiditytestutil.MsgCreatePoolExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeID),
		sdk.NewCoins(X, Y).String(),
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	// swap coins from the pool
	offerCoin := sdk.NewCoin(denomY, math.NewInt(50_000_000))
	batchOut, err := liquiditytestutil.MsgSwapWithinBatchExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", uint32(1)),
		fmt.Sprintf("%d", liquiditytypes.DefaultSwapTypeID),
		offerCoin.String(),
		denomX,
		fmt.Sprintf("%.3f", 0.019),
		fmt.Sprintf("%.3f", 0.003),
	)
	s.Require().NoError(err)
	batchHeight := s.txResult(val.ClientCtx, batchOut.Bytes()).Height

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			"with invalid pool id",
			[]string{
				"invalidpoolid",
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"with not supported pool id",
			[]string{
				fmt.Sprintf("%d", uint32(2)),
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"valid case",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryPoolBatchSwapMsg()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, append(tc.args, fmt.Sprintf("--%s=%d", flags.FlagHeight, batchHeight)))

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				var resp liquiditytypes.QueryPoolBatchSwapMsgResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &resp)
				s.Require().NoError(err)
				s.Require().Equal(val.Address.String(), resp.GetSwap().Msg.SwapRequesterAddress)
				s.Require().Equal(true, resp.GetSwap().Executed)
				s.Require().Equal(true, resp.GetSwap().Succeeded)
				s.Require().Equal(true, resp.GetSwap().ToBeDeleted)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetCircuitBreaker() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	// use two different tokens that are minted to the test account
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)
	X := sdk.NewCoin(denomX, math.NewInt(1_000_000_000))
	Y := sdk.NewCoin(denomY, math.NewInt(5_000_000_000))
	requireCode := func(out testutil.BufferWriter, err error, code uint32, log string) {
		s.T().Helper()
		s.Require().NoError(err)
		res := s.txResult(clientCtx, out.Bytes())
		s.Require().Equal(code, res.Code, res.RawLog)
		if log != "" {
			s.Require().Contains(res.RawLog, log)
		}
	}

	// liquidity pool should be created prior to test this integration test
	out, err := liquiditytestutil.MsgCreatePoolExec(clientCtx, val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeID), sdk.NewCoins(X, Y).String())
	requireCode(out, err, 0, "")

	// swap coins from the pool
	offerCoin := sdk.NewCoin(denomY, math.NewInt(50_000_000))
	swap := func() (testutil.BufferWriter, error) {
		return liquiditytestutil.MsgSwapWithinBatchExec(clientCtx, val.Address.String(), "1",
			fmt.Sprintf("%d", liquiditytypes.DefaultSwapTypeID), offerCoin.String(), denomX,
			fmt.Sprintf("%.3f", 0.019), fmt.Sprintf("%.3f", 0.003))
	}
	out, err = swap()
	requireCode(out, err, 0, "")

	// enable the circuit breaker through governance
	paramsOut, err := clitestutil.ExecTestCLICmd(clientCtx, cli.GetCmdQueryParams(), []string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)})
	s.Require().NoError(err)
	var params liquiditytypes.Params
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(paramsOut.Bytes(), &params))
	params.CircuitBreakerEnabled = true
	msgBz, err := clientCtx.Codec.MarshalInterfaceJSON(&liquiditytypes.MsgUpdateParams{
		Authority: authtypes.NewModuleAddress("gov").String(),
		Params:    params,
	})
	s.Require().NoError(err)
	proposal := fmt.Sprintf(`{"messages":[%s],"metadata":"","deposit":"%s","title":"enable circuit breaker","summary":"enable circuit breaker"}`,
		msgBz, sdk.NewCoin(s.cfg.BondDenom, v1.DefaultMinDepositTokens))
	out, err = clitestutil.ExecTestCLICmd(clientCtx, govcli.NewCmdSubmitProposal(),
		append([]string{testutil.WriteToNewTempFile(s.T(), proposal).Name(), fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String())}, liquiditytestutil.CommonArgs()...))
	requireCode(out, err, 0, "")
	out, err = liquiditytestutil.MsgVote(clientCtx, val.Address.String(), "1", "yes")
	requireCode(out, err, 0, "")
	s.Require().Eventually(func() bool {
		out, err := clitestutil.ExecTestCLICmd(clientCtx, cli.GetCmdQueryParams(), []string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)})
		return err == nil && strings.Contains(out.String(), `"circuit_breaker_enabled":true`)
	}, 30*time.Second, time.Second, "the circuit breaker was not enabled")

	// swaps, new pools and deposits fail because of the circuit breaker
	out, err = swap()
	requireCode(out, err, 40, "circuit breaker is triggered")
	out, err = liquiditytestutil.MsgCreatePoolExec(clientCtx, val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeID), sdk.NewCoins(X, Y).String())
	requireCode(out, err, 40, "circuit breaker is triggered")
	out, err = liquiditytestutil.MsgDepositWithinBatchExec(clientCtx, val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeID),
		sdk.NewCoins(sdk.NewCoin(denomX, math.NewInt(10_000_000)), sdk.NewCoin(denomY, math.NewInt(10_000_000))).String())
	requireCode(out, err, 40, "circuit breaker is triggered")

	// withdrawing the pool coins still works, until the pool is depleted
	poolCoinDenom := liquiditytypes.GetPoolCoinDenom(liquiditytypes.PoolName([]string{denomX, denomY}, liquiditytypes.DefaultPoolTypeID))
	withdraw := func(amount int64) (testutil.BufferWriter, error) {
		return liquiditytestutil.MsgWithdrawWithinBatchExec(clientCtx, val.Address.String(), "1",
			sdk.NewCoins(sdk.NewCoin(poolCoinDenom, math.NewInt(amount))).String())
	}
	for _, amount := range []int64{500000, 499999, 1} {
		out, err = withdraw(amount)
		requireCode(out, err, 0, "")
	}
	out, err = withdraw(1)
	requireCode(out, err, 39, "the pool is depleted of reserve coin")
}

func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchSwapMsgs() {

	val := s.network.Validators[0]

	// use two different tokens that are minted to the test account
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)
	X := sdk.NewCoin(denomX, math.NewInt(1_000_000_000))
	Y := sdk.NewCoin(denomY, math.NewInt(5_000_000_000))

	// liquidity pool should be created prior to test this integration test
	_, err := liquiditytestutil.MsgCreatePoolExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeID),
		sdk.NewCoins(X, Y).String(),
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	// swap coins from the pool
	offerCoin := sdk.NewCoin(denomY, math.NewInt(50_000_000))
	batchOut, err := liquiditytestutil.MsgSwapWithinBatchExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", uint32(1)),
		fmt.Sprintf("%d", liquiditytypes.DefaultSwapTypeID),
		offerCoin.String(),
		denomX,
		fmt.Sprintf("%.3f", 0.019),
		fmt.Sprintf("%.3f", 0.003),
	)
	s.Require().NoError(err)
	batchHeight := s.txResult(val.ClientCtx, batchOut.Bytes()).Height

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			"with invalid pool id",
			[]string{
				"invalidpoolid",
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"with not supported pool id",
			[]string{
				fmt.Sprintf("%d", uint32(2)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"valid case",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryPoolBatchSwapMsgs()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, append(tc.args, fmt.Sprintf("--%s=%d", flags.FlagHeight, batchHeight)))

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				var resps liquiditytypes.QueryPoolBatchSwapMsgsResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &resps)
				s.Require().NoError(err)

				for _, swap := range resps.GetSwaps() {
					s.Require().Equal(val.Address.String(), swap.Msg.SwapRequesterAddress)
					s.Require().Equal(true, swap.Executed)
					s.Require().Equal(true, swap.Succeeded)
					s.Require().Equal(true, swap.ToBeDeleted)
				}
			}
		})
	}
}

func (s *IntegrationTestSuite) TestInitGenesis() {

	testCases := []struct {
		name      string
		flags     func(dir string) []string
		expectErr bool
		err       error
	}{
		{
			name: "default genesis state",
			flags: func(dir string) []string {
				return []string{
					"liquidity-test",
				}
			},
			expectErr: false,
			err:       nil,
		},
	}
	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			testMbm := module.NewBasicManager(liquidity.AppModuleBasic{})

			home := s.T().TempDir()
			logger := log.NewNopLogger()
			cfg, err := genutiltest.CreateDefaultCometConfig(home)
			s.Require().NoError(err)

			serverCtx := server.NewContext(viper.New(), cfg, logger)
			interfaceRegistry := types.NewInterfaceRegistry()
			marshaler := codec.NewProtoCodec(interfaceRegistry)
			clientCtx := client.Context{}.
				WithCodec(marshaler).
				WithHomeDir(home)

			ctx := context.Background()
			ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)
			ctx = context.WithValue(ctx, server.ServerContextKey, serverCtx)

			cmd := genutilcli.InitCmd(testMbm, home)
			cmd.SetArgs(
				tc.flags(home),
			)

			if tc.expectErr {
				err := cmd.ExecuteContext(ctx)
				s.Require().EqualError(err, tc.err.Error())
			} else {
				s.Require().NoError(cmd.ExecuteContext(ctx))
			}
		})
	}
}

// txResult waits for the tx in out to be executed and returns its result.
func (s *IntegrationTestSuite) txResult(clientCtx client.Context, out []byte) *sdk.TxResponse {
	var resp sdk.TxResponse
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out, &resp), string(out))
	if resp.Code != 0 {
		return &resp // refused by CheckTx
	}
	for i := 0; i < 10; i++ {
		s.Require().NoError(s.network.WaitForNextBlock())
		res, err := clitestutil.ExecTestCLICmd(clientCtx, authcmd.QueryTxCmd(), []string{resp.TxHash, "--output=json"})
		if err != nil {
			continue
		}
		var result sdk.TxResponse
		s.Require().NoError(clientCtx.Codec.UnmarshalJSON(res.Bytes(), &result))
		return &result
	}
	s.FailNow("tx not found: " + resp.TxHash)
	return nil
}
