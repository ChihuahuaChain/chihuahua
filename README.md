<img alt='ChihuahuaChain logo' src="https://github.com/ChihuahuaChain/resources/blob/main/logo/logo_transparent_notext.png?raw=true" width="150"/>

# ChihuahuaChain
##### _The Cosmos MEME Coin_
Stay up to date with the latest news on our Socials
 - Join our [Telegram Community](https://t.me/chihuahua_cosmos)
 - Join our [Discord](https://discord.gg/chihuahuachain-878201449421619211)
 - Follow us on [Twitter](https://twitter.com/ChihuahuaChain)
 - Check out our [Medium](https://medium.com/@chihuahuachain)

# Node Installation

- #### Install Prerequisites

```bash
# update the local package list and install any available upgrades 
sudo apt-get update && sudo apt upgrade -y 

# install toolchain and ensure accurate time synchronization 
sudo apt-get install make build-essential gcc git jq chrony -y
```

- #### Install Go

```bash
# download the latest version
wget https://go.dev/dl/go1.21.1.linux-amd64.tar.gz

# remove old version (if any)
sudo rm -rf /usr/local/go

# install the new version
sudo tar -C /usr/local -xzf go1.21.1.linux-amd64.tar.gz
```

- #### Configure Environmental Variables
```bash
# run these commands
cat <<EOF >> ~/.profile
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export GO111MODULE=on
export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
EOF

source ~/.profile

go version
```
The output should be `go version go1.21.1 linux/amd64`

- #### Install Chihuahua from sources

```bash
# run these commands
git clone https://github.com/ChihuahuaChain/chihuahua.git
cd chihuahua
git fetch --tags
git checkout v7.0.2
make install
```

To verify the installation you can run `chihuahuad version` and it should return `v7.0.2`

- #### Initialize the Chain
Replace `$MONIKERNAME` with your choosen node name

`chihuahuad init $MONIKER_NAME --chain-id chihuahua-1`

- #### Download the Genesis

```bash
wget -O ~/.chihuahuad/config/genesis.json https://raw.githubusercontent.com/ChihuahuaChain/chihuahua/main/mainnet/genesis.json
```

- #### Add Seeds & Persistent Peers

```bash
seeds="77cbb35d1df17f48a42e9f157f12f55b691e9f5e@seeds.goldenratiostaking.net:1620,4936e377b4d4f17048f8961838a5035a4d21240c@chihuahua-seed-01.mercury-nodes.net:29540"
peers="b140eb36b20f3d201936c4757d5a1dcbf03a42f1@216.238.79.138:26656,19900e1d2b10be9c6672dae7abd1827c8e1aad1e@161.97.96.253:26656,c382a9a0d4c0606d785d2c7c2673a0825f7c53b2@88.99.94.120:26656,a5dfb048e4ed5c3b7d246aea317ab302426b37a1@137.184.250.180:26656,3bad0326026ca4e29c64c8d206c90a968f38edbe@128.199.165.78:26656,89b576c3eb72a4f0c66dc0899bec7c21552ea2a5@23.88.7.73:29538,38547b7b6868f93af1664d9ab0e718949b8853ec@54.184.20.240:30758,a9640eb569620d1f7be018a9e1919b0357a18b8c@38.146.3.160:26656,7e2239a0d4a0176fe4daf7a3fecd15ac663a8eb6@144.91.126.23:26656"
sed -i.bak -e "s/^seeds *=.*/seeds = \"$seeds\"/; s/^persistent_peers *=.*/persistent_peers = \"$peers\"/" ~/.chihuahuad/config/config.toml
```

- #### Update minimum-gas-price in app.toml

```bash
sed -i.bak 's/minimum-gas-prices =.*/minimum-gas-prices = "100uhuahua"/' $HOME/.chihuahuad/config/app.toml
```

- #### Setting up Cosmovisor

Install cosmovisor 
```bash
go install github.com/cosmos/cosmos-sdk/cosmovisor/cmd/cosmovisor@latest

which cosmovisor

# should return 
'/home/<your-user>/go/bin/cosmovisor'

# run these commands
cat <<EOF >> ~/.profile
export DAEMON_NAME=chihuahuad
export DAEMON_HOME=$HOME/.chihuahuad
EOF

source ~/.profile

echo $DAEMON_NAME

# should return
'chihuahuad'

# create the directories
mkdir -p $DAEMON_HOME/cosmovisor/genesis/bin
mkdir -p $DAEMON_HOME/cosmovisor/upgrades

# check the binary path with
which chihuahuad

# this should return
'/home/your-user/go/bin/chihuahuad'

# copy the binary into
cp $(which chihuahuad) $DAEMON_HOME/cosmovisor/genesis/bin
```
Set up the service file

```bash
sudo nano /etc/systemd/system/chihuahuad.service

# paste and edit <your-user> with your username
[Unit]
Description=Chihuahua Daemon (cosmovisor)
After=network-online.target

[Service]
User=<your-user>
ExecStart=/home/<your-user>/go/bin/cosmovisor run start
Restart=always
RestartSec=3
LimitNOFILE=4096
Environment="DAEMON_NAME=chihuahuad"
Environment="DAEMON_HOME=/home/<your-user>/.chihuahuad"
Environment="DAEMON_ALLOW_DOWNLOAD_BINARIES=false"
Environment="DAEMON_RESTART_AFTER_UPGRADE=true"
Environment="DAEMON_LOG_BUFFER_SIZE=512"

[Install]
WantedBy=multi-user.target
```

Enable the service

```bash
sudo -S systemctl daemon-reload
sudo -S systemctl enable chihuahuad
```

Get the latest [snapshot](https://polkachu.com/tendermint_snapshots/chihuahua) (_Thanks to [Polkachu](https://twitter.com/polka_chu)_) and follow the Pruning tips to save some GB

- #### Start the node

You can start the node by running
```bash
sudo systemctl start chihuahuad

# check the logs by running
journalctl -u chihuahuad -f
```
The node will take some time to catch-up with the blockchain.
You can follow the blocks being indexed by rnning

```bash
journalctl -u chihuahuad -f | grep indexed
```

# Join the Validators _(mainnet)_

ChihuahuaChain Governance [voted a proposal](https://www.mintscan.io/chihuahua/proposals/3) enabling the minimum 5% Commission enforced by the blockchain.

```bash
# create a new wallet for the validator

chihuahuad keys add <key-name>

# save the seed phrase (mnemonic) in a safe place
# copy the 'chihuahua...' address and send some HUAHUA
# in order to pay for the validator creation's transaction

# Make sure the Validator has fully synced before running 
chihuahuad tx staking create-validator \
  --from "<key-name>" \
  --amount "1000000uhuahua" \
  --pubkey "$(chihuahuad tendermint show-validator)" \
  --chain-id "chihuahua-1" \
  --moniker "<moniker>" \
  --commission-max-change-rate 0.01 \
  --commission-max-rate 0.20 \
  --commission-rate 0.10 \
  --min-self-delegation 1 \
  --details "<details>" \
  --security-contact "<contact>" \
  --website "<website>" \
  --gas-prices "1uhuahua"
  
# Make sure to backup the priv_validator_key.json file in your
# /home/<your-user>/.chihuahuad/config directory
# and store it in a safe place
```

**Congratulation!** Your Validator node should be up and running

_Make sure to join our [Discord](https://discord.gg/chihuahuachain-878201449421619211) and contact a moderator if you have a mainnet node so we can invite you to the validator's channel to follow up the latest updates and future upgrades._

---

# Chain Upgrades

- **v7.0.1** _(using v7.0.2 binary)_ - Block 13250000 - (2024-06-17 15:11:00)
  - [Upgrade Instruction](https://github.com/ChihuahuaChain/chihuahua/blob/main/mainnet/UPGRADES/v7.0.1)
- **v7** _(v7)_ - Block 12900000 - (2024-05-24 10:12:00)
  - [Upgrade Instruction](https://github.com/ChihuahuaChain/chihuahua/blob/main/mainnet/UPGRADES/v7)
- **v6** _(v6)_ - Block 10666000 - (2023-12-23 15:30:00)
  - [Upgrade Instruction](https://github.com/ChihuahuaChain/chihuahua/blob/main/mainnet/UPGRADES/v6)
- **v503** _(v5.0.4)_ - Block 9430000 - (2023-09-28 16:30:00) - (Halt height 9431130 due to wrong wasmd/wasmvm - Use v5.0.4 on directory v503)
  - [Upgrade Instruction](https://github.com/ChihuahuaChain/chihuahua/blob/main/mainnet/UPGRADES/v503)
- **v502** _(v5.0.2)_ - Block 9180000 - (2023-09-11 15:00:00)
  - [Upgrade Instruction](https://github.com/ChihuahuaChain/chihuahua/blob/main/mainnet/UPGRADES/v502)
- **v501** _(v5.0.1)_ - Block 8813000 - (2023-08-17 16:03:13)
  - [Upgrade Instruction](https://github.com/ChihuahuaChain/chihuahua/blob/main/mainnet/UPGRADES/v501)
- **v500** _(v5.0.0)_ - Block 8711111 - (2023-08-10 15:28:23)
  - [Upgrade Instruction](https://github.com/ChihuahuaChain/chihuahua/blob/main/mainnet/UPGRADES/v500)
- **v421** _(v4.2.1)_ - Block 6376376 - (2022-03-04 12:20:37)
  - [Upgrade Instruction](https://github.com/ChihuahuaChain/chihuahua/blob/main/mainnet/UPGRADES/v421)
- **v420** _(v4.2.0)_ - Block 6039999 - (2022-02-09 13:41:07)
  - [Upgrade Instruction](https://github.com/ChihuahuaChain/chihuahua/blob/main/mainnet/UPGRADES/v420)
- **v410** _(v4.1.0)_ - Block 4886666 - (2022-11-22 15:25:11)
  - [Upgrade Instruction](https://github.com/ChihuahuaChain/chihuahua/blob/main/mainnet/UPGRADES/v410)
- **v400** _(v4.0.0)_ - Block 4787878 - (2022-11-15 13:36:46)
  - [Upgrade Instruction](https://github.com/ChihuahuaChain/chihuahua/blob/main/mainnet/UPGRADES/v400)
- **v310** _(v3.1.0)_ - Block 4673333 - (2022-11-07 15:11:26)
  - [Upgrade Instruction](https://github.com/ChihuahuaChain/chihuahua/blob/main/mainnet/UPGRADES/v310)
- **iavl fast node** _(v2.4.*) - block unk -  state breaking upgrade that ensures that dragonberry is fully patched, go1.19 is used, and iavl fast node is properly configured.
- **burnmech** _(v2.3.0)_ - Block unk -  We won't be using this version, but we have it for historical purposes and for the creation of archive nodes.
- **burnmech** _(v2.2.2)_ - Block 4488444 - (2022-10-21 14:14:37) (retracted)
  - [Upgrade Instruction](https://github.com/ChihuahuaChain/chihuahua/blob/main/mainnet/UPGRADES/burnmech) 
- **authz** _(v2.1.0)_ - Block 4182410 - (2022-09-30 14:54:19)
  - [Upgrade Instruction](https://github.com/ChihuahuaChain/chihuahua/blob/main/mainnet/UPGRADES/authz)
- **minpropdeposit** _(v2.0.2)_ - Block 3654321 - (2022-08-25 13:00:26)
  - [Upgrade Instruction](https://github.com/ChihuahuaChain/chihuahua/blob/main/mainnet/UPGRADES/minpropdeposit)
- **Chiwawasm** _(v2.0.1)_ - Block 3000800 - (2022-07-11 17:02:14)
  - [Upgrade Instruction](https://github.com/ChihuahuaChain/chihuahua/blob/main/mainnet/UPGRADES/chiwawasm)
- **angryandy** _(v1.1.1)_ - Block 535000 - (2022-01-19 17:20:00)
  - [Upgrade Instruction](https://github.com/ChihuahuaChain/chihuahua/tree/main/mainnet/UPGRADES/angryandy)
