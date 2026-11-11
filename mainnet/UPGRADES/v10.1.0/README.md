# Chihuahua v10.1.0 Upgrade

The upgrade is scheduled for block **26239240**, expected around **2026-11-17 14:00 UTC** (block times vary, so it can come a little earlier or later). Countdown: https://explorer.chihuahua.wtf/block/26239240

v10.1.0 is a coordinated upgrade on top of v10.0.1. It does not change the Cosmos SDK, CometBFT, ibc-go or wasmd versions, and the upgrade handler runs module migrations only: there is no state migration.

## What changes

- **Toolchain and dependencies**: built with Go 1.26.10; non-consensus dependencies updated.
- **Liquidity**: hardened batch execution; swap orders removed from a batch at execution are refunded; order prices are bounded against the pool reserves at submission; withdrawal shares are computed without an intermediate size limit; depleted pools are accepted in genesis import.
- **Token factory**: active stakedrops are kept across genesis export and import; stakedrop heights are bounds-checked; builder addresses are compared by bytes in `BurnFrom` and `ForceTransfer`.
- **IBC hooks**: wasm memos on incoming ICS-20 transfers are executed and packet lifecycle callbacks are delivered; an `ibc_callback` is registered only when it names the packet sender; timeout callbacks run in an isolated context.
- **Groups**: decimals accepted from input are limited to 34 significant digits and an exponent within ±34.
- **Node operations**: zero-height export works again.

## For integrators and contract developers

- Contracts that use IBC lifecycle callbacks must set `ibc_callback` to their own address; callbacks naming another address are ignored.
- Group member weights and decision policy values outside the limits above are rejected.
- No API, query or denom changes.

## For node operators

- `app.toml` and `config.toml` from v10.0.1 keep working; `minimum-gas-prices` stays at least `500uhuahua`.
- Standard cosmovisor flow, same as v10.0.1.

## Validator checklist

Before the upgrade height:

1. Vote on the upgrade proposal.
2. Get v10.1.0 ready, in one of the ways below.
3. Make sure cosmovisor has `DAEMON_SHUTDOWN_GRACE=30s` and that `node` in `~/.chihuahuad/config/client.toml` points at your node's RPC.
4. Keep an eye on the node around the upgrade height: after the switch it should run v10.1.0 and keep signing.

Release binaries and sha256:

| Platform | sha256 |
|---|---|
| linux/amd64 | `6a3c064ee9bd6e158ccd3c2d9c017094b14a6977d5cec47dee1d4b78a69c2b5d` |
| linux/arm64 | `756087b9f8411622d732e6b1429de478113630c0d84c22e3159d96c5b9da28b7` |

## Upgrade with the Huahua Node Manager

Once the proposal has passed:

```bash
huahua-node upgrade
```

It finds the scheduled v10.1.0 upgrade, downloads the release, checks its sha256 and places it in cosmovisor, which switches at the upgrade height by itself.

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
curl -fLO https://github.com/ChihuahuaChain/chihuahua/releases/download/v10.1.0/chihuahuad_linux_$ARCH
curl -fsSL https://github.com/ChihuahuaChain/chihuahua/releases/download/v10.1.0/chihuahuad_sha256.txt | grep "chihuahuad_linux_$ARCH" | sha256sum -c
mkdir -p ~/.chihuahuad/cosmovisor/upgrades/v10.1.0/bin
install -m 755 chihuahuad_linux_$ARCH ~/.chihuahuad/cosmovisor/upgrades/v10.1.0/bin/chihuahuad
~/.chihuahuad/cosmovisor/upgrades/v10.1.0/bin/chihuahuad version   # v10.1.0
```

Or build it from source with Go 1.26.10:

```bash
cd chihuahua
git fetch --tags
git checkout v10.1.0
make install
chihuahuad version   # v10.1.0
mkdir -p ~/.chihuahuad/cosmovisor/upgrades/v10.1.0/bin
cp "$(which chihuahuad)" ~/.chihuahuad/cosmovisor/upgrades/v10.1.0/bin/
```

## After the upgrade

```bash
chihuahuad version                          # v10.1.0
chihuahuad q upgrade applied v10.1.0        # the height at which it was applied
journalctl -u chihuahuad -f                 # blocks being committed again
```

If the node stops at the upgrade height and doesn't restart, check that `~/.chihuahuad/cosmovisor/upgrades/v10.1.0/bin/chihuahuad version` prints `v10.1.0`, then restart the service. Questions: the validators' channel on [Discord](https://discord.gg/chihuahuachain-878201449421619211).
