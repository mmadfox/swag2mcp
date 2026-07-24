# version

## 目的

swag2mcp のバージョンを表示します。インストールされたバージョンの確認、バグ報告、互換性の確認に便利です。

## 使用するタイミング

- インストールされている swag2mcp のバージョンを確認したい場合
- バグを報告する際にバージョンを含める必要がある場合
- インストールが成功したことを確認したい場合

## 構文

```bash
swag2mcp version
swag2mcp --version
```

## 引数

なし。

## フラグ

なし。

## 仕組み

```bash
swag2mcp version
# swag2mcp v1.2.0

swag2mcp --version
# swag2mcp v1.2.0
```

## 出力形式

```
swag2mcp &lt;version&gt;
```

バージョンはビルド時に `ldflags` で設定されます。設定されていない場合、デフォルトで `"dev"` になります。

## ニュアンス

- **2 つの形式:** `swag2mcp version`（サブコマンド）と `swag2mcp --version`（グローバルフラグ）は同じ出力を生成します。
- **設定不要:** このコマンドはワークスペースや設定ファイルなしで動作します。
