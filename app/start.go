// Copyright © 2021 The Sanuscoin Team

package app

import (
	"fmt"
	"log"

	"sanus/sanus-sdk/sanus/daemon"

	"github.com/goava/di"
	"github.com/urfave/cli/v2"

	"sanus/sanus-sdk/app/context"
	"sanus/sanus-sdk/config"
)

// startCommand appends start action to application.
func startCommand(dic *di.Container, ctx context.Context, cfg *config.Config, app *App) {
	app.Commands = append(app.Commands, &cli.Command{
		Name:  "start",
		Usage: "Start node",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  config.AppTestnetName,
				Usage: fmt.Sprintf("Use %v (use config otherwise)", config.AppTestnetName),
			},
		},
		Before: func(cc *cli.Context) error {
			// determine network scope
			cfg.Net.Testnet = cfg.Net.Testnet || cc.Bool(config.AppTestnetName)
			// invoke config maker
			if err := dic.Invoke(config.Make); err != nil {
				return fmt.Errorf("make config: %w", err)
			}
			return nil
		},
		Action: func(c *cli.Context) error {

			go func() {
				daemon.Run(nil)
			}()

			return nil
		},
		After: func(cc *cli.Context) error {
			// wait while context canceled
			<-cc.Done()
			// wait while all workers finished
			ctx.Cancel()
			ctx.WgWait()
			log.Println("Application shutdown complete")

			return nil
		},
	})
}
