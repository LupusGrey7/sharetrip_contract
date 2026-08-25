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
)

// HTTP component tests for GET /api/v2/contracts/{contractId}
// (OpenAPI getContractByID, contract.yaml ~95). Handler: get_contract_by_id.go.
//
// This is NOT getActiveContractByCompanyId (yaml ~177). That lives in
// get_active_contract_by_company_id_test.go.

func TestGetContractByID_HTTP(t *testing.T) {
	now := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)

	activeContract := &domain.ContractOutput{
		ID:             5,
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
		wantDraft  bool
	}{
		{
			name:       "success returns active contract by id",
			path:       "/api/v2/contracts/5",
			service:    stubContractService{getResp: activeContract},
			wantStatus: http.StatusOK,
		},
		{
			// Get by id is not "get active". A draft contract is a valid 200.
			name: "success returns draft contract by id",
			path: "/api/v2/contracts/7",
			service: stubContractService{getResp: &domain.ContractOutput{
				ID:             7,
				ContractNumber: "C-7",
				CompanyID:      42,
				Status:         domain.ContractStatusDraft,
				StartDate:      now,
				EndDate:        now.AddDate(1, 0, 0),
				CreatedAt:      now,
				UpdatedAt:      now,
			}},
			wantStatus: http.StatusOK,
			wantDraft:  true,
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
				if tt.wantDraft {
					if got.ID != 7 || got.Status != string(domain.ContractStatusDraft) {
						t.Fatalf("unexpected draft response: %+v", got)
					}
					return
				}
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
