# Chihuahua v4.1.0 Upgrade (burn parameter)

The Upgrade is scheduled for block `4886666`. A countdown clock is [here](https://www.mintscan.io/chihuahua/blocks/4886666)

This guide assumes that you use cosmovisor to manage upgrades

```bash
# get the new version
cd chihuahua
git fetch --all
git checkout v4.1.0
make install
```

# check the version

```bash
# should be v4.1.0
chihuahuad version
# Should be commit 49a1b6d8f71bb0e981f6ff0fce5deae63e270324
chihuahuad version --long
```

# Make new directory and copy binary

```bash
mkdir -p $HOME/.chihuahuad/cosmovisor/upgrades/v410/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahuad/cosmovisor/upgrades/v410/bin
```

# check the version again

```bash
# should be v4.1.0
$HOME/.chihuahuad/cosmovisor/upgrades/v410/bin/chihuahuad version
```
