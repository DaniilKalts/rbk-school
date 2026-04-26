# Неделя 3 - Weather Tracker API

### Тема

Разработка REST API сервиса на Go с использованием PostgreSQL, Redis, Docker, Swagger UI и внешнего API погоды.

### Проект

Проект: сервис для отслеживания погоды по городам пользователя.

API позволяет:

- управлять пользователями;
- добавлять пользователю города для отслеживания;
- получать актуальную погоду по всем городам пользователя;
- сохранять каждый погодный запрос в историю;
- получать историю погоды с фильтрацией по городу;
- получать всю историю пользователя без фильтра по городу;
- использовать `limit` и `offset` для пагинации истории;
- кешировать погодные данные в Redis;
- выполнять параллельные запросы к Weather API.

Точка входа HTTP-сервера находится в `cmd/api/main.go`.
HTTP transport находится в `internal/adapters/transport/http/v1`.
Swagger UI статика находится в `web/swagger`.
OpenAPI спецификация находится в `api/v1/openapi.yaml`.

Для получения погоды используется Open-Meteo:

- Geocoding API: `https://geocoding-api.open-meteo.com/v1/search` - поиск города и получение координат;
- Forecast API: `https://api.open-meteo.com/v1/forecast` - получение текущей погоды по координатам.

Схема работы погоды:

`пользователь -> список городов -> координаты города -> текущая погода -> запись в weather_history -> JSON-ответ`

### Запуск

Перед запуском нужно подготовить файл окружения.

В проекте есть шаблон `3-week/.env.example`:

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

Создайте `.env` на основе шаблона:

```bash
cp .env.example .env
```

Запуск через Docker Compose:

```bash
docker compose up --build
```

После запуска API доступен на `http://localhost:8080`.

Swagger UI доступен по адресу:

```text
http://localhost:8080/swagger/
```

OpenAPI spec доступна по адресу:

```text
http://localhost:8080/api/v1/openapi.yaml
```

Локальный запуск без Docker тоже возможен, но PostgreSQL и Redis должны быть доступны отдельно:

```bash
go run ./cmd/api
```

### Ручки

### `GET /health`

Проверка, что сервис запущен и отвечает.

### `GET /swagger/`

Открывает Swagger UI в браузере.

### `GET /api/v1/openapi.yaml`

Возвращает OpenAPI spec в формате YAML.

### Пользователи

### `POST /api/v1/users`

Создает нового пользователя.

Тело запроса:

```json
{
  "first_name": "Daniil",
  "last_name": "Kalts",
  "email": "daniil@example.com"
}
```

### `GET /api/v1/users`

Возвращает список активных пользователей.

Удаленные пользователи не возвращаются, потому что используется soft delete через `deleted_at`.

### `GET /api/v1/users/{id}`

Возвращает активного пользователя по UUID.

### `PUT /api/v1/users/{id}`

Обновляет имя, фамилию и email активного пользователя.

### `DELETE /api/v1/users/{id}`

Мягко удаляет пользователя.

Запись остается в базе данных, но больше не возвращается в пользовательских ручках.

### Города Пользователя

### `POST /api/v1/users/{id}/cities`

Добавляет город в список отслеживания пользователя.

Тело запроса:

```json
{
  "city": "Almaty"
}
```

### `GET /api/v1/users/{id}/cities`

Возвращает список городов, за погодой в которых следит пользователь.

### `DELETE /api/v1/users/{id}/cities/{city_id}`

Удаляет город из списка отслеживания.

Это обычное удаление записи из `user_cities`. История погоды при этом остается в `weather_history`.

### Погода

### `GET /api/v1/users/{id}/weather`

Возвращает актуальную погоду по всем городам пользователя.

Что делает ручка:

- проверяет, что пользователь существует и не удален;
- получает список городов пользователя;
- параллельно запрашивает погоду по каждому городу;
- использует Redis-кеш для погодных данных;
- сохраняет каждую запись в `weather_history`;
- возвращает агрегированный JSON-ответ.

### `GET /api/v1/users/{id}/weather/history`

Возвращает историю погодных запросов пользователя.

Query-параметры:

- `city` - необязательный фильтр по городу;
- `limit` - необязательное количество записей;
- `offset` - необязательное количество записей, которое нужно пропустить.

Если `city` передан, API возвращает историю только по этому городу.

Если `city` не передан, API возвращает всю историю пользователя.

История сортируется по `requested_at DESC`.

### Примеры API

Ниже приведены `curl`-примеры для основных сценариев.

### Проверка Сервиса

```bash
curl http://localhost:8080/health
```

Ожидаемый ответ:

```text
ok
```

### Swagger UI

Открыть Swagger UI в браузере:

```text
http://localhost:8080/swagger/
```

Получить OpenAPI spec:

```bash
curl http://localhost:8080/api/v1/openapi.yaml
```

### Создать Пользователя

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"first_name":"Daniil","last_name":"Kalts","email":"daniil@example.com"}'
```

Пример ответа:

```json
{
  "id": "018f1b79-38f8-7a72-a64b-7c58a85f60f3",
  "first_name": "Daniil",
  "last_name": "Kalts",
  "email": "daniil@example.com",
  "created_at": "2026-04-26T10:00:00Z",
  "updated_at": "2026-04-26T10:00:00Z"
}
```

### Получить Пользователей

```bash
curl http://localhost:8080/api/v1/users
```

### Добавить Город Пользователю

Замените `USER_ID` на UUID пользователя.

```bash
curl -X POST http://localhost:8080/api/v1/users/USER_ID/cities \
  -H "Content-Type: application/json" \
  -d '{"city":"Almaty"}'
```

Пример ответа:

```json
{
  "id": "018f1b79-a83f-7f8d-93e1-31db8fcce245",
  "user_id": "018f1b79-38f8-7a72-a64b-7c58a85f60f3",
  "city": "Almaty",
  "created_at": "2026-04-26T10:20:00Z"
}
```

### Получить Города Пользователя

```bash
curl http://localhost:8080/api/v1/users/USER_ID/cities
```

### Получить Погоду По Городам Пользователя

Перед этим у пользователя должен быть хотя бы один город.

```bash
curl http://localhost:8080/api/v1/users/USER_ID/weather
```

Пример ответа:

```json
{
  "user_id": "018f1b79-38f8-7a72-a64b-7c58a85f60f3",
  "weather": [
    {
      "city": "Almaty",
      "temperature": 18.4,
      "feels_like": 17.9,
      "description": "partly cloudy",
      "requested_at": "2026-04-26T10:30:00Z"
    },
    {
      "city": "Astana",
      "temperature": 9.1,
      "feels_like": 5.8,
      "description": "rain",
      "requested_at": "2026-04-26T10:30:01Z"
    }
  ]
}
```

Что важно:

- внешний Open-Meteo API не дергается повторно, если погода по городу есть в Redis-кеше;
- TTL кеша настраивается через `REDIS_WEATHER_CACHE_TTL`;
- запись в `weather_history` создается при каждом запросе, даже если погода пришла из кеша.

### Получить Историю По Городу

```bash
curl "http://localhost:8080/api/v1/users/USER_ID/weather/history?city=Almaty"
```

Пример ответа:

```json
{
  "user_id": "018f1b79-38f8-7a72-a64b-7c58a85f60f3",
  "city": "Almaty",
  "history": [
    {
      "temperature": 18.4,
      "description": "partly cloudy",
      "requested_at": "2026-04-26T10:30:00Z"
    }
  ]
}
```

### Получить Всю Историю Пользователя

```bash
curl "http://localhost:8080/api/v1/users/USER_ID/weather/history"
```

Пример ответа:

```json
{
  "user_id": "018f1b79-38f8-7a72-a64b-7c58a85f60f3",
  "history": [
    {
      "city": "Almaty",
      "temperature": 18.4,
      "description": "partly cloudy",
      "requested_at": "2026-04-26T10:30:00Z"
    },
    {
      "city": "Astana",
      "temperature": 9.1,
      "description": "rain",
      "requested_at": "2026-04-26T10:25:00Z"
    }
  ]
}
```

### История С Limit И Offset

```bash
curl "http://localhost:8080/api/v1/users/USER_ID/weather/history?limit=10&offset=10"
```

Можно совмещать фильтр по городу и пагинацию:

```bash
curl "http://localhost:8080/api/v1/users/USER_ID/weather/history?city=Almaty&limit=10&offset=0"
```

### Примеры Ошибок

Невалидный UUID пользователя:

```bash
curl http://localhost:8080/api/v1/users/not-a-uuid
```

Ответ:

```json
{
  "code": 400,
  "message": "invalid user id"
}
```

Невалидный `limit`:

```bash
curl "http://localhost:8080/api/v1/users/USER_ID/weather/history?limit=0"
```

Ответ:

```json
{
  "code": 400,
  "message": "limit must be a positive number"
}
```

Невалидный `offset`:

```bash
curl "http://localhost:8080/api/v1/users/USER_ID/weather/history?offset=-1"
```

Ответ:

```json
{
  "code": 400,
  "message": "offset must be a non-negative number"
}
```

### Архитектура

Проект разделен на слои:

- `internal/adapters/transport/http/v1` - HTTP ручки, чтение path/query/body параметров, запись JSON-ответов;
- `internal/service` - бизнес-логика пользователей, городов, погоды и истории;
- `internal/repository` - работа с PostgreSQL через sqlc;
- `internal/adapters/client` - клиенты внешних API Open-Meteo;
- `internal/adapters/cache/redis` - Redis-кеш погодных данных;
- `internal/domain` - доменные модели и ошибки;
- `internal/config` - конфигурация приложения.

### База Данных

Используются таблицы:

- `users` - пользователи с soft delete через `deleted_at`;
- `user_cities` - города, которые отслеживает пользователь;
- `weather_history` - история погодных запросов.

Для `weather_history` есть индекс по `(user_id, city)`.

### TODO

### 1. Подготовка И База Данных

- [x] Создать миграции для пользователей.
- [x] Создать миграции для таблицы городов пользователя.
- [x] Создать миграции для таблицы `weather_history`.
- [x] Создать составной индекс по `(user_id, city)` в таблице `weather_history`.

### 2. Управление Пользователями

- [x] Реализовать ручку `POST /users`.
- [x] Реализовать ручку `GET /users`.
- [x] Реализовать ручку `GET /users/{id}`.
- [x] Реализовать ручку `PUT /users/{id}`.
- [x] Реализовать ручку `DELETE /users/{id}`.

### 3. Города Пользователя

- [x] Реализовать ручку `POST /users/{id}/cities`.
- [x] Реализовать ручку `GET /users/{id}/cities`.
- [x] Реализовать ручку `DELETE /users/{id}/cities/{city_id}`.

### 4. Погода

- [x] Интегрировать получение пользователя с проверкой на soft delete.
- [x] Реализовать запрос к внешнему API для получения погоды по городу.
- [x] Агрегировать данные по всем отслеживаемым городам пользователя.
- [x] Организовать сохранение результатов в таблицу `weather_history`.
- [x] Связать бизнес-логику с ручкой `GET /users/{id}/weather`.

### 5. История Погоды

- [x] Реализовать ручку `GET /users/{id}/weather/history`.
- [x] Добавить фильтрацию по query-параметру `city`.
- [x] Использовать безопасные SQL-запросы для защиты от инъекций.
- [x] Внедрить сортировку `ORDER BY requested_at DESC`.
- [x] Реализовать `limit`.
- [x] Собрать итоговый ответ в нужном JSON-формате.

### 6. Дополнительные Задания

- [x] Сделать параметр `city` необязательным.
- [x] Добавить поддержку пагинации через `offset`.
- [x] Реализовать кеширование погодных данных в Redis.
- [x] Реализовать параллельные запросы в Weather API.
