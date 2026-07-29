# swag2mcp-cli

Скилл **swag2mcp-cli** предоставляет вашему LLM полный справочник CLI swag2mcp — каждую команду, флаг, аргумент и опцию конфигурации. С этим скиллом LLM может точно отвечать на вопросы "как сделать...", не додумывая.

## Что покрывает

Все 13 команд CLI:

| Команда | Назначение |
|---------|------------|
| `init` | Инициализация workspace и конфига |
| `add` | Добавление спецификации или коллекции |
| `delete` | Удаление спецификации или коллекции |
| `ls` | Список настроенных спецификаций |
| `run` | Запуск API Explorer TUI |
| `validate` | Валидация файла конфигурации |
| `clean` | Очистка кэшированных данных |
| `update` | Обновление кэша из конфигурации |
| `mcp` | Запуск MCP-сервера |
| `version` | Показ версии |
| `info` | Показ информации о рантайме |
| `import` | Импорт workspace из ZIP-файла |
| `export` | Экспорт workspace в ZIP-файл |

Плюс все флаги, структура конфигурационного файла, методы аутентификации и расширенные опции.

## Установка через LLM-агента

Скопируйте этот запрос в вашу IDE с ИИ:

```
Создай директорию .agents/skills/swag2mcp-cli/ и добавь скилл из
https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md
```

## Ручная установка

```bash
mkdir -p .agents/skills/swag2mcp-cli
curl -o .agents/skills/swag2mcp-cli/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md
```

## Требуется перезагрузка

После добавления скилла перезапустите LLM-клиент или IDE (см. [Обзор](overview.md#требуется-перезагрузка)).
