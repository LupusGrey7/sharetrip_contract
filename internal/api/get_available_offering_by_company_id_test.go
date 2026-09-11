package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"job4j/sharetrip-contract/internal/contract/domain"
	"job4j/sharetrip-contract/internal/contract/usecase"

	"github.com/go-playground/validator/v10"
)

type stubCompanyService struct {
	resp *domain.AvailabilityOutput
	err  error
	got  *domain.GetAvailableOfferingByCompanyIDInput
}

func (s *stubCompanyService) GetAvailableOfferingByCompanyID(
	ctx context.Context,
	input *domain.GetAvailableOfferingByCompanyIDInput,
) (*domain.AvailabilityOutput, error) {
	s.got = input
	return s.resp, s.err
}

func TestGetAvailableOfferingByCompanyID_HTTP_200_Allowed(t *testing.T) {
	company := &stubCompanyService{
		resp: &domain.AvailabilityOutput{
			CompanyID:   42,
			ServiceCode: domain.ServiceCodeTripCreation,
			Allowed:     true,
		},
	}
	srv := &Server{
		Validator:       validator.New(validator.WithRequiredStructEnabled()),
		ContractService: stubContractService{},
		CompanyService:  company,
	}
	app := newRoutesApp(t, srv)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v2/companies/42/services/trip_creation/availability",
		nil,
	)
	req.Host = "localhost"
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

	var got AvailabilityResult
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.CompanyID != 42 || got.ServiceCode != "trip_creation" || !got.Allowed {
		t.Fatalf("unexpected: %+v", got)
	}
	if company.got == nil || company.got.CompanyID != 42 || company.got.ServiceCode != domain.ServiceCodeTripCreation {
		t.Fatalf("service got input: %+v", company.got)
	}
}

func TestGetAvailableOfferingByCompanyID_HTTP_200_Denied(t *testing.T) {
	srv := &Server{
		ContractService: stubContractService{},
		CompanyService: &stubCompanyService{
			resp: &domain.AvailabilityOutput{
				CompanyID:   42,
				ServiceCode: domain.ServiceCodeTripCreation,
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
	req.Host = "localhost"
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
		CompanyService:  &stubCompanyService{err: usecase.ErrCompanyNotFound},
	}
	app := newRoutesApp(t, srv)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v2/companies/99/services/trip_creation/availability",
		nil,
	)
	req.Host = "localhost"
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
		CompanyService:  &stubCompanyService{err: usecase.ErrServiceNotFound},
	}
	app := newRoutesApp(t, srv)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v2/companies/42/services/trip_creation/availability",
		nil,
	)
	req.Host = "localhost"
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
	var got ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.Code != "SERVICE_NOT_FOUND" {
		t.Fatalf("code=%q, want SERVICE_NOT_FOUND", got.Code)
	}
}

func TestGetAvailableOfferingByCompanyID_HTTP_400_BadCompanyID(t *testing.T) {
	company := &stubCompanyService{}
	srv := &Server{
		ContractService: stubContractService{},
		CompanyService:  company,
	}
	app := newRoutesApp(t, srv)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v2/companies/abc/services/trip_creation/availability",
		nil,
	)
	req.Host = "localhost"
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	if resp.StatusCode != http.StatusBadRequest {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("status=%d body=%s", resp.StatusCode, b)
	}
	if company.got != nil {
		t.Fatalf("service must not be called, got input: %+v", company.got)
	}
}

func TestGetAvailableOfferingByCompanyID_HTTP_400_CompanyIDBelowMinimum(t *testing.T) {
	company := &stubCompanyService{}
	srv := &Server{
		Validator:       validator.New(validator.WithRequiredStructEnabled()),
		ContractService: stubContractService{},
		CompanyService:  company,
	}
	app := newRoutesApp(t, srv)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v2/companies/0/services/trip_creation/availability",
		nil,
	)
	req.Host = "localhost"
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	if resp.StatusCode != http.StatusBadRequest {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("status=%d body=%s", resp.StatusCode, b)
	}
	if company.got != nil {
		t.Fatalf("service must not be called, got input: %+v", company.got)
	}
}

func TestGetAvailableOfferingByCompanyID_HTTP_400_InvalidServiceCode(t *testing.T) {
	company := &stubCompanyService{}
	srv := &Server{
		Validator:       validator.New(validator.WithRequiredStructEnabled()),
		ContractService: stubContractService{},
		CompanyService:  company,
	}
	app := newRoutesApp(t, srv)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v2/companies/42/services/unknown_code/availability",
		nil,
	)
	req.Host = "localhost"
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	if resp.StatusCode != http.StatusBadRequest {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("status=%d body=%s", resp.StatusCode, b)
	}
	if company.got != nil {
		t.Fatalf("service must not be called, got input: %+v", company.got)
	}
}
