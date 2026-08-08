package http

type GetInfoResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ContractStatus string

const (
	ContractStatusDraft      ContractStatus = "draft"
	ContractStatusActive     ContractStatus = "active"
	ContractStatusSuspended  ContractStatus = "suspended"
	ContractStatusTerminated ContractStatus = "terminated"
)

type CreateContractRequest struct {
	CompanyID      int            `json:"company_id"`
	ContractNumber string         `json:"contract_number"`
	Status         ContractStatus `json:"status"`
	StartDate      string         `json:"start_date"`
	EndDate        string         `json:"end_date"`
}

type ContractResponse struct {
	ID          int            `json:"id"`
	Number      string         `json:"number"`
	Date        string         `json:"date"`
	Status      ContractStatus `json:"status"`
	TotalAmount float64        `json:"total_amount"`
}

type ToContract struct {
	CompanyID      int            `json:"company_id"`
	ContractNumber string         `json:"contract_number"`
	Status         ContractStatus `json:"status"`
	StartDate      string         `json:"start_date"`
	EndDate        string         `json:"end_date"`
}

type ContractServiceItemResponse struct {
	ContractID int `json:"contract_id"`
	ServiceID  int `json:"service_id"`
	Quantity   int `json:"quantity"`
}
