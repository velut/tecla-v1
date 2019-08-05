package main

import (
	"log"

	"github.com/velut/tecla/server/pkg/app"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln("Tecla cannot start:", err)
	}
}

func run() error {
	options := app.DefaultOptions()
	options.ProductionMode = false

	app := app.NewApp(options)
	return app.Run()
}
