-- drop enum service_type
DROP TYPE IF EXISTS contracts.service_type CASCADE;

-- drop enum contract_status
DROP TYPE IF EXISTS contracts.contract_status CASCADE;

-- drop table companies
DROP TABLE IF EXISTS contracts.companies CASCADE;

-- drop table contracts
DROP TABLE IF EXISTS contracts.contracts CASCADE;

-- drop table contract_services
DROP TABLE IF EXISTS contracts.contract_services CASCADE;

-- drop table services
DROP TABLE IF EXISTS contracts.services CASCADE;

-- drop schema contracts
DROP SCHEMA IF EXISTS contracts CASCADE;

-- create schema contracts
CREATE SCHEMA IF NOT EXISTS contracts;
