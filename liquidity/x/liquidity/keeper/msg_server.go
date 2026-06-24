package keeper

// DONTCOVER

// Although written in msg_server_test.go, it is approached at the keeper level rather than at the msgServer level
// so is not included in the coverage.

import (
	"context"
	"fmt"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/Victor118/liquidity/x/liquidity/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the distribution MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// Message server, handler for CreatePool msg
func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.GetCircuitBreakerEnabled(ctx) {
		return nil, types.ErrCircuitBreakerEnabled
	}

	pool, err := k.Keeper.CreatePool(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreatePool,
			sdk.NewAttribute(types.AttributeValuePoolId, strconv.FormatUint(pool.Id, 10)),
			sdk.NewAttribute(types.AttributeValuePoolTypeId, fmt.Sprintf("%d", msg.PoolTypeId)),
			sdk.NewAttribute(types.AttributeValuePoolName, pool.Name()),
			sdk.NewAttribute(types.AttributeValueReserveAccount, pool.ReserveAccountAddress),
			sdk.NewAttribute(types.AttributeValueDepositCoins, msg.DepositCoins.String()),
			sdk.NewAttribute(types.AttributeValuePoolCoinDenom, pool.PoolCoinDenom),
		),
	})

	return &types.MsgCreatePoolResponse{}, nil
}

// Message server, handler for MsgDepositWithinBatch
func (k msgServer) DepositWithinBatch(goCtx context.Context, msg *types.MsgDepositWithinBatch) (*types.MsgDepositWithinBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.GetCircuitBreakerEnabled(ctx) {
		return nil, types.ErrCircuitBreakerEnabled
	}

	poolBatch, found := k.GetPoolBatch(ctx, msg.PoolId)
	if !found {
		return nil, types.ErrPoolBatchNotExists
	}

	batchMsg, err := k.Keeper.DepositWithinBatch(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDepositWithinBatch,
			sdk.NewAttribute(types.AttributeValuePoolId, strconv.FormatUint(batchMsg.Msg.PoolId, 10)),
			sdk.NewAttribute(types.AttributeValueBatchIndex, strconv.FormatUint(poolBatch.Index, 10)),
			sdk.NewAttribute(types.AttributeValueMsgIndex, strconv.FormatUint(batchMsg.MsgIndex, 10)),
			sdk.NewAttribute(types.AttributeValueDepositCoins, batchMsg.Msg.DepositCoins.String()),
		),
	})

	return &types.MsgDepositWithinBatchResponse{}, nil
}

// Message server, handler for MsgWithdrawWithinBatch
func (k msgServer) WithdrawWithinBatch(goCtx context.Context, msg *types.MsgWithdrawWithinBatch) (*types.MsgWithdrawWithinBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	poolBatch, found := k.GetPoolBatch(ctx, msg.PoolId)
	if !found {
		return nil, types.ErrPoolBatchNotExists
	}

	batchMsg, err := k.Keeper.WithdrawWithinBatch(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawWithinBatch,
			sdk.NewAttribute(types.AttributeValuePoolId, strconv.FormatUint(batchMsg.Msg.PoolId, 10)),
			sdk.NewAttribute(types.AttributeValueBatchIndex, strconv.FormatUint(poolBatch.Index, 10)),
			sdk.NewAttribute(types.AttributeValueMsgIndex, strconv.FormatUint(batchMsg.MsgIndex, 10)),
			sdk.NewAttribute(types.AttributeValuePoolCoinDenom, batchMsg.Msg.PoolCoin.Denom),
			sdk.NewAttribute(types.AttributeValuePoolCoinAmount, batchMsg.Msg.PoolCoin.Amount.String()),
		),
	})

	return &types.MsgWithdrawWithinBatchResponse{}, nil
}

// Message server, handler for MsgSwapWithinBatch
func (k msgServer) Swap(goCtx context.Context, msg *types.MsgSwapWithinBatch) (*types.MsgSwapWithinBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.GetCircuitBreakerEnabled(ctx) {
		return nil, types.ErrCircuitBreakerEnabled
	}

	poolBatch, found := k.GetPoolBatch(ctx, msg.PoolId)
	if !found {
		return nil, types.ErrPoolBatchNotExists
	}

	batchMsg, err := k.Keeper.SwapWithinBatch(ctx, msg, types.CancelOrderLifeSpan)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSwapWithinBatch,
			sdk.NewAttribute(types.AttributeValuePoolId, strconv.FormatUint(batchMsg.Msg.PoolId, 10)),
			sdk.NewAttribute(types.AttributeValueBatchIndex, strconv.FormatUint(poolBatch.Index, 10)),
			sdk.NewAttribute(types.AttributeValueMsgIndex, strconv.FormatUint(batchMsg.MsgIndex, 10)),
			sdk.NewAttribute(types.AttributeValueSwapTypeId, strconv.FormatUint(uint64(batchMsg.Msg.SwapTypeId), 10)),
			sdk.NewAttribute(types.AttributeValueOfferCoinDenom, batchMsg.Msg.OfferCoin.Denom),
			sdk.NewAttribute(types.AttributeValueOfferCoinAmount, batchMsg.Msg.OfferCoin.Amount.String()),
			sdk.NewAttribute(types.AttributeValueOfferCoinFeeAmount, batchMsg.Msg.OfferCoinFee.Amount.String()),
			sdk.NewAttribute(types.AttributeValueDemandCoinDenom, batchMsg.Msg.DemandCoinDenom),
			sdk.NewAttribute(types.AttributeValueOrderPrice, batchMsg.Msg.OrderPrice.String()),
		),
	})

	return &types.MsgSwapWithinBatchResponse{}, nil
}

func (k msgServer) DirectSwap(goCtx context.Context, msg *types.MsgDirectSwap) (*types.MsgDirectSwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.GetCircuitBreakerEnabled(ctx) {
		return nil, types.ErrCircuitBreakerEnabled
	}

	receiveAmount, err := k.Keeper.DirectSwapExecution(ctx, msg.PoolId, msg.OfferCoin, msg.DemandCoinDenom, msg.OrderPrice, msg.GetSwapRequester())
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDirectSwap,
			sdk.NewAttribute(types.AttributeValuePoolId, strconv.FormatUint(msg.PoolId, 10)),
			sdk.NewAttribute(types.AttributeValueOfferCoinDenom, msg.OfferCoin.Denom),
			sdk.NewAttribute(types.AttributeValueOfferCoinAmount, msg.OfferCoin.Amount.String()),
			sdk.NewAttribute(types.AttributeValueDemandCoinDenom, msg.DemandCoinDenom),
			sdk.NewAttribute(types.AttributeValueOrderPrice, msg.OrderPrice.String()),
			sdk.NewAttribute(types.AttributeValueExchangedDemandCoinAmount, receiveAmount.String()),
		),
	})

	return &types.MsgDirectSwapResponse{
		ReceivedAmount: sdk.NewCoin(msg.DemandCoinDenom, receiveAmount),
	}, nil
}

// UpdateParams implements the gRPC MsgServer interface. When an UpdateParams
// proposal passes, it updates the module parameters. The update can only be
// performed if the requested authority is the Cosmos SDK governance module
// account.
func (k msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if k.authority.String() != req.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority.String(), req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.SetParams(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
