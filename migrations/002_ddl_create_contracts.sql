-- create enum contract status
CREATE TYPE contracts.contract_status AS ENUM (
    'draft',
    'active',
    'suspended',
    'terminated'
);

-- create table companies
CREATE TABLE IF NOT EXISTS contracts.companies (
    id serial NOT NULL, -- Company ID
    name varchar(255) NOT NULL, -- Company name
    created_at timestamp NOT NULL DEFAULT now(), -- Company note: created at
    updated_at timestamp NOT NULL DEFAULT now(), -- Company note: updated at
    PRIMARY KEY (id)
);

-- create table contracts
CREATE TABLE IF NOT EXISTS contracts.contracts (
    id serial NOT NULL, -- Contract ID
    contract_number text NOT NULL, -- Contract number
    company_id int NOT NULL, -- Company ID
    status contracts.contract_status NOT NULL DEFAULT 'active', -- Contract status
    start_date timestamp NOT NULL DEFAULT now(), -- Contract start date
    end_date timestamp NOT NULL DEFAULT now() + interval '1 year', -- Contract end date
    created_at timestamp NOT NULL DEFAULT now(), -- Contract created at
    updated_at timestamp NOT NULL DEFAULT now(), -- Contract updated at
    CHECK (start_date <= end_date),
    PRIMARY KEY (id),
    CONSTRAINT fk_contracts_companies FOREIGN KEY (company_id) REFERENCES contracts.companies(id)
);

-- Add indexes to a table contracts
CREATE INDEX idx_contracts_company_id ON contracts.contracts (company_id);
CREATE INDEX idx_contracts_status ON contracts.contracts (status);
CREATE UNIQUE INDEX idx_contracts_contract_number ON contracts.contracts (contract_number);

-- Business rule: Only one active contract per company
-- Index checks for uniqueness of company_id ONLY for rows with status 'active'
CREATE UNIQUE INDEX IF NOT EXISTS idx_contracts_one_active_per_company
ON contracts.contracts (company_id)
WHERE status = 'active'; -- Ключевое бизнес-правило

-- Add description to a enum contract_status
COMMENT ON TYPE contracts.contract_status IS 'Статус договора';
-- check PostgreSQL version 16.2 or higher (if lower than 16.2, comment may be throw an error)
-- COMMENT ON value contracts.contract_status.draft IS 'Draft - Договор создан, но не подписан';
-- COMMENT ON value contracts.contract_status.active IS 'Active - Договор подписан и активен';
-- COMMENT ON value contracts.contract_status.suspended IS 'Suspended - Договор приостановлен по инициативе клиента или системы';
-- COMMENT ON value contracts.contract_status.terminated IS 'Terminated - Договор прекращен по инициативе клиента или системы';

-- Add a description to a table companies
COMMENT ON TABLE contracts.companies IS 'Таблица-справочник компаний';
COMMENT ON COLUMN contracts.companies.id IS 'ID компании';
COMMENT ON COLUMN contracts.companies.name IS 'Наименование компании';
COMMENT ON COLUMN contracts.companies.created_at IS 'Дата создания';
COMMENT ON COLUMN contracts.companies.updated_at IS 'Дата обновления';

-- Add a description to a table contracts
COMMENT ON TABLE contracts.contracts IS 'Таблица договоров';
COMMENT ON COLUMN contracts.contracts.id IS 'ID договора (внутренний идентификатор)';
COMMENT ON COLUMN contracts.contracts.contract_number IS 'Номер договора';
COMMENT ON COLUMN contracts.contracts.company_id IS 'ID компании';
COMMENT ON COLUMN contracts.contracts.status IS 'Статус договора';
COMMENT ON COLUMN contracts.contracts.start_date IS 'Дата начала действия договора';
COMMENT ON COLUMN contracts.contracts.end_date IS 'Дата окончания действия договора';
COMMENT ON COLUMN contracts.contracts.created_at IS 'Дата создания договора';
COMMENT ON COLUMN contracts.contracts.updated_at IS 'Дата обновления договора';

-- drop index if exists idx_contracts_company_id;
-- drop index if exists idx_contracts_one_active_per_company;
-- drop index if exists idx_contracts_status;

-- drop type if exists contracts.contract_status;
-- drop table if exists contracts.contracts CASCADE;
-- drop table if exists contracts.companies CASCADE;
