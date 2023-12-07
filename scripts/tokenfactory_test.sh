
BASE_FLAGS="--chain-id=local-1 --keyring-backend=test --home=$HOME/.huahua --yes"

# Create denom and verify
chihuahuad tx tokenfactory create-denom woof --from hua1 $BASE_FLAGS --gas=2100000
addr=$(chihuahuad keys show hua1 -a --keyring-backend=test --home=$HOME/.huahua); echo $addr
chihuahuad q tokenfactory denoms-from-creator $addr

# mint denom to self
chihuahuad tx tokenfactory mint 1000factory/chihuahua1hj5fveer5cjtn4wd6wstzugjfdxzl0xp9eludp/woof --from hua1 $BASE_FLAGS
chihuahuad q bank balances $addr
