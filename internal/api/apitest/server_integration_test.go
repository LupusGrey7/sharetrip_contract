// Package apitest is a separate Go package on purpose.
// Moving this file into internal/api (package api) would import application.New,
// and app already imports api — that is an import cycle.
// Keep the file here until we extract a cycle-free test helper.
package apitest

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"job4j/sharetrip-contract/internal/api"
	application "job4j/sharetrip-contract/internal/app"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// TestAPIWithPostgres verifies the real HTTP -> service -> usecase -> storage -> PostgreSQL chain.
// It does not read DATABASE_URL or .env.test: the DSN comes directly from Testcontainers.
func TestAPIWithPostgres(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	t.Cleanup(cancel)

	pgContainer, err := postgres.Run(
		ctx,
		"postgres:17",
		postgres.WithDatabase("contract_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("start PostgreSQL container: %v", err)
	}
	t.Cleanup(func() {
		terminateCtx, terminateCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer terminateCancel()
		if err := pgContainer.Terminate(terminateCtx); err != nil {
			t.Errorf("terminate PostgreSQL container: %v", err)
		}
	})

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("get PostgreSQL connection string: %v", err)
	}

	applyMigrations(t, ctx, dsn)

	pool, err := storage.NewPool(ctx, dsn)
	if err != nil {
		t.Fatalf("create PostgreSQL pool: %v", err)
	}
	t.Cleanup(pool.Close)

	fiberApp := application.New(pool)

	t.Run("healthcheck uses real database", func(t *testing.T) {
		resp := sendRequest(t, fiberApp, http.MethodGet, "/healthcheck", nil)
		defer closeBody(t, resp)

		requireStatus(t, resp, http.StatusOK)

		var got api.HealthcheckResponse
		decodeResponse(t, resp, &got)
		if got.Status != "ok" {
			t.Fatalf("status=%q, want ok", got.Status)
		}
	})

	t.Run("contract lifecycle enables trip creation", func(t *testing.T) {
		const companyID = 870001

		createBody := []byte(`{
			"company_id": 870001,
			"contract_number": "IT-CONTRACT-870001",
			"status": "active",
			"start_date": "2026-01-01T00:00:00Z",
			"end_date": "2027-01-01T00:00:00Z"
		}`)
		createResp := sendRequest(t, fiberApp, http.MethodPost, "/api/v2/contracts/", createBody)
		requireStatus(t, createResp, http.StatusCreated)

		var created api.ContractResponse
		decodeResponse(t, createResp, &created)
		closeBody(t, createResp)
		if created.ID <= 0 || created.CompanyID != companyID || created.Status != "active" {
			t.Fatalf("unexpected created contract: %+v", created)
		}

		getResp := sendRequest(
			t,
			fiberApp,
			http.MethodGet,
			fmt.Sprintf("/api/v2/contracts/%d", created.ID),
			nil,
		)
		requireStatus(t, getResp, http.StatusOK)

		var found api.ContractResponse
		decodeResponse(t, getResp, &found)
		closeBody(t, getResp)
		if found.ID != created.ID || found.ContractNumber != created.ContractNumber {
			t.Fatalf("unexpected fetched contract: %+v", found)
		}

		activeResp := sendRequest(
			t,
			fiberApp,
			http.MethodGet,
			fmt.Sprintf("/api/v2/contracts/active?companyId=%d", companyID),
			nil,
		)
		requireStatus(t, activeResp, http.StatusOK)
		var active api.ContractResponse
		decodeResponse(t, activeResp, &active)
		closeBody(t, activeResp)
		if active.ID != created.ID || active.Status != "active" || active.CompanyID != companyID {
			t.Fatalf("unexpected active contract: %+v", active)
		}

		upsertBody := fmt.Appendf(nil, `{
			"contract_id": %d,
			"services": [{"service_code": "trip_creation", "is_enabled": true}]
		}`, created.ID)
		upsertResp := sendRequest(t, fiberApp, http.MethodPut, "/api/v2/services", upsertBody)
		requireStatus(t, upsertResp, http.StatusOK)

		var upserted api.UpsertServicesResponse
		decodeResponse(t, upsertResp, &upserted)
		closeBody(t, upsertResp)
		if upserted.ContractID != created.ID || len(upserted.Services) != 1 || !upserted.Services[0].IsEnabled {
			t.Fatalf("unexpected upsert response: %+v", upserted)
		}

		availabilityResp := sendRequest(
			t,
			fiberApp,
			http.MethodGet,
			fmt.Sprintf(
				"/api/v2/companies/%d/services/trip_creation/availability",
				companyID,
			),
			nil,
		)
		defer closeBody(t, availabilityResp)
		requireStatus(t, availabilityResp, http.StatusOK)

		var availability api.AvailabilityResult
		decodeResponse(t, availabilityResp, &availability)
		if !availability.Allowed || availability.CompanyID != companyID {
			t.Fatalf("unexpected availability: %+v", availability)
		}
	})
}

func applyMigrations(t *testing.T, ctx context.Context, dsn string) {
	t.Helper()

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open migration database: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close migration database: %v", err)
		}
	})

	if err := goose.SetDialect("postgres"); err != nil {
		t.Fatalf("set goose dialect: %v", err)
	}
	if err := goose.UpContext(ctx, db, filepath.Join(moduleRoot(t), "migrations")); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
}

func moduleRoot(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("go.mod not found from %s", dir)
		}
		dir = parent
	}
}

func sendRequest(t *testing.T, app *fiber.App, method, path string, body []byte) *http.Response {
	t.Helper()

	req, err := http.NewRequest(method, path, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	if body != nil {
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	}

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	return resp
}

func requireStatus(t *testing.T, resp *http.Response, want int) {
	t.Helper()
	if resp.StatusCode == want {
		return
	}
	body, _ := io.ReadAll(resp.Body)
	t.Fatalf("status=%d, want=%d body=%s", resp.StatusCode, want, body)
}

func decodeResponse(t *testing.T, resp *http.Response, target any) {
	t.Helper()
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		t.Fatalf("decode response: %v", err)
	}
}

func closeBody(t *testing.T, resp *http.Response) {
	t.Helper()
	if err := resp.Body.Close(); err != nil {
		t.Errorf("close response body: %v", err)
	}
}
