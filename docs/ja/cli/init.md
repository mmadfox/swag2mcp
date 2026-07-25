# init

## 目的

`init` コマンドは **ワークスペース** — `swag2mcp.yaml` 設定ファイルとキャッシュ、spec、レスポンス、認証スクリプト用のサブディレクトリを持つディレクトリを作成します。swag2mcp をセットアップする際に最初に実行するコマンドです。

## 使用するタイミング

- swag2mcp を初めてセットアップする場合
- 特定のディレクトリに新しいワークスペースを作成したい場合
- 破損した、または存在しないワークスペースを再初期化する必要がある場合

## 構文

```bash
swag2mcp init [path] [flags]
```

## 引数

| 引数 | 位置 | 必須 | 説明 |
|------|------|------|------|
| `path` | 1 | いいえ | ワークスペースディレクトリ。省略時はデフォルトで `~/.swag2mcp`。 |

## フラグ

| フラグ | 省略形 | 型 | デフォルト | 説明 |
|-------|--------|-----|-----------|------|
| `--interactive` | `-i` | `bool` | `false` | 対話型 TUI ウィザードを実行 |
| `--force` | `-f` | `bool` | `false` | 空でないディレクトリの既存設定を上書き |

## 仕組み

### 非対話モード（デフォルト）

spec なしの最小限の `swag2mcp.yaml` を作成します。後で手動でファイルを編集します。

```bash
swag2mcp init
# ~/.swag2mcp/swag2mcp.yaml を作成

swag2mcp init ./my-project
# ./my-project/swag2mcp.yaml を作成

swag2mcp init /absolute/path
# /absolute/path/swag2mcp.yaml を作成
```

### 対話モード（`-i`）

18 ステップの TUI ウィザードを起動し、以下をガイドします：

1. ワークスペースディレクトリの選択
2. ドメイン、タイトル、ベース URL での spec 追加
3. location URL での collection 設定
4. 認証の設定（全 9 方式）
5. HTTP クライアント設定（タイムアウト、プロキシ、ヘッダーなど）

```bash
swag2mcp init -i
```

### 強制モード（`--force`）

デフォルトでは、`init` は空でないディレクトリでの実行を拒否します。上書きするには `--force` を使用します：

```bash
swag2mcp init -f
swag2mcp init ./existing-dir -f
```

## 作成されるもの

```
~/.swag2mcp/
├── swag2mcp.yaml       # 設定ファイル
├── cache/               # ダウンロードされたリモート spec ファイル
├── specs/               # ローカル spec ファイル
├── responses/           # 保存された API 呼び出しレスポンス
└── auth_scripts/        # 認証スクリプト（ScriptAuth タイプ用）
```

## コマンド実行後の確認

```bash
ls ~/.swag2mcp/swag2mcp.yaml
# ファイルが存在すれば、init は成功
```

## ニュアンス

- **パス解決:** `[path]` はファイルパスではなく**ワークスペースディレクトリ**です。CLI は自動的に `swag2mcp.yaml` を追加します。解決順序：明示的な `[path]` → カレントディレクトリ（`./`）→ `~/.swag2mcp/`。
- **空でないディレクトリのチェック:** `--force` なしでは、ターゲットディレクトリが存在し空でない場合、`init` はエラーを返します。これにより誤った上書きを防ぎます。
- **認証スクリプトのスタブ:** いずれかの spec が `ScriptAuth` を使用する場合、`init` は `auth_scripts/` にスタブスクリプトファイル（Unix では `.sh`、Windows では `.bat`）を作成します。
- **出力:** 成功時に、設定パスとヒントを表示：`"Next step: edit swag2mcp.yaml or run 'swag2mcp ls' to list configured specs"`。
