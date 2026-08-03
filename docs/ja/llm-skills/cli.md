# swag2mcp-cli

**swag2mcp-cli** スキルは、LLMに完全なswag2mcp CLIリファレンスを提供します — すべてのコマンド、フラグ、引数、設定オプション。このスキルにより、LLMは「どうやって...」という質問に推測せず正確に答えることができます。

## カバー範囲

13のCLIコマンドすべて：

| コマンド | 目的 |
|----------|------|
| `init` | ワークスペースと設定を初期化 |
| `add` | 仕様またはコレクションを追加 |
| `delete` | 仕様またはコレクションを削除 |
| `ls` | 設定済み仕様を一覧表示 |
| `run` | API Explorer TUIを起動 |
| `validate` | 設定ファイルを検証 |
| `clean` | キャッシュデータをクリア |
| `update` | 設定からキャッシュを更新 |
| `mcp` | MCPサーバーを起動 |
| `version` | バージョン情報を表示 |
| `info` | 実行時情報を表示 |
| `import` | spec ファイルのインポートまたはバックアップからワークスペースを復元 |
| `export` | ワークスペースをZIPファイルにエクスポート |

さらに、すべてのフラグ、設定ファイル構造、認証方法、高度なオプション。

## 直接リンク

<https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md>

## LLMエージェント経由でインストール

次のリクエストをAI搭載IDEにコピーしてください：

```
ディレクトリ .agents/skills/swag2mcp-cli/ を作成し、
https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md からスキルを追加
```

## 手動インストール

```bash
mkdir -p .agents/skills/swag2mcp-cli
curl -o .agents/skills/swag2mcp-cli/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md
```

## 再起動が必要です

スキルを追加した後、LLMクライアントまたはIDEを再起動してください（[概要](overview.md#再起動が必要です)を参照）。
