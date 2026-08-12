package model

import "time"

type ContractStatus string

const (
	ContractStatusDraft      ContractStatus = "draft"
	ContractStatusActive     ContractStatus = "active"
	ContractStatusSuspended  ContractStatus = "suspended"
	ContractStatusTerminated ContractStatus = "terminated"
)

// Contract — сущность договора (колонки = contract_management.contracts). Без JSON — это не HTTP DTO.
type Contract struct {
	ID             int
	ContractNumber string
	CompanyID      int
	Status         ContractStatus
	StartDate      time.Time
	EndDate        time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// CreateContractRequest — вход usecase/service (не wire JSON; JSON живёт в internal/http/dto.go).
type CreateContractRequest struct {
	CompanyID      int
	ContractNumber string
	Status         ContractStatus
	StartDate      time.Time
	EndDate        time.Time
}

// ContractResponse — результат сценария для отдачи наружу через http mapper.
type ContractResponse struct {
	ID             int
	ContractNumber string
	CompanyID      int
	Status         ContractStatus
	StartDate      time.Time
	EndDate        time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type GetContractByIDRequest struct {
	// tagm params — field for Fiber ParamsParser (name in URL :contractId).
	// tag validate — field for Fiber Validator (name in URL :contractId).
	ContractID int `params:"contractId" validate:"required,min=1"`
}

type ChangeContractStatusRequest struct {
	ContractID int
	Status     ContractStatus
}

type ChangeContractStatusResponse struct {
	ID             int
	ContractNumber string
	Status         ContractStatus
}

func ContractToResponse(c *Contract) *ContractResponse {
	if c == nil {
		return nil
	}
	return &ContractResponse{
		ID:             c.ID,
		ContractNumber: c.ContractNumber,
		CompanyID:      c.CompanyID,
		Status:         c.Status,
		StartDate:      c.StartDate,
		EndDate:        c.EndDate,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
}
