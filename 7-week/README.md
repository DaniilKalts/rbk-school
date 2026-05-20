# Неделя 7 — Docker, CI/CD и разделение проекта на сервисы

### Тема

Продолжаем работу над проектом из прошлого ДЗ, где добавляли тесты, logging middleware, structured logging, mock'и, dependency injection, integration test и т.д.

### Цель

Научиться:

- разделять backend-приложение на несколько сервисов;
- контейнеризировать Go-сервисы через Docker;
- настраивать общий `docker-compose.yml`;
- организовывать взаимодействие сервисов между собой;
- выносить работу с внешним интернет API в отдельный gateway-сервис.

### Основная задача

Разделить проект на 2 сервиса:

1. **API Service** — основной сервис с бизнес-логикой, REST API и работой с PostgreSQL.
2. **Gateway Service** — отдельный сервис, который принимает запросы и ходит во внешний интернет API.

### Что нужно реализовать

#### 1. API Service

Должен содержать:

- REST endpoints из прошлого проекта;
- подключение к PostgreSQL;
- service/repository слой;
- dependency injection;
- structured logging;
- logging middleware;
- unit тесты;
- integration test repository слоя.

#### 2. Gateway Service

Должен:

- принимать HTTP-запросы;
- ходить во внешний публичный API;
- использовать `http.Client` с timeout;
- не работать напрямую с БД;
- вызывать API Service через внутренний адрес Docker Compose.

### Docker требования

У каждого сервиса должен быть свой Dockerfile:

```
api-service/Dockerfile
gateway-service/Dockerfile
```

В корне проекта должен быть общий:

```
docker-compose.yml
```

Он должен поднимать:

- `api-service`
- `gateway-service`
- `postgres`

### Пример структуры проекта

```
project/
├── docker-compose.yml
├── api-service/
│   ├── Dockerfile
│   ├── go.mod
│   ├── cmd/
│   └── internal/
└── gateway-service/
    ├── Dockerfile
    ├── go.mod
    ├── cmd/
    └── internal/
```

### Проверка запуска

Проект должен запускаться одной командой:

```bash
docker compose up --build
```

После запуска должны работать:

```
GET http://localhost:8080/health
GET http://localhost:8081/health
```

### Что нельзя делать

- ходить в production БД;
- использовать реальные внешние API в тестах;
- писать бизнес-логику в handler;
- делать тесты зависимыми друг от друга;
- проверять только happy path.

### Критерии оценки

- проект разделен на 2 сервиса;
- оба сервиса запускаются через Docker Compose;
- Gateway ходит во внешний API;
- API Service вызывает Gateway;
- API Service сохраняет данные в PostgreSQL;
- есть тесты, DI и structured logging;
- README понятный и полный.

### Чеклист реализации

- [x] Проект разделен на 2 сервиса: `api-service` и `gateway-service`.
- [x] У каждого сервиса свой `Dockerfile`.
- [x] В корне проекта общий `docker-compose.yml`.
- [x] `docker-compose.yml` поднимает `api-service`, `gateway-service`, `postgres`.
- [x] Проект запускается одной командой `docker compose up --build`.
- [x] `GET http://localhost:8080/health` работает (API Service).
- [x] `GET http://localhost:8081/health` работает (Gateway Service).
- [x] API Service перенесен из прошлого проекта (endpoints, PostgreSQL, тесты, DI, logging).
- [x] Gateway Service ходит во внешний публичный API через `http.Client` с timeout.
- [x] Gateway Service не работает напрямую с БД.
- [x] Gateway Service вызывает API Service через внутренний адрес Docker Compose.
- [x] API Service вызывает Gateway.

---

## Проект

### Запуск

```bash
docker compose up --build
```

Поднимает: `api-service` (`:8080`), `gateway-service` (`:8081`), `postgres` (`:5432`), `redis` (`:6379`). Дополнительно: `pgadmin` (`:5050`) и `redis-commander` (`:8082`) — UI-утилиты для локальной отладки.

При старте api-service автоматически прогоняет миграции через `goose`.

### Поток данных

```
[client] → api-service (REST + JWT + Postgres) → gateway-service → open-meteo
                                                       ↑
                                  gateway периодически пингует api-service /health
                                  через /api/v1/ready (readiness check)
```

api-service хранит пользователей, города, историю запросов. Gateway-service — тонкий прокси к внешним публичным API (`geocoding-api.open-meteo.com`, `api.open-meteo.com`), без БД.

### Эндпоинты

**API Service (`:8080`)**

| Метод | Путь | Назначение |
|---|---|---|
| `GET` | `/health` | Healthcheck |
| `POST` | `/api/v1/auth/register` | Регистрация пользователя |
| `POST` | `/api/v1/auth/login` | Логин, возвращает JWT |
| `POST` | `/api/v1/auth/logout` | Logout (требует Bearer) |
| `GET` | `/api/v1/users/me` | Текущий пользователь (Bearer) |
| `GET` | `/api/v1/cities` | Список городов пользователя (Bearer) |
| `POST` | `/api/v1/cities` | Добавить город (Bearer) |
| `DELETE` | `/api/v1/cities/{city_id}` | Удалить город (Bearer) |
| `GET` | `/api/v1/weather` | Погода по городам пользователя (Bearer) — ходит через gateway |
| `GET` | `/api/v1/weather/history` | История запросов (Bearer) |
| `GET` | `/api/v1/users` | Список пользователей (admin) |
| `GET` | `/api/v1/swagger/*` | Swagger UI |

**Gateway Service (`:8081`)**

| Метод | Путь | Назначение |
|---|---|---|
| `GET` | `/health` | Healthcheck |
| `GET` | `/api/v1/weather?city=NAME` | Погода по городу (geocoding + openmeteo) |
| `GET` | `/api/v1/ready` | Deep readiness (пингует api-service `/health`) |

### Переменные окружения

Все настройки берутся из корневого `.env`. Ключевое:

**Общие**

| Переменная | Дефолт | Сервис |
|---|---|---|
| `LOG_LEVEL` | `info` | оба |
| `LOG_FORMAT` | `json` | оба |
| `SERVER_HTTP_TIMEOUT` | `15s` | api |
| `SERVER_HANDLER_TIMEOUT` | `10s` | оба |
| `SERVER_SHUTDOWN_TIMEOUT` | `15s` | оба |

**API Service**

| Переменная | Дефолт |
|---|---|
| `SERVER_ADDR` | `:8080` |
| `POSTGRES_HOST` / `PORT` / `DATABASE` / `USER` / `PASSWORD` | — |
| `POSTGRES_SSL_MODE` | `disable` |
| `REDIS_ADDR` / `DB` / `PASSWORD` | — |
| `REDIS_WEATHER_CACHE_TTL` | `10m` |
| `JWT_SECRET` | `change-me` |
| `JWT_ACCESS_TOKEN_TTL` | `15m` |
| `GATEWAY_BASE_URL` | `http://gateway-service:8081` |
| `GATEWAY_TIMEOUT` | `15s` |

**Gateway Service**

| Переменная | Дефолт |
|---|---|
| `SERVER_ADDR` | `:8081` |
| `EXTERNAL_GEOCODING_URL` | `https://geocoding-api.open-meteo.com/v1/search` |
| `EXTERNAL_OPENMETEO_URL` | `https://api.open-meteo.com/v1/forecast` |
| `EXTERNAL_TIMEOUT` | `15s` |
| `API_SERVICE_BASE_URL` | `http://api-service:8080` |
| `API_SERVICE_TIMEOUT` | `5s` |

### Структура

```
.
├── docker-compose.yml
├── .dockerignore
├── .env
├── go.mod / go.sum
├── api-service/
│   ├── Dockerfile
│   ├── api/v1/                       # OpenAPI спека
│   ├── cmd/api/main.go
│   ├── database/migrations/
│   ├── database/queries/             # sqlc вход
│   ├── sqlc.yaml
│   ├── web/swagger/                  # Swagger UI
│   └── internal/
│       ├── app/{app,container}.go    # DI + HTTP lifecycle
│       ├── config/                   # Server, Postgres, Redis, JWT, Logger, Gateway
│       ├── domain/{user,city,history,weather}/
│       ├── repository/{user,city,weather}/
│       ├── cache/{blacklist,weather}/
│       ├── service/{auth,user,city,weather}/
│       └── adapter/
│           ├── client/gateway/       # HTTP-клиент к gateway-service
│           ├── cache/redis/
│           ├── database/postgres/
│           └── transport/http/
│               ├── router.go
│               ├── middleware/       # auth, role (RequestID/Logger в pkg/)
│               └── v1/{auth,user,city,weather}/
├── gateway-service/
│   ├── Dockerfile
│   ├── cmd/gateway/main.go
│   └── internal/
│       ├── app/{app,container}.go
│       ├── config/                   # Server, Logger, External, APIService
│       ├── service/weather/          # geocoding + openmeteo оркестрация
│       └── adapter/
│           ├── client/{apiservice,geocoding,openmeteo}/
│           └── transport/http/
│               ├── router.go
│               └── v1/{weather,health}/
└── pkg/
    ├── logger/                       # Zap factory
    ├── jwt/                          # JWT manager
    ├── httpx/                        # JSON/Error helpers, request_id + claims context
    └── middleware/                   # RequestID, Logger
```

### Тесты

```bash
# unit + handler-level (offline)
go test ./api-service/... ./gateway-service/... ./pkg/...

# integration (поднимает Postgres/Redis в testcontainers, нужен Docker)
go test -tags=integration ./api-service/internal/repository/... ./api-service/internal/cache/...
```
