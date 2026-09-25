# Chihuahua v9.1.0 Upgrade

The Upgrade is scheduled for block `TBD`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/TBD)

This guide assumes that you use cosmovisor to manage upgrades.

## What changes

v9.1.0 shuts down the Alliance module, following proposal #99: the upgrade closes the ampGASH alliance and clears the Alliance state. It is the first step of the Chihuahua v10 rebuild (Cosmos SDK v0.54, IBC v11, CosmWasm 3), which will remove the Alliance module entirely.

After the upgrade the Alliance module is empty and stays inactive until v10 removes it. Validators lose the voting power that came from the Alliance virtual stake.

### Fixes

- stakedrops in HUAHUA skipped the start/end block checks, allowing a start in the past or start == end
- a stakedrop starting at the current block lost its first block payout
- the `fee_payer` event attribute contained raw bytes instead of the bech32 address
- tokenfactory mint and burn events reported the sender instead of the target address
- `chihuahuad export` failed with `module ibchooks does not exist`

The liquidity module is now maintained in this repository (`x/liquidity`), with no change to its behaviour or state.

## If you are syncing from 0 you need to apply v9.1.0 at height TBD

```bash
# get the new version
cd chihuahua
git fetch --all
git checkout v9.1.0
make install
```

# check the version

```bash
# should be v9.1.0
chihuahuad version
# Should be commit TBD
chihuahuad version --long | grep commit
```

# Make new directory and copy binary

```bash
mkdir -p $HOME/.chihuahuad/cosmovisor/upgrades/v9.1.0/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahuad/cosmovisor/upgrades/v9.1.0/bin
```

# check the version again

```bash
# should be v9.1.0
$HOME/.chihuahuad/cosmovisor/upgrades/v9.1.0/bin/chihuahuad version
```
