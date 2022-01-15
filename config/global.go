// Copyright © 2021 The Sanuscoin Team

package config

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"
)

const (
	AppName = "Sanuscoin"

	AppMainNetName = "mainnet"
	AppTestnetName = "testnet"

	appDataPathName   = "data"
	appDBPathName     = "ccdb"
	appLogsPathName   = "logs"
	appWalletPathName = "wallet"
	appRootPathName   = "sanuscoin"

	appConfigFilename = ".env"
)

var (
	appRootPath = getRootPath()

	appLogsPath = filepath.Join(appRootPath, appLogsPathName)
	appDataPath = filepath.Join(appRootPath, appDataPathName)
	appDBPath   = filepath.Join(appRootPath, appDataPathName)

	appConfigFile = filepath.Join(appRootPath, appConfigFilename)
)

// returns an operating system specific directory to be used for storing application data for an application.
// This unexported version takes an operating system argument primarily to enable the testing package to properly test
// the function by forcing an operating system that is not the currently one.
func getRootPath() string {
	// The caller really shouldn't prepend the name with a period, but if they do,
	// handle it gracefully by trimming it.
	appNameUpper := string(unicode.ToUpper(rune(AppName[0]))) + AppName[1:]
	appNameLower := string(unicode.ToLower(rune(AppName[0]))) + AppName[1:]

	// Get the OS specific home directory via the Go standard cc.
	var homeDir string
	usr, err := user.Current()
	if err == nil {
		homeDir = usr.HomeDir
	}

	// Fall back to standard HOME environment variable that works for most POSIX OSes if the directory from the
	// Go standard cc failed.
	if err != nil || homeDir == "" {
		homeDir = os.Getenv("HOME")
	}

	switch runtime.GOOS {
	// Attempt to use the LOCALAPPDATA or APPDATA environment variable on Windows.
	case "windows":
		// Windows XP and before didn't have a LOCALAPPDATA, so fallback
		// to transfer APPDATA when LOCALAPPDATA is not set.
		var appData string

		appData = os.Getenv("USERPROFILE")

		if appData == "" {
			appData = os.Getenv("LOCALAPPDATA")
		}

		if appData == "" {
			appData = os.Getenv("APPDATA")
		}

		if appData != "" {
			return filepath.Join(appData, "."+appNameLower)
		}

	case "darwin":
		if homeDir != "" {
			return filepath.Join(homeDir, "Library", "Application Support", appNameUpper)
		}

	case "plan9":
		if homeDir != "" {
			return filepath.Join(homeDir, appNameLower)
		}

	default:
		if homeDir != "" {
			return filepath.Join(homeDir, ".config", appNameLower)
		}
	}

	// Fall back to the current directory if all else fails.
	return "~/." + appNameLower
}

func InitPaths(cfg *Config) {
	subdir := cfg.Net.ScopeName()
	appRootPath = osAppRootPath()
	appLogsPath = filepath.Join(appRootPath, subdir, appLogsPathName)
	appDataPath = filepath.Join(appRootPath, subdir, appDataPathName)
}

// AppDataPath return path to application's data dir.
func AppDataPath() string {
	return appDataPath
}

// AppDBPath return path to application's kvdb dir.
func AppDBPath() string {
	return appDataPath
}

// AppLogsPath return path to application's logs dir.
func AppLogsPath() string {
	return appLogsPath
}

// AppLogPath return path to specific service log file
func AppLogPath(fName string) string {
	return filepath.Join(AppLogsPath(), fName)
}

// AppRootPath return path to application's root dir.
func AppRootPath() string {
	return appRootPath
}

// AppWalletPath return path to application's wallet dir.
func AppWalletPath() string {
	return filepath.Join(AppRootPath(), appWalletPathName)
}

func AppConfigFile() string {
	return appConfigFile
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
