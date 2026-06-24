package types_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/Victor118/liquidity/x/liquidity/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

func TestAlphabeticalDenomPair(t *testing.T) {
	denomA := "uCoinA"
	denomB := "uCoinB"
	afterDenomA, afterDenomB := types.AlphabeticalDenomPair(denomA, denomB)
	require.Equal(t, denomA, afterDenomA)
	require.Equal(t, denomB, afterDenomB)

	afterDenomA, afterDenomB = types.AlphabeticalDenomPair(denomB, denomA)
	require.Equal(t, denomA, afterDenomA)
	require.Equal(t, denomB, afterDenomB)
}

func TestGetReserveAcc(t *testing.T) {
	expectedReserveAcc := sdk.AccAddress(address.Module(types.ModuleName, []byte("D35A0CC16EE598F90B044CE296A405BA9C381E38837599D96F2F70C2F02A23A4")))
	tests := []struct {
		poolCoinDenom      string
		expectedReserveAcc sdk.AccAddress
		expPass            bool
	}{
		{
			poolCoinDenom:      "poolD35A0CC16EE598F90B044CE296A405BA9C381E38837599D96F2F70C2F02A23A4",
			expectedReserveAcc: expectedReserveAcc,
			expPass:            true,
		},
		{
			poolCoinDenom:      "poolD35A0CC16EE598F90B044CE296A405BA9C381E38837599D96F2F70C2F02A",
			expectedReserveAcc: nil,
			expPass:            false,
		},
		{
			poolCoinDenom:      "D35A0CC16EE598F90B044CE296A405BA9C381E38837599D96F2F70C2F02A23A4",
			expectedReserveAcc: nil,
			expPass:            false,
		},
		{
			poolCoinDenom:      "ibc/D35A0CC16EE598F90B044CE296A405BA9C381E38837599D96F2F70C2F02A23A4",
			expectedReserveAcc: nil,
			expPass:            false,
		},
	}

	for _, tc := range tests {
		reserveAcc, err := types.GetReserveAcc(tc.poolCoinDenom)
		if tc.expPass {
			require.Equal(t, tc.expectedReserveAcc, reserveAcc)
		} else {
			require.Nil(t, reserveAcc)
			require.ErrorIs(t, err, types.ErrInvalidDenom)
		}
	}
}

func TestSortDenoms(t *testing.T) {
	tests := []struct {
		denoms         []string
		expectedDenoms []string
	}{
		{
			denoms:         []string{"uCoinB", "uCoinA"},
			expectedDenoms: []string{"uCoinA", "uCoinB"},
		},
		{
			denoms:         []string{"uCoinC", "uCoinA", "uCoinB"},
			expectedDenoms: []string{"uCoinA", "uCoinB", "uCoinC"},
		},
		{
			denoms:         []string{"uCoinC", "uCoinA", "uCoinD", "uCoinB"},
			expectedDenoms: []string{"uCoinA", "uCoinB", "uCoinC", "uCoinD"},
		},
	}

	for _, tc := range tests {
		sortedDenoms := types.SortDenoms(tc.denoms)
		require.Equal(t, tc.expectedDenoms, sortedDenoms)
	}
}

func TestGetPoolInformation(t *testing.T) {
	testCases := []struct {
		reserveCoinDenoms     []string
		poolTypeID            uint32
		expectedPoolName      string
		expectedReserveAcc    string
		expectedPoolCoinDenom string
		len32                 bool
	}{
		{
			reserveCoinDenoms:     []string{"denomX", "denomY"},
			poolTypeID:            uint32(1),
			expectedPoolName:      "denomX/denomY/1",
			expectedReserveAcc:    sdk.AccAddress(address.Module(types.ModuleName, []byte("D35A0CC16EE598F90B044CE296A405BA9C381E38837599D96F2F70C2F02A23A4"))).String(),
			expectedPoolCoinDenom: "poolD35A0CC16EE598F90B044CE296A405BA9C381E38837599D96F2F70C2F02A23A4",
		},
		{
			reserveCoinDenoms:     []string{"stake", "token"},
			poolTypeID:            uint32(1),
			expectedPoolName:      "stake/token/1",
			expectedReserveAcc:    sdk.AccAddress(address.Module(types.ModuleName, []byte("E4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07"))).String(),
			expectedPoolCoinDenom: "poolE4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07",
		},
		{
			reserveCoinDenoms:     []string{"uatom", "uusd"},
			poolTypeID:            uint32(2),
			expectedPoolName:      "uatom/uusd/2",
			expectedReserveAcc:    sdk.AccAddress(address.Module(types.ModuleName, []byte("3036F43CB8131A1A63D2B3D3B11E9CF6FA2A2B6FEC17D5AD283C25C939614A8C"))).String(),
			expectedPoolCoinDenom: "pool3036F43CB8131A1A63D2B3D3B11E9CF6FA2A2B6FEC17D5AD283C25C939614A8C",
		},
		{
			reserveCoinDenoms:     []string{"uatom", "usdt"},
			poolTypeID:            uint32(3),
			expectedPoolName:      "uatom/usdt/3",
			expectedReserveAcc:    "cosmos1aqvez6g6wejw8hu35kplycf2taqsfkpj3ns3c5v4dhwazfdzhzastyr290",
			expectedPoolCoinDenom: "pool93E069B333B5ECEBFE24C6E1437E814003248E0DD7FF8B9F82119F4587449BA5",
		},
	}

	for _, tc := range testCases {
		poolName := types.PoolName(tc.reserveCoinDenoms, tc.poolTypeID)
		t.Logf("Poolname from reserviceCoinDenoms %v and pool type %d = %s", tc.reserveCoinDenoms, tc.poolTypeID, poolName)
		require.Equal(t, tc.expectedPoolName, poolName)

		reserveAcc := types.GetPoolReserveAcc(poolName)
		require.Equal(t, tc.expectedReserveAcc, reserveAcc.String())

		poolCoinDenom := types.GetPoolCoinDenom(poolName)
		require.Equal(t, tc.expectedPoolCoinDenom, poolCoinDenom)

		expectedReserveAcc, err := types.GetReserveAcc(poolCoinDenom)
		require.NoError(t, err)
		require.Equal(t, tc.expectedReserveAcc, expectedReserveAcc.String())
	}
}

func TestGetCoinsTotalAmount(t *testing.T) {
	testCases := []struct {
		coins        sdk.Coins
		expectResult math.Int
	}{
		{
			coins:        sdk.NewCoins(sdk.NewCoin("uCoinA", math.NewInt(100)), sdk.NewCoin("uCoinB", math.NewInt(100))),
			expectResult: math.NewInt(200),
		},
		{
			coins:        sdk.NewCoins(sdk.NewCoin("uCoinA", math.NewInt(100)), sdk.NewCoin("uCoinB", math.NewInt(300))),
			expectResult: math.NewInt(400),
		},
		{
			coins:        sdk.NewCoins(sdk.NewCoin("uCoinA", math.NewInt(500))),
			expectResult: math.NewInt(500),
		},
	}

	for _, tc := range testCases {
		totalAmount := types.GetCoinsTotalAmount(tc.coins)
		require.Equal(t, tc.expectResult, totalAmount)
	}
}

func TestValidateReserveCoinLimit(t *testing.T) {
	testCases := []struct {
		name                 string
		maxReserveCoinAmount math.Int
		depositCoins         sdk.Coins
		expectErr            bool
	}{
		{
			name:                 "valid case",
			maxReserveCoinAmount: math.ZeroInt(), // 0 means unlimited amount
			depositCoins:         sdk.NewCoins(sdk.NewCoin("uCoinA", math.NewInt(100_000_000_000)), sdk.NewCoin("uCoinB", math.NewInt(100))),
			expectErr:            false,
		},
		{
			name:                 "valid case",
			maxReserveCoinAmount: math.NewInt(1_000_000_000_000),
			depositCoins:         sdk.NewCoins(sdk.NewCoin("uCoinA", math.NewInt(500_000_000_000)), sdk.NewCoin("uCoinB", math.NewInt(500_000_000_000))),
			expectErr:            false,
		},
		{
			name:                 "negative value of max reserve coin amount",
			maxReserveCoinAmount: math.NewInt(-100),
			depositCoins:         sdk.NewCoins(sdk.NewCoin("uCoinA", math.NewInt(100_000_000_000)), sdk.NewCoin("uCoinB", math.NewInt(100))),
			expectErr:            true,
		},
		{
			name:                 "cannot exceed reserve coin limit amount",
			maxReserveCoinAmount: math.NewInt(1_000_000_000_000),
			depositCoins:         sdk.NewCoins(sdk.NewCoin("uCoinA", math.NewInt(1_000_000_000_000)), sdk.NewCoin("uCoinB", math.NewInt(100))),
			expectErr:            true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectErr {
				err := types.ValidateReserveCoinLimit(tc.maxReserveCoinAmount, tc.depositCoins)
				require.Equal(t, types.ErrExceededReserveCoinLimit, err)
			} else {
				err := types.ValidateReserveCoinLimit(tc.maxReserveCoinAmount, tc.depositCoins)
				require.NoError(t, err)
			}
		})
	}
}

func TestGetOfferCoinFee(t *testing.T) {
	testDenom := "test"
	testCases := []struct {
		name               string
		offerCoin          sdk.Coin
		swapFeeRate        math.LegacyDec
		expectOfferCoinFee sdk.Coin
	}{
		{
			name:               "case1",
			offerCoin:          sdk.NewCoin(testDenom, math.NewInt(1)),
			swapFeeRate:        types.DefaultSwapFeeRate,
			expectOfferCoinFee: sdk.NewCoin(testDenom, math.NewInt(1)),
		},
		{
			name:               "case2",
			offerCoin:          sdk.NewCoin(testDenom, math.NewInt(10)),
			swapFeeRate:        types.DefaultSwapFeeRate,
			expectOfferCoinFee: sdk.NewCoin(testDenom, math.NewInt(1)),
		},
		{
			name:               "case3",
			offerCoin:          sdk.NewCoin(testDenom, math.NewInt(100)),
			swapFeeRate:        types.DefaultSwapFeeRate,
			expectOfferCoinFee: sdk.NewCoin(testDenom, math.NewInt(1)),
		},
		{
			name:               "case4",
			offerCoin:          sdk.NewCoin(testDenom, math.NewInt(1000)),
			swapFeeRate:        types.DefaultSwapFeeRate,
			expectOfferCoinFee: sdk.NewCoin(testDenom, math.NewInt(2)),
		},
		{
			name:               "case5",
			offerCoin:          sdk.NewCoin(testDenom, math.NewInt(10000)),
			swapFeeRate:        types.DefaultSwapFeeRate,
			expectOfferCoinFee: sdk.NewCoin(testDenom, math.NewInt(15)),
		},
		{
			name:               "case6",
			offerCoin:          sdk.NewCoin(testDenom, math.NewInt(10001)),
			swapFeeRate:        types.DefaultSwapFeeRate,
			expectOfferCoinFee: sdk.NewCoin(testDenom, math.NewInt(16)),
		},
		{
			name:               "case7",
			offerCoin:          sdk.NewCoin(testDenom, math.NewInt(10700)),
			swapFeeRate:        types.DefaultSwapFeeRate,
			expectOfferCoinFee: sdk.NewCoin(testDenom, math.NewInt(17)),
		},
		{
			name:               "case8",
			offerCoin:          sdk.NewCoin(testDenom, math.NewInt(10000)),
			swapFeeRate:        math.LegacyZeroDec(),
			expectOfferCoinFee: sdk.NewCoin(testDenom, math.NewInt(0)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectOfferCoinFee, types.GetOfferCoinFee(tc.offerCoin, tc.swapFeeRate))
		})
	}
}

func TestCheckOverflow(t *testing.T) {
	testCases := []struct {
		name      string
		a         math.Int
		b         math.Int
		expectErr error
	}{
		{
			name:      "valid case",
			a:         math.NewInt(10000),
			b:         math.NewInt(100),
			expectErr: nil,
		},
		{
			name:      "overflow case",
			a:         math.NewInt(1_000_000_000_000_000_000).MulRaw(1_000_000),
			b:         math.NewInt(1_000_000_000_000_000_000).MulRaw(1_000_000_000_000_000_000).MulRaw(1_000_000_000_000_000_000),
			expectErr: types.ErrOverflowAmount,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := types.CheckOverflow(tc.a, tc.b)
			require.ErrorIs(t, err, tc.expectErr)
		})
	}
}

func TestCheckOverflowWithDec(t *testing.T) {
	testCases := []struct {
		name      string
		a         math.LegacyDec
		b         math.LegacyDec
		expectErr error
	}{
		{
			name:      "valid case",
			a:         math.LegacyMustNewDecFromStr("1.0"),
			b:         math.LegacyMustNewDecFromStr("0.0000001"),
			expectErr: nil,
		},
		{
			name:      "overflow case",
			a:         math.LegacyMustNewDecFromStr("100000000000000000000000000000000000000000000000000000000000.0").MulInt64(10),
			b:         math.LegacyMustNewDecFromStr("0.000000000000000001"),
			expectErr: types.ErrOverflowAmount,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := types.CheckOverflowWithDec(tc.a, tc.b)
			require.ErrorIs(t, err, tc.expectErr)
		})
	}
}
