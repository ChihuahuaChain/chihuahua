package simulation

// Simulation operation weights for liquidity messages.
const (
	DefaultWeightMsgCreatePool          int = 5
	DefaultWeightMsgDepositWithinBatch  int = 10
	DefaultWeightMsgWithdrawWithinBatch int = 10
	DefaultWeightMsgSwapWithinBatch     int = 85
)
