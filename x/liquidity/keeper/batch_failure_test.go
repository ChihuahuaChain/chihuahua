package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/ChihuahuaChain/chihuahua/x/liquidity"
	"github.com/ChihuahuaChain/chihuahua/x/liquidity/testutil/app"
	"github.com/ChihuahuaChain/chihuahua/x/liquidity/types"
)

func setupSwapBatch(t *testing.T) (*app.LiquidityApp, sdk.Context, types.Pool, sdk.AccAddress, sdk.Coin, sdk.Coin) {
	simapp, ctx := createTestInput(t)
	params := types.DefaultParams()
	simapp.LiquidityKeeper.SetParams(ctx, params)

	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)
	addrs := app.AddTestAddrsIncremental(simapp, ctx, 3, math.NewInt(10000))
	poolID := app.TestCreatePool(t, simapp, ctx, math.NewInt(1000000000), math.NewInt(1000000000), denomX, denomY, addrs[0])
	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolID)
	require.True(t, found)

	offerCoin := sdk.NewCoin(denomX, math.NewInt(10000))
	app.TestSwapPool(t, simapp, ctx, []sdk.Coin{offerCoin}, []math.LegacyDec{math.LegacyNewDecWithPrec(11, 1)}, addrs[1:2], poolID, false)
	return simapp, ctx, pool, addrs[1], offerCoin, types.GetOfferCoinFee(offerCoin, params.SwapFeeRate)
}

// A swap batch that cannot be executed is refunded instead of halting the chain.
func TestSwapExecutionFailureRefundsTheBatch(t *testing.T) {
	simapp, ctx, pool, requester, offerCoin, offerCoinFee := setupSwapBatch(t)

	// a depleted pool makes SwapExecution fail
	sink := app.AddTestAddrsIncremental(simapp, ctx, 1, math.ZeroInt())[0]
	reserve := simapp.LiquidityKeeper.GetReserveCoins(ctx, pool)
	require.NoError(t, simapp.BankKeeper.SendCoins(ctx, pool.GetReserveAccount(), sink, reserve))
	require.True(t, simapp.LiquidityKeeper.IsDepletedPool(ctx, pool))

	require.NotPanics(t, func() { liquidity.EndBlocker(ctx, simapp.LiquidityKeeper) })

	balance := simapp.BankKeeper.GetBalance(ctx, requester, offerCoin.Denom)
	require.Equal(t, offerCoin.Add(offerCoinFee).Amount, balance.Amount)

	states := simapp.LiquidityKeeper.GetAllPoolBatchSwapMsgStatesAsPointer(ctx, mustBatch(t, simapp, ctx, pool.Id))
	require.Len(t, states, 1)
	require.True(t, states[0].Executed)
	require.False(t, states[0].Succeeded)
	require.True(t, states[0].ToBeDeleted)

	batch := mustBatch(t, simapp, ctx, pool.Id)
	require.True(t, batch.Executed)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
	require.Empty(t, simapp.LiquidityKeeper.GetAllPoolBatchSwapMsgStatesAsPointer(ctx, mustBatch(t, simapp, ctx, pool.Id)))
}

// When a transfer of a matched swap fails, none of the transfers of the batch
// are applied: the requester doesn't receive the demand coin for free.
func TestSwapTransferFailureLeavesNoPartialTransfers(t *testing.T) {
	simapp, ctx, pool, requester, offerCoin, offerCoinFee := setupSwapBatch(t)

	// take the escrowed offer coin away: the escrow -> reserve transfer fails
	sink := app.AddTestAddrsIncremental(simapp, ctx, 1, math.ZeroInt())[0]
	escrowed := sdk.NewCoins(offerCoin.Add(offerCoinFee))
	require.NoError(t, simapp.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sink, escrowed))

	reserveBefore := simapp.LiquidityKeeper.GetReserveCoins(ctx, pool)
	demandBefore := simapp.BankKeeper.GetBalance(ctx, requester, pool.ReserveCoinDenoms[1])

	require.NotPanics(t, func() { liquidity.EndBlocker(ctx, simapp.LiquidityKeeper) })

	require.Equal(t, reserveBefore, simapp.LiquidityKeeper.GetReserveCoins(ctx, pool))
	require.Equal(t, demandBefore, simapp.BankKeeper.GetBalance(ctx, requester, pool.ReserveCoinDenoms[1]))

	// the refund fails too (empty escrow): the swap is left as it was
	states := simapp.LiquidityKeeper.GetAllPoolBatchSwapMsgStatesAsPointer(ctx, mustBatch(t, simapp, ctx, pool.Id))
	require.Len(t, states, 1)
	require.False(t, states[0].Executed)
	require.False(t, states[0].Succeeded)
}

// A deposit whose execution fails after the pool coins are minted leaves no
// minted pool coins and no partial transfers behind, and doesn't halt the chain.
func TestDepositFailureLeavesNoPartialTransfers(t *testing.T) {
	simapp, ctx := createTestInput(t)
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)
	X := math.NewInt(1000000000)
	Y := math.NewInt(500000000)
	addrs := app.AddTestAddrsIncremental(simapp, ctx, 3, math.NewInt(10000))
	poolID := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])
	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolID)
	require.True(t, found)

	app.TestDepositPool(t, simapp, ctx, X, Y, addrs[1:2], poolID, false)

	// take the escrowed deposit away: the escrow -> reserve transfer fails
	sink := app.AddTestAddrsIncremental(simapp, ctx, 1, math.ZeroInt())[0]
	escrowed := sdk.NewCoins(sdk.NewCoin(denomX, X), sdk.NewCoin(denomY, Y))
	require.NoError(t, simapp.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sink, escrowed))

	poolCoinSupply := simapp.BankKeeper.GetSupply(ctx, pool.PoolCoinDenom)
	reserveBefore := simapp.LiquidityKeeper.GetReserveCoins(ctx, pool)

	require.NotPanics(t, func() { liquidity.EndBlocker(ctx, simapp.LiquidityKeeper) })

	require.Equal(t, poolCoinSupply, simapp.BankKeeper.GetSupply(ctx, pool.PoolCoinDenom))
	require.Equal(t, reserveBefore, simapp.LiquidityKeeper.GetReserveCoins(ctx, pool))
	require.True(t, simapp.BankKeeper.GetBalance(ctx, addrs[1], pool.PoolCoinDenom).IsZero())

	states := simapp.LiquidityKeeper.GetAllPoolBatchDepositMsgs(ctx, mustBatch(t, simapp, ctx, poolID))
	require.Len(t, states, 1)
	require.True(t, states[0].Executed)
	require.False(t, states[0].Succeeded)
}

func mustBatch(t *testing.T, simapp *app.LiquidityApp, ctx sdk.Context, poolID uint64) types.PoolBatch {
	batch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolID)
	require.True(t, found)
	return batch
}
