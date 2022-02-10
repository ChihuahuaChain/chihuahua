#!/bin/bash
CHAIN="chihuahua"
GIT_LINK="https://raw.githubusercontent.com/cosmos/chain-registry/master"
ALTLINK="https://raw.githubusercontent.com/clemensgg/RELAYER-dev-crew/main/chains"
GOLINK="https://git.io/vQhTU"
DATABASE="goleveldb"
FIX_TRUST_PERIOD=448h0m0s
FIX_HEIGHT_DIFF=5000
FIX_VERSION="v1.1.1"
echo "---------------------- S T A T E - S Y N C ----------------------"


    #main
main() {
    install_dependencies
    install_go
    fetch_cr
    check_rpc
    build_init
    config
    start
}

    #helperfunction check unique vaulues
unique_values() {
    typeset i
    for i do
        [ "$1" = "$i" ] || return 1
    done
    return 0
}

    #install basic dependencies
install_dependencies() {
    echo "> updating dependencies..."
    sudo apt update -qq && sudo apt upgrade -qq
    sudo apt install -qq build-essential git curl jq wget
}

    #install go
install_go() {
    echo "> installing go..."
    sudo rm -rf /usr/local/go && sudo rm -rf $HOME/.go
    wget -q -O - $GOLINK | bash
}

    #fetch chain-registry
fetch_cr() {
    echo "> fetching chain-registry..."
    echo "-----------------------------------------------------------------"

    CHAIN_JSON=$(curl -s ${GIT_LINK}/$CHAIN/chain.json)
    NODE_HOME_DIR=$(echo $CHAIN_JSON | jq -r '.node_home') 
    NODE_HOME_DIR=$(eval echo $NODE_HOME_DIR)
    CHAIN_NAME=$(echo $CHAIN_JSON | jq -r '.chain_name')
    CHAIN_ID=$(echo $CHAIN_JSON | jq -r '.chain_id')
    NODED=$(echo $CHAIN_JSON | jq -r '.daemon_name')
    GEN_URL=$(echo $CHAIN_JSON | jq -r '.genesis.genesis_url')
    DPATH=$(echo $CHAIN_JSON | jq -r '.slip44')
    GIT_REPO=$(echo $CHAIN_JSON | jq -r '.codebase.git_repo')
    VERSION=$(echo $CHAIN_JSON | jq -r '.codebase.recommended_version')
    SEEDS=$(echo $CHAIN_JSON | jq -r '.peers.seeds')
    PEERS=$(echo $CHAIN_JSON | jq -r '.peers.persistent_peers')
    RPC_SERVERS=$(echo $CHAIN_JSON | jq -r '.apis.rpc')
    MEP2P=$(curl -s ifconfig.me):26656
    SEEDLIST=""
    PEERLIST=""
    RPCLIST=""

    readarray -t arr < <(jq -c '.[]' <<< $SEEDS)
    for item in ${arr[@]}; do
        ID=$(echo $item | jq -r '.id')
        ADD=$(echo $item | jq -r '.address')
        SEEDLIST="${SEEDLIST},${ID}@${ADD}"
    done
    for item in ${arr[@]}; do
        ID=$(echo $item | jq -r '.id')
        ADD=$(echo $item | jq -r '.address')
        PEERLIST="${PEERLIST},${ID}@${ADD}"
    done
    readarray -t arr < <(jq -c '.[]' <<< $RPC_SERVERS)
    for item in ${arr[@]}; do
        ADD=$(echo $item | jq -r '.address')
        RPCLIST="${RPCLIST},${ADD}"
    done
    SEEDLIST="${SEEDLIST:1}"
    PEERLIST="${PEERLIST:1},8cd0e06a36e618f22cfc827b184b0011d1bb8164@157.90.81.24:22110"
    RPCLIST="${RPCLIST:1},http://157.90.81.24:22111"

        #fix version, repo, homedir...
    if [ ! -z "$FIX_VERSION" ] ; then
        VERSION=$FIX_VERSION
    fi  
    if [ ! -z "$FIX_TRUST_PERIOD" ] ; then
        TRUST_PERIOD=$FIX_TRUST_PERIOD
    fi 
    if [ ! -z "$FIX_HEIGHT_DIFF" ] ; then
        HEIGHT_DIFF=$FIX_HEIGHT_DIFF
    fi 

    echo "home dir: $NODE_HOME_DIR"
    echo "chain name: $CHAIN_NAME"
    echo "chain id: $CHAIN_ID"
    echo "daemon name: $NODED"
    echo "genesis file url: $GEN_URL"
    echo "git repo: $GIT_REPO"
    echo "version: $VERSION"
    echo "seeds: $SEEDLIST"
    echo "rpc servers: $RPCLIST"
    echo "height diff: $HEIGHT_DIFF"
    echo "trust period: $TRUST_PERIOD"
    echo "-----------------------------------------------------------------"
}

    #check rpc connectivity, query trust hash
check_rpc(){
    HASHES=""
    echo "> checking RPC connectivity..."
    IFS=',' read -ra rpcarr <<< "$RPCLIST"
    for rpc in ${rpcarr[@]}; do
    re='.*[0-9].*'
        if [[ "$rpc" == *"https://"* ]] && [[ ! $rpc =~ $re ]] ; then
            if [ ${rpc: -1} = "/" ] ; then
                rpc=${rpc%?}
            fi
            rpc=$rpc:443
        fi
        if [[ "$rpc" == *"http://"* ]] && [[ ! $rpc =~ $re ]] ; then
            if [ ${rpc: -1} = "/" ] ; then
                rpc=${rpc%?}
            fi
            rpc=$rpc:80
        fi
        if [[ ! "$rpc" == *"http://"* ]] && [[ ! $rpc =~ $re ]] ; then
            if [ ${rpc: -1} = "/" ] ; then
                rpc=${rpc%?}
            fi
            rpc=$rpc:26657
        fi
        RPCNUM=$((RPCNUM+1))
        RES=$(curl -s $rpc/status --connect-timeout 3) || true
        if [ -z "$RES" ] || [[ "$RES" == *"Forbidden"* ]]; then
            echo "> $rpc didn't respond. dropping..."
        else
            HEIGHT=$(echo $RES | jq -r '.result.sync_info.latest_block_height')
            CHECKHEIGHT=$(($HEIGHT-$HEIGHT_DIFF))
            RES=$(curl -s "$rpc/commit?height=$CHECKHEIGHT")
            HASH=$(echo $RES | jq -r '.result.signed_header.commit.block_id.hash')
            TRUSTHASH=$HASH
            HASHES="${HASHES},${HASH}"
            RPCLIST_FINAL="${RPCLIST_FINAL},${rpc}"
            echo $rpc last height: $HEIGHT
        fi
    done
    
    if [ -z "$TRUSTHASH" ] ; then
        echo "> trust hash empty. couldn't connect to any RPCs. exiting..."
        exit
    fi
    HASHES="${HASHES:1}"

    RPCLIST_FINAL="${RPCLIST_FINAL:1}"
    if [[ ! "$RPCLIST_FINAL" == *","* ]] ; then
        RPCLIST_FINAL="${RPCLIST_FINAL},${RPCLIST_FINAL}"
    fi

    echo "working rpc list: $RPCLIST_FINAL"
    if unique_values "${HASHES[@]}"; then
        echo "> hash checks passed!"
        echo "> trust hash: $TRUSTHASH"
        echo "> trust height: $CHECKHEIGHT"
    else
        echo "> hash checks failed, exiting..."
        exit
    fi
    echo "-----------------------------------------------------------------"
}

    #build and initialize node
build_init(){
    echo "> building $NODED $VERSION from $GIT_REPO..."
    if [ -d "$HOME/$CHAIN_NAME-core" ] ; then
        cd $HOME/$CHAIN_NAME-core && git fetch
    else
        mkdir -p $HOME/$CHAIN_NAME-core
        git clone $GIT_REPO $HOME/$CHAIN_NAME-core && cd $HOME/$CHAIN_NAME-core
    fi

    git checkout $VERSION 
    if [[ "$DATABASE" == *"rocksdb"* ]] ; then
        BUILD_TAGS=rocksdb make install
    else
        make install
    fi
    cd

    RAND=$(echo $RANDOM | md5sum | head -c 6; echo;)
    echo "> initializing $NODED with moniker $RAND"
    echo "> home dir: ${NODE_HOME_DIR}"
    
    $NODED init $RAND --chain-id=$CHAIN_ID  --home $NODE_HOME_DIR -o

    echo "> downloading genesis from $GEN_URL..."
    rm ${NODE_HOME_DIR}/config/genesis.json
    if [[ "$GEN_URL" == *".gz"* ]]; then
        wget -q $GEN_URL -O genesis.json.gz
        gunzip -df genesis.json.gz && mv genesis.json ${NODE_HOME_DIR}/config/genesis.json
    elif [[ "$GEN_URL" == *".bz2"* ]]; then
        wget -q $GEN_URL -O genesis.tar.bz2
        tar -xjf genesis.tar.bz2 && mv genesis.json ${NODE_HOME_DIR}/config/genesis.json
    else
        wget -q $GEN_URL -O ${NODE_HOME_DIR}/config/genesis.json
    fi
}

    #configure state-sync
config() {
    echo "> configuring seeds & state-sync"
    sed -i '/rpc_servers = ""/c rpc_servers = "'$RPCLIST_FINAL'"' $NODE_HOME_DIR/config/config.toml
    sed -i 's/external_address = ""/external_address = "'$MEP2P'"/g' $NODE_HOME_DIR/config/config.toml
    sed -i 's/seeds = .*/seeds = "'$SEEDLIST'"/g' $NODE_HOME_DIR/config/config.toml
    sed -i 's/enable = false/enable = true/g' $NODE_HOME_DIR/config/config.toml
    sed -i 's/trust_height.*/trust_height = '$CHECKHEIGHT'/g' $NODE_HOME_DIR/config/config.toml
    sed -i 's/trust_hash.*/trust_hash = "'$TRUSTHASH'"/g' $NODE_HOME_DIR/config/config.toml
    sed -i 's/trust_period.*/trust_period = "'$TRUST_PERIOD'"/g' $NODE_HOME_DIR/config/config.toml
#    sed -i 's/persistent_peers = .*/persistent_peers = "'$PEERLIST'"/g' $NODE_HOME_DIR/config/config.toml
}

    #spin up node
start() {
    echo "> starting $NODED..."
    $NODED unsafe-reset-all --home $NODE_HOME_DIR
    $NODED start --home $NODE_HOME_DIR --x-crisis-skip-assert-invariants --db_backend $DATABASE || true
    echo "done!"
}


main; exit


