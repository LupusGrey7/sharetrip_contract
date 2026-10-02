package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"job4j/sharetrip-contract/internal/contract/domain"

	"github.com/gofiber/fiber/v2"
)

type stubContractService struct {
	createResp  *domain.ContractOutput
	createErr   error
	activeResp  *domain.ContractOutput
	activeErr   error
	activeInput **domain.GetActiveContractByCompanyIDInput
}

func (s stubContractService) CreateContract(ctx context.Context, input *domain.CreateContractInput) (*domain.ContractOutput, error) {
	return s.createResp, s.createErr
}

func (s stubContractService) GetActiveContractByCompanyID(
	ctx context.Context,
	input *domain.GetActiveContractByCompanyIDInput,
) (*domain.ContractOutput, error) {
	if s.activeInput != nil {
		*s.activeInput = input
	}
	return s.activeResp, s.activeErr
}

func newRoutesApp(t *testing.T, srv *Server) *fiber.App {
	t.Helper()
	app := fiber.New()
	srv.SetupRoutes(app)
	return app
}

func performRequest(t *testing.T, app interface {
	Test(*http.Request, ...int) (*http.Response, error)
}, method, path string, body []byte) *http.Response {
	t.Helper()

	req, err := http.NewRequest(method, path, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Host = "localhost"
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func decodeJSON(t *testing.T, resp *http.Response, target any) {
	t.Helper()
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		t.Fatalf("decode response: %v", err)
	}
}

func closeResponseBody(t *testing.T, resp *http.Response) {
	t.Helper()
	if err := resp.Body.Close(); err != nil {
		t.Errorf("close response body: %v", err)
	}
}
