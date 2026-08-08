package http

import "job4j/sharetrip-contract/internal/contract/model"

func toCreateContractInput(r *CreateContractRequest) *model.CreateContractRequest {
	if r == nil {
		return nil
	}
	return &model.CreateContractRequest{
		CompanyID:      r.CompanyID,
		ContractNumber: r.ContractNumber,
		Status:         model.ContractStatus(r.Status),
		StartDate:      r.StartDate,
		EndDate:        r.EndDate,
	}
}

func toContractResponse(c *model.ContractResponse) *ContractResponse {
	if c == nil {
		return nil
	}
	return &ContractResponse{
		ID:             c.ID,
		ContractNumber: c.ContractNumber,
		CompanyID:      c.CompanyID,
		Status:         string(c.Status),
		StartDate:      c.StartDate,
		EndDate:        c.EndDate,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
}

func toUpsertInput(r *UpsertServicesRequest) *model.UpsertContractServicesRequest {
	if r == nil {
		return nil
	}
	items := make([]model.ServiceItem, 0, len(r.Services))
	for _, s := range r.Services {
		items = append(items, model.ServiceItem{
			ServiceCode: s.ServiceCode,
			IsEnabled:   s.IsEnabled,
		})
	}
	return &model.UpsertContractServicesRequest{
		ContractID: r.ContractID,
		Services:   items,
	}
}

func toUpsertResponse(r *model.UpsertContractServicesResponse) *UpsertServicesResponse {
	if r == nil {
		return nil
	}
	items := make([]ServiceItemDTO, 0, len(r.Services))
	for _, s := range r.Services {
		items = append(items, ServiceItemDTO{
			ServiceCode: s.ServiceCode,
			IsEnabled:   s.IsEnabled,
		})
	}
	return &UpsertServicesResponse{
		ContractID: r.ContractID,
		Services:   items,
	}
}

func toHealthcheckResponse(r *model.GetHealthcheckInfoResponse) *HealthcheckResponse {
	if r == nil {
		return nil
	}
	return &HealthcheckResponse{Status: r.Status, Message: r.Message}
}

func toAvailabilityResponse(r *model.GetAvailableResultResponse) *AvailabilityResult {
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
