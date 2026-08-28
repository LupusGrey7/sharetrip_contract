package service

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"job4j/sharetrip-contract/internal/contract/domain"
	"job4j/sharetrip-contract/internal/observability/logctx"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

func (s *CompanyService) GetAvailableOfferingByCompanyID(
	ctx context.Context,
	input *domain.GetAvailableOfferingByCompanyIDInput,
) (res *domain.AvailabilityOutput, err error) {
	ctxSpc, span := otel.Tracer("CompanyService").Start(ctx, "CompanyService.GetAvailableOfferingByCompanyID")
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	logger := logctx.Logger(ctxSpc).With(
		slog.String("service", "CompanyService"),
		slog.String("operation", "GetAvailableOfferingByCompanyID"),
		slog.String("company_id", strconv.Itoa(input.CompanyID)),
		slog.String("service_code", string(input.ServiceCode)),
	)
	logger.Debug("get available offering started")

	if s.pool == nil {
		res, err = s.useCase.GetAvailableOfferingByCompanyID(
			ctxSpc, nil, s.offeringRepo, s.companyRepo, input,
		)
		if err != nil {
			return nil, fmt.Errorf("GetAvailableOfferingByCompanyID: %w", err)
		}
		return res, nil
	}

	txCtx, txSpan := otel.Tracer("database").Start(ctxSpc, "DB.Transaction")
	defer txSpan.End()

	res, err = tx(txCtx, s.pool, func(pgTx pgx.Tx) (*domain.AvailabilityOutput, error) {
		return s.useCase.GetAvailableOfferingByCompanyID(
			txCtx,
			pgTx,
			s.offeringRepo,
			s.companyRepo,
			input,
		)
	})
	if err != nil {
		return nil, fmt.Errorf("GetAvailableOfferingByCompanyID: %w", err)
	}

	logger.Debug("get available offering completed", slog.Bool("allowed", res.Allowed))
	return res, nil
}
