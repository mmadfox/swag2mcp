# 認証

## 概要

swag2mcp は認証が必要な API を扱うための **9 つの認証方式** をサポートしています。設定ファイルに一度設定するだけで、以降 `invoke` によるすべての API 呼び出しに自動的に適切なトークンとヘッダーが含まれます。

### 設定場所

認証は `swag2mcp.yaml` の **spec** レベルで設定します：

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: bearer
      config:
        token: "my-token"
```

### 仕組み

- 設定で認証タイプとパラメーターを指定します
- swag2mcp は `invoke` を呼び出すときに自動的にすべてのリクエストに適用します
- API を呼び出す前にトークンを要求する**必要はありません** — 自動的に行われます
- トークンが期限切れになった場合（OAuth2、Script）、swag2mcp が自動的に更新します

### 環境変数

機密データ（トークン、パスワード、キー）は `$(VAR_NAME)` 構文を使用して環境変数に保存できます：

```yaml
auth:
  type: bearer
  config:
    token: "$(MY_API_TOKEN)"
```

swag2mcp は起動時に `MY_API_TOKEN` の値を代入します。

### MCP auth ツール

LLM エージェントは `auth` MCP ツールを介してトークンやヘッダーを取得できます — 例えば、curl コマンドを構築したりユーザーに表示したりするためです。

**本番環境**では、このツールは `--disable-llm-auth`（デフォルトで有効）で無効にし、LLM がトークンにアクセスできないようにする必要があります。

### 方式一覧

| 方式 | 説明 | 最適な用途 |
|------|------|-----------|
| [`none`](/auth/none) | 認証なし | 公開 API |
| [`basic`](/auth/basic) | HTTP Basic（ユーザー名 + パスワード） | レガシー API、シンプルな認証 |
| [`bearer`](/auth/bearer) | Bearer トークン（JWT、トークン） | 最新の REST API |
| [`api-key`](/auth/api-key) | ヘッダーまたはクエリパラメーターの API キー | API キーを使用するサービス |
| [`digest`](/auth/digest) | HTTP Digest（ユーザー名 + パスワード） | レガシー API、Basic より安全 |
| [`hmac`](/auth/hmac) | HMAC-SHA256 署名（Binance スタイル） | 暗号通貨取引所 |
| [`oauth2-cc`](/auth/oauth2-cc) | OAuth2 Client Credentials | サーバー間、マイクロサービス |
| [`oauth2-pwd`](/auth/oauth2-pwd) | OAuth2 Password Grant | ユーザーログイン付きアプリ |
| [`script`](/auth/script) | トークン取得用の外部スクリプト | 任意のカスタム認証スキーム |
