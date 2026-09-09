package api

import "job4j/sharetrip-contract/internal/contract/domain"

func toCreateContractInput(r *CreateContractRequest) *domain.CreateContractInput {
	if r == nil {
		return nil
	}
	return &domain.CreateContractInput{
		CompanyID:      r.CompanyID,
		ContractNumber: r.ContractNumber,
		Status:         domain.ContractStatus(r.Status),
		StartDate:      r.StartDate,
		EndDate:        r.EndDate,
	}
}

func toGetActiveContractByCompanyIDInput(r *GetActiveContractByCompanyIDRequest) *domain.GetActiveContractByCompanyIDInput {
	if r == nil {
		return nil
	}
	return &domain.GetActiveContractByCompanyIDInput{CompanyID: r.CompanyID}
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
