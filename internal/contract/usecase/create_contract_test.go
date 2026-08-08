package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"job4j/sharetrip-contract/internal/contract/model"
	"job4j/sharetrip-contract/internal/contract/usecase"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5"
)

type stubContractRepo struct {
	contract *model.Contract
	err      error
}

func (s stubContractRepo) GetContractByIDTx(ctx context.Context, tx pgx.Tx, id int) (*model.Contract, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.contract, nil
}

func (s stubContractRepo) GetContractByIDForUpdateTx(ctx context.Context, tx pgx.Tx, id int) (*model.Contract, error) {
	return s.GetContractByIDTx(ctx, tx, id)
}

func (s stubContractRepo) CreateContractTx(ctx context.Context, tx pgx.Tx, contract *model.Contract) (*model.Contract, error) {
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

	got, err := uc.CreateContract(context.Background(), nil, stubContractRepo{}, &model.CreateContractRequest{
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
	if got.Status != model.ContractStatusDraft {
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

	_, err := uc.CreateContract(context.Background(), nil, stubContractRepo{}, &model.CreateContractRequest{
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
		&model.GetContractByIDRequest{ContractID: 1},
	)
	if !errors.Is(err, usecase.ErrContractNotFound) {
		t.Fatalf("got %v, want ErrContractNotFound", err)
	}
}
