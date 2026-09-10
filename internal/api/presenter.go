package api

import (
	"job4j/sharetrip-contract/gen"
	"job4j/sharetrip-contract/internal/contract/domain"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

func toCreateContractInput(r *gen.CreateContractRequest) *domain.CreateContractInput {
	if r == nil {
		return nil
	}
	status := domain.ContractStatus("")
	if r.Status != nil {
		status = domain.ContractStatus(*r.Status)
	}
	contractNumber := ""
	if r.ContractNumber != nil {
		contractNumber = *r.ContractNumber
	}
	return &domain.CreateContractInput{
		CompanyID:      int(r.CompanyId),
		ContractNumber: contractNumber,
		Status:         status,
		StartDate:      r.StartDate,
		EndDate:        r.ExpiredAt,
	}
}

func toCreateContractResponse(output *domain.ContractOutput) *gen.CreateContractResponse {
	if output == nil {
		return nil
	}
	id := openapi_types.UUID(output.ID)
	companyID := int64(output.CompanyID)
	contractNumber := output.ContractNumber
	status := gen.ContractStatus(output.Status)
	start := output.StartDate
	expired := output.EndDate
	created := output.CreatedAt
	updated := output.UpdatedAt
	return &gen.CreateContractResponse{
		Contract: gen.ContractResponse{
			Id:             &id,
			CompanyId:      &companyID,
			ContractNumber: &contractNumber,
			Status:         &status,
			StartDate:      &start,
			ExpiredAt:      &expired,
			CreatedAt:      &created,
			UpdatedAt:      &updated,
		},
	}
}

func toGetActiveContractByCompanyIDInput(r *gen.GetActiveContractByCompanyIdParams) *domain.GetActiveContractByCompanyIDInput {
	if r == nil {
		return nil
	}
	return &domain.GetActiveContractByCompanyIDInput{CompanyID: int(r.CompanyId)}
}

func toContractResponse(output *domain.ContractOutput) *ContractResponse {
	if output == nil {
		return nil
	}
	return &ContractResponse{
		ID:             output.ID,
		ContractNumber: output.ContractNumber,
		CompanyID:      output.CompanyID,
		Status:         string(output.Status),
		StartDate:      output.StartDate,
		ExpiateAt:      output.EndDate,
		CreatedAt:      output.CreatedAt,
		UpdatedAt:      output.UpdatedAt,
	}
}

func toHealthcheckResponse(r *domain.HealthcheckOutput) *HealthcheckResponse {
	if r == nil {
		return nil
	}
	return &HealthcheckResponse{Status: r.Status, Message: r.Message}
}

func toAvailabilityInput(r *GetAvailableOfferingByCompanyIDRequest) *domain.GetAvailableOfferingByCompanyIDInput {
	if r == nil {
		return nil
	}
	return &domain.GetAvailableOfferingByCompanyIDInput{
		CompanyID:   r.CompanyID,
		ServiceCode: domain.ServiceCode(r.ServiceCode),
	}
}

func toAvailabilityResponse(r *domain.AvailabilityOutput) *AvailabilityResult {
	if r == nil {
		return nil
	}
	return &AvailabilityResult{
		CompanyID:   r.CompanyID,
		ServiceCode: string(r.ServiceCode),
		Allowed:     r.Allowed,
		Reason:      r.Reason,
	}
}
