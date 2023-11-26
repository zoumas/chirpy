package app

import (
	"log"
	"net/http"

	"github.com/zoumas/chirpy/json/internal/database"
	"github.com/zoumas/chirpy/json/internal/env"
)

// App is used to implement stateful handlers. It groups global state.
type App struct {
	Env             *env.Env
	DB              *database.DB
	ChirpRepository database.ChirpRepository
	UserRepository  database.UserRepository

	// FileServerHits is used to count the number of times the website
	// has been viewed since the server started.
	FileServerHits int
}

func New(env *env.Env, db *database.DB) *App {
	return &App{
		Env:             env,
		DB:              db,
		ChirpRepository: NewJSONChirpResository(db),
		UserRepository:  NewJSONUserRepository(db),
	}
}

func (app *App) Run(server *http.Server) {
	// TODO: implement graceful shutdown here
	log.Printf("serving from %s on port:%s", app.Env.FileserverPath, app.Env.Port)
	log.Fatal(server.ListenAndServe())
}
