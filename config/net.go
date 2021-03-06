// Copyright © 2021 The Sanuscoin Team

package config

import (
	"time"
)

const (
	defaultHTTPPort = ":8080"
	defaultRPCPort  = ":8099"

	netDefaultInterval = time.Second * 30
	netDefaultTimeout  = time.Second * 3

	netCfgIsInitialized  = "configuration already initialized"
	netCfgNotInitialized = "configuration not initialized properly"
)

type (
	// netConfig describes network config.
	netConfig struct {
		Testnet bool
		Http    string
		RPC     string
	}
)

// newNetConfig returns network config instance.
func newNetConfig() *netConfig {
	return &netConfig{
		Http: defaultHTTPPort,
		RPC: defaultRPCPort,
	}
}

// ScopeName determines current scope name of the network like:
// * mainnet - the daemon network, that contents real data.
// * testnet - is an alternative network, to be used for testing.
func (c *netConfig) ScopeName() string {
	if c.Testnet {
		return AppTestnetName
	}
	return AppMainNetName
}
