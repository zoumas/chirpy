package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Env bundles all the environment variables the server needs to run.
type Env struct {
	port string
}

// Loadenv loads the environment variables to a struct.
// If the server is run with a -local flag then the environment is loaded from a .env file using godotenv.
func LoadEnv() (*Env, error) {
	local := flag.Bool("local", false, "Depend on the .env file for local development")
	flag.Parse()

	if *local {
		err := godotenv.Load()
		if err != nil {
			return nil, fmt.Errorf("failed to load environment from .env : %s", err)
		}
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		return nil, errors.New("PORT environment variable is not set")
	}

	return &Env{
		port: port,
	}, nil
}
