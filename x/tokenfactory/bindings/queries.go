package bindings

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	bindingstypes "github.com/ChihuahuaChain/chihuahua/x/tokenfactory/bindings/types"
	tokenfactorykeeper "github.com/ChihuahuaChain/chihuahua/x/tokenfactory/keeper"
	tftypes "github.com/ChihuahuaChain/chihuahua/x/tokenfactory/types"
)

type QueryPlugin struct {
	bankKeeper         bankkeeper.Keeper
	tokenFactoryKeeper *tokenfactorykeeper.Keeper
}

// NewQueryPlugin returns a reference to a new QueryPlugin.
func NewQueryPlugin(b bankkeeper.Keeper, tfk *tokenfactorykeeper.Keeper) *QueryPlugin {
	return &QueryPlugin{
		bankKeeper:         b,
		tokenFactoryKeeper: tfk,
	}
}

// GetDenomAdmin is a query to get denom admin.
func (qp QueryPlugin) GetDenomAdmin(ctx sdk.Context, denom string) (*bindingstypes.AdminResponse, error) {
	metadata, err := qp.tokenFactoryKeeper.GetAuthorityMetadata(ctx, denom)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin for denom: %s", denom)
	}
	return &bindingstypes.AdminResponse{Admin: metadata.Admin}, nil
}

// GetDenomAdmin is a query to get denom admin.
func (qp QueryPlugin) GetStakedropByDenom(ctx sdk.Context, denom string) (*bindingstypes.StakedropByDenomResponse, error) {
	stakedropsByDenom := []tftypes.Stakedrop{}
	iter, err := qp.tokenFactoryKeeper.ActiveStakedrop.Indexes.StakedropByDenom.MatchExact(ctx, denom)
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
	return &bindingstypes.StakedropByDenomResponse{Stakedrops: stakedropsByDenom}, nil
}

func (qp QueryPlugin) GetDenomsByCreator(ctx sdk.Context, creator string) (*bindingstypes.DenomsByCreatorResponse, error) {
	// TODO: validate creator address
	denoms := qp.tokenFactoryKeeper.GetDenomsFromCreator(ctx, creator)
	return &bindingstypes.DenomsByCreatorResponse{Denoms: denoms}, nil
}

func (qp QueryPlugin) GetMetadata(ctx sdk.Context, denom string) (*bindingstypes.MetadataResponse, error) {
	metadata, found := qp.bankKeeper.GetDenomMetaData(ctx, denom)
	var parsed *bindingstypes.Metadata
	if found {
		parsed = SdkMetadataToWasm(metadata)
	}
	return &bindingstypes.MetadataResponse{Metadata: parsed}, nil
}

func (qp QueryPlugin) GetParams(ctx sdk.Context) (*bindingstypes.ParamsResponse, error) {
	params := qp.tokenFactoryKeeper.GetParams(ctx)
	return &bindingstypes.ParamsResponse{
		Params: bindingstypes.Params{
			DenomCreationFee: ConvertSdkCoinsToWasmCoins(params.DenomCreationFee),
		},
	}, nil
}
