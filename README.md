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

Сервис не создает поездки, не управляет жизненным циклом поездки, не отправляет уведомления и не принимает решение, кто именно аутентифицирован в системе.

## Локальный стенд: куда ходить

Поднять инфраструктуру:

```powershell
# 1) ShareTrip — Jaeger + OTEL Collector (общие для всех микросервисов)
docker compose -f deploy/docker-compose.yml up -d jaeger otel-collector

# 2) Contract — своя БД + Swagger UI
make up
make migrate-up
make run
```

## Бизнес-сценарий (happy path)

В рамках сценария создания поездки поток выглядит так:

1. Пользователь отправляет запрос на создание поездки в ShareTrip.
2. ShareTrip определяет `company_id` пользователя.
3. ShareTrip обращается в Contract Service.
4. Contract Service проверяет:
   - есть ли у компании договор;
   - активен ли договор;
   - не истек ли срок действия договора;
   - разрешена ли услуга `trip_creation`.
5. Contract Service возвращает результат проверки.
6. Если услуга доступна, ShareTrip продолжает выполнение сценария.
7. Если услуга недоступна, ShareTrip возвращает отказ.

## Минимальные бизнес-правила в БД

- у компании может быть только один `active` договор;
- одна и та же услуга не дублируется в рамках одного договора;
- статус договора ограничен допустимым набором;
- даты договора валидны (`start_date <= end_date`).

Эти правила важны не только как описание логики, но и как ограничения, которые должны быть выражены в схеме БД через `PK`, `FK`, `UNIQUE`, `CHECK` и индексы.

## Основной бизнес-вопрос сервиса

Перед выполнением действия другой сервис должен иметь возможность спросить:

`Может ли компания использовать конкретную услугу сейчас?`

Если договор активен, срок его действия не истек и услуга разрешена, Contract Service возвращает положительный ответ.
Если договор отсутствует, неактивен (в статусе suspended) или истек( terminated) или услуга запрещена, сервис возвращает отказ.

### HTTP API Contract Service

```text
GET  /healthcheck
GET  /api/v2/contracts/openapi
POST /api/v2/contracts/
GET  /api/v2/contracts/active?companyId={companyId}
GET  /api/v2/companies/{companyId}/services/{serviceCode}/availability
```

### Трейсинг в Jaeger (Contract)

| Endpoint | Spans в Jaeger |
|----------|----------------|
| `POST /api/v2/contracts/` | ✅ handler → service → usecase → storage |
| `GET /api/v2/contracts/active` | ✅ полная цепочка |
| `GET /api/v2/companies/.../availability` | ✅ полная цепочка |
| `GET /api/v2/contracts/openapi` | ✅ |
| `GET /healthcheck` | ✅ |

В Jaeger выберите сервис **`sharetrip-contract`** (не путать с `share-trip`).

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

### Makefile

##### Проверка что приложение и БД активны и работают
В Makefile e2e  /healthcheck.

```powershell
# open terminal and setup
make e2e
```

Проверяет `GET http://localhost:8082/healthcheck` (порт из `HTTP_PORT` / `.env.dev`).

#### Автоматические тесты

Быстрые unit/component-тесты и интеграционный API-набор:

```powershell
make test
make test-integration
```

Integration-тест сам запускает PostgreSQL через Testcontainers. Jaeger для `go test` не нужен.

##### Бизнес-сценарий (sad path)
Возможные причины отказа:

При создании поездки с проверкой доступности услуги поток выглядит так:

1. Пользователь отправляет запрос на создание поездки в ShareTrip.
2. ShareTrip определяет `company_id` пользователя.
3. ShareTrip обращается в Contract Service.
   ```curl
   GET /companies/{companyId}/services/{serviceCode}/availability
   ```
4. Contract Service проверяет:
- валидация что данные ID договора корректны -> нет
 -> возвращает ошибку со статусом http code 400
   http code 400
   ```json
      {
      "code": "COMPANY_VALIDATE_ERROR",
      "message": "Company id is invalid"
      }
   ```
- найти договор по company_id -> наличие у данной компании договор -> нет 
  -> возвращает ошибку со статусом http code 404
   ```json
      {
      "code": "CONTRACT_NOT_FOUND",
      "message": "Active contract for the company was not found"
      }
   ```
   
 - активен ли договор 
   проверка статуса договора(
      - не находится ли договор в статусе suspended или terminated;
      - не истек ли срок действия договора(на текущую дату);
      - не начнется ли договор в будущем (start_date > now());
      ) -> нет
 ->  возвращает результат проверки  
для бизнес-отказ (договор suspended, expired, service disabled) → 200 + allowed: false + reason.

 - проверка на то что компании доступная данная услуга (serviceCode) -> нет
 -> возвращает результат проверки  
для бизнес-отказ (договор suspended, expired, service disabled) → 200 + allowed: false + reason.

5. Contract Service возвращает результат проверки в ShareTrip.

 -> возврат в вызывающий сервис ответа 
       со статусом http code 200
```json
{
  "company_id": 4829104857,
  "service_code": "trip_creation",
  "allowed": false,
  "reason": "The company contract status is suspended."
}

```

6. Если услуга недоступна, сервис ShareTrip возвращает отказ.


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

## Projects Ports 
Порты (не путать с ShareTrip):

| Сервис                         | Порт       |
|--------------------------------|------------|
| ShareTrip API                  | `8080`     |
| **ShareTrip Notification API** | **`8081`** |
| **Contract API**               | **`8082`** |
| Swagger UI (contract)          | `8086`     |
| **Kafka UI**                   | **`9000`** |

---

