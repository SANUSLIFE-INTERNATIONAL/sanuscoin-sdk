// Copyright © 2021 The Sanuscoin Team

package app

import (
	"fmt"

	"github.com/goava/di"
	"github.com/urfave/cli/v2"

	"sanuscoin/sanuscoin-sdk/app/context"
	"sanuscoin/sanuscoin-sdk/config"
)

const (
	// These constants define the application semantic
	// and follow the semantic versioning 2.0.0 spec (http://semver.org/).
	appMajor uint = 0
	appMinor uint = 1
	appPatch uint = 0

	// appPreRelease flag to append suffix contain "-dev"
	// per the semantic versioning spec.
	appPreRelease = true

	// appVerPrefix use this prefix on version printing.
	appVerPrefix = "v"
)

// versionCommand appends version action to application.
func versionCommand(_ *di.Container, _ context.Context, _ *config.Config, app *App) {
	app.Commands = append(app.Commands, &cli.Command{
		Name:    "version",
		Usage:   "Show version",
		Aliases: []string{"ver"},
		Action: func(*cli.Context) error {
			ver := version()
			fmt.Println(ver)

			return nil
		},
	})
}

// version returns the application version as a properly formed string
// per the semantic versioning 2.0.0 spec (http://semver.org/).
func version() string {
	// start with the major, minor, and patch versions.
	ver := fmt.Sprintf("%s%d.%d.%d", appVerPrefix, appMajor, appMinor, appPatch)
	// append pre-release if there is one.
	if appPreRelease {
		ver += "-dev"
	}
	return ver
}
