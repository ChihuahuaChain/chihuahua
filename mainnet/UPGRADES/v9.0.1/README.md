# Chihuahua v9.0.1 Upgrade

The Upgrade is scheduled for block `16623000`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/16623000)

This guide assumes that you use cosmovisor to manage upgrades.

## If you are syncing from 0 you need to apply v9.0.1 at height 16623000

```bash
# get the new version
cd chihuahua
git fetch --all
git checkout v9.0.1
make install
```

# check the version

```bash
# should be v9.0.1
chihuahuad version
# Should be commit 808edba9ec050437b13e877987d5386b15642254
chihuahuad version --long | grep commit
```

# Make new directory and copy binary

```bash
mkdir -p $HOME/.chihuahuad/cosmovisor/upgrades/v9.0.1/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahuad/cosmovisor/upgrades/v9.0.1/bin
```

# check the version again

```bash
# should be v9.0.1
$HOME/.chihuahuad/cosmovisor/upgrades/v9.0.1/bin/chihuahuad version
```
