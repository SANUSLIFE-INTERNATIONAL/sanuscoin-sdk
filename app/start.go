// Copyright © 2021 The Sanuscoin Team

package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"sanus/sanus-sdk/app/context"
	"sanus/sanus-sdk/config"
	"sanus/sanus-sdk/misc/disk"
	sanusHttp "sanus/sanus-sdk/network/http"
	sanusRPC "sanus/sanus-sdk/network/rpc"
	"sanus/sanus-sdk/sanus/daemon"

	"github.com/goava/di"
	"github.com/urfave/cli/v2"
)

const (
	defaultLogFile = "app.log"
)

// interruptSignals defines the default signals to catch in order to do a proper
// shutdown.  This may be modified during init depending on the platform.
var interruptSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

// startCommand appends start action to application.
func (application *App) startCommand(dic *di.Container, ctx context.Context, cfg *config.Config, app *App) {
	app.Commands = append(app.Commands, &cli.Command{
		Name:  "start",
		Usage: "Start node",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        config.AppTestnetName,
				Usage:       fmt.Sprintf("Use %v (use config otherwise)", config.AppTestnetName),
				Destination: &cfg.Net.Testnet,
			},
		},
		Before: func(cc *cli.Context) error {

			// determine network scope
			cfg.Net.Testnet = cfg.Net.Testnet || cc.Bool(config.AppTestnetName)

			if err := dic.Invoke(config.InitPaths); err != nil {
				return err
			}

			application.Debugf("Checking root directory (%v)", config.AppRootPath())
			application.Debugf("Checking config file (%v)", config.AppConfigFile())

			if !disk.DirExists(config.AppRootPath()) || !disk.FileExists(config.AppConfigFile()) {
				return fmt.Errorf("error caused when trying to read root path [%v]", config.AppRootPath())
			}

			application.initLogger()

			return nil
		},
		Action: func(c *cli.Context) error {
			var (
				interruptCh = application.interruptListener()

				btcStopSig    = make(chan struct{}, 1)
				httpStopSig   = make(chan struct{}, 1)
				rpcStopSig   = make(chan struct{}, 1)
				walletStopSig = make(chan struct{}, 1)

				appStopSig = make(chan struct{}, 1)
			)

			go func() {
				<-interruptCh
				close(btcStopSig)
				application.Info("Stopping BTCD service")

				close(walletStopSig)
				application.Info("Stopping WALLET service")

				close(httpStopSig)
				application.Info("Stopping HTTP server")

				close(rpcStopSig)
				application.Info("Stopping RPC server")

				close(appStopSig)
			}()

			var httpSrv *sanusHttp.HTTPServer
			if err := dic.Resolve(&httpSrv); err != nil {
				return fmt.Errorf("can't fetch http-server from system | %v", err)
			}
			go httpSrv.Serve(httpStopSig)

			var rpcSrv *sanusRPC.RPCServer
			if err := dic.Resolve(&rpcSrv); err != nil {
				return fmt.Errorf("can't fetch http-server from system | %v", err)
			}
			go rpcSrv.Serve(rpcStopSig)

			var btcdService *daemon.BTCDaemon
			if err := dic.Resolve(&btcdService); err != nil {
				return fmt.Errorf("can't fetch btcd service from system | %v", err)
			}
			if err := btcdService.Start(btcStopSig); err != nil {
				return fmt.Errorf("can't start btcd service | %v", err)
			}
			<-appStopSig
			return nil
		},
		After: func(cc *cli.Context) error {
			// wait while context canceled
			<-cc.Done()
			// wait while all workers finished
			ctx.Cancel()
			ctx.WgWait()
			//logger.Info("Application shutdown complete")
			return nil
		},
	})
}

// interruptListener listens for OS Signals such as SIGINT (Ctrl+C) and shutdown
// requests from shutdownRequestChannel.  It returns a channel that is closed
// when either signal is received.
func (application *App) interruptListener() <-chan struct{} {
	c := make(chan struct{})
	go func() {
		interruptChannel := make(chan os.Signal, 1)
		signal.Notify(interruptChannel, interruptSignals...)

		// Listen for initial shutdown signal and close the returned
		// channel to notify the caller.
		sig := <-interruptChannel
		application.Infof("Received %s signal, shutting down...", sig)
		close(c)
	}()
	return c
}
