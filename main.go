// Copyright © 2021 The Sanuscoin Team

package main

import (
	"log"

	"sanus/sanus-sdk/sanus/daemon"
	"sanus/sanus-sdk/sanus/sdk"

	"github.com/goava/di"

	"sanus/sanus-sdk/app"
	"sanus/sanus-sdk/app/context"
	"sanus/sanus-sdk/config"

	sanusHttp "sanus/sanus-sdk/network/http"
)

func main() {
	// create the application DI-container
	c, err := di.New(
		// provide the application
		di.Provide(app.NewApp),
		// provide the application's context
		di.Provide(context.NewContext),
		// provide the application's config
		di.Provide(config.NewConfig),
		// provide the application wallet
		di.Provide(sdk.NewWallet),
		// provide the application http server
		di.Provide(sanusHttp.NewHTTP),
		// provide the application btcd service
		di.Provide(daemon.NewBTCDaemon),
	)
	if err != nil {
		log.Fatal(err)
	}

	// invoke application starter
	if err = c.Invoke(app.Start); err != nil {
		log.Fatal(err)
	}
}
