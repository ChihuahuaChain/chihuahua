#!/bin/sh

set -o errexit -o nounset

CHAINID="test"

# Build genesis file incl account for passed address
coins="100000000000stake,100000000000samoleans"
chihuahuad init $CHAINID --chain-id $CHAINID

# modify the gov params to be super short
sed -i 's#"max_deposit_period": "172800s"#"max_deposit_period": "1s"#g' ~/.chihuahua/config/genesis.json
sed -i 's#"voting_period": "172800s"#"voting_period": "30s"#g' ~/.chihuahua/config/genesis.json

chihuahuad keys add validator --keyring-backend="test"
chihuahuad add-genesis-account $(chihuahuad keys show validator -a --keyring-backend="test") $coins
chihuahuad gentx validator 5000000000stake --keyring-backend="test" --chain-id $CHAINID --commission-rate "0.01"
chihuahuad collect-gentxs

# Set proper defaults and change ports
sed -i 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ~/.chihuahua/config/config.toml
sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ~/.chihuahua/config/config.toml
sed -i 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ~/.chihuahua/config/config.toml
sed -i 's/index_all_keys = false/index_all_keys = true/g' ~/.chihuahua/config/config.toml

# Start the chihuahua
chihuahuad start