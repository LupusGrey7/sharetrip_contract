package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_readsTypedEnvironment(t *testing.T) {
	t.Setenv("ENV", "test")
	t.Setenv("HTTP_PORT", "9090")
	t.Setenv("DB_HOST", "postgres.test")
	t.Setenv("DB_PORT", "7777")
	t.Setenv("DB_USER", "test_user")
	t.Setenv("DB_PASSWORD", "test_password")
	t.Setenv("DB_NAME", "contract_test")
	t.Setenv("DB_SSLMODE", "require")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if cfg.Environment != "test" {
		t.Fatalf("Environment = %q, want test", cfg.Environment)
	}
	if cfg.HTTPPort != "9090" {
		t.Fatalf("HTTPPort = %q, want 9090", cfg.HTTPPort)
	}
	if cfg.Database.Port != 7777 {
		t.Fatalf("Database.Port = %d, want 7777", cfg.Database.Port)
	}

	const wantDSN = "postgres://test_user:test_password@postgres.test:7777/contract_test?sslmode=require"
	if got := cfg.Database.DSN(); got != wantDSN {
		t.Fatalf("Database.DSN() = %q, want %q", got, wantDSN)
	}
}

func TestLoad_usesDefaultsWithoutLocalEnvFile(t *testing.T) {
	unsetEnv(t,
		"ENV",
		"HTTP_PORT",
		"DB_DRIVER",
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"DB_SSLMODE",
	)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if cfg.HTTPPort != "8080" {
		t.Fatalf("HTTPPort = %q, want 8080", cfg.HTTPPort)
	}
	if cfg.Database.Host != "localhost" {
		t.Fatalf("Database.Host = %q, want localhost", cfg.Database.Host)
	}
	if cfg.Database.Port != 6547 {
		t.Fatalf("Database.Port = %d, want 6547", cfg.Database.Port)
	}
}

func TestLoad_rejectsInvalidTypedValue(t *testing.T) {
	t.Setenv("DB_PORT", "not-a-number")

	if _, err := Load(); err == nil {
		t.Fatal("expected an error for invalid DB_PORT")
	}
}

func TestLoadWithEnvFile_readsOptionalLocalFile(t *testing.T) {
	unsetEnv(t, "HTTP_PORT", "DB_HOST", "DB_PORT")

	envPath := filepath.Join(t.TempDir(), ".env.test")
	content := []byte("HTTP_PORT=8081\nDB_HOST=postgres.from.file\nDB_PORT=7654\n")
	if err := os.WriteFile(envPath, content, 0o600); err != nil {
		t.Fatalf("write env file: %v", err)
	}

	cfg, err := LoadWithEnvFile(envPath)
	if err != nil {
		t.Fatalf("LoadWithEnvFile: %v", err)
	}

	if cfg.HTTPPort != "8081" {
		t.Fatalf("HTTPPort = %q, want 8081", cfg.HTTPPort)
	}
	if cfg.Database.Host != "postgres.from.file" {
		t.Fatalf("Database.Host = %q, want postgres.from.file", cfg.Database.Host)
	}
	if cfg.Database.Port != 7654 {
		t.Fatalf("Database.Port = %d, want 7654", cfg.Database.Port)
	}
}

func TestLoadWithEnvFile_processEnvironmentWins(t *testing.T) {
	t.Setenv("DB_HOST", "postgres.from.process")

	envPath := filepath.Join(t.TempDir(), ".env.test")
	if err := os.WriteFile(envPath, []byte("DB_HOST=postgres.from.file\n"), 0o600); err != nil {
		t.Fatalf("write env file: %v", err)
	}

	cfg, err := LoadWithEnvFile(envPath)
	if err != nil {
		t.Fatalf("LoadWithEnvFile: %v", err)
	}

	if cfg.Database.Host != "postgres.from.process" {
		t.Fatalf("Database.Host = %q, want process environment value", cfg.Database.Host)
	}
}

func TestLoadWithEnvFile_missingFileUsesEnvironment(t *testing.T) {
	t.Setenv("DB_HOST", "postgres.from.process")

	cfg, err := LoadWithEnvFile(filepath.Join(t.TempDir(), "missing.env"))
	if err != nil {
		t.Fatalf("LoadWithEnvFile: %v", err)
	}
	if cfg.Database.Host != "postgres.from.process" {
		t.Fatalf("Database.Host = %q, want postgres.from.process", cfg.Database.Host)
	}
}

func unsetEnv(t *testing.T, keys ...string) {
	t.Helper()

	for _, key := range keys {
		value, existed := os.LookupEnv(key)
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("unset %s: %v", key, err)
		}

		restoreKey := key
		restoreValue := value
		restoreExisted := existed
		t.Cleanup(func() {
			if restoreExisted {
				_ = os.Setenv(restoreKey, restoreValue)
				return
			}
			_ = os.Unsetenv(restoreKey)
		})
	}
}
