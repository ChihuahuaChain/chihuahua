# Chihuahua v5.0.3 Upgrade (v5.0.4 is being used)

The Upgrade is scheduled for block `9430000`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/9430000)

This guide assumes that you use cosmovisor to manage upgrades.

## If you are syncing from 0 you need to apply v5.0.4 at height 9431130

```bash
# get the new version
cd chihuahua
git fetch --all
git checkout v5.0.4
make install
```

# check the version

```bash
# should be v5.0.4
chihuahuad version
# Should be commit ffc89d161eefe4c0333cd005e73f5db545852e5d
chihuahuad version --long | grep commit
```

# Make new directory and copy binary (yes, it's v5.0.3 in v503 directory)

```bash
mkdir -p $HOME/.chihuahuad/cosmovisor/upgrades/v503/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahuad/cosmovisor/upgrades/v503/bin
```

# check the version again

```bash
# should be v5.0.4
$HOME/.chihuahuad/cosmovisor/upgrades/v503/bin/chihuahuad version
```
