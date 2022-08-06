## Installation Steps

### Install Prerequisites 

The following are necessary to build chihuahua from source. 

#### 1. Basic Packages
```bash:
# update the local package list and install any available upgrades 
sudo apt-get update && sudo apt upgrade -y 
# install toolchain and ensure accurate time synchronization 
sudo apt-get install make build-essential gcc git jq chrony -y
```

#### 2. Install Go
Follow the instructions [here](https://golang.org/doc/install) to install Go.

Alternatively, for Ubuntu LTS, you can do:
```bash:
wget -c https://go.dev/dl/go1.18.3.linux-amd64.tar.gz && sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.18.3.linux-amd64.tar.gz && sudo rm -rf go1.18.3.linux-amd64.tar.gz
```

Unless you want to configure in a non standard way, then set these in the `.profile` in the user's home (i.e. `~/`) folder.

```bash:
echo 'export GOROOT=/usr/local/go' >> $HOME/.bash_profile
echo 'export GOPATH=$HOME/go' >> $HOME/.bash_profile
echo 'export GO111MODULE=on' >> $HOME/.bash_profile
echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> $HOME/.bash_profile && . $HOME/.bash_profile
go version
```
Output should be: `go version go1.18.3 linux/amd64`

### 3. Install Chihuahua from source

```bash:
git clone https://github.com/ChihuahuaChain/chihuahua.git
cd chihuahua
git checkout main
make install
```
Note: there is no tag to build off of, just use main for now

### Init chain
```bash:
chihuahuad init $MONIKER_NAME --chain-id chihuahua-1
```

### Download Genesis
```bash:
cp mainnet/genesis.json ~/.chihuahua/config/genesis.json
```

### Add/recover keys
```bash:
# To create new keypair - make sure you save the mnemonics!
chihuahuad keys add <key-name> 

# Restore existing odin wallet with mnemonic seed phrase. 
# You will be prompted to enter mnemonic seed. 
chihuahuad keys add <key-name> --recover
```

## Instructions for post-genesis validators

### Create the validator

Note that proposal #1 agrees that all validators set commission to at
least 5% (this rule is now automatically enforced)

```bash:
chihuahuad tx staking create-validator \
  --from "<key-name>" \
  --amount "10000000uhuahua" \
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
  --gas-prices "0.025uhuahua"
```

### Backup critical files
```bash:
priv_validator_key.json
```
