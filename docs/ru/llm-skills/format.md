# swag2mcp-format

Скилл **swag2mcp-format** обучает вашего LLM отображать ответы MCP-инструментов swag2mcp в чистом, компактном, человекочитаемом markdown-формате. Без этого скилла LLM сам решает, как форматировать ответы — часто многословно и непоследовательно.

## Что покрывает

Все MCP-инструменты swag2mcp:

- `spec_list`, `spec_by_id` — обзор и детали спецификаций
- `collection_by_spec`, `collection_by_id` — коллекции с тегами
- `tag_by_spec`, `tag_by_collection`, `tag_by_id` — списки тегов
- `endpoint_by_spec`, `endpoint_by_collection`, `endpoint_by_tag`, `endpoint_by_id` — списки эндпоинтов
- `search` — результаты поиска
- `inspect` — полная информация об операции с компактными схемами
- `invoke` — результаты API-вызовов
- `auth` — информация об аутентификации
- `info` — информация о рантайме

## Прямая ссылка

<https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md>

## Установка через LLM-агента

Скопируйте этот запрос в вашу IDE с ИИ:

```
Создай директорию .agents/skills/swag2mcp-format/ и добавь скилл из
https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md
```

## Ручная установка

```bash
mkdir -p .agents/skills/swag2mcp-format
curl -o .agents/skills/swag2mcp-format/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md
```

## Требуется перезагрузка

После добавления скилла перезапустите LLM-клиент или IDE (см. [Обзор](overview.md#требуется-перезагрузка)).
