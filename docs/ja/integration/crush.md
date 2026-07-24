# Crush 統合

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

## その他

お使いのクライアントが見つかりませんか？すべての MCP 統合は同じパターンに従います：
- コマンドを `swag2mcp`、引数を `mcp` に設定
- オプションでワークスペースパスを追加：`mcp /path/to/workspace`
- 正確な設定ファイルの場所と形式については、クライアントのドキュメントを確認

ほとんどの MCP クライアントは stdio トランスポートをサポートし、一部は HTTP（SSE / Streamable HTTP）をサポートしています。
