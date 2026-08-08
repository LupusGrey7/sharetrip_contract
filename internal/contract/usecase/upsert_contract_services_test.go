package usecase_test

import (
	"context"
	"errors"
	"testing"

	"job4j/sharetrip-contract/internal/contract/model"
	"job4j/sharetrip-contract/internal/contract/usecase"
	"job4j/sharetrip-contract/internal/storage"

	"github.com/jackc/pgx/v5"
)

type stubOfferingDictRepo struct {
	missing []string
	err     error
}

func (s stubOfferingDictRepo) ExistServiceCodesTx(ctx context.Context, tx pgx.Tx, codes []string) ([]string, error) {
	return s.missing, s.err
}

func (s stubOfferingDictRepo) GetOfferingByCodeTx(ctx context.Context, tx pgx.Tx, code string) (*model.Offering, error) {
	return nil, storage.ErrOfferingNotFound
}

type stubLinkRepo struct {
	items []model.ServiceItem
	err   error
}

func (s *stubLinkRepo) UpsertContractServiceTx(ctx context.Context, tx pgx.Tx, contractID int, item model.ServiceItem) error {
	if s.err != nil {
		return s.err
	}
	for i := range s.items {
		if s.items[i].ServiceCode == item.ServiceCode {
			s.items[i] = item
			return nil
		}
	}
	s.items = append(s.items, item)
	return nil
}

func (s *stubLinkRepo) ListContractServicesTx(ctx context.Context, tx pgx.Tx, contractID int) ([]model.ServiceItem, error) {
	return s.items, s.err
}

func TestUpsertContractServices_OK(t *testing.T) {
	uc := usecase.NewOfferingUseCase()
	link := &stubLinkRepo{}
	contract := stubContractRepo{contract: &model.Contract{ID: 7}}

	resp, err := uc.UpsertContractServices(
		context.Background(),
		nil,
		contract,
		stubOfferingDictRepo{},
		link,
		&model.UpsertContractServicesRequest{
			ContractID: 7,
			Services: []model.ServiceItem{
				{ServiceCode: "trip_creation", IsEnabled: true},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ContractID != 7 || len(resp.Services) != 1 {
		t.Fatalf("unexpected resp: %+v", resp)
	}
	if !resp.Services[0].IsEnabled {
		t.Fatal("expected is_enabled=true")
	}
}

func TestUpsertContractServices_ContractNotFound(t *testing.T) {
	uc := usecase.NewOfferingUseCase()
	_, err := uc.UpsertContractServices(
		context.Background(),
		nil,
		stubContractRepo{err: storage.ErrContractNotFound},
		stubOfferingDictRepo{},
		&stubLinkRepo{},
		&model.UpsertContractServicesRequest{
			ContractID: 1,
			Services:   []model.ServiceItem{{ServiceCode: "trip_creation", IsEnabled: true}},
		},
	)
	if !errors.Is(err, usecase.ErrContractNotFound) {
		t.Fatalf("got %v, want ErrContractNotFound", err)
	}
}

func TestUpsertContractServices_ServiceNotFound(t *testing.T) {
	uc := usecase.NewOfferingUseCase()
	_, err := uc.UpsertContractServices(
		context.Background(),
		nil,
		stubContractRepo{contract: &model.Contract{ID: 1}},
		stubOfferingDictRepo{missing: []string{"trip_creation"}},
		&stubLinkRepo{},
		&model.UpsertContractServicesRequest{
			ContractID: 1,
			Services:   []model.ServiceItem{{ServiceCode: "trip_creation", IsEnabled: true}},
		},
	)
	if !errors.Is(err, usecase.ErrServiceNotFound) {
		t.Fatalf("got %v, want ErrServiceNotFound", err)
	}
}
