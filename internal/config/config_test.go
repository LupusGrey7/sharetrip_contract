package config

import (
	"os"
	"path/filepath"
	"testing"
)

// clearConfigEnvVars drops keys that LoadWithEnvFile fills from a dotenv file.
// Needed when tests run under `make test`: Makefile includes .env.dev and exports
// DATABASE_URL/HTTP_PORT/ENV, and loadEnvFile does not override non-empty process env.
func clearConfigEnvVars() {
	_ = os.Unsetenv("DATABASE_URL")
	_ = os.Unsetenv("HTTP_PORT")
	_ = os.Unsetenv("ENV")
}

func TestMain(m *testing.M) {
	clearConfigEnvVars()
	if _, err := LoadTestEnvFile(); err != nil {
		_, _ = os.Stderr.WriteString("config tests: load " + DefaultEnvTestFile + ": " + err.Error() + "\n")
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestLoad_fromEnvTest(t *testing.T) {
	// Isolate from Make-exported .env.dev (and from other tests).
	t.Setenv("DATABASE_URL", "")
	t.Setenv("HTTP_PORT", "")
	t.Setenv("ENV", "")

	// go test runs with cwd = package dir; resolve from module root.
	path, err := ResolveEnvFile(DefaultEnvTestFile)
	if err != nil {
		t.Fatalf("ResolveEnvFile: %v", err)
	}

	cfg, err := LoadWithEnvFile(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	const wantURL = "postgres://test_user:test_password@localhost:6547/contract_test?sslmode=disable"
	if cfg.DatabaseURL != wantURL {
		t.Fatalf("DatabaseURL = %q, want %q", cfg.DatabaseURL, wantURL)
	}
	if cfg.HTTPPort != "8081" {
		t.Fatalf("HTTPPort = %q, want 8081 from .env.test", cfg.HTTPPort)
	}
	if os.Getenv("ENV") != "test" {
		t.Fatalf("ENV = %q, want test from .env.test", os.Getenv("ENV"))
	}
}

func TestLoadTest_explicitFile(t *testing.T) {
	path, err := ResolveEnvFile(DefaultEnvTestFile)
	if err != nil {
		t.Fatalf("ResolveEnvFile: %v", err)
	}

	cfg, err := LoadWithEnvFile(path)
	if err != nil {
		t.Fatalf("LoadWithEnvFile: %v", err)
	}
	if cfg.DatabaseURL == "" {
		t.Fatal("expected DATABASE_URL from " + DefaultEnvTestFile)
	}
}

func TestLoadWithEnvFile_missingFileUsesProcessEnv(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://test_user:test_password@localhost:6547/contract_test?sslmode=disable")
	t.Setenv("HTTP_PORT", "")

	cfg, err := LoadWithEnvFile(filepath.Join(t.TempDir(), "missing.env"))
	if err != nil {
		t.Fatalf("LoadWithEnvFile: %v", err)
	}

	if cfg.DatabaseURL != "postgres://test_user:test_password@localhost:6547/contract_test?sslmode=disable" {
		t.Fatalf("DatabaseURL = %q", cfg.DatabaseURL)
	}
	if cfg.HTTPPort != "8080" {
		t.Fatalf("HTTPPort = %q, want default 8080", cfg.HTTPPort)
	}
}

// for example, with description and debug logs to stderr
func TestLoadWithEnvFile_fromDotEnv(t *testing.T) {
	// === ПОДГОТОВКА (Arrange) ===
	// Создаём временную папку и .env файл с кастомными значениями
	// Это нужно, чтобы тестировать загрузку конфига из конкретного файла
	dir := t.TempDir()
	envPath := filepath.Join(dir, "custom.env")

	// Формируем содержимое файла с переменными окружения
	// Каждая переменная должна быть на отдельной строке в формате KEY=VALUE
	content := "DATABASE_URL=postgres://test_user:test_password@localhost:6547/contract_test?sslmode=disable\nHTTP_PORT=9090\n"
	if err := os.WriteFile(envPath, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Logf("Created env file at: %s", envPath) //test logs

	// ВАЖНО: очищаем переменные окружения ПЕРЕД загрузкой конфига!
	// Это гарантирует, что LoadWithEnvFile будет читать ВСЕ значения ТОЛЬКО из файла,
	// а не из предыдущих тестов или системного окружения.
	// Используем t.Setenv() которая автоматически восстановит значения после теста.
	t.Setenv("DATABASE_URL", "")
	t.Setenv("HTTP_PORT", "")

	// LoadWithEnvFile fills only missing/empty keys (process env wins).
	// Empty Setenv above lets values from the file apply.
	cfg, err := LoadWithEnvFile(envPath)
	if err != nil {
		t.Fatalf("LoadWithEnvFile: %v", err)
	}
	t.Logf("Loaded config successfully: DATABASE_URL=%q, HTTP_PORT=%q", cfg.DatabaseURL, cfg.HTTPPort)

	// === ПРОВЕРКА (Assert) ===
	// Проверяем что DATABASE_URL был успешно загружен из файла
	// Важно проверить полный URL с базой данных и параметрами подключения
	wantURL := "postgres://test_user:test_password@localhost:6547/contract_test?sslmode=disable"
	if cfg.DatabaseURL != wantURL {
		t.Fatalf("DatabaseURL = %q, want %q", cfg.DatabaseURL, wantURL)
	}

	// Проверяем что HTTP_PORT был успешно загружен из файла
	// и имеет значение 9090, а не значение по умолчанию (8080)
	if cfg.HTTPPort != "9090" {
		t.Fatalf("HTTPPort = %q, want 9090", cfg.HTTPPort)
	}
}

func TestLoadWithEnvFile_missingDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("HTTP_PORT", "8080")

	_, err := LoadWithEnvFile("")
	if err == nil {
		t.Fatal("expected error when DATABASE_URL is empty")
	}
}

func TestLoadWithEnvFile_processEnvOverridesDotEnv(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, "override.env")
	if err := os.WriteFile(envPath, []byte("DATABASE_URL=postgres://from-file:test_password@localhost:6547/contract_test?sslmode=disable\nHTTP_PORT=3000\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	// re init env
	t.Setenv("DATABASE_URL", "postgres://from-process")
	t.Setenv("HTTP_PORT", "4000")

	cfg, err := LoadWithEnvFile(envPath)
	if err != nil {
		t.Fatalf("LoadWithEnvFile: %v", err)
	}
	t.Logf("After load susscesfully: DATABASE_URL=%q, HTTP_PORT=%q", cfg.DatabaseURL, cfg.HTTPPort)

	if cfg.DatabaseURL != "postgres://from-process" {
		t.Fatalf("DatabaseURL = %q, want process env to win", cfg.DatabaseURL)
	}
	if cfg.HTTPPort != "4000" {
		t.Fatalf("HTTPPort = %q, want process env to win", cfg.HTTPPort)
	}
}
