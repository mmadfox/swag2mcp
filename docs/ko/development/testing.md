# 테스트

## 명령어

```bash
# 단위 테스트
go test ./...

# 특정 패키지
go test ./internal/service/...

# 통합 테스트
make integration-tests

# 커버리지
make cover

# 모든 테스트
make testall
```

## 테스트 구조

```
tests/
├── main_test.go              # 진입점
├── suite_test.go             # 스위트 설정
├── suite_auth_test.go        # 인증 테스트
├── suite_config_test.go      # 설정 테스트
├── suite_mcp_tools_test.go   # MCP 도구 테스트
├── suite_search_test.go      # 검색 테스트
├── suite_ratelimit_test.go   # 속도 제한 테스트
├── suite_response_test.go    # 응답 테스트
├── suite_export_test.go      # 내보내기 테스트
├── suite_import_test.go      # 가져오기 테스트
├── suite_parsing_test.go     # 파싱 테스트
├── suite_transport_test.go   # 전송 테스트
├── suite_mock_test.go        # 모의 서버 테스트
├── suite_workspace_test.go   # 워크스페이스 테스트
├── suite_errors_test.go      # 오류 테스트
└── suite_version_test.go     # 버전 테스트
```

## 커버리지

목표: 핵심 패키지 80%+:

- `auth`
- `cache`
- `config`
- `env`
- `httpclient`
- `id`
- `index`
- `server/mcp`
- `service`
- `spec`
- `workspace`

## 모의

MCP 서버 테스트에 `go.uber.org/mock` 사용:

```bash
go generate ./...
```

`handler.go`에서 `internal/server/mcp/mock_svc_test.go`를 생성합니다.

## 테이블 기반 테스트

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"유효한 입력", "hello", "HELLO", false},
        {"빈 입력", "", "", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := DoSomething(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.Equal(t, tt.want, got)
        })
    }
}
```
