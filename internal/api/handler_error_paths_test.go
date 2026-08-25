package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"job4j/sharetrip-contract/internal/contract/usecase"

	"github.com/go-playground/validator/v10"
)

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
