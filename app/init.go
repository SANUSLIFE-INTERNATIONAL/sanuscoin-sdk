// Copyright © 2021 The Sanuscoin Team

package app

import (
	"fmt"

	"github.com/goava/di"
	"github.com/urfave/cli/v2"

	"sanus/sanus-sdk/app/context"
	"sanus/sanus-sdk/config"
)

// initCommand appends initialize action to cli app.
func (application *App) initCommand(dic *di.Container, _ context.Context, cfg *config.Config, app *App) {
	app.Commands = append(app.Commands, &cli.Command{
		Name:  "init",
		Usage: "Init config",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        config.AppTestnetName,
				Usage:       fmt.Sprintf("Use %v (mainnet otherwise)", config.AppTestnetName),
				Destination: &cfg.Net.Testnet,
			},
		},
		Before: func(*cli.Context) error {
			// invoke config maker
			if err := dic.Invoke(config.Make); err != nil {
				return fmt.Errorf("make config: %w", err)
			}
			application.SetOutput(defaultLogFile, "APP")
			application.Infof("Creating configuration directory %v  \n", config.AppRootPath())
			return nil
		},
		Action: func(*cli.Context) error {
			application.Info("Initialize configuration...")
			// invoke config initializer
			if err := dic.Invoke(config.Init); err != nil {
				return fmt.Errorf("init config: %w", err)
			}
			application.Info("Initial configuration complete")

			return nil
		},
	})
}
