# Chihuahua v9.0.5 Upgrade

The Upgrade is scheduled for block `18504000`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/18504000)

This guide assumes that you use cosmovisor to manage upgrades.

## If you are syncing from 0 you need to apply v9.0.5 at height 18504000

```bash
# get the new version
cd chihuahua
git fetch --all
git checkout v9.0.5
make install
```

# check the version

```bash
# should be v9.0.5
chihuahuad version
# Should be commit 286ef3b9b63837fd99c51d473b3f75ad824c18d5
chihuahuad version --long | grep commit
```

# Make new directory and copy binary

```bash
mkdir -p $HOME/.chihuahuad/cosmovisor/upgrades/v9.0.5/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahuad/cosmovisor/upgrades/v9.0.5/bin
```

# check the version again

```bash
# should be v9.0.5
$HOME/.chihuahuad/cosmovisor/upgrades/v9.0.5/bin/chihuahuad version
```
