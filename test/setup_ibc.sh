#!/bin/bash
set -eux
####################### Config variables & functions #######################
# Common
VALIDATOR="validator"
NODE_IP="localhost"

# Chi configs
CHI_CHAIN_ID="chihuahua"
CHI_MONIKER="chihuahua"
CHI_HOME="$HOME/.chihuahuad"
CHI_BINARY="chihuahuad --home=$CHI_HOME"
CHI_TX_FLAGS="--keyring-backend test --chain-id $CHI_CHAIN_ID --from $VALIDATOR -y --fees=1000uhuahua"
CHI_RPC_LADDR="$NODE_IP:26657"
CHI_P2P_LADDR="$NODE_IP:26656"
CHI_GRPC_ADDR="$NODE_IP:9090"

# Osmo configs
OSMO_CHAIN_ID="osmosis"
OSMO_MONIKER="osmosis"
OSMO_HOME="$HOME/.osmosisd"
OSMO_BINARY="osmosisd --home=$OSMO_HOME"
OSMO_TX_FLAGS="--keyring-backend test --chain-id $OSMO_CHAIN_ID --from $VALIDATOR -y --fees=1000uosmo"
OSMO_RPC_LADDR="$NODE_IP:26658"
OSMO_P2P_LADDR="$NODE_IP:26646"
OSMO_GRPC_ADDR="$NODE_IP:9091"


####################### Initializate chains #######################
echo "==============> Starting chain initialization...<=============="
# Clean start
killall $CHI_BINARY &> /dev/null || true
killall $OSMO_BINARY &> /dev/null || true
killall rly 2> /dev/null || true
rm -rf $CHI_HOME
rm -rf $OSMO_HOME
rm -rf ./test/relayer/keys
rm -rf ./test/logs
mkdir ./test/logs
cp ./test/relayer/config/config_temp.yaml ./test/relayer/config/config.yaml

# Chi chain init
$CHI_BINARY init --chain-id $CHI_CHAIN_ID $CHI_MONIKER
sed -i '' 's/"voting_period": "172800s"/"voting_period": "30s"/g' $CHI_HOME/config/genesis.json
sed -i '' 's/"max_deposit_period": "172800s"/"max_deposit_period": "30s"/g' $CHI_HOME/config/genesis.json
sed -i '' 's/stake/uhuahua/g' $CHI_HOME/config/genesis.json
sed -i -E "s|keyring-backend = \".*\"|keyring-backend = \"test\"|g" $CHI_HOME/config/client.toml
sed -i -E "s|minimum-gas-prices = \".*\"|minimum-gas-prices = \"0uhuahua\"|g" $CHI_HOME/config/app.toml

$CHI_BINARY keys add $VALIDATOR --keyring-backend=test
$CHI_BINARY genesis add-genesis-account $($CHI_BINARY keys show $VALIDATOR --keyring-backend=test -a) 1000000000000000000uhuahua
$CHI_BINARY genesis gentx validator 10000000000uhuahua --keyring-backend=test --chain-id=$CHI_CHAIN_ID
$CHI_BINARY genesis collect-gentxs 

# Osmo chain init
$OSMO_BINARY init --chain-id $OSMO_CHAIN_ID $OSMO_MONIKER
sed -i '' 's/"voting_period": "172800s"/"voting_period": "30s"/g' $OSMO_HOME/config/genesis.json
sed -i '' 's/"max_deposit_period": "172800s"/"max_deposit_period": "30s"/g' $OSMO_HOME/config/genesis.json
sed -i '' 's/stake/uosmo/g' $OSMO_HOME/config/genesis.json
sed -i -E "s|keyring-backend = \".*\"|keyring-backend = \"test\"|g" $OSMO_HOME/config/client.toml
sed -i -E "s|minimum-gas-prices = \".*\"|minimum-gas-prices = \"0uosmo\"|g" $OSMO_HOME/config/app.toml
sed -i -E "s|chain-id = \"\"|chain-id = \"${OSMO_CHAIN_ID}\"|g" $OSMO_HOME/config/client.toml
sed -i -E "s|node = \".*\"|node = \"tcp://${OSMO_RPC_LADDR}\"|g" $OSMO_HOME/config/client.toml

$OSMO_BINARY keys add $VALIDATOR --keyring-backend=test
$OSMO_BINARY add-genesis-account $($OSMO_BINARY keys show $VALIDATOR --keyring-backend=test -a) 1000000000000000000uosmo
$OSMO_BINARY gentx validator 10000000000uosmo --keyring-backend=test --chain-id=$OSMO_CHAIN_ID
$OSMO_BINARY collect-gentxs 


####################### Start chains #######################
echo "==============> Starting chihuahua...<=============="
$CHI_BINARY start \
       --rpc.laddr tcp://${CHI_RPC_LADDR} \
       --grpc.address ${CHI_GRPC_ADDR} \
       --p2p.laddr tcp://${CHI_P2P_LADDR} \
       --grpc-web.enable=false \
       --log_level trace \
       --trace \
       &> ./test/logs/chi &
( tail -f -n0 ./test/logs/chi & ) | grep -q "finalizing commit of block"
echo "Chain started"

echo "==============> Starting osmosis...<=============="
$OSMO_BINARY start \
       --rpc.laddr tcp://${OSMO_RPC_LADDR} \
       --grpc.address ${OSMO_GRPC_ADDR} \
       --p2p.laddr tcp://${OSMO_P2P_LADDR} \
       --grpc-web.enable=false \
       --log_level trace \
       --trace \
       &> ./test/logs/osmo &
( tail -f -n0 ./test/logs/osmo & ) | grep -q "finalizing commit of block"
echo "Chain started"

####################### Start relayer #######################

echo "==============> Funding relayers...<=============="
RELAYER_DIR="./test/relayer"
# chihuahua1hnuduzstgj2ze7l7g5rk5x4qlw8hd7es5s8xjd
MNEMONIC_1="vessel resist soda upset gadget spread sock egg soft barely hotel local weather image gaze core game once swarm nurse target fame stay small"
# osmo1sxcqadwaef3q2t40zr2qfaegemt4jndx4nm2d6
MNEMONIC_2="lyrics matter source business aisle naive ripple kidney honey brown carpet execute kite steak come system below erupt arch neither pond sort horse satisfy"

# send tokens to relayers
$CHI_BINARY tx bank send $VALIDATOR chihuahua1hnuduzstgj2ze7l7g5rk5x4qlw8hd7es5s8xjd 1000000uhuahua $CHI_TX_FLAGS
sleep 5
$OSMO_BINARY tx bank send $VALIDATOR osmo1sxcqadwaef3q2t40zr2qfaegemt4jndx4nm2d6 1000000uosmo $OSMO_TX_FLAGS
sleep 5



echo "==============> Restoring relayer accounts...<=============="
rly keys restore chihuahua rly1 "$MNEMONIC_1" --home $RELAYER_DIR
rly keys restore osmosis rly2 "$MNEMONIC_2" --home $RELAYER_DIR
rly transact link chi-osmo --home $RELAYER_DIR

echo "==============> Starting relayers...<=============="
sleep 5
rly start --home $RELAYER_DIR &> ./test/logs/rly