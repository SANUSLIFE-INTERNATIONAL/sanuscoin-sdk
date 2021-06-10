// Copyright © 2021 The Sanuscoin Team

package disk

import (
	"os"
)

const (
	// DefaultDirMode controls the default permissions
	// on any paths created by using MakeDirs.
	DefaultDirMode = os.FileMode(0o755)
)

// MakeDirs ensures that the full path you wanted exists.
func MakeDirs(dirs ...string) error {
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, DefaultDirMode); err != nil {
			return err
		}
	}
	return nil
}

// DirExists checks if specific directory exists
func DirExists(dirName string) bool {
	info, err := os.Stat(dirName)
	if err != nil {
		return false
	}
	return info.IsDir()
}
