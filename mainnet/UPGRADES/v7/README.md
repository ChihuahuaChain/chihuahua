# Chihuahua v7 Upgrade

The Upgrade is scheduled for block `12900000`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/12900000)

This guide assumes that you use cosmovisor to manage upgrades.

## If you are syncing from 0 you need to apply v7 at height 12900000

```bash
# get the new version
cd chihuahua
git fetch --all
git checkout v7
make install
```

# check the version

```bash
# should be v7
chihuahuad version
# Should be commit 2022cdb82b605a12c8ada3acde5b32b9e9306a8c
chihuahuad version --long | grep commit
```

# Make new directory and copy binary

```bash
mkdir -p $HOME/.chihuahuad/cosmovisor/upgrades/v7/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahuad/cosmovisor/upgrades/v7/bin
```

# check the version again

```bash
# should be v7
$HOME/.chihuahuad/cosmovisor/upgrades/v7/bin/chihuahuad version
```
