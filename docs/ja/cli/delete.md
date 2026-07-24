# delete

## 目的

設定から **spec**（API サービス）または **collection**（spec ファイル）を削除します。`add` の逆の操作です。

## 使用するタイミング

- API が不要になった場合
- spec から特定の spec ファイルを削除したい場合
- ワークスペースをクリーンアップしている場合

## 構文

```bash
swag2mcp delete spec [path]
swag2mcp delete collection [path]
```

## 引数

| 引数 | 位置 | 必須 | 説明 |
|------|------|------|------|
| `path` | 1 | いいえ | ワークスペースディレクトリ。省略時はパス解決ルールに従います。 |

## フラグ

なし。両方のサブコマンドは純粋に対話型です。

## 仕組み

### Spec の削除

リストから spec を選択するよう促し、削除前に確認を求めます。

```bash
swag2mcp delete spec
```

### Collection の削除

spec を選択し、次にその spec 内の collection を選択するよう促し、確認を求めます。

```bash
swag2mcp delete collection
```

## ID の確認

対話型プロンプトは ID ではなく人間が読める名前を表示します。ID が必要な場合：

```bash
# すべての spec を ID 付きで一覧表示
swag2mcp ls

# 特定の spec の collection を一覧表示
swag2mcp ls --tags
```

## コマンド実行後の確認

```bash
swag2mcp ls [path]
# 削除された spec または collection が表示されなくなるはずです
```

## ニュアンス

- **TTY 必須:** 両方のコマンドは対話型ターミナルが必要です。CI/CD パイプライン、cron ジョブ、または非対話型スクリプトでは**動作しません**。
- **`--force` や `--yes` なし:** 確認プロンプトをスキップする方法はありません。これは誤削除を防ぐための意図的な設計です。
- **自動初期化:** 設定ファイルが存在しない場合、`delete` は自動的に init ウィザードを実行します。
- **YAML モードなし:** `add` とは異なり、`--yaml` フラグはありません。削除は常に対話型です。
