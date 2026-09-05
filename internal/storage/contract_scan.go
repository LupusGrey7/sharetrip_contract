package storage

import (
	"job4j/sharetrip-contract/internal/contract/domain"

	"github.com/jackc/pgx/v5"
)

func scanContractRow(row pgx.Row) (*domain.ContractEntity, error) {
	var c domain.ContractEntity
	var status string
	err := row.Scan(
		&c.ID,
		&c.ContractNumber,
		&c.CompanyID,
		&status,
		&c.StartDate,
		&c.EndDate,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	c.Status = domain.ContractStatus(status)
	return &c, nil
}
