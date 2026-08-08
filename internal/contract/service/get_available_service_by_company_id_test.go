package service_test

import (
	"context"
	"errors"
	"testing"

	"job4j/sharetrip-contract/internal/contract/model"
	"job4j/sharetrip-contract/internal/contract/service"
	"job4j/sharetrip-contract/internal/contract/usecase"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5"
)

// stubCompanyUseCase is a stub implementation of the CompanyUseCase interface.
type stubCompanyUseCase struct {
	resp *model.GetAvailableResultResponse
	err  error

	gotCompanyID   int
	gotServiceCode model.ServiceCodeType
	called         bool
}

// GetAvailableOfferingByCompanyID is a stub implementation of the GetAvailableOfferingByCompanyID method.
func (s *stubCompanyUseCase) GetAvailableOfferingByCompanyID(
	ctx context.Context,
	tx pgx.Tx,
	offeringRepo storage.BaseTxOfferingRepository,
	companyRepo storage.BaseCompanyRepository,
	req *model.GetAvailableOfferingByCompanyIDRequest,
) (*model.GetAvailableResultResponse, error) {
	s.called = true
	if req != nil {
		s.gotCompanyID = req.CompanyID
		s.gotServiceCode = req.ServiceCode
	}
	return s.resp, s.err
}

func TestCompanyService_GetAvailableOfferingByCompanyID_OK(t *testing.T) {
	uc := &stubCompanyUseCase{
		resp: &model.GetAvailableResultResponse{
			CompanyID:   42,
			ServiceCode: model.ServiceCodeTripCreation,
			Allowed:     true,
		},
	}
	svc := service.NewCompanyService(nil, nil, nil, uc)

	got, err := svc.GetAvailableOfferingByCompanyID(context.Background(), &model.GetAvailableOfferingByCompanyIDRequest{
		CompanyID:   42,
		ServiceCode: model.ServiceCodeTripCreation,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !uc.called || uc.gotCompanyID != 42 || uc.gotServiceCode != model.ServiceCodeTripCreation {
		t.Fatalf("use case call: called=%v company=%d code=%s", uc.called, uc.gotCompanyID, uc.gotServiceCode)
	}
	if got == nil || !got.Allowed {
		t.Fatalf("got %+v", got)
	}
}

func TestCompanyService_GetAvailableOfferingByCompanyID_CompanyNotFound(t *testing.T) {
	uc := &stubCompanyUseCase{err: usecase.ErrCompanyNotFound}
	svc := service.NewCompanyService(nil, nil, nil, uc)

	_, err := svc.GetAvailableOfferingByCompanyID(context.Background(), &model.GetAvailableOfferingByCompanyIDRequest{
		CompanyID:   99,
		ServiceCode: model.ServiceCodeTripCreation,
	})
	if !errors.Is(err, usecase.ErrCompanyNotFound) {
		t.Fatalf("got %v, want ErrCompanyNotFound", err)
	}
}

func TestCompanyService_GetAvailableOfferingByCompanyID_ServiceNotFound(t *testing.T) {
	uc := &stubCompanyUseCase{err: usecase.ErrServiceNotFound}
	svc := service.NewCompanyService(nil, nil, nil, uc)

	_, err := svc.GetAvailableOfferingByCompanyID(context.Background(), &model.GetAvailableOfferingByCompanyIDRequest{
		CompanyID:   42,
		ServiceCode: model.ServiceCodeTripCreation,
	})
	if !errors.Is(err, usecase.ErrServiceNotFound) {
		t.Fatalf("got %v, want ErrServiceNotFound", err)
	}
}

func TestCompanyService_GetAvailableOfferingByCompanyID_Denied(t *testing.T) {
	uc := &stubCompanyUseCase{
		resp: &model.GetAvailableResultResponse{
			CompanyID:   42,
			ServiceCode: model.ServiceCodeTripCreation,
			Allowed:     false,
			Reason:      "service disabled for active contract",
		},
	}
	svc := service.NewCompanyService(nil, nil, nil, uc)

	got, err := svc.GetAvailableOfferingByCompanyID(context.Background(), &model.GetAvailableOfferingByCompanyIDRequest{
		CompanyID:   42,
		ServiceCode: model.ServiceCodeTripCreation,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Allowed || got.Reason == "" {
		t.Fatalf("want denied with reason, got %+v", got)
	}
}
