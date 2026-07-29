-- create enum contract status
CREATE TYPE contract_status AS ENUM (
    'draft',
    'active',
    'suspended',
    'terminated'
);

-- create table companies
CREATE TABLE IF NOT EXISTS companies (
    id serial NOT NULL, -- Company ID
    name varchar(255) NOT NULL, -- Company name
    created_at timestamp NOT NULL DEFAULT now(), -- Company note: created at
    updated_at timestamp NOT NULL DEFAULT now(), -- Company note: updated at
    PRIMARY KEY (id)
);

-- create table contracts
CREATE TABLE IF NOT EXISTS contracts (
    id serial NOT NULL, -- Contract ID
    company_id int NOT NULL, -- Company ID
    status contract_status NOT NULL DEFAULT 'active', -- Contract status
    start_date timestamp NOT NULL DEFAULT now(), -- Contract start date
    end_date timestamp NOT NULL DEFAULT now() + interval '1 year', -- Contract end date
    -- services jsonb NOT NULL,  -- List of available services
    services_available boolean NOT NULL DEFAULT TRUE,  -- Flag of availability of services
    created_at timestamp NOT NULL DEFAULT now(), -- Contract created at
    updated_at timestamp NOT NULL DEFAULT now(), -- Contract updated at
    CHECK (status IN ('draft', 'active', 'suspended', 'terminated')), -- Check if status is valid
    CHECK (start_date <= end_date),
    CHECK (services_available IS TRUE OR services_available IS FALSE),
    PRIMARY KEY (id),
    CONSTRAINT fk_contracts_companies FOREIGN KEY (company_id) REFERENCES companies(id)
);

-- Add description to a enum contract_status
COMMENT ON enum contract_status IS 'Статус договора';
COMMENT ON value contract_status.draft IS 'Draft - Договор создан, но не подписан';
COMMENT ON value contract_status.active IS 'Active - Договор подписан и активен';
COMMENT ON value contract_status.suspended IS 'Suspended - Договор приостановлен по инициативе клиента или системы';
COMMENT ON value contract_status.terminated IS 'Terminated - Договор прекращен по инициативе клиента или системы';

-- Add a description to a table companies
COMMENT ON TABLE companies IS 'Таблица-справочник компаний';
COMMENT ON COLUMN companies.id IS 'ID компании';
COMMENT ON COLUMN companies.name IS 'Наименование компании';
COMMENT ON COLUMN companies.created_at IS 'Дата создания';
COMMENT ON COLUMN companies.updated_at IS 'Дата обновления';

-- Add a description to a table contracts
COMMENT ON TABLE contracts IS 'Таблица договоров';
COMMENT ON COLUMN contracts.id IS 'ID договора';
COMMENT ON COLUMN contracts.company_id IS 'ID компании';
COMMENT ON COLUMN contracts.status IS 'Статус договора';
COMMENT ON COLUMN contracts.start_date IS 'Дата начала действия договора';
COMMENT ON COLUMN contracts.end_date IS 'Дата окончания действия договора';
COMMENT ON COLUMN contracts.services IS 'Список доступных услуг';
COMMENT ON COLUMN contracts.services_available IS 'Флаг доступности услуг';
COMMENT ON COLUMN contracts.created_at IS 'Дата создания договора';
COMMENT ON COLUMN contracts.updated_at IS 'Дата обновления договора';

-- Add an index to a column contracts
CREATE INDEX idx_contracts_company_id ON contracts (company_id);
CREATE INDEX idx_contracts_status ON contracts (status);

-- Business rule: Only one active contract per company
-- Index checks for uniqueness of company_id ONLY for rows with status 'active'
-- Индекс проверяет уникальность company_id ТОЛЬКО для строк со статусом 'active'
CREATE UNIQUE INDEX IF NOT EXISTS idx_contracts_one_active_per_company
ON contracts (company_id)
WHERE status = 'active';

-- drop index idx_contracts_company_id;
-- drop index idx_contracts_one_active_per_company;
-- drop table contracts;
-- drop table companies;
-- drop enum contract_status;
