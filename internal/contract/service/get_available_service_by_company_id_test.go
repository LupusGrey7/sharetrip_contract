package service_test

import (
	"context"
	"errors"
	"testing"

	"job4j/sharetrip-contract/internal/contract/domain"
	"job4j/sharetrip-contract/internal/contract/service"
	"job4j/sharetrip-contract/internal/contract/usecase"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5"
)

type stubCompanyUseCase struct {
	resp *domain.AvailabilityOutput
	err  error

	gotCompanyID   int
	gotServiceCode domain.ServiceCode
	called         bool
}

func (s *stubCompanyUseCase) GetAvailableOfferingByCompanyID(
	ctx context.Context,
	tx pgx.Tx,
	offeringRepo storage.BaseTxOfferingRepository,
	companyRepo storage.BaseCompanyRepository,
	input *domain.GetAvailableOfferingByCompanyIDInput,
) (*domain.AvailabilityOutput, error) {
	s.called = true
	if input != nil {
		s.gotCompanyID = input.CompanyID
		s.gotServiceCode = input.ServiceCode
	}
	return s.resp, s.err
}

func TestCompanyService_GetAvailableOfferingByCompanyID_OK(t *testing.T) {
	uc := &stubCompanyUseCase{
		resp: &domain.AvailabilityOutput{
			CompanyID:   42,
			ServiceCode: domain.ServiceCodeTripCreation,
			Allowed:     true,
		},
	}
	svc := service.NewCompanyService(nil, nil, nil, uc)

	got, err := svc.GetAvailableOfferingByCompanyID(context.Background(), &domain.GetAvailableOfferingByCompanyIDInput{
		CompanyID:   42,
		ServiceCode: domain.ServiceCodeTripCreation,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !uc.called || uc.gotCompanyID != 42 || uc.gotServiceCode != domain.ServiceCodeTripCreation {
		t.Fatalf("use case call: called=%v company=%d code=%s", uc.called, uc.gotCompanyID, uc.gotServiceCode)
	}
	if got == nil || !got.Allowed {
		t.Fatalf("got %+v", got)
	}
}

func TestCompanyService_GetAvailableOfferingByCompanyID_CompanyNotFound(t *testing.T) {
	uc := &stubCompanyUseCase{err: usecase.ErrCompanyNotFound}
	svc := service.NewCompanyService(nil, nil, nil, uc)

	_, err := svc.GetAvailableOfferingByCompanyID(context.Background(), &domain.GetAvailableOfferingByCompanyIDInput{
		CompanyID:   99,
		ServiceCode: domain.ServiceCodeTripCreation,
	})
	if !errors.Is(err, usecase.ErrCompanyNotFound) {
		t.Fatalf("got %v, want ErrCompanyNotFound", err)
	}
}

func TestCompanyService_GetAvailableOfferingByCompanyID_ServiceNotFound(t *testing.T) {
	uc := &stubCompanyUseCase{err: usecase.ErrServiceNotFound}
	svc := service.NewCompanyService(nil, nil, nil, uc)

	_, err := svc.GetAvailableOfferingByCompanyID(context.Background(), &domain.GetAvailableOfferingByCompanyIDInput{
		CompanyID:   42,
		ServiceCode: domain.ServiceCodeTripCreation,
	})
	if !errors.Is(err, usecase.ErrServiceNotFound) {
		t.Fatalf("got %v, want ErrServiceNotFound", err)
	}
}

func TestCompanyService_GetAvailableOfferingByCompanyID_Denied(t *testing.T) {
	uc := &stubCompanyUseCase{
		resp: &domain.AvailabilityOutput{
			CompanyID:   42,
			ServiceCode: domain.ServiceCodeTripCreation,
			Allowed:     false,
			Reason:      "service disabled for active contract",
		},
	}
	svc := service.NewCompanyService(nil, nil, nil, uc)

	got, err := svc.GetAvailableOfferingByCompanyID(context.Background(), &domain.GetAvailableOfferingByCompanyIDInput{
		CompanyID:   42,
		ServiceCode: domain.ServiceCodeTripCreation,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Allowed || got.Reason == "" {
		t.Fatalf("want denied with reason, got %+v", got)
	}
}
