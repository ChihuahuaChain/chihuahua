# Chihuahua v7.0.1 Upgrade (USE v7.0.2 binary instead)

The Upgrade is scheduled for block `13250000`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/13250000)

This guide assumes that you use cosmovisor to manage upgrades.

## If you are syncing from 0 you need to apply v7.0.2 at height 13250000

```bash
# get the new version
cd chihuahua
git fetch --all
git checkout v7.0.2
make install
```

# check the version

```bash
# should be v7.0.2
chihuahuad version
# Should be commit 5b818dfb534994c1dd6542665f9c7504b25a185e
chihuahuad version --long | grep commit
```

# Make new directory and copy binary (yes, we are using v7.0.1 directory, not a typo)

```bash
mkdir -p $HOME/.chihuahuad/cosmovisor/upgrades/v7.0.1/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahuad/cosmovisor/upgrades/v7.0.1/bin
```

# check the version again

```bash
# should be v7.0.2
$HOME/.chihuahuad/cosmovisor/upgrades/v7.0.1/bin/chihuahuad version
```
