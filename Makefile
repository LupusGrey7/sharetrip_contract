# ============================================================
# Env profile: APP_ENV → automatically ENV_FILE=.env.$(APP_ENV)
# ============================================================
# Persist choice in .active-app-env (local, gitignored).
#   make env-default              → back to dev / .env.dev
#   make env-use APP_ENV=prod     → prod / .env.prod (remembered)
# One-shot (does not change saved profile):
#   make migrate-status APP_ENV=test
ACTIVE_PROFILE_FILE := .active-app-env

ifeq ($(origin APP_ENV),command line)
    # one-shot or env-use: APP_ENV came from the command line
else ifneq ($(wildcard $(ACTIVE_PROFILE_FILE)),)
    APP_ENV := $(strip $(file < $(ACTIVE_PROFILE_FILE)))
    ifeq ($(APP_ENV),)
        APP_ENV := dev
    endif
else
    APP_ENV := dev
endif

# Always derived from profile — do not pick ENV_FILE by hand
ENV_FILE := .env.$(APP_ENV)

ifneq (,$(wildcard $(ENV_FILE)))
    include $(ENV_FILE)
    export
else
    $(warning Env file '$(ENV_FILE)' not found. Create it or: make env-use APP_ENV=dev|prod|test)
endif

export APP_ENV
export ENV_FILE

# ============================================================
# Check Operation System PC — единственная логика что остаётся в Makefile
# ============================================================
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    YAML_CHECK_SCRIPT := scripts\checkos\yaml-check.bat
    RM_RF = if exist "$(BUILD_DIR)" rmdir /s /q "$(BUILD_DIR)"
else
    DETECTED_OS := $(shell uname -s)
    YAML_CHECK_SCRIPT := bash scripts/checkos/yaml-check.sh
    RM_RF = rm -rf
endif

# ============================================================
# Important Variables
# ============================================================
GO := go
GO_PKG := ./...
APP_NAME=sharetrip_contract
BUILD_DIR=./build
MAIN_FILE=cmd/contract/main.go
HTTP_PORT ?= 8082
DB_DSN=$(DB_DRIVER)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
MIGRATIONS_DIR = ./migrations
SEEDS_DIR := ./scripts/seeds
GOOSE_SEED_TABLE := goose_seed_version
DEPLOY_DIR := ./deploy
DC := $(DEPLOY_DIR)/docker/docker-compose.yml

# Version and ldflags (for embedding the version into the binary)
VERSION=1.0.0
LDFLAGS=-ldflags "-X main.Version=${VERSION}"

# Task - Default target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  deps        	  	 					- install tools or check their availability"
	@echo "  fmt         	  	 					- code formating"
	@echo "  lint        	  	 					- run the linter"
	@echo "  test        	  		 				- run all tests"
	@echo "  test-integration							- run API integration tests with Testcontainers"
	@echo "  build         	 					- build a binary file"
	@echo "  run         	 	 					- run the application locally"
	@echo "  e2e         	  	 					- end to end check an application locally"
	@echo "  up          	 	 					- raise app infrastructure docker image"
	@echo "  start         	 					    - start app infrastructure docker image"
	@echo "  stop        	  	 					- stop app infrastructure docker image"
	@echo "  restart        	 					- restart app infrastructure docker image"
	@echo "  clean-image    	 					- clean all app infrastructure docker image"
	@echo "  down        	  	 					- down app infrastructure docker image"
	@echo "  migrate-up  	  	 					- apply migrations"
	@echo "  migrate-down	  	 					- roll back the last migration"
	@echo "  migrate-status  	 					- check migration status"
	@echo "  seed            	 					- apply dev seed data (demo company 42 + trip_creation)"
	@echo "  seed-down       	 					- roll back last seed"
	@echo "  seed-status     	 					- check seed version status"
	@echo "  check           	 					- run all checks: formatting, linter, tests, coverage, vulnerability detection"
	@echo "  coverage    	  	 					- run tests and generate HTML coverage report"
	@echo "  cover       	  						- alias for coverage"
	@echo "  vulncheck       	 					- run vulnerability detection tool"
	@echo "  all         	 						- run all checks: lint, tests, coverage, vulnerability detection"
	@echo "  yaml-check     	 					- run yaml check tool"
	@echo "  info                       				 	- show information about the OS and yq"
	@echo "  help                        					- show this help"
	@echo "  env-default                 					- reset profile to default (APP_ENV=dev → .env.dev)"
	@echo "  env-use        						- switch profile; file = .env.<profile> (APP_ENV={prod|dev|test})"
	@echo "  env-info                    					- show active APP_ENV / ENV_FILE / DB target (password masked)"

# Task - Prepare environment (installation of tools)
# If you don't want to install locally (for example, again), comment out the commands inside
# ============================================================
# Prepare Environment (installation of tools)
# ============================================================
.PHONY: deps
deps:
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5
	$(GO) install github.com/pressly/goose/v3/cmd/goose@latest

# Task - Formats source code
.PHONY: fmt
fmt:
	$(GO) fmt $(GO_PKG)

# Task - Check code with golangci-lint (start with check OS),
.PHONY: lint
lint:
ifeq ($(OS),Windows_NT)
	golangci-lint run
else
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "❌ golangci-lint is not installed. Please install it:"; \
		echo "   https://golangci-lint.run"; \
		exit 1; \
	fi
	golangci-lint run
endif

# Task - Run tests (detailed output (test names and PASS/FAIL))
.PHONY: test
test:
	$(GO) test -v $(GO_PKG)

.PHONY: test-integration
test-integration:
	$(GO) test -count=1 -v ./internal/api/apitest

# Task - Clean builds (Deletes compiled files)
.PHONY: clean
clean:
	@echo "Cleaning..."
ifeq ($(OS),Windows_NT)
	$(RM_RF)
else
	$(RM_RF) $(BUILD_DIR)
endif

# Task - Compile Go code to a binary file
.PHONY: build
build:
	@echo "Building..."
	$(GO) build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)

# Task - Local application run
.PHONY: run
run:
	$(GO) run $(MAIN_FILE)

# Task - Local application check (including adding a check for the response body:)
.PHONY: e2e
e2e:
ifeq ($(OS),Windows_NT)
	powershell -NoProfile -Command "$$r = Invoke-WebRequest -Uri http://localhost:$(HTTP_PORT)/healthcheck -UseBasicParsing; if ($$r.StatusCode -ne 200) { exit 1 }; if ($$r.Content -notmatch '\"status\"\s*:\s*\"ok\"') { exit 1 }"
else
	curl -sf http://localhost:$(HTTP_PORT)/healthcheck | grep -q '"status":"ok"'
endif

# ============================================================
# Task - Raise infrastructure (PostgreSQL in Docker)
# Compose file: ./deploy/docker/docker-compose.yml
# ============================================================
.PHONY: up down restart start stop logs clean-image
# Task - Start all services
up:
	docker-compose -f $(DC) --project-name $(APP_NAME) up -d
# Task - Start specific service
start:
	docker-compose -f $(DC) --project-name $(APP_NAME) start

# Task - Stop all services
stop:
	docker-compose -f $(DC) --project-name $(APP_NAME) stop

# Task - Show logs of all services
logs:
	docker-compose -f $(DC) --project-name $(APP_NAME) logs -f

# Task - Restart all services
restart:
	docker-compose -f $(DC) --project-name $(APP_NAME) restart

# Task - Down all services
down:
	docker-compose -f $(DC) --project-name $(APP_NAME) down

# Task - Clean all: containers, networks and volumes (carefully!) will be deleted ALL!
clean-image:
	docker-compose -f $(DC) --project-name $(APP_NAME) down -v
ifeq ($(OS),Windows_NT)
	$(RM_RF)
else
	$(RM_RF) $(BUILD_DIR)
endif

# Task - Reset active profile to development (default)
.PHONY: env-default
env-default:
	@echo dev>$(ACTIVE_PROFILE_FILE)
	@echo Active profile reset to default: APP_ENV=dev → ENV_FILE=.env.dev

# Task - Select and remember profile (ENV_FILE follows automatically)
# Usage: make env-use APP_ENV=dev|prod|test
.PHONY: env-use
env-use:
ifeq ($(origin APP_ENV),command line)
	@echo $(APP_ENV)>$(ACTIVE_PROFILE_FILE)
	@echo Active profile set: APP_ENV=$(APP_ENV) → ENV_FILE=.env.$(APP_ENV)
else
	$(error Usage: make env-use APP_ENV=dev|prod|test)
endif

# Task - Show which env profile Make will use (safe: no password printed)
.PHONY: env-info
env-info:
	@echo "APP_ENV=$(APP_ENV)"
	@echo "ENV_FILE=$(ENV_FILE)"
	@echo "DB_DRIVER=$(DB_DRIVER)"
	@echo "DB_HOST=$(DB_HOST)"
	@echo "DB_PORT=$(DB_PORT)"
	@echo "DB_NAME=$(DB_NAME)"
	@echo "DB_USER=$(DB_USER)"
	@echo "DB_SSLMODE=$(DB_SSLMODE)"
	@echo "DB_DSN=$(DB_DRIVER)://$(DB_USER):***@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)"

# Task - Apply all pending migrations
.PHONY: migrate-up
migrate-up:
	@echo "Using ENV_FILE=$(ENV_FILE)"
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" up

# Task - Roll back the last migration
.PHONY: migrate-down
migrate-down:
	@echo "Using ENV_FILE=$(ENV_FILE)"
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" down

#Task - Check migration status
.PHONY: migrate-status
migrate-status:
	@echo "Using ENV_FILE=$(ENV_FILE)"
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" status

# Dev/demo data (separate from DDL). Dir: scripts/seeds. Table: goose_seed_version (not goose_db_version).
# Precheck inside scripts/seeds/*.sql (IF NOT EXISTS / RAISE) — goose has no Liquibase preConditions.
.PHONY: seed
seed:
	@echo "Using ENV_FILE=$(ENV_FILE) SEEDS_DIR=$(SEEDS_DIR)"
	goose -dir $(SEEDS_DIR) -table $(GOOSE_SEED_TABLE) postgres "$(DB_DSN)" up

.PHONY: seed-down
seed-down:
	@echo "Using ENV_FILE=$(ENV_FILE) SEEDS_DIR=$(SEEDS_DIR)"
	goose -dir $(SEEDS_DIR) -table $(GOOSE_SEED_TABLE) postgres "$(DB_DSN)" down

.PHONY: seed-status
seed-status:
	@echo "Using ENV_FILE=$(ENV_FILE) SEEDS_DIR=$(SEEDS_DIR)"
	goose -dir $(SEEDS_DIR) -table $(GOOSE_SEED_TABLE) postgres "$(DB_DSN)" status

# Task - A full run, like in CI: formatting, linter, tests
.PHONY: check
check:
	$(MAKE) fmt
	$(MAKE) lint
	$(MAKE) test
	$(MAKE) coverage

#Task - Vulnerability detection tool
.PHONY: vulncheck
vulncheck:
	$(GO)vulncheck $(GO_PKG)

# Task - Generate coverage report in HTML format
.PHONY: coverage
coverage:
	$(GO) test -coverprofile=coverage.out $(GO_PKG)
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: file://$(shell pwd)/coverage.html"

# Task - Output coverage to the terminal (optional)
.PHONY: cover-report
cover-report:
	$(GO) test -cover $(GO_PKG)

# Task - Default - Run all checks from the list
.PHONY: all
all: 
	$(MAKE) lint
	$(MAKE) test
	$(MAKE) coverage

# ============================================================
# Task - One line call - the script does everything itself
# ============================================================
.PHONY: yaml-check
yaml-check:
	@echo "OS: $(DETECTED_OS)"
	$(YAML_CHECK_SCRIPT)

# ============================================================
# Task - Show information about the OS and yq
# ============================================================
.PHONY: info
info:
	@echo "OS: $(DETECTED_OS)"
ifeq ($(OS),Windows_NT)
	@where yq >nul 2>&1 && yq --version || echo yq: not installed
else
	@echo "yq: $$(yq --version 2>/dev/null || echo 'not installed')"
endif

# DETECTED_OS: Windows | Linux | Darwin — для yaml-check и info