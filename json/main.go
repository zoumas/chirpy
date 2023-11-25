package main

import (
	"github.com/zoumas/chirpy/json/internal/app"
	"github.com/zoumas/chirpy/json/internal/env"
)

func main() {
	app := app.New(env.MustLoad())
	server := ConfiguredServer(app)
	app.Run(server)
}
