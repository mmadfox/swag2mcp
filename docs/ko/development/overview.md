# 개발 개요

## 이 프로젝트에 대해

swag2mcp는 OpenAPI/Swagger/Postman 명세를 Model Context Protocol(MCP)을 통해 LLM 에이전트와 연결하는 Go 프로젝트입니다. Go 1.23+로 빌드되었으며 80개 이상의 린터가 적용된 엄격한 코딩 규칙을 따릅니다.

이 섹션은 코드베이스를 이해하고, 기여하거나, 새 인증 방법, MCP 도구 또는 통합으로 swag2mcp를 확장하려는 **엔지니어**를 위해 작성되었습니다.

## 개발 스킬

프로젝트에는 프로젝트의 규칙과 패턴을 인코딩하는 두 가지 개발 스킬이 포함되어 있습니다. 사용하거나 무시할 수 있습니다 — 도구일 뿐 규칙이 아닙니다.

### godeveloper

[godeveloper 스킬](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/godeveloper/SKILL.md)은 프로젝트의 모든 코드 규칙을 정의합니다:

- **이름** — 패키지, 파일, 타입, 인터페이스, 리시버, 상수
- **포맷팅** — gofmt/gofumpt/goimports/gci, 120줄 제한, 임포트 순서
- **오류 처리** — 8개 오류 코드가 있는 `LLMError`, 센티넬 오류, 오류 래핑
- **인터페이스** — 작은 인터페이스, 구성, 소비자 측 정의
- **동시성** — 뮤텍스 세분성, 고루틴 수명, 컨텍스트 전달
- **테스트** — 테이블 기반 테스트, `newTestService()`/`seedTestData()` 헬퍼, 모의 생성
- **프로젝트 패턴** — 서비스 레이어, 요청/응답 구조체, 함수형 옵션, MCP 핸들러 패턴

### swag2mcp-cli

[swag2mcp-cli 스킬](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md)은 구문, 플래그, 인수, 예시와 함께 모든 CLI 명령어를 문서화합니다. CLI 명령어 작업이나 문서 작성 시 유용합니다.

## 주요 아키텍처 결정

### 서비스 레이어 패턴

모든 기능은 동일한 3단계 패턴을 따릅니다:

1. **요청 검증** `s.validateRequest(req)` 사용 (`go-playground/validator` 사용)
2. **엔티티 조회** 메모리 내 인덱스에서 (`not_found` 코드로 `LLMError` 반환)
3. **비즈니스 로직 실행** 및 타입화된 응답 또는 `LLMError` 반환

```go
func (s *Service) Search(ctx context.Context, req SearchRequest) (SearchResponse, error) {
    if err := s.validateRequest(req); err != nil {
        return SearchResponse{}, NewLLMError(validationFailedErrCode, err.Error())
    }
    results, err := s.index.Search(req.Query, req.Limit)
    if err != nil {
        return SearchResponse{}, NewLLMError(invokeErrorCode, err.Error())
    }
    return SearchResponse{Results: results}, nil
}
```

### 요청/응답 구조체

각 메서드에는 전용 `{Method}Request` 및 `{Method}Response` 구조체가 있습니다. 요청 구조체는 검증을 위해 `validate` 태그를, 문서화를 위해 `jsonschema` 태그를 사용합니다:

```go
type SearchRequest struct {
    Query string `json:"query" validate:"required,min=1" jsonschema:"description=필드 필터를 지원하는 검색 쿼리"`
    Limit int    `json:"limit" validate:"required,min=1,max=50" jsonschema:"description=최대 결과 수"`
}

type SearchResponse struct {
    Results []EndpointSearchItem `json:"results"`
}
```

### 함수형 옵션

설정은 함수형 옵션 패턴을 사용합니다:

```go
type Option func(*Service)

func New(opts ...Option) (*Service, error)

func WithDisableLLMAuth(disable bool) Option {
    return func(s *Service) {
        s.disableLLMAuth.Store(disable)
    }
}
```

### MCP 핸들러 패턴

MCP 서버는 구성된 인터페이스 패턴을 사용합니다. `internal/server/mcp/handler.go`의 `Svc` 인터페이스는 더 작은 인터페이스(`CatalogReader`, `EndpointExplorer`, `EndpointExecutor`, `SystemInfo`, `ResponseManager`)로 구성됩니다. 각 핸들러 메서드는 서비스 레이어에 위임합니다:

```go
type handler struct {
    service Svc
}

func (h *handler) handleSearch(ctx context.Context, _ *sdkmcp.CallToolRequest, req service.SearchRequest) (*sdkmcp.CallToolResult, any, error) {
    resp, err := h.service.Search(ctx, req)
    if err != nil {
        return nil, nil, err
    }
    return &sdkmcp.CallToolResult{StructuredContent: resp}, nil, nil
}
```

### LLMError

LLM에 반환되는 모든 오류는 8개 코드 중 하나와 함께 `LLMError` 타입을 사용합니다:

| 코드 | 발생 시기 |
|------|----------|
| `validation_failed` | 잘못된 입력 (잘못된 ID 형식, 필수 필드 누락) |
| `not_found` | 인덱스에서 엔티티를 찾을 수 없음 |
| `rate_limit` | 엔드포인트별 10초 쿨다운 초과 |
| `invoke_error` | HTTP 요청/응답 실패 |
| `config_error` | 설정 로드 또는 검증 실패 |
| `workspace_error` | 워크스페이스 디렉토리 또는 파일 작업 실패 |
| `parse_error` | 명세 파일 파싱 실패 |
| `auth_error` | 인증 토큰 검색 실패 |

메시지는 무엇이 잘못되었는지와 **다음에 무엇을 해야 하는지**를 LLM 소비자에게 적합한 평이한 언어로 설명해야 합니다.

### ID 생성

모든 ID는 결정론적 MD5 해시입니다:

```go
id.Domain("meteo")                          // 32자 16진수
id.Collection("meteo", "Forecast")          // 32자 16진수
id.Tag("meteo", "Forecast", "pets")         // 32자 16진수
id.Method("meteo", "Forecast", "pets", "GET", "/v2/pet/{petId}")
```

### 설정 계단식

설정은 **전역 → spec → collection**의 세 수준으로 계단식으로 적용됩니다. 각 수준이 이전 수준을 재정의합니다. 모든 `http_client` 설정은 모든 수준에서 재정의할 수 있습니다. 헤더와 쿠키는 병합되고, 단순 값은 대체됩니다.

## 빠른 참조

| 영역 | 규칙 |
|------|------|
| **Go 버전** | 1.23+ |
| **포맷터** | gofmt, gofumpt, goimports, gci |
| **줄 길이** | 120자 |
| **린터** | `.golangci.yml`에 80개 이상 |
| **오류 타입** | 8개 코드의 `LLMError` |
| **모의 프레임워크** | `go.uber.org/mock` |
| **테스트 헬퍼** | `newTestService()`, `seedTestData()` |
| **설정 형식** | 계단식 YAML |
| **인증 디스패치** | `UnmarshalYAML`이 `type` 필드 읽기 |
| **ID 생성** | MD5 기반 (`id.Domain()`, `id.Collection()` 등) |
| **속도 제한** | `invoke`용 엔드포인트당 10초 |
| **응답 크기** | 기본 1MB, 초과 시 파일에 저장 |
| **커버리지 목표** | 핵심 패키지 80%+ |
| **빌드** | `make build` |
| **린트** | `make lint` |
| **테스트** | `go test ./...` |
| **생성** | `go generate ./...` |
