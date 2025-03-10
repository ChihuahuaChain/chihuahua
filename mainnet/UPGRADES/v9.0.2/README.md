# Chihuahua v9.0.2 Upgrade

The Upgrade is scheduled for block `17073000`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/17073000)

This guide assumes that you use cosmovisor to manage upgrades.

## If you are syncing from 0 you need to apply v9.0.2 at height 17073000

```bash
# get the new version
cd chihuahua
git fetch --all
git checkout v9.0.2
make install
```

# check the version

```bash
# should be v9.0.2
chihuahuad version
# Should be commit ada240ea99bc0420182128dfa39f5ebc1f0a5e8b
chihuahuad version --long | grep commit
```

# Make new directory and copy binary

```bash
mkdir -p $HOME/.chihuahuad/cosmovisor/upgrades/v9.0.2/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahuad/cosmovisor/upgrades/v9.0.2/bin
```

# check the version again

```bash
# should be v9.0.2
$HOME/.chihuahuad/cosmovisor/upgrades/v9.0.2/bin/chihuahuad version
```
