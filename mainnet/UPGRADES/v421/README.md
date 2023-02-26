# Chihuahua v4.2.1 Upgrade

The Upgrade is scheduled for block `6376376`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/6376376)

This guide assumes that you use cosmovisor to manage upgrades

```bash
# get the new version
cd chihuahua
git fetch --all
git checkout v4.2.1
make install
```

# check the version

```bash
# should be v4.2.1
chihuahuad version
# Should be commit 3c9d81b261b6e53dabeb65043019baaaed325428
chihuahuad version --long |grep commit
```

# Make new directory and copy binary

```bash
mkdir -p $HOME/.chihuahuad/cosmovisor/upgrades/v421/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahuad/cosmovisor/upgrades/v421/bin
```

# check the version again

```bash
# should be v4.2.1
$HOME/.chihuahuad/cosmovisor/upgrades/v421/bin/chihuahuad version
```

