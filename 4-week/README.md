# Неделя 4 - Weather + Users API (Auth + Security)

### Тема

Доработка сервиса Weather Tracker API из предыдущей недели: добавить аутентификацию, авторизацию, middleware и защищенные маршруты.

### Цель

Нужно расширить существующий REST API и реализовать:

- JWT-аутентификацию;
- middleware для чтения и проверки токена;
- защищенные маршруты;
- роли пользователей `user` и `admin`;
- безопасную работу с паролями;
- получение текущего пользователя из JWT, а не из `user_id` в URL.

Точка входа сервера находится в `cmd/api/main.go`.
HTTP transport находится в `internal/adapters/transport/http/v1`.
Swagger UI статика находится в `web/swagger`.
OpenAPI спецификация находится в `api/v1/openapi.yaml`.

### Что нужно изменить по техзаданию

#### 1. Аутентификация

Добавить публичные endpoints:

- `POST /auth/register`
- `POST /auth/login`

Требования:

- при регистрации email должен быть уникальным;
- пароль нужно хранить только в виде bcrypt-хеша;
- `POST /auth/login` должен возвращать `access_token` в формате JWT.

#### 2. JWT

JWT payload должен содержать:

- `user_id`
- `email`
- `role`
- `exp`

#### 3. Middleware

Нужно реализовать `AuthMiddleware`, который:

- читает `Authorization` header;
- проверяет формат `Bearer <token>`;
- валидирует подпись токена;
- проверяет `exp`;
- парсит claims;
- кладет текущего пользователя в `context`.

#### 4. Защищенные маршруты

Для пользовательских операций нужно убрать `user_id` из URL.

Все операции должны работать через пользователя, извлеченного из JWT.

Целевые маршруты:

- `POST /cities`
- `GET /cities`
- `DELETE /cities/{city_id}`
- `GET /weather`
- `GET /weather/history`
- `GET /users/me`

#### 5. Роли

Поддерживаемые роли:

- `user`
- `admin`

Только `admin` должен иметь доступ к маршрутам:

- `GET /users`
- `GET /users/{id}`
- `DELETE /users/{id}`

Для этого нужно использовать отдельную проверку роли, например `RequireRole("admin")`.

#### 6. Безопасность

Обязательные требования:

- использовать bcrypt для паролей;
- хранить JWT secret в `.env`;
- не хранить пароль в plain text;
- проверять срок жизни токена;
- проверять подпись токена.

#### 7. Архитектура

Сохраняем слоистую структуру:

- handler
- service
- repository

Текущий пользователь должен попадать в handler/service только через JWT claims и `context`.

#### 8. Flow запроса

Целевой сценарий:

`Request -> AuthMiddleware -> Handler -> Service -> Repository`

### Рекомендуемая модель пользователя

Для week 4 в модели пользователя должны появиться поля:

- `password_hash`
- `role`

Рекомендуемая роль по умолчанию при регистрации: `user`.

### Рекомендуемые env-переменные

Текущий `.env.example` уже содержит настройки сервера, PostgreSQL и Redis.
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
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DATABASE=weather_api
POSTGRES_SSL_MODE=disable
POSTGRES_MAX_CONNS=10
POSTGRES_MIN_CONNS=1
POSTGRES_MAX_CONN_LIFETIME=1h
POSTGRES_MAX_CONN_IDLE_TIME=30m

REDIS_ADDR=redis:6379
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_DIAL_TIMEOUT=5s
REDIS_READ_TIMEOUT=3s
REDIS_WRITE_TIMEOUT=3s
REDIS_WEATHER_CACHE_TTL=10m
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

### Что должно получиться после выполнения week 4

Минимальный ожидаемый результат:

- пользователь регистрируется через `POST /auth/register`;
- пользователь логинится через `POST /auth/login`;
- сервис возвращает JWT access token;
- защищенные ручки читают пользователя из токена;
- маршруты `/cities`, `/weather`, `/weather/history`, `/users/me` работают без `user_id` в URL;
- admin получает доступ к просмотру и удалению пользователей;
- обычный `user` не может выполнять admin-only операции.

### Чеклист реализации

- [x] Добавить поля `password_hash` и `role` в модель и БД.
- [ ] Реализовать `POST /auth/register`.
- [ ] Реализовать `POST /auth/login`.
- [x] Добавить генерацию и валидацию JWT.
- [ ] Реализовать `AuthMiddleware`.
- [ ] Реализовать middleware или helper для проверки роли admin.
- [ ] Убрать `user_id` из пользовательских weather/city маршрутов.
- [ ] Добавить `GET /users/me`.
- [ ] Обновить OpenAPI спецификацию.
- [ ] Обновить Swagger UI.
- [x] Дополнить `.env.example` JWT-настройками.
