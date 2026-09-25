package app

import (
	"context"

	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// removedModules are the modules dropped by the v10 upgrade, together with
// their stores:
//   - alliance: shut down in v9.1.0
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

// RegisterUpgradeHandlers registers the upgrade handlers
func (app *App) RegisterUpgradeHandlers(cfg module.Configurator) {
	app.UpgradeKeeper.SetUpgradeHandler(UpgradeName, func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		for _, name := range removedModules {
			delete(fromVM, name)
		}
		return app.mm.RunMigrations(ctx, cfg, fromVM)
	})
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
