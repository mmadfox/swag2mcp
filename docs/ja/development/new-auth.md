# 新しい認証方法の追加

## 手順

1. **認証クライアントを作成** `internal/auth/<name>.go` に
2. **`Authenticator` インターフェースを実装**
3. **型定数を追加** `internal/auth/auth.go` に
4. **YAML デコーダーを追加** `internal/config/auth.go` に
5. **デコーダーを登録** `authDecoders` マップに
6. **テストを書く**

## 1. 認証クライアント

`internal/auth/my_auth.go` を作成：

```go
package auth

import "net/http"

type MyAuthClient struct {
    Token string `yaml:"token" validate:"required"`
}

func (c *MyAuthClient) New() error {
    c.Token = resolveEnv(c.Token)
    return nil
}

func (c *MyAuthClient) Type() Type {
    return MyAuth
}

func (c *MyAuthClient) Apply(req *http.Request, out *Info) error {
    if c.Token == "" {
        return nil
    }
    setAuthHeader(req, out, "X-My-Auth", c.Token)
    return nil
}

func (c *MyAuthClient) Validate() error {
    return authValidator.Struct(c)
}
```

## 2. Authenticator インターフェース

すべての認証クライアントは以下を実装する必要があります：

```go
type Authenticator interface {
    New() error                    // 初期化、環境変数の解決
    Type() Type                    // 認証型識別子を返す
    Apply(req *http.Request, out *Info) error  // リクエストに認証を適用
    Validate() error               // 必須フィールドを検証
}
```

## 3. 型定数

`internal/auth/auth.go` に追加：

```go
const MyAuth Type = "my-auth"
```

## 4. YAML デコーダー

`internal/config/auth.go` にデコーダー関数を追加します。デコーダーは `*yaml.Node` を受け取り、認証クライアント構造体にデコードする必要があります：

```go
func decodeMyAuth(node *yaml.Node) (auth.Authenticator, error) {
    var client auth.MyAuthClient
    if err := decodeConfig(node, &client); err != nil {
        return nil, err
    }
    return &client, nil
}
```

`decodeConfig` ヘルパーは共通パターンを処理します：ノードが空でないことを確認し、YAML を構造体にデコードし、失敗時に説明的なエラーを返します。

## 5. デコーダーの登録

`internal/config/auth.go` の `authDecoders` マップにデコーダーを追加：

```go
var authDecoders = map[string]authDecoder{
    // ... 既存のデコーダー
    auth.MyAuth.String(): decodeMyAuth,
}
```

`Auth` の `UnmarshalYAML` メソッドは YAML から `type` フィールドを読み取り、アンダースコアをハイフンに正規化し、`authDecoders` でデコーダーを検索し、`config` ノードで呼び出します。これが swag2mcp が各スペックに対してどの認証クライアントをインスタンス化するかを認識する方法です。

## 6. テスト

`internal/auth/my_auth_test.go` を作成し、以下をカバーするテーブル駆動テストを記述：

- `New()` が環境変数を正しく解決する
- `Type()` が正しい型を返す
- `Apply()` が正しいヘッダー/クエリパラメータを設定する
- `Apply()` が空の値を適切に処理する
- `Validate()` が有効な設定で成功する
- `Validate()` が必須フィールド欠落で失敗する
