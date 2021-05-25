// Copyright © 2021 The Sanuscoin Team

package main

import (
	"log"

	"github.com/goava/di"

	"sanuscoin/sanuscoin-sdk/app"
	"sanuscoin/sanuscoin-sdk/app/context"
	"sanuscoin/sanuscoin-sdk/config"
)

func main() {
	// create the application DI-container
	c, err := di.New(
		// provide the application
		di.Provide(app.NewApp),
		// provide the application's context
		di.Provide(context.NewContext),
		// application's providers
		di.Provide(config.NewConfig), // provide the application's config
		// application's invokers
		di.Invoke(config.Load), // invoke config loader
	)
	if err != nil {
		log.Fatal(err)
	}

	// invoke application starter
	if err = c.Invoke(app.Start); err != nil {
		log.Fatal(err)
	}
}
