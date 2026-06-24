package types

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// TestValidateBuildersAddresses_RejectsInvalid guards against a chain-halt bug:
// builder commissions are paid in the liquidity batch EndBlocker via
// SendAmountToBuilders, which converts each builder address. An empty/invalid
// address must be rejected at param-validation time so it can never be
// persisted and trigger a panic during block production.
func TestValidateBuildersAddresses_RejectsInvalid(t *testing.T) {
	// Empty address must be rejected (used to be accepted, then panicked on payout).
	require.Error(t, validateBuildersAddresses([]WeightedAddress{
		{Address: "", Weight: math.LegacyOneDec()},
	}))

	// Malformed (non-bech32) address must be rejected.
	require.Error(t, validateBuildersAddresses([]WeightedAddress{
		{Address: "not-a-bech32-address", Weight: math.LegacyOneDec()},
	}))

	// Empty list is allowed: builder commission is simply not routed.
	require.NoError(t, validateBuildersAddresses([]WeightedAddress{}))

	// A valid address whose weights sum to 1 must pass.
	valid := sdk.AccAddress([]byte("builder_address_0001")).String()
	require.NoError(t, validateBuildersAddresses([]WeightedAddress{
		{Address: valid, Weight: math.LegacyOneDec()},
	}))
}
