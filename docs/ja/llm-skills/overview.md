# Skills for LLM

スキルは、LLMエージェントがswag2mcpをより効果的に使用する方法を教えるMarkdownファイルです。エージェントのシステムプロンプトの一部として読み込まれ、LLMにレスポンスのフォーマットとCLIコマンドの理解のための正確な指示を提供します。

## 利用可能なスキル

| スキル | 説明 | ソース |
|--------|------|--------|
| **swag2mcp-format** | すべてのMCPツールのレスポンスをコンパクトで読みやすいMarkdownテーブルにフォーマット | `swag2mcp-format/SKILL.md` |
| **swag2mcp-cli** | 完全なCLIリファレンス — LLMはすべてのコマンド、フラグ、設定オプションを把握 | `swag2mcp-cli/SKILL.md` |

## スキルが重要な理由

フォーマットスキルがないと、LLMはツール結果の表示方法を自分で決定します — 多くの場合、冗長で一貫性がありません。フォーマットスキルにより、すべてのレスポンスが同じクリーンなパターンに従います：リストにはコンパクトなテーブル、詳細にはインラインヘッダー、コンパクトなスキーマ。

CLIスキルにより、LLMはswag2mcpコマンドに関する「どうやって...」という質問に推測せず正確に答えることができます。

## LLMエージェント経由でインストール

次のリクエストをAI搭載IDE（OpenCode、Cursor、Claude Desktop、VS Codeなど）にコピーしてください：

```
swag2mcpスキルをプロジェクトに追加してください：

1. ディレクトリ .agents/skills/swag2mcp-format/ を作成し、https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md からスキルを追加
2. ディレクトリ .agents/skills/swag2mcp-cli/ を作成し、https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md からスキルを追加
```

エージェントが両方のスキルファイルをダウンロードし、正しい場所に配置します。

## 手動インストール

LLMクライアントがエージェントベースのセットアップをサポートしていない場合は、手動でファイルをダウンロードしてください：

```bash
mkdir -p .agents/skills/swag2mcp-format
mkdir -p .agents/skills/swag2mcp-cli

curl -o .agents/skills/swag2mcp-format/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md
curl -o .agents/skills/swag2mcp-cli/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md
```

## LLMクライアントの設定

OpenCodeの場合は、`opencode.json`にスキルを追加します：

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

## 再起動が必要です

**スキルを追加した後、LLMクライアントまたはIDEを再起動してください。** 一部のツールは起動時にのみスキルを読み込みます。スキルが反映されない場合は、以下を試してください：

- **OpenCode**：アプリケーションを再起動するか、opencodeコマンドを再実行
- **Cursor**：ウィンドウを閉じて再度開く（`Cmd+Shift+W` / `Ctrl+Shift+W`）
- **Claude Desktop**：アプリケーションを終了して再起動
- **VS Code**：ウィンドウをリロード（`Ctrl+Shift+P` → "Developer: Reload Window"）
