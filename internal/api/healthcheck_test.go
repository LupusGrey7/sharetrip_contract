package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"job4j/sharetrip-contract/internal/contract/domain"
	"job4j/sharetrip-contract/internal/contract/service"

	"github.com/gofiber/fiber/v2"
)

type stubInfoUseCase struct {
	resp *domain.HealthcheckOutput
	err  error
}

func (s stubInfoUseCase) GetHealthcheckInfo(ctx context.Context) (*domain.HealthcheckOutput, error) {
	return s.resp, s.err
}

func newTestApp(t *testing.T, uc stubInfoUseCase) *fiber.App {
	t.Helper()
	healthcheckSvc := service.NewHealthcheckService(uc)
	srv := NewServer(nil, healthcheckSvc, nil, nil)
	fiberApp := fiber.New()
	srv.SetupRoutes(fiberApp)
	return fiberApp
}

func TestHealthcheck_OK(t *testing.T) {
	fiberApp := newTestApp(t, stubInfoUseCase{
		resp: &domain.HealthcheckOutput{Status: "ok", Message: "application and database are healthy"},
	})

	req, err := http.NewRequest(http.MethodGet, "/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Host = "localhost"
	resp, err := fiberApp.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var got HealthcheckResponse
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("json: %v body=%s", err, body)
	}
	if got.Status != "ok" {
		t.Fatalf("status = %q, want ok", got.Status)
	}
}

func TestHealthcheck_DBDown_Returns503(t *testing.T) {
	fiberApp := newTestApp(t, stubInfoUseCase{
		err: errors.New("database healthcheck failed: connection refused"),
	})

	req, err := http.NewRequest(http.MethodGet, "/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Host = "localhost"
	resp, err := fiberApp.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", resp.StatusCode)
	}
}
