<p align="center">
  <img src="https://github.com/ChihuahuaChain/resources/blob/main/logo/logo_transparent_notext.png?raw=true" width="140" alt="Chihuahua">
</p>

<h1 align="center">Chihuahua</h1>

<p align="center">
  The Cosmos meme coin that actually built a chain.<br>
  <a href="https://chihuahua.wtf">chihuahua.wtf</a> ·
  <a href="https://explorer.chihuahua.wtf">Explorer</a> ·
  <a href="https://github.com/ChihuahuaChain/huahua-node-manager">Node Manager</a> ·
  <a href="https://t.me/chihuahua_wtf">Telegram</a> ·
  <a href="https://discord.gg/chihuahuachain-878201449421619211">Discord</a> ·
  <a href="https://x.com/ChihuahuaChain">X</a>
</p>

<p align="center">
  <a href="https://github.com/ChihuahuaChain/chihuahua/releases/latest"><img src="https://img.shields.io/github/v/release/ChihuahuaChain/chihuahua?color=f5c518&label=release" alt="Latest release"></a>
  <img src="https://img.shields.io/badge/chain-chihuahua--1-f5c518" alt="chihuahua-1">
  <img src="https://img.shields.io/badge/go-1.23.9-555" alt="Go 1.23.9">
</p>

Chihuahua is a proof-of-stake blockchain built with the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk), live since 2021. `chihuahuad` is its node. On top of the usual Cosmos modules (staking, governance, IBC) the chain has:

- **permissionless CosmWasm**: anyone can upload and run smart contracts;
- **fee burn** (`x/feeburn`): a share of every transaction fee is burned;
- **token factory** (`x/tokenfactory`): anyone can create a native token;
- **stakedrops**: airdrops paid block by block to HUAHUA stakers;
- **liquidity pools** (`x/liquidity`).

## Network

| Field | Value |
|---|---|
| Chain ID | `chihuahua-1` |
| Token | HUAHUA, base denom `uhuahua` (1 HUAHUA = 1,000,000 uhuahua) |
| Address prefix | `chihuahua` |
| Minimum gas price | `100uhuahua` (wallets default to `1250uhuahua`) |
| RPC | https://rpc.chihuahua.wtf |
| REST | https://api.chihuahua.wtf |
| Explorer | https://explorer.chihuahua.wtf |
| Snapshots | https://snapshots.chihuahua.wtf |
| Current release | see [Releases](https://github.com/ChihuahuaChain/chihuahua/releases) and [Chain upgrades](#chain-upgrades) |

## Run a node or a validator

The quickest way is the **[Huahua Node Manager](https://github.com/ChihuahuaChain/huahua-node-manager)**. One command sets up a full node or a validator on Linux (amd64 or arm64):

```sh
bash <(curl -fsSL https://raw.githubusercontent.com/ChihuahuaChain/huahua-node-manager/main/huahua-node.sh)
```

Choose **Easy**, answer three questions, and a few minutes later the node is:

- synced from the latest snapshot, or with state sync;
- running under cosmovisor, so it upgrades itself at each upgrade height;
- started at boot, as a systemd service or a Docker container.

It downloads the release binary and checks every download against its sha256, finds live peers, and guides you through `create-validator`. Run the same command again, or `huahua-node`, to open its live dashboard.

To read the script before running it:

```sh
curl -fsSLO https://raw.githubusercontent.com/ChihuahuaChain/huahua-node-manager/main/huahua-node.sh
less huahua-node.sh
bash huahua-node.sh
```

The rest of this section is the manual way.

### 1. Get `chihuahuad`

**Release binary.** Every [release](https://github.com/ChihuahuaChain/chihuahua/releases) has static Linux binaries for amd64 and arm64, with their sha256 in `chihuahuad_sha256.txt`:

```sh
VERSION=v9.0.7
ARCH=amd64   # or arm64
curl -fLO https://github.com/ChihuahuaChain/chihuahua/releases/download/$VERSION/chihuahuad_linux_$ARCH
curl -fsSL https://github.com/ChihuahuaChain/chihuahua/releases/download/$VERSION/chihuahuad_sha256.txt | grep "chihuahuad_linux_$ARCH" | sha256sum -c
install -m 755 chihuahuad_linux_$ARCH ~/go/bin/chihuahuad   # or any directory in your PATH
chihuahuad version
```

**From source.** You need Go 1.23.9 and a C toolchain (`make`, `gcc`):

```sh
git clone https://github.com/ChihuahuaChain/chihuahua.git
cd chihuahua
git checkout v9.0.7
make install
chihuahuad version   # v9.0.7
```

Build with the exact Go version: a different one can produce a different app hash.

### 2. Configure

```sh
chihuahuad init <your-node-name> --chain-id chihuahua-1
curl -fsSL https://raw.githubusercontent.com/ChihuahuaChain/chihuahua/main/mainnet/genesis.json -o ~/.chihuahuad/config/genesis.json
sed -i 's/^minimum-gas-prices *=.*/minimum-gas-prices = "100uhuahua"/' ~/.chihuahuad/config/app.toml
```

Add a few live peers to `persistent_peers` in `~/.chihuahuad/config/config.toml`. Peers change over time: ask in [Discord](https://discord.gg/chihuahuachain-878201449421619211), or let the Node Manager find them.

### 3. Sync

Syncing from genesis takes a very long time, because it replays every upgrade. Start from a snapshot instead. [snapshots.chihuahua.wtf](https://snapshots.chihuahua.wtf/latest.json) publishes one every 12 hours, with its height and sha256:

```sh
meta=$(curl -fsSL https://snapshots.chihuahua.wtf/latest.json)
file=$(echo "$meta" | jq -r .file)
curl -fLO "https://snapshots.chihuahua.wtf/$file"
echo "$(echo "$meta" | jq -r .sha256)  $file" | sha256sum -c
cp ~/.chihuahuad/data/priv_validator_state.json /tmp/   # keep your signing state
rm -rf ~/.chihuahuad/data ~/.chihuahuad/wasm
lz4 -dc "$file" | tar -x -C ~/.chihuahuad
cp /tmp/priv_validator_state.json ~/.chihuahuad/data/
```

The snapshot includes the contracts' code (`data/wasm`), which a CosmWasm node needs.

### 4. Run it with cosmovisor

[Cosmovisor](https://github.com/cosmos/cosmos-sdk/tree/main/tools/cosmovisor) runs `chihuahuad` and swaps the binary when an upgrade height is reached.

```sh
go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@latest
mkdir -p ~/.chihuahuad/cosmovisor/genesis/bin
cp "$(which chihuahuad)" ~/.chihuahuad/cosmovisor/genesis/bin/
```

`/etc/systemd/system/chihuahuad.service`, with your user in place of `<user>`:

```ini
[Unit]
Description=Chihuahua node (cosmovisor)
After=network-online.target

[Service]
User=<user>
ExecStart=/home/<user>/go/bin/cosmovisor run start
Restart=always
RestartSec=3
LimitNOFILE=65535
Environment="DAEMON_NAME=chihuahuad"
Environment="DAEMON_HOME=/home/<user>/.chihuahuad"
Environment="DAEMON_ALLOW_DOWNLOAD_BINARIES=true"
Environment="DAEMON_RESTART_AFTER_UPGRADE=true"

[Install]
WantedBy=multi-user.target
```

```sh
sudo systemctl daemon-reload
sudo systemctl enable --now chihuahuad
journalctl -u chihuahuad -f
```

With `DAEMON_ALLOW_DOWNLOAD_BINARIES=true`, cosmovisor downloads the new binary from the upgrade proposal and checks its sha256 (see [Chain upgrades](#chain-upgrades)). Before switching, cosmovisor 1.7 checks the height with `chihuahuad status`, which asks the RPC set as `node` in `~/.chihuahuad/config/client.toml`: keep it pointed at this node (`tcp://localhost:26657`, or your RPC port). To place it yourself instead, put it in `~/.chihuahuad/cosmovisor/upgrades/<upgrade name>/bin/` before the upgrade height.

### 5. Become a validator

Wait until the node has caught up (`chihuahuad status | jq .sync_info.catching_up` shows `false`), then create a key and fund it with a little HUAHUA for the transaction and your self-delegation:

```sh
chihuahuad keys add validator   # write the mnemonic down somewhere safe
```

`validator.json`:

```json
{
  "pubkey": <the output of: chihuahuad comet show-validator>,
  "amount": "1000000uhuahua",
  "moniker": "<your validator name>",
  "identity": "<keybase id, optional>",
  "website": "<website, optional>",
  "security": "<security contact, optional>",
  "details": "<description, optional>",
  "commission-rate": "0.10",
  "commission-max-rate": "0.20",
  "commission-max-change-rate": "0.01",
  "min-self-delegation": "1"
}
```

```sh
chihuahuad tx staking create-validator validator.json \
  --from validator --chain-id chihuahua-1 \
  --gas auto --gas-adjustment 1.4 --gas-prices 1250uhuahua
```

The chain enforces a minimum commission of 5% ([proposal 3](https://explorer.chihuahua.wtf/proposal/3)).

Back up `~/.chihuahuad/config/priv_validator_key.json` offline. Never run two nodes with the same key at the same time: double signing gets the validator slashed and jailed for good.

Then say hello in [Discord](https://discord.gg/chihuahuachain-878201449421619211): the validators' channel is where upgrades are coordinated.

## Chain upgrades

Upgrades are voted on chain. Each software upgrade proposal names the upgrade, its height, and the release binaries with their sha256, so cosmovisor can install them by itself. The upgrade notes are in [`mainnet/UPGRADES`](mainnet/UPGRADES), and the explorer shows a countdown to any block (`https://explorer.chihuahua.wtf/block/<height>`).

| Upgrade | Binary | Height | Date (UTC) |
|---|---|---|---|
| — | v9.0.7 | no upgrade height: patch release (Cosmos SDK v0.50.15, CometBFT v0.38.23) | 2026-05-13 |
| [v9.0.6](mainnet/UPGRADES/v9.0.6) | v9.0.6 | 20,523,000 | 2025-10-29 13:40 |
| [v9.0.5](mainnet/UPGRADES/v9.0.5) | v9.0.5 | 18,504,000 | 2025-06-16 13:20 |
| [v9.0.4](mainnet/UPGRADES/v9.0.4) | v9.0.4 | 18,385,000 | 2025-06-08 13:00 |
| — | v9.0.3 | no upgrade height: replaces v9.0.2 | |
| [v9.0.2](mainnet/UPGRADES/v9.0.2) | v9.0.2 | 17,073,000 | 2025-03-12 13:56 |
| [v9.0.1](mainnet/UPGRADES/v9.0.1) | v9.0.1 | 16,623,000 | 2025-02-11 01:49 |
| [v9.0.0](mainnet/UPGRADES/v9.0.0) | v9.0.0 | 16,529,000 | 2025-02-03 13:11 |
| [v8.0.2](mainnet/UPGRADES/v8.0.2) | v8.0.2 | 15,103,000 | 2024-10-22 13:50 |
| [v8.0.0](mainnet/UPGRADES/v8.0.0) | v8.0.0 | 14,762,000 | 2024-09-29 13:45 |
| [v7.0.1](mainnet/UPGRADES/v7.0.1) | v7.0.2 | 13,250,000 | 2024-06-17 15:11 |
| [v7](mainnet/UPGRADES/v7) | v7 | 12,900,000 | 2024-05-24 10:12 |
| [v6](mainnet/UPGRADES/v6) | v6 | 10,666,000 | 2023-12-23 15:30 |
| [v503](mainnet/UPGRADES/v503) | v5.0.4 | 9,430,000 | 2023-09-28 16:30 |
| [v502](mainnet/UPGRADES/v502) | v5.0.2 | 9,180,000 | 2023-09-11 15:00 |
| [v501](mainnet/UPGRADES/v501) | v5.0.1 | 8,813,000 | 2023-08-17 16:03 |
| [v500](mainnet/UPGRADES/v500) | v5.0.0 | 8,711,111 | 2023-08-10 15:28 |
| [v421](mainnet/UPGRADES/v421) | v4.2.1 | 6,376,376 | 2023-03-04 12:20 |
| [v420](mainnet/UPGRADES/v420) | v4.2.0 | 6,039,999 | 2023-02-09 13:41 |
| [v410](mainnet/UPGRADES/v410) | v4.1.0 | 4,886,666 | 2022-11-22 15:25 |
| [v400](mainnet/UPGRADES/v400) | v4.0.0 | 4,787,878 | 2022-11-15 13:36 |
| [v310](mainnet/UPGRADES/v310) | v3.1.0 | 4,673,333 | 2022-11-07 15:11 |
| iavl fast node | v2.4.x | | |
| [burnmech](mainnet/UPGRADES/burnmech) | v2.2.2 (retracted) | 4,488,444 | 2022-10-21 14:14 |
| [authz](mainnet/UPGRADES/authz) | v2.1.0 | 4,182,410 | 2022-09-30 14:54 |
| [minpropdeposit](mainnet/UPGRADES/minpropdeposit) | v2.0.2 | 3,654,321 | 2022-08-25 13:00 |
| [chiwawasm](mainnet/UPGRADES/chiwawasm) | v2.0.1 | 3,000,800 | 2022-07-11 17:02 |
| [angryandy](mainnet/UPGRADES/angryandy) | v1.1.1 | 535,000 | 2022-01-19 17:20 |

Notes:

- **v503:** the chain halted at 9,431,130 because of a wrong wasmd/wasmvm; use the v5.0.4 binary in the `v503` directory.
- **iavl fast node (v2.4.x):** state-breaking upgrade that completed the Dragonberry patch, moved to go1.19 and configured the iavl fast node.
- **burnmech:** v2.3.0 was never used on chain; it is kept for archive nodes.

To sync from genesis, run each binary up to its upgrade height, in order.

## Development

```sh
make install          # build and install chihuahuad (needs Go 1.23.9)
go test ./...         # unit and integration tests
make proto-gen        # regenerate the protobuf code (needs Docker)
scripts/test_node.sh  # single-node local chain, with funded test keys
```

### Releases

Pushing a `v*` tag runs the [release workflow](.github/workflows/release.yml). It:

1. builds static `chihuahuad` binaries for linux/amd64 and linux/arm64;
2. checks that each one is static and reports the tag as its version;
3. drafts a GitHub release with the binaries, `chihuahuad_sha256.txt` and `upgrade-info.json`.

`upgrade-info.json` is the cosmovisor download info: each binary's URL with its sha256. To propose the upgrade:

```sh
scripts/upgrade-proposal.sh -n <upgrade name> -t <tag> -s summary.md -a <target UTC time>
chihuahuad tx gov submit-proposal proposal.json --from <key> --chain-id chihuahua-1 --gas auto --gas-adjustment 1.4 --gas-prices 1250uhuahua
```

The script downloads the binaries, checks them against `upgrade-info.json`, estimates the height for the target time, and writes `proposal.json`. Publish the draft release before submitting the proposal, or cosmovisor can't download the binaries.

## Community

- Website: [chihuahua.wtf](https://chihuahua.wtf)
- Telegram: [t.me/chihuahua_wtf](https://t.me/chihuahua_wtf)
- Discord: [ChihuahuaChain](https://discord.gg/chihuahuachain-878201449421619211)
- X: [@ChihuahuaChain](https://x.com/ChihuahuaChain)
- Medium: [@chihuahuachain](https://medium.com/@chihuahuachain)
