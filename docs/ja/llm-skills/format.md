# swag2mcp-format

**swag2mcp-format** スキルは、LLMにswag2mcp MCPツールの応答をクリーンでコンパクトな読みやすいMarkdown形式で表示する方法を教えます。このスキルがないと、LLMは応答のフォーマット方法を自分で決定します — 多くの場合、冗長で一貫性がありません。

## カバー範囲

すべてのswag2mcp MCPツール：

- `spec_list`, `spec_by_id` — 仕様の概要と詳細
- `collection_by_spec`, `collection_by_id` — タグ付きコレクション
- `tag_by_spec`, `tag_by_collection`, `tag_by_id` — タグ一覧
- `endpoint_by_spec`, `endpoint_by_collection`, `endpoint_by_tag`, `endpoint_by_id` — エンドポイント一覧
- `search` — 検索結果
- `inspect` — コンパクトなスキーマ付き完全な操作詳細
- `invoke` — API呼び出し結果
- `auth` — 認証情報
- `info` — 実行時情報

## LLMエージェント経由でインストール

次のリクエストをAI搭載IDEにコピーしてください：

```
ディレクトリ .agents/skills/swag2mcp-format/ を作成し、
https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md からスキルを追加
```

## 手動インストール

```bash
mkdir -p .agents/skills/swag2mcp-format
curl -o .agents/skills/swag2mcp-format/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md
```

## 再起動が必要です

スキルを追加した後、LLMクライアントまたはIDEを再起動してください（[概要](overview.md#再起動が必要です)を参照）。
