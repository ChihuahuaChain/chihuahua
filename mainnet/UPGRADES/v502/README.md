# Chihuahua v5.0.2 Upgrade

The Upgrade is scheduled for block `9180000`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/9180000)

This guide assumes that you use cosmovisor to manage upgrades.

```bash
# get the new version
cd chihuahua
git fetch --all
git checkout v5.0.2
make install
```

# check the version

```bash
# should be v5.0.2
chihuahuad version
# Should be commit 3c20e93c7dc16f82fa2584382402f9149ec1b431
chihuahuad version --long
```

# Make new directory and copy binary

```bash
mkdir -p $HOME/.chihuahuad/cosmovisor/upgrades/v502/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahuad/cosmovisor/upgrades/v502/bin
```

# check the version again

```bash
# should be v5.0.2
$HOME/.chihuahuad/cosmovisor/upgrades/v502/bin/chihuahuad version
```
