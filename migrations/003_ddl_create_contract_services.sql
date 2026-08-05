-- create enum service type
CREATE TYPE contracts.service_type AS ENUM (
    'trip_creation',
    'trip_participants',
    'notifications',
    'premium_support'
);

-- Create dictionary table services
CREATE TABLE IF NOT EXISTS contracts.services (
    id serial NOT NULL,
    service_code contracts.service_type NOT NULL DEFAULT 'trip_creation',
    description text,
    is_active boolean NOT NULL DEFAULT TRUE,
    PRIMARY KEY (id)
);

-- Create a table for contract services (services - these are the services that are provided as part of the contract)
CREATE TABLE IF NOT EXISTS contracts.contract_services (
    id serial NOT NULL,
    contract_id int NOT NULL,
    service_id int NOT NULL,
    is_enabled boolean NOT NULL DEFAULT TRUE, -- Флаг доступности услуг в рамках договора
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    PRIMARY KEY (id),
    CONSTRAINT fk_service_id FOREIGN KEY (service_id) REFERENCES contracts.services(id),
    CONSTRAINT fk_contract_id FOREIGN KEY (contract_id) REFERENCES contracts.contracts(id)
);

-- Index for a table
CREATE UNIQUE INDEX IF NOT EXISTS idx_services_service_code ON contracts.services (service_code);

-- index FRK
CREATE INDEX IF NOT EXISTS idx_contract_services_contract_id ON contracts.contract_services (contract_id);
CREATE INDEX IF NOT EXISTS idx_contract_services_service_id ON contracts.contract_services (service_id);

-- Business rule: Only one service per contract
-- Index checks for uniqueness of contract_id and service_id
CREATE UNIQUE INDEX IF NOT EXISTS idx_contract_services_unique
ON contracts.contract_services (contract_id, service_id);

-- Add a data to a table services
INSERT INTO contracts.services (service_code, description, is_active)
VALUES ('trip_creation', 'создание поездок', TRUE) ON CONFLICT (service_code) DO NOTHING;
INSERT INTO contracts.services (service_code, description, is_active)
VALUES ('trip_participants', 'добавление участников поездки', TRUE) ON CONFLICT (service_code) DO NOTHING;
INSERT INTO contracts.services (service_code, description, is_active)
VALUES ('notifications', 'отправка уведомлений', TRUE) ON CONFLICT (service_code) DO NOTHING;
INSERT INTO contracts.services (service_code, description, is_active)
VALUES ('premium_support', 'расширенная поддержка', TRUE) ON CONFLICT (service_code) DO NOTHING;

-- Add a description to a table
COMMENT ON TABLE contracts.contract_services IS 'Таблица-справочник услуг доступные в рамках договора';
COMMENT ON COLUMN contracts.contract_services.id IS 'ID Технический идентификатор связи договора и услуги';
COMMENT ON COLUMN contracts.contract_services.contract_id IS 'ID договора';
COMMENT ON COLUMN contracts.contract_services.service_id IS 'ID услуги';
COMMENT ON COLUMN contracts.contract_services.created_at IS 'Дата создания';
COMMENT ON COLUMN contracts.contract_services.updated_at IS 'Дата обновления';
COMMENT ON COLUMN contracts.contract_services.is_enabled IS 'Флаг доступности услуг в рамках договора';

-- Add a description to a table
COMMENT ON TABLE contracts.services IS 'Таблица-справочник услуг';
COMMENT ON COLUMN contracts.services.id IS 'ID услуги';
COMMENT ON COLUMN contracts.services.service_code IS 'Код услуги';
COMMENT ON COLUMN contracts.services.description IS 'Описание услуги';
COMMENT ON COLUMN contracts.services.is_active IS 'Флаг активности услуги, в данный момент доступна';


-- drop index if exists idx_services_name;
-- drop index if exists idx_services_service_code;
-- drop index if exists idx_contract_services_contract_id;
-- drop index if exists idx_contract_services_service_id;
-- drop index if exists idx_contracts_one_active_per_company;
-- drop index if exists idx_contract_services_unique;

-- drop table if exists contracts.services CASCADE;
-- drop table if exists contracts.contract_services;
