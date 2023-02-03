# Chihuahua v4.2.0 Upgrade

The Upgrade is scheduled for block `6039999`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/6039999)

This guide assumes that you use cosmovisor to manage upgrades

```bash
# get the new version
cd chihuahua
git fetch --all
git checkout v4.2.0
make install
```

# check the version

```bash
# should be v4.2.0
chihuahuad version
# Should be commit b3909d0eb5dc2ef0b81261a89b47bcd58c558154
chihuahuad version --long
```

# Make new directory and copy binary

```bash
mkdir -p $HOME/.chihuahuad/cosmovisor/upgrades/v420/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahuad/cosmovisor/upgrades/v420/bin
```

# check the version again

```bash
# should be v4.2.0
$HOME/.chihuahuad/cosmovisor/upgrades/v420/bin/chihuahuad version
```
