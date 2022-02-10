#!/bin/sh

# --chain-id test --broadcast-mode block --keyring-backend test --from validator

chihuahuad tx gov submit-proposal software-upgrade angryandy --chain-id test --broadcast-mode block --keyring-backend test --from validator -y --deposit 10000000stake --title "upgrade" --description "1234" --upgrade-height 100

sleep 3

chihuahuad tx gov vote 1 VOTE_OPTION_YES --chain-id test --broadcast-mode block --keyring-backend test --from validator -y
