package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/joho/godotenv"
)

const (
	DefaultEnvDevFile  = ".env.dev"  // DefaultEnvDevFile is an optional local development file.
	DefaultEnvProdFile = ".env.prod" // DefaultEnvProdFile is optional locally; production should inject process env.
	DefaultEnvTestFile = ".env.test" // DefaultEnvTestFile is optional; unit tests must not depend on a local file.

	// AppEnvVar selects which dotenv file to load (dev|prod|test). Same as Make APP_ENV.
	AppEnvVar     = "APP_ENV"
	DefaultAppEnv = "dev"
)

// EnvFileForAppEnv maps APP_ENV short name to dotenv filename.
// Accepted: dev|development → .env.dev, prod|production → .env.prod, test → .env.test.
func EnvFileForAppEnv(appEnv string) string {
	switch appEnv {
	case "prod", "production":
		return DefaultEnvProdFile
	case "test":
		return DefaultEnvTestFile
	case "dev", "development", "":
		return DefaultEnvDevFile
	default:
		// Allow custom profiles like staging → .env.staging
		return ".env." + appEnv
	}
}

// CurrentAppEnv returns APP_ENV from process env or DefaultAppEnv.
func CurrentAppEnv() string {
	v := os.Getenv(AppEnvVar)
	if v == "" {
		return DefaultAppEnv
	}
	return v
}

// LoadEnvFile copies keys from an optional dotenv file into the process environment.
// Load (not Overload): variables already set in the process (K8s envFrom, shell, CI) win over the file.
// A missing file is not an error: in Docker/Kubernetes there is no .env file at all.
func LoadEnvFile(path string) (bool, error) {
	if err := godotenv.Load(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("load env file %q: %w", path, err)
	}
	return true, nil
}
