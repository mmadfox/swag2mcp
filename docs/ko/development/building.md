# 빌드

## 요구 사항

- Go 1.26+
- Make

## 명령어

```bash
# 빌드
make build

# 버전과 함께 빌드
make build VERSION=1.0.0

# 린트
make lint

# 테스트
go test ./...

# 모든 테스트
make testall
```

## GoReleaser

릴리스용:

```bash
goreleaser release --snapshot --clean
```

## 플랫폼

| 플랫폼 | 아키텍처 |
|--------|---------|
| Linux | amd64, arm64 |
| macOS | amd64, arm64 |
| Windows | amd64 |

## 린트

```bash
make lint
```

80개 이상의 린터와 함께 `golangci-lint`를 사용합니다. 설정은 `.golangci.yml`에 있습니다.
