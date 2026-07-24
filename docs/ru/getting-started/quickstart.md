# Быстрый старт

Запустите swag2mcp за 2 минуты.

## 1. Инициализация

### Домашняя директория (рекомендовано)

Одноразовая настройка для всей системы. Конфиг хранится в вашей домашней папке.

::: code-group

```bash [macOS / Linux]
swag2mcp init
# Создаёт ~/.swag2mcp/swag2mcp.yaml
```

```powershell [Windows]
swag2mcp.exe init
# Создаёт %USERPROFILE%\.swag2mcp\swag2mcp.yaml
```

:::

### Директория проекта

Для изолированной рабочей области внутри вашего проекта.

::: code-group

```bash [macOS / Linux]
mkdir -p ./swag2mcp && swag2mcp init ./swag2mcp
```

```powershell [Windows]
mkdir ./swag2mcp; swag2mcp.exe init ./swag2mcp
```

:::

### Из ZIP

Если у вас уже есть готовая рабочая область (например, от коллеги):

```bash
swag2mcp import --from-zip workspace.zip
```

## 2. Установка навыков агента (рекомендовано)

Установите навыки swag2mcp, чтобы обучить вашего AI-агента всем командам, флагам, формату конфига и реальным примерам.

Попросите агента:

```bash
"Создай директорию .agents/skills/swag2mcp-cli и добавь навык из https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md в .agents/skills/swag2mcp-cli/SKILL.md"
"Создай директорию .agents/skills/swag2mcp-format и добавь навык из https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md в .agents/skills/swag2mcp-format/SKILL.md"
```

> Некоторым IDE требуется перезапуск после добавления навыков.

## 3. Настройка LLM-клиента / IDE

Настройте вашу IDE для подключения к swag2mcp. IDE будет автоматически запускать MCP-сервер при необходимости.

::: code-group

```json [OpenCode]
{
  "mcp": {
    "swag2mcp": {
      "type": "local",
      "command": ["swag2mcp", "mcp"],
      "enabled": true
    }
  }
}
```

```json [Claude Desktop]
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp"]
    }
  }
}
```

```json [Crush]
{
  "mcp": {
    "swag2mcp": {
      "type": "stdio",
      "command": "swag2mcp",
      "args": ["mcp"]
    }
  }
}
```

:::

Для других IDE (Cursor, VS Code, JetBrains) смотрите [руководство по интеграции](../integration/opencode.md).

> Если вы инициализировали рабочую область в нестандартном пути (например, `./swag2mcp`), укажите полный путь в команде:
> `"command": ["swag2mcp", "mcp", "/абсолютный/путь/до/swag2mcp"]`

> **После любого изменения конфига перезапустите MCP-сервер**, чтобы изменения вступили в силу.

## 4. Запуск MCP-сервера

### stdio (по умолчанию) — для локальной IDE

Ничего настраивать не нужно. Ваша IDE запускает swag2mcp автоматически через конфиг выше.

```bash
swag2mcp mcp
```

### SSE / Streamable HTTP — для удалённого доступа

```bash
swag2mcp mcp --transport sse --http-addr :8080
```

Или настройте в `swag2mcp.yaml`:

```yaml
mcp:
  transport: sse
  addr: ":8080"
  path: "/mcp"
```

Все флаги смотрите в [справочнике MCP-сервера](../configuration/mcp-server.md).

### Фильтрация спецификаций по тегам

```bash
swag2mcp mcp --tags weather,public
```

Только спецификации с соответствующими тегами будут доступны LLM.

### Проверка работы

После подключения спросите вашего LLM-агента:

```bash
"Какие MCP-инструменты ты поддерживаешь?"
```

Если агент перечислит инструменты swag2mcp (`spec_list`, `search`, `invoke` и т.д.) — всё работает.

### Примеры запросов

| Спросите агента | Что произойдёт |
|-------|-------------|
| "Какая погода в Нью-Йорке?" | `invoke` — вызов API прогноза Open-Meteo |
| "Какая текущая цена BTC?" | `invoke` — вызов Binance ticker API |
| "Расскажи шутку про папу" | `invoke` — вызов icanhazdadjoke API |
| "Покажи Пикачу" | `invoke` — вызов PokéAPI по имени |
| "Кто такой Рик Санчес?" | `invoke` — вызов API персонажей Рика и Морти |
| "Какое качество воздуха в Пекине?" | `invoke` — вызов Open-Meteo air quality API |
| "Насколько высоки волны у берегов Португалии?" | `invoke` — вызов Open-Meteo marine API |
| "Найди шутки про собак" | `invoke` — поиск шуток через dadjoke |
| "Список всех покемонов" | `invoke` — список через PokéAPI |
| "Какая высота Эвереста?" | `invoke` — вызов Open-Meteo elevation API |

## 5. Что дальше?

- [Концепции](../concepts/overview.md) — понимание архитектуры
- [Конфигурация](../configuration/config-file.md) — настройка параметров
- [Команды CLI](../cli/overview.md) — полный справочник команд
