# Chihuahua v9.5.0 Upgrade

The upgrade is scheduled for block `TBD`. Countdown: https://explorer.chihuahua.wtf/block/TBD

**v9.5.0 is an intermediate upgrade.** It prepares the chain for **v10**, which follows a few days later: Cosmos SDK v0.54, IBC v11, CosmWasm 3.

## What changes

### Alliance retired

Following proposal #99, v9.5.0 retires the Alliance module, which has been inactive since the Migaloo/White Whale project wound down.

- The HUAHUA held by the module is reintroduced into the ecosystem.
- The staked ampGASH and the rewards left unclaimed in the module are burned: they no longer have any use or value.

The module is left empty, and v10 removes it for good.

### Bug fixes

- Stakedrops: fixed the start and end block checks, and the missing payout for the first block.
- Event attributes: `fee_payer` now reports the bech32 address; tokenfactory mint and burn events report the correct address.
- `chihuahuad export` works again.

## Upgrade with cosmovisor (automatic)

The upgrade proposal carries the release binaries (linux/amd64, linux/arm64) with their sha256. With cosmovisor v1.7, set in the service environment:

```ini
Environment="DAEMON_ALLOW_DOWNLOAD_BINARIES=true"
Environment="DAEMON_RESTART_AFTER_UPGRADE=true"
Environment="DAEMON_SHUTDOWN_GRACE=30s"
```

- `DAEMON_ALLOW_DOWNLOAD_BINARIES`: cosmovisor downloads the binary from the proposal and verifies its sha256.
- `DAEMON_SHUTDOWN_GRACE`: cosmovisor stops the old binary cleanly before starting the new one. Without it the old process is killed and the new one can fail to open the database.

Keep `node` in `~/.chihuahuad/config/client.toml` pointed at your node's RPC (`tcp://localhost:26657`, or your RPC port): cosmovisor checks the height through it before switching binaries.

## Upgrade with cosmovisor (manual)

Download the release binary and check it:

```bash
ARCH=amd64   # or arm64
curl -fLO https://github.com/ChihuahuaChain/chihuahua/releases/download/v9.5.0/chihuahuad_linux_$ARCH
curl -fsSL https://github.com/ChihuahuaChain/chihuahua/releases/download/v9.5.0/chihuahuad_sha256.txt | grep "chihuahuad_linux_$ARCH" | sha256sum -c
mkdir -p ~/.chihuahuad/cosmovisor/upgrades/v9.5.0/bin
install -m 755 chihuahuad_linux_$ARCH ~/.chihuahuad/cosmovisor/upgrades/v9.5.0/bin/chihuahuad
~/.chihuahuad/cosmovisor/upgrades/v9.5.0/bin/chihuahuad version   # v9.5.0
```

Or build it from source with Go 1.23.9:

```bash
cd chihuahua
git fetch --tags
git checkout v9.5.0
make install
chihuahuad version   # v9.5.0
mkdir -p ~/.chihuahuad/cosmovisor/upgrades/v9.5.0/bin
cp "$(which chihuahuad)" ~/.chihuahuad/cosmovisor/upgrades/v9.5.0/bin/
```

## Syncing from genesis

Apply v9.5.0 at the upgrade height above, after v9.0.6 (the v9.0.7 binary also runs up to it).
