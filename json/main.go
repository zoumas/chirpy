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

	server := &http.Server{
		Addr:    ":" + env.port,
		Handler: withLogger(withCORS(http.NewServeMux())),
	}
	log.Printf("STATUS: chirpy serving on port:%s", env.port)
	log.Fatal(server.ListenAndServe())
}
