// app.go

// Register the AppModule for the Interchain Accounts module and the authentication module
// Note: No `icaauth` exists, this must be substituted with an actual Interchain Accounts authentication module
ModuleBasics = module.NewBasicManager(
    ...
    ica.AppModuleBasic{},
    icaauth.AppModuleBasic{},
    ...
)

... 

// Add module account permissions for the Interchain Accounts module
// Only necessary for host chain functionality
// Each Interchain Account created on the host chain is derived from the module account created
maccPerms = map[string][]string{
    ...
    icatypes.ModuleName:            nil,
}

...

// Add Interchain Accounts Keepers for each submodule used and the authentication module
// If a submodule is being statically disabled, the associated Keeper does not need to be added. 
type App struct {
    ...

    ICAControllerKeeper icacontrollerkeeper.Keeper
    ICAHostKeeper       icahostkeeper.Keeper
    //ICAAuthKeeper       icaauthkeeper.Keeper

    ...
}

...

// Create store keys for each submodule Keeper and the authentication module
keys := sdk.NewKVStoreKeys(
    ...
    icacontrollertypes.StoreKey,
    icahosttypes.StoreKey,
    icaauthtypes.StoreKey,
    ...
)

... 

// Create the scoped keepers for each submodule keeper and authentication keeper
scopedICAControllerKeeper := app.CapabilityKeeper.ScopeToModule(icacontrollertypes.SubModuleName)
scopedICAHostKeeper := app.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
//scopedICAAuthKeeper := app.CapabilityKeeper.ScopeToModule(icaauthtypes.ModuleName)

...

// Create the Keeper for each submodule
app.ICAControllerKeeper = icacontrollerkeeper.NewKeeper(
    appCodec, keys[icacontrollertypes.StoreKey], app.GetSubspace(icacontrollertypes.SubModuleName),
    app.IBCKeeper.ChannelKeeper, // may be replaced with middleware such as ics29 fee
    app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
    scopedICAControllerKeeper, app.MsgServiceRouter(),
)
app.ICAHostKeeper = icahostkeeper.NewKeeper(
    appCodec, keys[icahosttypes.StoreKey], app.GetSubspace(icahosttypes.SubModuleName),
    app.IBCKeeper.ChannelKeeper, // may be replaced with middleware such as ics29 fee
    app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
    app.AccountKeeper, scopedICAHostKeeper, app.MsgServiceRouter(),
)

// Create Interchain Accounts AppModule
icaModule := ica.NewAppModule(&app.ICAControllerKeeper, &app.ICAHostKeeper)

// Create your Interchain Accounts authentication module
//app.ICAAuthKeeper = icaauthkeeper.NewKeeper(appCodec, keys[icaauthtypes.StoreKey], app.MsgServiceRouter())

// ICA auth AppModule
//icaAuthModule := icaauth.NewAppModule(appCodec, app.ICAAuthKeeper)

// Create controller IBC application stack and host IBC module as desired
icaControllerStack := icacontroller.NewIBCMiddleware(nil, app.ICAControllerKeeper)
icaHostIBCModule := icahost.NewIBCModule(app.ICAHostKeeper)

// Register host and authentication routes
ibcRouter.
    AddRoute(icacontrollertypes.SubModuleName, icaControllerStack).
    AddRoute(icahosttypes.SubModuleName, icaHostIBCModule)
...

// Register Interchain Accounts and authentication module AppModule's
app.moduleManager = module.NewManager(
    ...
    icaModule,
    //icaAuthModule,
)

...

// Add Interchain Accounts to begin blocker logic
app.moduleManager.SetOrderBeginBlockers(
    ...
    icatypes.ModuleName,
    ...
)

// Add Interchain Accounts to end blocker logic
app.moduleManager.SetOrderEndBlockers(
    ...
    icatypes.ModuleName,
    ...
)

// Add Interchain Accounts module InitGenesis logic
app.moduleManager.SetOrderInitGenesis(
    ...
    icatypes.ModuleName,
    ...
)

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramskeeper.Keeper {
    ...
    paramsKeeper.Subspace(icahosttypes.SubModuleName)
    paramsKeeper.Subspace(icacontrollertypes.SubModuleName)
    ...
}
