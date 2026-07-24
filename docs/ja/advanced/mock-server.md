# モックサーバー

## 概要

モックサーバーは OpenAPI スキーマに基づいて偽の API レスポンスを生成します。実際の HTTP 呼び出しを行わずに API 統合をテストできます。開発、LLM エージェントのテスト、デモンストレーションに便利です。

モックサーバーは**別のバイナリ** — `swag2mcp-mock` です。メインの `swag2mcp` バイナリには含まれておらず、別途インストールする必要があります。

## インストール

```bash
# オプション 1: GitHub Releases からダウンロード
# swag2mcp-mock_<version>_<os>_<arch>.tar.gz を探す

# オプション 2: Go でインストール
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

## 設定

設定でモックサーバーを有効にします：

```yaml
mock_enabled: true
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092

specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
        base_mock_url: "127.0.0.1:9090"
```

## パラメーター

### mock_enabled

- **型:** `bool`
- **デフォルト:** `false`
- **効果:** `true` の場合、すべてのアクティブな collection に `base_mock_url` が設定されている必要があります。モックサーバーは各 collection の HTTP サーバーを起動します。

### mock_auth

モック認証サーバーのポート。これらは OAuth2、Digest、HMAC 認証エンドポイントをシミュレートし、実際の認証情報なしで認証付き API をテストできます。

| フィールド | デフォルト | 説明 |
|-----------|-----------|------|
| `oauth2_port` | `9090` | モック OAuth2 トークンサーバーのポート |
| `digest_port` | `9091` | モック Digest 認証サーバーのポート |
| `hmac_port` | `9092` | モック HMAC 認証サーバーのポート |

### base_mock_url（collection ごと）

- **型:** `string`
- **必須:** はい（`mock_enabled: true` の場合）
- **形式:** `host:port`（例：`localhost:8080`、`127.0.0.1:9000`）
- **効果:** 各 collection はこのアドレスで独自の HTTP サーバーを取得します。サーバーは spec で定義されたすべてのエンドポイントにランダムに生成されたデータで応答します。

## モックサーバーの起動

```bash
# デフォルト設定で起動
swag2mcp-mock mockserver

# TLS で起動
swag2mcp-mock mockserver --tls

# カスタム TLS 証明書で起動
swag2mcp-mock mockserver --tls --tls-cert cert.pem --tls-key key.pem
```

### TLS フラグ

| フラグ | 説明 |
|-------|------|
| `--tls` | 自己署名証明書で TLS を有効化 |
| `--tls-cert` | TLS 証明書ファイルへのパス |
| `--tls-key` | TLS キーファイルへのパス |

`--tls` が `--tls-cert` と `--tls-key` なしで設定された場合、`localhost` 用の自己署名証明書が自動的に生成されます。

## モックサーバーの動作

モックサーバーを起動すると、次の処理が行われます：

1. **すべての spec ファイルを解析** — 各 collection の OpenAPI/Swagger spec を読み取ります
2. **ハンドラーを登録** — spec で定義されたすべてのパスとメソッドの HTTP ハンドラーを作成します
3. **偽のデータを生成** — レスポンススキーマに一致するランダムに生成されたデータ（正しい型、形式、構造）で応答します
4. **認証サーバーを起動** — テスト用に OAuth2、Digest、HMAC 認証エンドポイントをシミュレートします

### モックのテスト

```bash
# 1 つのターミナルで：
swag2mcp-mock mockserver

# 別のターミナルで：
curl http://localhost:8080/pets
# → [{"id":1,"name":"Pet_name","status":"available"}]
```

## 偽のデータの生成方法

モックサーバーは OpenAPI スキーマに基づいて現実的な偽のデータを生成します：

- **文字列** — ランダムな単語、文章、または形式固有の値（email、URL、UUID、日付、電話番号など）
- **数値** — 指定された範囲内のランダムな整数と浮動小数点数
- **ブール値** — ランダムな true/false
- **配列** — 1〜3 個のランダムな項目
- **オブジェクト** — すべてのプロパティがランダムな値で埋められます
- **Enum** — enum リストからランダムな値
- **Null 許容フィールド** — 時々 `null` を返します（約 10% の確率）

## ユースケース

- **開発** — 実際の API アクセスなしで統合をテスト
- **LLM エージェントのテスト** — LLM がエンドポイントを発見、調査、呼び出しできることを確認
- **デモンストレーション** — 実際の API を設定せずに swag2mcp の動作を表示
- **負荷テスト** — 実際の API にアクセスせずに MCP サーバーを負荷テスト

## 重要な注意点

- **別のバイナリ** — `swag2mcp-mock` はメインの `swag2mcp` バイナリに含まれていません。別途インストールしてください。
- **各 collection に独自のポート** — collection ごとに `base_mock_url` を設定します
- **認証モックサーバーはグローバル** — OAuth2、Digest、HMAC サーバーは collection の数に関係なく設定されたポートで実行されます
- **Spec 解析の失敗は致命的ではありません** — collection の spec が解析できない場合、警告とともにスキップされます
- **自己署名 TLS** — `--tls` を証明書なしで使用する場合、localhost のみの自己署名証明書が生成されます
