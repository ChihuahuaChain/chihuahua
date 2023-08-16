# Chihuahua v5.0.1 Upgrade

The Upgrade is scheduled for block `8813000`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/8813000)

This guide assumes that you use cosmovisor to manage upgrades

## Ensure you are using at least go1.20 before compiling (you will get an error anyways so, upgrade your go!) 

```bash
# get the new version
cd chihuahua
git fetch --all
git checkout v5.0.1
make install
```

# check the version

```bash
# should be v5.0.1
chihuahuad version
# Should be commit 12875e31481439c41125c8fae92a2d1e84eb0ce2
chihuahuad version --long
```

# Make new directory and copy binary

```bash
mkdir -p $HOME/.chihuahuad/cosmovisor/upgrades/v501/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahuad/cosmovisor/upgrades/v501/bin
```

# check the version again

```bash
# should be v5.0.1
$HOME/.chihuahuad/cosmovisor/upgrades/v501/bin/chihuahuad version
```
