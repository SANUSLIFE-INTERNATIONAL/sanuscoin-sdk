// Copyright © 2021 The Sanuscoin Team

package app

import (
	"log"

	"github.com/goava/di"
	"github.com/urfave/cli/v2"

	"sanus/sanus-sdk/app/context"
	"sanus/sanus-sdk/config"
)

type (
	// App describes cli application.
	App struct {
		*cli.App
	}

	// command describes func for append an action to app.
	command func(*di.Container, context.Context, *config.Config, *App)
)

// command appends a command to application actions list.
func (app *App) command(
	dic *di.Container,
	ctx context.Context,
	cfg *config.Config,
	add command,
) {
	add(dic, ctx, cfg, app)
}

// newAppCli returns application instance.
func newAppCli(cfg *config.Config) *App {
	return &App{
		App: &cli.App{
			Usage: "Sanuscoin regular node CLI",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "debug",
					Usage:       "Show debug info",
					Aliases:     []string{"D"},
					Destination: &cfg.App.Debug,
				},
				&cli.BoolFlag{
					Name:        "verbose",
					Usage:       "Show additional info",
					Aliases:     []string{"v"},
					Destination: &cfg.App.Verbose,
				},
			},
			ExitErrHandler: func(_ *cli.Context, err error) {
				if err != nil {
					log.Fatalln(err)
				}
			},
		},
	}
}
