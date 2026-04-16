# Неделя 2 - Weather API

## Тема

Разработка HTTP API на Go: роутинг через `chi`, работа с JSON, разделение приложения на слои и интеграция с внешним API.

## Проект

Проект: приложение о погоде.

Нужно расширить Weather API и добавить ручки, которые возвращают погоду по городу, погоду по городам страны и топ самых теплых городов в стране.

API должно получать данные из Open-Meteo, обрабатывать их на стороне сервиса и отдавать клиенту JSON-ответ.

Точка входа HTTP-сервера находится в `cmd/api/main.go`.
HTTP transport находится в `internal/transport/http/v1`.

Для получения погоды нужно использовать Open-Meteo:

- Geocoding API: `https://geocoding-api.open-meteo.com/v1/search` - поиск города и получение координат;
- Forecast API: `https://api.open-meteo.com/v1/forecast` - получение текущей погоды по координатам.

Схема работы: название города -> координаты через Geocoding API -> текущая температура через Forecast API -> JSON-ответ приложения.

## Запуск

Перед запуском нужно подготовить файл окружения.

В проекте есть шаблон `2-week/.env.example`:

```env
COUNTRY_STATE_CITY_API_KEY=ВАШ_КЛЮЧ
SERVER_ADDR=:8080
SERVER_HTTP_TIMEOUT=15s
```

Создайте на его основе файл `.env` и подставьте ваш API key для CountryStateCity.

Дополнительно можно настроить параметры HTTP-сервера:

- `SERVER_ADDR` - адрес сервера, по умолчанию `:8080`;
- `SERVER_HTTP_TIMEOUT` - timeout для HTTP-клиента и HTTP-сервера, по умолчанию `15s`.

Получить ключ и ознакомиться с инструкцией можно здесь:

`https://app.countrystatecity.in/`

Установить зависимости:

```bash
go mod tidy
```

Запустить приложение:

```bash
go run ./cmd/api
```

После запуска сервер должен быть доступен на `http://localhost:8080`.

Пример запроса:

```bash
curl http://localhost:8080/api/v1/weather/Almaty
```

## Endpoints

### `GET /health`

Проверка, что сервис запущен и отвечает.

### `GET /api/v1/weather/{city}`

Возвращает погоду по названию города.

Дополнительно в ответе возвращается рекомендация, что надеть в зависимости от температуры:

- холодно - теплая одежда;
- прохладно - куртка;
- тепло - легкая одежда.

### `GET /api/v1/weather/country/{country}`

Возвращает информацию о погоде по городам страны.

### `GET /api/v1/weather/country/{country}/top`

Возвращает топ-3 самых теплых городов в стране.

Дополнительно можно передать query-параметр `limit`, чтобы получить больше или меньше результатов.

## Примеры API

Ниже приведены `curl`-примеры, которые показывают возможности текущего API.

Примечание: значения температуры меняются со временем, потому что API запрашивает живые погодные данные.

### Проверка Сервиса

```bash
curl http://localhost:8080/health
```

Ожидаемый ответ:

```text
ok
```

### Погода По Городу

Запрос погоды по городу:

```bash
curl "http://localhost:8080/api/v1/weather/almaty"
```

Пример ответа:

```json
{
  "city": "Almaty",
  "conditions": {
    "temperature": 10.1,
    "feelsLike": 9.3
  },
  "recommendation": "It is cool outside, a jacket will help."
}
```

Что важно:

- можно передавать город в любом регистре, например `almaty`, `ALMATY`, `Almaty`;
- в ответе название города нормализуется.

Дополнительный пример:

```bash
curl "http://localhost:8080/api/v1/weather/Shymkent"
```

### Погода По Стране

Запрос погоды по городам страны:

```bash
curl "http://localhost:8080/api/v1/weather/country/KZ"
```

Тот же endpoint с кодом страны в нижнем регистре тоже работает:

```bash
curl "http://localhost:8080/api/v1/weather/country/kz"
```

Пример ответа:

```json
[
  {
    "city": "Abai",
    "conditions": {
      "temperature": -2.3,
      "feelsLike": -7.9
    },
    "recommendation": "It is cold outside, wear warm clothes."
  },
  {
    "city": "Almaty",
    "conditions": {
      "temperature": 11.5,
      "feelsLike": 10.9
    },
    "recommendation": "It is cool outside, a jacket will help."
  },
  {
    "city": "Shymkent",
    "conditions": {
      "temperature": 11.8,
      "feelsLike": 8.4
    },
    "recommendation": "It is cool outside, a jacket will help."
  }
]
```

`/api/v1/weather/country/KZ` возвращает полный список найденных городов страны. В примере выше показан только фрагмент ответа.

### Топ Самых Теплых Городов В Стране

Топ-3 самых теплых городов в стране:

```bash
curl "http://localhost:8080/api/v1/weather/country/KZ/top"
```

Пример ответа:

```json
[
  {
    "city": "Jambyl",
    "conditions": {
      "temperature": 15.9,
      "feelsLike": 10.3
    },
    "recommendation": "It is cool outside, a jacket will help."
  },
  {
    "city": "Atyrau",
    "conditions": {
      "temperature": 13.9,
      "feelsLike": 10.4
    },
    "recommendation": "It is cool outside, a jacket will help."
  },
  {
    "city": "Turkistan",
    "conditions": {
      "temperature": 12.5,
      "feelsLike": 8.6
    },
    "recommendation": "It is cool outside, a jacket will help."
  }
]
```

### Топ Городов С Пользовательским Limit

Можно передать query-параметр `limit`:

```bash
curl "http://localhost:8080/api/v1/weather/country/KZ/top?limit=5"
```

Пример ответа:

```json
[
  {
    "city": "Jambyl",
    "conditions": {
      "temperature": 15.9,
      "feelsLike": 10.3
    },
    "recommendation": "It is cool outside, a jacket will help."
  },
  {
    "city": "Atyrau",
    "conditions": {
      "temperature": 13.9,
      "feelsLike": 10.4
    },
    "recommendation": "It is cool outside, a jacket will help."
  },
  {
    "city": "Turkistan",
    "conditions": {
      "temperature": 12.5,
      "feelsLike": 8.6
    },
    "recommendation": "It is cool outside, a jacket will help."
  },
  {
    "city": "Shymkent",
    "conditions": {
      "temperature": 11.8,
      "feelsLike": 8.4
    },
    "recommendation": "It is cool outside, a jacket will help."
  },
  {
    "city": "Almaty",
    "conditions": {
      "temperature": 11.5,
      "feelsLike": 10.9
    },
    "recommendation": "It is cool outside, a jacket will help."
  }
]
```

### Примеры Ошибок

Невалидный `limit`:

```bash
curl "http://localhost:8080/api/v1/weather/country/KZ/top?limit=0"
```

Ответ:

```json
{
  "error": "Limit must be a positive number."
}
```

Пустой `country` path parameter не совпадает с маршрутом, поэтому такой запрос вернет `404 Not Found`:

```bash
curl "http://localhost:8080/api/v1/weather/country/"
```

## Требования

- Использовать `net/http` и `chi`.
- Все ответы возвращать в формате JSON.
- Разделить код на слои:
  - `internal/transport/http/v1` - HTTP endpoints, чтение параметров запроса, запись ответа.
  - `internal/service` - бизнес-логика, сортировка, рекомендации по одежде.
  - `internal/client` - работа с Open-Meteo.
- Добавить обработку ошибок.
- Использовать Open-Meteo для получения погодных данных.

## TODO

### 1. Подготовка проекта

- [x] Создать точку входа `cmd/api/main.go`.
- [x] Создать базовую структуру проекта: `internal/handler`, `internal/service`, `internal/client`.
- [x] Подключить `chi` и настроить HTTP-сервер.

### 2. Интеграция с Weather API

- [x] Реализовать клиент для Open-Meteo: поиск координат города и получение текущей погоды.
- [x] Обработать ошибки внешнего API.

### 3. Бизнес-логика

- [x] Добавить рекомендацию по одежде на основе температуры.
- [x] Реализовать получение погоды по городам страны.
- [x] Реализовать сортировку и выбор топ-3 самых теплых городов.

### 4. HTTP endpoints

- [x] Реализовать `GET /weather/{city}`.
- [x] Реализовать `GET /weather/country/{country}`.
- [x] Реализовать `GET /weather/country/{country}/top`.
- [x] Настроить JSON-ответы для успешных и ошибочных сценариев.

### 5. Проверка

- [x] Проверить запуск приложения через `go run ./cmd/api`.
- [x] Проверить основные endpoints через `curl`.
- [x] Проверить сценарии ошибок.
