package api

import (
	"time"

	"github.com/google/uuid"
)

// HTTP DTO still used where gen has no Fiber params-struct with validate tags,
// or for response shapes not yet switched in handlers.
// Create / Sign / Active query params → package gen (see api/contract.yaml x-oapi-codegen-extra-tags).

type ContractResponse struct {
	ID             uuid.UUID `json:"id" db:"id"`
	ContractNumber string    `json:"contract_number"`
	CompanyID      int       `json:"company_id"`
	Status         string    `json:"status"`
	StartDate      time.Time `json:"start_date"`
	ExpiateAt      time.Time `json:"expired_at"`
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

// AvailabilityResult = OpenAPI AvailabilityResult (response; no inbound validate tags in old dto).
type AvailabilityResult struct {
	CompanyID   int    `json:"company_id"`
	ServiceCode string `json:"service_code"`
	Allowed     bool   `json:"allowed"`
	Reason      string `json:"reason,omitempty"`
}

// GetAvailableOfferingByCompanyIDRequest — path args from RegisterHandlers are separate
// method params (no gen struct). Same validate rules as YAML / former dto.
type GetAvailableOfferingByCompanyIDRequest struct {
	CompanyID   int    `validate:"required,min=1"`
	ServiceCode string `validate:"required,oneof=trip_start trip_creation trip_participants notifications premium_support"`
}
