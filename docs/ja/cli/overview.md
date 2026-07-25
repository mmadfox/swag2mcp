# CLI コマンド

## 概要

`swag2mcp` CLI は、すべての操作の単一エントリポイントです — ワークスペースの初期化や API 仕様の管理から、LLM 統合のための MCP サーバーの起動まで。OpenAPI/Swagger/Postman の spec を扱う完全なライフサイクルをカバーする **13 のコマンド** を提供します。

### CLI が解決するもの

- **ワークスペースライフサイクル** — 作成（`init`）、検査（`info`、`ls`）、クリーンアップ（`clean`）、更新（`update`）、削除（`delete`）
- **Spec と Collection の管理** — API 仕様とその collection の追加（`add`）、一覧表示（`ls`）、削除（`delete`）
- **実行モード** — LLM ツールアクセス用の MCP サーバー起動（`mcp`）、または対話型 TUI エクスプローラー起動（`run`）
- **診断** — 設定の検証（`validate`）、バージョン表示（`version`）、ランタイム情報表示（`info`）
- **バックアップと復元** — ZIP による完全なワークスペースラウンドトリップ（`export`、`import`）

### 主要なニュアンス

- **パス解決** — `[path]` を受け付けるコマンドは**ワークスペースディレクトリ**（ファイルパスではない）を期待します。解決順序：明示的な `[path]` → カレントディレクトリ（`./`）→ `~/.swag2mcp/`。CLI は自動的に `swag2mcp.yaml` を追加します。サービスとして実行する場合や IDE 設定では、間違ったワークスペースを読み込まないよう常に明示的なパスを渡してください。
- **Spec と Collection の違い** — **spec** は論理的な API サービス（例：「Open-Meteo API」）を表し、**collection** は 1 つの OpenAPI/Swagger/Postman ファイルです。1 つの spec は複数の collection を持つことができます。
- **`--version`** はフラグ（`swag2mcp --version`）とサブコマンド（`swag2mcp version`）の両方としてサポートされています。
- **`add spec` / `add collection`** は `--yaml`（インライン文字列または標準入力の `-`）を介して YAML 入力を受け付けます。ファイルまたはヒアドキュメントからのパイプは、特殊文字によるシェルの引用符問題を回避します。
- **`delete`** は TTY（対話型ターミナル）が必要です。`--force` や `--yes` フラグはありません — 常に選択と確認を促します。
- **`mcp`** は LLM 統合の主要コマンドです。3 つのトランスポートをサポート：`stdio`（デフォルト）、`sse`、`streamable-http`。`--disable-llm-auth` フラグ（デフォルト：`true`）は MCP ツールリストから `auth` ツールを削除し、LLM がトークンを表示したり要求したりするのを防ぎます。認証は引き続き機能します — トークンは LLM 経由ではなく、標準の設定メカニズムを通じて取得されます。このモードは**本番環境**に推奨されます（LLM は認証情報にアクセスできません）。**デバッグ**や短命トークンを使用する場合は、`--disable-llm-auth=false` を設定して LLM が `auth` ツールを介して新しいトークンを要求できるようにします。
- **`validate`** は YAML 構文、設定構造、spec ファイルの存在、URL の到達可能性、spec 形式（OpenAPI/Swagger/Postman）、認証設定、HTTP クライアントの正確性をチェックします。認証エンドポイントや API エンドポイントの可用性は**テストしません**。
- **`export` / `import`** は完全なワークスペースラウンドトリップを提供します — 設定ファイル、spec ファイル、キャッシュ、認証スクリプトがすべて ZIP アーカイブに含まれます。
- **`clean`** は `cache/` と `responses/` ディレクトリを削除しますが、`specs/` と `auth_scripts/` は保持します。古いレスポンス（48 時間以上）は `mcp` 起動時に自動的にクリーンアップされます。

## コマンド一覧

| コマンド | 説明 |
|---------|------|
| [`init`](/cli/init) | デフォルト設定でワークスペースディレクトリを初期化 |
| [`add`](/cli/add) | 設定に spec または collection を追加 |
| [`delete`](/cli/delete) | 対話的に spec または collection を削除 |
| [`ls`](/cli/ls) | すべての spec とその collection を一覧表示 |
| [`run`](/cli/run) | 対話型 TUI API エクスプローラーを起動 |
| [`validate`](/cli/validate) | 設定と spec ファイルを検証 |
| [`clean`](/cli/clean) | キャッシュされた spec と呼び出しレスポンスをクリア |
| [`update`](/cli/update) | すべての spec を再検証、再キャッシュ、再インデックス化 |
| [`mcp`](/cli/mcp) | LLM ツールアクセス用の MCP サーバーを起動 |
| [`version`](/cli/version) | swag2mcp バージョンを表示 |
| [`info`](/cli/info) | 詳細な設定とランタイム情報を表示 |
| [`import`](/cli/import) | spec ファイルをインポート、または ZIP からワークスペースを復元 |
| [`export`](/cli/export) | ワークスペースをポータブルな ZIP バックアップとしてエクスポート |

## グローバルフラグ

| フラグ | 説明 |
|-------|------|
| `--version` | バージョンを表示（`version` サブコマンドと同じ） |
| `--help` | 任意のコマンドのヘルプを表示 |
