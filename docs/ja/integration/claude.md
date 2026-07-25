# Claude Desktop 統合

## stdio

`claude_desktop_config.json` に設定：

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

## カスタムワークスペース

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/path/to/workspace"]
    }
  }
}
```

## 使用方法

Claude Desktop を再起動後、以下のことが可能です：

- "すべての API のリストを表示して"
- "注文を作成するエンドポイントを探して"
- "モスクワの天気 API を呼び出して"

## その他

お使いのクライアントが見つかりませんか？すべての MCP 統合は同じパターンに従います：
- コマンドを `swag2mcp`、引数を `mcp` に設定
- オプションでワークスペースパスを追加：`mcp /path/to/workspace`
- 正確な設定ファイルの場所と形式については、クライアントのドキュメントを確認

ほとんどの MCP クライアントは stdio トランスポートをサポートし、一部は HTTP（SSE / Streamable HTTP）をサポートしています。
