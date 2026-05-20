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
- [ ] Gateway Service вызывает API Service через внутренний адрес Docker Compose.
- [ ] API Service вызывает Gateway.
