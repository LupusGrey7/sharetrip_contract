-- +goose Up
-- +goose StatementBegin

CREATE SCHEMA IF NOT EXISTS contract_management;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP SCHEMA IF EXISTS contract_management RESTRICT;

-- +goose StatementEnd
