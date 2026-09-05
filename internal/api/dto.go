package api

import (
	"github.com/google/uuid"
	"time"
)

// HTTP DTO = OpenAPI / JSON / path params on the API boundary.
// Service/usecase expose Input / Output. Storage entities never reach this package.

type CreateContractRequest struct {
	CompanyID      int       `json:"company_id" validate:"required,min=1"`
	ContractNumber string    `json:"contract_number" validate:"omitempty,min=1"`
	Status         string    `json:"status" validate:"omitempty,oneof=draft active suspended terminated"`
	StartDate      time.Time `json:"start_date" validate:"required"`
	EndDate        time.Time `json:"end_date" validate:"required"`
}

type ContractResponse struct {
	ID             uuid.UUID `json:"id" db:"id"`
	ContractNumber string    `json:"contract_number"`
	CompanyID      int       `json:"company_id"`
	Status         string    `json:"status"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type HealthcheckResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// AvailabilityResult = OpenAPI AvailabilityResult.
type AvailabilityResult struct {
	CompanyID   int    `json:"company_id"`
	ServiceCode string `json:"service_code"`
	Allowed     bool   `json:"allowed"`
	Reason      string `json:"reason,omitempty"`
}

// GetAvailableOfferingByCompanyIDRequest — Fiber ParamsParser (path).
type GetAvailableOfferingByCompanyIDRequest struct {
	CompanyID   int    `params:"companyId" validate:"required,min=1"`
	ServiceCode string `params:"serviceCode" validate:"required,oneof=trip_start trip_creation trip_participants notifications premium_support"`
}

// GetActiveContractByCompanyIDRequest — Fiber QueryParser (OpenAPI query companyId).
type GetActiveContractByCompanyIDRequest struct {
	CompanyID int `query:"companyId" validate:"required,min=1"`
}
