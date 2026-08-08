package model

import "time"

// ServiceItem — связь договор↔услуга (домен/application). JSON-форма — http.ServiceItemDTO.
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

// Offering — словарь services.
type Offering struct {
	ServiceCode string
	Description string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
