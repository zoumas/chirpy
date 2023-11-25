package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// ConfiguredRouter bundles the definitions and handling of the endpoints the server supports.
// It returns a *chi.Mux which is a http.Handler ready for use.
func ConfiguredRouter(app *App) http.Handler {
	mainRouter := chi.NewRouter()

	mainRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodOptions,
			http.MethodPut,
			http.MethodDelete,
		},
	}))
	mainRouter.Use(middleware.Logger)
	mainRouter.Use(middleware.Recoverer)

	setupFileserver(mainRouter, app)

	mainRouter.Post("/reset", app.ResetMetrics)
	mainRouter.Get("/metrics", app.ReportMetrics)

	mainRouter.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	})

	mainRouter.Post("/panic", func(_ http.ResponseWriter, _ *http.Request) {
		panic("testing middleware.Recoverer")
	})

	return mainRouter
}

func setupFileserver(router *chi.Mux, app *App) {
	root := http.Dir(app.env.fileserverPath)
	fileserver := app.IncrementMetrics(http.StripPrefix("/app", http.FileServer(root)))

	// Chi Behavior:
	// A request to /app/assets creates a duplicate request on both /app/assets and /app/assets/
	// when both /app and /app/* are handled, incrementing the fileserver hits twice. This is a bug.

	// Solved: Using a custom ResponseWriter that captures the statusCode from the returned function.
	// If the statusCode is a 301 MovedPermanently then we shouldn't increment the file server hits.

	router.Handle("/app", fileserver)
	router.Handle("/app/*", fileserver)
}
