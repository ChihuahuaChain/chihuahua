package app_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"

	//"github.com/ChihuahuaChain/chihuahua/app/encoding"
	"github.com/ChihuahuaChain/chihuahua/app"
	"github.com/ChihuahuaChain/chihuahua/testutil/network"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"

	"github.com/ChihuahuaChain/chihuahua/x/payment/types" // what to do about this?
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	cosmosnet "github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	"google.golang.org/grpc"
)

// BasicFuncTestSuite is used to check for basic functionality from the cosmos-sdk
type BasicFuncTestSuite struct {
	suite.Suite

	cfg cosmosnet.Config
	//encCfg   encoding.EncodingConfig
	network  *cosmosnet.Network
	kr       keyring.Keyring
	accounts []string
	conn     *grpc.ClientConn
}

func NewBasicFuncTestSuite(cfg cosmosnet.Config) *BasicFuncTestSuite {
	return &BasicFuncTestSuite{cfg: cfg}
}

func (s *BasicFuncTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	if testing.Short() {
		s.T().Skip("skipping test in unit-tests mode.")
	}

	numAccounts := 100
	s.accounts = make([]string, numAccounts)
	for i := 0; i < numAccounts; i++ {
		s.accounts[i] = tmrand.Str(20)
	}

	net := network.New(s.T(), s.cfg, s.accounts...)

	connection, err := network.GRPCConn(net)
	s.conn = connection
	s.Require().NoError(err)
	s.network = net
	s.kr = net.Validators[0].ClientCtx.Keyring
	//s.encCfg = encoding.MakeEncodingConfig(app.ModuleEncodingRegisters...)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *BasicFuncTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func TestBasicFuncTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.EnableLogging = false
	cfg.MinGasPrices = "0samoleans"
	cfg.NumValidators = 2
	suite.Run(t, NewBasicFuncTestSuite(cfg))
}

func (s *BasicFuncTestSuite) TestBankModule() {
	require := s.Require()
	assert := s.Assert()
	val := s.network.Validators[0]

	type test struct {
		name   string
		msgGen func(client.Context) sdk.Msg
	}

	tests := []test{
		{
			"submit a text governance proposal",
			func(c client.Context) sdk.Msg {
				msg, err := createSendMsg(c, s.accounts[0], s.accounts[1], 1000000)
				require.NoError(err)
				return msg
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			clientCtx := val.ClientCtx
			kr := clientCtx.Keyring
			node, err := clientCtx.GetNode()
			require.NoError(err)

			// quick check of balances
			bals, err := queryForBalance(clientCtx, s.accounts[0], s.conn)
			require.NoError(err)
			fmt.Println(bals)

			signer := types.NewKeyringSigner(kr, s.accounts[0], clientCtx.ChainID)

			err = signer.UpdateAccountFromClient(clientCtx)
			require.NoError(err)

			coin := sdk.Coin{
				Denom:  app.BondDenom,
				Amount: sdk.NewInt(1000000),
			}

			opts := []types.TxBuilderOption{
				types.SetFeeAmount(sdk.NewCoins(coin)),
				types.SetGasLimit(1000000000),
			}

			builder := signer.NewTxBuilder(opts...)

			msg := tc.msgGen(clientCtx)

			tx, err := signer.BuildSignedTx(builder, msg)
			require.NoError(err)

			rawTx, err := s.cfg.TxConfig.TxEncoder()(tx)
			require.NoError(err)

			res, err := val.ClientCtx.BroadcastTxSync(rawTx)
			require.NoError(err)
			assert.Equal(abci.CodeTypeOK, res.Code)
			hexHash := res.TxHash

			// wait a block to clear the txs
			require.NoError(s.network.WaitForNextBlock())
			require.NoError(s.network.WaitForNextBlock())

			hash, err := hex.DecodeString(hexHash)
			require.NoError(err)

			qres, err := node.Tx(context.Background(), hash, false)
			require.NoError(err)

			fmt.Println(qres.TxResult.Code, qres.TxResult.Log, qres.TxResult.GasUsed, qres.TxResult.GasWanted)
		})
	}
}

func queryForBalance(c client.Context, acc string, conn *grpc.ClientConn) (string, error) { // make sure this works
	kr := c.Keyring
	rec, err := kr.Key(acc)
	if err != nil {
		return "", err
	}

	addr := rec.GetAddress()
	// if err != nil {
	// 	return "", err
	// }

	qc := banktypes.NewQueryClient(conn)
	res, err := qc.AllBalances(context.Background(), &banktypes.QueryAllBalancesRequest{
		Address: addr.String(),
	})
	if err != nil {
		return "", err
	}

	return res.Balances.String(), nil
}

func createSendMsg(c client.Context, acc1, acc2 string, amount int64) (sdk.Msg, error) {
	kr := c.Keyring
	rec1, err := kr.Key(acc1)
	if err != nil {
		return nil, err
	}
	addr1 := rec1.GetAddress()
	// if err != nil {
	// 	return nil, err
	// }
	rec2, err := kr.Key(acc2)
	if err != nil {
		return nil, err
	}
	addr2 := rec2.GetAddress()
	// if err != nil {
	// 	return nil, err
	// }
	coins := sdk.NewCoins(
		sdk.NewCoin(app.BondDenom, sdk.NewInt(amount)),
	)
	return banktypes.NewMsgSend(addr1, addr2, coins), nil
}
