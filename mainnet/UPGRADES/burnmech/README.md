# Chihuahua 2.2.2 Upgrade (burnmech)
# IMPORTANT
## It was previously announced as v2.2.1, the upgrade will use version v2.2.2 instead!

The Upgrade is scheduled for block `4488444`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/4488444)

This guide assumes that you use cosmovisor to manage upgrades

```bash
# get the new version
cd chihuahua
git pull
git checkout v2.2.2
make install
```

# check the version

```bash
# should be v2.2.2
chihuahuad version
# Should be commit 37d8d33091079d02ab6858c7f7b09e8ba4702fe0
chihuahuad version --long
```

# Make new directory and copy binary

```bash
mkdir -p $HOME/.chihuahua/cosmovisor/upgrades/burnmech/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahua/cosmovisor/upgrades/burnmech/bin
```

# check the version again

```bash
# should be v2.2.2
$HOME/.chihuahua/cosmovisor/upgrades/burnmech/bin/chihuahuad version
```
