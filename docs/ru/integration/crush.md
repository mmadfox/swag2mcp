# Интеграция с Crush

## stdio

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

## HTTP

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "--transport", "sse", "--http-addr", "127.0.0.1:8080"]
    }
  }
}
```

## Другие клиенты

Не нашли свой клиент? Все MCP-интеграции следуют одному шаблону:
- Укажите команду `swag2mcp` с аргументом `mcp`
- При необходимости добавьте путь к рабочей области: `mcp /path/to/workspace`
- Проверьте документацию вашего клиента для точного расположения и формата файла конфигурации

Большинство MCP-клиентов поддерживают stdio-транспорт, а некоторые — HTTP (SSE / Streamable HTTP).
