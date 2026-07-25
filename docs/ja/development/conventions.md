# コード規約

## Go

- **Go 1.26+**
- **gofmt** / **gofumpt** / **goimports** / **gci**
- 1 行 **120 文字**
- ネストされた if の代わりに**ガード節**
- **命名**：プライベートは `camelCase`、エクスポートは `PascalCase`

## エラー

LLM から見えるエラーには `LLMError` を使用：

```go
type LLMError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}
```

エラーコード：
- `validation_failed` — 無効なパラメータ
- `not_found` — リソースが見つからない
- `rate_limit` — レート制限超過
- `invoke_error` — API 呼び出しエラー

## インターフェース

- 小さなインターフェース（1〜3 メソッド）
- インターフェース合成
- 設定には関数型オプション

## テスト

- テーブル駆動テスト
- テストヘルパー（`newTestService()`、`seedTestData()`）
- `go.uber.org/mock` によるモック
- コアパッケージで 80% 以上のカバレッジ

## 設定

- YAML 形式
- カスケード：グローバル → スペック → コレクション
- `go-playground/validator` による検証
- `$(VAR)` による環境変数
