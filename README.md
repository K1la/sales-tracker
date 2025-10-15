## Sales Tracker

Простой сервис для учёта финансовых записей (доходы/расходы) с CRUD-операциями и аналитикой по данным: сумма, среднее, количество, медиана и 90-й перцентиль. Данные хранятся в PostgreSQL, агрегаты считаются SQL-запросами. В комплекте — минимальный фронтенд (HTML + JS) для добавления/просмотра записей, фильтрации/сортировки, аналитики и экспорта в CSV.

Фронтенд доступен по корню `/` (сервируется приложением), API — под префиксом `/api`.

---

## Доступные API endpoints

Базовый префикс: `/api`

- POST `/api/items` — создать запись
  - Тело: `{ "type": "income|expense", "amount": number, "date": "YYYY-MM-DD", "category": string, "description"?: string }`
  - Ответ: созданный объект
- GET `/api/items` — список записей
  - Query: `from=YYYY-MM-DD`, `to=YYYY-MM-DD`, `sort_by=<field> <ASC|DESC>` (по умолчанию `date DESC`)
- GET `/api/items/:id` — запись по ID
- PUT `/api/items/:id` — обновить запись (частично)
  - Тело: подмножество `{ type, amount, date, category, description }`
- DELETE `/api/items/:id` — удалить запись

Аналитика:
- GET `/api/analytics?from=YYYY-MM-DD&to=YYYY-MM-DD`
  - Ответ: `{ "sum": number, "average": number, "count": number, "median": number, "percentile": number }`

---

## Запуск через Docker Compose

Требования: установлен Docker и Docker Compose.

1) Создайте `.env` в корне проекта (пример ниже).
2) Запустите сборку и контейнеры:

```bash
docker compose build
docker compose up -d
```

Проверка логов бэкенда:
```bash
docker compose logs -f backend
```

Приложение: `http://localhost:8080/`

PostgreSQL: `localhost:5432`

Миграции применяются автоматически контейнером `migrator` (goose) при старте.

### Пример `.env`

```dotenv
# База данных
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=sales-tracker

# Goose (миграции)
GOOSE_DRIVER=postgres
GOOSE_MIGRATION_DIR=/migrations
```

Примечания:
- Основные значения для Postgres и HTTP берутся из `env/config.yaml`, пароль БД — из переменной окружения `DB_PASSWORD` (см. `internal/config/config.go`).
- Строка подключения для мигратора собирается из переменных `.env` (см. `docker-compose.yml`).

Остановка:
```bash
docker compose down
# Удалить volume с данными Postgres при необходимости
docker compose down -v
```

---

## Архитектура проекта

- `cmd/sales-tracker/` — точка входа (main)
- `internal/api/handlers/` — HTTP-обработчики (handler)
  - `items/` — CRUD по записям
  - `analytics/` — расчёт агрегатов
- `internal/api/router/` — маршрутизация (router), выдача статических файлов из `web/`
- `internal/api/server/` — HTTP-сервер (server)
- `internal/service/` — слой бизнес-логики (service)
  - `items/`
  - `analytics/`
- `internal/repository/` — работа с базой данных PostgreSQL (postgres / repository, SQL)
  - `db.go` — подключение/утилиты
  - `items/` — репозиторий записей
  - `analytics/` — репозиторий аналитики
- `internal/dto/` — DTO: структуры запросов/ответов API (dto)
- `internal/model/` — доменные/БД-модели (model)
- `internal/config/` — конфигурация (yaml + env)
- `migrations/` — SQL-миграции (goose)
- `env/config.yaml` — конфиг HTTP/DB по умолчанию
- `web/` — фронтенд (статические файлы: `index.html`, `styles.css`, `app.js`)
- `Dockerfile`, `docker-compose.yml` — контейнеризация
- `go.mod`, `go.sum` — зависимости

Слои по назначению:
- **model** — модели БД/домена
- **dto** — контракты API (вход/выход)
- **handler** — HTTP-эндпоинты и валидация
- **service** — бизнес-правила
- **postgres (repository)** — SQL и доступ к БД
- **router** — маршруты и статика
- **server** — запуск HTTP-сервера

### Дерево проекта

```text
sales-tracker/
├── cmd/
│   └── sales-tracker/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   ├── items/
│   │   │   │   ├── handler.go
│   │   │   │   └── methods.go
│   │   │   └── analytics/
│   │   │       ├── handler.go
│   │   │       └── methods.go
│   │   ├── router/
│   │   │   └── router.go
│   │   └── server/
│   │       └── server.go
│   ├── config/
│   │   ├── config.go
│   │   └── types.go
│   ├── dto/
│   │   └── dto.go
│   ├── model/
│   │   └── model.go
│   ├── repository/
│   │   ├── db.go
│   │   ├── items/
│   │   │   ├── repo.go
│   │   │   └── methods.go
│   │   └── analytics/
│   │       ├── repo.go
│   │       └── methods.go
│   └── service/
│       ├── items/
│       │   ├── service.go
│       │   └── methods.go
│       └── analytics/
│           ├── service.go
│           └── methods.go
├── migrations/
│   └── 20251014195537_items_table.sql
├── env/
│   └── config.yaml
├── web/
│   ├── index.html
│   ├── styles.css
│   └── app.js
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── go.sum
```

---

## Примеры запросов

### Создание записи

Запрос:
```bash
curl -X POST http://localhost:8080/api/items \
  -H "Content-Type: application/json" \
  -d '{
    "type": "income",
    "amount": 1234.56,
    "date": "2025-10-15",
    "category": "salary",
    "description": "October"
  }'
```

Ответ 200 OK:
```json
{
  "id": 1,
  "type": "income",
  "amount": 1234.56,
  "date": "2025-10-15",
  "category": "salary",
  "description": "October",
  "created_at": "2025-10-15T10:00:00Z",
  "updated_at": "2025-10-15T10:00:00Z"
}
```

### Аналитика за период

Запрос:
```bash
curl "http://localhost:8080/api/analytics?from=2025-10-01&to=2025-10-31"
```

Ответ 200 OK:
```json
{
  "sum": 12345.67,
  "average": 4115.22,
  "count": 3,
  "median": 4000,
  "percentile": 9000
}
```
