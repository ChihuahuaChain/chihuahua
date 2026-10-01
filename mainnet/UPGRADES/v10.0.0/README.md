# Chihuahua v10.0.0 Upgrade

The upgrade is scheduled for block **TBD**, expected around **TBD UTC** (block times vary, so it can come a little earlier or later). Countdown: https://explorer.chihuahua.wtf/block/TBD

The upgrade proposal is **expedited**: its voting period is 24 hours.

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

- `alliance` (retired in v9.5.0)
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
- `x/feeburn` keeps the total of burned transaction fees on chain, starting from the total published on burn.chihuahua.wtf at the upgrade: `chihuahuad query feeburn total-burned`, REST `/chihuahua/feeburn/total_burned`.
- CosmWasm stack updated to wasmd v0.70.4 and wasmvm v3.0.8.

## For integrators and contract developers

- The IBC transfer query `DenomTrace` (`/ibc.applications.transfer.v1.Query/DenomTrace`) no longer exists in ibc-go. Use `Denom` (`/ibc.applications.transfer.v1.Query/Denom`). Contracts using the `DenomTrace` stargate query must switch to `Denom`.
- The REST and gRPC endpoints of the removed modules (`alliance`, `params`, `crisis`, `nft`, `circuit`, IBC fee) are gone.
- The CosmWasm capabilities now go up to `cosmwasm_2_2`, plus `ibc2` for IBC v2.

## For node operators

- chihuahuad is now built with Go 1.26.
- `app.toml` and `config.toml` from v9 keep working. New options get their defaults; the experimental CometBFT features (LibP2P, AdaptiveSync) stay disabled unless enabled explicitly.
- The `--x-crisis-skip-assert-invariants` start flag no longer exists: remove it from your service files.
- wasmvm v3 uses a new compiled module cache in `data/wasm/cache/modules`: contracts are recompiled on their first execution after the upgrade, so the first calls to each contract are slower. The previous cache (the `*-wasmer7` directory, several GB) is no longer used and can be deleted once the node runs v10.

## Validator checklist

Before the upgrade height:

1. Vote on the upgrade proposal (expedited: 24 hours).
2. Get v10.0.0 ready, in one of the three ways below.
3. Make sure cosmovisor has `DAEMON_SHUTDOWN_GRACE=30s` and that `node` in `~/.chihuahuad/config/client.toml` points at your node's RPC.
4. Remove `--x-crisis-skip-assert-invariants` from the start command of your service, if present: v10 no longer accepts it.
5. Keep an eye on the node around the upgrade height: after the switch it should run v10.0.0 and keep signing.

Release binaries and sha256:

| Platform | sha256 |
|---|---|
| linux/amd64 | `TBD` |
| linux/arm64 | `TBD` |

## Upgrade with the Huahua Node Manager

Once the proposal has passed:

```bash
huahua-node upgrade
```

It finds the scheduled v10.0.0 upgrade, downloads the release, checks its sha256 and places it in cosmovisor, which switches at the upgrade height by itself.

## Upgrade with cosmovisor (automatic)

The upgrade proposal carries the release binaries (linux/amd64, linux/arm64) with their sha256. With cosmovisor v1.7, set in the service environment:

```ini
Environment="DAEMON_ALLOW_DOWNLOAD_BINARIES=true"
Environment="DAEMON_RESTART_AFTER_UPGRADE=true"
Environment="DAEMON_SHUTDOWN_GRACE=30s"
```

Keep `node` in `~/.chihuahuad/config/client.toml` pointed at your node's RPC (`tcp://localhost:26657`, or your RPC port): cosmovisor checks the height through it before switching binaries.

## Upgrade with cosmovisor (binary placed by hand)

Download the release binary and check it:

```bash
ARCH=amd64   # or arm64
curl -fLO https://github.com/ChihuahuaChain/chihuahua/releases/download/v10.0.0/chihuahuad_linux_$ARCH
curl -fsSL https://github.com/ChihuahuaChain/chihuahua/releases/download/v10.0.0/chihuahuad_sha256.txt | grep "chihuahuad_linux_$ARCH" | sha256sum -c
mkdir -p ~/.chihuahuad/cosmovisor/upgrades/v10.0.0/bin
install -m 755 chihuahuad_linux_$ARCH ~/.chihuahuad/cosmovisor/upgrades/v10.0.0/bin/chihuahuad
~/.chihuahuad/cosmovisor/upgrades/v10.0.0/bin/chihuahuad version   # v10.0.0
```

Or build it from source with Go 1.26:

```bash
cd chihuahua
git fetch --tags
git checkout v10.0.0
make install
chihuahuad version   # v10.0.0
mkdir -p ~/.chihuahuad/cosmovisor/upgrades/v10.0.0/bin
cp "$(which chihuahuad)" ~/.chihuahuad/cosmovisor/upgrades/v10.0.0/bin/
```

## Syncing from genesis

Apply v10.0.0 at height TBD, after v9.5.0 at height 25571500.

## After the upgrade

```bash
chihuahuad version                          # v10.0.0
chihuahuad q upgrade applied v10.0.0        # the height at which it was applied
chihuahuad q feeburn total-burned           # the burned fees total, now on chain
journalctl -u chihuahuad -f                 # blocks being committed again
```

The first calls to each contract are slower while wasmvm v3 recompiles it. If the node stops at the upgrade height and doesn't restart, check that `~/.chihuahuad/cosmovisor/upgrades/v10.0.0/bin/chihuahuad version` prints `v10.0.0`, then restart the service. Questions: the validators' channel on [Discord](https://discord.gg/chihuahuachain-878201449421619211).
