package env

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Env bundles all the environment variables the server needs to run.
type Env struct {
	Port           string
	FileserverPath string
}

// Load loads the environment variables into a struct.
// If the server is run with a -local flag then the environment is loaded from a .env file using godotenv.
func Load() (*Env, error) {
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
		return nil, envNotFound("PORT")
	}

	fileserverPath, ok := os.LookupEnv("FS_PATH")
	if !ok {
		return nil, envNotFound("FS_PATH")
	}

	return &Env{
		Port:           port,
		FileserverPath: fileserverPath,
	}, nil
}

// MustLoad uses the same semantics as Load but panics on error.
func MustLoad() *Env {
	local := flag.Bool("local", false, "Depend on the .env file for local development")
	flag.Parse()

	if *local {
		err := godotenv.Load()
		if err != nil {
			panic(fmt.Errorf("failed to load environment from .env : %s", err))
		}
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		panic(envNotFound("PORT"))
	}

	fileserverPath, ok := os.LookupEnv("FS_PATH")
	if !ok {
		panic(envNotFound("FS_PATH"))
	}

	return &Env{
		Port:           port,
		FileserverPath: fileserverPath,
	}
}

func envNotFound(name string) error {
	return fmt.Errorf("%s environment variable is not set", name)
}
