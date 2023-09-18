# Chihuahua v5.0.3 Upgrade

The Upgrade is scheduled for block `9430000`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/9430000)

This guide assumes that you use cosmovisor to manage upgrades.

```bash
# get the new version
cd chihuahua
git fetch --all
git checkout v5.0.3
make install
```

# check the version

```bash
# should be v5.0.3
chihuahuad version
# Should be commit 7511c95ad8b91a9230167b0c2a66f1db871cc2de
chihuahuad version --long | grep commit
```

# Make new directory and copy binary

```bash
mkdir -p $HOME/.chihuahuad/cosmovisor/upgrades/v503/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahuad/cosmovisor/upgrades/v503/bin
```

# check the version again

```bash
# should be v5.0.3
$HOME/.chihuahuad/cosmovisor/upgrades/v503/bin/chihuahuad version
```
