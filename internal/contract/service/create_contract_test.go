package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"job4j/sharetrip-contract/internal/contract/model"
	"job4j/sharetrip-contract/internal/contract/service"
	"job4j/sharetrip-contract/internal/contract/usecase"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5"
)

// stubContractUseCase — заглушка BaseContractUseCase для unit-тестов сервиса.
// Не вызывает repo: на слое service проверяем wiring + проброс ошибок, не SQL.
type stubContractUseCase struct {
	contract *model.Contract
	err      error

	gotReq *model.CreateContractRequest
	called bool
}

func (s *stubContractUseCase) CreateContract(
	ctx context.Context,
	tx pgx.Tx,
	repo storage.BaseTxContractRepository,
	request *model.CreateContractRequest,
) (*model.Contract, error) {
	s.called = true
	s.gotReq = request
	if s.err != nil {
		return nil, s.err
	}
	return s.contract, nil
}

func (s *stubContractUseCase) GetContractByID(
	ctx context.Context,
	tx pgx.Tx,
	repo storage.BaseTxContractRepository,
	request *model.GetContractByIDRequest,
) (*model.Contract, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.contract, nil
}

func TestContractService_CreateContract_OK(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

	uc := &stubContractUseCase{
		contract: &model.Contract{
			ID:             7,
			ContractNumber: "C-10-1",
			CompanyID:      10,
			Status:         model.ContractStatusDraft,
			StartDate:      start,
			EndDate:        end,
			CreatedAt:      now,
			UpdatedAt:      now,
		},
	}
	// pool=nil → без БД; repo можно nil — stub use case его не трогает
	svc := service.NewContractService(nil, nil, uc)

	got, err := svc.CreateContract(context.Background(), &model.CreateContractRequest{
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
	if uc.gotReq == nil || uc.gotReq.CompanyID != 10 {
		t.Fatalf("use case request: %+v", uc.gotReq)
	}
	if got == nil || got.ID != 7 || got.CompanyID != 10 || got.Status != model.ContractStatusDraft {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestContractService_CreateContract_InvalidRequest(t *testing.T) {
	uc := &stubContractUseCase{err: usecase.ErrInvalidRequest}
	svc := service.NewContractService(nil, nil, uc)

	start := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	_, err := svc.CreateContract(context.Background(), &model.CreateContractRequest{
		CompanyID: 10,
		StartDate: start,
		EndDate:   start.AddDate(0, -1, 0), // end before start — на реальном UC; здесь просто stub err
	})
	if !errors.Is(err, usecase.ErrInvalidRequest) {
		t.Fatalf("got %v, want ErrInvalidRequest", err)
	}
}
