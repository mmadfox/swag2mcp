# FAQ

## 一般

### swag2mcp とは何ですか？どのような問題を解決しますか？

swag2mcp は OpenAPI/Swagger/Postman の API 仕様と LLM エージェントを Model Context Protocol (MCP) を介して連携させます。各 API を AI エージェントに接続するためのカスタムコードを書く代わりに、YAML ファイルに一度設定するだけで、LLM が API を発見、調査、呼び出すための 19 のツールを利用できるようになります。

### 他の API-to-LLM ツールとの違いは？

- **コーディング不要** — YAML で API を設定、統合コードは不要
- **19 の MCP ツール** — 発見から呼び出し、大規模レスポンス処理まで完全なツールキット
- **9 つの認証方式** — あらゆる API 認証スキームに対応
- **全文検索** — bluge エンジンによる全エンドポイントの検索
- **TUI エクスプローラー** — ブラウジングとテストのための対話型ターミナルインターフェース
- **モックサーバー** — 実際の API 呼び出しなしでテスト可能

### 対応している API 仕様フォーマットは？

OpenAPI 3.x、Swagger 2.0、Postman Collections v2.1 に対応しています。

### spec と collection の違いは？

**Spec** は論理的な API サービス（例：「Open-Meteo Weather APIs」）を表します。**Collection** は 1 つの OpenAPI/Swagger/Postman ファイルです。1 つの spec は複数の collection を持つことができます — 例えば、API が異なるサービス（予報、大気質、海洋）用に別々の spec ファイルを持つ場合などです。

### 対応している MCP トランスポートは？

3 つのトランスポート：`stdio`（デフォルト、ローカル LLM クライアント用）、`sse`（リモートクライアント用 Server-Sent Events）、`streamable-http`（最新の HTTP ストリーミング）。

### 任意の LLM で swag2mcp を使用できますか？

はい、MCP プロトコルをサポートする任意の LLM クライアントで使用できます：Claude Desktop、VS Code、Cursor、Windsurf、JetBrains IDE、OpenCode など。

## インストール

### swag2mcp をインストールするには？

```bash
# オプション 1: GitHub Releases からダウンロード
# https://github.com/mmadfox/swag2mcp/releases/latest にアクセス
# お使いの OS とアーキテクチャに合ったアーカイブをダウンロード

# オプション 2: Go でインストール
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

### Go のインストールは必要ですか？

いいえ。Linux (amd64, arm64)、macOS (amd64, arm64)、Windows (amd64) 向けのプリビルドバイナリが [GitHub Releases ページ](https://github.com/mmadfox/swag2mcp/releases) で入手可能です。

### モックサーバーをインストールするには？

モックサーバーは別のバイナリです：

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

または GitHub Releases から `swag2mcp-mock_<version>_<os>_<arch>.tar.gz` をダウンロードしてください。

## はじめに

### すぐに始めるには？

```bash
# 1. ワークスペースを初期化
mkdir -p .swag2mcp && swag2mcp init ./.swag2mcp

# 2. MCP サーバーを起動（init 後、公開サンプル仕様が含まれています）
swag2mcp mcp
```

`init` 後、ワークスペースには既にいくつかの公開サンプル仕様（icanhazdadjoke、Open-Meteo、Binance、PokéAPI）が含まれています。すぐに MCP サーバーを起動できます — 手動で仕様を追加する必要はありません。

独自の API を追加する場合：

```bash
swag2mcp add spec --yaml - <<EOF
domain: dadjoke
llm_title: icanhazdadjoke API
base_url: https://icanhazdadjoke.com
collections:
  - llm_title: Jokes
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
EOF
```

### swag2mcp を IDE に接続するには？

**VS Code** (`.vscode/settings.json`):
```json
{
  "mcp": {
    "servers": {
      "swag2mcp": {
        "command": "swag2mcp",
        "args": ["mcp", "/absolute/path/to/.swag2mcp"]
      }
    }
  }
}
```

**Cursor** (`~/.cursor/mcp.json`):
```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/absolute/path/to/.swag2mcp"]
    }
  }
}
```

**Claude Desktop** (`claude_desktop_config.json`):
```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/absolute/path/to/.swag2mcp"]
    }
  }
}
```

ワークスペースディレクトリには常に絶対パスを使用してください。

## 設定

### 設定ファイルの場所は？

デフォルト：`~/.swag2mcp/swag2mcp.yaml`。任意のディレクトリに作成し、コマンドにパスを渡すこともできます。

### API を追加するには？

```bash
# 対話モード
swag2mcp add spec

# YAML を使用（スクリプト推奨）
swag2mcp add spec --yaml - <<EOF
domain: my-api
llm_title: My API
base_url: https://api.example.com/v1
collections:
  - llm_title: Main
    location: https://example.com/spec.yaml
EOF
```

### 既存の spec に collection を追加するには？

```bash
swag2mcp add collection --yaml - <<EOF
spec_domain: meteo
llm_title: Air Quality
location: https://example.com/air-quality.yaml
EOF
```

### spec を一時的に無効にするには？

spec 設定で `disable: true` を設定します。その spec は読み込まれず、インデックスも作成されません。

### 読み込む spec をフィルタリングできますか？

はい、`--tags` フラグを使用します：`swag2mcp mcp --tags=public`。一致するタグを持つ spec のみが読み込まれます。

### シークレットに環境変数を使用するには？

認証フィールドで `$(VAR_NAME)` 構文を使用します：

```yaml
auth:
  type: bearer
  config:
    token: "$(MY_API_TOKEN)"
```

起動前に変数を設定します：`export MY_API_TOKEN="eyJhbGci..."`

## 認証

### 対応している認証方式は？

9 つの方式：`none`、`basic`、`bearer`、`digest`、`hmac`、`oauth2-cc`（クライアントクレデンシャル）、`oauth2-pwd`（パスワードグラント）、`api-key`、`script`。

### トークンを渡すには？

設定ファイルまたは環境変数を介して：

```yaml
auth:
  type: bearer
  config:
    token: "$(MY_TOKEN)"
```

### invoke の前に auth を呼び出す必要がありますか？

いいえ。`invoke` ツールは spec の設定から自動的に認証を適用します。`auth` MCP ツールは、ユーザーにトークンを表示したい場合（例：curl コマンド用）にのみ必要です。

### auth ツールが表示されないのはなぜ？

`auth` ツールはデフォルトで無効になっています（`--disable-llm-auth=true`）。これは本番環境向けのセキュリティ対策です。有効にするには：`swag2mcp mcp --disable-llm-auth=false`。

### OAuth2 トークンはどのように更新されますか？

OAuth2 Client Credentials と Password Grant のトークンは、期限切れ時に自動的に更新されます。Bearer トークンは静的であり、手動で更新する必要があります。

## MCP サーバー

### MCP サーバーを起動するには？

```bash
# デフォルト（stdio トランスポート）
swag2mcp mcp

# HTTP トランスポート
swag2mcp mcp --transport sse --http-addr :8080
```

### ポートを変更するには？

```bash
swag2mcp mcp --transport sse --http-addr 0.0.0.0:9090
```

### MCP HTTP エンドポイントを保護するには？

Bearer トークンを設定します：

```bash
swag2mcp mcp --transport sse --http-addr :8080 --auth-token "my-secret"
```

LLM クライアントはすべてのリクエストに `Authorization: Bearer my-secret` を含める必要があります。

### HTTP トランスポートの MCP ハンドシェイクとは？

SSE および Streamable HTTP トランスポートでは、MCP プロトコルは 3 ステップのハンドシェイクを必要とします：

```
Step 1: POST /mcp → {"method":"initialize", ...}
Step 2: POST /mcp → {"method":"notifications/initialized"}
Step 3: POST /mcp → {"method":"tools/list", ...}  ← これで動作
```

初期化前のツール呼び出しは失敗します。

## 使用法

### エンドポイントを検索するには？

`search` MCP ツールまたは TUI（`swag2mcp run`）を使用します。検索はフィールドフィルター（`method:GET`、`tag:pets`）、あいまい検索、ワイルドカード、ブール演算子をサポートしています。

### API を呼び出すには？

LLM が `invoke` MCP ツールを使用します。最初にエンドポイントを調査して必要なパラメーターを理解してください：

```
inspect(endpointId: "...")  → 契約を理解
invoke(endpointId: "...", parameters: {...})  → 呼び出し
```

### レスポンスが大きすぎる場合は？

`max_response_size`（デフォルト 1 MB）を超えるレスポンスはディスクに保存されます。LLM はファイル参照を受け取り、`response_outline`、`response_compress`、`response_slice` ツールで探索できます。

### レートリミッターはどのように動作しますか？

各エンドポイントには 10 秒のクールダウンがあります。LLM が 10 秒以内に同じエンドポイントを 2 回呼び出すと、2 回目の呼び出しは静かにブロックされます。設定で無効化または調整できます。

### 実際の API 呼び出しなしでテストできますか？

はい、モックサーバーを使用します：

```bash
swag2mcp-mock mockserver
```

OpenAPI スキーマに基づいて偽のレスポンスを生成します。

## ワークスペース管理

### 設定をバックアップするには？

```bash
swag2mcp export --output ~/backups/swag2mcp-2026-07-24.zip
```

### 別のマシンに移行するには？

```bash
# 古いマシンで
swag2mcp export --output swag2mcp.zip

# ZIP をコピーし、新しいマシンで
swag2mcp import --from-zip swag2mcp.zip
```

### spec ファイルを更新するには？

```bash
swag2mcp update
```

設定を再検証し、キャッシュをクリアし、すべての spec ファイルを再ダウンロードします。

### ディスク容量をクリーンアップするには？

```bash
swag2mcp clean
```

キャッシュされた spec ファイルと保存された API レスポンスを削除します。古いレスポンス（48 時間以上）は MCP サーバー起動時に自動的にクリーンアップされます。

## TUI

### TUI とは？どう使うの？

TUI（Terminal User Interface）は対話型 API エクスプローラーです。`swag2mcp run` で起動します。3 つのモードがあります：Search（全文検索）、Browse（ツリー移動：Spec → Collection → Tag → Endpoint）、Auth（トークン表示）。

### キーボードショートカットは？

| キー | アクション |
|------|-----------|
| `↑/↓` | 移動 |
| `Enter` | 選択 |
| `Esc` | 戻る |
| `Tab` | モード切替 |
| `/` | 検索 |
| `N/P` | 次/前のページ |
| `q` | 終了 |

## 高度な設定

### プロキシを使用できますか？

はい、`http_client.proxy` で設定します：

```yaml
http_client:
  proxy:
    url: "http://proxy.company.com:8080"
    username: "$(PROXY_USER)"
    password: "$(PROXY_PASS)"
    bypass:
      - "localhost"
      - "*.internal.com"
```

### カスタム認証方式を追加できますか？

はい、`internal/auth/` に `Authenticator` インターフェースを実装し、設定パーサーに登録します。詳細は開発セクションを参照してください。

### カスタム MCP ツールを追加できますか？

はい、`Svc` インターフェースにメソッドを追加し、サービスレイヤーに実装し、ハンドラーを追加して登録します。詳細は開発セクションを参照してください。

### `swag2mcp` と `swag2mcp-mock` の違いは？

`swag2mcp` は CLI コマンドと MCP サーバーを持つメインバイナリです。`swag2mcp-mock` は、実際の API 呼び出しなしでテストするためのモックサーバーを起動する別のバイナリです。
