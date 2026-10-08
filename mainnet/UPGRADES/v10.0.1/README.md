# Chihuahua v10.0.1 Upgrade

The upgrade is scheduled for block **25742780**, expected around **2026-10-14 12:00 UTC** (block times vary, so it can come a little earlier or later). Countdown: https://explorer.chihuahua.wtf/block/25742780

v10.0.1 is a coordinated upgrade on top of v10.0.0. It adds protocol features and hardening across x/feeburn, x/liquidity and CosmWasm. It does not change the Cosmos SDK, CometBFT, ibc-go or wasmd versions shipped in v10.0.0.

## What changes

### Protocol minimum gas price

The chain now enforces a minimum gas price **in consensus**, set by governance and read from x/feeburn parameters. Until now the minimum was only a per-validator mempool setting, so it varied from node to node; the floor is the same on every node and is applied both when a transaction enters the mempool and when the block executes.

- The upgrade sets the floor to **500 uhuahua per unit of gas**.
- Governance can change it at any time with `MsgUpdateParams` on x/feeburn (native denom only).
- A validator's local `minimum-gas-prices` may still be set higher; it can no longer be *effectively* lower, because the protocol floor wins.

### Public burn address

The x/feeburn module account becomes the chain's **public burn address**. Anything sent to it, by a community-pool spend or by anyone else, is permanently removed from the supply at the end of the block and added to the on-chain burned total:

- Burnable: **uhuahua** and any **token factory** denom (locally minted, native tokens).
- Left untouched: IBC vouchers (`ibc/...`) and every other denom, so a burn never strands an asset escrowed on another chain.
- The burned total is on chain: `chihuahuad query feeburn total-burned`, REST `/chihuahua/feeburn/total_burned`, and feeds the figure on burn.chihuahua.wtf.

### Larger CosmWasm contracts

The maximum contract size accepted by a store-code message goes from **800 KiB to 1600 KiB**. This simply makes room for bigger, more capable contracts — including something we are building that we think you will enjoy. More on that soon. 🐾

### Liquidity robustness

Pool batch execution (swaps, deposits, withdrawals) is now **atomic** and fails safely: a batch that cannot execute is rolled back and refunded rather than interrupting block production. One misbehaving batch can no longer stop the chain.

### Stability and consensus hardening

The release also carries internal fixes that make upgrades, state migrations and node recovery more robust. These are transparent to users and integrators.

## For integrators and contract developers

- **Minimum fee:** a transaction must pay at least the protocol floor (currently `500uhuahua` per unit of gas, i.e. `ceil(500 × gasLimit)` uhuahua). Set your clients and relayers to use a gas price at or above the floor, or transactions will be rejected with `insufficient fee`.
- **Contract size:** store-code now accepts contracts up to 1600 KiB.
- No API, query or denom changes beyond the above.

## For node operators

- Set `minimum-gas-prices` in `app.toml` to at least `500uhuahua` (or pass `--minimum-gas-prices`). The node prints a warning at startup if the local value is below the protocol floor; a lower local value only makes this node's mempool accept transactions that then fail at block execution.
- `app.toml` and `config.toml` from v10.0.0 keep working otherwise.
- Standard cosmovisor flow, same as v10.0.0.

## Validator checklist

Before the upgrade height:

1. Vote on the upgrade proposal.
2. Get v10.0.1 ready, in one of the ways below.
3. Make sure cosmovisor has `DAEMON_SHUTDOWN_GRACE=30s` and that `node` in `~/.chihuahuad/config/client.toml` points at your node's RPC.
4. Set `minimum-gas-prices` to at least `500uhuahua`.
5. Keep an eye on the node around the upgrade height: after the switch it should run v10.0.1 and keep signing.

Release binaries and sha256:

| Platform | sha256 |
|---|---|
| linux/amd64 | `f6be02c6f5118cfa44ad5d85a9f461960990781001746f88be5cf3dad174da4f` |
| linux/arm64 | `926cbda1610a21f6dc913976a0ffe338cedf5379d2ee1a40cabedea7b8f72540` |

## Upgrade with the Huahua Node Manager

Once the proposal has passed:

```bash
huahua-node upgrade
```

It finds the scheduled v10.0.1 upgrade, downloads the release, checks its sha256 and places it in cosmovisor, which switches at the upgrade height by itself.

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
curl -fLO https://github.com/ChihuahuaChain/chihuahua/releases/download/v10.0.1/chihuahuad_linux_$ARCH
curl -fsSL https://github.com/ChihuahuaChain/chihuahua/releases/download/v10.0.1/chihuahuad_sha256.txt | grep "chihuahuad_linux_$ARCH" | sha256sum -c
mkdir -p ~/.chihuahuad/cosmovisor/upgrades/v10.0.1/bin
install -m 755 chihuahuad_linux_$ARCH ~/.chihuahuad/cosmovisor/upgrades/v10.0.1/bin/chihuahuad
~/.chihuahuad/cosmovisor/upgrades/v10.0.1/bin/chihuahuad version   # v10.0.1
```

Or build it from source with Go 1.26:

```bash
cd chihuahua
git fetch --tags
git checkout v10.0.1
make install
chihuahuad version   # v10.0.1
mkdir -p ~/.chihuahuad/cosmovisor/upgrades/v10.0.1/bin
cp "$(which chihuahuad)" ~/.chihuahuad/cosmovisor/upgrades/v10.0.1/bin/
```

## After the upgrade

```bash
chihuahuad version                          # v10.0.1
chihuahuad q upgrade applied v10.0.1        # the height at which it was applied
chihuahuad q feeburn params                 # min_gas_prices is now set
chihuahuad q feeburn total-burned           # the burned total, now including the burn address
journalctl -u chihuahuad -f                 # blocks being committed again
```

If the node stops at the upgrade height and doesn't restart, check that `~/.chihuahuad/cosmovisor/upgrades/v10.0.1/bin/chihuahuad version` prints `v10.0.1`, then restart the service. Questions: the validators' channel on [Discord](https://discord.gg/chihuahuachain-878201449421619211).
