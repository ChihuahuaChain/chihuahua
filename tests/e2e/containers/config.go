package containers

// ImageConfig contains all images and their respective tags
// needed for running e2e tests.
type ImageConfig struct {
	InitRepository string
	InitTag        string

	ChihuahuaRepository string
	ChihuahuaTag        string

	RelayerRepository string
	RelayerTag        string

	PriceFeederRepository string
	PriceFeederTag        string
}

//nolint:deadcode
const (
	// Current Git branch chihuahua repo/version. It is meant to be built locally.
	// It is used when skipping upgrade by setting HUA_E2E_SKIP_UPGRADE to true).
	// This image should be pre-built with `make docker-build-debug` either in CI or locally.
	CurrentBranchRepository = "chihuahua"
	CurrentBranchTag        = "debug"
	// Pre-upgrade chihuahua repo/tag to pull.
	// It should be uploaded to Docker Hub. HUA_E2E_SKIP_UPGRADE should be unset
	// for this functionality to be used.
	previousVersionRepository = "ghcr.io/cosmoscontracts/chihuahua"
	previousVersionTag        = "11.0.0-e2e"
	// Pre-upgrade repo/tag for chihuahua initialization (this should be one version below upgradeVersion)
	previousVersionInitRepository = "ghcr.io/cosmoscontracts/chihuahua"
	previousVersionInitTag        = "11.0.0-e2e-init-chain"
	// Hermes repo/version for relayer
	relayerRepository = "ghcr.io/cosmoscontracts/hermes"
	relayerTag        = "0.13.0"
	// Price-Feeder tool repo/tag for Oracle
	pricefeederRepository = "ghcr.io/cosmoscontracts/price-feeder"
	pricefeederTag        = "0.0.1"
)

// Returns ImageConfig needed for running e2e test.
// If isUpgrade is true, returns images for running the upgrade
// If isFork is true, utilizes provided fork height to initiate fork logic
func NewImageConfig(isUpgrade, isFork bool) ImageConfig {
	config := ImageConfig{
		RelayerRepository:     relayerRepository,
		RelayerTag:            relayerTag,
		PriceFeederRepository: pricefeederRepository,
		PriceFeederTag:        pricefeederTag,
	}

	if !isUpgrade {
		// If upgrade is not tested, we do not need InitRepository and InitTag
		// because we directly call the initialization logic without
		// the need for Docker.
		config.ChihuahuaRepository = CurrentBranchRepository
		config.ChihuahuaTag = CurrentBranchTag
		return config
	}

	// If upgrade is tested, we need to utilize InitRepository and InitTag
	// to initialize older state with Docker
	config.InitRepository = previousVersionInitRepository
	config.InitTag = previousVersionInitTag

	if isFork {
		// Forks are state compatible with earlier versions before fork height.
		// Normally, validators switch the binaries pre-fork height
		// Then, once the fork height is reached, the state breaking-logic
		// is run.
		config.ChihuahuaRepository = CurrentBranchRepository
		config.ChihuahuaTag = CurrentBranchTag
	} else {
		// Upgrades are run at the time when upgrade height is reached
		// and are submitted via a governance proposal. Thefore, we
		// must start running the previous Chihuahua version. Then, the node
		// should auto-upgrade, at which point we can restart the updated
		// Chihuahua validator container.
		config.ChihuahuaRepository = previousVersionRepository
		config.ChihuahuaTag = previousVersionTag
	}

	return config
}
