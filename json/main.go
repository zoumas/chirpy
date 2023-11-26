package main

import (
	"log"

	"github.com/zoumas/chirpy/json/internal/app"
	"github.com/zoumas/chirpy/json/internal/database"
	"github.com/zoumas/chirpy/json/internal/env"
)

func main() {
	env, err := env.Load()
	if err != nil {
		log.Fatalf("failed to load configuration : %s", err)
	}
	db, err := database.New(env.DSN)
	if err != nil {
		log.Fatalf("failed to connect to database : %s", err)
	}

	app := app.New(env, db)
	server := ConfiguredServer(app)
	app.Run(server)
}
