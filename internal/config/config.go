package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// Config contains all process-level application settings.
type Config struct {
	Environment string `env:"ENV" env-default:"development"`
	HTTPPort    string `env:"HTTP_PORT" env-default:"8082"`
	Database    DatabaseConfig
}

// DatabaseConfig contains PostgreSQL connection settings.
type DatabaseConfig struct {
	Driver   string `env:"DB_DRIVER" env-default:"postgres"`
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Port     int    `env:"DB_PORT" env-default:"6547"`
	User     string `env:"DB_USER" env-default:"postgres"`
	Password string `env:"DB_PASSWORD" env-default:"password"`
	Name     string `env:"DB_NAME" env-default:"sharetrip_contract"`
	SSLMode  string `env:"DB_SSLMODE" env-default:"disable"`
}

// DSN builds a pgx-compatible PostgreSQL connection string.
func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s?sslmode=%s",
		c.Driver,
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		c.SSLMode,
	)
}

// Load parses typed application settings from process environment variables.
// CI and production should provide settings through the process environment.
func Load() (Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return Config{}, fmt.Errorf("read config from environment: %w", err)
	}

	return cfg, nil
}

// LoadWithEnvFile loads an optional local dotenv file and then parses Config.
// Existing process variables win over file values, which is required for CI.
func LoadWithEnvFile(envFile string) (Config, error) {
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil && !os.IsNotExist(err) {
			return Config{}, fmt.Errorf("load env file %q: %w", envFile, err)
		}
	}

	return Load()
}
