package usecase_test

import (
	"context"
	"errors"
	"testing"

	"job4j/sharetrip-contract/internal/contract/usecase"
)

type stubHealthRepo struct {
	err error
}

func (s stubHealthRepo) CheckHealth(ctx context.Context) error {
	return s.err
}

func TestHealthcheckUseCase_GetHealthcheckInfo_OK(t *testing.T) {
	uc := usecase.NewHealthcheckUseCase(stubHealthRepo{})
	resp, err := uc.GetHealthcheckInfo(context.Background())
	if err != nil {
		t.Fatalf("GetHealthcheckInfo: %v", err)
	}
	if resp.Status != "ok" {
		t.Fatalf("status = %q, want ok", resp.Status)
	}
}

func TestHealthcheckUseCase_GetHealthcheckInfo_DBError(t *testing.T) {
	uc := usecase.NewHealthcheckUseCase(stubHealthRepo{err: errors.New("boom")})
	_, err := uc.GetHealthcheckInfo(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}
