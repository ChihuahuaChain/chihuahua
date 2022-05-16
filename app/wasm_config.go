package app

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
)

const (
	// DefaultChihuahuaInstanceCost is initially set the same as in wasmd
	DefaultChihuahuaInstanceCost uint64 = 60_000
	// DefaultChihuahuaCompileCost set to a large number for testing
	DefaultChihuahuaCompileCost uint64 = 100
)

// ChihuahuaGasRegisterConfig is defaults plus a custom compile amount
func ChihuahuaGasRegisterConfig() wasmkeeper.WasmGasRegisterConfig {
	gasConfig := wasmkeeper.DefaultGasRegisterConfig()
	gasConfig.InstanceCost = DefaultChihuahuaInstanceCost
	gasConfig.CompileCost = DefaultChihuahuaCompileCost

	return gasConfig
}

func NewChihuahuaWasmGasRegister() wasmkeeper.WasmGasRegister {
	return wasmkeeper.NewWasmGasRegister(ChihuahuaGasRegisterConfig())
}
