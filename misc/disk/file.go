// Copyright © 2021 The Sanuscoin Team

package disk

import (
	"os"
)

const (
	// DefaultFileMode controls the default permissions on any file.
	DefaultFileMode = os.FileMode(0o644)
)

// Create creates or truncates the named file.
func Create(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, DefaultFileMode)
}

func FileExists(fName string) bool {
	if _, err := os.Stat(fName); err == nil {
		return true
	}
	return false
}
