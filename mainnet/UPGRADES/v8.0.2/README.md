# Chihuahua v8.0.2 Upgrade

The Upgrade is scheduled for block `15103000`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/15103000)

This guide assumes that you use cosmovisor to manage upgrades.

## If you are syncing from 0 you need to apply v8.0.2 at height 15103000

```bash
# get the new version
cd chihuahua
git fetch --all
git checkout v8.0.2
make install
```

# check the version

```bash
# should be v8.0.2
chihuahuad version
# Should be commit e1fe631a7ae1d29387ca1269449010c398efc973
chihuahuad version --long | grep commit
```

# Make new directory and copy binary

```bash
mkdir -p $HOME/.chihuahuad/cosmovisor/upgrades/v8.0.2/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahuad/cosmovisor/upgrades/v8.0.2/bin
```

# check the version again

```bash
# should be v8.0.2
$HOME/.chihuahuad/cosmovisor/upgrades/v8.0.2/bin/chihuahuad version
```
