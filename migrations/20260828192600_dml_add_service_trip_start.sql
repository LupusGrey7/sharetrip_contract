-- +goose Up
-- +goose StatementBegin
--
-- Dictionary row for ShareTrip integration (company 123246 / trip_start).
-- Safe on DBs where 20260805161334 already applied without trip_start.
--

INSERT INTO contract_management.services (service_code, description, is_active)
VALUES (
    'trip_start',
    'старт поездки (ShareTrip moveTripPublished-ToStarted)',
    TRUE
)
ON CONFLICT (service_code) DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM contract_management.services
WHERE service_code = 'trip_start'
  AND NOT EXISTS (
      SELECT 1
      FROM contract_management.contract_services cs
      WHERE cs.service_code = 'trip_start'
  );

-- +goose StatementEnd
