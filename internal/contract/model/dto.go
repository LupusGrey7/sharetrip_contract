package model

import "time"

type ContractStatus string

const (
	ContractStatusDraft      ContractStatus = "draft"
	ContractStatusActive     ContractStatus = "active"
	ContractStatusSuspended  ContractStatus = "suspended"
	ContractStatusTerminated ContractStatus = "terminated"
)

type InfoResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type CreateContractRequest struct {
	CompanyID      int                   `json:"company_id" validate:"required,min=1"`
	ContractNumber string                `json:"contract_number" validate:"required,min=1"`
	Status         ContractStatus        `json:"status" validate:"required,oneof=draft active suspended terminated"`
	StartDate      string                `json:"start_date" validate:"required,date"`
	EndDate        string                `json:"end_date" validate:"required,date"`
	Services       []ContractServiceItem `json:"services" validate:"required,min=1"`
}

type Contract struct {
	ID          int     `json:"id"`
	Number      string  `json:"number"`
	Date        string  `json:"date"`
	Status      string  `json:"status"`
	TotalAmount float64 `json:"total_amount"`
}

type Service struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ContractServiceItem struct {
	ContractID  int    `json:"contract_id" validate:"required,min=1"`
	ServiceCode string `json:"service_code" validate:"required,oneof=trip_creation trip_participants notifications premium_support"`
	IsEnabled   bool   `json:"is_enabled" validate:"required,boolean" default:"false"`
}

type GetContractByIDRequest struct {
	ContractID int `json:"contract_id" validate:"required,min=1"`
}

type GetContractByIDResponse struct {
	ContractID          int     `json:"contract_id"`
	ContractNumber      string  `json:"contract_number"`
	ContractDate        string  `json:"contract_date"`
	ContractStatus      string  `json:"contract_status"`
	ContractTotalAmount float64 `json:"contract_total_amount"`
}

type GetOfferingByIDRequest struct {
	ID int `json:"id" validate:"required,min=1"`
}

type GetOfferingByIDResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Offering struct {
	ServiceCode string    // enum: trip_creation, trip_participants, notifications, premium_support
	Description string    // description of the service
	IsActive    bool      // is the service active
	CreatedAt   time.Time // created at
	UpdatedAt   time.Time // updated at
}

type UpsertContractOfferingsRequest struct {
	ContractID int                   `json:"contract_id" validate:"required,min=1"`
	Services   []ContractServiceItem `json:"services" validate:"required,min=1"`
}

type UpsertContractOfferingsResponse struct {
	ContractID int                   `json:"contract_id"`
	Services   []ContractServiceItem `json:"services"`
}

type GetOfferingByCodeRequest struct {
	ServiceCode string `json:"service_code" validate:"required,oneof=trip_creation trip_participants notifications premium_support"`
}

type GetOfferingByCodeResponse struct {
	ServiceCode        string    `json:"service_code"`
	ServiceName        string    `json:"service_name"`
	ServiceDescription string    `json:"service_description"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type ChangeContractStatusRequest struct {
	ContractID int            `json:"contract_id" validate:"required,min=1"`
	Status     ContractStatus `json:"status" validate:"required,oneof=draft active suspended terminated"`
}

type ChangeContractStatusResponse struct {
	ContractID int            `json:"contract_id"`
	Status     ContractStatus `json:"status"`
}

type ContractOfferingItemResponse struct {
	ContractID  int    `json:"contract_id"`
	ServiceCode string `json:"service_code"`
	IsEnabled   bool   `json:"is_enabled"`
}

type GetHealthcheckInfoResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
