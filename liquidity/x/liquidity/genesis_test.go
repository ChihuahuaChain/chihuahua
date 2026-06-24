package liquidity_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/Victor118/liquidity/app"
	"github.com/Victor118/liquidity/x/liquidity"
	"github.com/Victor118/liquidity/x/liquidity/types"
)

func TestGenesisState(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	types.RegisterLegacyAminoCodec(cdc)
	simapp := app.Setup(t, false)

	ctx := simapp.BaseApp.NewContext(false)
	genesis := types.DefaultGenesisState()

	liquidity.InitGenesis(ctx, simapp.LiquidityKeeper, *genesis)

	defaultGenesisExported := liquidity.ExportGenesis(ctx, simapp.LiquidityKeeper)

	require.Equal(t, genesis, defaultGenesisExported)

	// define test denom X, Y for Liquidity Pool
	denomX, denomY := types.AlphabeticalDenomPair("denomX", "denomY")

	X := math.NewInt(1000000000)
	Y := math.NewInt(1000000000)

	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, math.NewInt(10000))
	poolID := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])

	// begin block, init
	app.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolID, true)
	app.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolID, true)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	price, _ := math.LegacyNewDecFromStr("1.1")
	offerCoins := []sdk.Coin{sdk.NewCoin(denomX, math.NewInt(10000))}
	orderPrices := []math.LegacyDec{price}
	orderAddrs := addrs[1:2]
	_, _ = app.TestSwapPool(t, simapp, ctx, offerCoins, orderPrices, orderAddrs, poolID, false)
	_, _ = app.TestSwapPool(t, simapp, ctx, offerCoins, orderPrices, orderAddrs, poolID, false)
	_, _ = app.TestSwapPool(t, simapp, ctx, offerCoins, orderPrices, orderAddrs, poolID, true)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
	_, _ = app.TestSwapPool(t, simapp, ctx, offerCoins, orderPrices, orderAddrs, poolID, true)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	genesisExported := liquidity.ExportGenesis(ctx, simapp.LiquidityKeeper)
	bankGenesisExported := simapp.BankKeeper.ExportGenesis(ctx)

	simapp2 := app.Setup(t, false)

	ctx2 := simapp2.BaseApp.NewContext(false)
	ctx2 = ctx2.WithBlockHeight(1)

	simapp2.BankKeeper.InitGenesis(ctx2, bankGenesisExported)
	liquidity.InitGenesis(ctx2, simapp2.LiquidityKeeper, *genesisExported)
	simapp2GenesisExported := liquidity.ExportGenesis(ctx2, simapp2.LiquidityKeeper)
	require.Equal(t, genesisExported, simapp2GenesisExported)
}
