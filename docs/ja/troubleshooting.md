# トラブルシューティング

## インストールの問題

### swag2mcp: command not found

バイナリが PATH にありません。

```bash
# Go がインストールされているか確認
go version

# Go がバイナリをインストールする場所を確認
go env GOPATH
# 通常は ~/go または ~/go/bin

# PATH に追加（~/.zshrc または ~/.bashrc に追加）
export PATH=$PATH:$(go env GOPATH)/bin

# またはフルパスを使用
~/go/bin/swag2mcp --version
```

GitHub Releases からバイナリをダウンロードした場合、PATH が通ったディレクトリにあることを確認してください：

```bash
# /usr/local/bin に移動（macOS/Linux）
sudo mv swag2mcp /usr/local/bin/
```

### permission denied

バイナリに実行権限がありません。

```bash
# go install の場合（所有権を修正）
sudo chown -R $(whoami) $(go env GOPATH)

# ダウンロードしたバイナリの場合
chmod +x /path/to/swag2mcp
```

### Go のバージョンが古すぎる

swag2mcp には Go 1.23+ が必要です。

```bash
go version
# バージョンが 1.23 未満の場合、Go を更新：
# https://go.dev/dl/
```

### モックサーバーが見つからない

モックサーバーは別のバイナリです。明示的にインストールしてください：

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

## 設定の問題

### 設定ファイルが見つからない

swag2mcp が `swag2mcp.yaml` を見つけられません。

```bash
# 新しい設定を作成
swag2mcp init

# またはパスを明示的に指定
swag2mcp mcp /path/to/workspace
swag2mcp ls /path/to/workspace
```

**よくある原因：** ランダムなディレクトリから `swag2mcp mcp` を実行し、プロジェクトのワークスペースではなく `~/.swag2mcp/` を探しました。常にパスを明示的に指定してください。

### 間違ったワークスペースが読み込まれた

期待とは異なるワークスペースが読み込まれました。

**解決順序：** 明示的な `[path]` → カレントディレクトリ（`./`）→ `~/.swag2mcp/`。パスなしで `swag2mcp mcp` を `swag2mcp.yaml` のないディレクトリから実行すると、`~/.swag2mcp/` にフォールバックします。

**修正：** 常にワークスペースパスを指定：`swag2mcp mcp /path/to/your/workspace`

### YAML 解析エラー

設定ファイルに無効な YAML 構文があります。

```bash
# 設定を検証
swag2mcp validate

# よくある間違い：
# - スペースの代わりにタブ（YAML はスペースが必要）
# - ネストされたフィールドのインデント不足
# - 特殊文字を含む引用符なしの文字列（: # & {）
```

**ヒント：** YAML リンターまたは YAML 対応エディターを使用して構文エラーを発見してください。

### 検証エラー：「no specifications defined」

設定ファイルは存在しますが、spec がありません。

```bash
# spec を追加
swag2mcp add spec

# または swag2mcp.yaml を編集して少なくとも 1 つの spec を追加
```

### 検証エラー：「duplicate domain」

2 つの spec が同じ `domain` 値を持っています。ドメインは一意である必要があります。

```bash
# 現在の spec を一覧表示
swag2mcp ls

# swag2mcp.yaml で重複ドメインを確認
```

### 検証エラー：「invalid spec location」

`location` URL またはファイルパスにアクセスできないか、有効な spec ファイルではありません。

```bash
# URL にアクセスできるか確認
curl -I https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml

# ローカルファイルが存在するか確認
ls -la ./specs/my-api.yaml

# ファイルが有効な OpenAPI/Swagger/Postman か確認
# （単なる JSON や HTML ページではない）
```

**よくある原因：** `location` フィールドが spec ファイル URL ではなく API エンドポイント自体（例：`https://api.example.com/v1/users`）を指しています。location は OpenAPI/Swagger/Postman ファイルを指す必要があります。

## MCP サーバーの問題

### ポートが既に使用中

別のプロセスがポートを使用しています。

```bash
# プロセスを特定
lsof -i :8080

# 強制終了
kill <PID>

# または別のポートを使用
swag2mcp mcp --transport sse --http-addr :9090
```

### 接続が拒否された

MCP サーバーが実行されていないか、到達できません。

```bash
# サーバーが実行中であることを確認
swag2mcp mcp --transport sse --http-addr 127.0.0.1:8080

# 別のターミナルでヘルスエンドポイントを確認
curl http://127.0.0.1:8080/health

# カスタムパスを使用している場合
curl http://127.0.0.1:8080/custom-path/health
```

### MCP ツールが LLM クライアントに表示されない

LLM クライアントがツールを認識できません。

```bash
# spec が読み込まれているか確認
swag2mcp ls

# spec が無効になっていないか確認
swag2mcp validate

# サーバーログを確認
swag2mcp mcp --logfile /tmp/swag2mcp.log
cat /tmp/swag2mcp.log

# IDE 設定のワークスペースパスが正しいか確認
# （絶対パスである必要があります）
```

**よくある原因：**
- IDE 設定のワークスペースパスが間違っている
- すべての spec に `disable: true` が設定されている
- `--tags` で spec がフィルタリングされている
- 指定されたパスに設定ファイルが存在しない

### MCP ハンドシェイクが失敗する（HTTP トランスポート）

SSE および Streamable HTTP トランスポートでは、ツール呼び出しが機能する前に MCP プロトコルの初期化が必要です。

```
Step 1: POST /mcp → {"method":"initialize", ...}
Step 2: POST /mcp → {"method":"notifications/initialized"}
Step 3: POST /mcp → {"method":"tools/list", ...}  ← これで動作
```

LLM クライアントがツールを呼び出す前にハンドシェイクを完了していることを確認してください。

### ヘルスチェックが 404 を返す

ヘルスエンドポイントのパスが MCP パスと異なる場合があります。

```bash
# デフォルトのヘルスエンドポイント
curl http://127.0.0.1:8080/health

# MCP パスを変更しても、ヘルスは /health のまま
# （--http-path の影響を受けません）
```

### Auth ツールが利用できない

`auth` MCP ツールが表示されません。

`auth` ツールは**デフォルトで無効**です（`--disable-llm-auth=true`）。これは本番環境のセキュリティのための意図的な設定です。

```bash
# auth ツールを有効化
swag2mcp mcp --disable-llm-auth=false
```

## 認証の問題

### 401 Unauthorized

認証情報がないか無効なため、API がリクエストを拒否しました。

```bash
# 認証が設定されているか確認
swag2mcp info

# 設定を検証
swag2mcp validate

# 環境変数が設定されているか確認
echo $MY_TOKEN

# トークンが期限切れでないか確認（bearer トークンは静的）
```

**よくある原因：**
- トークンがないか空
- 環境変数が設定されていない
- トークンが期限切れ（bearer トークンは自動更新されない）
- 間違った認証タイプが設定されている

### 403 Forbidden

権限不足のため API がリクエストを拒否しました。

- トークンに必要なスコープがない可能性があります
- API キーがこのリソースにアクセスできない可能性があります
- API ドキュメントで必要な権限を確認してください

### OAuth2 トークンエンドポイントに到達できない

swag2mcp が OAuth2 トークン URL に到達できません。

```bash
# 設定の token_url を確認
# URL が正しく、到達可能か確認
curl -X POST https://auth.example.com/oauth/token \
  -d "grant_type=client_credentials" \
  -d "client_id=test" \
  -d "client_secret=test"

# ネットワーク接続を確認
# 企業プロキシの背後にある場合はプロキシ設定を確認
```

### Digest 認証が失敗する

swag2mcp が Digest 認証ハンドシェイクを完了できません。

- サーバーは 401 レスポンスとともに `WWW-Authenticate: Digest ...` ヘッダーを返す必要があります
- チャレンジは 5 分間キャッシュされます — サーバーが nonce を変更した場合、キャッシュの期限切れを待ちます
- ユーザー名とパスワードが正しいか確認してください

### HMAC 署名の不一致

API が HMAC 署名付きリクエストを拒否しました。

- `api_key` と `secret_key` が正しいか確認
- API が Binance スタイルの HMAC-SHA256 署名を使用しているか確認
- 一部の取引所は異なる署名方法を使用します — HMAC 認証は特に Binance 互換 API 向けです

### Script 認証が失敗する

外部認証スクリプトが失敗しました。

```bash
# スクリプトが存在するか確認
ls -la ~/.swag2mcp/auth_scripts/my-domain.sh

# スクリプトを手動で実行してテスト
sh ~/.swag2mcp/auth_scripts/my-domain.sh

# スクリプトの出力形式を確認（JSON である必要があります：{"token": "...", "expires_in": 3600}）
# スクリプトが 30 秒以内に完了するか確認
# スクリプトに実行権限があるか確認
chmod +x ~/.swag2mcp/auth_scripts/my-domain.sh
```

## 検索の問題

### 検索結果がない

検索でエンドポイントが見つかりませんでした。

```bash
# spec が読み込まれているか確認
swag2mcp ls

# spec が無効になっていないか確認
swag2mcp validate

# よりシンプルなクエリを試す
# メソッドで検索：method:GET
# タグで検索：tag:pets

# インデックスは MCP サーバー起動時に再構築されます
# spec を追加したばかりの場合は、サーバーを再起動してください
```

### 検索結果が無関係

クエリが広すぎるか曖昧です。

- フィールドフィルターを使用して絞り込む：`method:GET +tag:pets`
- 正確なフレーズを使用：`"find pet by status"`
- `limit` パラメーターを使用してより焦点を絞った結果を得る

## API 呼び出しの問題

### invoke がエラーを返す

API 呼び出しが失敗しました。

```bash
# エラーメッセージを確認 — HTTP ステータスコードが含まれています
# 4xx エラー：パラメーター、認証、権限を確認
# 5xx エラー：API サーバーに問題があります

# 呼び出し前に必ずエンドポイントを調査
inspect(endpointId: "...")

# すべての必須パラメーターが指定されているか確認
# パラメーターの型（string、number、boolean）を確認
```

### レート制限エラー

LLM が同じエンドポイントを短時間に呼び出しすぎました。

各エンドポイントには 10 秒のクールダウンがあります。再呼び出し前に待機するか、レートリミッターを無効にします：

```yaml
disable_ratelimiter: true
```

### レスポンスが大きすぎる（fileRef が返された）

レスポンスが `max_response_size` を超えました。

これは正常です。レスポンスツールを使用してデータを探索します：

```
1. response_outline(path) → 構造を理解
2. response_compress(path, mode: "first_of_array") → サンプルを取得
3. response_slice(path, jsonPath: "data.0") → 特定のデータを取得
```

または制限を増やします：

```yaml
http_client:
  max_response_size: 4194304  # 4 MB
```

### API レスポンスが遅い

API の応答に時間がかかりすぎています。

```yaml
http_client:
  timeout: 120s  # デフォルトの 30s から増加
```

## ワークスペースの問題

### swag2mcp init が失敗する：「directory is not empty」

ターゲットディレクトリに既にファイルがあります。

```bash
# --force で上書き
swag2mcp init --force

# または別のディレクトリを使用
swag2mcp init ./new-workspace
```

### swag2mcp update が失敗する

1 つ以上の spec ファイルをダウンロードできませんでした。

```bash
# エラーメッセージでどの URL が失敗したか確認
# URL にアクセスできるか確認
curl -I <failed-url>

# ネットワーク接続を確認
# プロキシ設定を確認
```

### エクスポートで ZIP が作成されない

`[output]` 引数はディレクトリではなく `.zip` で終わるファイルパスである必要があります。

```bash
# 正しい
swag2mcp export /path/to/workspace /path/to/backup.zip

# 間違い（ZIP は作成されません）
swag2mcp export /path/to/workspace /some/directory
```

### インポートが失敗する：「not a valid swag2mcp backup」

ZIP ファイルが `swag2mcp export` で作成されたものではありません。

`swag2mcp export` で作成された ZIP アーカイブのみインポートできます。アーカイブは特定の内部構造（`swag2mcp.yaml`、`specs/`、`auth_scripts/`）を持っています。

## TUI の問題

### TUI が正しく表示されない

ターミナルが小さすぎるか、必要な機能をサポートしていません。

- 最小ターミナルサイズ：80×24 文字
- TUI は Bubbletea を使用し、最新のほとんどのターミナルで動作します
- ターミナルウィンドウのサイズを変更してみてください
- 別のターミナルエミュレーターを試してみてください

### TUI に「no specs found」と表示される

ワークスペースに設定された spec がありません。

```bash
# spec を確認
swag2mcp ls

# spec を追加
swag2mcp add spec
```

## モックサーバーの問題

### モックサーバーが起動しない

```bash
# 設定で mock_enabled: true を確認
# すべての collection に base_mock_url が設定されているか確認
# ポートが使用中でないか確認
lsof -i :9090

# モックサーバーログを確認
swag2mcp-mock mockserver
```

### モックサーバーが空のレスポンスを返す

spec ファイルにレスポンススキーマが定義されていない可能性があります。

- モックサーバーはレスポンススキーマからデータを生成します
- スキーマが見つからない場合、`{}` を返します
- OpenAPI spec に `responses` と `schema` が定義されているか確認してください

## ネットワークの問題

### プロキシ接続が失敗した

swag2mcp が設定されたプロキシに接続できません。

```bash
# プロキシ URL の形式を確認（スキームを含める必要があります：http://、https://、socks5://）
# プロキシ認証情報を確認
# バイパスリストを確認 — ターゲットがバイパスリストにある可能性があります
# curl でプロキシをテスト
curl -x http://proxy.company.com:8080 https://api.example.com
```

### TLS/SSL エラー

証明書の検証に失敗しました。

- MCP サーバーに自己署名証明書を使用する場合、クライアントがそれを信頼する必要があります
- `--tls` を使用したモックサーバーの場合、自己署名証明書が自動的に生成されます
- API 呼び出しの場合、swag2mcp はシステムの証明書ストアを使用します

## その他の問題

### ディスク使用量が多い

キャッシュとレスポンスディレクトリは時間とともに増大する可能性があります。

```bash
# すべてをクリーンアップ
swag2mcp clean

# 古いレスポンス（48 時間以上）は MCP サーバー起動時に自動的にクリーンアップされます
# キャッシュファイルは 1〜48 時間の間でランダムに期限切れになります
```

### go install 後に「command not found」

`go install` ディレクトリが PATH にありません。

```bash
# Go がバイナリをインストールする場所を確認
go env GOPATH
# PATH に追加
export PATH=$PATH:$(go env GOPATH)/bin
```

### LLM がツールを正しく使用しない

LLM にはより良い指示やフォーマットスキルが必要な場合があります。

- spec 設定で `llm_instruction` を使用して API の機能を説明する
- 一貫した出力フォーマットには [swag2mcp-format スキル](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md) の使用を検討する
- LLM レスポンスの品質はモデルと受け取る指示に依存します

### バグを報告するには？

[GitHub](https://github.com/mmadfox/swag2mcp/issues) で Issue を開き、以下を含めてください：
- swag2mcp バージョン（`swag2mcp --version`）
- お使いの OS とアーキテクチャ
- 実行した正確なコマンド
- 完全なエラーメッセージ
- 設定ファイル（シークレットは削除してください）
