package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"job4j/sharetrip-contract/internal/contract/model"
	"job4j/sharetrip-contract/internal/contract/usecase"

	"github.com/go-playground/validator/v10"
)

type stubCompanyService struct {
	resp *model.GetAvailableResultResponse
	err  error
	got  *model.GetAvailableOfferingByCompanyIDRequest
}

func (s *stubCompanyService) GetAvailableOfferingByCompanyID(
	ctx context.Context,
	req *model.GetAvailableOfferingByCompanyIDRequest,
) (*model.GetAvailableResultResponse, error) {
	s.got = req
	return s.resp, s.err
}

func TestGetAvailableOfferingByCompanyID_HTTP_200_Allowed(t *testing.T) {
	company := &stubCompanyService{
		resp: &model.GetAvailableResultResponse{
			CompanyID:   42,
			ServiceCode: model.ServiceCodeTripCreation,
			Allowed:     true,
		},
	}
	srv := &Server{
		Validator:       validator.New(validator.WithRequiredStructEnabled()),
		ContractService: stubContractService{},
		OfferingService: stubOfferingService{},
		CompanyService:  company,
	}
	app := newRoutesApp(t, srv)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v2/companies/42/services/trip_creation/availability",
		nil,
	)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("status=%d body=%s", resp.StatusCode, b)
	}

	var got AvailabilityResult
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.CompanyID != 42 || got.ServiceCode != "trip_creation" || !got.Allowed {
		t.Fatalf("unexpected: %+v", got)
	}
	if company.got == nil || company.got.CompanyID != 42 || company.got.ServiceCode != model.ServiceCodeTripCreation {
		t.Fatalf("service got request: %+v", company.got)
	}
}

func TestGetAvailableOfferingByCompanyID_HTTP_200_Denied(t *testing.T) {
	srv := &Server{
		ContractService: stubContractService{},
		OfferingService: stubOfferingService{},
		CompanyService: &stubCompanyService{
			resp: &model.GetAvailableResultResponse{
				CompanyID:   42,
				ServiceCode: model.ServiceCodeTripCreation,
				Allowed:     false,
				Reason:      "no active contract with this service",
			},
		},
	}
	app := newRoutesApp(t, srv)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v2/companies/42/services/trip_creation/availability",
		nil,
	)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d", resp.StatusCode)
	}
	var got AvailabilityResult
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.Allowed || got.Reason == "" {
		t.Fatalf("want allowed=false with reason, got %+v", got)
	}
}

func TestGetAvailableOfferingByCompanyID_HTTP_404_Company(t *testing.T) {
	srv := &Server{
		ContractService: stubContractService{},
		OfferingService: stubOfferingService{},
		CompanyService:  &stubCompanyService{err: usecase.ErrCompanyNotFound},
	}
	app := newRoutesApp(t, srv)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v2/companies/99/services/trip_creation/availability",
		nil,
	)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status=%d", resp.StatusCode)
	}
	var got ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.Code != "COMPANY_NOT_FOUND" {
		t.Fatalf("code=%q, want COMPANY_NOT_FOUND", got.Code)
	}
}

func TestGetAvailableOfferingByCompanyID_HTTP_404_Service(t *testing.T) {
	srv := &Server{
		ContractService: stubContractService{},
		OfferingService: stubOfferingService{},
		CompanyService:  &stubCompanyService{err: usecase.ErrServiceNotFound},
	}
	app := newRoutesApp(t, srv)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v2/companies/42/services/trip_creation/availability",
		nil,
	)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status=%d", resp.StatusCode)
	}
	var got ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.Code != "SERVICE_NOT_FOUND" {
		t.Fatalf("code=%q, want SERVICE_NOT_FOUND", got.Code)
	}
}

func TestGetAvailableOfferingByCompanyID_HTTP_400_BadCompanyID(t *testing.T) {
	srv := &Server{
		ContractService: stubContractService{},
		OfferingService: stubOfferingService{},
		CompanyService:  &stubCompanyService{},
	}
	app := newRoutesApp(t, srv)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v2/companies/abc/services/trip_creation/availability",
		nil,
	)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("status=%d body=%s", resp.StatusCode, b)
	}
}

func TestGetAvailableOfferingByCompanyID_HTTP_400_InvalidServiceCode(t *testing.T) {
	srv := &Server{
		Validator:       validator.New(validator.WithRequiredStructEnabled()),
		ContractService: stubContractService{},
		OfferingService: stubOfferingService{},
		CompanyService:  &stubCompanyService{},
	}
	app := newRoutesApp(t, srv)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v2/companies/42/services/unknown_code/availability",
		nil,
	)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("status=%d body=%s", resp.StatusCode, b)
	}
}
