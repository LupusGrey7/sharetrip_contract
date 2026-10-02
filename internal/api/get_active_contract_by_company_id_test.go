package api

import (
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"job4j/sharetrip-contract/internal/contract/domain"
	"job4j/sharetrip-contract/internal/contract/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// HTTP component tests for GET /api/v2/contracts/active?companyId=.
// RegisterHandlers' generated wrapper binds the query into gen params before
// calling the handler; the handler validates and maps those params to domain input.
//
// Review asked for get_active_contract_by_id_test.go. That name mixed two YAML
// operations. This file is get-active-by-company (yaml ~177).

func TestGetActiveContractByCompanyID_HTTP(t *testing.T) {
	now := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	contractID := uuid.MustParse("00000000-0000-0000-0000-000000000005")

	active := &domain.ContractOutput{
		ID:             contractID,
		ContractNumber: "C-5",
		CompanyID:      42,
		Status:         domain.ContractStatusActive,
		StartDate:      now,
		EndDate:        now.AddDate(1, 0, 0),
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	tests := []struct {
		name       string
		path       string
		service    stubContractService
		wantStatus int
		wantCode   string
		wantInput  int
	}{
		{
			name:       "success returns active contract for company",
			path:       "/api/v2/contracts/active?companyId=42",
			service:    stubContractService{activeResp: active},
			wantStatus: http.StatusOK,
			wantInput:  42,
		},
		{
			name:       "missing companyId query",
			path:       "/api/v2/contracts/active",
			service:    stubContractService{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "companyId fails validation",
			path:       "/api/v2/contracts/active?companyId=0",
			service:    stubContractService{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "companyId is not a number",
			path:       "/api/v2/contracts/active?companyId=not-a-number",
			service:    stubContractService{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "active contract not found",
			path:       "/api/v2/contracts/active?companyId=99",
			service:    stubContractService{activeErr: usecase.ErrContractNotFound},
			wantStatus: http.StatusNotFound,
			wantCode:   "CONTRACT_NOT_FOUND",
			wantInput:  99,
		},
		{
			name:       "unexpected service error",
			path:       "/api/v2/contracts/active?companyId=42",
			service:    stubContractService{activeErr: errors.New("database unavailable")},
			wantStatus: http.StatusInternalServerError,
			wantCode:   "INTERNAL_SERVER_ERROR",
			wantInput:  42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotInput *domain.GetActiveContractByCompanyIDInput
			service := tt.service
			service.activeInput = &gotInput
			srv := &Server{
				Validator:       validator.New(validator.WithRequiredStructEnabled()),
				ContractService: service,
				CompanyService:  &stubCompanyService{},
			}

			resp := performRequest(t, newRoutesApp(t, srv), http.MethodGet, tt.path, nil)
			defer closeResponseBody(t, resp)

			if resp.StatusCode != tt.wantStatus {
				body, _ := io.ReadAll(resp.Body)
				t.Fatalf("status=%d, want=%d body=%s", resp.StatusCode, tt.wantStatus, body)
			}
			if tt.wantInput == 0 {
				if gotInput != nil {
					t.Fatalf("service called unexpectedly with input: %+v", gotInput)
				}
			} else if gotInput == nil || gotInput.CompanyID != tt.wantInput {
				t.Fatalf("service input=%+v, want company_id=%d", gotInput, tt.wantInput)
			}

			if tt.wantStatus == http.StatusOK {
				var got ContractResponse
				decodeJSON(t, resp, &got)
				if got.ID != contractID || got.CompanyID != 42 || got.Status != string(domain.ContractStatusActive) {
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
