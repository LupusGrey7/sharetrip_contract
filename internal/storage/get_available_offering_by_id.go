package storage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"job4j/sharetrip-contract/internal/contract/model"

	"github.com/jackc/pgx/v5"
)

// Финальная проверка availability (после того как company известна и code есть в словаре).
// Пустой результат ≠ ErrCompanyNotFound / ErrOfferingNotFound — это бизнес-ответ allowed=false.
const getAvailableOfferingByCompanyIDQuery = `
SELECT cs.is_enabled
FROM contract_management.contracts c
INNER JOIN contract_management.contract_services cs
	ON cs.contract_id = c.id
WHERE c.company_id = $1
  AND cs.service_code = $2
  AND c.status_id = 'active'
LIMIT 1
`

func (r *CompanyRepository) GetAvailableOfferingByCompanyIDTx(
	ctx context.Context,
	tx pgx.Tx,
	companyID int,
	serviceCode model.ServiceCodeType,
) (*model.GetAvailableResultResponse, error) {
	log := slog.With(
		slog.String("company_id", strconv.Itoa(companyID)),
		slog.String("service_code", string(serviceCode)),
	)
	log.Debug("get available offering by company id started")

	var enabled bool
	err := tx.QueryRow(ctx, getAvailableOfferingByCompanyIDQuery, companyID, string(serviceCode)).Scan(&enabled)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Компания и код уже проверены выше по flow — здесь нет активной связки.
			return &model.GetAvailableResultResponse{
				CompanyID:   companyID,
				ServiceCode: serviceCode,
				Allowed:     false,
				Reason:      "no active contract with this service",
			}, nil
		}
		return nil, fmt.Errorf("GetAvailableOfferingByCompanyIDTx: %w", err)
	}

	resp := &model.GetAvailableResultResponse{
		CompanyID:   companyID,
		ServiceCode: serviceCode,
		Allowed:     enabled,
	}
	if !enabled {
		resp.Reason = "service disabled for active contract"
	}

	log.Debug("get available offering by company id success", slog.Bool("allowed", enabled))
	return resp, nil
}
