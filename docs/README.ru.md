# swag2mcp

**swag2mcp** — это CLI-инструмент и MCP (Model Context Protocol) сервер, который связывает OpenAPI/Swagger/Postman спецификации с LLM-агентами (Opencode, Crush, Copilot, Cursor и другими).

Он индексирует ваши API-спецификации в полнотекстовый поисковый движок, предоставляет 14 MCP-инструментов и позволяет LLM находить, изучать и вызывать реальные API-эндпоинты — без единой строки интеграционного кода.

---

## Содержание

- [Быстрый старт](#быстрый-старт)
- [Конфигурация](#конфигурация)
- [CLI Команды](#cli-команды)
- [MCP Сервер](#mcp-сервер)
- [Поиск](#поиск)
- [Рабочая директория (Workspace)](#рабочая-директория-workspace)
- [Кэширование](#кэширование)
- [Разработка](#разработка)

---

## Быстрый старт

```bash
# Установка
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest

# Инициализация рабочей директории
swag2mcp init

# Запуск MCP сервера (для LLM-агентов)
swag2mcp mcp

# Или интерактивный проводник
swag2mcp run
```

---

## Конфигурация

### YAML Схема

```yaml
http_client:                        # опционально, глобальные настройки HTTP
  headers:                          # опционально
    X-API-Version: "2"
  cookies: []                       # опционально
  user_agent: ""                    # опционально
  timeout: 0s                       # опционально
  follow_redirects: true            # опционально
  max_redirects: 10                 # опционально
  max_response_size: 1048           # опционально, байт (по умолч. 1KB, макс 1MB)

specs:
  - domain: petstore                    # обязательно, 1-60 символов, [a-zA-Z0-9_-]
    llm_title: Petstore API             # обязательно, 5-120 символов
    llm_instruction: |                  # опционально, макс 500 символов
      Используй это API для управления питомцами, заказами и пользователями.
    base_url: https://petstore.swagger.io/v2  # обязательно, валидный URL
    disable: false                      # опционально
    tags: [public, demo]                # опционально, для фильтрации
    http_client:                        # опционально, переопределяет глобальный
      headers:
        X-API-Version: "2"
    auth:                               # опционально
      type: bearer                      # см. Методы аутентификации
      config:
        token: $(TOKEN_AUTH)
    collections:
      - llm_title: Petstore Swagger     # опционально, макс 120 символов
        llm_instruction: |             # опционально, макс 360 символов
          Основные эндпоинты Petstore
        title: ""                      # опционально, авто-заполняется из spec
        location: https://petstore.swagger.io/v2/swagger.json  # обязательно, 5-250 символов
        disable: false                  # опционально
        base_url: ""                    # опционально, переопределяет base_url спецификации
        http_client: {}                 # опционально, переопределяет spec
```

### Теги — фильтрация спецификаций по проектам

Теги позволяют группировать спецификации по проектам, окружениям или командам. При запуске MCP сервера используйте `--tags` для загрузки только нужных спецификаций:

```bash
# Запуск сервера только с публичными спецификациями
swag2mcp mcp --tags=public

# Запуск с несколькими тегами
swag2mcp mcp --tags=public,internal

# Запуск нескольких серверов для разных проектов
swag2mcp mcp --tags=project-alpha --logfile=/tmp/swag2mcp-alpha.log
swag2mcp mcp --tags=project-beta  --logfile=/tmp/swag2mcp-beta.log
```

Это позволяет запускать отдельные MCP серверы для разных проектов из одного конфигурационного файла.

### Методы аутентификации

| Тип | Поля | Пример конфига |
|-----|------|----------------|
| `none` | — | `type: none` |
| `basic` | `username`, `password` | `username: $(USER)`, `password: $(PASS)` |
| `bearer` | `token` | `token: $(TOKEN)` |
| `digest` | `username`, `password` | `username: admin`, `password: secret` |
| `api-key` | `key`, `value`, `in` (header/query) | `key: X-API-Key`, `value: $(KEY)`, `in: header` |
| `oauth2-cc` | `client_id`, `client_secret`, `token_url`, `scopes` | `client_id: $(ID)`, `token_url: https://auth.example.com/token` |
| `oauth2-pwd` | `username`, `password`, `client_id`, `client_secret`, `token_url`, `scopes` | `username: $(USER)`, `token_url: https://auth.example.com/token` |
| `script` | `source` | `source: path/to/auth.sh` |

Все строковые поля поддерживают синтаксис `$(ENV_VAR)` — разрешается в рантайме из переменных окружения.

---

## CLI Команды

Все команды, принимающие `[path]`, используют одинаковое разрешение пути:

```
swag2mcp <command>          → ~/.swag2mcp/
swag2mcp <command> ./       → {cwd}/.swag2mcp/
swag2mcp <command> path/to  → {cwd}/path/to/.swag2mcp/
```

### `init [path]`

Инициализация рабочей директории и конфигурации.

| Флаг | Сокращение | По умолчанию | Описание |
|------|------------|--------------|----------|
| `--interactive` | `-i` | `false` | Интерактивный мастер |
| `--force` | `-f` | `false` | Перезаписать существующий конфиг |

```bash
swag2mcp init              # создать ~/.swag2mcp/swag2mcp.yaml
swag2mcp init ./           # создать ./.swag2mcp/swag2mcp.yaml
swag2mcp init -i           # интерактивный мастер
```

### `add spec [path]` / `add collection [path]`

Добавить спецификацию или коллекцию в конфиг.

| Флаг | Сокращение | По умолчанию | Описание |
|------|------------|--------------|----------|
| `--yaml` | `-y` | `""` | YAML ввод (используйте `-` для stdin) |
| `--example` | `-e` | `false` | Показать пример YAML |

```bash
swag2mcp add spec
swag2mcp add spec --yaml 'domain: petstore\nllm_title: Petstore API\nbase_url: https://...'
cat spec.yaml | swag2mcp add spec --yaml -
swag2mcp add spec --example
```

### `delete spec [path]` / `delete collection [path]`

Удалить спецификацию или коллекцию. Интерактивный выбор.

```bash
swag2mcp delete spec
swag2mcp delete collection
```

### `ls [path]`

Список спецификаций и коллекций.

| Флаг | Сокращение | По умолчанию | Описание |
|------|------------|--------------|----------|
| `--tags` | `-t` | `""` | Фильтр по тегам (через запятую) |

```bash
swag2mcp ls
swag2mcp ls --tags=public,internal
```

### `run [path]`

Интерактивный проводник API (TUI). Поиск, просмотр, изучение и сохранение эндпоинтов.

```bash
swag2mcp run
```

### `validate [path]`

Валидация конфигурации и проверка доступности всех коллекций.

| Флаг | Сокращение | По умолчанию | Описание |
|------|------------|--------------|----------|
| `--tags` | `-t` | `""` | Фильтр спецификаций по тегам |

```bash
swag2mcp validate
swag2mcp validate --tags=public
```

### `clean [path]`

Удалить всё содержимое директорий `cache/` и `responses/`.

```bash
swag2mcp clean
```

### `update [path]`

Валидация конфига, очистка кэша, перекэширование всех spec-файлов.

```bash
swag2mcp update
```

### `mcp [path]`

Запуск MCP сервера в headless-режиме (stdio транспорт). Основная production-команда для интеграции с LLM.

| Флаг | Сокращение | По умолчанию | Описание |
|------|------------|--------------|----------|
| `--logfile` | `-f` | `""` | Путь к лог-файлу |
| `--tags` | `-t` | `""` | Фильтр спецификаций по тегам |
| `--disable-llm-auth` | | `true` | `true` — аутентификация под капотом (LLM не видит токены). `false` — LLM может запрашивать токены через инструмент `auth` |
| `--dump-dir` | | `""` | Директория для дампа HTTP запросов (отладка) |

```bash
swag2mcp mcp
swag2mcp mcp --tags=public --logfile=/var/log/swag2mcp.log
swag2mcp mcp --disable-llm-auth=false
swag2mcp mcp --dump-dir=/tmp/dump
```

---

## MCP Сервер

MCP сервер предоставляет 14 инструментов через stdio транспорт. LLM-агенты (Opencode, Crush, Copilot, Cursor и др.) подключаются автоматически после настройки.

### Иерархия инструментов

```
spec_list                       — список всех доступных спецификаций
  └─ spec_by_id                 — детали спецификации по ID
       └─ collection_by_spec    — коллекции в спецификации
            └─ tag_by_collection     — теги в коллекции
                 └─ endpoint_by_tag  — эндпоинты в теге
                      └─ inspect          — полная OpenAPI операция
                           └─ invoke       — выполнение API вызова

search                          — полнотекстовый поиск по всем эндпоинтам
```

### Справочник инструментов

| Инструмент | Аргументы | Возвращает | Описание |
|------------|-----------|------------|----------|
| `spec_list` | — | `Spec[]` | Все доступные спецификации |
| `spec_by_id` | `id` | Spec + Collections | Детали спецификации |
| `collection_by_spec` | `specId` | Collections | Коллекции в спецификации |
| `collection_by_id` | `id` | Collection + Tags | Детали коллекции |
| `tag_by_collection` | `collectionId` | Tags | Теги в коллекции |
| `tag_by_spec` | `specId` | Tags | Все теги спецификации |
| `tag_by_id` | `id` | Tag | Метаданные тега |
| `endpoint_by_tag` | `tagId` | Endpoints | Эндпоинты в теге |
| `endpoint_by_collection` | `collectionId` | Endpoints | Все эндпоинты коллекции |
| `endpoint_by_spec` | `specId` | Endpoints | Все эндпоинты спецификации |
| `endpoint_by_id` | `id` | Endpoint | Краткая сводка эндпоинта |
| `search` | `query`, `limit` | Endpoints | Полнотекстовый поиск |
| `inspect` | `endpointId` | Full Operation | Полный объект OpenAPI операции |
| `invoke` | `endpointId`, `parameters`, `requestBody` | Response | Выполнение реального API вызова |
| `auth` | `specId` | Token | Получение auth токена для спецификации |

---

## Поиск

### Синтаксис запросов

| Возможность | Синтаксис | Пример |
|-------------|-----------|--------|
| Термин | `термин` | `питомцы` |
| Фраза | `"фраза"` | `"добавить питомца"` |
| Поле: method | `method:термин` | `method:post` |
| Поле: tag | `tag:термин` | `tag:auth` |
| Поле: path | `path:термин` | `path:/users` |
| Поле: summary | `summary:термин` | `summary:login` |
| Обязательно (AND) | `+термин` | `+method:post +tag:user` |
| Исключить (NOT) | `-термин` | `-deprecated` |
| Wildcard | `*` | `path:*/v2/*` |
| Нечёткий поиск | `термин~` | `watex~` |
| Регулярка | `/паттерн/` | `/user(s\|sessions)/` |
| Повышение веса | `термин^N` | `tag:pet^5` |
| Всё подряд | `*` | `*` |

### Примеры

```
# Найти POST эндпоинты в теге auth
+method:post +tag:auth

# Поиск эндпоинтов, связанных с логином
summary:"login"~

# Найти все пути с пользователями, исключить устаревшие
path:*/users/* -deprecated

# Сложный запрос
+method:get +tag:pet summary:"find by status"
```

### Индексируемые поля

| Поле | Тип | Содержимое |
|------|-----|------------|
| `method` | text | HTTP метод (в нижнем регистре) |
| `tag` | text | Имя тега (в нижнем регистре) |
| `path` | text | Путь API (в нижнем регистре) |
| `summary` | text (анализируемый) | Описание эндпоинта (в нижнем регистре) |
| `_all` | text (анализируемый) | method + path + tag + summary |

---

## Рабочая директория (Workspace)

### Структура директорий

```
~/.swag2mcp/                    # или {project}/.swag2mcp/
├── swag2mcp.yaml               # Файл конфигурации
├── cache/                      # Кэш удалённых спецификаций
│   ├── {hash}.spec             # Содержимое spec-файла
│   └── {hash}.meta             # JSON метаданные
├── specs/                      # Локальные spec-файлы (управляются пользователем)
├── responses/                  # Файлы ответов на вызовы
└── auth_scripts/               # Скрипты аутентификации
```

### Разрешение пути

```
swag2mcp <command>          → ~/.swag2mcp/
swag2mcp <command> ./       → {cwd}/.swag2mcp/
swag2mcp <command> path/to  → {cwd}/path/to/.swag2mcp/
```

### .gitignore

В `.gitignore` должны попадать только временные данные:

```
.swag2mcp/cache/*
.swag2mcp/responses/*
```

Конфиг `.swag2mcp/swag2mcp.yaml` и spec-файлы в `.swag2mcp/specs/` **должны быть в репозитории**.

### Рекомендация

Все spec-файлы храните в `.swag2mcp/specs/` — это единственный способ гарантировать, что они не будут скопированы в кэш и будут использоваться напрямую.

---

### Правила

| Источник | Поведение |
|----------|-----------|
| HTTP/HTTPS URL | Всегда кэшируется. TTL: случайный 1-48ч. |
| Локальный путь внутри `specs/` | Используется напрямую, не кэшируется. |
| Локальный путь вне `specs/` | Копируется в кэш при первом доступе. |
| `file://` URL | Обрабатывается как локальный путь. |

### Ключ кэша

SHA-256 хеш нормализованного location (первые 16 байт = 32 hex символа).

### Логика попадания в кэш

1. Чтение `.meta` файла — истёк или отсутствует → промах
2. Для локальных источников: `ModTime` изменился → промах
3. `.spec` файл отсутствует → промах
4. Иначе → попадание

---

## Разработка

```bash
# Сборка
go build ./cmd/swag2mcp/

# Тесты
go test ./...

# Линтер
make lint

# Запуск
go run ./cmd/swag2mcp/main.go
```
