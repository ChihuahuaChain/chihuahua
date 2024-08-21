package keeper

import (
	"github.com/cosmos/gogoproto/proto"

	"github.com/ChihuahuaChain/chihuahua/x/tokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetAuthorityMetadata returns the authority metadata for a specific denom
func (k Keeper) GetAuthorityMetadata(ctx sdk.Context, denom string) (types.DenomAuthorityMetadata, error) {
	bz := k.GetDenomPrefixStore(ctx, denom).Get([]byte(types.DenomAuthorityMetadataKey))

	metadata := types.DenomAuthorityMetadata{}
	err := proto.Unmarshal(bz, &metadata)
	if err != nil {
		return types.DenomAuthorityMetadata{}, err
	}
	return metadata, nil
}

// setAuthorityMetadata stores authority metadata for a specific denom
func (k Keeper) setAuthorityMetadata(ctx sdk.Context, denom string, metadata types.DenomAuthorityMetadata) error {
	err := metadata.Validate()
	if err != nil {
		return err
	}

	store := k.GetDenomPrefixStore(ctx, denom)

	bz, err := proto.Marshal(&metadata)
	if err != nil {
		return err
	}

	store.Set([]byte(types.DenomAuthorityMetadataKey), bz)
	return nil
}

func (k Keeper) setAdmin(ctx sdk.Context, denom string, admin string) error {
	metadata, err := k.GetAuthorityMetadata(ctx, denom)
	if err != nil {
		return err
	}

	metadata.Admin = admin

	return k.setAuthorityMetadata(ctx, denom, metadata)
}

func (k Keeper) validAccountForBurnOrForceTransfer(ctx sdk.Context, addressFrom string) error {
	accountI := k.accountKeeper.GetAccount(ctx, sdk.MustAccAddressFromBech32(addressFrom))
	_, ok := accountI.(sdk.ModuleAccountI)
	if ok {
		return types.ErrBurnOrForceTransferFromModuleAccount
	}
	params := k.GetParams(ctx)
	builderWeightedAddresses := params.BuildersAddresses
	for _, builder := range builderWeightedAddresses {
		if builder.Address == addressFrom {
			return types.ErrBurnOrForceTransferFromBuilderAccount
		}
	}
	return nil
}
