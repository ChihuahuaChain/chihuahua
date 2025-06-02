# Chihuahua v9.0.4 Upgrade

The Upgrade is scheduled for block `18385000`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/18385000)

This guide assumes that you use cosmovisor to manage upgrades.

## If you are syncing from 0 you need to apply v9.0.4 at height 18385000

```bash
# get the new version
cd chihuahua
git fetch --all
git checkout v9.0.4
make install
```

# check the version

```bash
# should be v9.0.4
chihuahuad version
# Should be commit d79ba427d0842401052272551bf465d4aa1a8fe0
chihuahuad version --long | grep commit
```

# Make new directory and copy binary

```bash
mkdir -p $HOME/.chihuahuad/cosmovisor/upgrades/v9.0.4/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahuad/cosmovisor/upgrades/v9.0.4/bin
```

# check the version again

```bash
# should be v9.0.4
$HOME/.chihuahuad/cosmovisor/upgrades/v9.0.4/bin/chihuahuad version
```
