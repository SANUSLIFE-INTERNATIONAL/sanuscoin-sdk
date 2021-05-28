// Copyright © 2021 The Sanuscoin Team

package app

import (
	"os"

	"github.com/goava/di"

	"sanus/sanus-sdk/app/context"
	"sanus/sanus-sdk/config"
)

// NewApp returns an application.
func NewApp(dic *di.Container, ctx context.Context, cfg *config.Config) *App {
	// create application
	app := newAppCli(cfg)

	// append cli commands
	app.command(dic, ctx, cfg, initCommand)
	app.command(dic, ctx, cfg, startCommand)
	app.command(dic, ctx, cfg, versionCommand)

	return app
}

// Start starts the application.
func Start(ctx context.Context, app *App) error {
	return app.RunContext(ctx, os.Args)
}
