package main

import (
	"context"
	"job4j/share_trip_contract/internal/config"
	"job4j/share_trip_contract/internal/storage"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// init is invoked before main()
func init() {
	// Same selector as Makefile: APP_ENV=dev|prod|test → .env.dev / .env.prod / .env.test
	appEnv := config.CurrentAppEnv()
	envName := config.EnvFileForAppEnv(appEnv)

	cwd, err := os.Getwd()
	envFile := envName
	if err == nil {
		envFile = filepath.Join(cwd, envName)
	}
	if loadErr := godotenv.Overload(envFile); loadErr != nil {
		log.Printf("No env file at %s (APP_ENV=%s): %v", envFile, appEnv, loadErr)
	} else {
		log.Printf("Loaded env file %s (APP_ENV=%s)", envFile, appEnv)
	}
}

func main() {
	ctx := context.Background()

	cfg := readCfg()

	pool, err := storage.NewPool(ctx, cfg.DSN())
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	// logging connection to DB
	if pingErr := pool.Ping(ctx); pingErr != nil {
		log.Fatalf("failed to ping database: %v", pingErr)
	}
	log.Printf("Connected to database successfully")
}

func readCfg() storage.Config {
	return storage.Config{
		Host:     config.Env("DB_HOST", "localhost"),
		Port:     config.EnvInt("DB_PORT", 6547),
		User:     config.Env("DB_USER", "postgres"),
		Password: config.Env("DB_PASSWORD", "password"),
		DBName:   config.Env("DB_NAME", "sharetrip_contract"),
		SSLMode:  config.Env("DB_SSLMODE", "disable"),
	}
}
