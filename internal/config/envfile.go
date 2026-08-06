package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	DefaultEnvDevFile  = ".env.dev"  // DefaultEnvDevFile is used by the running application (development secrets).
	DefaultEnvProdFile = ".env.prod" // DefaultEnvProdFile is used by the running application (production secrets).
	DefaultEnvTestFile = ".env.test" // DefaultEnvTestFile is committed; used by unit/integration tests (no real secrets).

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

// ModuleRoot returns the directory that contains go.mod by walking up from the working directory.
func ModuleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getwd: %w", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %s", dir)
		}
		dir = parent
	}
}

// ResolveEnvFile returns an absolute path to the named dotenv file in the module root.
func ResolveEnvFile(name string) (string, error) {
	root, err := ModuleRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, name), nil
}

// LoadTest loads variables from .env.test at the module root, then builds Config.
func LoadTestEnvFile() (Config, error) {
	path, err := ResolveEnvFile(DefaultEnvTestFile)
	if err != nil {
		return Config{}, err
	}
	return LoadWithEnvFile(path)
}
