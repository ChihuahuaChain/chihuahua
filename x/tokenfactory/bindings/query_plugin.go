package bindings

import (
	"encoding/json"
	"fmt"

	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"

	errorsmod "cosmossdk.io/errors"
	bindingstypes "github.com/ChihuahuaChain/chihuahua/x/tokenfactory/bindings/types"
	tftypes "github.com/ChihuahuaChain/chihuahua/x/tokenfactory/types"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types" //nolint:staticcheck
	ibcconnectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
)

// CustomQuerier dispatches custom CosmWasm bindings queries.
func CustomQuerier(qp *QueryPlugin) func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var contractQuery bindingstypes.TokenFactoryQuery
		if err := json.Unmarshal(request, &contractQuery); err != nil {
			return nil, errorsmod.Wrap(err, "osmosis query")
		}

		switch {
		case contractQuery.FullDenom != nil:
			creator := contractQuery.FullDenom.CreatorAddr
			subdenom := contractQuery.FullDenom.Subdenom

			fullDenom, err := GetFullDenom(creator, subdenom)
			if err != nil {
				return nil, errorsmod.Wrap(err, "osmo full denom query")
			}

			res := bindingstypes.FullDenomResponse{
				Denom: fullDenom,
			}

			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "failed to marshal FullDenomResponse")
			}

			return bz, nil

		case contractQuery.Admin != nil:
			res, err := qp.GetDenomAdmin(ctx, contractQuery.Admin.Denom)
			if err != nil {
				return nil, err
			}

			bz, err := json.Marshal(res)
			if err != nil {
				return nil, fmt.Errorf("failed to JSON marshal AdminResponse: %w", err)
			}

			return bz, nil

		case contractQuery.Metadata != nil:
			res, err := qp.GetMetadata(ctx, contractQuery.Metadata.Denom)
			if err != nil {
				return nil, err
			}

			bz, err := json.Marshal(res)
			if err != nil {
				return nil, fmt.Errorf("failed to JSON marshal MetadataResponse: %w", err)
			}

			return bz, nil

		case contractQuery.DenomsByCreator != nil:
			res, err := qp.GetDenomsByCreator(ctx, contractQuery.DenomsByCreator.Creator)
			if err != nil {
				return nil, err
			}

			bz, err := json.Marshal(res)
			if err != nil {
				return nil, fmt.Errorf("failed to JSON marshal DenomsByCreatorResponse: %w", err)
			}

			return bz, nil

		case contractQuery.Params != nil:

			res, err := qp.GetParams(ctx)
			if err != nil {
				return nil, err
			}

			bz, err := json.Marshal(res)
			if err != nil {
				return nil, fmt.Errorf("failed to JSON marshal ParamsResponse: %w", err)
			}
			return bz, nil
		case contractQuery.Stakedrop != nil:
			stakedropsByDenom := []tftypes.Stakedrop{}
			iter, err := qp.tokenFactoryKeeper.ActiveStakedrop.Indexes.StakedropByDenom.MatchExact(ctx, contractQuery.Stakedrop.Denom)
			if err != nil {
				return nil, err
			}
			for ; iter.Valid(); iter.Next() {
				key, err := iter.PrimaryKey()
				if err != nil {
					return nil, err
				}
				stakedrop, err := qp.tokenFactoryKeeper.ActiveStakedrop.Get(ctx, key)
				if err != nil {
					return nil, err
				}
				stakedropsByDenom = append(stakedropsByDenom, stakedrop)
			}
			if err != nil {
				return nil, err
			}

			bz, err := json.Marshal(stakedropsByDenom)
			if err != nil {
				return nil, fmt.Errorf("failed to JSON marshal ParamsResponse: %w", err)
			}

			return bz, nil

		default:
			return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown token query variant"}
		}
	}
}

// ConvertSdkCoinsToWasmCoins converts sdk type coins to wasm vm type coins
func ConvertSdkCoinsToWasmCoins(coins []sdk.Coin) []wasmvmtypes.Coin {
	var toSend []wasmvmtypes.Coin
	for _, coin := range coins {
		c := ConvertSdkCoinToWasmCoin(coin)
		toSend = append(toSend, c)
	}
	return toSend
}

// ConvertSdkCoinToWasmCoin converts a sdk type coin to a wasm vm type coin
func ConvertSdkCoinToWasmCoin(coin sdk.Coin) wasmvmtypes.Coin {
	return wasmvmtypes.Coin{
		Denom: coin.Denom,
		// Note: tokenfactory tokens have 18 decimal places, so 10^22 is common, no longer in u64 range
		Amount: coin.Amount.String(),
	}
}

func AcceptedStargateQueries() wasmkeeper.AcceptedQueries {
	return wasmkeeper.AcceptedQueries{
		// ibc
		"/ibc.core.client.v1.Query/ClientState":         &ibcclienttypes.QueryClientStateResponse{},
		"/ibc.core.client.v1.Query/ConsensusState":      &ibcclienttypes.QueryConsensusStateResponse{},
		"/ibc.core.connection.v1.Query/Connection":      &ibcconnectiontypes.QueryConnectionResponse{},
		"/ibc.core.channel.v1.Query/ChannelClientState": &ibcchanneltypes.QueryChannelClientStateResponse{},

		// token factory
		"/osmosis.tokenfactory.v1beta1.Query/Params":                 &tftypes.QueryParamsResponse{},
		"/osmosis.tokenfactory.v1beta1.Query/DenomAuthorityMetadata": &tftypes.QueryDenomAuthorityMetadataResponse{},
		"/osmosis.tokenfactory.v1beta1.Query/DenomsFromCreator":      &tftypes.QueryDenomsFromCreatorResponse{},

		// transfer
		"/ibc.applications.transfer.v1.Query/DenomTrace":    &ibctransfertypes.QueryDenomTraceResponse{},
		"/ibc.applications.transfer.v1.Query/EscrowAddress": &ibctransfertypes.QueryEscrowAddressResponse{},

		// auth
		"/cosmos.auth.v1beta1.Query/Account": &authtypes.QueryAccountResponse{},
		"/cosmos.auth.v1beta1.Query/Params":  &authtypes.QueryParamsResponse{},

		// bank
		"/cosmos.bank.v1beta1.Query/Balance":       &banktypes.QueryBalanceResponse{},
		"/cosmos.bank.v1beta1.Query/DenomMetadata": &banktypes.QueryDenomsMetadataResponse{},
		"/cosmos.bank.v1beta1.Query/Params":        &banktypes.QueryParamsResponse{},
		"/cosmos.bank.v1beta1.Query/SupplyOf":      &banktypes.QuerySupplyOfResponse{},
	}
}
