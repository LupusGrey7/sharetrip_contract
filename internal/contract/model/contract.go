package model

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

// Contract — entity of contract (columns = contract_management.contracts). Without JSON — this is not HTTP DTO.
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

// CreateContractRequest — enter for usecase/service (not wire JSON; JSON lives in internal/http/dto.go).
type CreateContractRequest struct {
	CompanyID      int
	ContractNumber string
	Status         ContractStatus
	StartDate      time.Time
	EndDate        time.Time
}

// ContractResponse — result of scenario for sending outside through http mapper.
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
	// params — name of segment in URL (:contractId) for Fiber ParamsParser.
	// validate — rules go-playground/validator after parsing.
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
