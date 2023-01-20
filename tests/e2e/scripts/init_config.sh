#!/bin/sh
sed -i '114,117d' config.toml
address=`sed '230!d' /tmp/chihuahua-e2e-testnet/chihuahua-test-a/chihuahua-test-a-node-prune-default-snapshot-state-sync-from/config/genesis.json`
validator=`sed '231!d' /tmp/chihuahua-e2e-testnet/chihuahua-test-a/chihuahua-test-a-node-prune-default-snapshot-state-sync-from/config/genesis.json`
var1=${address#*chihuahua}
var2=${var1%\"*}
var3=${validator#*chihuahua}
var4=${var3%\"*}
echo "address = \"chihuahua$var2\"" >> config.toml
echo "chain_id = \"chihuahua-test-a\"" >> config.toml
echo "validator = \"chihuahua$var4\"" >> config.toml
echo "prefix = \"chihuahua\"" >> config.toml