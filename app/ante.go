package app

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	ibcante "github.com/cosmos/ibc-go/v11/modules/core/ante"
	ibckeeper "github.com/cosmos/ibc-go/v11/modules/core/keeper"

	corestoretypes "cosmossdk.io/core/store"
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	feeburnAnte "github.com/ChihuahuaChain/chihuahua/x/feeburn/ante"
	feeburnkeeper "github.com/ChihuahuaChain/chihuahua/x/feeburn/keeper"
)

// HandlerOptions are the options required for constructing a default SDK AnteHandler.
type HandlerOptions struct {
	ante.HandlerOptions
	BankKeeper            bankkeeper.Keeper
	FeeburnKeeper         *feeburnkeeper.Keeper
	IBCKeeper             *ibckeeper.Keeper
	NodeConfig            *wasmtypes.NodeConfig
	WasmKeeper            *wasmkeeper.Keeper
	TXCounterStoreService corestoretypes.KVStoreService
}

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer, burning part of them (see x/feeburn).
func NewAnteHandler(options HandlerOptions) (sdk.AnteHandler, error) {
	if options.AccountKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "account keeper is required for ante builder")
	}
	if options.BankKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "bank keeper is required for ante builder")
	}
	if options.FeeburnKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "fee burn keeper is required for ante builder")
	}
	if options.SignModeHandler == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	}
	if options.IBCKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "ibc keeper is required for ante builder")
	}
	if options.NodeConfig == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "wasm config is required for ante builder")
	}
	if options.WasmKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "wasm keeper is required for ante builder")
	}
	if options.TXCounterStoreService == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "wasm store service is required for ante builder")
	}

	anteDecorators := []sdk.AnteDecorator{
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		wasmkeeper.NewLimitSimulationGasDecorator(options.NodeConfig.SimulationGasLimit), // after setup context to enforce limits early
		wasmkeeper.NewCountTXDecorator(options.TXCounterStoreService),
		wasmkeeper.NewGasRegisterDecorator(options.WasmKeeper.GetGasRegister()),
		wasmkeeper.NewTxContractsDecorator(),
		ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		feeburnAnte.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker, *options.FeeburnKeeper),
		ante.NewSetPubKeyDecorator(options.AccountKeeper), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ibcante.NewRedundantRelayDecorator(options.IBCKeeper),
	}

	return sdk.ChainAnteDecorators(anteDecorators...), nil
}
