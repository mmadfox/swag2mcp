# ビルド

## 要件

- Go 1.26+
- Make

## コマンド

```bash
# ビルド
make build

# バージョン指定でビルド
make build VERSION=1.0.0

# リンター
make lint

# テスト
go test ./...

# 全テスト
make testall
```

## GoReleaser

リリース用：

```bash
goreleaser release --snapshot --clean
```

## プラットフォーム

| プラットフォーム | アーキテクチャ |
|----------|-------------|
| Linux | amd64、arm64 |
| macOS | amd64、arm64 |
| Windows | amd64 |

## リンター

```bash
make lint
```

`golangci-lint` を 80 以上のリンターとともに使用します。設定は `.golangci.yml` にあります。
