// Copyright © 2021 The Sanuscoin Team

package app

import (
	"sanus/sanus-sdk/misc/log"

	"github.com/goava/di"
	"github.com/urfave/cli/v2"

	"sanus/sanus-sdk/app/context"
	"sanus/sanus-sdk/config"
)

type (
	// App describes cli application.
	App struct {
		*cli.App

		*log.Logger
	}

	// command describes func for append an action to app.
	command func(*di.Container, context.Context, *config.Config, *App)
)

// command appends a command to application actions list.
func (application *App) command(
	dic *di.Container,
	ctx context.Context,
	cfg *config.Config,
	add command,
) {
	add(dic, ctx, cfg, application)
}

func (application *App) initLogger() {
	application.SetOutput(defaultLogFile, "APP")
}

// newAppCli returns application instance.
func newAppCli(cfg *config.Config) *App {
	logger := log.NewLogger(cfg)
	return &App{
		Logger: logger,
		App: &cli.App{
			Usage: "Sanuscoin transfer node CLI",
		},
	}
}
