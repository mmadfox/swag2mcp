# Cursor 統合

## stdio

Cursor の設定で MCP サーバーを追加：

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

## 使用方法

接続後、Cursor AI エージェントは以下のことが可能です：

- API を探索する
- 関連するエンドポイントを見つける
- API を呼び出して結果を表示する
- リクエストのデバッグを支援する

## その他

お使いのクライアントが見つかりませんか？すべての MCP 統合は同じパターンに従います：
- コマンドを `swag2mcp`、引数を `mcp` に設定
- オプションでワークスペースパスを追加：`mcp /path/to/workspace`
- 正確な設定ファイルの場所と形式については、クライアントのドキュメントを確認

ほとんどの MCP クライアントは stdio トランスポートをサポートし、一部は HTTP（SSE / Streamable HTTP）をサポートしています。
