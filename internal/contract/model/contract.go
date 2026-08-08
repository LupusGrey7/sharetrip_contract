package model

import "time"

type ContractStatus string

const (
	ContractStatusDraft      ContractStatus = "draft"
	ContractStatusActive     ContractStatus = "active"
	ContractStatusSuspended  ContractStatus = "suspended"
	ContractStatusTerminated ContractStatus = "terminated"
)

// Contract — сущность договора (колонки = contract_management.contracts).
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

// CreateContractRequest = OpenAPI ContractRequest (без списка услуг — услуги через PUT /services).
type CreateContractRequest struct {
	CompanyID      int            `json:"company_id" validate:"required,min=1"`
	ContractNumber string         `json:"contract_number" validate:"omitempty,min=1"`
	Status         ContractStatus `json:"status" validate:"omitempty,oneof=draft active suspended terminated"`
	StartDate      time.Time      `json:"start_date" validate:"required"`
	EndDate        time.Time      `json:"end_date" validate:"required"`
}

// ContractResponse = OpenAPI ContractResponse / ContractStatusResponse.
type ContractResponse struct {
	ID             int            `json:"id"`
	ContractNumber string         `json:"contract_number"`
	CompanyID      int            `json:"company_id"`
	Status         ContractStatus `json:"status"`
	StartDate      time.Time      `json:"start_date"`
	EndDate        time.Time      `json:"end_date"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type GetContractByIDRequest struct {
	ContractID int `json:"contract_id" validate:"required,min=1"`
}

type ChangeContractStatusRequest struct {
	ContractID int            `json:"contract_id" validate:"required,min=1"`
	Status     ContractStatus `json:"status" validate:"required,oneof=draft active suspended terminated"`
}

type ChangeContractStatusResponse struct {
	ID             int            `json:"id"`
	ContractNumber string         `json:"contract_number"`
	Status         ContractStatus `json:"status"`
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
