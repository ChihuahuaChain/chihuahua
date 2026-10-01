# Chihuahua v10 Upgrade

The Upgrade is scheduled for block `TBD`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/TBD)

This guide assumes that you use cosmovisor to manage upgrades.

## What changes

v10 brings Chihuahua to the current Cosmos stack:

| Component | v9 | v10 |
|---|---|---|
| Cosmos SDK | v0.50.15 | v0.54.4 |
| CometBFT | v0.38.23 | v0.39.4 |
| CosmWasm (wasmd / wasmvm) | v0.53.4 / v2.2.4 | v0.70.4 / v3.0.8 |
| IBC (ibc-go) | v8 | v11 |
| Go | 1.23 | 1.26 |

### Removed modules

The following modules are removed and their state is deleted by the upgrade:

- `alliance` (shut down in v9.1.0)
- `capability` and the IBC fee middleware (`feeibc`, ICS-29): removed from ibc-go
- `crisis`, `nft`, `circuit`: no longer maintained by the Cosmos SDK and unused on Chihuahua
- `params`: every module manages its own parameters

The historical governance proposals that used these modules can still be queried.

### Kept and maintained in this repository

- `x/liquidity`, previously imported from an external repository
- `x/group`: from SDK v0.54 the Cosmos SDK distributes it under a non-commercial license. Chihuahua keeps its own Apache-2.0 copy, and the existing groups keep working unchanged
- `x/ibc-hooks`

### Other changes

- IBC denom traces are migrated to the new ibc-go denom format. `ibc/...` denoms and balances do not change.
- The ante handler adds the CosmWasm simulation gas limit, transaction counter and gas register decorators, and the IBC redundant relay decorator.
- `chihuahuad bark`.
- CosmWasm security fixes: wasmd v0.70.4 and wasmvm v3.0.8 include the patches for the CosmWasm advisories CWA-2026-003, CWA-2026-004, CWA-2026-005 and CWA-2026-006, which affect wasmd v0.53.4 and wasmvm v2.2.4 used by v9.

## For integrators and contract developers

- The IBC transfer query `DenomTrace` (`/ibc.applications.transfer.v1.Query/DenomTrace`) no longer exists in ibc-go. Use `Denom` (`/ibc.applications.transfer.v1.Query/Denom`). Contracts using the `DenomTrace` stargate query must switch to `Denom`.
- The REST and gRPC endpoints of the removed modules (`alliance`, `params`, `crisis`, `nft`, `circuit`, IBC fee) are gone.
- The CosmWasm capabilities now go up to `cosmwasm_2_2`, plus `ibc2` for IBC v2.

## For node operators

- chihuahuad is now built with Go 1.26.
- `app.toml` and `config.toml` from v9 keep working. New options get their defaults; the experimental CometBFT features (LibP2P, AdaptiveSync) stay disabled unless enabled explicitly.
- The `--x-crisis-skip-assert-invariants` start flag no longer exists: remove it from your service files.
- wasmvm v3 uses a new compiled module cache in `data/wasm/cache/modules`: contracts are recompiled on their first execution after the upgrade, so the first calls to each contract are slower. The previous cache (the `*-wasmer7` directory, several GB) is no longer used and can be deleted once the node runs v10.

## If you are syncing from 0 you need to apply v10 at height TBD

```bash
# get the new version
cd chihuahua
git fetch --all
git checkout v10.0.0
make install
```

# check the version

```bash
# should be v10.0.0
chihuahuad version
# Should be commit TBD
chihuahuad version --long | grep commit
```

# Make new directory and copy binary

```bash
mkdir -p $HOME/.chihuahuad/cosmovisor/upgrades/v10/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahuad/cosmovisor/upgrades/v10/bin
```

# check the version again

```bash
# should be v10.0.0
$HOME/.chihuahuad/cosmovisor/upgrades/v10/bin/chihuahuad version
```
