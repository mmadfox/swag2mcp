# エクスポートとインポート

## 概要

swag2mcp は ZIP アーカイブによる完全なワークスペースのラウンドトリップをサポートしています。ワークスペース全体（設定、spec ファイル、認証スクリプト）を ZIP ファイルにエクスポートし、別のマシンで復元できます。

## エクスポート

ワークスペースのポータブルな ZIP バックアップを作成します。

```bash
# デフォルトファイルにエクスポート（swag2mcp-backup-<timestamp>.zip）
swag2mcp export

# カスタムパスでエクスポート
swag2mcp export --output ~/backups/swag2mcp-backup.zip

# 特定の spec のみエクスポート
swag2mcp export --spec meteo
swag2mcp export --spec meteo,store
```

### エクスポートに含まれるもの

| 項目 | 説明 |
|------|------|
| `swag2mcp.yaml` | 設定ファイル |
| `specs/` | すべての spec ファイル（OpenAPI/Swagger/Postman） |
| `auth_scripts/` | 認証スクリプト |
| `swag2mcp.meta` | メタデータ（互換性のためのバージョン情報） |

キャッシュとレスポンスは**エクスポートされません** — これらは一時的なデータであり、復元時には古くなっています。

### デフォルトのファイル名

出力パスを指定しない場合、ファイルはカレントディレクトリに `swag2mcp-backup-<YYYY-MM-DD-HHMMSS>.zip` として保存されます（UTC タイムスタンプ）。

## インポート

ZIP バックアップからワークスペースを復元するか、spec ファイルをインポートします。

### ZIP からの復元

```bash
# ワークスペース全体を復元
swag2mcp import --from-zip /path/to/backup.zip

# 上書きして復元
swag2mcp import --from-zip /path/to/backup.zip -f
```

ZIP は `swag2mcp export` で作成されたものである必要があります — 任意の ZIP ファイルは機能しません。

### 単一の spec ファイルをインポート

spec ファイルをダウンロードしてワークスペースに追加します：

```bash
swag2mcp import https://example.com/spec.yaml myspec
swag2mcp import /path/to/workspace https://example.com/spec.yaml myspec
```

### 既存の設定からの一括インポート

指定された spec（ドメイン）のすべての collection spec ファイルをダウンロードします：

```bash
swag2mcp import --spec meteo
swag2mcp import /path/to/workspace --spec meteo,store
```

各 collection の spec ファイルをダウンロードし、`specs/` に保存し、設定をローカルコピーを指すように更新します。

## ユースケース

### バックアップ

```bash
swag2mcp export --output swag2mcp-$(date +%Y-%m-%d).zip
```

### 別のマシンへの転送

```bash
# 古いマシンで
swag2mcp export --output swag2mcp.zip

# ZIP を新しいマシンにコピーし、次を実行：
swag2mcp import --from-zip swag2mcp.zip
```

### 設定の共有

```bash
swag2mcp init
swag2mcp export --output template.zip
# template.zip を同僚と共有
```

## エクスポート後の確認

ZIP ファイルが作成されたことを常に確認してください：

```bash
ls -la swag2mcp-backup-*.zip
```

## 重要な注意点

- **出力は `.zip` で終わるファイルパスである必要があります** — ディレクトリを渡さないでください
- **キャッシュとレスポンスは除外されます** — 設定、spec、認証スクリプトのみが保存されます
- **ZIP は自己完結型です** — swag2mcp がインストールされた任意のマシンで復元できます
- **Spec フィルター** — `--spec` を使用して特定の spec のみをエクスポートまたはインポートします
