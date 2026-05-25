# Signal Peer Management Service

Система управления пирами (peers) с ролями `master` и `slave`. Позволяет регистрировать узлы, отслеживать их статус онлайн, назначать мастеров и получать списки активных узлов. Состояние хранится в PostgreSQL, фоновая задача автоматически помечает неактивные пиры как офлайн.

## Особенности

- Регистрация пира с указанием роли (`master`/`slave`)
- Автоматическое отслеживание статуса онлайн через `heartbeat` (`PATCH /heartbeat`)
- Периодическая проверка и отключение неактивных пиров (по умолчанию каждые 30 секунд)
- Получение списков всех мастеров и слейвов
- Назначение пира мастером через API
- Работа с адресами в формате `host:port`
- Постгрес для персистентности
- Docker Compose для разворачивания БД

## Технологии

- [Go](https://go.dev/) (1.26.1)
- [Gin](https://gin-gonic.com/) — веб-фреймворк
- [pgx](https://github.com/jackc/pgx) — драйвер и пул соединений с PostgreSQL
- [golang-migrate](https://github.com/golang-migrate/migrate) — миграции БД
- PostgreSQL 16+
- Docker / Docker Compose

## Установка и запуск

### 1. Требования

- Go 1.21+
- Docker и Docker Compose

### 3. Переменные окружения

Создайте файл `.env` в корне проекта, например:

```env
# PostgreSQL для приложения
DATABASE_URL=postgres://user:password@localhost:5432/signal_db

# Docker-контейнер с PostgreSQL
DOCKER_PG_USER=user
DOCKER_PG_PASSWORD=password
DOCKER_PG_DB=signal_db
DOCKER_DB_PORT=5432

# Путь к корню проекта (для скрипта start.sh)
ROOT=/absolute/path/to/project

# Путь к исполняемому файлу migrate (если используется скрипт)
GOLANG_MIGRATE=/path/to/migrate
```

### 4. Запуск базы данных

```bash
docker-compose up -d
```

### 5. Применение миграций

```bash
migrate -path ./migrations -database "$DATABASE_URL?sslmode=disable" up
```

### 6. Запуск приложения

```bash
go run ./cmd/app
```

По умолчанию сервер запускается на `:8080` (стандартный порт Gin). При необходимости можно переопределить через `gin.Run(":порт")` — в текущей реализации `s.router.Run()` без аргументов использует `:8080`.

### 7. Альтернативный запуск через скрипт

```bash
chmod +x ./Scripts/start.sh
./Scripts/start.sh
```

Скрипт поднимает Docker, накатывает миграции, запускает приложение и после нажатия `read` откатывает миграции и останавливает контейнеры. Подходит для разработки.

## API Endpoints

Базовый URL: `http://localhost:8080`

### GET `/ping`
Проверка доступности сервера.

**Ответ (200):**
```json
{"message": "pong"}
```

### POST `/connect`
Регистрация нового пира. Адрес пира извлекается из `RemoteAddr` запроса.

**Тело:**
```json
{"role": "master"}   // или "slave"
```

**Ответ (201):**
```json
{"status": "ok", "addr": "127.0.0.1:54321"}
```

### DELETE `/disconnect`
Удаление пира (по `RemoteAddr` запроса).

**Ответ (200):** *пустое тело*

### PATCH `/heartbeat`
Обновление статуса пира на онлайн (и сброс таймера неактивности). Использует `RemoteAddr` для идентификации.

**Ответ (200):**
```json
{"message": "ok"}
```

### GET `/get_masters`
Возвращает список адресов всех мастеров.

**Ответ (200):**
```json
{"message": "successful", "master_addrs": ["10.0.0.1:8080", "10.0.0.2:8081"]}
```

### GET `/get_slaves`
Возвращает список адресов всех слейвов.

**Ответ (200):**
```json
{"message": "successful", "slave_addrs": ["10.0.0.3:8082", "10.0.0.4:8083"]}
```

### POST `/set_master`
Назначить существующему пиру роль `master`.

**Тело:**
```json
{"addr_port": "192.168.1.10:9090"}
```

**Ответ (200):**
```json
{"message": "ok"}
```

## Примеры запросов (cURL)

```bash
# Регистрация мастера
curl -X POST http://localhost:8080/connect \
  -H "Content-Type: application/json" \
  -d '{"role":"master"}'

# Heartbeat
curl -X PATCH http://localhost:8080/heartbeat

# Получить мастеров
curl -X GET http://localhost:8080/get_masters

# Назначить мастером другого пира
curl -X POST http://localhost:8080/set_master \
  -H "Content-Type: application/json" \
  -d '{"addr_port":"192.168.1.5:8080"}'

# Удалить себя
curl -X DELETE http://localhost:8080/disconnect
```

## Структура проекта

```
.
├── cmd/
│   └── app/                # точка входа main.go
├── internal/
│   ├── api/                # HTTP-обработчики и сервер
│   │   └── server.go
│   ├── domain/             # модели и доменные типы
│   │   └── peer.go
│   ├── repository/         # работа с БД (pgx)
│   │   └── postgres.go
│   └── service/            # бизнес-логика
│       └── peer.go
├── migrations/             # SQL-миграции
│   ├── 000001_init.up.sql
│   └── 000001_init.down.sql
├── Scripts/
│   └── start.sh            # вспомогательный скрипт запуска
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

## Примечания. Некоторые из пунктов будут исправлены в следующей версии.

- **Идентификация пиров**: все запросы (`/heartbeat`, `/disconnect`) используют `c.Request.RemoteAddr`, что в среде с прокси может давать неверный IP. Рекомендуется передавать адрес в явном виде через заголовок или тело запроса.
- **Graceful shutdown**: в приложении отсутствует обработка сигналов ОС для корректного закрытия соединений с БД. При необходимости доработать.
- **Роль мастера**: логика множественных мастеров не запрещена – система допускает любое количество мастеров. При необходимости добавьте ограничения в сервисном слое.
- **Данный сервер является учебным!**: Сервер не является готовым к использованию в профессиональный среде. Я написал его лишь для того, чтобы развивать мой проект *Pigeon*, который нуждается в сигнальном сервере.

## Лицензия
Apache-2.0
