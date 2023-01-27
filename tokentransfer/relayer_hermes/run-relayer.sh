#!/bin/sh

hermes keys add --chain chihuahuaa --mnemonic-file "alice.json"
hermes keys add --chain chihuahuab --mnemonic-file "bob.json"

hermes create channel --a-chain chihuahuaa --b-chain chihuahuab --a-port transfer --b-port transfer --new-client-connection
hermes start