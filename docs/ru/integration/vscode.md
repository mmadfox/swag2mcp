# Интеграция с VS Code

## Via .vscode/mcp.json

1. Установите MCP-расширение для VS Code (например, MCP Client от org.mcp или аналогичное).
2. Создайте `.vscode/mcp.json` в корне проекта:

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "${workspaceFolder}"]
    }
  }
}
```

> // "${{workspaceFolder}}" будет передан как путь к рабочей области

3. Перезагрузите окно VS Code (Ctrl+Shift+P → "Reload Window").
4. Используйте AI-ассистента — теперь он будет знать о ваших API.

## Альтернатива: через настройки VS Code

Также можно настроить в `.vscode/settings.json`:

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp"]
    }
  }
}
```

## Использование

После настройки AI-ассистент VS Code сможет работать с вашими API через swag2mcp.

## Другие клиенты

Не нашли свой клиент? Все MCP-интеграции следуют одному шаблону:
- Укажите команду `swag2mcp` с аргументом `mcp`
- При необходимости добавьте путь к рабочей области: `mcp /path/to/workspace`
- Проверьте документацию вашего клиента для точного расположения и формата файла конфигурации

Большинство MCP-клиентов поддерживают stdio-транспорт, а некоторые — HTTP (SSE / Streamable HTTP).
