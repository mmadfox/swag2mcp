# Workspace

ワークスペースは、swag2mcp がすべてのデータ（設定、キャッシュされた spec、ローカル spec ファイル、保存されたレスポンス、認証スクリプト）を保存するディレクトリです。

## 構造

```
~/.swag2mcp/                          # ワークスペースルート（デフォルト）
├── swag2mcp.yaml                     # 設定ファイル
├── cache/                            # キャッシュされたリモート spec ファイル
│   ├── a1b2c3d4e5f6...spec          # キャッシュされた spec の内容
│   └── a1b2c3d4e5f6...meta          # キャッシュメタデータ（JSON）
├── specs/                            # ローカル spec ファイル
│   └── my-api.yaml
├── responses/                        # 保存された API レスポンス（大きなレスポンス）
│   ├── meteo-get-forecast-abc123.json
│   └── response-fragment-def456.json
└── auth_scripts/                     # 認証スクリプト
    ├── meteo.sh                      # Unix シェルスクリプト
    └── meteo.bat                     # Windows バッチスクリプト
```

## デフォルトパス

- **Linux/macOS**: `~/.swag2mcp/`
- **Windows**: `%USERPROFILE%\.swag2mcp\`

## カスタムパス

```bash
swag2mcp mcp /path/to/workspace
swag2mcp mcp ./my-workspace
```

## ディレクトリ

### cache/

ダウンロードされたリモート spec ファイルを保存します。各ファイルは URL の SHA-256 ハッシュをファイル名としてキャッシュされます：

- `{hash}.spec` — キャッシュされた spec ファイルの内容
- `{hash}.meta` — JSON メタデータ（ソース URL、キャッシュ時間、TTL）

各キャッシュファイルには 1 時間から 48 時間の間でランダムな TTL があります。キャッシュは起動のたびに自動的にチェックされます — 有効な（期限切れでない）エントリが存在する場合、ダウンロードせずに再利用されます。

**コマンド:**
- `swag2mcp update` — キャッシュをクリアし、すべての spec を再ダウンロード
- `swag2mcp clean` — キャッシュとレスポンスをクリア

### specs/

collection が `location: specs/{name}` を介して指すローカル spec ファイルを保存します。ここにあるファイルはキャッシュなしで直接使用されます。

このディレクトリは以下によって作成されます：
- `swag2mcp import &lt;source&gt; &lt;name&gt;` — リモート spec をダウンロードしてここに保存
- `swag2mcp export` — ここからエクスポート ZIP に spec をコピー
- 手動配置 — 自分で spec ファイルをコピーできます

### responses/

`max_response_size` 制限（デフォルト 1 MB）を超える API レスポンスを保存します。LLM がエンドポイントを呼び出し、レスポンスが大きすぎる場合、swag2mcp はここに保存し、代わりにファイル参照を返します。

命名規則：`{domain}-{method}-{path_with_underscores}-{6char_hex}.json`

古いレスポンスは MCP サーバー起動後 48 時間で自動的にクリーンアップされます。

### auth_scripts/

`script` 認証タイプの認証スクリプトを保存します。各スクリプトは spec のドメインにちなんで命名されます。

#### 命名規則

| プラットフォーム | ファイル名 | 例 |
|-------------|----------|-----|
| Unix（Linux、macOS） | `{domain}.sh` | `meteo.sh` |
| Windows | `{domain}.bat` | `meteo.bat` |

ドメインに `/` または `\` 文字を含めることはできません。

#### スクリプトの仕組み

1. swag2mcp は 30 秒のタイムアウトでスクリプトを実行します
2. スクリプトは有効な JSON を stdout に出力する必要があります
3. swag2mcp は JSON を解析し、API リクエストにトークンを使用します

#### 期待される出力形式

```json
{
  "token": "your-token-here",
  "expires_in": 3600
}
```

| フィールド | 型 | 必須 | 説明 |
|-----------|------|------|------|
| `token` | string | ✅ | 認証トークン |
| `access_token` | string | ❌ | `token` の代替（最初にチェック） |
| `token_type` | string | ❌ | トークンタイプ（例："Bearer"） |
| `expires_in` | number | ❌ | トークンの有効期間（秒）（デフォルト：3600） |

#### 実行

| プラットフォーム | コマンド |
|-------------|---------|
| Unix | `sh {domain}.sh` |
| Windows | `cmd /c {domain}.bat` |

#### トークンキャッシュ

トークンは期限が切れるまでメモリ内にキャッシュされます。API 呼び出しのたびに、swag2mcp は最初にキャッシュをチェックします — スクリプトはキャッシュされたトークンが期限切れになった場合のみ実行されます。

#### スタブの作成

`auth: { type: script, config: { domain: "myapi" } }` を設定すると、swag2mcp は自動的にスタブスクリプトを作成します：

**Unix（`auth_scripts/myapi.sh`）:**
```bash
#!/bin/sh
echo '{"token": "your-token-here", "expires_in": 3600}'
```

**Windows（`auth_scripts/myapi.bat`）:**
```bat
@echo off
echo {"token": "your-token-here", "expires_in": 3600}
```

プレースホルダートークンを実際の認証ロジックに置き換えてください。

#### 孤立スクリプトのクリーンアップ

spec を削除すると、その認証スクリプトは孤立します。swag2mcp は以下で自動的に孤立スクリプトを削除します：
- `swag2mcp update`
- `swag2mcp clean`

## コマンド

### update

```bash
swag2mcp update [path]
```

設定を検証し、キャッシュとレスポンスをクリアし、すべての spec ファイルを再ダウンロードします。また、認証スクリプトの存在を確認し、孤立スクリプトを削除します。

このコマンドは以下を行った後に使用します：
- collection の追加または削除
- collection の location の変更
- 再キャッシュが必要な spec ファイルの編集

### clean

```bash
swag2mcp clean [path]
```

`cache/` と `responses/` のすべての内容と、孤立した認証スクリプトを削除します。spec の再キャッシュは**行いません** — それには `update` を使用してください。

### validate

```bash
swag2mcp validate [path]
```

すべての collection location を含む設定を検証します。[CLI: validate](../cli/validate.md) を参照してください。

## エクスポートとインポート

```bash
# ワークスペースを ZIP にエクスポート（デフォルト名：swag2mcp-backup-{date}.zip）
swag2mcp export

# 特定のパスにエクスポート
swag2mcp export /path/to/workspace /path/to/backup.zip

# 特定の spec のみをエクスポート
swag2mcp export --spec meteo

# バックアップから復元
swag2mcp import --from-zip /path/to/backup.zip
swag2mcp import /path/to/workspace /path/to/backup.zip
```

エクスポートに含まれるもの：`swag2mcp.yaml`、`specs/`、`auth_scripts/`。キャッシュとレスポンスは除外されます（これらはローカルデータです）。

## .gitignore

ワークスペースが Git リポジトリ内にある場合、以下のエントリを `.gitignore` に追加します：

```gitignore
# swag2mcp — ローカルデータのみ
.swag2mcp/cache/
.swag2mcp/responses/
```

`cache/` と `responses/` ディレクトリには、コミットすべきでないローカルでマシン固有のデータが含まれています。その他（`swag2mcp.yaml`、`specs/`、`auth_scripts/`）は、設定をチーム全体で共有するためにリポジトリに含めるべきです。
