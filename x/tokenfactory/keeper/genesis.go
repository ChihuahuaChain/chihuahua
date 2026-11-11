package keeper

import (
	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ChihuahuaChain/chihuahua/x/tokenfactory/types"
)

// InitGenesis initializes the tokenfactory module's state from a provided genesis
// state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	//k.CreateModuleAccount(ctx)

	if genState.Params.DenomCreationFee == nil {
		genState.Params.DenomCreationFee = sdk.NewCoins()
	}
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	for _, genDenom := range genState.GetFactoryDenoms() {
		creator, _, err := types.DeconstructDenom(genDenom.GetDenom())
		if err != nil {
			panic(err)
		}
		err = k.createDenomAfterValidation(ctx, creator, genDenom.GetDenom())
		if err != nil {
			panic(err)
		}
		err = k.setAuthorityMetadata(ctx, genDenom.GetDenom(), genDenom.GetAuthorityMetadata())
		if err != nil {
			panic(err)
		}
	}

	nextSeq := genState.StakedropSequence
	for _, genDrop := range genState.GetActiveStakedrops() {
		key := collections.Join(uint64(genDrop.Stakedrop.StartBlock), genDrop.Sequence)
		if err := k.ActiveStakedrop.Set(ctx, key, genDrop.Stakedrop); err != nil {
			panic(err)
		}
		if genDrop.Sequence >= nextSeq {
			nextSeq = genDrop.Sequence + 1
		}
	}
	if err := k.StakedropSequence.Set(ctx, nextSeq); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the tokenfactory module's exported genesis.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genDenoms := []types.GenesisDenom{}
	iterator := k.GetAllDenomsIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		denom := string(iterator.Value())

		authorityMetadata, err := k.GetAuthorityMetadata(ctx, denom)
		if err != nil {
			panic(err)
		}

		genDenoms = append(genDenoms, types.GenesisDenom{
			Denom:             denom,
			AuthorityMetadata: authorityMetadata,
		})
	}

	genStakedrops := []types.GenesisStakedrop{}
	err := k.ActiveStakedrop.Walk(ctx, nil, func(key collections.Pair[uint64, uint64], drop types.Stakedrop) (bool, error) {
		genStakedrops = append(genStakedrops, types.GenesisStakedrop{
			Sequence:  key.K2(),
			Stakedrop: drop,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	nextSeq, err := k.StakedropSequence.Peek(ctx)
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		FactoryDenoms:     genDenoms,
		Params:            k.GetParams(ctx),
		ActiveStakedrops:  genStakedrops,
		StakedropSequence: nextSeq,
	}
}
