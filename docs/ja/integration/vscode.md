# VS Code 統合

## Via .vscode/mcp.json

1. VS Code 用の MCP 拡張機能をインストールします（例：org.mcp の MCP Client など）。
2. プロジェクトルートに `.vscode/mcp.json` を作成：

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

> // "${{workspaceFolder}}" はワークスペースパスとして渡されます

3. VS Code ウィンドウを再読み込みします（Ctrl+Shift+P → "Reload Window"）。
4. AI アシスタントを使用 — これで API を認識します。

## 代替：VS Code 設定経由

`.vscode/settings.json` で設定することもできます：

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

設定後、VS Code AI アシスタントは swag2mcp を通じて API を操作できます。

## その他

お使いのクライアントが見つかりませんか？すべての MCP 統合は同じパターンに従います：
- コマンドを `swag2mcp`、引数を `mcp` に設定
- オプションでワークスペースパスを追加：`mcp /path/to/workspace`
- 正確な設定ファイルの場所と形式については、クライアントのドキュメントを確認

ほとんどの MCP クライアントは stdio トランスポートをサポートし、一部は HTTP（SSE / Streamable HTTP）をサポートしています。
