# Chihuahua 2.2.1 Upgrade (burnmech)

The Upgrade is scheduled for block `4488444`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/4488444)

This guide assumes that you use cosmovisor to manage upgrades

```bash
# get the new version
cd chihuahua
git pull
git checkout v2.2.1
make install
```

# check the version

```bash
# should be v2.2.1
chihuahuad version
# Should be commit ffe1ed09b770acd9637eadfef941aff48f87b07b
chihuahuad version --long
```

# Make new directory and copy binary

```bash
mkdir -p $HOME/.chihuahua/cosmovisor/upgrades/burnmech/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahua/cosmovisor/upgrades/burnmech/bin
```

# check the version again

```bash
# should be v2.2.1
$HOME/.chihuahua/cosmovisor/upgrades/burnmech/bin/chihuahuad version
```
