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

	mux := http.NewServeMux()

	fileserver := http.FileServer(http.Dir(env.fileserverPath))
	mux.Handle("/", fileserver)

	server := &http.Server{
		Addr:    ":" + env.port,
		Handler: withLogger(withCORS(mux)),
	}
	log.Printf("STATUS: chirpy serving on port:%s", env.port)
	log.Fatal(server.ListenAndServe())
}
