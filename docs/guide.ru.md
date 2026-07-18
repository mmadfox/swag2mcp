# swag2mcp

**swag2mcp** — это CLI-инструмент и MCP (Model Context Protocol) сервер, который связывает OpenAPI/Swagger/Postman спецификации с LLM-агентами (Opencode, Cursor, Claude, Copilot, Crush и другими).

Он индексирует ваши API-спецификации в полнотекстовый поисковый движок, предоставляет **16 MCP-инструментов** и позволяет LLM находить, изучать и вызывать реальные API-эндпоинты — без единой строки интеграционного кода.

---

## Содержание

- [Быстрый старт](#быстрый-старт)
- [Конфигурация](#конфигурация)
- [CLI Команды](#cli-команды)
- [MCP Сервер](#mcp-сервер)
- [Интеграция](#интеграция)
- [Поиск](#поиск)
- [Рабочая директория (Workspace)](#рабочая-директория-workspace)
- [Кэширование](#кэширование)
- [Разработка](#разработка)

---

## Быстрый старт

### Вариант 1 — Скачать с GitHub Releases (рекомендуется)

1. Откройте https://github.com/mmadfox/swag2mcp/releases/latest
2. Найдите архив для вашей системы:

   | ОС | Архитектура | Архив |
   |----|-------------|-------|
   | Linux | x86_64 | `swag2mcp_<version>_linux_amd64.tar.gz` |
   | Linux | ARM64 | `swag2mcp_<version>_linux_arm64.tar.gz` |
   | macOS | Intel | `swag2mcp_<version>_darwin_amd64.tar.gz` |
   | macOS | Apple Silicon | `swag2mcp_<version>_darwin_arm64.tar.gz` |
   | Windows | x86_64 | `swag2mcp_<version>_windows_amd64.zip` |

3. Скачайте и установите:

   **Linux / macOS:**
   ```bash
   tar -xzf swag2mcp_<version>_<os>_<arch>.tar.gz
   sudo mv swag2mcp /usr/local/bin/
   swag2mcp --version
   ```

   **Windows (PowerShell):**
   ```powershell
   Expand-Archive swag2mcp_<version>_windows_amd64.zip -DestinationPath .
   move swag2mcp.exe C:\Windows\System32\
   swag2mcp --version
   ```

4. (Опционально) Повторите для mock-сервера — скачайте `swag2mcp-mock_<version>_<os>_<arch>.tar.gz`

### Вариант 2 — Установка через Go

Если у вас установлен Go:

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

### После установки

```bash
# Инициализация рабочей директории
swag2mcp init

# Запуск MCP сервера (для LLM-агентов)
swag2mcp mcp

# Или интерактивный проводник
swag2mcp run
```---

## Example LLM Queries

After setup, try asking your agent:

| Query | What happens |
|-------|-------------|
| "Show me all available APIs" | `spec_list` — lists petstore, binance, countries |
| "What endpoints does Binance have?" | `endpoint_by_spec` — shows 4 market data endpoints |
| "Find endpoints related to pets" | `search("pet")` — finds petstore endpoints |
| "What tags are in the Petstore API?" | `tag_by_spec` — shows "pets" tag |
| "Show me the GET /pets endpoint details" | `inspect` — shows parameters and response schema |
| "Get the current BTC price from Binance" | `invoke` — real API call to Binance |
| "Find countries in Europe" | `invoke` — calls REST Countries API |

---

---

## Конфигурация

### YAML Схема

```yaml
mock_enabled: true                    # опционально, включает режим мок-сервера

http_client:                        # опционально, глобальные настройки HTTP
  random: false                     # опционально, случайные browser-like заголовки
  proxy:                            # опционально
    url: socks5h://127.0.0.1:1080   # http, https, socks5, socks5h
    username: ""                    # опционально
    password: ""                    # опционально
    bypass: []                      # опционально, напр. ["*.local", "10.0.0.0/8"]
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
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/petstore.json  # обязательно, 5-250 символов
        disable: false                  # опционально
        base_url: ""                    # опционально, переопределяет base_url спецификации
        base_mock_url: localhost:8080   # опционально, формат "host:port" или "host:port/path"
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
| `hmac` | `api_key`, `secret_key` | `api_key: $(API_KEY)`, `secret_key: $(SECRET_KEY)` |
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

Запуск MCP сервера в headless-режиме. Основная production-команда для интеграции с LLM.

| Флаг | Сокращение | По умолчанию | Описание |
|------|------------|--------------|----------|
| `--logfile` | `-f` | `""` | Путь к лог-файлу |
| `--tags` | `-t` | `""` | Фильтр спецификаций по тегам |
| `--disable-llm-auth` | | `true` | `true` — аутентификация под капотом. `false` — LLM может запрашивать токены через `auth` |
| `--dump-dir` | | `""` | Директория для дампа HTTP запросов (отладка) |
| `--transport` | | `"stdio"` | Транспорт MCP: `stdio`, `sse`, `streamable-http` |
| `--http-addr` | | `":8080"` | Адрес HTTP сервера (для sse/streamable-http) |
| `--http-path` | | `"/mcp"` | Путь HTTP для MCP handler |
| `--auth-token` | | `""` | Bearer токен для HTTP транспорта |

```bash
swag2mcp mcp                                    # stdio (локальный агент)
swag2mcp mcp --tags=public --logfile=/var/log/swag2mcp.log
swag2mcp mcp --disable-llm-auth=false
swag2mcp mcp --dump-dir=/tmp/dump
swag2mcp mcp --transport streamable-http --http-addr :9090  # удалённый агент
```

### `version`

Показать версию swag2mcp. Также доступен как флаг `--version`.

```bash
swag2mcp version
swag2mcp --version
```

### `info [path]`

Показать детальную информацию о конфигурации и рантайме в формате JSON.

```bash
swag2mcp info
swag2mcp info ./
```

### `import [path] [source] [name]`

Импортировать spec-файлы в рабочую директорию.

| Флаг | Сокращение | По умолчанию | Описание |
|------|------------|--------------|----------|
| `--spec` | `-s` | `nil` | Импорт коллекций из указанных spec'ов (через запятую) |
| `--from-zip` | | `""` | Восстановить workspace из ZIP-архива |

```bash
swag2mcp import https://example.com/spec.yaml myspec
swag2mcp import --spec petstore
swag2mcp import --from-zip /path/to/backup.zip
```

### `export [path] [output]`

Экспортировать workspace в ZIP-архив.

| Флаг | Сокращение | По умолчанию | Описание |
|------|------------|--------------|----------|
| `--spec` | `-s` | `nil` | Экспортировать только указанные spec'ы (через запятую) |

```bash
swag2mcp export
swag2mcp export /path/to/workspace /path/to/backup.zip
swag2mcp export --spec petstore
```

### `mockserver [path]`

Запускает мок-серверы для всех API спецификаций (отдельный бинарник: `swag2mcp-mock`).

| Флаг | По умолчанию | Описание |
|------|-------------|----------|
| `--tls` | `false` | Включить TLS с самоподписанным сертификатом |
| `--tls-cert` | `""` | Путь к TLS сертификату |
| `--tls-key` | `""` | Путь к TLS ключу |

```bash
swag2mcp-mock
swag2mcp-mock --tls
```

**Рабочий процесс:**
1. Добавить `mock_enabled: true` и `base_mock_url` в конфиг
2. Запустить мок-сервер: `swag2mcp-mock`
3. Запустить MCP сервер: `swag2mcp mcp` — invoke будет использовать `base_mock_url`
4. Аутентификация применяется автоматически: OAuth2/Digest используют мок-серверы на портах 9090/9091; остальные типы применяют credentials напрямую

### Mock аутентификация

Когда в конфиге указан `auth`, MCP сервер применяет аутентификацию автоматически.
Только два типа аутентификации требуют отдельного мок-сервера:

| Тип | Мок-эндпоинт | Поведение |
|-----|-------------|-----------|
| `oauth2-cc` / `oauth2-pwd` | `POST /token` на порту 9090 | Принимает любые `client_id`/`username`+`password`, возвращает `{"access_token":"<random>","token_type":"Bearer","expires_in":3600}` |
| `digest` | `GET /` на порту 9091 | Отправляет 401 challenge с `algorithm=MD5`, принимает любой Digest response, возвращает `{"status":"authenticated","method":"digest"}` |

Остальные типы (`basic`, `bearer`, `api-key`, `hmac`, `script`) **не требуют** мок-сервера —
MCP сервер сам применяет настроенные credentials к каждому запросу.

---

## MCP Сервер

MCP сервер предоставляет **16 инструментов** через stdio или HTTP транспорт. LLM-агенты (Opencode, Cursor, Claude, Copilot, Crush и др.) подключаются автоматически после настройки.

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
| `info` | — | Runtime info | Версия swag2mcp, конфиг, статистика |

---

## Интеграция

swag2mcp работает по протоколу MCP и совместим с любым MCP-клиентом.

### Локально (stdio) — агент на той же машине

Запустите сервер:

```bash
swag2mcp mcp
```

| Клиент | Файл конфига | Содержимое |
|--------|-------------|------------|
| **OpenCode** | `opencode.json` | `{"mcp":{"swag2mcp":{"type":"local","command":["swag2mcp","mcp"]}}}` |
| **Cursor** | `.cursor/mcp.json` | `{"mcpServers":{"swag2mcp":{"command":"swag2mcp","args":["mcp"]}}}` |
| **Claude Desktop** | `claude_desktop_config.json` | `{"mcpServers":{"swag2mcp":{"command":"swag2mcp","args":["mcp"]}}}` |
| **VS Code** | `.vscode/mcp.json` | `{"servers":{"swag2mcp":{"type":"stdio","command":"swag2mcp","args":["mcp"]}}}` |
| **Crush** | `crush.json` | `{"mcp":{"swag2mcp":{"type":"stdio","command":"swag2mcp","args":["mcp"]}}}` |

### Удалённо (HTTP) — агент в облаке / на другой машине

Запустите сервер с HTTP транспортом:

```bash
swag2mcp mcp --transport streamable-http --http-addr :8080 --auth-token my-secret
```

Или настройте в `swag2mcp.yaml`:

```yaml
mcp:
  transport: streamable-http
  addr: ":8080"
  path: "/mcp"
  auth_token: $(MCP_AUTH_TOKEN)
```

| Клиент | Файл конфига | Содержимое |
|--------|-------------|------------|
| **OpenCode** | `opencode.json` | `{"mcp":{"swag2mcp":{"type":"remote","url":"http://localhost:8080/mcp","headers":{"Authorization":"Bearer ${MCP_AUTH_TOKEN}"}}}}` |
| **VS Code** | `.vscode/mcp.json` | `{"servers":{"swag2mcp":{"type":"http","url":"http://localhost:8080/mcp"}}}` |

> **Проверка здоровья** (работает без MCP handshake):
> ```bash
> curl http://localhost:8080/health
> # → {"status":"ok","version":"v1.1.3"}
> ```

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

## Кэширование

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
