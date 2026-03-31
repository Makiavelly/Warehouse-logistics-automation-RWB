# Warehouse Logistics Automation — RWB

Система автоматического вызова транспорта на склады на основе прогноза отгрузок.

## Архитектура

```
┌──────────────┐   REST/JSON    ┌──────────────────────┐   gRPC/protobuf  ┌──────────────────┐
│   Frontend   │ ─────────────► │  logistics-service   │ ───────────────► │   ml-service     │
└──────────────┘                │  (Go / port 8080)    │                  │ (Python / :50051)│
                                └──────────┬───────────┘                  └──────────────────┘
┌──────────────┐   REST + API-Key        │
│ External API │ ───────────────────────►│ /api/data/ingest
│ (статусы)    │                         │ SQL
└──────────────┘              ┌──────────▼──────────┐
                              │      PostgreSQL      │
                              └─────────────────────┘
```

### Компоненты
- **logistics-service** — Go backend: бизнес-логика, REST API, автоматический вызов транспорта
- **ml-service** — Python сервис: прогнозная модель, gRPC сервер на порту `50051` (разрабатывается отдельно)
- **frontend** — веб-интерфейс (разрабатывается отдельно)
- **PostgreSQL** — основная БД

устанавливаем admin:
```docker exec postgres-logistics psql -U postgres -d logistics -c "INSERT INTO users (username, password_hash, role) VALUES ('admin', crypt('admin123', gen_salt('bf')), 'admin');"```
Подключитесь к БД: docker exec -it postgres-logistics psql -U postgres -d logistics

В Postman посмотреть токен админа:
POST http://localhost:8080/api/auth/login и в body поставьте {"username":"admin","password":"admin123"}
В ответе увидите токен админа

Сохраните токен в Authorization Bear Token (Это в Postman) 
Открыть BD: http://localhost:18888

# из папки backend/
cd D:\Warehouse-logistics-automation-RWB\backend

# Загрузить обучающие данные (с target_2h)
go run ./cmd/seed -file ../test_team_track.parquet

# Загрузить тестовые данные
go run ./cmd/seed -file ../test_team_track.parquet

### Почему gRPC между backend и ML-сервисом

| | HTTP/REST | gRPC |
|---|---|---|
| Протокол | JSON (текст) | Protobuf (бинарный, ~5× компактнее) |
| Контракт | Нет гарантий | `.proto` файл — единый источник правды для Go и Python |
| Типизация | Ручная десериализация | Кодогенерация на обеих сторонах |
| Streaming | Нет | Поддерживается (для batch-прогнозов в будущем) |

**Proto-контракт:** `proto/logistics/logistics.proto`

```protobuf
service MLService {
  rpc Ping(Empty)              returns (Empty);
  rpc Predict(PredictRequest)  returns (PredictReply);    // прогноз отгрузок
  rpc Retrain(RetrainRequest)  returns (RetrainReply);    // запуск дообучения
  rpc RetrainStatus(Empty)     returns (RetrainStatusReply); // статус обучения
}
```

Переключение HTTP ↔ gRPC через env: `ML_USE_GRPC=true` (gRPC, по умолчанию) / `ML_USE_GRPC=false` (HTTP fallback).

## Роли и возможности

### Admin (менеджер)
| Действие | Метод | URL |
|---|---|---|
| Войти в систему | `POST` | `/api/auth/login` |
| Добавить пользователя/водителя | `POST` | `/api/auth/register` |
| Добавить склад | `POST` | `/api/warehouses` |
| Список складов | `GET` | `/api/warehouses` |
| Удалить склад | `DELETE` | `/api/warehouses/{warehouse_id}` |
| Добавить маршрут к складу | `POST` | `/api/warehouses/{warehouse_id}/routes` |
| Список маршрутов склада | `GET` | `/api/warehouses/{warehouse_id}/routes` |
| Удалить маршрут | `DELETE` | `/api/warehouses/{warehouse_id}/routes/{route_id}` |
| Установить порог вызова для (склад, маршрут) | `PUT` | `/api/thresholds` |
| Список порогов | `GET` | `/api/thresholds?warehouse_id=&route_id=` |
| Запросить прогноз отгрузок | `POST` | `/api/forecasts/predict` |
| История прогнозов | `GET` | `/api/forecasts?warehouse_id=&route_id=&from=&to=` |
| История вызовов машин | `GET` | `/api/truck-calls?warehouse_id=&route_id=` |
| Точность вызовов (аналитика) | `GET` | `/api/truck-calls/accuracy?warehouse_id=&route_id=` |
| Список водителей | `GET` | `/api/drivers` |
| Назначить водителя на маршрут | `PUT` | `/api/drivers/assign` |
| Дообучить модель | `POST` | `/api/model/retrain` |

### Driver (водитель)
| Действие | Метод | URL |
|---|---|---|
| Получить сигнал вызова на погрузку | `GET` | `/api/driver/signal` |
| Отметить своевременность вызова | `POST` | `/api/driver/truck-calls/{truck_call_id}/timeliness` |
| Статистика по своим вызовам | `GET` | `/api/driver/stats` |

### External integration (внешний сервис)
| Действие | Метод | URL | Auth |
|---|---|---|---|
| Загрузить новые данные для дообучения | `POST` | `/api/data/ingest` | `X-API-Key` header |

## Логика работы системы

1. Менеджер создаёт склады и маршруты, задаёт пороговые значения для каждой пары `(склад, маршрут)`.
2. По расписанию (или вручную) запрашивается прогноз у ML-сервиса для нужной пары и горизонта планирования (по умолчанию 2 часа).
3. Если `predicted_count >= threshold.value` — система автоматически создаёт запись `truck_call` со статусом `pending`.
4. Водитель, назначенный на эту пару, видит сигнал через `GET /api/driver/signal`.
5. После выполнения водитель отмечает своевременность (`on_time` / `late` / `early`) и фактическое количество ёмкостей.
6. Менеджер видит аналитику по точности вызовов: % своевременных, среднее прогноз vs факт.

## Бизнес-допущения

- Все машины одинакового объёма (единица — «ёмкость»/контейнер).
- Горизонт прогноза — 2 часа (как в train/test данных `target_2h`).
- Один водитель назначен на одну пару `(склад, маршрут)`.
- Сигнал водителю — pull-модель (водитель опрашивает `/signal`); push-уведомления реализуются на стороне frontend/mobile.
- ML-сервис реализует gRPC-сервер (`MLService`) на порту `50051`; при `ML_USE_GRPC=false` ожидается HTTP на `POST /predict` и `POST /retrain`.

## Метрики качества системы

- **WAPE + |Relative Bias|** — основная ML-метрика (соответствует условиям хакатона).
- **Accuracy Rate (%)** — доля своевременных вызовов машин (`on_time / total`).
- **ΔForecast** — среднее отклонение прогноза от фактического количества ёмкостей.
- **Missed calls** — вызовы со статусом `missed` (машина не приехала).

## Запуск

### Требования
- Docker + Docker Compose

### Быстрый старт

```bash
# Клонировать репозиторий
git clone <repo-url>
cd Warehouse-logistics-automation-RWB

# Запустить backend + PostgreSQL
docker compose up --build -d

# Проверить работу
curl http://localhost:8080/health
# {"status":"ok"}
```

Frontend доступен по адресу `http://localhost:3000`.

### Переменные окружения (backend)

| Переменная | По умолчанию | Описание |
|---|---|---|
| `DB_ADDRESS` | `postgres://postgres:password@postgres:5432/logistics` | Строка подключения к PostgreSQL |
| `AUTH_SECRET_KEY` | `change-me` | JWT secret (менять в prod!) |
| `AUTH_API_KEY` | `internal-api-key` | API key для внешних интеграций |
| `ML_ADDRESS` | `ml-service:50051` | Адрес ML сервиса (gRPC host:port) |
| `ML_USE_GRPC` | `true` | `true` — gRPC клиент, `false` — HTTP fallback |
| `APP_LOG_LEVEL` | `info` | Уровень логирования (`debug`/`info`/`warn`/`error`) |

### Регенерация gRPC кода (при изменении proto)

```bash
cd backend

# Установить плагины (один раз)
make proto-install

# Регенерировать pb.go из .proto файла
make proto
```

> Сгенерированные файлы (`*.pb.go`) коммитятся в репозиторий — разработчикам не нужно устанавливать `protoc` для сборки.

### Локальная разработка (без Docker)

```bash
cd backend

# Запустить только postgres
docker compose up postgres -d

# Собрать и запустить
go build -o main ./logistics-service/main.go
./main -config logistics-service/config.yml
```

## Структура репозитория

```
proto/
  logistics/
    logistics.proto           # gRPC контракт backend ↔ ml-service
    logistics.pb.go           # Сгенерировано protoc (не редактировать)
    logistics_grpc.pb.go      # Сгенерировано protoc (не редактировать)

backend/
  proto/logistics/            # Копия сгенерированных pb.go для сборки Go модуля
  logistics-service/
    main.go                   # Точка входа, регистрация роутов
    config.yml                # Конфигурация по умолчанию
    config/                   # Структуры конфигурации
    core/
      models/                 # Доменные модели
      ports/                  # Интерфейсы (ports/adapters pattern)
      service/                # Бизнес-логика
      errors/                 # Ошибки домена
    adapters/
      db/                     # PostgreSQL + SQL-миграции
        migrations/           # .up.sql / .down.sql
      rest/                   # HTTP handlers (один файл на группу)
      JWT/                    # JWT генерация и валидация
      middleware/             # Role-based auth + API-key auth + logging
      ml_client/              # HTTP fallback клиент к ML сервису
      ml_client_grpc/         # gRPC клиент к ML сервису (основной)
      http_server/            # HTTP сервер с graceful shutdown
      logger/                 # slog обёртка
      mux/                    # HTTP multiplexer
      validator/              # Валидация запросов (go-playground)
  Dockerfile
  Makefile
  compose.yml

frontend/                     # (разрабатывается отдельно)
ml-service/                   # (разрабатывается отдельно)
compose.yaml                  # Root compose (весь стек)
```

## API — примеры запросов

```bash
# Логин (получить JWT)
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"secret"}'

# Создать склад (admin JWT в заголовке)
curl -X POST http://localhost:8080/api/warehouses \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Коледино","office_from_id":"WH-001","address":"Московская обл."}'

# Установить порог для пары (склад, маршрут)
curl -X PUT http://localhost:8080/api/thresholds \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"warehouse_id":"<uuid>","route_id":"<uuid>","value":50}'

# Запросить прогноз
curl -X POST http://localhost:8080/api/forecasts/predict \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"warehouse_id":"<uuid>","route_id":"<uuid>","forecast_time":"2026-04-01T10:00:00Z","horizon_hours":2}'

# Загрузить данные для дообучения (API key)
curl -X POST http://localhost:8080/api/data/ingest \
  -H "X-API-Key: internal-api-key" \
  -H "Content-Type: application/json" \
  -d '{"data_points":[{"route_id":"r1","office_from_id":"WH-001","timestamp":"2026-03-01T12:00:00Z","status_1":10,"target_2h":45}]}'
```



Шаг 0 — Запуск

# В корне проекта D:\Warehouse-logistics-automation-RWB\
make up
# или напрямую:
docker compose up --build -d

Что поднимается:
- localhost:8080 — Go API (backend)
- localhost:5432 — PostgreSQL
- localhost:18888 — pgAdmin (визуальный просмотр БД)

  ---
Шаг 1 — Авторизация в Postman

Логин

POST http://localhost:8080/api/auth/login
Content-Type: application/json

{
"username": "admin",
"password": "admin123"
}
В ответ получишь token. Скопируй его — он нужен для всех следующих запросов.

Как использовать токен

В каждом следующем запросе добавляй заголовок:
Authorization: Bearer <вставь_токен_сюда>

В Postman удобно: вкладка Authorization → Bearer Token → вставить токен.

  ---
Шаг 2 — Загрузка данных из parquet

cd D:\Warehouse-logistics-automation-RWB\backend

# Сначала обучающие данные (содержат target_2h — нужны ML-модели)
go run ./cmd/seed -file ../train_team_track.parquet

# Потом тестовые
go run ./cmd/seed -file ../test_team_track.parquet

После этого данные попадают в таблицу raw_data в PostgreSQL.

  ---
Шаг 3 — Узнать какие склады и маршруты есть в данных

# Посмотреть уникальные склады из загруженных данных
docker exec postgres-logistics psql -U postgres -d logistics -c \
"SELECT DISTINCT office_from_id FROM raw_data ORDER BY office_from_id;"

# Посмотреть уникальные маршруты
docker exec postgres-logistics psql -U postgres -d logistics -c \
"SELECT DISTINCT route_id, office_from_id FROM raw_data ORDER BY office_from_id, route_id;"

  ---
Шаг 4 — Создать склады через API

Для каждого office_from_id из данных создаёшь склад:

POST http://localhost:8080/api/warehouses
Authorization: Bearer <token>
Content-Type: application/json

{
"name": "Склад 1",
"office_from_id": "<значение из данных>",
"address": ""
}

Получишь id склада — сохрани, он нужен для маршрутов.

Посмотреть все склады:
GET http://localhost:8080/api/warehouses
Authorization: Bearer <token>

  ---
Шаг 5 — Создать маршруты

Для каждого route_id привязываешь его к складу:

POST http://localhost:8080/api/warehouses/<warehouse_id>/routes
Authorization: Bearer <token>
Content-Type: application/json

{
"route_id": "<значение из данных>",
"name": "Маршрут 1"
}

  ---
Шаг 6 — Посмотреть данные в pgAdmin

Открой http://localhost:18888

- Email: admin@logistics.local
- Password: password

Подключение к серверу:
- Host: postgres
- Port: 5432
- Database: logistics
- Username: postgres
- Password: password

Там видно все таблицы: raw_data, warehouses, routes, users и т.д.

  ---
Шаг 7 — Установить пороги (thresholds)

Порог — при каком прогнозируемом количестве вызывать грузовик:

PUT http://localhost:8080/api/thresholds
Authorization: Bearer <token>
Content-Type: application/json

{
"warehouse_id": "<id склада>",
"route_id": "<id маршрута>",
"value": 10.0
}

  ---
Шаг 8 — Прогноз (нужен ML-сервис)

POST http://localhost:8080/api/forecasts/predict
Authorization: Bearer <token>
Content-Type: application/json

{
"warehouse_id": "<id склада>",
"route_id": "<id маршрута>",
"forecast_time": "2024-01-15T10:00:00Z",
"horizon_hours": 2
}

▎ ⚠️ Это требует работающего ML-сервиса (ml-service). Сейчас он не запущен — это следующий шаг разработки.

  ---
Что сейчас работает без ML

┌──────────────────────────────────────┬─────────────────────────┐
│               Endpoint               │       Что делает        │
├──────────────────────────────────────┼─────────────────────────┤
│ POST /api/auth/login                 │ Логин                   │
├──────────────────────────────────────┼─────────────────────────┤
│ GET/POST /api/warehouses             │ Склады                  │
├──────────────────────────────────────┼─────────────────────────┤
│ GET/POST /api/warehouses/{id}/routes │ Маршруты                │
├──────────────────────────────────────┼─────────────────────────┤
│ GET/PUT /api/thresholds              │ Пороги                  │
├──────────────────────────────────────┼─────────────────────────┤
│ POST /api/data/ingest                │ Загрузка данных         │
├──────────────────────────────────────┼─────────────────────────┤
│ GET /health                          │ Проверка что сервер жив │
└──────────────────────────────────────┴─────────────────────────┘
