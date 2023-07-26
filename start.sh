rm -rf $HOME/.chihuahuad
make install
chihuahuad init val1 --chain-id=testing --home=$HOME/.chihuahuad
chihuahuad keys add validator --keyring-backend=test --home=$HOME/.chihuahuad
chihuahuad genesis add-genesis-account $(chihuahuad keys show validator --keyring-backend=test --home=$HOME/.chihuahuad -a) 1000000000000000000uhuahua --home=$HOME/.chihuahuad
chihuahuad genesis gentx validator 10000000000uhuahua --keyring-backend=test --chain-id=testing
chihuahuad genesis collect-gentxs 
sed -i '' 's/"voting_period": "172800s"/"voting_period": "20s"/g' $HOME/.chihuahuad/config/genesis.json
sed -i '' 's/stake/uhuahua/g' $HOME/.chihuahuad/config/genesis.json
sed -i -E "s|minimum-gas-prices = \".*\"|minimum-gas-prices = \"0uhuahua\"|g" $HOME/.chihuahuad/config/app.toml

chihuahuad start --home=$HOME/.chihuahuad