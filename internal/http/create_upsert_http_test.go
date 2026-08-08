package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"job4j/sharetrip-contract/internal/contract/model"
	"job4j/sharetrip-contract/internal/contract/usecase"

	"github.com/gofiber/fiber/v2"
)

type stubContractService struct {
	createResp *model.ContractResponse
	createErr  error
	getResp    *model.ContractResponse
	getErr     error
}

func (s stubContractService) CreateContract(ctx context.Context, request *model.CreateContractRequest) (*model.ContractResponse, error) {
	return s.createResp, s.createErr
}

func (s stubContractService) GetContractByID(ctx context.Context, request *model.GetContractByIDRequest) (*model.ContractResponse, error) {
	return s.getResp, s.getErr
}

type stubOfferingService struct {
	resp *model.UpsertContractServicesResponse
	err  error
}

func (s stubOfferingService) UpsertContractServices(ctx context.Context, request *model.UpsertContractServicesRequest) (*model.UpsertContractServicesResponse, error) {
	return s.resp, s.err
}

func newRoutesApp(t *testing.T, srv *Server) *fiber.App {
	t.Helper()
	app := fiber.New()
	srv.SetupRoutes(app)
	return app
}

func TestCreateContract_HTTP_201(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	srv := &Server{
		ContractService: stubContractService{
			createResp: &model.ContractResponse{
				ID: 5, ContractNumber: "C-1", CompanyID: 10, Status: model.ContractStatusDraft,
				StartDate: now, EndDate: now.AddDate(1, 0, 0), CreatedAt: now, UpdatedAt: now,
			},
		},
		OfferingService: stubOfferingService{},
	}
	app := newRoutesApp(t, srv)

	body := []byte(`{"company_id":10,"start_date":"2026-01-01T00:00:00Z","end_date":"2027-01-01T00:00:00Z"}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/v2/contracts/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("status=%d body=%s", resp.StatusCode, b)
	}
}

func TestUpsertContractServices_HTTP_200(t *testing.T) {
	srv := &Server{
		ContractService: stubContractService{},
		OfferingService: stubOfferingService{
			resp: &model.UpsertContractServicesResponse{
				ContractID: 5,
				Services:   []model.ServiceItem{{ServiceCode: "trip_creation", IsEnabled: true}},
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
	defer resp.Body.Close()
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
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status=%d", resp.StatusCode)
	}
}
