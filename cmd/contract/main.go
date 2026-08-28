package main

import (
	"context"
	"fmt"
	"job4j/sharetrip-contract/internal/api"
	"job4j/sharetrip-contract/internal/app"
	"job4j/sharetrip-contract/internal/config"
	"job4j/sharetrip-contract/internal/storage"
	"log/slog"
	"os"
	"time"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	if err := run(); err != nil {
		slog.Error("application stopped", slog.Any("error", err))
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	appEnv := config.CurrentAppEnv()
	envFile := config.EnvFileForAppEnv(appEnv)
	cfg, err := config.LoadWithEnvFile(envFile)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	pool, err := storage.NewPool(ctx, cfg.Database.DSN())
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()

	slog.Info("connected to database")

	// init Tracing (OpenTelemetry → otel-collector → Jaeger)
	tp, err := app.InitTracing(ctx)
	if err != nil {
		slog.Error("init tracing failed", "error", err)
		os.Exit(1)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if shutdownErr := tp.Shutdown(shutdownCtx); shutdownErr != nil {
			slog.Error("shutdown tracing failed", "error", shutdownErr)
		}
	}()

	fiberApp := app.New(pool)
	addr := ":" + cfg.HTTPPort

	api.LogRegisteredRoutes(addr)

	slog.Info("server listening", slog.String("address", addr))
	if err := fiberApp.Listen(addr); err != nil {
		return fmt.Errorf("listen on %s: %w", addr, err)
	}

	return nil
}
