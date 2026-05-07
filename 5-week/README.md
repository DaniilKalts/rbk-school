# Неделя 5 - Рефакторинг проекта под Clean Architecture + подготовка к Unit Tests

### Тема

Рефакторинг существующего проекта Weather Tracker API: разделить код по слоям Clean Architecture и подготовить проект к написанию unit tests.

### Цель

Переработать существующий проект так, чтобы код был разделён по слоям:

`handler -> service -> repository`

После этого на следующем занятии можно будет начать писать unit tests для service layer.

Точка входа сервера находится в `cmd/api/main.go`.
HTTP transport находится в `internal/adapters/transport/http/v1`.
Swagger UI статика находится в `web/swagger`.
OpenAPI спецификация находится в `api/v1/openapi.yaml`.

### Что нужно изменить по техзаданию

#### 1. Структура проекта

Основная задача — рефакторинг архитектуры.

Нужно разнести код по слоям:

```
internal/
  handler/
  service/
  repository/
  model/
  dto/
  middleware/
  config/
```

#### 2. Handler

Handler должен:

- принимать HTTP request;
- читать path/query/body;
- вызывать service;
- возвращать HTTP response.

Handler НЕ должен:

- писать SQL;
- содержать бизнес-логику;
- работать напрямую с БД.

#### 3. Service

Service должен содержать бизнес-логику:

- регистрация пользователя;
- проверка пароля;
- генерация JWT;
- работа с погодой;
- проверки прав доступа.

Service НЕ должен знать про HTTP.

#### 4. Repository

Repository отвечает только за работу с БД.

Repository НЕ должен:

- генерировать JWT;
- проверять роли;
- содержать бизнес-логику.

#### 5. Интерфейсы

Service слой должен зависеть от interface.

Пример:

```go
type UserRepository interface {
    GetByEmail(ctx context.Context, email string) (*User, error)
}
```

#### 6. DTO и Model разделить

Пароль никогда не должен попадать в response.

Нужно создать отдельные DTO для API response.

#### 7. Middleware

`AuthMiddleware` должен:

- читать `Authorization` header;
- проверять JWT;
- класть user в `context`.

#### 8. Обязательные endpoints

- `POST /auth/register`
- `POST /auth/login`
- `POST /cities`
- `GET /cities`
- `DELETE /cities/{city_id}`
- `GET /weather`
- `GET /weather/history`
- `GET /users/me`

#### 9. Частые ошибки

- SQL в handler;
- business logic в middleware;
- giant services;
- shared DTO;
- Service зависит от HTTP;
- Repository знает про роли.

#### 10. Flow запроса

Целевой сценарий:

`Request -> AuthMiddleware -> Handler -> Service -> Repository`

### Рекомендуемые env-переменные

Текущий `.env.example` уже содержит настройки сервера, PostgreSQL, Redis, pgAdmin и Redis Commander.
Для week 4 его нужно дополнить как минимум JWT-настройками, например:

```env
JWT_SECRET=change-me
JWT_ACCESS_TOKEN_TTL=15m
```

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

### Минимальные требования

Проект должен:

- использовать JWT;
- использовать bcrypt;
- иметь слои handler/service/repository;
- использовать interfaces;
- быть готовым к unit tests.

### Что должно получиться после выполнения week 5

Минимальный ожидаемый результат:

- код разделён по слоям: handler, service, repository, model, dto, middleware, config;
- handler не содержит SQL и бизнес-логики;
- service содержит бизнес-логику и не знает про HTTP;
- repository отвечает только за работу с БД;
- service зависит от интерфейсов, а не от конкретных реализаций;
- DTO и model разделены — пароль не попадает в response;
- все обязательные endpoints работают;
- проект готов к написанию unit tests для service layer.

### Чеклист реализации

- [x] Разнести код по слоям: handler/service/repository/model/dto/middleware/config.
- [x] Handler принимает request, вызывает service, возвращает response.
- [x] Service содержит бизнес-логику, не знает про HTTP.
- [x] Repository отвечает только за работу с БД.
- [x] Service зависит от интерфейсов (interfaces).
- [x] Разделить DTO и Model — пароль не попадает в response.
- [x] AuthMiddleware читает JWT и кладёт user в context.
- [x] Все обязательные endpoints реализованы.
- [x] Проект готов к unit tests.
