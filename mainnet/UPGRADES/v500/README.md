# Chihuahua v5.0.0 Upgrade

The Upgrade is scheduled for block `8711111`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/8711111)

This guide assumes that you use cosmovisor to manage upgrades

## Ensure you are using at least go1.20 before compiling (you will get an error anyways so, upgrade your go!) 

```bash
# get the new version
cd chihuahua
git fetch --all
git checkout v5.0.0
make install
```

# check the version

```bash
# should be v5.0.0
chihuahuad version
# Should be commit 1184218fd760897018f870eee78c14cef1ee0379
chihuahuad version --long
```

# Make new directory and copy binary

```bash
mkdir -p $HOME/.chihuahuad/cosmovisor/upgrades/v500/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahuad/cosmovisor/upgrades/v500/bin
```

# check the version again

```bash
# should be v5.0.0
$HOME/.chihuahuad/cosmovisor/upgrades/v500/bin/chihuahuad version
```
