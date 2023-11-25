package main

import (
	"net/http"
	"time"

	"github.com/zoumas/chirpy/json/internal/app"
)

func ConfiguredServer(app *app.App) *http.Server {
	return &http.Server{
		Addr:    ":" + app.Env.Port,
		Handler: ConfiguredRouter(app),

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  time.Minute,
	}
}
