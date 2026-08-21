package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"job4j/sharetrip-contract/internal/contract/domain"
	"job4j/sharetrip-contract/internal/contract/usecase"

	"github.com/go-playground/validator/v10"
)

func TestGetContractByID_HTTP(t *testing.T) {
	now := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		path       string
		service    stubContractService
		wantStatus int
		wantCode   string
	}{
		{
			name: "success",
			path: "/api/v2/contracts/5",
			service: stubContractService{getResp: &domain.ContractOutput{
				ID:             5,
				ContractNumber: "C-5",
				CompanyID:      42,
				Status:         domain.ContractStatusActive,
				StartDate:      now,
				EndDate:        now.AddDate(1, 0, 0),
				CreatedAt:      now,
				UpdatedAt:      now,
			}},
			wantStatus: http.StatusOK,
		},
		{
			name:       "path parameter is not a number",
			path:       "/api/v2/contracts/not-a-number",
			service:    stubContractService{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "path parameter fails validation",
			path:       "/api/v2/contracts/0",
			service:    stubContractService{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "contract not found",
			path:       "/api/v2/contracts/99",
			service:    stubContractService{getErr: usecase.ErrContractNotFound},
			wantStatus: http.StatusNotFound,
			wantCode:   "CONTRACT_NOT_FOUND",
		},
		{
			name:       "unexpected service error",
			path:       "/api/v2/contracts/5",
			service:    stubContractService{getErr: errors.New("database unavailable")},
			wantStatus: http.StatusInternalServerError,
			wantCode:   "INTERNAL_SERVER_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Server{
				Validator:       validator.New(validator.WithRequiredStructEnabled()),
				ContractService: tt.service,
				OfferingService: stubOfferingService{},
				CompanyService:  &stubCompanyService{},
			}

			resp := performRequest(t, newRoutesApp(t, srv), http.MethodGet, tt.path, nil)
			defer closeResponseBody(t, resp)

			if resp.StatusCode != tt.wantStatus {
				body, _ := io.ReadAll(resp.Body)
				t.Fatalf("status=%d, want=%d body=%s", resp.StatusCode, tt.wantStatus, body)
			}

			if tt.wantStatus == http.StatusOK {
				var got ContractResponse
				decodeJSON(t, resp, &got)
				if got.ID != 5 || got.CompanyID != 42 || got.Status != string(domain.ContractStatusActive) {
					t.Fatalf("unexpected response: %+v", got)
				}
				return
			}

			if tt.wantCode != "" {
				var got ErrorResponse
				decodeJSON(t, resp, &got)
				if got.Code != tt.wantCode {
					t.Fatalf("code=%q, want=%q", got.Code, tt.wantCode)
				}
			}
		})
	}
}

func TestCreateContract_HTTP_ErrorPaths(t *testing.T) {
	validBody := []byte(`{
		"company_id": 42,
		"contract_number": "C-42",
		"status": "active",
		"start_date": "2026-01-01T00:00:00Z",
		"end_date": "2027-01-01T00:00:00Z"
	}`)

	tests := []struct {
		name       string
		body       []byte
		serviceErr error
		wantStatus int
		wantCode   string
	}{
		{
			name:       "malformed JSON",
			body:       []byte(`{"company_id":`),
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "request validation failed",
			body:       []byte(`{}`),
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "usecase rejected request",
			body:       validBody,
			serviceErr: usecase.ErrInvalidRequest,
			wantStatus: http.StatusBadRequest,
			wantCode:   "COMPANY_VALIDATE_ERROR",
		},
		{
			name:       "unexpected service error",
			body:       validBody,
			serviceErr: errors.New("database unavailable"),
			wantStatus: http.StatusInternalServerError,
			wantCode:   "INTERNAL_SERVER_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Server{
				Validator:       validator.New(validator.WithRequiredStructEnabled()),
				ContractService: stubContractService{createErr: tt.serviceErr},
				OfferingService: stubOfferingService{},
				CompanyService:  &stubCompanyService{},
			}

			resp := performRequest(t, newRoutesApp(t, srv), http.MethodPost, "/api/v2/contracts/", tt.body)
			defer closeResponseBody(t, resp)

			if resp.StatusCode != tt.wantStatus {
				body, _ := io.ReadAll(resp.Body)
				t.Fatalf("status=%d, want=%d body=%s", resp.StatusCode, tt.wantStatus, body)
			}
			if tt.wantCode != "" {
				var got ErrorResponse
				decodeJSON(t, resp, &got)
				if got.Code != tt.wantCode {
					t.Fatalf("code=%q, want=%q", got.Code, tt.wantCode)
				}
			}
		})
	}
}

func TestUpsertContractServices_HTTP_ErrorPaths(t *testing.T) {
	validBody := []byte(`{
		"contract_id": 5,
		"services": [{"service_code": "trip_creation", "is_enabled": true}]
	}`)

	tests := []struct {
		name       string
		body       []byte
		serviceErr error
		wantStatus int
		wantCode   string
	}{
		{
			name:       "malformed JSON",
			body:       []byte(`{"contract_id":`),
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "request validation failed",
			body:       []byte(`{"contract_id": 0, "services": []}`),
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "service code not found",
			body:       validBody,
			serviceErr: usecase.ErrServiceNotFound,
			wantStatus: http.StatusNotFound,
			wantCode:   "SERVICE_NOT_FOUND",
		},
		{
			name:       "unexpected service error",
			body:       validBody,
			serviceErr: errors.New("database unavailable"),
			wantStatus: http.StatusInternalServerError,
			wantCode:   "INTERNAL_SERVER_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Server{
				Validator:       validator.New(validator.WithRequiredStructEnabled()),
				ContractService: stubContractService{},
				OfferingService: stubOfferingService{err: tt.serviceErr},
				CompanyService:  &stubCompanyService{},
			}

			resp := performRequest(t, newRoutesApp(t, srv), http.MethodPut, "/api/v2/services", tt.body)
			defer closeResponseBody(t, resp)

			if resp.StatusCode != tt.wantStatus {
				body, _ := io.ReadAll(resp.Body)
				t.Fatalf("status=%d, want=%d body=%s", resp.StatusCode, tt.wantStatus, body)
			}
			if tt.wantCode != "" {
				var got ErrorResponse
				decodeJSON(t, resp, &got)
				if got.Code != tt.wantCode {
					t.Fatalf("code=%q, want=%q", got.Code, tt.wantCode)
				}
			}
		})
	}
}

func performRequest(t *testing.T, app interface {
	Test(*http.Request, ...int) (*http.Response, error)
}, method, path string, body []byte) *http.Response {
	t.Helper()

	req, err := http.NewRequest(method, path, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
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
