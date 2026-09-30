package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnvFileForAppEnv(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"dev", DefaultEnvDevFile},
		{"development", DefaultEnvDevFile},
		{"", DefaultEnvDevFile},
		{"prod", DefaultEnvProdFile},
		{"production", DefaultEnvProdFile},
		{"test", DefaultEnvTestFile},
		{"staging", ".env.staging"},
	}
	for _, tt := range tests {
		if got := EnvFileForAppEnv(tt.in); got != tt.want {
			t.Fatalf("EnvFileForAppEnv(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestLoadEnvFile_ShouldFillProcessEnvFromFile(t *testing.T) {
	setAppEnv(t, nil)
	unsetForTest(t, httpPortEnv, databaseDSNEnv)
	path := writeEnvFile(t, "HTTP_PORT=8081\nDATABASE_DSN=postgres://file\n")

	loaded, err := LoadEnvFile(path)
	if err != nil || !loaded {
		t.Fatalf("LoadEnvFile = (%v, %v), want (true, nil)", loaded, err)
	}

	cfg, err := LoadAppConfig()
	if err != nil {
		t.Fatalf("LoadAppConfig: %v", err)
	}
	if cfg.HTTPPort != "8081" || cfg.DatabaseDSN != "postgres://file" {
		t.Fatalf("got HTTPPort=%q DatabaseDSN=%q, want values from file", cfg.HTTPPort, cfg.DatabaseDSN)
	}
}

func TestLoadEnvFile_ProcessEnvWinsOverFile(t *testing.T) {
	setAppEnv(t, map[string]string{databaseDSNEnv: "postgres://process"})
	path := writeEnvFile(t, "DATABASE_DSN=postgres://file\n")

	if _, err := LoadEnvFile(path); err != nil {
		t.Fatalf("LoadEnvFile: %v", err)
	}

	cfg, err := LoadAppConfig()
	if err != nil {
		t.Fatalf("LoadAppConfig: %v", err)
	}
	if cfg.DatabaseDSN != "postgres://process" {
		t.Fatalf("DatabaseDSN = %q, want process value", cfg.DatabaseDSN)
	}
}

func TestLoadEnvFile_MissingFileIsNotAnError(t *testing.T) {
	loaded, err := LoadEnvFile(filepath.Join(t.TempDir(), "missing.env"))
	if err != nil || loaded {
		t.Fatalf("LoadEnvFile = (%v, %v), want (false, nil)", loaded, err)
	}
}

func writeEnvFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), ".env.test")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write env file: %v", err)
	}
	return path
}

// unsetForTest removes keys so godotenv.Load can set them; t.Setenv registers the restore first.
func unsetForTest(t *testing.T, keys ...string) {
	t.Helper()
	for _, key := range keys {
		t.Setenv(key, "")
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("unset %s: %v", key, err)
		}
	}
}
