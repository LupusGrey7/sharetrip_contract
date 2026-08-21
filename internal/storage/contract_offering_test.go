package storage_test

import (
	"context"
	"testing"

	"job4j/sharetrip-contract/internal/contract/domain"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type codeRows struct {
	codes []string
	i     int
}

func (r *codeRows) Close()                                       {}
func (r *codeRows) Err() error                                   { return nil }
func (r *codeRows) CommandTag() pgconn.CommandTag                { return pgconn.CommandTag{} }
func (r *codeRows) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (r *codeRows) Next() bool {
	if r.i >= len(r.codes) {
		return false
	}
	r.i++
	return true
}
func (r *codeRows) Scan(dest ...any) error {
	*(dest[0].(*string)) = r.codes[r.i-1]
	return nil
}
func (r *codeRows) Values() ([]any, error) { return nil, nil }
func (r *codeRows) RawValues() [][]byte    { return nil }
func (r *codeRows) Conn() *pgx.Conn        { return nil }

type offeringFakeTx struct {
	queryRows pgx.Rows
	execErr   error
}

func (t offeringFakeTx) Begin(ctx context.Context) (pgx.Tx, error) { return t, nil }
func (t offeringFakeTx) BeginFunc(ctx context.Context, f func(pgx.Tx) error) error {
	return f(t)
}
func (t offeringFakeTx) Commit(ctx context.Context) error   { return nil }
func (t offeringFakeTx) Rollback(ctx context.Context) error { return nil }
func (t offeringFakeTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return 0, nil
}
func (t offeringFakeTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults { return nil }
func (t offeringFakeTx) LargeObjects() pgx.LargeObjects                               { return pgx.LargeObjects{} }
func (t offeringFakeTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	return nil, nil
}
func (t offeringFakeTx) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, t.execErr
}
func (t offeringFakeTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return t.queryRows, nil
}
func (t offeringFakeTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return nil
}
func (t offeringFakeTx) Conn() *pgx.Conn { return nil }

func TestExistServiceCodesTx_ReportsMissing(t *testing.T) {
	repo := storage.NewOfferingRepository(nil)
	tx := offeringFakeTx{queryRows: &codeRows{codes: []string{"trip_creation"}}}

	missing, err := repo.ExistServiceCodesTx(context.Background(), tx, []string{"trip_creation", "notifications"})
	if err != nil {
		t.Fatal(err)
	}
	if len(missing) != 1 || missing[0] != "notifications" {
		t.Fatalf("missing=%v, want [notifications]", missing)
	}
}

func TestUpsertContractServiceTx_OK(t *testing.T) {
	repo := storage.NewContractOfferingRepository(nil)
	tx := offeringFakeTx{}
	err := repo.UpsertContractServiceTx(context.Background(), tx, 1, domain.ServiceItemEntity{
		ServiceCode: "trip_creation",
		IsEnabled:   true,
	})
	if err != nil {
		t.Fatal(err)
	}
}
