package model

import "time"

// ServiceItem = OpenAPI ServiceItem (связь договора с услугой, без contract_id в item).
type ServiceItem struct {
	ServiceCode string `json:"service_code" validate:"required,oneof=trip_creation trip_participants notifications premium_support"`
	IsEnabled   bool   `json:"is_enabled" validate:"required,boolean"`
}

// UpsertContractServicesRequest = OpenAPI UpsertServicesRequest.
type UpsertContractServicesRequest struct {
	ContractID int           `json:"contract_id" validate:"required,min=1"`
	Services   []ServiceItem `json:"services" validate:"required,min=1,dive"`
}

// UpsertContractServicesResponse = OpenAPI UpsertServicesResponse.
type UpsertContractServicesResponse struct {
	ContractID int           `json:"contract_id" validate:"required,min=1"`
	Services   []ServiceItem `json:"services" validate:"required,min=1,dive"`
}

// Offering — словарь services (не путать с contract_services).
type Offering struct {
	ServiceCode string
	Description string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
