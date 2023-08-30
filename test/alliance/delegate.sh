#!/bin/bash

echo ""
echo "#################################################"
echo "# Alliance: bridge funds and create an alliance #"
echo "#################################################"
echo ""

VALIDATOR="validator"
MAIN_HOME="$HOME/.chihuahuad"
MAIN_BINARY="chihuahuad --home=$MAIN_HOME"
MAIN_CHAIN_ID="chihuahua"
MAIN_DENOM="uhuahua"
MAIN_TX_FLAGS="--keyring-backend test --chain-id $MAIN_CHAIN_ID --from $VALIDATOR -y --fees=1000$MAIN_DENOM"

COUNTER_HOME="$HOME/.terrad"
COUNTER_BINARY="terrad --home=$COUNTER_HOME"
COUNTER_CHAIN_ID="terra"
COUNTER_DENOM="uluna"
COUNTER_TX_FLAGS="--keyring-backend test --chain-id $COUNTER_CHAIN_ID --from $VALIDATOR -y --fees=1000$COUNTER_DENOM --broadcast-mode block"

AMOUNT_TO_DELEGATE=10000000000
VAL_WALLET_1=$($COUNTER_BINARY keys show validator -a --keyring-backend test)
VAL_WALLET_2=$($MAIN_BINARY keys show validator -a --keyring-backend test)

echo "Sending tokens from validator wallet on counter chain to validator wallet on main"
IBC_TRANSFER=$($COUNTER_BINARY tx ibc-transfer transfer transfer channel-0 $VAL_WALLET_2 $AMOUNT_TO_DELEGATE$COUNTER_DENOM $COUNTER_TX_FLAGS -o json | jq -r '.raw_log' )

if [[ "$IBC_TRANSFER" == "failed to execute message"* ]]; then
    echo "Error: IBC transfer failed, with error: $IBC_TRANSFER"
    exit 1
fi

ACCOUNT_BALANCE=""
IBC_DENOM=""
while [ "$ACCOUNT_BALANCE" == "" ]; do
    IBC_DENOM=$($MAIN_BINARY q bank balances $VAL_WALLET_2 -o json | jq -r '.balances[0].denom')
    if [ "$IBC_DENOM" != "$MAIN_DENOM" ]; then
        ACCOUNT_BALANCE=$($MAIN_BINARY q bank balances $VAL_WALLET_2 -o json | jq -r '.balances[0].amount')
    fi
    sleep 2
done

echo "Creating an alliance with the denom $IBC_DENOM"
PROPOSAL_HEIGHT=$($MAIN_BINARY tx gov submit-legacy-proposal create-alliance $IBC_DENOM 5 0 5 0 0.99 1s --deposit 10000000000$MAIN_DENOM $MAIN_TX_FLAGS -o json | jq -r '.height')
sleep 5
PROPOSAL_ID=$($MAIN_BINARY query gov proposals --count-total -o json --output json | jq .proposals[-1].id -r)
VOTE_RES=$($MAIN_BINARY tx gov vote $PROPOSAL_ID yes $MAIN_TX_FLAGS -o json)

ALLIANCE="null"
while [ "$ALLIANCE" == "null" ]; do
    echo "Waiting for alliance with denom $IBC_DENOM to be created"
    ALLIANCE=$($MAIN_BINARY q alliance alliances -o json | jq -r '.alliances[0]')
    sleep 2
done

echo "Delegating $AMOUNT_TO_DELEGATE to the alliance $IBC_DENOM"
VAL_ADDR=$($MAIN_BINARY query staking validators --output json | jq .validators[0].operator_address --raw-output)
DELEGATE_RES=$($MAIN_BINARY tx alliance delegate $VAL_ADDR $AMOUNT_TO_DELEGATE$IBC_DENOM $MAIN_TX_FLAGS -o json)
sleep 5

DELEGATIONS=$($MAIN_BINARY query alliance delegation $VAL_WALLET_2 $VAL_ADDR $IBC_DENOM -o json | jq -r '.delegation.balance.amount')
if [[ "$DELEGATIONS" == "0" ]]; then
    echo "Error: Alliance delegations expected to be bigger than 0"
    exit 1
fi

echo "Query bank balance after alliance creation"
TOTAL_SUPPLY_BEFORE_ALLIANCE=$($MAIN_BINARY query bank total --denom $MAIN_DENOM --height $PROPOSAL_HEIGHT -o json | jq -r '.amount')
sleep 10
TOTAL_SUPPLY_AFTER_ALLIANCE=$($MAIN_BINARY query bank total --denom $MAIN_DENOM -o json | jq -r '.amount')
TOTAL_SUPPLY_INCREMENT=$(($TOTAL_SUPPLY_BEFORE_ALLIANCE - $TOTAL_SUPPLY_AFTER_ALLIANCE))

if [ "$TOTAL_SUPPLY_INCREMENT" -gt 100000 ] && [ "$TOTAL_SUPPLY_INCREMENT" -lt 1000000 ]; then
    echo "Error: Something went wrong, total supply of $MAIN_DENOM has increased out of range 100_000 between 1_000_000. current value $TOTAL_SUPPLY_INCREMENT"
    exit 1
fi

echo ""
echo "#########################################################"
echo "# Success: Alliance bridge funds and create an alliance #"
echo "#########################################################"
echo ""