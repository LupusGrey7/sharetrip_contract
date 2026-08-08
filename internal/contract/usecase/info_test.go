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

func TestInfoUseCase_GetHealthcheckInfo_OK(t *testing.T) {
	uc := usecase.NewInfoUseCase(stubHealthRepo{})
	resp, err := uc.GetHealthcheckInfo(context.Background())
	if err != nil {
		t.Fatalf("GetHealthcheckInfo: %v", err)
	}
	if resp.Status != "ok" {
		t.Fatalf("status = %q, want ok", resp.Status)
	}
}

func TestInfoUseCase_GetHealthcheckInfo_DBError(t *testing.T) {
	uc := usecase.NewInfoUseCase(stubHealthRepo{err: errors.New("boom")})
	_, err := uc.GetHealthcheckInfo(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}
