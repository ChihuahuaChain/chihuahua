package keeper_test

import (
	"testing"

	testkeeper "github.com/ChihuahuaChain/chihuahua/testutil/keeper"

	"github.com/ChihuahuaChain/chihuahua/x/feeburn/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.FeeburnKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
	require.EqualValues(t, params.TxFeeBurnPercent, k.TxFeeBurnPercent(ctx))
}
