package config

import (
	"errors"
	"os"
	"strings"
)

const (
	environmentEnv     = "ENV"
	environmentDefault = "development"
	httpPortEnv        = "HTTP_PORT"
	httpPortDefault    = "8082"
	databaseDSNEnv     = "DATABASE_DSN"

	otelServiceNameEnv          = "OTEL_SERVICE_NAME"
	otelServiceNameDefault      = "sharetrip-contract"
	otelServiceVersionEnv       = "OTEL_SERVICE_VERSION"
	otelServiceVersionDefault   = "1.0.0"
	otelEnvironmentEnv          = "OTEL_ENVIRONMENT"
	otelEnvironmentDefault      = "local"
	otelExporterEndpointEnv     = "OTEL_EXPORTER_ENDPOINT"
	otelExporterEndpointDefault = "localhost:4319"
)

// ErrDatabaseDSNRequired is returned when DATABASE_DSN is missing or blank.
var ErrDatabaseDSNRequired = errors.New(databaseDSNEnv + " is required")

// Config is the runtime configuration loaded from process environment.
// Same keys work for local (.env.<APP_ENV> via godotenv), Docker, and Kubernetes envFrom.
// DB_* keys are not read here: they are used only by Makefile/goose.
type Config struct {
	Environment string
	HTTPPort    string
	DatabaseDSN string
	Tracing     TracingConfig
}

// TracingConfig holds OpenTelemetry exporter settings.
type TracingConfig struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	Endpoint       string
}

// LoadAppConfig reads process env; DATABASE_DSN has no default and must be set (fail-fast on start).
func LoadAppConfig() (Config, error) {
	cfg := Config{
		Environment: getEnv(environmentEnv, environmentDefault),
		HTTPPort:    getEnv(httpPortEnv, httpPortDefault),
		DatabaseDSN: getEnv(databaseDSNEnv, ""),
		Tracing: TracingConfig{
			ServiceName:    getEnv(otelServiceNameEnv, otelServiceNameDefault),
			ServiceVersion: getEnv(otelServiceVersionEnv, otelServiceVersionDefault),
			Environment:    getEnv(otelEnvironmentEnv, otelEnvironmentDefault),
			Endpoint:       getEnv(otelExporterEndpointEnv, otelExporterEndpointDefault),
		},
	}

	if cfg.DatabaseDSN == "" {
		return Config{}, ErrDatabaseDSNRequired
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
