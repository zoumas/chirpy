package main

import "net/http"

// registerRoutes bundles the definitions and handling of the endpoints the server supports.
func registerRoutes(mux *http.ServeMux, env *Env) {
	fileserver := http.StripPrefix("/app", http.FileServer(http.Dir(env.fileserverPath)))
	mux.Handle("/app/", fileserver)

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	})
}
