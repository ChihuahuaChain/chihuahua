# Chihuahua 2.1.0 Upgrade (authz)

The Upgrade is scheduled for block `4182410`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/4182410)

This guide assumes that you use cosmovisor to manage upgrades

```bash
# get the new version
cd chihuahua
git pull
git checkout v2.1.0
make install
```

# check the version

```bash
# should be v2.1.0
chihuahuad version
# Should be commit 35b9a17713b87a4691a9430c28eef37bf3bbb4c1
chihuahuad version --long
```

# Make new directory and copy binary

```bash
mkdir -p $HOME/.chihuahua/cosmovisor/upgrades/authz/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahua/cosmovisor/upgrades/authz/bin
```

# check the version again

```bash
# should be v2.1.0
$HOME/.chihuahua/cosmovisor/upgrades/authz/bin/chihuahuad version
```
