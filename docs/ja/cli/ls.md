# ls

## 目的

設定されているすべての **spec** とその **collection** を人間が読める形式で一覧表示します。ワークスペースで利用可能な API を確認するための主要な方法です。

## 使用するタイミング

- 設定されている API を確認したい場合
- spec または collection の ID を見つける必要がある場合
- 各 collection のエンドポイント数を確認したい場合
- タグで spec をフィルタリングしたい場合

## 構文

```bash
swag2mcp ls [path] [flags]
```

## 引数

| 引数 | 位置 | 必須 | 説明 |
|------|------|------|------|
| `path` | 1 | いいえ | ワークスペースディレクトリ。省略時はパス解決ルールに従います。 |

## フラグ

| フラグ | 省略形 | 型 | デフォルト | 説明 |
|-------|--------|-----|-----------|------|
| `--tags` | `-t` | `string` | `""` | タグで spec をフィルタリング（カンマ区切り） |

## 仕組み

### すべての spec を一覧表示

各 spec をドメイン、collection、エンドポイント数とともに表示します：

```bash
swag2mcp ls
```

出力例：

```
Specifications:
  dadjoke (https://icanhazdadjoke.com)
    jokes (3 endpoints)
  meteo (https://meteo.swagger.io/v2)
    forecast (5 endpoints)
    current (8 endpoints)
  binance (https://api.binance.com)
    market-data (12 endpoints)
```

### タグでフィルタリング

指定されたタグを持つ spec のみを表示します：

```bash
swag2mcp ls --tags=public
swag2mcp ls --tags=public,internal
```

## コマンド実行後の確認

`add`、`delete`、`update`、`import` の後に `ls` を使用して、ワークスペースの状態が期待通りであることを確認します。

## ニュアンス

- **自動初期化:** 設定ファイルが存在しない場合、`ls` は自動的に init ウィザードを実行します。
- **タグフィルタリング:** タグはカンマ区切りです。指定されたタグの**いずれか**に一致する spec が表示されます（OR 論理）。
- **出力形式:** 出力はプレーンテキストで、JSON ではありません。機械可読な出力には `info` を使用してください。
