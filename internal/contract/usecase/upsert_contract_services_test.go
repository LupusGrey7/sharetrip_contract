package usecase_test

import (
	"context"
	"errors"
	"testing"

	"job4j/sharetrip-contract/internal/contract/domain"
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

func (s stubOfferingDictRepo) GetOfferingByCodeTx(ctx context.Context, tx pgx.Tx, code string) (*domain.OfferingEntity, error) {
	return nil, storage.ErrOfferingNotFound
}

func (s stubOfferingDictRepo) IsOfferingExistsByCodeTx(ctx context.Context, tx pgx.Tx, code string) error {
	if s.err != nil {
		return s.err
	}
	if len(s.missing) > 0 {
		return storage.ErrOfferingNotFound
	}
	return nil
}

type stubLinkRepo struct {
	items []domain.ServiceItemEntity
	err   error
}

func (s *stubLinkRepo) UpsertContractServiceTx(ctx context.Context, tx pgx.Tx, contractID int, item domain.ServiceItemEntity) error {
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

func (s *stubLinkRepo) ListContractServicesTx(ctx context.Context, tx pgx.Tx, contractID int) ([]domain.ServiceItemEntity, error) {
	return s.items, s.err
}

func TestUpsertContractServices_OK(t *testing.T) {
	uc := usecase.NewOfferingUseCase()
	link := &stubLinkRepo{}
	contract := stubContractRepo{contract: &domain.ContractEntity{ID: 7}}

	resp, err := uc.UpsertContractServices(
		context.Background(),
		nil,
		contract,
		stubOfferingDictRepo{},
		link,
		&domain.UpsertContractServicesInput{
			ContractID: 7,
			Services: []domain.ServiceItemInput{
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
		&domain.UpsertContractServicesInput{
			ContractID: 1,
			Services:   []domain.ServiceItemInput{{ServiceCode: "trip_creation", IsEnabled: true}},
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
		stubContractRepo{contract: &domain.ContractEntity{ID: 1}},
		stubOfferingDictRepo{missing: []string{"trip_creation"}},
		&stubLinkRepo{},
		&domain.UpsertContractServicesInput{
			ContractID: 1,
			Services:   []domain.ServiceItemInput{{ServiceCode: "trip_creation", IsEnabled: true}},
		},
	)
	if !errors.Is(err, usecase.ErrServiceNotFound) {
		t.Fatalf("got %v, want ErrServiceNotFound", err)
	}
}
