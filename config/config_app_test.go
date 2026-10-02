package config

import (
	"errors"
	"reflect"
	"testing"
)

var appEnvKeys = []string{
	environmentEnv, httpPortEnv, databaseDSNEnv,
	otelServiceNameEnv, otelServiceVersionEnv, otelEnvironmentEnv, otelExporterEndpointEnv,
}

// t.Setenv isolates tests from make/shell .env.dev; it cannot be combined with t.Parallel.
func setAppEnv(t *testing.T, env map[string]string) {
	t.Helper()
	for _, key := range appEnvKeys {
		t.Setenv(key, env[key])
	}
}

func TestLoadAppConfig_WhenDatabaseDSNEmptyOrBlank_ShouldFail(t *testing.T) {
	for _, dsn := range []string{"", "   "} {
		setAppEnv(t, map[string]string{databaseDSNEnv: dsn})

		if _, err := LoadAppConfig(); !errors.Is(err, ErrDatabaseDSNRequired) {
			t.Fatalf("DATABASE_DSN=%q: error = %v, want ErrDatabaseDSNRequired", dsn, err)
		}
	}
}

func TestLoadAppConfig_ShouldReadEnvOrDefaults(t *testing.T) {
	const dsn = "postgres://u:p@db:5432/db?sslmode=disable"
	tests := []struct {
		name string
		env  map[string]string
		want Config
	}{
		{
			name: "only DSN set - defaults",
			env:  map[string]string{databaseDSNEnv: dsn},
			want: Config{
				Environment: "development", HTTPPort: "8082", DatabaseDSN: dsn,
				Tracing: TracingConfig{
					ServiceName: "sharetrip-contract", ServiceVersion: "1.0.0",
					Environment: "local", Endpoint: "localhost:4319",
				},
			},
		},
		{
			name: "all set - env wins, values trimmed",
			env: map[string]string{
				environmentEnv: "production", httpPortEnv: " 9090 ", databaseDSNEnv: dsn,
				otelServiceNameEnv: "contract-test", otelServiceVersionEnv: "2.0.0",
				otelEnvironmentEnv: "k8s-local", otelExporterEndpointEnv: "host.docker.internal:4319",
			},
			want: Config{
				Environment: "production", HTTPPort: "9090", DatabaseDSN: dsn,
				Tracing: TracingConfig{
					ServiceName: "contract-test", ServiceVersion: "2.0.0",
					Environment: "k8s-local", Endpoint: "host.docker.internal:4319",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setAppEnv(t, tt.env)

			got, err := LoadAppConfig()
			if err != nil {
				t.Fatalf("LoadAppConfig: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}
