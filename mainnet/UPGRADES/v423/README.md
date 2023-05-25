# Chihuahua v4.2.3 Upgrade - Huckleberry Patch

This NON CONSENSUS BREAKING upgrade can be applied at your own convenience, we recommend doing that at the earliest opportunity.

 guide assumes that you use cosmovisor to manage upgrades

```bash
# get the new version
cd chihuahua
git fetch --all
git checkout v4.2.3
make install
```

# check the version

```bash
# should be v4.2.3
chihuahuad version
# Should be commit 21989a2b212e904dfe3f7bd4bc679b8c493692d8
chihuahuad version --long |grep commit
```

# Make new directory and copy binary

```bash
mkdir -p $HOME/.chihuahuad/cosmovisor/upgrades/v421/bin
cp $HOME/go/bin/chihuahuad $HOME/.chihuahuad/cosmovisor/upgrades/v421/bin
```

# check the version again

```bash
# should be v4.2.3
$HOME/.chihuahuad/cosmovisor/upgrades/v421/bin/chihuahuad version
```

