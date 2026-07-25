# import

## 目的

spec ファイルをワークスペースにインポートするか、ZIP バックアップからワークスペース全体を復元します。3 つのモードが異なるシナリオをカバーします：単一 spec の追加、既存設定からの一括インポート、または完全なワークスペースの復元。

## 使用するタイミング

- spec URL またはファイルがあり、ワークスペースに追加したい場合
- 設定で参照されているすべての spec ファイルをダウンロードしたい場合
- `export` で作成された ZIP バックアップからワークスペースを復元する必要がある場合
- swag2mcp を別のマシンに移行する場合

## 構文

```bash
swag2mcp import [path] [source] [name] [flags]
```

## 引数

| 引数 | 位置 | 必須 | 説明 |
|------|------|------|------|
| `path` | 1 | いいえ | ワークスペースディレクトリ。省略時はパス解決ルールに従います。 |
| `source` | 2 | 場合による | spec ファイルの URL またはローカルパス、または ZIP アーカイブへのパス |
| `name` | 3 | 場合による | 新しい spec のドメイン名 |

## フラグ

| フラグ | 省略形 | 型 | デフォルト | 説明 |
|-------|--------|-----|-----------|------|
| `--spec` | `-s` | `stringSlice` | `nil` | 指定された spec から collection をインポート（カンマ区切り） |
| `--from-zip` | | `string` | `""` | swag2mcp バックアップ ZIP からワークスペースを復元 |

## 仕組み

### モード 1 — URL またはファイルからの単一インポート

spec ファイルをダウンロードし、ドメイン名を付けてワークスペースに追加します：

```bash
swag2mcp import https://example.com/spec.yaml myspec
swag2mcp import /path/to/workspace https://example.com/spec.yaml myspec
swag2mcp import ./local-spec.yaml myspec
```

spec ファイルは `specs/` に保存され、設定が新しい spec エントリで更新されます。

### モード 2 — 既存設定からの一括インポート

指定されたドメインのすべての collection を、設定された URL からダウンロードします：

```bash
swag2mcp import --spec meteo
swag2mcp import /path/to/workspace --spec meteo,store
```

各 collection の spec ファイルがダウンロードされ、`specs/` に保存されます。設定はローカルコピーを指すように更新されます。

### モード 3 — ZIP バックアップからの復元

`swag2mcp export` で作成された ZIP アーカイブからワークスペース全体を復元します：

```bash
swag2mcp import --from-zip /path/to/backup.zip
swag2mcp import /path/to/workspace /path/to/backup.zip
```

> **ZIP は `swag2mcp export` で作成されたものである必要があります。** 任意の ZIP ファイルは機能しません — アーカイブは特定の内部構造（`swag2mcp.yaml`、`specs/`、`auth_scripts/`）を持っています。

## コマンド実行後の確認

```bash
# 単一または一括インポート
swag2mcp ls [path]
# 新しい spec がリストに表示されるはずです

# ZIP 復元
swag2mcp ls [path]
# バックアップのすべての spec が表示されるはずです
```

## ニュアンス

- **一括モードには設定が必要:** `--spec` を使用する場合、設定ファイルが存在する必要があります。必要に応じて最初に `init` を実行してください。
- **単一インポートはワークスペースを作成:** ワークスペースが存在しない場合、自動的に作成されます。
- **ZIP 検出:** `.zip` で終わる位置引数は ZIP ソースとして扱われます。`--from-zip` フラグは位置検出より優先されます。
- **`--force`:** ZIP 復元時に既存のワークスペースを上書きするために使用できます。
- **HTTP クライアント:** 設定からのグローバル HTTP クライアント設定がインポート中に適用されます（タイムアウト、プロキシ、ヘッダーなど）。
