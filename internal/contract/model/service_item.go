package model

import "time"

// ServiceItem — contract↔service (domain/application). JSON-form — http.ServiceItemDTO.
type ServiceItem struct {
	ServiceCode string
	IsEnabled   bool
}

type UpsertContractServicesRequest struct {
	ContractID int
	Services   []ServiceItem
}

type UpsertContractServicesResponse struct {
	ContractID int
	Services   []ServiceItem
}

// Offering — dictionary of services.
type Offering struct {
	ServiceCode string
	Description string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
