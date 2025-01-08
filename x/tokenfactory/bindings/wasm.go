package bindings

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cosmos/cosmos-sdk/baseapp"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	tokenfactorykeeper "github.com/ChihuahuaChain/chihuahua/x/tokenfactory/keeper"
	liquiditykeeper "github.com/Victor118/liquidity/x/liquidity/keeper"
)

func RegisterCustomPlugins(
	grpcRouter *baseapp.GRPCQueryRouter,
	appCodec codec.Codec,
	bank bankkeeper.Keeper,
	tokenFactory *tokenfactorykeeper.Keeper,
	liquidity *liquiditykeeper.Keeper,
) []wasmkeeper.Option {
	wasmQueryPlugin := NewQueryPlugin(bank, tokenFactory)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom:   CustomQuerier(wasmQueryPlugin),
		Stargate: wasmkeeper.AcceptListStargateQuerier(AcceptedStargateQueries(), grpcRouter, appCodec),
		Grpc:     wasmkeeper.AcceptListGrpcQuerier(AcceptedStargateQueries(), grpcRouter, appCodec),
	})
	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		CustomMessageDecorator(bank, tokenFactory, liquidity),
	)

	return []wasmkeeper.Option{
		queryPluginOpt,
		messengerDecoratorOpt,
	}
}
