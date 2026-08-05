-- +goose Up
-- +goose StatementBegin

-- create enum contract status
CREATE TYPE contract_management.contract_status AS ENUM (
    'draft',
    'active',
    'suspended',
    'terminated'
);

-- create table companies
CREATE TABLE IF NOT EXISTS contract_management.companies (
    id serial NOT NULL, -- Company ID
    name varchar(255) NOT NULL, -- Company name
    created_at timestamp NOT NULL DEFAULT now(), -- Company note: created at
    updated_at timestamp NOT NULL DEFAULT now(), -- Company note: updated at
    PRIMARY KEY (id)
);

-- create table contracts
CREATE TABLE IF NOT EXISTS contract_management.contracts (
    id serial NOT NULL, -- Contract ID
    contract_number text NOT NULL, -- Contract number
    company_id int NOT NULL, -- Company ID
    status contract_management.contract_status NOT NULL DEFAULT 'draft', -- Contract status
    start_date timestamp NOT NULL DEFAULT now(), -- Contract start date
    end_date timestamp NOT NULL DEFAULT now() + interval '1 year', -- Contract end date
    created_at timestamp NOT NULL DEFAULT now(), -- Contract created at
    updated_at timestamp NOT NULL DEFAULT now(), -- Contract updated at
    CHECK (start_date <= end_date),
    PRIMARY KEY (id),
    CONSTRAINT fk_contracts_companies FOREIGN KEY (company_id) REFERENCES contract_management.companies(id)
);

-- Add indexes to a table contracts
CREATE INDEX idx_contracts_company_id ON contract_management.contracts (company_id);
CREATE INDEX idx_contracts_status ON contract_management.contracts (status);
CREATE UNIQUE INDEX idx_contracts_contract_number ON contract_management.contracts (contract_number);

-- Business rule: Only one active contract per company
-- Index checks for uniqueness of company_id ONLY for rows with status 'active'
CREATE UNIQUE INDEX IF NOT EXISTS idx_contracts_one_active_per_company
ON contract_management.contracts (company_id)
WHERE status = 'active'; -- Ключевое бизнес-правило

-- Add description to a enum contract_status
COMMENT ON TYPE contract_management.contract_status IS 'Статус договора';
-- check PostgreSQL version 16.2 or higher (if lower than 16.2, comment may be throw an error)
-- COMMENT ON value contract_management.contract_status.draft IS 'Draft - Договор создан, но не подписан';
-- COMMENT ON value contract_management.contract_status.active IS 'Active - Договор подписан и активен';
-- COMMENT ON value contract_management.contract_status.suspended IS 'Suspended - Договор приостановлен по инициативе клиента или системы';
-- COMMENT ON value contract_management.contract_status.terminated IS 'Terminated - Договор прекращен по инициативе клиента или системы';

-- Add a description to a table companies
COMMENT ON TABLE contract_management.companies IS 'Таблица-справочник компаний';
COMMENT ON COLUMN contract_management.companies.id IS 'ID компании';
COMMENT ON COLUMN contract_management.companies.name IS 'Наименование компании';
COMMENT ON COLUMN contract_management.companies.created_at IS 'Дата создания';
COMMENT ON COLUMN contract_management.companies.updated_at IS 'Дата обновления';

-- Add a description to a table contracts
COMMENT ON TABLE contract_management.contracts IS 'Таблица договоров';
COMMENT ON COLUMN contract_management.contracts.id IS 'ID договора (внутренний идентификатор)';
COMMENT ON COLUMN contract_management.contracts.contract_number IS 'Номер договора';
COMMENT ON COLUMN contract_management.contracts.company_id IS 'ID компании';
COMMENT ON COLUMN contract_management.contracts.status IS 'Статус договора';
COMMENT ON COLUMN contract_management.contracts.start_date IS 'Дата начала действия договора';
COMMENT ON COLUMN contract_management.contracts.end_date IS 'Дата окончания действия договора';
COMMENT ON COLUMN contract_management.contracts.created_at IS 'Дата создания договора';
COMMENT ON COLUMN contract_management.contracts.updated_at IS 'Дата обновления договора';

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS contract_management.idx_contracts_company_id;
DROP INDEX IF EXISTS contract_management.idx_contracts_status;
DROP INDEX IF EXISTS contract_management.idx_contracts_contract_number;
DROP INDEX IF EXISTS contract_management.idx_contracts_one_active_per_company;
ALTER TABLE IF EXISTS contract_management.contracts DROP CONSTRAINT IF EXISTS fk_contracts_companies;

DROP TABLE IF EXISTS contract_management.contracts;
DROP TABLE IF EXISTS contract_management.companies;
DROP TYPE IF EXISTS contract_management.contract_status;

-- +goose StatementEnd
