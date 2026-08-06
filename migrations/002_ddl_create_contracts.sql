-- create dictionary contract status
CREATE TABLE IF NOT EXISTS contract_management.contract_status (
    id VARCHAR(32) DEFAULT 'draft', -- 'draft', 'active', 'suspended', 'terminated'
    description TEXT NOT NULL,
    is_editable BOOLEAN NOT NULL DEFAULT FALSE, -- metadata for business logic
    created_at TIMESTAMP NOT NULL DEFAULT NOW(), -- Contract created at
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(), -- Contract updated at
    PRIMARY KEY (id)
);

-- create table contracts
CREATE TABLE IF NOT EXISTS contract_management.contracts (
    id SERIAL NOT NULL, -- Contract ID
    contract_number VARCHAR(32) NOT NULL, -- Contract number
    company_id INT NOT NULL, -- Company ID
    status_id VARCHAR(32) NOT NULL, -- ID Contract Status
    start_date TIMESTAMP NOT NULL DEFAULT NOW(), -- Contract start date
    end_date TIMESTAMP NOT NULL DEFAULT NOW() + INTERVAL '1 YEAR', -- Contract end date
    created_at TIMESTAMP NOT NULL DEFAULT NOW(), -- Contract created at
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(), -- Contract updated at
    CHECK (start_date <= end_date),
    PRIMARY KEY (id),
    CONSTRAINT fk_contracts_status_id FOREIGN KEY (status_id) REFERENCES contract_management.contract_status(id)
);

-- Add indexes to a table contracts
CREATE INDEX idx_contracts_company_id ON contract_management.contracts (company_id);
CREATE INDEX idx_contracts_status_id ON contract_management.contracts (status_id);
CREATE UNIQUE INDEX idx_contracts_contract_number ON contract_management.contracts (contract_number);

-- Business rule: Only one active contract per company
-- Index checks for uniqueness of company_id ONLY for rows with status 'active'
CREATE UNIQUE INDEX IF NOT EXISTS idx_contracts_one_active_per_company
ON contract_management.contracts (company_id)
WHERE status_id = 'active'; -- Ключевое бизнес-правило


-- Insert data to a table contract_status
INSERT INTO contract_management.contract_status (id, description, is_editable)
VALUES ('draft', 'Draft - Договор создан, но не подписан', FALSE) ON CONFLICT (id) DO NOTHING;
INSERT INTO contract_management.contract_status (id, description, is_editable)
VALUES ('active', 'Active - Договор подписан и активен', TRUE) ON CONFLICT (id) DO NOTHING;
INSERT INTO contract_management.contract_status (id, description, is_editable)
VALUES ('suspended', 'Suspended - Договор приостановлен по инициативе клиента или системы', FALSE) ON CONFLICT (id) DO NOTHING;
INSERT INTO contract_management.contract_status (id, description, is_editable)
VALUES ('terminated', 'Terminated - Договор прекращен по инициативе клиента или системы', FALSE) ON CONFLICT (id) DO NOTHING;

-- Add description to a table contract_status
COMMENT ON TABLE contract_management.contract_status IS 'Статус договора';
COMMENT ON COLUMN contract_management.contract_status.id IS 'ID статуса';
COMMENT ON COLUMN contract_management.contract_status.description IS 'Описание';
COMMENT ON COLUMN contract_management.contract_status.is_editable IS 'Можно ли редактировать';
COMMENT ON COLUMN contract_management.contract_status.created_at IS 'Дата создания статуса';
COMMENT ON COLUMN contract_management.contract_status.updated_at IS 'Дата обновления статуса';

-- Add a description to a table contracts
COMMENT ON TABLE contract_management.contracts IS 'Таблица договоров';
COMMENT ON COLUMN contract_management.contracts.id IS 'ID договора (внутренний идентификатор)';
COMMENT ON COLUMN contract_management.contracts.contract_number IS 'Номер договора';
COMMENT ON COLUMN contract_management.contracts.company_id IS 'ID компании';
COMMENT ON COLUMN contract_management.contracts.status_id IS 'ID статуса договора';
COMMENT ON COLUMN contract_management.contracts.start_date IS 'Дата начала действия договора';
COMMENT ON COLUMN contract_management.contracts.end_date IS 'Дата окончания действия договора';
COMMENT ON COLUMN contract_management.contracts.created_at IS 'Дата создания договора (timestamp)';
COMMENT ON COLUMN contract_management.contracts.updated_at IS 'Дата обновления договора (timestamp)';


-- DROP INDEX IF EXISTS contract_management.idx_contracts_company_id;
-- DROP INDEX IF EXISTS contract_management.idx_contracts_status_id;
-- DROP INDEX IF EXISTS contract_management.idx_contracts_contract_number;
-- DROP INDEX IF EXISTS contract_management.idx_contracts_one_active_per_company;
-- ALTER TABLE IF EXISTS contract_management.contracts DROP CONSTRAINT IF EXISTS fk_contracts_status_id;

-- DROP TABLE IF EXISTS contract_management.contracts;
-- DROP TABLE IF EXISTS contract_management.contract_status;
