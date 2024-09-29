# Chihuahua v8.0.0 Upgrade

The Upgrade is scheduled for block `14762000`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/14762000)

This guide assumes that you use cosmovisor to manage upgrades.

## If you are syncing from 0 you need to apply v8.0.0 at height 14762000

```bash
# get the new version
cd chihuahua
git fetch --all
git checkout v8.0.0
make install
```

# check the version

```bash
# should be v8.0.0
chihuahuad version
# Should be commit fc42649d573ee6c06225f52162cf325db3b3b7db
chihuahuad version --long | grep commit
```

# Make new directory and copy binary

```bash
mkdir -p $HOME/.chihuahuad/cosmovisor/upgrades/v8.0.0/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahuad/cosmovisor/upgrades/v8.0.0/bin
```

# check the version again

```bash
# should be v8.0.0
$HOME/.chihuahuad/cosmovisor/upgrades/v8.0.0/bin/chihuahuad version
```
