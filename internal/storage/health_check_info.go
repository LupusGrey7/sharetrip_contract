package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// HealthQuerier is satisfied by *pgxpool.Pool; extracted for unit tests.
type HealthQuerier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// HealthRepository checks that PostgreSQL is reachable.
type HealthRepository interface {
	CheckHealth(ctx context.Context) error
}

type InfoRepository struct {
	db HealthQuerier
}

func NewInfoRepository(pool *pgxpool.Pool) *InfoRepository {
	return &InfoRepository{db: pool}
}

// NewInfoRepositoryForTest wires a fake querier for unit tests (no real Postgres).
func NewInfoRepositoryForTest(db HealthQuerier) *InfoRepository {
	return &InfoRepository{db: db}
}

// CheckHealth verifies DB connectivity with a trivial round-trip.
// Prefer SELECT 1 / Ping over pg_stat_activity: we only need "is DB up", not session metrics.
func (r *InfoRepository) CheckHealth(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var one int
	if err := r.db.QueryRow(ctx, `SELECT 1`).Scan(&one); err != nil {
		return fmt.Errorf("database healthcheck failed: %w", err)
	}
	if one != 1 {
		return fmt.Errorf("database healthcheck unexpected result: %d", one)
	}
	return nil
}
