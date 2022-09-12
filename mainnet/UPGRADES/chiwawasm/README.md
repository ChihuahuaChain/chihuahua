# Chihuahua 2.0.1 Upgrade

The Upgrade is scheduled for block `3000800`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/3000800)

This guide assumes that you use cosmovisor to manage upgrades

```bash
# get the new version
cd chihuahua
git pull
git checkout v2.0.1
make install
```

# check the version

```bash
# should be v2.0.1
chihuahuad version
# Should be commit 937a38f8d39b3e6fb3976132014f2888737a9790
chihuahuad version --long
```

# Make new directory and copy binary

```bash
mkdir -p $HOME/.chihuahua/cosmovisor/upgrades/chiwawasm/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahua/cosmovisor/upgrades/chiwawasm/bin
```

# check the version again

```bash
# should be v2.0.1
$HOME/.chihuahua/cosmovisor/upgrades/chiwawasm/bin/chihuahuad version
```
