# swag2mcp

**swag2mcp** は、OpenAPI/Swagger/Postman API仕様とLLMエージェント（Opencode、Crush、Copilot、Cursorなど）を橋渡しするCLIツールおよびMCP（Model Context Protocol）サーバーです。

API仕様を全文検索エンジンにインデックスし、16のMCPツールとして公開し、LLMが実際のAPIエンドポイントを発見、検査、呼び出しできるようにします——統合コードを1行も書く必要はありません。

---

## 目次

- [クイックスタート](#クイックスタート)
- [設定](#設定)
- [CLIコマンド](#cliコマンド)
- [MCPサーバー](#mcpサーバー)
- [検索](#検索)
- [ワークスペース](#ワークスペース)
- [キャッシュ](#キャッシュ)
- [開発](#開発)

---

## クイックスタート

### オプション1 — GitHub Releasesからダウンロード（推奨）

1. https://github.com/mmadfox/swag2mcp/releases/latest を開く
2. お使いのシステムに合ったアーカイブを見つける：

   | OS | アーキテクチャ | アーカイブ |
   |----|-------------|-----------|
   | Linux | x86_64 | `swag2mcp_<version>_linux_amd64.tar.gz` |
   | Linux | ARM64 | `swag2mcp_<version>_linux_arm64.tar.gz` |
   | macOS | Intel | `swag2mcp_<version>_darwin_amd64.tar.gz` |
   | macOS | Apple Silicon | `swag2mcp_<version>_darwin_arm64.tar.gz` |
   | Windows | x86_64 | `swag2mcp_<version>_windows_amd64.zip` |

3. ダウンロードしてインストール：

   **Linux / macOS:**
   ```bash
   tar -xzf swag2mcp_<version>_<os>_<arch>.tar.gz
   sudo mv swag2mcp /usr/local/bin/
   swag2mcp --version
   ```

   **Windows (PowerShell):**
   ```powershell
   Expand-Archive swag2mcp_<version>_windows_amd64.zip -DestinationPath .
   move swag2mcp.exe C:\Windows\System32\
   swag2mcp --version
   ```

4. （オプション）モックサーバーも同様に — `swag2mcp-mock_<version>_<os>_<arch>.tar.gz` をダウンロード

### オプション2 — Goでインストール

Goがインストールされている場合：

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

### インストール後

```bash
# ワークスペースを初期化
swag2mcp init

# MCPサーバーを起動（LLMエージェント用）
swag2mcp mcp

# またはインタラクティブエクスプローラー
swag2mcp run
```---

## Example LLM Queries

After setup, try asking your agent:

| Query | What happens |
|-------|-------------|
| "Show me all available APIs" | `spec_list` — lists petstore, binance, dadjoke, pokeapi |
| "What endpoints does Binance have?" | `endpoint_by_spec` — shows 4 market data endpoints |
| "Find endpoints related to pets" | `search("pet")` — finds petstore endpoints |
| "What tags are in the Petstore API?" | `tag_by_spec` — shows "pets" tag |
| "Show me the GET /pets endpoint details" | `inspect` — shows parameters and response schema |
| "Get the current BTC price from Binance" | `invoke` — real API call to Binance |
| "Get a random dad joke" | `invoke` — calls icanhazdadjoke API |

---

---

## 設定

### YAMLスキーマ

```yaml
mock_enabled: true                    # オプション、モックサーバーモードを有効化

http_client:                        # オプション、グローバルHTTPデフォルト
  headers:                          # オプション
    X-API-Version: "2"
  cookies: []                       # オプション
  user_agent: ""                    # オプション
  timeout: 0s                       # オプション
  follow_redirects: true            # オプション
  max_redirects: 10                 # オプション
  max_response_size: 1048           # オプション、バイト（デフォルト1KB、最大1MB）

specs:
  - domain: petstore                    # 必須、1-60文字、[a-zA-Z0-9_-]
    llm_title: Petstore API             # 必須、5-120文字
    llm_instruction: |                  # オプション、最大500文字
      このAPIを使用してペット、注文、ユーザーを管理します。
    base_url: https://petstore.swagger.io/v2  # 必須、有効なURL
    disable: false                      # オプション
    tags: [public, demo]                # オプション、フィルタリング用
    http_client:                        # オプション、グローバル設定を上書き
      headers:
        X-API-Version: "2"
    auth:                               # オプション
      type: bearer                      # 認証方法を参照
      config:
        token: $(TOKEN_AUTH)
    collections:
      - llm_title: Petstore Swagger     # オプション、最大120文字
        llm_instruction: |             # オプション、最大360文字
          Petstoreの主要エンドポイント
        title: ""                      # オプション、specから自動入力
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/petstore.json  # 必須、5-250文字
        disable: false                  # オプション
        base_url: ""                    # オプション、specのbase_urlを上書き
        base_mock_url: localhost:8080   # オプション、形式 "host:port" または "host:port/path"
        http_client: {}                 # オプション、specを上書き
```

### タグ — プロジェクトによる仕様のフィルタリング

タグを使用すると、仕様をプロジェクト、環境、チームごとにグループ化できます。MCPサーバーの起動時に `--tags` を使用して、一致する仕様のみを読み込みます：

```bash
# 公開仕様のみでサーバーを起動
swag2mcp mcp --tags=public

# 複数のタグでサーバーを起動
swag2mcp mcp --tags=public,internal

# 異なるプロジェクト用に複数のサーバーを実行
swag2mcp mcp --tags=project-alpha --logfile=/tmp/swag2mcp-alpha.log
swag2mcp mcp --tags=project-beta  --logfile=/tmp/swag2mcp-beta.log
```

これにより、単一の設定ファイルから異なるプロジェクト用に個別のMCPサーバーを実行できます。

### 認証方法

| タイプ | フィールド | 設定例 |
|--------|-----------|--------|
| `none` | — | `type: none` |
| `basic` | `username`, `password` | `username: $(USER)`, `password: $(PASS)` |
| `bearer` | `token` | `token: $(TOKEN)` |
| `digest` | `username`, `password` | `username: admin`, `password: secret` |
| `hmac` | `api_key`, `secret_key` | `api_key: $(API_KEY)`, `secret_key: $(SECRET_KEY)` |
| `api-key` | `key`, `value`, `in` (header/query) | `key: X-API-Key`, `value: $(KEY)`, `in: header` |
| `oauth2-cc` | `client_id`, `client_secret`, `token_url`, `scopes` | `client_id: $(ID)`, `token_url: https://auth.example.com/token` |
| `oauth2-pwd` | `username`, `password`, `client_id`, `client_secret`, `token_url`, `scopes` | `username: $(USER)`, `token_url: https://auth.example.com/token` |
| `script` | `source` | `source: パス/to/auth.sh` |

すべての文字列フィールドは `$(ENV_VAR)` 構文をサポートしています——実行時に環境変数から解決されます。

---

## CLIコマンド

`[path]` を受け付けるすべてのコマンドは、同じパス解決を使用します：

```
swag2mcp <command>          → ~/.swag2mcp/
swag2mcp <command> ./       → {cwd}/.swag2mcp/
swag2mcp <command> path/to  → {cwd}/path/to/.swag2mcp/
```

### `init [path]`

ワークスペースと設定を初期化します。

| フラグ | 短縮 | デフォルト | 説明 |
|--------|------|-----------|------|
| `--interactive` | `-i` | `false` | インタラクティブウィザードを実行 |
| `--force` | `-f` | `false` | 既存の設定を上書き |

```bash
swag2mcp init              # ~/.swag2mcp/swag2mcp.yaml を作成
swag2mcp init ./           # ./.swag2mcp/swag2mcp.yaml を作成
swag2mcp init -i           # インタラクティブウィザード
```

### `add spec [path]` / `add collection [path]`

設定に仕様またはコレクションを追加します。

| フラグ | 短縮 | デフォルト | 説明 |
|--------|------|-----------|------|
| `--yaml` | `-y` | `""` | YAML入力（stdinには `-` を使用） |
| `--example` | `-e` | `false` | YAML例を表示 |

```bash
swag2mcp add spec
swag2mcp add spec --yaml 'domain: petstore\nllm_title: Petstore API\nbase_url: https://...'
cat spec.yaml | swag2mcp add spec --yaml -
swag2mcp add spec --example
```

### `delete spec [path]` / `delete collection [path]`

設定から仕様またはコレクションを削除します。インタラクティブな選択。

```bash
swag2mcp delete spec
swag2mcp delete collection
```

### `ls [path]`

仕様とコレクションを一覧表示します。

| フラグ | 短縮 | デフォルト | 説明 |
|--------|------|-----------|------|
| `--tags` | `-t` | `""` | タグでフィルタリング（カンマ区切り） |

```bash
swag2mcp ls
swag2mcp ls --tags=public,internal
```

### `run [path]`

インタラクティブAPIエクスプローラー（TUI）。エンドポイントの検索、閲覧、検査、保存。

```bash
swag2mcp run
```

### `validate [path]`

設定を検証し、すべてのコレクションの場所がアクセス可能か確認します。

| フラグ | 短縮 | デフォルト | 説明 |
|--------|------|-----------|------|
| `--tags` | `-t` | `""` | タグで仕様をフィルタリング |

```bash
swag2mcp validate
swag2mcp validate --tags=public
```

### `clean [path]`

`cache/` および `responses/` ディレクトリのすべての内容を削除します。

```bash
swag2mcp clean
```

### `update [path]`

設定を検証し、キャッシュをクリアし、すべてのspecファイルを再キャッシュします。

```bash
swag2mcp update
```

### `mcp [path]`

ヘッドレスモードでMCPサーバーを起動します（stdioトランスポート）。LLM統合のための主要な本番コマンドです。

| フラグ | 短縮 | デフォルト | 説明 |
|--------|------|-----------|------|
| `--logfile` | `-f` | `""` | ログファイルのパス |
| `--tags` | `-t` | `""` | タグで仕様をフィルタリング |
| `--disable-llm-auth` | | `true` | `true` — 認証はバックグラウンドで実行（LLMはトークンを見ません）。`false` — LLMは `auth` ツールでトークンを要求可能 |
| `--dump-dir` | | `""` | HTTPリクエストダンプディレクトリ（デバッグ用） |

```bash
swag2mcp mcp
swag2mcp mcp --tags=public --logfile=/var/log/swag2mcp.log
swag2mcp mcp --disable-llm-auth=false
swag2mcp mcp --dump-dir=/tmp/dump
```

### `mockserver [path]`

すべてのAPI仕様に対してモックHTTPサーバーを起動します。各コレクションは独自の
HTTPサーバーを持ち、OpenAPIレスポンススキーマに一致するランダムデータを生成します。

| フラグ | デフォルト | 説明 |
|--------|-----------|------|
| `--tls` | `false` | 自己署名証明書でTLSを有効化 |
| `--tls-cert` | `""` | TLS証明書ファイルのパス |
| `--tls-key` | `""` | TLSキーファイルのパス |

```bash
swag2mcp-mock
swag2mcp-mock --tls
```

**ワークフロー：**
1. `mock_enabled: true` と `base_mock_url` を設定に追加
2. モックサーバーを起動：`swag2mcp-mock`
3. MCPサーバーを起動：`swag2mcp mcp` — invokeは `base_url` の代わりに `base_mock_url` を使用します
4. 認証は自動的に適用されます：OAuth2/Digestはポート9090/9091のモックサーバーを使用し、その他のタイプは認証情報を直接適用します

### モック認証

仕様で `auth` が設定されている場合、MCPサーバーが自動的に認証を適用します。
専用のモックサーバーが必要なのは次の2つの認証タイプのみです：

| 認証タイプ | モックエンドポイント | 動作 |
|------------|---------------------|------|
| `oauth2-cc` / `oauth2-pwd` | ポート9090の `POST /token` | 任意の `client_id`/`username`+`password` を受け付け、`{"access_token":"<random>","token_type":"Bearer","expires_in":3600}` を返す |
| `digest` | ポート9091の `GET /` | `algorithm=MD5` の401チャレンジを送信、任意のDigestレスポンスを受け付け、`{"status":"authenticated","method":"digest"}` を返す |

その他の認証タイプ（`basic`、`bearer`、`api-key`、`hmac`、`script`）はモックサーバーを
**必要としません** — MCPサーバーが設定された認証情報を各リクエストに自動的に適用します。

---

## インテグレーション

swag2mcpはModel Context Protocol (MCP) に対応しており、MCP互換クライアントで動作します。

### ローカル (stdio) — 同じマシンのエージェント

サーバーを起動：

```bash
swag2mcp mcp
```

| クライアント | 設定ファイル | 内容 |
|------------|-------------|------|
| **OpenCode** | `opencode.json` | `{"mcp":{"swag2mcp":{"type":"local","command":["swag2mcp","mcp"]}}}` |
| **Cursor** | `.cursor/mcp.json` | `{"mcpServers":{"swag2mcp":{"command":"swag2mcp","args":["mcp"]}}}` |
| **Claude Desktop** | `claude_desktop_config.json` | `{"mcpServers":{"swag2mcp":{"command":"swag2mcp","args":["mcp"]}}}` |
| **VS Code** | `.vscode/mcp.json` | `{"servers":{"swag2mcp":{"type":"stdio","command":"swag2mcp","args":["mcp"]}}}` |
| **Crush** | `crush.json` | `{"mcp":{"swag2mcp":{"type":"stdio","command":"swag2mcp","args":["mcp"]}}}` |

### リモート (HTTP) — クラウド / 別のマシンのエージェント

HTTPトランスポートでサーバーを起動：

```bash
swag2mcp mcp --transport streamable-http --http-addr :8080 --auth-token my-secret
```


> **Note:** If you initialized the workspace at a custom path (e.g. `swag2mcp init ./my-project`), you must specify the path when starting the MCP server: `swag2mcp mcp ./my-project`. The IDE configuration must also use the full path to the config file.

または `swag2mcp.yaml` で設定：

```yaml
mcp:
  transport: streamable-http
  addr: ":8080"
  path: "/mcp"
  auth_token: $(MCP_AUTH_TOKEN)
```

| クライアント | 設定ファイル | 内容 |
|------------|-------------|------|
| **OpenCode** | `opencode.json` | `{"mcp":{"swag2mcp":{"type":"remote","url":"http://localhost:8080/mcp","headers":{"Authorization":"Bearer ${MCP_AUTH_TOKEN}"}}}}` |
| **VS Code** | `.vscode/mcp.json` | `{"servers":{"swag2mcp":{"type":"http","url":"http://localhost:8080/mcp"}}}` |

> **ヘルスチェック**（MCPハンドシェイク不要）：
> ```bash
> curl http://localhost:8080/health
> # → {"status":"ok","version":"v1.1.3"}
> ```

---

## MCPサーバー

MCPサーバーは、stdioまたはHTTPトランスポートを介して16のツールを公開します。LLMエージェント（Opencode、Cursor、Claude、Copilotなど）は、設定後に自動的に接続します。

### ツール階層

```
spec_list                       — 利用可能なすべての仕様を一覧表示
  └─ spec_by_id                 — IDで仕様の詳細を取得
       └─ collection_by_spec    — 仕様内のコレクション
            └─ tag_by_collection     — コレクション内のタグ
                 └─ endpoint_by_tag  — タグ内のエンドポイント
                      └─ inspect          — 完全なOpenAPI操作
                           └─ invoke       — API呼び出しを実行

search                          — 全エンドポイントの全文検索
```

### ツールリファレンス

| ツール | 引数 | 戻り値 | 説明 |
|--------|------|--------|------|
| `spec_list` | — | `Spec[]` | 利用可能なすべての仕様 |
| `spec_by_id` | `id` | Spec + Collections | 仕様の詳細 |
| `collection_by_spec` | `specId` | Collections | 仕様内のコレクション |
| `collection_by_id` | `id` | Collection + Tags | コレクションの詳細 |
| `tag_by_collection` | `collectionId` | Tags | コレクション内のタグ |
| `tag_by_spec` | `specId` | Tags | 仕様内のすべてのタグ |
| `tag_by_id` | `id` | Tag | 単一タグのメタデータ |
| `endpoint_by_tag` | `tagId` | Endpoints | タグ内のエンドポイント |
| `endpoint_by_collection` | `collectionId` | Endpoints | コレクション内の全エンドポイント |
| `endpoint_by_spec` | `specId` | Endpoints | 仕様内の全エンドポイント |
| `endpoint_by_id` | `id` | Endpoint | エンドポイントの概要 |
| `search` | `query`, `limit` | Endpoints | 全文検索 |
| `inspect` | `endpointId` | Full Operation | 完全なOpenAPI操作オブジェクト |
| `invoke` | `endpointId`, `parameters`, `requestBody` | Response | 実際のAPI呼び出しを実行 |
| `auth` | `specId` | Token | 仕様の認証トークンを取得 |

---

## 検索

### クエリ構文

| 機能 | 構文 | 例 |
|------|------|-----|
| 用語 | `用語` | `ペット` |
| フレーズ | `"フレーズ"` | `"ペットを追加"` |
| フィールド: method | `method:用語` | `method:post` |
| フィールド: tag | `tag:用語` | `tag:auth` |
| フィールド: path | `path:用語` | `path:/users` |
| フィールド: summary | `summary:用語` | `summary:login` |
| 必須 (AND) | `+用語` | `+method:post +tag:user` |
| 除外 (NOT) | `-用語` | `-deprecated` |
| ワイルドカード | `*` | `path:*/v2/*` |
| あいまい | `用語~` | `watex~` |
| 正規表現 | `/パターン/` | `/user(s\|sessions)/` |
| ブースト | `用語^N` | `tag:pet^5` |
| すべて一致 | `*` | `*` |

### 例

```
# authタグ内のPOSTエンドポイントを検索
+method:post +tag:auth

# ログイン関連のエンドポイントを検索
summary:"login"~

# ユーザー関連のパスをすべて検索、非推奨を除外
path:*/users/* -deprecated

# 複合クエリ
+method:get +tag:pet summary:"find by status"
```

### インデックス化されたフィールド

| フィールド | タイプ | 内容 |
|-----------|--------|------|
| `method` | text | HTTPメソッド（小文字） |
| `tag` | text | タグ名（小文字） |
| `path` | text | APIパス（小文字） |
| `summary` | text（分析済み） | エンドポイントの概要/説明（小文字） |
| `_all` | text（分析済み） | method + path + tag + summary |

---

## ワークスペース

### ディレクトリ構造

```
~/.swag2mcp/                    # または {プロジェクト}/.swag2mcp/
├── swag2mcp.yaml               # 設定ファイル
├── cache/                      # キャッシュされたリモート仕様
│   ├── {hash}.spec             # 仕様ファイルの内容
│   └── {hash}.meta             # JSONメタデータ
├── specs/                      # ローカル仕様ファイル（ユーザー管理）
├── responses/                  # 呼び出し応答ファイル
└── auth_scripts/               # 認証スクリプト
```

### パス解決

```
swag2mcp <command>          → ~/.swag2mcp/
swag2mcp <command> ./       → {cwd}/.swag2mcp/
swag2mcp <command> path/to  → {cwd}/path/to/.swag2mcp/
```

### .gitignore

一時データのみを無視する必要があります：

```
.swag2mcp/cache/*
.swag2mcp/responses/*
```

設定ファイル `.swag2mcp/swag2mcp.yaml` と `.swag2mcp/specs/` 内のspecファイルは**リポジトリに含める必要があります**。

### 推奨

すべてのspecファイルを `.swag2mcp/specs/` に保存してください——キャッシュにコピーされず直接使用されることを保証する唯一の方法です。

---

## キャッシュ

### ルール

| ソース | 動作 |
|--------|------|
| HTTP/HTTPS URL | 常にキャッシュ。TTL：ランダム1-48時間。 |
| `specs/` 内のローカルパス | 直接使用、キャッシュされない。 |
| `specs/` 外のローカルパス | 初回アクセス時にキャッシュにコピー。 |
| `file://` URL | ローカルパスとして扱われる。 |

### キャッシュキー

正規化された場所のSHA-256ハッシュ（最初の16バイト = 32 hex文字）。

### キャッシュヒットロジック

1. `.meta` ファイルを読み取り——期限切れまたは欠落 → ミス
2. ローカルソースの場合：`ModTime` が変更 → ミス
3. `.spec` ファイルが欠落 → ミス
4. それ以外 → ヒット

---

## 開発

```bash
# ビルド
go build ./cmd/swag2mcp/

# テスト
go test ./...

# リンター
make lint

# 実行
go run ./cmd/swag2mcp/main.go
```
