# Chihuahua v6 Upgrade

The Upgrade is scheduled for block `10666000`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/9430000)

This guide assumes that you use cosmovisor to manage upgrades.

## If you are syncing from 0 you need to apply v6 at height 10666000

```bash
# get the new version
cd chihuahua
git fetch --all
git checkout v6
make install
```

# check the version

```bash
# should be v6
chihuahuad version
# Should be commit c49c60854b6d167ea7c912ac4fd368d1546fa7dd
chihuahuad version --long | grep commit
```

# Make new directory and copy binary

```bash
mkdir -p $HOME/.chihuahuad/cosmovisor/upgrades/v6/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahuad/cosmovisor/upgrades/v6/bin
```

# check the version again

```bash
# should be v6
$HOME/.chihuahuad/cosmovisor/upgrades/v6/bin/chihuahuad version
```
