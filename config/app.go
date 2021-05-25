// Copyright © 2021 The Sanuscoin Team

package config

const (
	appDefaultName = "SANUSCOIN"
)

type (
	// appConfig describes application config.
	appConfig struct {
		Debug   bool
		Name    string
		Verbose bool
	}
)

// newAppConfig returns app config instance.
func newAppConfig() *appConfig {
	return &appConfig{
		Name: appDefaultName,
	}
}
