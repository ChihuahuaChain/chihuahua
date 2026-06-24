package app

// DONTCOVER

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"testing"
	"time"

	"cosmossdk.io/log"
	mathsdk "cosmossdk.io/math"
	pruningtypes "cosmossdk.io/store/pruning/types"
	abci "github.com/cometbft/cometbft/abci/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	tmtypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/Victor118/liquidity/x/liquidity"
	"github.com/Victor118/liquidity/x/liquidity/keeper"
	"github.com/Victor118/liquidity/x/liquidity/types"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
)

// DefaultConsensusParams defines the default Tendermint consensus params used in
// LiquidityApp testing.
var DefaultConsensusParams = tmtypes.ConsensusParams{
	Block: tmtypes.BlockParams{
		MaxBytes: 200000,
		MaxGas:   2000000,
	},
	Evidence: tmtypes.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
	},
	Validator: tmtypes.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

// SetupOptions defines arguments that are passed into `Simapp` constructor.
type SetupOptions struct {
	LoadLatest bool
	Logger     log.Logger
	DB         *dbm.MemDB
	AppOpts    servertypes.AppOptions
}

func setup(withGenesis bool, invCheckPeriod uint) (*LiquidityApp, GenesisState) {
	db := dbm.NewMemDB()
	appOptions := make(simtestutil.AppOptionsMap, 0)
	appOptions[flags.FlagHome] = DefaultNodeHome
	appOptions[server.FlagInvCheckPeriod] = invCheckPeriod
	app := NewLiquidityApp(log.NewNopLogger(), db, nil, true, appOptions)
	if withGenesis {
		return app, NewDefaultGenesisState()
	}
	return app, GenesisState{}
}

// Setup initializes a new SimApp. A Nop logger is set in SimApp.
func Setup(t *testing.T, isCheckTx bool) *LiquidityApp {
	t.Helper()

	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)

	// create validator set with single validator
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, mathsdk.NewInt(100000000000000))),
	}

	app := SetupWithGenesisValSet(t, valSet, []authtypes.GenesisAccount{acc}, balance)
	app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height: app.LastBlockHeight() + 1,
		Hash:   app.LastCommitID().Hash,
	})
	return app
}

// SetupWithGenesisValSet initializes a new SimApp with a validator set and genesis accounts
// that also act as delegators. For simplicity, each validator is bonded with a delegation
// of one consensus engine unit in the default token of the simapp from first genesis
// account. A Nop logger is set in SimApp.
func SetupWithGenesisValSet(t *testing.T, valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount, balances ...banktypes.Balance) *LiquidityApp {
	t.Helper()

	app, genesisState := setup(true, 5)
	genesisState, err := simtestutil.GenesisStateWithValSet(app.AppCodec(), genesisState, valSet, genAccs, balances...)
	require.NoError(t, err)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)

	// init chain will set the validator set and initialize the genesis accounts
	app.InitChain(
		&abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: simtestutil.DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

	// commit genesis changes
	app.Commit()

	return app
}

type GenerateAccountStrategy func(int) []sdk.AccAddress

// AddRandomTestAddr creates new account with random address.
func AddRandomTestAddr(app *LiquidityApp, ctx sdk.Context, initCoins sdk.Coins) sdk.AccAddress {
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	SaveAccount(app, ctx, addr, initCoins)
	return addr
}

// AddTestAddrs constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrs(app *LiquidityApp, ctx sdk.Context, accNum int, initCoins sdk.Coins) []sdk.AccAddress {
	testAddrs := simtestutil.CreateIncrementalAccounts(accNum)
	for _, addr := range testAddrs {
		if err := FundAccount(app, ctx, addr, initCoins); err != nil {
			panic(err)
		}
	}
	return testAddrs
}

// permission of minting, create a "faucet" account. (@fdymylja)
func FundAccount(app *LiquidityApp, ctx sdk.Context, addr sdk.AccAddress, amounts sdk.Coins) error {
	if err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, amounts); err != nil {
		return err
	}
	return app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, amounts)
}

// AddTestAddrs constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrsIncremental(app *LiquidityApp, ctx sdk.Context, accNum int, accAmt mathsdk.Int) []sdk.AccAddress {
	return addTestAddrs(app, ctx, accNum, accAmt, simtestutil.CreateIncrementalAccounts)
}

func addTestAddrs(app *LiquidityApp, ctx sdk.Context, accNum int, accAmt mathsdk.Int, strategy GenerateAccountStrategy) []sdk.AccAddress {
	testAddrs := strategy(accNum)
	bondDenom, _ := app.StakingKeeper.BondDenom(ctx)
	initCoins := sdk.NewCoins(sdk.NewCoin(bondDenom, accAmt))

	for _, addr := range testAddrs {
		if err := FundAccount(app, ctx, addr, initCoins); err != nil {
			panic(err)
		}
	}

	return testAddrs
}

// SaveAccount saves the provided account into the simapp with balance based on initCoins.
func SaveAccount(app *LiquidityApp, ctx sdk.Context, addr sdk.AccAddress, initCoins sdk.Coins) {
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	app.AccountKeeper.SetAccount(ctx, acc)
	if initCoins.IsAllPositive() {
		err := FundAccount(app, ctx, addr, initCoins)
		if err != nil {
			panic(err)
		}
	}
}

// NewSimappWithCustomOptions initializes a new SimApp with custom options.
func NewLiquidityAppWithCustomOptions(t *testing.T, isCheckTx bool, options SetupOptions) *LiquidityApp {
	t.Helper()

	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)
	// create validator set with single validator
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, mathsdk.NewInt(100000000000000))),
	}

	app := NewLiquidityApp(options.Logger, options.DB, nil, options.LoadLatest, options.AppOpts)
	genesisState := app.DefaultGenesis()
	genesisState, err = simtestutil.GenesisStateWithValSet(app.AppCodec(), genesisState, valSet, []authtypes.GenesisAccount{acc}, balance)
	require.NoError(t, err)

	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		stateBytes, err := tmjson.MarshalIndent(genesisState, "", " ")
		require.NoError(t, err)
		// Initialize the chain
		app.InitChain(
			&abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simtestutil.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}

func SaveAccountWithFee(app *LiquidityApp, ctx sdk.Context, addr sdk.AccAddress, initCoins sdk.Coins, offerCoin sdk.Coin) {
	SaveAccount(app, ctx, addr, initCoins)
	params := app.LiquidityKeeper.GetParams(ctx)
	offerCoinFee := types.GetOfferCoinFee(offerCoin, params.SwapFeeRate)
	err := FundAccount(app, ctx, addr, sdk.NewCoins(offerCoinFee))
	if err != nil {
		panic(err)
	}
}

func TestAddr(addr string, bech string) (sdk.AccAddress, error) {
	res, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, err
	}
	bechexpected := res.String()
	if bech != bechexpected {
		return nil, fmt.Errorf("bech encoding doesn't match reference")
	}

	bechres, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(bechres, res) {
		return nil, err
	}

	return res, nil
}

// CreateTestInput returns a simapp with custom LiquidityKeeper to avoid
// messing with the hooks.
func CreateTestInput(t *testing.T) (*LiquidityApp, sdk.Context) {
	cdc := codec.NewLegacyAmino()
	types.RegisterLegacyAminoCodec(cdc)
	keeper.BatchLogicInvariantCheckFlag = true

	app := Setup(t, false)
	if app.BaseApp == nil {
		t.Log("Something goes wrong !")
	}

	ctx := app.BaseApp.NewContext(false)

	appCodec := app.AppCodec()

	app.LiquidityKeeper = keeper.NewKeeper(
		appCodec,
		app.GetKey(types.StoreKey),
		app.BankKeeper,
		app.AccountKeeper,
		app.DistrKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName),
	)

	app.DistrKeeper.FeePool.Set(ctx, distrtypes.InitialFeePool())
	app.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	app.StakingKeeper.SetParams(ctx, stakingtypes.DefaultParams())
	return app, ctx
}

func GetRandPoolAmt(r *rand.Rand, minInitDepositAmt mathsdk.Int) (x, y mathsdk.Int) {
	x = GetRandRange(r, int(minInitDepositAmt.Int64()), 100000000000000).MulRaw(int64(math.Pow10(r.Intn(10))))
	y = GetRandRange(r, int(minInitDepositAmt.Int64()), 100000000000000).MulRaw(int64(math.Pow10(r.Intn(10))))
	return
}

func GetRandRange(r *rand.Rand, min, max int) mathsdk.Int {
	return mathsdk.NewInt(int64(r.Intn(max-min) + min))
}

func GetRandomSizeOrders(denomX, denomY string, x, y mathsdk.Int, r *rand.Rand, sizeXToY, sizeYToX int32) (xToY, yToX []*types.MsgSwapWithinBatch) {
	randomSizeXtoY := int(r.Int31n(sizeXToY))
	randomSizeYtoX := int(r.Int31n(sizeYToX))
	return GetRandomOrders(denomX, denomY, x, y, r, randomSizeXtoY, randomSizeYtoX)
}

func GetRandomOrders(denomX, denomY string, x, y mathsdk.Int, r *rand.Rand, sizeXToY, sizeYToX int) (xToY, yToX []*types.MsgSwapWithinBatch) {
	currentPrice := x.ToLegacyDec().Quo(y.ToLegacyDec())

	for len(xToY) < sizeXToY {
		orderPrice := currentPrice.Mul(mathsdk.LegacyNewDecFromIntWithPrec(GetRandRange(r, 991, 1009), 3))
		orderAmt := mathsdk.LegacyZeroDec()
		if r.Intn(2) == 1 {
			orderAmt = x.ToLegacyDec().Mul(mathsdk.LegacyNewDecFromIntWithPrec(GetRandRange(r, 1, 100), 4))
		} else {
			orderAmt = mathsdk.LegacyNewDecFromIntWithPrec(GetRandRange(r, 1000, 10000), 0)
		}
		if orderAmt.Quo(orderPrice).TruncateInt().IsZero() {
			continue
		}
		orderCoin := sdk.NewCoin(denomX, orderAmt.Ceil().TruncateInt())

		xToY = append(xToY, &types.MsgSwapWithinBatch{
			OfferCoin:       orderCoin,
			DemandCoinDenom: denomY,
			OrderPrice:      orderPrice,
		})
	}

	for len(yToX) < sizeYToX {
		orderPrice := currentPrice.Mul(mathsdk.LegacyNewDecFromIntWithPrec(GetRandRange(r, 991, 1009), 3))
		orderAmt := mathsdk.LegacyZeroDec()
		if r.Intn(2) == 1 {
			orderAmt = y.ToLegacyDec().Mul(mathsdk.LegacyNewDecFromIntWithPrec(GetRandRange(r, 1, 100), 4))
		} else {
			orderAmt = mathsdk.LegacyNewDecFromIntWithPrec(GetRandRange(r, 1000, 10000), 0)
		}
		if orderAmt.Mul(orderPrice).TruncateInt().IsZero() {
			continue
		}
		orderCoin := sdk.NewCoin(denomY, orderAmt.Ceil().TruncateInt())

		yToX = append(yToX, &types.MsgSwapWithinBatch{
			OfferCoin:       orderCoin,
			DemandCoinDenom: denomX,
			OrderPrice:      orderPrice,
		})
	}
	return xToY, yToX
}

func TestCreatePool(t *testing.T, simapp *LiquidityApp, ctx sdk.Context, x, y mathsdk.Int, denomX, denomY string, addr sdk.AccAddress) uint64 {
	deposit := sdk.NewCoins(sdk.NewCoin(denomX, x), sdk.NewCoin(denomY, y))
	params := simapp.LiquidityKeeper.GetParams(ctx)
	// set accounts for creator, depositor, withdrawer, balance for deposit
	SaveAccount(simapp, ctx, addr, deposit.Add(params.PoolCreationFee...)) // pool creator
	depositX := simapp.BankKeeper.GetBalance(ctx, addr, denomX)
	depositY := simapp.BankKeeper.GetBalance(ctx, addr, denomY)
	depositBalance := sdk.NewCoins(depositX, depositY)
	require.Equal(t, deposit, depositBalance)

	// create Liquidity pool
	poolTypeID := types.DefaultPoolTypeID
	poolID := simapp.LiquidityKeeper.GetNextPoolID(ctx)
	msg := types.NewMsgCreatePool(addr, poolTypeID, depositBalance)
	_, err := simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.NoError(t, err)

	// verify created liquidity pool
	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolID)
	require.True(t, found)
	require.Equal(t, poolID, pool.Id)
	require.Equal(t, denomX, pool.ReserveCoinDenoms[0])
	require.Equal(t, denomY, pool.ReserveCoinDenoms[1])

	// verify minted pool coin
	poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)
	creatorBalance := simapp.BankKeeper.GetBalance(ctx, addr, pool.PoolCoinDenom)
	require.Equal(t, poolCoin, creatorBalance.Amount)
	return poolID
}

func TestDepositPool(t *testing.T, simapp *LiquidityApp, ctx sdk.Context, x, y mathsdk.Int, addrs []sdk.AccAddress, poolID uint64, withEndblock bool) {
	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolID)
	require.True(t, found)
	denomX, denomY := pool.ReserveCoinDenoms[0], pool.ReserveCoinDenoms[1]
	deposit := sdk.NewCoins(sdk.NewCoin(denomX, x), sdk.NewCoin(denomY, y))

	moduleAccAddress := simapp.AccountKeeper.GetModuleAddress(types.ModuleName)
	moduleAccEscrowAmtX := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, denomX)
	moduleAccEscrowAmtY := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, denomY)
	iterNum := len(addrs)
	for i := 0; i < iterNum; i++ {
		SaveAccount(simapp, ctx, addrs[i], deposit) // pool creator

		depositMsg := types.NewMsgDepositWithinBatch(addrs[i], poolID, deposit)
		_, err := simapp.LiquidityKeeper.DepositWithinBatch(ctx, depositMsg)
		require.NoError(t, err)

		depositorBalanceX := simapp.BankKeeper.GetBalance(ctx, addrs[i], pool.ReserveCoinDenoms[0])
		depositorBalanceY := simapp.BankKeeper.GetBalance(ctx, addrs[i], pool.ReserveCoinDenoms[1])
		require.Equal(t, denomX, depositorBalanceX.Denom)
		require.Equal(t, denomY, depositorBalanceY.Denom)

		// check escrow balance of module account
		moduleAccEscrowAmtX = moduleAccEscrowAmtX.Add(deposit[0])
		moduleAccEscrowAmtY = moduleAccEscrowAmtY.Add(deposit[1])
		moduleAccEscrowAmtXAfter := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, denomX)
		moduleAccEscrowAmtYAfter := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, denomY)
		require.Equal(t, moduleAccEscrowAmtX, moduleAccEscrowAmtXAfter)
		require.Equal(t, moduleAccEscrowAmtY, moduleAccEscrowAmtYAfter)
	}
	batch, bool := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolID)
	require.True(t, bool)

	// endblock
	if withEndblock {
		liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)
		msgs := simapp.LiquidityKeeper.GetAllPoolBatchDepositMsgs(ctx, batch)
		for i := 0; i < iterNum; i++ {
			// verify minted pool coin
			poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)
			depositorPoolCoinBalance := simapp.BankKeeper.GetBalance(ctx, addrs[i], pool.PoolCoinDenom)
			require.NotEqual(t, mathsdk.ZeroInt(), depositorPoolCoinBalance)
			require.NotEqual(t, mathsdk.ZeroInt(), poolCoin)

			require.True(t, msgs[i].Executed)
			require.True(t, msgs[i].Succeeded)
			require.True(t, msgs[i].ToBeDeleted)

			// error balance after endblock
			depositorBalanceX := simapp.BankKeeper.GetBalance(ctx, addrs[i], pool.ReserveCoinDenoms[0])
			depositorBalanceY := simapp.BankKeeper.GetBalance(ctx, addrs[i], pool.ReserveCoinDenoms[1])
			require.Equal(t, denomX, depositorBalanceX.Denom)
			require.Equal(t, denomY, depositorBalanceY.Denom)
		}
	}
}

func TestWithdrawPool(t *testing.T, simapp *LiquidityApp, ctx sdk.Context, poolCoinAmt mathsdk.Int, addrs []sdk.AccAddress, poolID uint64, withEndblock bool) {
	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolID)
	require.True(t, found)
	moduleAccAddress := simapp.AccountKeeper.GetModuleAddress(types.ModuleName)
	moduleAccEscrowAmtPool := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, pool.PoolCoinDenom)

	iterNum := len(addrs)
	for i := 0; i < iterNum; i++ {
		balancePoolCoin := simapp.BankKeeper.GetBalance(ctx, addrs[i], pool.PoolCoinDenom)
		require.True(t, balancePoolCoin.Amount.GTE(poolCoinAmt))

		withdrawCoin := sdk.NewCoin(pool.PoolCoinDenom, poolCoinAmt)
		withdrawMsg := types.NewMsgWithdrawWithinBatch(addrs[i], poolID, withdrawCoin)
		_, err := simapp.LiquidityKeeper.WithdrawWithinBatch(ctx, withdrawMsg)
		require.NoError(t, err)

		moduleAccEscrowAmtPoolAfter := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, pool.PoolCoinDenom)
		moduleAccEscrowAmtPool.Amount = moduleAccEscrowAmtPool.Amount.Add(withdrawMsg.PoolCoin.Amount)
		require.Equal(t, moduleAccEscrowAmtPool, moduleAccEscrowAmtPoolAfter)

		balancePoolCoinAfter := simapp.BankKeeper.GetBalance(ctx, addrs[i], pool.PoolCoinDenom)
		if balancePoolCoin.Amount.Equal(withdrawCoin.Amount) {

		} else {
			require.Equal(t, balancePoolCoin.Sub(withdrawCoin).Amount, balancePoolCoinAfter.Amount)
		}

	}

	if withEndblock {
		poolCoinBefore := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)

		// endblock
		liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

		batch, bool := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolID)
		require.True(t, bool)

		// verify burned pool coin
		poolCoinAfter := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)
		fmt.Println(poolCoinAfter, poolCoinBefore)
		require.True(t, poolCoinAfter.LT(poolCoinBefore))

		for i := 0; i < iterNum; i++ {
			withdrawerBalanceX := simapp.BankKeeper.GetBalance(ctx, addrs[i], pool.ReserveCoinDenoms[0])
			withdrawerBalanceY := simapp.BankKeeper.GetBalance(ctx, addrs[i], pool.ReserveCoinDenoms[1])
			require.True(t, withdrawerBalanceX.IsPositive())
			require.True(t, withdrawerBalanceY.IsPositive())

			withdrawMsgs := simapp.LiquidityKeeper.GetAllPoolBatchWithdrawMsgStates(ctx, batch)
			require.True(t, withdrawMsgs[i].Executed)
			require.True(t, withdrawMsgs[i].Succeeded)
			require.True(t, withdrawMsgs[i].ToBeDeleted)
		}
	}
}

func TestSwapPool(t *testing.T, simapp *LiquidityApp, ctx sdk.Context, offerCoins []sdk.Coin, orderPrices []mathsdk.LegacyDec,
	addrs []sdk.AccAddress, poolID uint64, withEndblock bool) ([]*types.SwapMsgState, types.PoolBatch) {
	if len(offerCoins) != len(orderPrices) || len(orderPrices) != len(addrs) {
		require.True(t, false)
	}

	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolID)
	require.True(t, found)

	moduleAccAddress := simapp.AccountKeeper.GetModuleAddress(types.ModuleName)

	var swapMsgStates []*types.SwapMsgState

	params := simapp.LiquidityKeeper.GetParams(ctx)

	iterNum := len(addrs)
	for i := 0; i < iterNum; i++ {
		moduleAccEscrowAmtPool := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, offerCoins[i].Denom)
		currentBalance := simapp.BankKeeper.GetBalance(ctx, addrs[i], offerCoins[i].Denom)
		if currentBalance.IsLT(offerCoins[i]) {
			SaveAccountWithFee(simapp, ctx, addrs[i], sdk.NewCoins(offerCoins[i]), offerCoins[i])
		}
		var demandCoinDenom string
		if pool.ReserveCoinDenoms[0] == offerCoins[i].Denom {
			demandCoinDenom = pool.ReserveCoinDenoms[1]
		} else if pool.ReserveCoinDenoms[1] == offerCoins[i].Denom {
			demandCoinDenom = pool.ReserveCoinDenoms[0]
		} else {
			require.True(t, false)
		}

		swapMsg := types.NewMsgSwapWithinBatch(addrs[i], poolID, types.DefaultSwapTypeID, offerCoins[i], demandCoinDenom, orderPrices[i], params.SwapFeeRate)
		batchPoolSwapMsg, err := simapp.LiquidityKeeper.SwapWithinBatch(ctx, swapMsg, 0)
		require.NoError(t, err)

		swapMsgStates = append(swapMsgStates, batchPoolSwapMsg)
		moduleAccEscrowAmtPoolAfter := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, offerCoins[i].Denom)
		moduleAccEscrowAmtPool.Amount = moduleAccEscrowAmtPool.Amount.Add(offerCoins[i].Amount).Add(types.GetOfferCoinFee(offerCoins[i], params.SwapFeeRate).Amount)
		require.Equal(t, moduleAccEscrowAmtPool, moduleAccEscrowAmtPoolAfter)

	}
	batch, _ := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolID)

	if withEndblock {
		// endblock
		liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

		batch, found = simapp.LiquidityKeeper.GetPoolBatch(ctx, poolID)
		require.True(t, found)
	}
	return swapMsgStates, batch
}

func GetSwapMsg(t *testing.T, simapp *LiquidityApp, ctx sdk.Context, offerCoins []sdk.Coin, orderPrices []mathsdk.LegacyDec,
	addrs []sdk.AccAddress, poolID uint64) []*types.MsgSwapWithinBatch {
	if len(offerCoins) != len(orderPrices) || len(orderPrices) != len(addrs) {
		require.True(t, false)
	}

	var msgs []*types.MsgSwapWithinBatch
	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolID)
	require.True(t, found)

	params := simapp.LiquidityKeeper.GetParams(ctx)

	iterNum := len(addrs)
	for i := 0; i < iterNum; i++ {
		currentBalance := simapp.BankKeeper.GetBalance(ctx, addrs[i], offerCoins[i].Denom)
		if currentBalance.IsLT(offerCoins[i]) {
			SaveAccountWithFee(simapp, ctx, addrs[i], sdk.NewCoins(offerCoins[i]), offerCoins[i])
		}
		var demandCoinDenom string
		if pool.ReserveCoinDenoms[0] == offerCoins[i].Denom {
			demandCoinDenom = pool.ReserveCoinDenoms[1]
		} else if pool.ReserveCoinDenoms[1] == offerCoins[i].Denom {
			demandCoinDenom = pool.ReserveCoinDenoms[0]
		} else {
			require.True(t, false)
		}

		msgs = append(msgs, types.NewMsgSwapWithinBatch(addrs[i], poolID, types.DefaultSwapTypeID, offerCoins[i], demandCoinDenom, orderPrices[i], params.SwapFeeRate))
	}
	return msgs
}

// NewTestNetworkFixture returns a new simapp AppConstructor for network simulation tests
func NewTestNetworkFixture() network.TestFixture {
	dir, err := os.MkdirTemp("", "liquidityapp")
	if err != nil {
		panic(fmt.Sprintf("failed creating temporary directory: %v", err))
	}
	defer os.RemoveAll(dir)

	app := NewLiquidityApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, simtestutil.NewAppOptionsWithFlagHome(dir))

	appCtr := func(val network.ValidatorI) servertypes.Application {
		return NewLiquidityApp(
			val.GetCtx().Logger, dbm.NewMemDB(), nil, true,
			simtestutil.NewAppOptionsWithFlagHome(val.GetCtx().Config.RootDir),
			bam.SetPruning(pruningtypes.NewPruningOptionsFromString(val.GetAppConfig().Pruning)),
			bam.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
			bam.SetChainID(val.GetCtx().Viper.GetString(flags.FlagChainID)),
		)
	}

	return network.TestFixture{
		AppConstructor: appCtr,
		GenesisState:   app.DefaultGenesis(),
		EncodingConfig: testutil.TestEncodingConfig{
			InterfaceRegistry: app.InterfaceRegistry(),
			Codec:             app.AppCodec(),
			TxConfig:          app.TxConfig(),
			Amino:             app.LegacyAmino(),
		},
	}
}

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(o string) interface{} {
	return nil
}
