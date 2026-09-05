-- +goose Up
-- +goose StatementBegin

-- create dictionary service type
CREATE TABLE IF NOT EXISTS contract_management.services (
    service_code VARCHAR(32) NOT NULL, -- 'trip_creation', 'trip_participants', 'notifications', 'premium_support'
    description TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE, -- metadata for business logic
    created_at TIMESTAMP NOT NULL DEFAULT NOW(), -- Service created at
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(), -- Service updated at
    PRIMARY KEY (service_code)
);

-- Create a table for contract services (services - these are the services that are provided as part of the contract)
CREATE TABLE IF NOT EXISTS contract_management.contract_services (
    id SERIAL NOT NULL,
    contract_id UUID NOT NULL,
    service_code VARCHAR(32) NOT NULL, -- 'trip_creation', 'trip_participants', 'notifications', 'premium_support'
    is_enabled BOOLEAN NOT NULL DEFAULT TRUE, -- Flag of availability of services within the contract
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    CONSTRAINT fk_contract_services_service_code
        FOREIGN KEY (service_code) REFERENCES contract_management.services(service_code),
    CONSTRAINT fk_contract_services_contract_id
        FOREIGN KEY (contract_id) REFERENCES contract_management.contracts(id)
);

-- index FKs
CREATE INDEX IF NOT EXISTS idx_contract_services_contract_id ON contract_management.contract_services (contract_id);
CREATE INDEX IF NOT EXISTS idx_contract_services_service_code ON contract_management.contract_services (service_code);

-- Business rule: a restriction that the same service is not duplicated within the framework of one contract
-- Index checks for uniqueness of contract_id and service_code
CREATE UNIQUE INDEX IF NOT EXISTS idx_contract_services_unique ON contract_management.contract_services (
    contract_id, service_code
);

-- Add a data to a table services
INSERT INTO contract_management.services (service_code, description, is_active)
VALUES ('trip_creation', 'создание поездок', TRUE) ON CONFLICT (service_code) DO NOTHING;
INSERT INTO contract_management.services (service_code, description, is_active)
VALUES ('trip_start', 'старт поездки (ShareTrip moveTripPublished-ToStarted)', TRUE) ON CONFLICT (service_code) DO NOTHING;
INSERT INTO contract_management.services (service_code, description, is_active)
VALUES ('trip_participants', 'добавление участников поездки', TRUE) ON CONFLICT (service_code) DO NOTHING;
INSERT INTO contract_management.services (service_code, description, is_active)
VALUES ('notifications', 'отправка уведомлений', TRUE) ON CONFLICT (service_code) DO NOTHING;
INSERT INTO contract_management.services (service_code, description, is_active)
VALUES ('premium_support', 'extended support', TRUE) ON CONFLICT (service_code) DO NOTHING;

-- Add a description to a table
COMMENT ON TABLE contract_management.contract_services IS 'Таблица связи договора и услуг';
COMMENT ON COLUMN contract_management.contract_services.id IS 'ID Технический идентификатор связи договора и услуги';
COMMENT ON COLUMN contract_management.contract_services.contract_id IS 'ID договора';
COMMENT ON COLUMN contract_management.contract_services.service_code IS 'Код услуги';
COMMENT ON COLUMN contract_management.contract_services.created_at IS 'Дата создания';
COMMENT ON COLUMN contract_management.contract_services.updated_at IS 'Дата обновления';
COMMENT ON COLUMN contract_management.contract_services.is_enabled IS 'Флаг доступности услуг в рамках договора';

-- Add a description to a table
COMMENT ON TABLE contract_management.services IS 'Таблица-справочник услуг';
COMMENT ON COLUMN contract_management.services.service_code IS 'Код услуги';
COMMENT ON COLUMN contract_management.services.description IS 'Описание услуги';
COMMENT ON COLUMN contract_management.services.is_active IS 'Флаг активности услуги, в данный момент доступна';
COMMENT ON COLUMN contract_management.services.created_at IS 'Дата создания услуги (timestamp)';
COMMENT ON COLUMN contract_management.services.updated_at IS 'Дата обновления услуги (timestamp)';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS contract_management.idx_contract_services_service_code;
DROP INDEX IF EXISTS contract_management.idx_contract_services_contract_id;
DROP INDEX IF EXISTS contract_management.idx_contract_services_unique;

ALTER TABLE IF EXISTS contract_management.contract_services DROP CONSTRAINT fk_contract_services_service_code;
ALTER TABLE IF EXISTS contract_management.contract_services DROP CONSTRAINT fk_contract_services_contract_id;

DROP TABLE IF EXISTS contract_management.contract_services;
DROP TABLE IF EXISTS contract_management.services;

-- +goose StatementEnd
