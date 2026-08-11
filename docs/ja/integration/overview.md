# 任意の MCP クライアント

swag2mcp は **MCP サーバー**（Model Context Protocol）です。つまり、このセクションに記載されているものだけでなく、**あらゆる MCP クライアント**で動作します。エディタ、IDE、エージェントが MCP プロトコルをサポートしていれば、swag2mcp を接続できます。

## ユニバーサルパターン

すべての MCP クライアントは同じ基本設定を使用します。swag2mcp を MCP サーバーとして追加します：

- **コマンド：** `swag2mcp`
- **引数：** `mcp`（オプションでワークスペースパス：`mcp /path/to/workspace`）

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

設定ファイルの正確な場所と形式（JSON、TOML、GUI 設定）はクライアントによって異なります。**クライアントの MCP ドキュメント**で設定場所を確認してください。

## トランスポート

- **stdio** — どこでも動作します。ほとんどの MCP クライアントがサポート
- **HTTP（SSE / Streamable HTTP）** — HTTP トランスポートオプションを持つクライアントでサポート

トランスポートフラグについては [`mcp` コマンド](/ja/cli/mcp) のリファレンスを参照してください。

## テスト済みの統合

| クライアント | ガイド |
|--------------|--------|
| OpenCode | [OpenCode](/ja/integration/opencode) |
| Cursor | [Cursor](/ja/integration/cursor) |
| Claude Desktop | [Claude Desktop](/ja/integration/claude) |
| VS Code | [VS Code](/ja/integration/vscode) |
| Crush | [Crush](/ja/integration/crush) |

> リストにないクライアントでも、**サポートされていない**という意味ではありません。MCP プロトコルを話す限り、上記のユニバーサルパターンを使用し、クライアントのマニュアルに従って設定場所を確認してください。
