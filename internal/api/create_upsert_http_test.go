package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"job4j/sharetrip-contract/internal/contract/domain"
	"job4j/sharetrip-contract/internal/contract/usecase"

	"github.com/gofiber/fiber/v2"
)

type stubContractService struct {
	createResp *domain.ContractOutput
	createErr  error
	getResp    *domain.ContractOutput
	getErr     error
	activeResp *domain.ContractOutput
	activeErr  error
}

func (s stubContractService) CreateContract(ctx context.Context, input *domain.CreateContractInput) (*domain.ContractOutput, error) {
	return s.createResp, s.createErr
}

func (s stubContractService) GetContractByID(ctx context.Context, input *domain.GetContractByIDInput) (*domain.ContractOutput, error) {
	return s.getResp, s.getErr
}

func (s stubContractService) GetActiveContractByCompanyID(
	ctx context.Context,
	input *domain.GetActiveContractByCompanyIDInput,
) (*domain.ContractOutput, error) {
	return s.activeResp, s.activeErr
}

type stubOfferingService struct {
	resp *domain.UpsertContractServicesOutput
	err  error
}

func (s stubOfferingService) UpsertContractServices(ctx context.Context, input *domain.UpsertContractServicesInput) (*domain.UpsertContractServicesOutput, error) {
	return s.resp, s.err
}

func newRoutesApp(t *testing.T, srv *Server) *fiber.App {
	t.Helper()
	app := fiber.New()
	srv.SetupRoutes(app)
	return app
}

func TestUpsertContractServices_HTTP_200(t *testing.T) {
	srv := &Server{
		ContractService: stubContractService{},
		OfferingService: stubOfferingService{
			resp: &domain.UpsertContractServicesOutput{
				ContractID: 5,
				Services:   []domain.ServiceItemOutput{{ServiceCode: "trip_creation", IsEnabled: true}},
			},
		},
	}
	app := newRoutesApp(t, srv)

	body := []byte(`{"contract_id":5,"services":[{"service_code":"trip_creation","is_enabled":true}]}`)
	req, _ := http.NewRequest(http.MethodPut, "/api/v2/services", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("status=%d body=%s", resp.StatusCode, b)
	}
	var got UpsertServicesResponse
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.ContractID != 5 || len(got.Services) != 1 || got.Services[0].ServiceCode != "trip_creation" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestUpsertContractServices_HTTP_404_Contract(t *testing.T) {
	srv := &Server{
		ContractService: stubContractService{},
		OfferingService: stubOfferingService{err: usecase.ErrContractNotFound},
	}
	app := newRoutesApp(t, srv)

	body := []byte(`{"contract_id":99,"services":[{"service_code":"trip_creation","is_enabled":true}]}`)
	req, _ := http.NewRequest(http.MethodPut, "/api/v2/services", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status=%d", resp.StatusCode)
	}
}
