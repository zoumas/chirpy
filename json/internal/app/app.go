package app

import (
	"log"
	"net/http"

	"github.com/zoumas/chirpy/json/internal/env"
)

// App is used to implement stateful handlers. It groups global state.
type App struct {
	Env            *env.Env
	FileServerHits int
}

func New(env *env.Env) *App {
	return &App{Env: env}
}

func (app *App) Run(server *http.Server) {
	log.Printf("serving from %s on port:%s", app.Env.FileserverPath, app.Env.Port)
	log.Fatal(server.ListenAndServe())
}
