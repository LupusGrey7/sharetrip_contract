package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	DevelopmentEnv  = "development"
	ProductionEnv   = "production"
	defaultHTTPPort = "8080"
)

// Config holds application settings loaded from environment variables.
type Config struct {
	DatabaseURL string
	HTTPPort    string
}

// Load reads .env from the working directory (if present), then builds Config from env.
// DATABASE_URL is required. HTTP_PORT defaults to 8080 when unset.
func LoadWithDevEnvFile() (Config, error) {
	return LoadWithEnvFile(DefaultEnvDevFile)
}

// LoadWithEnvFile loads the given dotenv file when it exists, then reads process env.
// Useful in tests. Pass "" to skip loading a file and use only os.Getenv.
func LoadWithEnvFile(envFile string) (Config, error) {
	if envFile != "" {
		if err := loadEnvFile(envFile); err != nil {
			return Config{}, err
		}
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = defaultHTTPPort
	}

	return Config{
		DatabaseURL: dbURL,
		HTTPPort:    port,
	}, nil
}

func loadEnvFile(path string) error {
	// First check if the file exists
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			// File does not exist - this is normal (optional file)
			return nil
		}
		return fmt.Errorf("stat env file %q: %w", path, err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read env file %q: %w", path, err)
	}

	parsed, err := godotenv.Unmarshal(string(data))
	if err != nil {
		return fmt.Errorf("parse env file %q: %w", path, err)
	}

	for k, v := range parsed {
		if _, exists := os.LookupEnv(k); !exists || os.Getenv(k) == "" {
			if err := os.Setenv(k, v); err != nil {
				return fmt.Errorf("set env %q: %w", k, err)
			}
		}
	}

	// Debug logging
	_, _ = fmt.Fprintf(os.Stderr, "DEBUG: Successfully loaded env from %s\n", path)
	_, _ = fmt.Fprintf(os.Stderr, "DEBUG: After load - DATABASE_URL=%q, HTTP_PORT=%q\n", os.Getenv("DATABASE_URL"), os.Getenv("HTTP_PORT"))
	return nil
}
