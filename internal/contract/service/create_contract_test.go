package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"job4j/sharetrip-contract/internal/contract/domain"
	"job4j/sharetrip-contract/internal/contract/service"
	"job4j/sharetrip-contract/internal/contract/usecase"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5"
)

// stubContractUseCase — заглушка BaseContractUseCase для unit-тестов сервиса.
type stubContractUseCase struct {
	output *domain.ContractOutput
	err    error

	gotInput *domain.CreateContractInput
	called   bool
}

func (s *stubContractUseCase) CreateContract(
	ctx context.Context,
	tx pgx.Tx,
	repo storage.BaseTxContractRepository,
	input *domain.CreateContractInput,
) (*domain.ContractOutput, error) {
	s.called = true
	s.gotInput = input
	if s.err != nil {
		return nil, s.err
	}
	return s.output, nil
}

func (s *stubContractUseCase) GetContractByID(
	ctx context.Context,
	tx pgx.Tx,
	repo storage.BaseTxContractRepository,
	input *domain.GetContractByIDInput,
) (*domain.ContractOutput, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.output, nil
}

func (s *stubContractUseCase) GetActiveContractByCompanyID(
	ctx context.Context,
	tx pgx.Tx,
	repo storage.BaseTxContractRepository,
	input *domain.GetActiveContractByCompanyIDInput,
) (*domain.ContractOutput, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.output, nil
}

func TestContractService_CreateContract_OK(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

	uc := &stubContractUseCase{
		output: &domain.ContractOutput{
			ID:             7,
			ContractNumber: "C-10-1",
			CompanyID:      10,
			Status:         domain.ContractStatusDraft,
			StartDate:      start,
			EndDate:        end,
			CreatedAt:      now,
			UpdatedAt:      now,
		},
	}
	svc := service.NewContractService(nil, nil, uc)

	got, err := svc.CreateContract(context.Background(), &domain.CreateContractInput{
		CompanyID:      10,
		ContractNumber: "C-10-1",
		StartDate:      start,
		EndDate:        end,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !uc.called {
		t.Fatal("expected use case to be called")
	}
	if uc.gotInput == nil || uc.gotInput.CompanyID != 10 {
		t.Fatalf("use case input: %+v", uc.gotInput)
	}
	if got == nil || got.ID != 7 || got.CompanyID != 10 || got.Status != domain.ContractStatusDraft {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestContractService_CreateContract_InvalidRequest(t *testing.T) {
	uc := &stubContractUseCase{err: usecase.ErrInvalidRequest}
	svc := service.NewContractService(nil, nil, uc)

	start := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	_, err := svc.CreateContract(context.Background(), &domain.CreateContractInput{
		CompanyID: 10,
		StartDate: start,
		EndDate:   start.AddDate(0, -1, 0),
	})
	if !errors.Is(err, usecase.ErrInvalidRequest) {
		t.Fatalf("got %v, want ErrInvalidRequest", err)
	}
}
