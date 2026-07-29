-- drop schema contracts
DROP SCHEMA IF EXISTS contracts CASCADE AUTHORIZATION postgres;

-- create schema contracts
CREATE SCHEMA IF NOT EXISTS contracts AUTHORIZATION postgres;

-- drop table companies
DROP TABLE IF EXISTS contracts.companies CASCADE AUTHORIZATION postgres;

-- drop table contracts
DROP TABLE IF EXISTS contracts.contracts CASCADE AUTHORIZATION postgres;

-- drop table contract_services
DROP TABLE IF EXISTS contracts.contract_services CASCADE AUTHORIZATION postgres;

-- drop table services
DROP TABLE IF EXISTS contracts.services CASCADE AUTHORIZATION postgres;

-- drop enum contract_status
DROP TYPE IF EXISTS contracts.contract_status CASCADE AUTHORIZATION postgres;
