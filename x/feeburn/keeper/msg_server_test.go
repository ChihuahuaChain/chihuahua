package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/ChihuahuaChain/chihuahua/testutil/keeper"

	"github.com/ChihuahuaChain/chihuahua/x/feeburn/keeper"
	"github.com/ChihuahuaChain/chihuahua/x/feeburn/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.FeeburnKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
