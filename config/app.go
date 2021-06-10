// Copyright © 2021 The Sanuscoin Team

package config

const (
	appDefaultName      = "SANUSCOIN"
	appDefaultDebugMode = "info"
)

type (
	// appConfig describes application config.
	appConfig struct {
		Debug string
		Name  string
	}
)

// newAppConfig returns app config instance.
func newAppConfig() *appConfig {
	return &appConfig{
		Name:  appDefaultName,
		Debug: appDefaultDebugMode,
	}
}
