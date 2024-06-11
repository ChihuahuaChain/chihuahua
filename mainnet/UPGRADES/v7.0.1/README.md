# Chihuahua v7.0.1 Upgrade

The Upgrade is scheduled for block `13250000`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/13250000)

This guide assumes that you use cosmovisor to manage upgrades.

## If you are syncing from 0 you need to apply v7.0.1 at height 13250000

```bash
# get the new version
cd chihuahua
git fetch --all
git checkout v7.0.1
make install
```

# check the version

```bash
# should be v7.0.1
chihuahuad version
# Should be commit ed5ad95a6a79c84ec65c252a3b438566cb418fd7
chihuahuad version --long | grep commit
```

# Make new directory and copy binary

```bash
mkdir -p $HOME/.chihuahuad/cosmovisor/upgrades/v7.0.1/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahuad/cosmovisor/upgrades/v7.0.1/bin
```

# check the version again

```bash
# should be v7.0.1
$HOME/.chihuahuad/cosmovisor/upgrades/v7.0.1/bin/chihuahuad version
```
