# ShareTrip Contract Service

Проект ShareTrip состоит из нескольких сервисов:

1. ShareTrip — основной сервис поездок. Он управляет созданием поездки, маршрутом, участниками, статусами и жизненным циклом поездки.
2. Notification Service — сервис уведомлений. Он отвечает за доставку уведомлений пользователям через разные каналы.
3. Keycloak — внешний сервис аутентификации и авторизации. Он отвечает на вопрос, кто пользователь и какие у него роли.
4. ShareTrip Contract Service — текущий модуль.

## Назначение сервиса

Это отдельный сервис договоров со своей бизнес-ответственностью.

Основная задача `sharetrip_contract` — управление договорными отношениями между платформой ShareTrip и компаниями-клиентами.

## Ответственность сервиса

Сервис отвечает за:
- хранит договоры;
- хранение статуса и срока действия договора;
- хранение услуг, доступных в рамках договора;
- проверку доступности услуги для компании;
- предоставление API другим сервисам;
- владение своей схемой данных.

## Граница ответственности

### В scope - Что входит в границу сервиса

Contract Service отвечает только за договорные условия и доступность услуг в рамках договора.

### Вне scope - Что не входит в границу сервиса

- создание и управление поездками;
- доставка уведомлений;
- хранение пользовательских профилей;
- аутентификация и роли пользователя (это зона Keycloak и вызывающего сервиса).

---

## Локальный стенд: куда ходить

Поднять инфраструктуру:

```powershell
# 1) ShareTrip — Jaeger + OTEL Collector (общие для всех микросервисов)
cd ..\share_trip_260626\share_trip_june270626\share_trip
docker compose -f deploy/docker-compose.yml up -d jaeger otel-collector

# 2) Contract — своя БД + Swagger UI
cd ..\sharetrip_contract
make up
make migrate-up
make run
```

### Порты и URL

| Что | URL / host:port |
|-----|-----------------|
| **ShareTrip API** | http://localhost:8080 |
| **Contract API** | http://localhost:8082 |
| **Contract healthcheck** | http://localhost:8082/healthcheck |
| **Swagger UI (contract yaml)** | http://localhost:8086 |
| **OpenAPI yaml с running app** | http://localhost:8082/api/v2/contracts/openapi |
| **Jaeger UI (общий)** | http://localhost:16686 |
| **OTEL Collector (экспорт с хоста)** | `localhost:4319` |
| **Contract PostgreSQL** | `localhost:6547` |
| **ShareTrip PostgreSQL** | `localhost:6543` |
| **Keycloak (ShareTrip compose)** | http://localhost:8087 |
| **Grafana (ShareTrip compose)** | http://localhost:3000 |
| **Kafka UI (ShareTrip compose)** | http://localhost:9000 |

Подробно про единый Jaeger: [`.docs/cheatsheets/shared-jaeger-otel-cheatsheet.md`](.docs/cheatsheets/shared-jaeger-otel-cheatsheet.md).

### HTTP API Contract Service

```text
GET  /healthcheck
GET  /api/v2/contracts/openapi
POST /api/v2/contracts/
GET  /api/v2/contracts/active?companyId={companyId}
GET  /api/v2/companies/{companyId}/services/{serviceCode}/availability
```

Для `trip_creation → allowed: true` без `PUT /services` используй `make seed` (или SQL в `scripts/seeds`).

`PUT /services`, `GET /contracts/{contractId}` и `PATCH` статуса убраны из кода (остались в `api/contract.yaml`).

### Трейсинг в Jaeger (Contract)

| Endpoint | Spans в Jaeger |
|----------|----------------|
| `POST /api/v2/contracts/` | ✅ handler → service → usecase → storage |
| `GET /api/v2/contracts/active` | ✅ полная цепочка |
| `GET /api/v2/companies/.../availability` | ✅ полная цепочка |
| `GET /api/v2/contracts/openapi` | ❌ |
| `GET /healthcheck` | ❌ |

В Jaeger выбери сервис **`sharetrip-contract`** (не путать с `share-trip`).

---

## Makefile

### Проверка что приложение и БД активны

```powershell
make e2e
```

Проверяет `GET http://localhost:8082/healthcheck` (порт из `HTTP_PORT` / `.env.dev`).

### Автоматические тесты

```powershell
make test
make test-integration
```

Integration-тест сам запускает PostgreSQL через Testcontainers. Jaeger для `go test` не нужен.

Подробно: [component / integration / E2E](.docs/cheatsheets/component-integration-e2e-testing-cheatsheet.md).

---

## Алгоритм проверки доступности услуги

ShareTrip вызывает:

```text
GET /api/v2/companies/{companyId}/services/{serviceCode}/availability
```

Contract Service проверяет: компания известна → код услуги в словаре → активный договор + `contract_services`.

Успех: `200` + `allowed: true/false` + `reason`.  
Ошибки: `400` валидация, `404` company/service not found.

---

## Технологический стек

- Go
- PostgreSQL
- SQL-миграции (goose)
- Docker / docker-compose
- OpenTelemetry → OTEL Collector → Jaeger (общий стек ShareTrip)
- Testcontainers for Go

## Документация проекта

- API-контракт: `api/contract.yaml`
- Карта веток: `.docs/git-branches.md`
- Observability: `.docs/cheatsheets/shared-jaeger-otel-cheatsheet.md`
- Postman (ручные сценарии): `.docs/cheatsheets/postman-collections-cheatsheet.md`
