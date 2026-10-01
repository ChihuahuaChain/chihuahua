# Chihuahua v9.5.0 Upgrade

The upgrade is scheduled for block **25571500**, expected around **2026-10-02 12:37 UTC** (block times vary, so it can come a little earlier or later). Countdown: https://explorer.chihuahua.wtf/block/25571500

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

## Validator checklist

Before height 25571500:

1. Vote on the upgrade proposal.
2. Get v9.5.0 ready, in one of the three ways below.
3. Make sure cosmovisor has `DAEMON_SHUTDOWN_GRACE=30s` and that `node` in `~/.chihuahuad/config/client.toml` points at your node's RPC.
4. Keep an eye on the node around the upgrade height: after the switch it should run v9.5.0 and keep signing.

Release binaries and sha256:

| Platform | sha256 |
|---|---|
| linux/amd64 | `8a06e6c3b9e5577ac1f0374fef15cfb87880576406aef41a9f504e7894bb62cb` |
| linux/arm64 | `18e9839fd3ffb051d6ea22b6f84faa958e40dc6161eec6cb58f5604684763027` |

## Upgrade with the Huahua Node Manager

Once the proposal has passed:

```bash
huahua-node upgrade
```

It finds the scheduled v9.5.0 upgrade, downloads the release, checks its sha256 and places it in cosmovisor, which switches at the upgrade height by itself.

Nodes installed before 2026-09-26 don't have `DAEMON_SHUTDOWN_GRACE` in their service yet. Add it once, before the upgrade:

```bash
sudo sed -i '/^Environment=DAEMON_RESTART_AFTER_UPGRADE=true/a Environment=DAEMON_SHUTDOWN_GRACE=30s' /etc/systemd/system/chihuahuad.service
sudo systemctl daemon-reload && sudo systemctl restart chihuahuad
```

With Docker, add `- DAEMON_SHUTDOWN_GRACE=30s` to the `environment` of the node in its `docker-compose.yml`, then `docker compose up -d`.

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

## Upgrade with cosmovisor (binary placed by hand)

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

Apply v9.5.0 at height 25571500, after v9.0.6 (the v9.0.7 binary also runs up to it).

## After the upgrade

```bash
chihuahuad version                          # v9.5.0
chihuahuad q upgrade applied v9.5.0         # the height at which it was applied
journalctl -u chihuahuad -f                 # blocks being committed again
```

If the node stops at the upgrade height and doesn't restart, check that `~/.chihuahuad/cosmovisor/upgrades/v9.5.0/bin/chihuahuad version` prints `v9.5.0`, then restart the service. Questions: the validators' channel on [Discord](https://discord.gg/chihuahuachain-878201449421619211).
