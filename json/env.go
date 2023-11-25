package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Env bundles all the environment variables the server needs to run.
type Env struct {
	port           string
	fileserverPath string
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
		return nil, EnvNotFound("PORT")
	}

	fileserverPath, ok := os.LookupEnv("FS_PATH")
	if !ok {
		return nil, EnvNotFound("FS_PATH")
	}

	return &Env{
		port:           port,
		fileserverPath: fileserverPath,
	}, nil
}

func EnvNotFound(name string) error {
	return fmt.Errorf("%s environment variable is not set", name)
}
