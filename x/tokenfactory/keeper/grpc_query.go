package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ChihuahuaChain/chihuahua/x/tokenfactory/types"
	"github.com/cosmos/cosmos-sdk/types/query"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(ctx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(sdkCtx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Keeper) DenomAuthorityMetadata(ctx context.Context, req *types.QueryDenomAuthorityMetadataRequest) (*types.QueryDenomAuthorityMetadataResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	authorityMetadata, err := k.GetAuthorityMetadata(sdkCtx, req.GetDenom())
	if err != nil {
		return nil, err
	}

	return &types.QueryDenomAuthorityMetadataResponse{AuthorityMetadata: authorityMetadata}, nil
}

func (k Keeper) DenomsFromCreator(ctx context.Context, req *types.QueryDenomsFromCreatorRequest) (*types.QueryDenomsFromCreatorResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	denoms := k.GetDenomsFromCreator(sdkCtx, req.GetCreator())
	return &types.QueryDenomsFromCreatorResponse{Denoms: denoms}, nil
}

func (k Keeper) StakeDrops(ctx context.Context, req *types.QueryStakeDropsRequest) (*types.QueryStakeDropsResponse, error) {

	if req.Pagination != nil {
		stakedrops, pageResp, err := query.CollectionPaginate(ctx, k.ActiveStakedrop, req.Pagination, func(key collections.Pair[uint64, uint64], value types.Stakedrop) (types.Stakedrop, error) {
			return value, nil
		})

		if err != nil {

			return nil, err
		}

		return &types.QueryStakeDropsResponse{Stakedrops: stakedrops, Pagination: pageResp}, nil
	} else {
		iter, err := k.ActiveStakedrop.Iterate(ctx, nil)
		if err != nil {
			return nil, err
		}
		stakedrops, err := iter.Values()
		if err != nil {
			return nil, err
		}
		return &types.QueryStakeDropsResponse{Stakedrops: stakedrops, Pagination: nil}, nil
	}

}

func (k Keeper) StakeDropsFromDenom(ctx context.Context, req *types.QueryStakeDropFromDenomRequest) (*types.QueryStakeDropFromDenomResponse, error) {
	stakedropsByDenom := []types.Stakedrop{}
	iter, err := k.ActiveStakedrop.Indexes.StakedropByDenom.MatchExact(ctx, req.Denom)
	if err != nil {
		return nil, err
	}
	for ; iter.Valid(); iter.Next() {
		key, err := iter.PrimaryKey()
		if err != nil {
			return nil, err
		}
		stakedrop, err := k.ActiveStakedrop.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		stakedropsByDenom = append(stakedropsByDenom, stakedrop)
	}

	return &types.QueryStakeDropFromDenomResponse{Stakedrops: stakedropsByDenom}, nil

}
