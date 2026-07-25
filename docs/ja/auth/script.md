# Script 認証

## 目的

外部スクリプトによる認証 — 最も柔軟な方法。任意の言語（bash、Python など）でスクリプトを記述し、好きな方法でトークンを取得して swag2mcp に返すことができます。

## 使用するタイミング

- カスタムまたは非標準の認証スキーム
- 複雑なトークン取得ロジック（マルチステップ、追加チェック付き）
- 標準の方法がニーズに合わない場合

## 設定

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: script
      config:
        domain: "my-auth"
```

## パラメーター

| パラメーター | 必須 | 説明 |
|-----------|------|------|
| `domain` | はい | スクリプトファイル名（拡張子なし） |

## スクリプトの場所

スクリプトはワークスペースの `auth_scripts` ディレクトリに配置する必要があります：

- **Linux / macOS:** `{workspace}/auth_scripts/{domain}.sh`
- **Windows:** `{workspace}/auth_scripts/{domain}.bat`

## スクリプトの出力形式

スクリプトはトークンとその有効期限を JSON で stdout に出力する必要があります：

```bash
#!/bin/bash
# auth_scripts/my-auth.sh

TOKEN=$(curl -s -X POST https://auth.example.com/token \
  -d "grant_type=client_credentials" \
  -d "client_id=$CLIENT_ID" \
  -d "client_secret=$CLIENT_SECRET" | jq -r '.access_token')

echo "{\"token\": \"$TOKEN\", \"expires_in\": 3600}"
```

### JSON フィールド

| フィールド | 必須 | 説明 |
|-----------|------|------|
| `token` | はい | 認証トークン |
| `expires_in` | いいえ | トークンの有効期間（秒）（デフォルト：3600） |

## 注意点

- swag2mcp はキャッシュされたトークンが期限切れの場合、リクエストごとにスクリプトを実行します
- スクリプトは 30 秒以内に完了する必要があります
- トークンは有効期限までキャッシュされます
- スクリプトファイル名 = `{domain}.sh`（Unix）または `{domain}.bat`（Windows）
- `domain` に `/` または `\` を含めることはできません
