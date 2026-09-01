package api

import "time"

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
	ID             int       `json:"id"`
	ContractNumber string    `json:"contract_number"`
	CompanyID      int       `json:"company_id"`
	Status         string    `json:"status"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UpsertServiceItemRequest struct {
	ServiceCode string `json:"service_code" validate:"required,oneof=trip_creation trip_participants notifications premium_support"`
	IsEnabled   bool   `json:"is_enabled"`
}

type UpsertServicesRequest struct {
	ContractID int                        `json:"contract_id" validate:"required,min=1"`
	Services   []UpsertServiceItemRequest `json:"services" validate:"required,min=1,dive"`
}

type UpsertServiceItemResponse struct {
	ServiceCode string `json:"service_code"`
	IsEnabled   bool   `json:"is_enabled"`
}

type UpsertServicesResponse struct {
	ContractID int                         `json:"contract_id"`
	Services   []UpsertServiceItemResponse `json:"services"`
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
	ServiceCode string `params:"serviceCode" validate:"required,oneof=trip_creation trip_participants notifications premium_support"`
}

// GetContractByIDRequest — Fiber ParamsParser (path).
type GetContractByIDRequest struct {
	ContractID int `params:"contractId" validate:"required,min=1"`
}

// GetActiveContractByCompanyIDRequest — Fiber QueryParser (OpenAPI query companyId).
type GetActiveContractByCompanyIDRequest struct {
	CompanyID int `query:"companyId" validate:"required,min=1"`
}
