# Неделя 6 - Тестирование и логирование в Go

### Тема

Добавить в существующий проект Weather Tracker API unit-тесты, integration-тест, structured logging и logging middleware.

### Цель

Научиться:

- писать unit-тесты;
- тестировать HTTP handlers;
- использовать mock'и;
- внедрять dependency injection;
- добавлять structured logging;
- писать middleware;
- проверять негативные сценарии.

Точка входа сервера находится в `cmd/api/main.go`.
HTTP transport находится в `internal/adapter/transport/http/v1`.
Middleware находится в `internal/adapter/transport/http/middleware`.
Swagger UI статика находится в `web/swagger`.
OpenAPI спецификация находится в `api/v1/openapi.yaml`.

### Что нужно реализовать по техзаданию

#### 1. Dependency Injection

- handler получает service через конструктор;
- service получает repository через интерфейс;
- зависимости не создаются внутри business logic.

Плохо:

```go
func NewService() *Service {
    db := postgres.New()
    return &Service{db: db}
}
```

Хорошо:

```go
func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}
```

#### 2. Unit тесты Service слоя

Написать unit-тесты минимум для одного service метода.

Обязательно проверить:

- **happy path** — корректный успешный сценарий;
- **негативные сценарии** — пустые данные, invalid ID, entity not found, ошибки repository.

Использовать:

- `testify/assert`;
- `testify/require`;
- mock repository.

#### 3. Mock Repository

Создать mock repository через `github.com/stretchr/testify/mock`.

Пример:

```go
type MockUserRepository struct {
    mock.Mock
}
```

#### 4. Unit тесты Handler слоя

Написать тесты минимум для:

- одного GET endpoint;
- одного POST endpoint.

Проверять:

- HTTP status code;
- JSON response;
- validation;
- ошибки.

Использовать:

- `httptest.NewRequest`;
- `httptest.NewRecorder`;
- `router.ServeHTTP`.

#### 5. Logging Middleware

Добавить middleware логирования запросов.

Что должно логироваться:

- HTTP method;
- path;
- status code;
- duration;
- request_id.

Использовать Uber Zap.

#### 6. Structured Logging

Все новые логи должны быть structured.

Плохо:

```go
log.Println("user not found")
```

Хорошо:

```go
logger.Error(
    "user not found",
    zap.Int64("user_id", id),
)
```

#### 7. Integration Test

Написать минимум один integration test для repository.

Проверить:

- сохранение данных;
- чтение данных;
- работу SQL.

Можно использовать:

- SQLite in-memory;
- PostgreSQL;
- Docker;
- `testcontainers`.

#### 8. Проверка ошибок

Обязательно протестировать:

- bad request;
- invalid JSON;
- not found;
- internal error.

#### 9. Coverage

Покрытие тестами: минимум 60% service слоя.

Проверить командой:

```bash
go test -cover ./...
```

### Что нельзя делать

- ходить в production БД;
- использовать реальные внешние API;
- писать бизнес-логику в handler;
- делать тесты зависимыми друг от друга;
- проверять только happy path.

### Рекомендуемые env-переменные

Текущий шаблон окружения:

```env
SERVER_ADDR=:8080
SERVER_HTTP_TIMEOUT=15s

POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_DATABASE=weather_api
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_SSL_MODE=disable
POSTGRES_MIN_CONNS=1
POSTGRES_MAX_CONNS=10
POSTGRES_MAX_CONN_LIFETIME=1h
POSTGRES_MAX_CONN_IDLE_TIME=30m

PGADMIN_DEFAULT_EMAIL=admin@example.com
PGADMIN_DEFAULT_PASSWORD=admin
PGADMIN_PORT=5050

REDIS_ADDR=redis:6379
REDIS_PORT=6379
REDIS_DB=0
REDIS_PASSWORD=
REDIS_DIAL_TIMEOUT=5s
REDIS_READ_TIMEOUT=3s
REDIS_WRITE_TIMEOUT=3s
REDIS_WEATHER_CACHE_TTL=10m

REDIS_COMMANDER_PORT=8081
REDIS_COMMANDER_USER=admin
REDIS_COMMANDER_PASSWORD=admin

JWT_SECRET=change-me
JWT_ACCESS_TOKEN_TTL=15m
```

### Запуск

Создать `.env` на основе шаблона:

```bash
cp .env.example .env
```

Запуск через Docker Compose:

```bash
docker compose up --build
```

После запуска сервис доступен на `http://localhost:8080`.

pgAdmin:

```text
http://localhost:5050
```

Подключение к PostgreSQL из pgAdmin:

```text
Host: postgres
Port: 5432
Database: weather_api
Username: postgres
Password: postgres
```

Redis Commander:

```text
http://localhost:8081
```

Логин и пароль Redis Commander задаются через `REDIS_COMMANDER_USER` и `REDIS_COMMANDER_PASSWORD`.

Swagger UI:

```text
http://localhost:8080/swagger/
```

OpenAPI spec:

```text
http://localhost:8080/api/v1/openapi.yaml
```

Локальный запуск без Docker:

```bash
go run ./cmd/api
```

Для локального запуска PostgreSQL и Redis должны быть доступны отдельно.

### Запуск тестов

Все тесты:

```bash
go test -v ./...
```

Coverage:

```bash
go test -cover ./...
```

### Что должно быть в Pull Request

1. **Код проекта** — с тестами и logging middleware.
2. **Скриншот выполнения** — `go test -v ./...`.
3. **Скриншот coverage** — `go test -cover ./...`.

### Критерии оценки

**Отлично:**

- чистая архитектура;
- хорошие тесты;
- проверены edge cases;
- structured logging;
- понятный код.

**Хорошо:**

- тесты работают;
- middleware есть;
- mock'и используются.

**Нужно доработать:**

- тестируется только happy path;
- нет mock'ов;
- нет DI;
- handler содержит бизнес-логику;
- отсутствуют негативные сценарии.

### Что должно получиться после выполнения week 6

Минимальный ожидаемый результат:

- handler получает service через конструктор, service получает repository через интерфейс;
- написаны unit-тесты минимум для одного service метода с happy path и негативными сценариями;
- создан mock repository через `testify/mock`;
- написаны unit-тесты минимум для одного GET и одного POST handler'а;
- добавлен logging middleware на базе Uber Zap, логирующий method, path, status, duration, request_id;
- все новые логи structured;
- написан минимум один integration test для repository;
- проверены сценарии bad request, invalid JSON, not found, internal error;
- coverage service слоя минимум 60%.

### Чеклист реализации

- [x] handler получает service через конструктор.
- [x] service получает repository через интерфейс.
- [x] зависимости не создаются внутри business logic.
- [ ] Unit-тесты для одного service метода (happy path).
- [ ] Unit-тесты для одного service метода (негативные сценарии: пустые данные, invalid ID, not found, ошибки repository).
- [ ] Mock repository через `testify/mock`.
- [ ] Unit-тест для одного GET endpoint.
- [ ] Unit-тест для одного POST endpoint.
- [ ] В handler-тестах проверены HTTP status code, JSON response, validation, ошибки.
- [x] Logging middleware на Uber Zap (method, path, status, duration, request_id).
- [x] Все новые логи structured (Uber Zap).
- [ ] Integration test для repository (сохранение, чтение, SQL).
- [ ] Проверены bad request, invalid JSON, not found, internal error.
- [ ] Coverage service слоя ≥ 60% (`go test -cover ./...`).
