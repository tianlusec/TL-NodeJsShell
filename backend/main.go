package main

import (
	"log"
	"NodeJsshell/app"
	"NodeJsshell/config"
)

func main() {
	cfg := config.Load()
	app := app.NewApp(cfg)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
