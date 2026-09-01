package storage_test

import (
	"context"
	"errors"
	"testing"

	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5"
)

// --- fakes for pgx.Row / QueryRow without a real database ---

type fakeRow struct {
	err error
	val int
}

func (r fakeRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	if len(dest) == 0 {
		return errors.New("no dest")
	}
	p, ok := dest[0].(*int)
	if !ok {
		return errors.New("dest[0] is not *int")
	}
	*p = r.val
	return nil
}

type fakeDB struct {
	row pgx.Row
}

func (f fakeDB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return f.row
}

func newRepoWithRow(row pgx.Row) *storage.InfoRepository {
	return storage.NewInfoRepositoryForTest(fakeDB{row: row})
}

func TestInfoRepository_CheckHealth_OK(t *testing.T) {
	repo := newRepoWithRow(fakeRow{val: 1})
	if err := repo.CheckHealth(context.Background()); err != nil {
		t.Fatalf("CheckHealth: %v", err)
	}
}

func TestInfoRepository_CheckHealth_QueryError(t *testing.T) {
	repo := newRepoWithRow(fakeRow{err: errors.New("connection refused")})
	err := repo.CheckHealth(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestInfoRepository_CheckHealth_UnexpectedValue(t *testing.T) {
	repo := newRepoWithRow(fakeRow{val: 0})
	err := repo.CheckHealth(context.Background())
	if err == nil {
		t.Fatal("expected error for unexpected SELECT 1 result")
	}
}
