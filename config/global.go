// Copyright © 2021 The Sanuscoin Team

package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	AppMainnetName = "mainnet"
	AppTestnetName = "testnet"

	appDataPathName = "data"
	appLogsPathName = "logs"
	appRootPathName = "sanuscoin"
)

var (
	appRootPath string
	appLogsPath string
	appDataPath string
)

// AppDataPath return path to application's data dir.
func AppDataPath() string {
	return appDataPath
}

// AppLogsPath return path to application's logs dir.
func AppLogsPath() string {
	return appLogsPath
}

// AppDataPath return path to application's root dir.
func AppRootPath() string {
	return appRootPath
}

// osAppRootPath wraps func os.UserConfigDir.
func osAppRootPath() string {
	path := appRootPathName
	dir, _ := os.UserConfigDir()
	switch runtime.GOOS {
	case "windows", "darwin":
		path = strings.ToUpper(path[:1]) + path[1:]
	}

	return filepath.Join(dir, path)
}
