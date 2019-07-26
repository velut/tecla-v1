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
	app := app.NewApp(app.DefaultOptions())
	return app.Run()
}
