package ante

import (
	"math"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// checkTxFee enforces the effective minimum gas price:
//
//	effective = per-denom max( global min_gas_prices (on-chain feeburn param),
//	                           validator-local ctx.MinGasPrices() )
//
// The global floor is enforced in BOTH CheckTx and DeliverTx, so it is a true
// consensus floor. The validator-local component is mempool-only and is folded
// in only during CheckTx; it is NEVER read in the DeliverTx path, so the
// consensus decision depends only on on-chain params + tx data (deterministic).
func (dfd DeductFeeDecorator) checkTxFee(ctx sdk.Context, tx sdk.Tx) (sdk.Coins, int64, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return nil, 0, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	feeCoins := feeTx.GetFee()
	gas := feeTx.GetGas()

	// Global, governance-settable floor read from feeburn params (deterministic
	// across all nodes).
	minGasPrices := dfd.feeburnKeeper.GetParams(ctx).MinGasPrices

	// Fold in the validator-local floor for mempool admission only. Do NOT read
	// ctx.MinGasPrices() outside IsCheckTx: it is node-local and would break
	// consensus if consulted in DeliverTx.
	if ctx.IsCheckTx() {
		minGasPrices = maxDecCoins(minGasPrices, ctx.MinGasPrices())
	}

	// Skip the floor during genesis/gentx processing (height 0) and when no
	// floor is configured. Simulation is already excluded by AnteHandle.
	if ctx.BlockHeight() > 0 && !minGasPrices.IsZero() {
		requiredFees := make(sdk.Coins, len(minGasPrices))

		// Determine the required fees by multiplying each required minimum gas
		// price by the gas limit, where fee = ceil(minGasPrice * gasLimit).
		glDec := sdkmath.LegacyNewDec(int64(gas))
		for i, gp := range minGasPrices {
			fee := gp.Amount.Mul(glDec)
			requiredFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
		}

		if !feeCoins.IsAnyGTE(requiredFees) {
			return nil, 0, errorsmod.Wrapf(sdkerrors.ErrInsufficientFee, "insufficient fees; got: %s required: %s", feeCoins, requiredFees)
		}
	}

	priority := getTxPriority(feeCoins, int64(gas))
	return feeCoins, priority, nil
}

// maxDecCoins returns the per-denom maximum of two DecCoins: the union of
// denoms, each carrying the larger of the two amounts, sorted. Only used in
// CheckTx, so it never influences the consensus decision.
func maxDecCoins(a, b sdk.DecCoins) sdk.DecCoins {
	denoms := make(map[string]struct{}, len(a)+len(b))
	for _, c := range a {
		denoms[c.Denom] = struct{}{}
	}
	for _, c := range b {
		denoms[c.Denom] = struct{}{}
	}

	out := sdk.DecCoins{}
	for denom := range denoms {
		amt := a.AmountOf(denom)
		if bAmt := b.AmountOf(denom); bAmt.GT(amt) {
			amt = bAmt
		}
		out = out.Add(sdk.NewDecCoinFromDec(denom, amt))
	}
	return out.Sort()
}

// getTxPriority returns a naive tx priority based on the amount of the smallest denomination of the gas price
// provided in a transaction.
// NOTE: This implementation should be used with a great consideration as it opens potential attack vectors
// where txs with multiple coins could not be prioritize as expected.
func getTxPriority(fee sdk.Coins, gas int64) int64 {
	var priority int64
	for _, c := range fee {
		p := int64(math.MaxInt64)
		gasPrice := c.Amount.QuoRaw(gas)
		if gasPrice.IsInt64() {
			p = gasPrice.Int64()
		}
		if priority == 0 || p < priority {
			priority = p
		}
	}

	return priority
}
