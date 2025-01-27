# Chihuahua v9.0.0 Upgrade

The Upgrade is scheduled for block `16529000`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/16529000)

This guide assumes that you use cosmovisor to manage upgrades.

## If you are syncing from 0 you need to apply v9.0.0 at height 16529000

```bash
# get the new version
cd chihuahua
git fetch --all
git checkout v9.0.0
make install
```

# check the version

```bash
# should be v9.0.0
chihuahuad version
# Should be commit 17c7358f09dd399ce0106762deef79eafa246217
chihuahuad version --long | grep commit
```

# Make new directory and copy binary

```bash
mkdir -p $HOME/.chihuahuad/cosmovisor/upgrades/v9.0.0/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahuad/cosmovisor/upgrades/v9.0.0/bin
```

# check the version again

```bash
# should be v9.0.0
$HOME/.chihuahuad/cosmovisor/upgrades/v9.0.0/bin/chihuahuad version
```
