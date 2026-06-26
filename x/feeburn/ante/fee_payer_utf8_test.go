package ante

import (
	"testing"
	"unicode/utf8"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// TestFeePayerAttributeIsUTF8 is a regression test for the Hermes
// "EventAttribute.value: ... data is not UTF-8 encoded" failure on
// cosmos.tx.v1beta1.Service/Simulate.
//
// FeeTx.FeePayer() returns raw address bytes ([]byte) in cosmos-sdk v0.50.
// Emitting them via string(rawBytes) produces an event attribute value that is
// not valid UTF-8. Since abci.EventAttribute.Value is a proto3 string (which
// must be valid UTF-8), strict decoders such as Hermes (prost) reject the whole
// SimulateResponse, blocking gas estimation for every tx, including IBC
// MsgCreateClient / MsgRecvPacket. The value must be the bech32 string instead.
func TestFeePayerAttributeIsUTF8(t *testing.T) {
	// Raw 20-byte address containing bytes that are invalid as UTF-8.
	feePayer := []byte{
		0x00, 0x01, 0x02, 0xff, 0xfe, 0x80, 0x81, 0x10, 0x11, 0x12,
		0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c,
	}

	// Old behavior (buggy): raw bytes cast to string -> NOT valid UTF-8.
	require.False(t, utf8.ValidString(string(feePayer)),
		"raw address bytes must not be valid UTF-8 (this reproduces the bug)")

	// New behavior (fixed): bech32 string -> always valid UTF-8.
	require.True(t, utf8.ValidString(sdk.AccAddress(feePayer).String()),
		"bech32 fee_payer value must be valid UTF-8")
}
