# Chihuahua 2.0.2 Upgrade (minpropdeposit)

The Upgrade is scheduled for block `3654321`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/3654321)

This guide assumes that you use cosmovisor to manage upgrades

```bash
# get the new version
cd chihuahua
git pull
git checkout v2.0.2
make install
```

# check the version

```bash
# should be v2.0.2
chihuahuad version
# Should be commit eeb863cfbcf2731c20b1e9970f99a786507f0332
chihuahuad version --long
```

# Make new directory and copy binary

```bash
mkdir -p $HOME/.chihuahua/cosmovisor/upgrades/minpropdeposit/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahua/cosmovisor/upgrades/minpropdeposit/bin
```

# check the version again

```bash
# should be v2.0.2
$HOME/.chihuahua/cosmovisor/upgrades/minpropdeposit/bin/chihuahuad version
```
