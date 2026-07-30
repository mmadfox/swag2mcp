# Skills for LLM

Skills — это markdown-файлы, которые обучают вашего LLM-агента эффективнее работать с swag2mcp. Они загружаются как часть системного промпта агента и дают LLM точные инструкции для форматирования ответов и понимания CLI-команд.

## Доступные скиллы

| Скилл | Описание | Скачать |
|-------|----------|---------|
| **swag2mcp-format** | Форматирует ответы MCP-инструментов в компактные человекочитаемые markdown-таблицы | [SKILL.md](https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md) |
| **swag2mcp-cli** | Полный CLI-справочник — LLM знает каждую команду, флаг и опцию конфига | [SKILL.md](https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md) |

## Зачем нужны скиллы

Без скилла форматирования LLM сам решает, как отображать результаты инструментов — часто многословно и непоследовательно. Скилл форматирования обеспечивает единый чистый стиль: компактные таблицы для списков, inline-заголовки для деталей и компактные схемы.

CLI-скилл позволяет LLM точно отвечать на вопросы "как сделать..." о командах swag2mcp, не додумывая и не ошибаясь.

## Установка через LLM-агента

Скопируйте этот запрос в вашу IDE с ИИ (OpenCode, Cursor, Claude Desktop, VS Code и т.д.):

```
Добавь скиллы swag2mcp в мой проект:

1. Создай директорию .agents/skills/swag2mcp-format/ и добавь скилл из https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md
2. Создай директорию .agents/skills/swag2mcp-cli/ и добавь скилл из https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md
```

Агент скачает оба файла скиллов и разместит их в правильных директориях.

## Ручная установка

Если ваш LLM-клиент не поддерживает установку через агента, скачайте файлы вручную:

```bash
mkdir -p .agents/skills/swag2mcp-format
mkdir -p .agents/skills/swag2mcp-cli

curl -o .agents/skills/swag2mcp-format/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md
curl -o .agents/skills/swag2mcp-cli/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md
```

## Настройка LLM-клиента

У каждого LLM-клиента и IDE свой способ установки скиллов. Пример ниже для **OpenCode** — обратитесь к документации вашего клиента для правильного метода.

Для OpenCode добавьте скиллы в `opencode.json`:

```json
{
  "skills": [
    {
      "name": "swag2mcp-format",
      "sourceURL": "https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md"
    },
    {
      "name": "swag2mcp-cli",
      "sourceURL": "https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md"
    }
  ]
}
```

## Требуется перезагрузка

**После добавления скиллов перезапустите LLM-клиент или IDE.** Некоторые инструменты загружают скиллы только при запуске. Если скиллы не работают, попробуйте:

- **OpenCode**: Перезапустите приложение или выполните команду opencode заново
- **Cursor**: Закройте и откройте окно (`Cmd+Shift+W` / `Ctrl+Shift+W`)
- **Claude Desktop**: Выйдите и запустите приложение снова
- **VS Code**: Перезагрузите окно (`Ctrl+Shift+P` → "Developer: Reload Window")
