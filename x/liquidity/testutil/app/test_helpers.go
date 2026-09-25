package app

// DONTCOVER

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"testing"

	chihuahuaapp "github.com/ChihuahuaChain/chihuahua/app"

	mathsdk "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/ChihuahuaChain/chihuahua/x/liquidity"
	"github.com/ChihuahuaChain/chihuahua/x/liquidity/keeper"
	"github.com/ChihuahuaChain/chihuahua/x/liquidity/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
)

// Setup initializes a new chihuahua app for the liquidity tests.
func Setup(t *testing.T, _ bool) *LiquidityApp {
	t.Helper()
	return chihuahuaapp.Setup(t)
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

// AllowPoolCreator adds addr to the pool_permissioned_creator_addresses param.
// The upstream tests predate the pool creator whitelist, so every account
// created through the helpers below is allowed to create pools.
func AllowPoolCreator(app *LiquidityApp, ctx sdk.Context, addr sdk.AccAddress) {
	params := app.LiquidityKeeper.GetParams(ctx)
	for _, a := range params.PoolPermissionedCreatorAddresses {
		if a == addr.String() {
			return
		}
	}
	params.PoolPermissionedCreatorAddresses = append(params.PoolPermissionedCreatorAddresses, addr.String())
	if err := app.LiquidityKeeper.SetParams(ctx, params); err != nil {
		panic(err)
	}
}

// permission of minting, create a "faucet" account. (@fdymylja)
func FundAccount(app *LiquidityApp, ctx sdk.Context, addr sdk.AccAddress, amounts sdk.Coins) error {
	AllowPoolCreator(app, ctx, addr)
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
	AllowPoolCreator(app, ctx, addr)
	if initCoins.IsAllPositive() {
		err := FundAccount(app, ctx, addr, initCoins)
		if err != nil {
			panic(err)
		}
	}
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

// NewTestNetworkFixture returns the chihuahua app fixture for network tests
func NewTestNetworkFixture() network.TestFixture {
	return chihuahuaapp.NewTestNetworkFixture()
}
