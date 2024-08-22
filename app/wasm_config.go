package app

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

const (
	// DefaultChihuahuaInstanceCost is initially set the same as in wasmd
	DefaultChihuahuaInstanceCost uint64 = 60_000
	// DefaultChihuahuaCompileCost set to a large number for testing
	DefaultChihuahuaCompileCost uint64 = 100
)

// ChihuahuaGasRegisterConfig is defaults plus a custom compile amount
func ChihuahuaGasRegisterConfig() wasmtypes.WasmGasRegisterConfig {
	gasConfig := wasmtypes.DefaultGasRegisterConfig()
	gasConfig.InstanceCost = DefaultChihuahuaInstanceCost
	gasConfig.CompileCost = DefaultChihuahuaCompileCost

	return gasConfig
}

func NewChihuahuaWasmGasRegister() wasmtypes.WasmGasRegister {
	return wasmtypes.NewWasmGasRegister(ChihuahuaGasRegisterConfig())
}
