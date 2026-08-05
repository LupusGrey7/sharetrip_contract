-- +goose Up
-- +goose StatementBegin

-- create enum service type
CREATE TYPE contract_management.service_type AS ENUM (
    'trip_creation',
    'trip_participants',
    'notifications',
    'premium_support'
);

-- Create dictionary table services
CREATE TABLE IF NOT EXISTS contract_management.services (
    id serial NOT NULL,
    service_code contract_management.service_type NOT NULL DEFAULT 'trip_creation',
    description text,
    is_active boolean NOT NULL DEFAULT TRUE,
    PRIMARY KEY (id)
);

-- Create a table for contract services (services - these are the services that are provided as part of the contract)
CREATE TABLE IF NOT EXISTS contract_management.contract_services (
    id serial NOT NULL,
    contract_id int NOT NULL,
    service_id int NOT NULL,
    is_enabled boolean NOT NULL DEFAULT TRUE, -- Flag of availability of services within the contract
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    PRIMARY KEY (id),
    CONSTRAINT fk_service_id FOREIGN KEY (service_id) REFERENCES contract_management.services(id),
    CONSTRAINT fk_contract_id FOREIGN KEY (contract_id) REFERENCES contract_management.contracts(id)
);

-- Index for a table
CREATE UNIQUE INDEX IF NOT EXISTS idx_services_service_code ON contract_management.services (service_code);

-- index FRK
CREATE INDEX IF NOT EXISTS idx_contract_services_contract_id
ON contract_management.contract_services (contract_id);
CREATE INDEX IF NOT EXISTS idx_contract_services_service_id
ON contract_management.contract_services (service_id);

-- Business rule: Only one service per contract
-- Index checks for uniqueness of contract_id and service_id
CREATE UNIQUE INDEX IF NOT EXISTS idx_contract_services_unique
ON contract_management.contract_services (contract_id, service_id);

-- Add a data to a table services
INSERT INTO contract_management.services (service_code, description, is_active)
VALUES ('trip_creation', 'создание поездок', TRUE) ON CONFLICT (service_code) DO NOTHING;
INSERT INTO contract_management.services (service_code, description, is_active)
VALUES ('trip_participants', 'добавление участников поездки', TRUE) ON CONFLICT (service_code) DO NOTHING;
INSERT INTO contract_management.services (service_code, description, is_active)
VALUES ('notifications', 'отправка уведомлений', TRUE) ON CONFLICT (service_code) DO NOTHING;
INSERT INTO contract_management.services (service_code, description, is_active)
VALUES ('premium_support', 'extended support', TRUE) ON CONFLICT (service_code) DO NOTHING;

-- Add a description to a table
COMMENT ON TABLE contract_management.contract_services IS 'Таблица-справочник услуг доступные в рамках договора';
COMMENT ON COLUMN contract_management.contract_services.id IS 'ID Технический идентификатор связи договора и услуги';
COMMENT ON COLUMN contract_management.contract_services.contract_id IS 'ID договора';
COMMENT ON COLUMN contract_management.contract_services.service_id IS 'ID услуги';
COMMENT ON COLUMN contract_management.contract_services.created_at IS 'Дата создания';
COMMENT ON COLUMN contract_management.contract_services.updated_at IS 'Дата обновления';
COMMENT ON COLUMN contract_management.contract_services.is_enabled IS 'Флаг доступности услуг в рамках договора';

-- Add a description to a table
COMMENT ON TABLE contract_management.services IS 'Таблица-справочник услуг';
COMMENT ON COLUMN contract_management.services.id IS 'ID услуги';
COMMENT ON COLUMN contract_management.services.service_code IS 'Код услуги';
COMMENT ON COLUMN contract_management.services.description IS 'Описание услуги';
COMMENT ON COLUMN contract_management.services.is_active IS 'Флаг активности услуги, в данный момент доступна';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS contract_management.idx_services_service_code;
DROP INDEX IF EXISTS contract_management.idx_contract_services_contract_id;
DROP INDEX IF EXISTS contract_management.idx_contract_services_service_id;
DROP INDEX IF EXISTS contract_management.idx_contract_services_unique;

ALTER TABLE IF EXISTS contract_management.contract_services DROP CONSTRAINT fk_service_id;
ALTER TABLE IF EXISTS contract_management.contract_services DROP CONSTRAINT fk_contract_id;
DROP TABLE IF EXISTS contract_management.contract_services;

DROP TABLE IF EXISTS contract_management.services;
DROP TYPE IF EXISTS contract_management.service_type;

-- +goose StatementEnd
