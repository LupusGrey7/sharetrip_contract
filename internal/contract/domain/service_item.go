package domain

import "time"

// ServiceItemEntity — contract↔service entity used by usecase and storage.
type ServiceItemEntity struct {
	ServiceCode string
	IsEnabled   bool
}

// ServiceItemInput — one item in a service/usecase input.
type ServiceItemInput struct {
	ServiceCode string
	IsEnabled   bool
}

// ServiceItemOutput — one item returned above the usecase boundary.
type ServiceItemOutput struct {
	ServiceCode string
	IsEnabled   bool
}

// UpsertContractServicesInput — input for service/usecase.
type UpsertContractServicesInput struct {
	ContractID int
	Services   []ServiceItemInput
}

// UpsertContractServicesOutput — output from usecase/service.
type UpsertContractServicesOutput struct {
	ContractID int
	Services   []ServiceItemOutput
}

// OfferingEntity — dictionary entity used by usecase and storage.
type OfferingEntity struct {
	ServiceCode string
	Description string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func ServiceItemEntitiesToOutput(entities []ServiceItemEntity) []ServiceItemOutput {
	output := make([]ServiceItemOutput, 0, len(entities))
	for _, entity := range entities {
		output = append(output, ServiceItemOutput{
			ServiceCode: entity.ServiceCode,
			IsEnabled:   entity.IsEnabled,
		})
	}
	return output
}

func ServiceItemInputToEntity(input ServiceItemInput) ServiceItemEntity {
	return ServiceItemEntity{
		ServiceCode: input.ServiceCode,
		IsEnabled:   input.IsEnabled,
	}
}
