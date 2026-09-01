// transaction function - package service

package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func tx[T interface{}](
	ctx context.Context,
	pool *pgxpool.Pool,
	block func(tx pgx.Tx) (*T, error),
) (*T, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
		slog.String("layer", "transaction"),
	)
	logger.Debug("begin transaction")

	txBegin, err := pool.Begin(ctx)
	if err != nil {
		logger.Error("failed to begin transaction", slog.Any("error", err))
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	res, err := block(txBegin)
	if err != nil {
		logger.Error("transaction block failed", slog.Any("error", err))
		if rbErr := txBegin.Rollback(ctx); rbErr != nil {
			logger.Error("rollback transaction failed", slog.Any("error", rbErr))
		}
		return nil, err
	}

	if err = txBegin.Commit(ctx); err != nil {
		logger.Error("failed to commit transaction", slog.Any("error", err))
		if rbErr := txBegin.Rollback(ctx); rbErr != nil {
			logger.Error("rollback transaction failed", slog.Any("error", rbErr))
		}
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Info("commit transaction")
	return res, nil
}
