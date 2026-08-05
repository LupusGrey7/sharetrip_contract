-- +goose Up
-- +goose StatementBegin
-- 1. Change the column type to text and simultaneously nullify its data

CREATE SCHEMA IF NOT EXISTS contract_management;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Create schema contracts

DROP SCHEMA IF EXISTS contract_management;

-- +goose StatementEnd
