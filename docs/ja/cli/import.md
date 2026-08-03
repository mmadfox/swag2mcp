# import

## 目的

spec ファイルをワークスペースの `specs/` ディレクトリにローカル保存するか、ZIP バックアップからワークスペース全体を復元します。3 つのモードが異なるシナリオをカバーします：単一 spec の追加、既存設定からの一括インポート、または完全なワークスペースの復元。

## 使用するタイミング

- spec URL またはファイルがあり、ワークスペースにローカル保存したい場合
- 設定内のすべての collection spec ファイルをダウンロードして、ワークスペースを自己完結型にしたい場合
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
| `name` | 3 | 場合による | 保存するファイル名（例：`example-api.yaml`）。省略時は URL から自動生成されます。 |

## フラグ

| フラグ | 省略形 | 型 | デフォルト | 説明 |
|-------|--------|-----|-----------|------|
| `--spec` | `-s` | `string` | `""` | 設定から collection spec ファイルをダウンロード。値なしですべての spec、または `--spec meteo,github` のようにドメインを指定 |
| `--force` | `-f` | `bool` | `false` | 既存の spec ファイルをエラーなしで上書き |
| `--from-zip` | | `string` | `""` | swag2mcp バックアップ ZIP からワークスペースを復元 |

## 仕組み

### モード 1 — URL またはファイルからの単一インポート

spec ファイルをダウンロードして `specs/` に保存します：

```bash
swag2mcp import https://example.com/spec.yaml example-api.yaml
swag2mcp import /path/to/workspace https://example.com/spec.yaml example-api.yaml
swag2mcp import ./local-spec.yaml example-api.yaml
```

`name` を省略すると、URL のファイル名から自動生成されます：
```bash
swag2mcp import https://example.com/specs/petstore.yaml
# → petstore.yaml として保存
```

既存のファイルを `--force` で上書き：
```bash
swag2mcp import --force https://example.com/spec.yaml example-api.yaml
```

インポート後、出力にはワークスペースパス、保存されたファイル、`swag2mcp.yaml` に追加する YAML テンプレートが表示されます：

```
✅ Imported to /path/to/workspace
   specs/example-api.yaml

   Add to swag2mcp.yaml:
     specs:
       - domain: <your-domain>
         collections:
           - location: specs/example-api.yaml
```

### モード 2 — 既存設定からの一括インポート（`--spec`）

指定されたドメインのすべての collection spec ファイルを、設定された `location` URL からダウンロードし、`specs/` に保存して、設定をローカルコピーを指すように更新します：

```bash
swag2mcp import --spec                # すべての spec
swag2mcp import --spec meteo           # 特定の spec
swag2mcp import --spec meteo,github    # 複数の spec
swag2mcp import /path/to/workspace --spec meteo
```

指定されたドメインが設定に存在しない場合、コマンドはエラーを返します：
```
Error: import_no_match
  Spec "nonexistent" not found in config.
```

これによりワークスペースが自己完結型になります — インポート後はリモート spec URL は不要です。

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
- **HTTP クライアント:** 設定からのグローバル HTTP クライアント設定がインポート中に適用されます（タイムアウト、プロキシ、ヘッダーなど）。
