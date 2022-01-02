# chihuahua
**chihuahua** is a blockchain bilt using Cosmos SDK and Tendermint and created with [Starport](https://github.com/tendermint/starport).

## Get started

```
starport chain serve
```

`serve` command installs dependencies, builds, initializes, and starts your blockchain in development.

### Configure

Your blockchain in development can be configured with `config.yml`. To learn more, see the [Starport docs](https://docs.starport.network).

### Launch

To launch your blockchain live on multiple nodes, use `starport network` commands. Learn more about [Starport Network](https://github.com/tendermint/spn).

### Web Frontend

Starport has scaffolded a Vue.js-based web app in the `vue` directory. Run the following commands to install dependencies and start the app:

```
cd vue
npm install
npm run serve
```

The frontend app is built using the `@starport/vue` and `@starport/vuex` packages. For details, see the [monorepo for Starport front-end development](https://github.com/tendermint/vue).

## Release
To release a new version of your blockchain, create and push a new tag with `v` prefix. A new draft release with the configured targets will be created.

```
git tag v0.1
git push origin v0.1
```

After a draft release is created, make your final changes from the release page and publish it.

### Install
To install the latest version of your blockchain node's binary, execute the following command on your machine:
# MAINNET: CHIHUAHUA-1

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
wget https://golang.org/dl/go1.17.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.17.5.linux-amd64.tar.gz
```

Unless you want to configure in a non standard way, then set these in the `.profile` in the user's home (i.e. `~/`) folder.

```bash:
cat <<EOF >> ~/.profile
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export GO111MODULE=on
export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
EOF
source ~/.profile
go version
```
Output should be: `go version go1.17.5 linux/amd64`

### 3. Install Chihuahua from source

```bash:
git clone https://github.com/ChihuahuaChain/chihuahua.git
cd chihuahua
make install
```
Note: there is no tag to build off of, just use master for now

### Init chain
```bash:
chihuahuad init $MONIKER_NAME --chain-id chihuahua-1
```

### Download Genesis
```bash:
wget -O ~/.chihuahua/config/genesis.json https://raw.githubusercontent.com/ChihuahuaChain/mainnet/main/genesis.json
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
least 5%!

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

## Instructions for Genesis Validators

### GenTx Creation

### Create Gentx

#### Add genesis account:
```
chihuahuad add-genesis-account <key-name> 5000001000000uhuahua
```
Note: if you receive message: `failed to get address from Keybase:`, add `--keyring-backend os`

#### Create Gentx
```
chihuahuad gentx <key-name> 5000000000000uhuahua \
--chain-id chihuahua-1 \
--moniker="<moniker>" \
--commission-max-change-rate=0.01 \
--commission-max-rate=0.20 \
--commission-rate=0.05 \
--details="XXXXXXXX" \
--security-contact="XXXXXXXX" \
--website="XXXXXXXX"
```

### Submit PR with Gentx and peer id
1. Copy the contents of ${HOME}/.chihuahua/config/gentx/gentx-XXXXXXXX.json.
2. Fork the repository
3. Create a file gentx-{{VALIDATOR_NAME}}.json under the /gentxs folder in the forked repo, paste the copied text into the file.
4. Create a Pull Request to the main branch of the repository


### Backup critical files
```bash:
priv_validator_key.json
```

```
curl https://get.starport.network/ChihuahuaChain/chihuahua@latest! | sudo bash
```
`ChihuahuaChain/chihuahua` should match the `username` and `repo_name` of the Github repository to which the source code was pushed. Learn more about [the install process](https://github.com/allinbits/starport-installer).

## Learn more

- [Starport](https://github.com/tendermint/starport)
- [Starport Docs](https://docs.starport.network)
- [Cosmos SDK documentation](https://docs.cosmos.network)
- [Cosmos SDK Tutorials](https://tutorials.cosmos.network)
- [Discord](https://discord.gg/cosmosnetwork)



