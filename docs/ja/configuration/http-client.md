# HTTP クライアント

swag2mcp はすべての API 呼び出しに設定可能な HTTP クライアントを使用します。これらの設定はグローバルに定義され、spec および collection レベルで上書きできます。

## 設定

```yaml
http_client:
  timeout: 30s
  max_response_size: 1048576
  user_agent: "swag2mcp-global/1.0"
  follow_redirects: true
  max_redirects: 10
  random: false
  proxy:
    url: ""
    username: ""
    password: ""
    bypass: []
  headers:
    "Accept": "application/json"
  cookies:
    - name: "session"
      value: "abc123"
      domain: ".example.com"
      path: "/"
```

## Timeout

swag2mcp が API レスポンスを待機する時間を制御します。

- **型:** 期間（Go 形式：`30s`、`60s`、`2m`）
- **デフォルト:** `30s`
- **範囲:** 1 秒〜5 分
- **効果:** API がこの時間内に応答しない場合、リクエストはタイムアウトエラーで失敗します。
- **増やすタイミング:** 遅い API、大きなペイロード、信頼性の低いネットワーク。
- **減らすタイミング:** 内部 API、ヘルスチェック、高速失敗シナリオ。

```yaml
http_client:
  timeout: 60s
```

## Max Response Size

レスポンスが swag2mcp によってディスクに保存され、インラインで LLM に返されなくなるサイズの制限。

- **型:** `int`（バイト）
- **デフォルト:** `1048576`（1 MB）
- **範囲:** 256 〜 10,485,760 バイト（10 MB）
- **効果:** レスポンスがこの制限を超えると、`{workspace}/responses/` に JSON ファイルとして保存されます。LLM はファイル参照を受け取り、`response_outline`、`response_compress`、`response_slice` ツールで探索できます。
- **増やすタイミング:** 大規模なデータセットを返す API（レポート、ログ、分析）。
- **減らすタイミング:** LLM のコンテキストウィンドウが限られている場合、またはすべてのレスポンスにファイルベースのアクセスを希望する場合。

```yaml
http_client:
  max_response_size: 4194304  # 4 MB
```

## User-Agent

すべてのリクエストとともに送信される `User-Agent` ヘッダー。一部の API は特定のユーザーエージェントを要求したり、既知のボットユーザーエージェントをブロックしたりします。

- **型:** `string`
- **デフォルト:** `"swag2mcp-global/1.0"`
- **効果:** API サーバーにアプリケーションを識別します。
- **変更するタイミング:** API が特定のユーザーエージェントを要求する場合、または分析のためにアプリケーションを識別したい場合。

```yaml
http_client:
  user_agent: "MyApp/1.0"
```

## Follow Redirects

swag2mcp が HTTP リダイレクト（3xx ステータスコード）を自動的に追跡するかどうかを制御します。

- **型:** `bool`
- **デフォルト:** `true`
- **効果:** `true` の場合、swag2mcp は `max_redirects` 回までリダイレクトを追跡します。`false` の場合、リダイレクトレスポンスがそのまま返されます。
- **無効にするタイミング:** ループでリダイレクトする API、リダイレクトターゲットを手動で検査したいセキュリティ重視のエンドポイント。

```yaml
http_client:
  follow_redirects: false
```

## Max Redirects

swag2mcp が停止するまでに追跡するリダイレクト数を制限します。

- **型:** `int`
- **デフォルト:** `10`
- **範囲:** 0 〜 50
- **効果:** API がこの制限よりも多くリダイレクトした場合、リクエストは失敗します。
- **変更するタイミング:** 長いリダイレクトチェーンがある API、またはリダイレクトループでの高速失敗のために減らす。

```yaml
http_client:
  max_redirects: 5
```

## Randomizer

各リクエストにブラウザ風のランダムヘッダーを追加して、フィンガープリンティングとブロックを回避します。

- **型:** `bool`
- **デフォルト:** `false`
- **効果:** `true` の場合、swag2mcp は各リクエストにランダムなヘッダーを生成します：`User-Agent`（実際のブラウザ文字列のプールから）、`Accept`、`Accept-Language`、`Accept-Encoding`、`Cache-Control`。これにより `user_agent` 設定が上書きされます。
- **有効にするタイミング:** User-Agent またはヘッダーパターンに基づいてリクエストをブロックする API、スクレイピングシナリオ。

```yaml
http_client:
  random: true
```

## Proxy

プロキシサーバーは swag2mcp とターゲット API の間の仲介役として機能します。すべての HTTP トラフィックはそれを経由してルーティングされます。

**プロキシが必要な場合：**
- **企業ネットワーク** — すべてのアウトバウンドトラフィックが会社のプロキシを通過する必要がある
- **地理的制限** — 一部の API は地域ロックされており、適切な地域のプロキシがこれを回避する
- **静的 IP** — IP 許可リストが必要な API
- **匿名性** — ターゲット API から発信元 IP を隠す

### Proxy URL

- **型:** `string`
- **デフォルト:** `""`（プロキシなし）
- **サポートされるスキーム:** `http`、`https`、`socks5`、`socks5h`
- **`$(VAR)` をサポート:** ✅ 実行時に解決

| スキーム | 説明 | ユースケース |
|---------|------|-----------|
| `http` | HTTP トラフィック用 HTTP プロキシ | 企業プロキシ、基本的なプロキシ |
| `https` | HTTPS プロキシ（CONNECT トンネル） | セキュアな企業プロキシ |
| `socks5` | SOCKS5 プロキシ（DNS はローカルで解決） | 汎用、任意のプロトコル |
| `socks5h` | SOCKS5 プロキシ（DNS はプロキシ上で解決） | プロキシの DNS 解決が優れている場合 |

### Proxy 認証

プロキシが認証を必要とする場合、`username` と `password` を指定します：

- **`$(VAR)` をサポート:** ✅ 3 つのフィールドすべて（`url`、`username`、`password`）で実行時に解決

```yaml
http_client:
  proxy:
    url: "http://proxy.example.com:8080"
    username: "proxyuser"
    password: "$(PROXY_PASSWORD)"
```

### Proxy バイパス

プロキシを**経由すべきでない**ドメインのリスト。内部サービス、localhost、または直接のみアクセス可能な API に便利です。

```yaml
http_client:
  proxy:
    url: "http://proxy.example.com:8080"
    bypass:
      - "localhost"
      - "127.0.0.1"
      - "*.internal.company.com"
      - "api.local"
```

バイパスはワイルドカードパターンをサポートします（`*.example.com` は任意のサブドメインに一致）。

## Headers

すべてのリクエストに追加されるカスタム HTTP ヘッダー。ヘッダーはカスケードレベル間でマージされます：

```
Global headers → Spec headers（マージ） → Collection headers（マージ）
```

Collection ヘッダーは spec ヘッダーを上書きし、spec ヘッダーは同じキーのグローバルヘッダーを上書きします。

```yaml
http_client:
  headers:
    "Accept": "application/json"
    "Accept-Language": "en-US"
```

ヘッダー値は `$(ENV_VAR)` 解決をサポートします。

## Cookies

すべてのリクエストとともに送信される Cookie。Cookie はカスケードレベル間でマージされます（低いレベルが同じ Cookie 名のグローバルを上書き）。

```yaml
http_client:
  cookies:
    - name: "session"
      value: "abc123"
      domain: ".example.com"
      path: "/"
      secure: false
      http_only: false
```

### Cookie フィールド

| フィールド | 必須 | 説明 |
|-----------|------|------|
| `name` | はい | Cookie 名 |
| `value` | はい | Cookie 値（`$(ENV_VAR)` 解決をサポート） |
| `domain` | いいえ | ドメインスコープ（例：`.example.com`） |
| `path` | いいえ | パススコープ（例：`/`） |
| `secure` | いいえ | HTTPS 経由でのみ送信 |
| `http_only` | いいえ | JavaScript からアクセス不可 |

## Spec レベルでのカスタムヘッダー

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    http_client:
      headers:
        "Accept": "application/json"
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Spec レベルでの Cookie

```yaml
specs:
  - domain: example
    llm_title: Example API
    base_url: https://api.example.com
    http_client:
      cookies:
        - name: "session"
          value: "abc123"
        - name: "csrf"
          value: "$(CSRF_TOKEN)"
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## カスケード

HTTP クライアント設定はグローバルから spec、collection へとカスケードされます。すべての設定はすべてのレベルで上書き可能：

```
Global (http_client)
    ↓ 上書き（すべての設定）
Spec (specs[].http_client)
    ↓ 上書き（すべての設定）
Collection (specs[].collections[].http_client)
```

**すべての HTTP クライアント設定**（タイムアウト、プロキシ、ユーザーエージェント、リダイレクト、レスポンスサイズ、ランダマイザー、ヘッダー、Cookie）は spec と collection の両方のレベルで上書きできます。

詳細は [Configuration Cascade](./cascade) を参照してください。
