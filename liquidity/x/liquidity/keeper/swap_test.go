package keeper_test

import (
	"math/rand"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/Victor118/liquidity/app"
	"github.com/Victor118/liquidity/x/liquidity"
	"github.com/Victor118/liquidity/x/liquidity/types"
)

func TestSimulationSwapExecutionFindEdgeCase(t *testing.T) {
	for seed := int64(0); seed < 20; seed++ {
		r := rand.New(rand.NewSource(seed))

		simapp, ctx := createTestInput(t)
		params := simapp.LiquidityKeeper.GetParams(ctx)

		// define test denom X, Y for Liquidity Pool
		denomX := "denomX"
		denomY := "denomY"
		denomX, denomY = types.AlphabeticalDenomPair(denomX, denomY)

		// get random X, Y amount for create pool
		param := simapp.LiquidityKeeper.GetParams(ctx)
		X, Y := app.GetRandPoolAmt(r, param.MinInitDepositAmount)
		deposit := sdk.NewCoins(sdk.NewCoin(denomX, X), sdk.NewCoin(denomY, Y))

		// set pool creator account, balance for deposit
		addrs := app.AddTestAddrs(simapp, ctx, 3, params.PoolCreationFee)
		app.SaveAccount(simapp, ctx, addrs[0], deposit) // pool creator
		depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomX)
		depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomY)
		depositBalance := sdk.NewCoins(depositA, depositB)
		require.Equal(t, deposit, depositBalance)

		// create Liquidity pool
		poolTypeID := types.DefaultPoolTypeID
		msg := types.NewMsgCreatePool(addrs[0], poolTypeID, depositBalance)
		_, err := simapp.LiquidityKeeper.CreatePool(ctx, msg)
		require.NoError(t, err)

		for i := 0; i < 20; i++ {
			ctx = ctx.WithBlockHeight(int64(i))
			testSwapEdgeCases(t, r, simapp, ctx, X, Y, depositBalance, addrs)
		}
	}
}

func TestSwapExecution(t *testing.T) {
	for seed := int64(0); seed < 50; seed++ {
		s := rand.NewSource(seed)
		r := rand.New(s)
		simapp, ctx := createTestInput(t)
		simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
		params := simapp.LiquidityKeeper.GetParams(ctx)

		// define test denom X, Y for Liquidity Pool
		denomX := "denomX"
		denomY := "denomY"
		denomX, denomY = types.AlphabeticalDenomPair(denomX, denomY)

		// get random X, Y amount for create pool
		X, Y := app.GetRandPoolAmt(r, params.MinInitDepositAmount)
		deposit := sdk.NewCoins(sdk.NewCoin(denomX, X), sdk.NewCoin(denomY, Y))

		// set pool creator account, balance for deposit
		addrs := app.AddTestAddrs(simapp, ctx, 3, params.PoolCreationFee)
		app.SaveAccount(simapp, ctx, addrs[0], deposit) // pool creator
		depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomX)
		depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomY)
		depositBalance := sdk.NewCoins(depositA, depositB)
		require.Equal(t, deposit, depositBalance)

		// create Liquidity pool
		poolTypeID := types.DefaultPoolTypeID
		msg := types.NewMsgCreatePool(addrs[0], poolTypeID, depositBalance)
		_, err := simapp.LiquidityKeeper.CreatePool(ctx, msg)
		require.NoError(t, err)

		// verify created liquidity pool
		pools := simapp.LiquidityKeeper.GetAllPools(ctx)
		poolID := pools[0].Id
		require.Equal(t, 1, len(pools))
		require.Equal(t, uint64(1), poolID)
		require.Equal(t, denomX, pools[0].ReserveCoinDenoms[0])
		require.Equal(t, denomY, pools[0].ReserveCoinDenoms[1])

		// verify minted pool coin
		poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pools[0])
		creatorBalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], pools[0].PoolCoinDenom)
		require.Equal(t, poolCoin, creatorBalance.Amount)

		var xToY []*types.MsgSwapWithinBatch // buying Y from X
		var yToX []*types.MsgSwapWithinBatch // selling Y for X

		// make random orders, set buyer, seller accounts for the orders
		xToY, yToX = app.GetRandomSizeOrders(denomX, denomY, X, Y, r, 250, 250)
		buyerAddrs := app.AddTestAddrsIncremental(simapp, ctx, len(xToY), math.NewInt(0))
		sellerAddrs := app.AddTestAddrsIncremental(simapp, ctx, len(yToX), math.NewInt(0))

		for i, msg := range xToY {
			app.SaveAccountWithFee(simapp, ctx, buyerAddrs[i], sdk.NewCoins(msg.OfferCoin), msg.OfferCoin)
			msg.SwapRequesterAddress = buyerAddrs[i].String()
			msg.PoolId = poolID
			msg.OfferCoinFee = types.GetOfferCoinFee(msg.OfferCoin, params.SwapFeeRate)
		}
		for i, msg := range yToX {
			app.SaveAccountWithFee(simapp, ctx, sellerAddrs[i], sdk.NewCoins(msg.OfferCoin), msg.OfferCoin)
			msg.SwapRequesterAddress = sellerAddrs[i].String()
			msg.PoolId = poolID
			msg.OfferCoinFee = types.GetOfferCoinFee(msg.OfferCoin, params.SwapFeeRate)
		}

		// begin block, delete and init pool batch
		liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

		// handle msgs, set order msgs to batch
		for _, msg := range xToY {
			_, err := simapp.LiquidityKeeper.SwapWithinBatch(ctx, msg, 0)
			require.NoError(t, err)
		}
		for _, msg := range yToX {
			_, err := simapp.LiquidityKeeper.SwapWithinBatch(ctx, msg, 0)
			require.NoError(t, err)
		}

		// verify pool batch
		liquidityPoolBatch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolID)
		require.True(t, found)
		require.NotNil(t, liquidityPoolBatch)

		// end block, swap execution
		liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)
	}
}

func testSwapEdgeCases(t *testing.T, r *rand.Rand, simapp *app.LiquidityApp, ctx sdk.Context, X, Y math.Int, depositBalance sdk.Coins, addrs []sdk.AccAddress) {
	//simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)

	denomX := depositBalance[0].Denom
	denomY := depositBalance[1].Denom

	// verify created liquidity pool
	pools := simapp.LiquidityKeeper.GetAllPools(ctx)
	poolID := pools[0].Id
	require.Equal(t, 1, len(pools))
	require.Equal(t, uint64(1), poolID)
	require.Equal(t, denomX, pools[0].ReserveCoinDenoms[0])
	require.Equal(t, denomY, pools[0].ReserveCoinDenoms[1])

	// verify minted pool coin
	poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pools[0])
	creatorBalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], pools[0].PoolCoinDenom)
	require.Equal(t, poolCoin, creatorBalance.Amount)

	var xToY []*types.MsgSwapWithinBatch // buying Y from X
	var yToX []*types.MsgSwapWithinBatch // selling Y for X

	batch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolID)
	require.True(t, found)

	remainingSwapMsgs := simapp.LiquidityKeeper.GetAllNotProcessedPoolBatchSwapMsgStates(ctx, batch)
	if ctx.BlockHeight() == 0 || len(remainingSwapMsgs) == 0 {
		// make random orders, set buyer, seller accounts for the orders
		xToY, yToX = app.GetRandomSizeOrders(denomX, denomY, X, Y, r, 100, 100)
		buyerAddrs := app.AddTestAddrsIncremental(simapp, ctx, len(xToY), math.NewInt(0))
		sellerAddrs := app.AddTestAddrsIncremental(simapp, ctx, len(yToX), math.NewInt(0))

		for i, msg := range xToY {
			app.SaveAccountWithFee(simapp, ctx, buyerAddrs[i], sdk.NewCoins(msg.OfferCoin), msg.OfferCoin)
			msg.SwapRequesterAddress = buyerAddrs[i].String()
			msg.PoolId = poolID
			msg.OfferCoinFee = types.GetOfferCoinFee(msg.OfferCoin, params.SwapFeeRate)
		}
		for i, msg := range yToX {
			app.SaveAccountWithFee(simapp, ctx, sellerAddrs[i], sdk.NewCoins(msg.OfferCoin), msg.OfferCoin)
			msg.SwapRequesterAddress = sellerAddrs[i].String()
			msg.PoolId = poolID
			msg.OfferCoinFee = types.GetOfferCoinFee(msg.OfferCoin, params.SwapFeeRate)
		}
	}

	// begin block, delete and init pool batch
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	// handle msgs, set order msgs to batch
	for _, msg := range xToY {
		_, err := simapp.LiquidityKeeper.SwapWithinBatch(ctx, msg, int64(r.Intn(4)))
		require.NoError(t, err)
	}
	for _, msg := range yToX {
		_, err := simapp.LiquidityKeeper.SwapWithinBatch(ctx, msg, int64(r.Intn(4)))
		require.NoError(t, err)
	}

	// verify pool batch
	liquidityPoolBatch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolID)
	require.True(t, found)
	require.NotNil(t, liquidityPoolBatch)

	// end block, swap execution
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)
}

func TestBadSwapExecution(t *testing.T) {
	r := rand.New(rand.NewSource(0))

	simapp, ctx := app.CreateTestInput(t)
	params := simapp.LiquidityKeeper.GetParams(ctx)
	denomX, denomY := types.AlphabeticalDenomPair("denomX", "denomY")

	// add pool creator account
	X, Y := app.GetRandPoolAmt(r, params.MinInitDepositAmount)
	deposit := sdk.NewCoins(sdk.NewCoin(denomX, X), sdk.NewCoin(denomY, Y))
	creatorAddr := app.AddRandomTestAddr(simapp, ctx, deposit.Add(params.PoolCreationFee...))
	balanceX := simapp.BankKeeper.GetBalance(ctx, creatorAddr, denomX)
	balanceY := simapp.BankKeeper.GetBalance(ctx, creatorAddr, denomY)
	creatorBalance := sdk.NewCoins(balanceX, balanceY)
	require.Equal(t, deposit, creatorBalance)

	// create pool
	createPoolMsg := types.NewMsgCreatePool(creatorAddr, types.DefaultPoolTypeID, creatorBalance)
	_, err := simapp.LiquidityKeeper.CreatePool(ctx, createPoolMsg)
	require.NoError(t, err)

	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	offerCoin := sdk.NewCoin(denomX, math.NewInt(10000))
	offerCoinFee := types.GetOfferCoinFee(offerCoin, params.SwapFeeRate)
	testAddr := app.AddRandomTestAddr(simapp, ctx, sdk.NewCoins(offerCoin.Add(offerCoinFee)))

	currentPrice := X.ToLegacyDec().Quo(Y.ToLegacyDec())
	swapMsg := types.NewMsgSwapWithinBatch(testAddr, 0, types.DefaultSwapTypeID, offerCoin, denomY, currentPrice, params.SwapFeeRate)
	_, err = simapp.LiquidityKeeper.SwapWithinBatch(ctx, swapMsg, 0)
	require.ErrorIs(t, err, types.ErrPoolNotExists)

	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)
}

func TestBalancesAfterSwap(t *testing.T) {
	for price := int64(9800); price < 10000; price++ {
		simapp, ctx := app.CreateTestInput(t)
		params := simapp.LiquidityKeeper.GetParams(ctx)
		denomX, denomY := types.AlphabeticalDenomPair("denomx", "denomy")
		X, Y := math.NewInt(100_000_000), math.NewInt(100_000_000)

		creatorCoins := sdk.NewCoins(sdk.NewCoin(denomX, X), sdk.NewCoin(denomY, Y))
		creatorAddr := app.AddRandomTestAddr(simapp, ctx, creatorCoins.Add(params.PoolCreationFee...))

		orderPrice := math.LegacyNewDecWithPrec(price, 4)
		aliceCoin := sdk.NewCoin(denomY, math.NewInt(10_000_000))
		aliceAddr := app.AddRandomTestAddr(simapp, ctx, sdk.NewCoins(aliceCoin))

		pool, err := simapp.LiquidityKeeper.CreatePool(ctx, types.NewMsgCreatePool(creatorAddr, types.DefaultPoolTypeID, creatorCoins))
		require.NoError(t, err)

		liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

		offerAmt := aliceCoin.Amount.ToLegacyDec().Quo(math.LegacyMustNewDecFromStr("1.0015")).TruncateInt()
		offerCoin := sdk.NewCoin(denomY, offerAmt)

		_, err = simapp.LiquidityKeeper.SwapWithinBatch(ctx, types.NewMsgSwapWithinBatch(
			aliceAddr, pool.Id, types.DefaultSwapTypeID, offerCoin, denomX, orderPrice, params.SwapFeeRate), 0)
		require.NoError(t, err)

		liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

		deltaX := simapp.BankKeeper.GetBalance(ctx, aliceAddr, denomX).Amount
		deltaY := simapp.BankKeeper.GetBalance(ctx, aliceAddr, denomY).Amount.Sub(aliceCoin.Amount)
		require.Truef(t, !deltaX.IsNegative(), "deltaX should not be negative: %s", deltaX)
		require.Truef(t, deltaY.IsNegative(), "deltaY should be negative: %s", deltaY)

		deltaXWithoutFee := deltaX.ToLegacyDec().Quo(math.LegacyMustNewDecFromStr("0.9985"))
		deltaYWithoutFee := deltaY.ToLegacyDec().Quo(math.LegacyMustNewDecFromStr("1.0015"))
		effectivePrice := deltaXWithoutFee.Quo(deltaYWithoutFee.Neg())
		priceDiffRatio := orderPrice.Sub(effectivePrice).Abs().Quo(orderPrice)
		require.Truef(t, priceDiffRatio.LT(math.LegacyMustNewDecFromStr("0.01")), "effectivePrice differs too much from orderPrice")
	}
}

func TestRefundEscrow(t *testing.T) {
	for seed := int64(0); seed < 100; seed++ {
		r := rand.New(rand.NewSource(seed))

		X := math.NewInt(1_000_000)
		Y := app.GetRandRange(r, 10_000_000_000_000_000, 1_000_000_000_000_000_000)

		simapp, ctx := createTestInput(t)
		params := simapp.LiquidityKeeper.GetParams(ctx)

		addr := app.AddRandomTestAddr(simapp, ctx, sdk.NewCoins())

		pool, err := createPool(simapp, ctx, X, Y, DenomX, DenomY)
		require.NoError(t, err)

		for i := 0; i < 100; i++ {
			poolBalances := simapp.BankKeeper.GetAllBalances(ctx, pool.GetReserveAccount())
			RX := poolBalances.AmountOf(DenomX)
			RY := poolBalances.AmountOf(DenomY)
			poolPrice := RX.ToLegacyDec().Quo(RY.ToLegacyDec())

			offerAmt := RY.ToLegacyDec().Mul(math.LegacyNewDecFromIntWithPrec(app.GetRandRange(r, 1, 100_000_000_000_000_000), math.LegacyPrecision)) // RY * (0, 0.1)
			offerAmtWithFee := offerAmt.Quo(math.LegacyOneDec().Add(params.SwapFeeRate.QuoInt64(2))).TruncateInt()                                    // offerAmt / (1 + swapFeeRate/2)
			orderPrice := poolPrice.Mul(math.LegacyNewDecFromIntWithPrec(app.GetRandRange(r, 1, 1_000_000_000_000_000_000), math.LegacyPrecision))    // poolPrice * (0, 1)

			app.SaveAccount(simapp, ctx, addr, sdk.NewCoins(sdk.NewCoin(DenomY, offerAmt.Ceil().TruncateInt())))

			liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

			_, err := simapp.LiquidityKeeper.SwapWithinBatch(ctx, types.NewMsgSwapWithinBatch(
				addr, pool.Id, types.DefaultSwapTypeID, sdk.NewCoin(DenomY, offerAmtWithFee), DenomX, orderPrice, params.SwapFeeRate), 0)
			require.NoError(t, err)

			liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)
		}

		require.True(t, simapp.BankKeeper.GetAllBalances(ctx, simapp.AccountKeeper.GetModuleAddress(types.ModuleName)).IsZero(), "there must be no remaining coin escrow")
	}
}

func TestSwapWithDepletedPool(t *testing.T) {
	simapp, ctx, pool, creatorAddr, err := createTestPool(t, sdk.NewInt64Coin(DenomX, 1000000), sdk.NewInt64Coin(DenomY, 1000000))
	require.NoError(t, err)
	params := simapp.LiquidityKeeper.GetParams(ctx)

	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
	pc := simapp.BankKeeper.GetBalance(ctx, creatorAddr, pool.PoolCoinDenom)
	_, err = simapp.LiquidityKeeper.WithdrawWithinBatch(ctx, types.NewMsgWithdrawWithinBatch(creatorAddr, pool.Id, pc))
	require.NoError(t, err)
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	addr := app.AddRandomTestAddr(simapp, ctx, sdk.NewCoins(sdk.NewInt64Coin(DenomX, 100000)))
	offerCoin := sdk.NewInt64Coin(DenomX, 10000)
	orderPrice := math.LegacyMustNewDecFromStr("1.0")
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
	_, err = simapp.LiquidityKeeper.SwapWithinBatch(
		ctx,
		types.NewMsgSwapWithinBatch(addr, pool.Id, types.DefaultSwapTypeID, offerCoin, DenomY, orderPrice, params.SwapFeeRate),
		0)
	require.ErrorIs(t, err, types.ErrDepletedPool)
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)
}

func TestDirectSwap(t *testing.T) {
	simapp, ctx := createTestInput(t)
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)
	t.Logf("%v", params)
	poolTypeID := types.DefaultPoolTypeID
	addrs := app.AddTestAddrs(simapp, ctx, 3, params.PoolCreationFee)

	denomA := "uETH"
	denomB := "uUSD"
	denomA, denomB = types.AlphabeticalDenomPair(denomA, denomB)

	deposit := sdk.NewCoins(sdk.NewCoin(denomA, math.NewInt(102*1000000)), sdk.NewCoin(denomB, math.NewInt(2000*1000000)))
	app.SaveAccount(simapp, ctx, addrs[0], deposit)

	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomA)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomB)
	t.Logf("Initial user account %v , %v", depositA, depositB)
	depositBalance := sdk.NewCoins(depositA.SubAmount(math.NewInt(2*1000000)), depositB)

	msg := types.NewMsgCreatePool(addrs[0], poolTypeID, depositBalance)
	pool, err := simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.NoError(t, err)
	reserveCoins := simapp.LiquidityKeeper.GetReserveCoins(ctx, pool)
	initialX := reserveCoins[0].Amount.ToLegacyDec()
	initialY := reserveCoins[1].Amount.ToLegacyDec()
	currentPoolPrice := initialX.Quo(initialY)
	slippageAmount := currentPoolPrice.Mul(math.LegacyNewDecWithPrec(1, 2))
	maxOrderPrice := currentPoolPrice.Add(slippageAmount)
	receivedAmount, error := simapp.LiquidityKeeper.DirectSwapExecution(ctx, pool.Id, sdk.NewCoin(denomA, math.NewInt(1*1000000)), denomB, maxOrderPrice, addrs[0])
	require.NoError(t, error)
	reserveCoins = simapp.LiquidityKeeper.GetReserveCoins(ctx, pool)
	newX := reserveCoins[0].Amount.ToLegacyDec()
	newY := reserveCoins[1].Amount.ToLegacyDec()
	require.Equal(t, newX, math.LegacyNewDec(101000000)) // initial balance + offer coin + offer coin fee

	balanceB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomB)
	require.Equal(t, receivedAmount, balanceB.Amount)
	t.Logf("New Balance of uUSD : %v", balanceB)
	require.Equal(t, newY, initialY.Sub(balanceB.Amount.ToLegacyDec()))
}

func TestDirectSwapWithBuilders(t *testing.T) {
	simapp, ctx := createTestInput(t)
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)
	builderAddr1 := types.WeightedAddress{
		Address: "cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg36er2cp",
		Weight:  math.LegacyNewDec(1),
	}
	buildersAddresses := []types.WeightedAddress{builderAddr1}
	params.BuildersAddresses = buildersAddresses
	simapp.LiquidityKeeper.SetParams(ctx, params)
	t.Logf("%v", params)
	poolTypeID := types.DefaultPoolTypeID
	addrs := app.AddTestAddrs(simapp, ctx, 3, params.PoolCreationFee)

	denomA := "uETH"
	denomB := "uUSD"
	denomA, denomB = types.AlphabeticalDenomPair(denomA, denomB)

	deposit := sdk.NewCoins(sdk.NewCoin(denomA, math.NewInt(102*1000000)), sdk.NewCoin(denomB, math.NewInt(2000*1000000)))
	app.SaveAccount(simapp, ctx, addrs[0], deposit)

	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomA)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomB)
	t.Logf("Initial user account %v , %v", depositA, depositB)
	depositBalance := sdk.NewCoins(depositA.SubAmount(math.NewInt(2*1000000)), depositB)

	msg := types.NewMsgCreatePool(addrs[0], poolTypeID, depositBalance)
	pool, err := simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.NoError(t, err)
	reserveCoins := simapp.LiquidityKeeper.GetReserveCoins(ctx, pool)
	initialX := reserveCoins[0].Amount.ToLegacyDec()
	initialY := reserveCoins[1].Amount.ToLegacyDec()
	currentPoolPrice := initialX.Quo(initialY)
	slippageAmount := currentPoolPrice.Mul(math.LegacyNewDecWithPrec(1, 2))
	maxOrderPrice := currentPoolPrice.Add(slippageAmount)
	receivedAmount, error := simapp.LiquidityKeeper.DirectSwapExecution(ctx, pool.Id, sdk.NewCoin(denomA, math.NewInt(1*1000000)), denomB, maxOrderPrice, addrs[0])
	require.NoError(t, error)
	reserveCoins = simapp.LiquidityKeeper.GetReserveCoins(ctx, pool)
	newX := reserveCoins[0].Amount.ToLegacyDec()
	newY := reserveCoins[1].Amount.ToLegacyDec()
	builderBalanceB := simapp.BankKeeper.GetBalance(ctx, sdk.MustAccAddressFromBech32("cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg36er2cp"), denomB)
	t.Logf("builderBalanceB %v", builderBalanceB)
	require.True(t, builderBalanceB.Amount.Equal(math.NewInt(11881)))
	require.True(t, builderBalanceB.Amount.GT(math.NewInt(0)))
	require.Equal(t, newX, math.LegacyNewDec(101000000)) // initial balance + offer coin

	balanceB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomB)
	require.Equal(t, receivedAmount, balanceB.Amount)
	t.Logf("New Balance of uUSD : %v", balanceB)
	require.Equal(t, newY, initialY.Sub(balanceB.Amount.ToLegacyDec()).Sub(builderBalanceB.Amount.ToLegacyDec()))
}

func TestDirectSwapKoMaxOrderRatio(t *testing.T) {
	simapp, ctx := createTestInput(t)
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)
	t.Logf("%v", params)
	poolTypeID := types.DefaultPoolTypeID
	addrs := app.AddTestAddrs(simapp, ctx, 3, params.PoolCreationFee)

	denomA := "uETH"
	denomB := "uUSD"
	denomA, denomB = types.AlphabeticalDenomPair(denomA, denomB)

	deposit := sdk.NewCoins(sdk.NewCoin(denomA, math.NewInt(102*1000000)), sdk.NewCoin(denomB, math.NewInt(2000*1000000)))
	app.SaveAccount(simapp, ctx, addrs[0], deposit)

	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomA)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomB)
	t.Logf("Initial user account %v , %v", depositA, depositB)
	depositBalance := sdk.NewCoins(depositA.SubAmount(math.NewInt(2*1000000)), depositB)

	msg := types.NewMsgCreatePool(addrs[0], poolTypeID, depositBalance)
	pool, err := simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.NoError(t, err)
	reserveCoins := simapp.LiquidityKeeper.GetReserveCoins(ctx, pool)
	initialX := reserveCoins[0].Amount.ToLegacyDec()
	initialY := reserveCoins[1].Amount.ToLegacyDec()
	currentPoolPrice := initialX.Quo(initialY)
	slippageAmount := currentPoolPrice.Mul(math.LegacyNewDecWithPrec(1, 2))
	maxOrderPrice := currentPoolPrice.Add(slippageAmount)
	_, error := simapp.LiquidityKeeper.DirectSwapExecution(ctx, pool.Id, sdk.NewCoin(denomA, math.NewInt(70*1000000)), denomB, maxOrderPrice, addrs[0])
	require.Error(t, error)
	require.Equal(t, error.Error(), "can not exceed max order ratio of reserve coins that can be ordered at a order")

}

func TestDirectSwapKoMaxSlippageReached(t *testing.T) {
	simapp, ctx := createTestInput(t)
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)
	t.Logf("%v", params)
	poolTypeID := types.DefaultPoolTypeID
	addrs := app.AddTestAddrs(simapp, ctx, 3, params.PoolCreationFee)

	denomA := "uETH"
	denomB := "uUSD"
	denomA, denomB = types.AlphabeticalDenomPair(denomA, denomB)

	deposit := sdk.NewCoins(sdk.NewCoin(denomA, math.NewInt(102*1000000)), sdk.NewCoin(denomB, math.NewInt(2000*1000000)))
	app.SaveAccount(simapp, ctx, addrs[0], deposit)

	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomA)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomB)
	t.Logf("Initial user account %v , %v", depositA, depositB)
	depositBalance := sdk.NewCoins(depositA.SubAmount(math.NewInt(2*1000000)), depositB)

	msg := types.NewMsgCreatePool(addrs[0], poolTypeID, depositBalance)
	pool, err := simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.NoError(t, err)
	reserveCoins := simapp.LiquidityKeeper.GetReserveCoins(ctx, pool)
	initialX := reserveCoins[0].Amount.ToLegacyDec()
	initialY := reserveCoins[1].Amount.ToLegacyDec()
	currentPoolPrice := initialX.Quo(initialY)
	//using currentPoolPrice as OfferPrice should provide the error "Max slippage reached"
	_, error := simapp.LiquidityKeeper.DirectSwapExecution(ctx, pool.Id, sdk.NewCoin(denomA, math.NewInt(1*1000000)), denomB, currentPoolPrice, addrs[0])
	require.Error(t, error)
	require.Equal(t, error.Error(), "Max slippage reached")

}

func createPool(simapp *app.LiquidityApp, ctx sdk.Context, X, Y math.Int, denomX, denomY string) (types.Pool, error) {
	params := simapp.LiquidityKeeper.GetParams(ctx)

	coins := sdk.NewCoins(sdk.NewCoin(denomX, X), sdk.NewCoin(denomY, Y))
	addr := app.AddRandomTestAddr(simapp, ctx, coins.Add(params.PoolCreationFee...))

	return simapp.LiquidityKeeper.CreatePool(ctx, types.NewMsgCreatePool(addr, types.DefaultPoolTypeID, coins))
}
