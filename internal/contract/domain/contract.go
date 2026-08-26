package domain

import "time"

// ContractStatus — status of contract.
type ContractStatus string

// ContractStatus values.
const (
	ContractStatusDraft      ContractStatus = "draft"
	ContractStatusActive     ContractStatus = "active"
	ContractStatusSuspended  ContractStatus = "suspended"
	ContractStatusTerminated ContractStatus = "terminated"
)

// ContractEntity — entity used by usecase and storage. It never reaches API.
type ContractEntity struct {
	ID             int
	ContractNumber string
	CompanyID      int
	Status         ContractStatus
	StartDate      time.Time
	EndDate        time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// ContractOutput — output from usecase/service to API presenter.
type ContractOutput struct {
	ID             int
	ContractNumber string
	CompanyID      int
	Status         ContractStatus
	StartDate      time.Time
	EndDate        time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// CreateContractInput — input for service/usecase (not wire JSON).
type CreateContractInput struct {
	CompanyID      int
	ContractNumber string
	Status         ContractStatus
	StartDate      time.Time
	EndDate        time.Time
}

// GetActiveContractByCompanyIDInput — OpenAPI getActiveContractByCompanyId (query companyId).
type GetActiveContractByCompanyIDInput struct {
	CompanyID int
}

// ChangeContractStatusInput — input for service/usecase.
type ChangeContractStatusInput struct {
	ContractID int
	Status     ContractStatus
}

// ChangeContractStatusOutput — output from service/usecase.
type ChangeContractStatusOutput struct {
	ID             int
	ContractNumber string
	Status         ContractStatus
}

// ContractEntityToOutput closes the entity boundary inside the usecase layer.
func ContractEntityToOutput(entity *ContractEntity) *ContractOutput {
	if entity == nil {
		return nil
	}
	return &ContractOutput{
		ID:             entity.ID,
		ContractNumber: entity.ContractNumber,
		CompanyID:      entity.CompanyID,
		Status:         entity.Status,
		StartDate:      entity.StartDate,
		EndDate:        entity.EndDate,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
	}
}
