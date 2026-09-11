package api

import (
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"job4j/sharetrip-contract/gen"
	"job4j/sharetrip-contract/internal/contract/domain"
	"job4j/sharetrip-contract/internal/contract/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// HTTP component tests for POST /api/v2/contracts/ (handler: create_contract.go).
//
// What is real here: Fiber router, BodyParser, validator, presenter, HTTP status + JSON.
// What is stubbed: ContractService — we do not hit PostgreSQL in this file.
//
// Earlier these cases lived in create_upsert_http_test.go and handler_error_paths_test.go.
// Review asked for a file named create_contract_test.go next to the handler.

func TestCreateContract_HTTP(t *testing.T) {
	now := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	contractID := uuid.MustParse("00000000-0000-0000-0000-000000000005")

	validBody := []byte(`{
		"company_id": 42,
		"client_id": "11111111-1111-1111-1111-111111111111",
		"contract_number": "C-42",
		"status": "active",
		"start_date": "2026-01-01T00:00:00Z",
		"expired_at": "2027-01-01T00:00:00Z"
	}`)

	tests := []struct {
		name       string
		body       []byte
		service    stubContractService
		wantStatus int
		wantCode   string
	}{
		{
			name: "success returns 201",
			body: validBody,
			service: stubContractService{createResp: &domain.ContractOutput{
				ID:             contractID,
				ContractNumber: "C-42",
				CompanyID:      42,
				Status:         domain.ContractStatusActive,
				StartDate:      now,
				EndDate:        now.AddDate(1, 0, 0),
				CreatedAt:      now,
				UpdatedAt:      now,
			}},
			wantStatus: http.StatusCreated,
		},
		{
			// BodyParser cannot unmarshal broken JSON → 400 before the service is called.
			name:       "malformed JSON",
			body:       []byte(`{"company_id":`),
			service:    stubContractService{},
			wantStatus: http.StatusBadRequest,
		},
		{
			// Validator tags on gen.CreateContractRequest (from x-oapi-codegen-extra-tags).
			name:       "request validation failed",
			body:       []byte(`{}`),
			service:    stubContractService{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "usecase rejected request",
			body:       validBody,
			service:    stubContractService{createErr: usecase.ErrInvalidRequest},
			wantStatus: http.StatusBadRequest,
			wantCode:   "COMPANY_VALIDATE_ERROR",
		},
		{
			// Unknown errors must not leak internals; HandleError maps them to 500.
			name:       "unexpected service error",
			body:       validBody,
			service:    stubContractService{createErr: errors.New("database unavailable")},
			wantStatus: http.StatusInternalServerError,
			wantCode:   "INTERNAL_SERVER_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Server{
				Validator:       validator.New(validator.WithRequiredStructEnabled()),
				ContractService: tt.service,
				CompanyService:  &stubCompanyService{},
			}

			resp := performRequest(t, newRoutesApp(t, srv), http.MethodPost, "/api/v2/contracts", tt.body)
			defer closeResponseBody(t, resp)

			if resp.StatusCode != tt.wantStatus {
				body, _ := io.ReadAll(resp.Body)
				t.Fatalf("status=%d, want=%d body=%s", resp.StatusCode, tt.wantStatus, body)
			}

			if tt.wantStatus == http.StatusCreated {
				var got gen.CreateContractResponse
				decodeJSON(t, resp, &got)
				c := got.Contract
				if c.Id == nil || *c.Id != openapi_types.UUID(contractID) {
					t.Fatalf("unexpected id: %+v", got)
				}
				if c.CompanyId == nil || *c.CompanyId != 42 {
					t.Fatalf("unexpected company_id: %+v", got)
				}
				if c.Status == nil || *c.Status != gen.Active {
					t.Fatalf("unexpected status: %+v", got)
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
