package simulation

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectestutil "github.com/cosmos/cosmos-sdk/codec/testutil"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"

	"github.com/ChihuahuaChain/chihuahua/x/liquidity/types"
)

// simTxConfig returns a tx config using the bech32 prefixes of the global sdk config.
func simTxConfig() client.TxConfig {
	cfg := sdk.GetConfig()
	registry := codectestutil.CodecOptions{
		AccAddressPrefix: cfg.GetBech32AccountAddrPrefix(),
		ValAddressPrefix: cfg.GetBech32ValidatorAddrPrefix(),
	}.NewInterfaceRegistry()
	std.RegisterInterfaces(registry)
	types.RegisterInterfaces(registry)
	return tx.NewTxConfig(codec.NewProtoCodec(registry), tx.DefaultSignModes)
}
