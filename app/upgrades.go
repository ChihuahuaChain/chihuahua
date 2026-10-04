package app

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"strings"

	ibctransfertypes "github.com/cosmos/ibc-go/v11/modules/apps/transfer/types"
	"google.golang.org/protobuf/encoding/protowire"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/runtime"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	feeburnmoduletypes "github.com/ChihuahuaChain/chihuahua/x/feeburn/types"
)

// removedModules are the modules dropped by the v10 upgrade, together with
// their stores:
//   - alliance: retired in v9.5.0
//   - capability, feeibc: removed from ibc-go
//   - crisis, nft, circuit: no longer maintained by the Cosmos SDK, unused on chain
//   - params: every module manages its own params
var removedModules = []string{
	"alliance",
	"capability",
	"feeibc",
	"crisis",
	"nft",
	"circuit",
	"params",
}

// initialTotalBurned seeds the on-chain burned fees counter introduced in v10:
// the burned total published on burn.chihuahua.wtf, frozen on 2026-10-02.
var initialTotalBurned = sdk.NewCoins(sdk.NewCoin("uhuahua", sdkmath.NewInt(482_750_566_000_000)))

// RegisterUpgradeHandlers registers the upgrade handlers
func (app *App) RegisterUpgradeHandlers(cfg module.Configurator) {
	app.UpgradeKeeper.SetUpgradeHandler(UpgradeName, func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		for _, name := range removedModules {
			delete(fromVM, name)
		}

		// the transfer migration from denom traces to denoms panics on traces
		// whose base denom looks like a path: set those aside and migrate them here
		denoms, err := app.takeAmbiguousDenomTraces(sdkCtx)
		if err != nil {
			return nil, err
		}

		vm, err := app.mm.RunMigrations(ctx, cfg, fromVM)
		if err != nil {
			return nil, err
		}

		for _, denom := range denoms {
			app.TransferKeeper.SetDenom(sdkCtx, denom)
			sdkCtx.Logger().Info("migrated ambiguous denom trace", "denom", denom.Path(), "ibc_denom", denom.IBCDenom())
		}

		if err := app.FeeburnKeeper.SetTotalBurned(sdkCtx, initialTotalBurned); err != nil {
			return nil, err
		}
		return vm, nil
	})

	app.UpgradeKeeper.SetUpgradeHandler(PatchUpgradeName, func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		// Materialize the feeburn module account so the public burn address
		// exists with its burner permission from the first post-upgrade block,
		// even before it first receives funds. If a plain account already squats
		// the address (someone sent to it pre-upgrade), drop that record first:
		// GetModuleAccount panics on a non-module account, and x/bank keeps the
		// balance by address, so it survives the account record being replaced.
		burnAddr := authtypes.NewModuleAddress(feeburnmoduletypes.ModuleName)
		if acc := app.AccountKeeper.GetAccount(sdkCtx, burnAddr); acc != nil {
			if _, ok := acc.(authtypes.ModuleAccountI); !ok {
				app.AccountKeeper.RemoveAccount(sdkCtx, acc)
			}
		}
		app.AccountKeeper.GetModuleAccount(sdkCtx, feeburnmoduletypes.ModuleName)

		// Set the chain-wide minimum gas price floor enforced by x/feeburn in the
		// ante handler. GetParams preserves the existing TxFeeBurnPercent (50% on
		// mainnet); only MinGasPrices is introduced here. Governance can tune it
		// later via MsgUpdateParams.
		minGasPrices, err := sdk.ParseDecCoins(RecommendedMinGasPrices)
		if err != nil {
			return nil, err
		}
		params := app.FeeburnKeeper.GetParams(sdkCtx)
		params.MinGasPrices = minGasPrices
		if err := app.FeeburnKeeper.SetParams(sdkCtx, params); err != nil {
			return nil, err
		}

		return app.mm.RunMigrations(ctx, cfg, fromVM)
	})
}

// takeAmbiguousDenomTraces removes from the transfer store the denom traces
// that ibc-go cannot convert to denoms, and returns them converted.
//
// ibc-go v10+ rebuilds a denom from the full path of the trace, reading
// "port/channel" pairs as hops. A base denom containing a slash, like
// "IRO/heilelonmusk_653667-1", is then read as one more hop, leaving the base
// denom blank. Here the hops come from the trace path alone, which keeps the
// ibc/ hash, and so the balances, unchanged.
func (app *App) takeAmbiguousDenomTraces(ctx sdk.Context) ([]ibctransfertypes.Denom, error) {
	store := runtime.KVStoreAdapter(runtime.NewKVStoreService(app.keys[ibctransfertypes.StoreKey]).OpenKVStore(ctx))
	iter := storetypes.KVStorePrefixIterator(store, ibctransfertypes.DenomTraceKey)

	var (
		keys   [][]byte
		denoms []ibctransfertypes.Denom
	)
	for ; iter.Valid(); iter.Next() {
		path, base, err := decodeDenomTrace(iter.Value())
		if err != nil {
			_ = iter.Close()
			return nil, err
		}
		fullPath := base
		if path != "" {
			fullPath = path + "/" + base
		}
		hash := sha256.Sum256([]byte(fullPath))

		parsed := ibctransfertypes.ExtractDenomFromPath(fullPath)
		if parsed.Validate() == nil && bytes.Equal(parsed.Hash(), hash[:]) {
			continue // the ibc-go migration handles it
		}

		denom, err := denomFromTrace(path, base)
		if err != nil {
			_ = iter.Close()
			return nil, fmt.Errorf("denom trace %q: %w", fullPath, err)
		}
		if !bytes.Equal(denom.Hash(), hash[:]) {
			_ = iter.Close()
			return nil, fmt.Errorf("denom trace %q: converted denom %s changes the ibc denom", fullPath, denom.Path())
		}
		keys = append(keys, append([]byte{}, iter.Key()...))
		denoms = append(denoms, denom)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	for _, key := range keys {
		store.Delete(key)
	}
	return denoms, nil
}

// denomFromTrace builds a denom from a trace path made of "port/channel" pairs.
func denomFromTrace(path, base string) (ibctransfertypes.Denom, error) {
	var hops []ibctransfertypes.Hop
	if path != "" {
		parts := strings.Split(path, "/")
		if len(parts)%2 != 0 {
			return ibctransfertypes.Denom{}, fmt.Errorf("invalid trace path %q", path)
		}
		for i := 0; i < len(parts); i += 2 {
			hops = append(hops, ibctransfertypes.NewHop(parts[i], parts[i+1]))
		}
	}
	denom := ibctransfertypes.NewDenom(base, hops...)
	return denom, denom.Validate()
}

// decodeDenomTrace decodes an ibc.applications.transfer.v1.DenomTrace:
// path (field 1) and base_denom (field 2).
func decodeDenomTrace(bz []byte) (path, base string, err error) {
	for len(bz) > 0 {
		num, typ, n := protowire.ConsumeTag(bz)
		if n < 0 {
			return "", "", protowire.ParseError(n)
		}
		bz = bz[n:]
		if typ != protowire.BytesType {
			n = protowire.ConsumeFieldValue(num, typ, bz)
			if n < 0 {
				return "", "", protowire.ParseError(n)
			}
			bz = bz[n:]
			continue
		}
		v, n := protowire.ConsumeBytes(bz)
		if n < 0 {
			return "", "", protowire.ParseError(n)
		}
		bz = bz[n:]
		switch num {
		case 1:
			path = string(v)
		case 2:
			base = string(v)
		}
	}
	return path, base, nil
}

// setUpgradeStoreLoader applies the store upgrades of the pending upgrade
func (app *App) setUpgradeStoreLoader() {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}
	if upgradeInfo.Name != UpgradeName || app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}
	storeUpgrades := storetypes.StoreUpgrades{
		Deleted: removedModules,
	}
	// configure store loader that checks if version == upgradeHeight and applies store upgrades
	app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
}
