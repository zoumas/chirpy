package main

import (
	"log"
	"net/http"
)

func main() {
	env, err := LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	app := &App{
		env: env,
	}

	mux := http.NewServeMux()
	registerRoutes(mux, app)

	server := &http.Server{
		Addr:    ":" + env.port,
		Handler: withLogger(withCORS(mux)),
	}
	log.Printf("STATUS: chirpy serving on port:%s", env.port)
	log.Fatal(server.ListenAndServe())
}
