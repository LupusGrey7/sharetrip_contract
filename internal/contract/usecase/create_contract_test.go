package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"job4j/sharetrip-contract/internal/contract/domain"
	"job4j/sharetrip-contract/internal/contract/usecase"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5"
)

type stubContractRepo struct {
	contract *domain.ContractEntity
	err      error
}

func (s stubContractRepo) GetContractByIDTx(ctx context.Context, tx pgx.Tx, id int) (*domain.ContractEntity, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.contract, nil
}

func (s stubContractRepo) GetContractByIDForUpdateTx(ctx context.Context, tx pgx.Tx, id int) (*domain.ContractEntity, error) {
	return s.GetContractByIDTx(ctx, tx, id)
}

func (s stubContractRepo) GetActiveContractByCompanyIDTx(
	ctx context.Context,
	tx pgx.Tx,
	companyID int,
) (*domain.ContractEntity, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.contract, nil
}

func (s stubContractRepo) CreateContractTx(ctx context.Context, tx pgx.Tx, contract *domain.ContractEntity) (*domain.ContractEntity, error) {
	if s.err != nil {
		return nil, s.err
	}
	out := *contract
	out.ID = 42
	out.CreatedAt = time.Now()
	out.UpdatedAt = out.CreatedAt
	return &out, nil
}

func TestCreateContract_OK(t *testing.T) {
	uc := usecase.NewContractUseCase()
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0)

	got, err := uc.CreateContract(context.Background(), nil, stubContractRepo{}, &domain.CreateContractInput{
		CompanyID: 10,
		StartDate: start,
		EndDate:   end,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != 42 {
		t.Fatalf("id = %d, want 42", got.ID)
	}
	if got.Status != domain.ContractStatusDraft {
		t.Fatalf("status = %s, want draft", got.Status)
	}
	if got.ContractNumber == "" {
		t.Fatal("expected generated contract_number")
	}
}

func TestCreateContract_InvalidDates(t *testing.T) {
	uc := usecase.NewContractUseCase()
	start := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	_, err := uc.CreateContract(context.Background(), nil, stubContractRepo{}, &domain.CreateContractInput{
		CompanyID: 10,
		StartDate: start,
		EndDate:   end,
	})
	if !errors.Is(err, usecase.ErrInvalidRequest) {
		t.Fatalf("got %v, want ErrInvalidRequest", err)
	}
}

func TestGetContractByID_NotFound(t *testing.T) {
	uc := usecase.NewContractUseCase()
	_, err := uc.GetContractByID(
		context.Background(),
		nil,
		stubContractRepo{err: storage.ErrContractNotFound},
		&domain.GetContractByIDInput{ContractID: 1},
	)
	if !errors.Is(err, usecase.ErrContractNotFound) {
		t.Fatalf("got %v, want ErrContractNotFound", err)
	}
}

func TestGetActiveContractByCompanyID_OK(t *testing.T) {
	uc := usecase.NewContractUseCase()
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	got, err := uc.GetActiveContractByCompanyID(
		context.Background(),
		nil,
		stubContractRepo{contract: &domain.ContractEntity{
			ID:        5,
			CompanyID: 42,
			Status:    domain.ContractStatusActive,
			StartDate: now,
			EndDate:   now.AddDate(1, 0, 0),
		}},
		&domain.GetActiveContractByCompanyIDInput{CompanyID: 42},
	)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != 5 || got.CompanyID != 42 || got.Status != domain.ContractStatusActive {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestGetActiveContractByCompanyID_NotFound(t *testing.T) {
	uc := usecase.NewContractUseCase()
	_, err := uc.GetActiveContractByCompanyID(
		context.Background(),
		nil,
		stubContractRepo{err: storage.ErrContractNotFound},
		&domain.GetActiveContractByCompanyIDInput{CompanyID: 99},
	)
	if !errors.Is(err, usecase.ErrContractNotFound) {
		t.Fatalf("got %v, want ErrContractNotFound", err)
	}
}
