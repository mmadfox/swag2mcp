# Cursor 統合

## stdio

### Via Settings UI

1. Cursor設定を開く (Cmd+, / Ctrl+,)
2. Go to **MCPサーバー**
3. Click **新規サーバーを追加**
4. Fill in:
   - **名前:** `swag2mcp`
   - **タイプ:** `command`
   - **コマンド:** `swag2mcp mcp`
5. Click **保存**

### 設定ファイル経由

In `~/.cursor/mcp.json`:

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

After connecting, Cursor AI Agent can:

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
