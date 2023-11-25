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

	server := &http.Server{
		Addr:    ":" + env.port,
		Handler: ConfiguredRouter(app),
	}

	log.Printf("STATUS: chirpy serving on port:%s", env.port)
	log.Fatal(server.ListenAndServe())
}
