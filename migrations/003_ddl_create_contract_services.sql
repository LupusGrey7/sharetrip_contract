-- Create dictionary table services
CREATE TABLE IF NOT EXISTS services (
    id serial NOT NULL,
    name varchar(255) NOT NULL,
    description text,
    is_active boolean NOT NULL DEFAULT TRUE,
    PRIMARY KEY (id)
);

-- Create a table for contract services (services - these are the services that are provided as part of the contract)
CREATE TABLE IF NOT EXISTS contract_services (
    id serial NOT NULL,
    contract_id int NOT NULL,
    service_id int NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    PRIMARY KEY (contract_id, service_id),
    CONSTRAINT fk_service_id FOREIGN KEY (service_id) REFERENCES services(id),
    CONSTRAINT fk_contract_id FOREIGN KEY (contract_id) REFERENCES contracts(id)
);

-- Add a data to a table services
INSERT INTO services (name, description, is_active) VALUES ('trip_creation', 'создание поездок', TRUE) ON CONFLICT DO NOTHING;
INSERT INTO services (name, description, is_active) VALUES ('trip_participants', 'добавление участников поездки', TRUE) ON CONFLICT DO NOTHING;
INSERT INTO services (name, description, is_active) VALUES ('notifications', 'отправка уведомлений', TRUE) ON CONFLICT DO NOTHING;
INSERT INTO services (name, description, is_active) VALUES ('premium_support', 'расширенная поддержка', TRUE) ON CONFLICT DO NOTHING;

-- Add a description to a table
COMMENT ON TABLE contract_services IS 'Таблица-справочник услуг доступные в рамках договора';
COMMENT ON COLUMN contract_services.id IS 'ID услуги в рамках договора';
COMMENT ON COLUMN contract_services.contract_id IS 'ID договора';
COMMENT ON COLUMN contract_services.service_id IS 'ID услуги';
COMMENT ON COLUMN contract_services.created_at IS 'Дата создания';
COMMENT ON COLUMN contract_services.updated_at IS 'Дата обновления';

-- Add a description to a table
COMMENT ON TABLE services IS 'Таблица-справочник услуг';
COMMENT ON COLUMN services.id IS 'ID услуги';
COMMENT ON COLUMN services.name IS 'Наименование услуги';
COMMENT ON COLUMN services.description IS 'Описание услуги';
COMMENT ON COLUMN services.is_active IS 'Флаг активности услуги, в данный момент доступна';

-- Index for a table
CREATE INDEX IF NOT EXISTS idx_services_name ON services (name);

-- index FRK
CREATE INDEX IF NOT EXISTS idx_contract_services_contract_id ON contract_services (contract_id);
CREATE INDEX IF NOT EXISTS idx_contract_services_service_id ON contract_services (service_id);

-- Business rule: Only one active contract per company
-- Index checks for uniqueness of company_id ONLY for rows with status 'active'
CREATE UNIQUE INDEX IF NOT EXISTS idx_contracts_one_active_per_company
ON contracts (company_id)
WHERE status = 'active'; -- Ключевое бизнес-правило

-- Business rule: Only one service per contract
-- Index checks for uniqueness of contract_id and service_id
CREATE UNIQUE INDEX IF NOT EXISTS idx_contract_services_unique
ON contract_services (contract_id, service_id);

-- drop index idx_services_name;
-- drop index idx_contract_services_contract_id;
-- drop index idx_contract_services_service_id;
-- drop index idx_contracts_one_active_per_company;
-- drop index idx_contract_services_unique;
-- drop table services;
-- drop table contract_services;
