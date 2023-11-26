package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/zoumas/chirpy/json/internal/app"
)

// ConfiguredRouter bundles the definitions and handling of the endpoints the server supports.
// It returns a *chi.Mux which is a http.Handler ready for use.
func ConfiguredRouter(app *app.App) *chi.Mux {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
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
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Mount("/api", ConfiguredApiRouter(app))
	router.Mount("/app", ConfiguredAppRouter(app))
	router.Mount("/admin", ConfiguredAdminRouter(app))

	return router
}

func ConfiguredApiRouter(app *app.App) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/reset", app.ResetMetrics)

	router.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	})

	router.Post("/chirps", app.CreateChirp)
	router.Get("/chirps", app.GetChirps)

	return router
}

func ConfiguredAppRouter(app *app.App) *chi.Mux {
	router := chi.NewRouter()

	root := http.Dir(app.Env.FileserverPath)
	fileserver := app.IncrementMetrics(http.StripPrefix("/app", http.FileServer(root)))

	router.Handle("/", fileserver)
	router.Handle("/*", fileserver)

	return router
}

func ConfiguredAdminRouter(app *app.App) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/metrics", app.ReportMetrics)

	return router
}
