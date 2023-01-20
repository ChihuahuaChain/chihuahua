package configurer

import (
	"os"
	"testing"
	"time"

	"github.com/ChihuahuaChain/chihuahua/tests/e2e/configurer/chain"
	"github.com/ChihuahuaChain/chihuahua/tests/e2e/containers"
	"github.com/ChihuahuaChain/chihuahua/tests/e2e/initialization"
)

type CurrentBranchConfigurer struct {
	baseConfigurer
}

var _ Configurer = (*CurrentBranchConfigurer)(nil)

func NewCurrentBranchConfigurer(t *testing.T, chainConfigs []*chain.Config, setupTests setupFn, containerManager *containers.Manager) Configurer {
	return &CurrentBranchConfigurer{
		baseConfigurer: baseConfigurer{
			chainConfigs:     chainConfigs,
			containerManager: containerManager,
			setupTests:       setupTests,
			syncUntilHeight:  defaultSyncUntilHeight,
			t:                t,
		},
	}
}

func (cb *CurrentBranchConfigurer) ConfigureChains() error {
	for _, chainConfig := range cb.chainConfigs {
		if err := cb.ConfigureChain(chainConfig); err != nil {
			return err
		}
	}
	return nil
}

func (cb *CurrentBranchConfigurer) ConfigureChain(chainConfig *chain.Config) error {
	cb.t.Logf("starting e2e infrastructure from current branch for chain-id: %s", chainConfig.ID)
	err := os.MkdirAll("chihuahua-e2e-testnet", 0o755)
	if err != nil {
		return err
	}
	cb.t.Logf("temp directory for chain-id %v: %v", chainConfig.ID, "/tmp/chihuahua-e2e-testnet")

	initializedChain, err := initialization.InitChain(chainConfig.ID, "/tmp/chihuahua-e2e-testnet", chainConfig.ValidatorInitConfigs, time.Duration(chainConfig.VotingPeriod*1000000000), time.Duration(chainConfig.ExpeditedVotingPeriod*1000000000), 0)
	if err != nil {
		return err
	}
	cb.initializeChainConfigFromInitChain(initializedChain, chainConfig)
	return nil
}

func (cb *CurrentBranchConfigurer) RunSetup() error {
	return cb.setupTests(cb)
}
