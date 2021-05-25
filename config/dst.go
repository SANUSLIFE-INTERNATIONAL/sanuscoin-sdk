// Copyright © 2021 The Sanuscoin Team

package config

const (
	dstDefaultPathName = "storage"
)

type (
	// dstConfig describes data store config.
	dstConfig struct {
		Path string
	}
)

// newDstConfig returns data store config instance.
func newDstConfig() *dstConfig {
	return &dstConfig{
		Path: dstDefaultPathName,
	}
}
