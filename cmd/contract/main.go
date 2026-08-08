package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"job4j/sharetrip-contract/internal/app"
	"job4j/sharetrip-contract/internal/config"
	httpserver "job4j/sharetrip-contract/internal/http"
	"job4j/sharetrip-contract/internal/storage"

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

	log.Printf("Connected to database successfully")

	fiberApp := app.New(pool)
	addr := fmt.Sprintf(":%s", config.Env("HTTP_PORT", "8080"))

	httpserver.LogRegisteredRoutes(addr)

	log.Printf("listening on %s", addr)
	if err := fiberApp.Listen(addr); err != nil {
		log.Fatal(err)
	}
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
