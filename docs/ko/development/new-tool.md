# 새 MCP 도구 추가

## 단계

1. **도구 이름 상수 추가** `internal/service/service.go`
2. **요청/응답 타입 생성** `internal/service/types.go`
3. **서비스 구현** `internal/service/` (새 파일 또는 기존 파일에 추가)
4. **마크다운 정의 생성** `internal/service/definitions/` — `MakeToolDefinitions`가 읽는 파일
5. **`Svc` 인터페이스에 메서드 추가** `internal/server/mcp/handler.go`
6. **핸들러 추가** `handler.go`
7. **도구 등록** `mcp.go`의 `registerTools`
8. **모의 생성**: `go generate ./...`
9. **테스트 작성**

## 1. 도구 이름 상수

`internal/service/service.go`에 상수 추가:

```go
const MyNewTool = "my_new_tool"
```

## 2. 요청/응답 타입

`internal/service/types.go`에 정의:

```go
type MyNewToolRequest struct {
    Param1 string `json:"param1" validate:"required" jsonschema:"required,param1 설명"`
}

type MyNewToolResponse struct {
    Result string `json:"result"`
}
```

## 3. 서비스 구현

`internal/service/my_new_tool.go`를 생성하거나 기존 서비스 파일에 추가합니다. 표준 서비스 패턴을 따르세요: 검증 → 조회 → 실행 → 반환:

```go
func (s *Service) MyNewTool(ctx context.Context, req MyNewToolRequest) (MyNewToolResponse, error) {
    if err := s.validateRequest(req); err != nil {
        return MyNewToolResponse{}, NewLLMError(validationFailedErrCode, err.Error())
    }
    // 비즈니스 로직
    return MyNewToolResponse{Result: "ok"}, nil
}
```

## 4. 마크다운 정의

`internal/service/definitions/my_new_tool.md`를 생성합니다. 이 파일은 `MakeToolDefinitions()`에서 읽어 바이너리에 포함됩니다. 프론트매터의 `name:` 필드는 상수와 일치해야 합니다:

```markdown
---
name: my_new_tool
---

# my_new_tool

도구 설명.

## 매개변수

| 매개변수 | 타입 | 설명 |
|---------|------|------|
| `param1` | string | 설명 |
```

`MakeToolDefinitions()` 함수는 `tools.go`에서 포함된 `definitions/` 디렉토리의 모든 `.md` 파일을 읽고, YAML 프론트매터에서 `name` 필드를 파싱하며, 본문을 도구 설명으로 사용합니다. `instruction.md` 파일은 특별히 처리되어 LLM의 시스템 지침이 됩니다.

## 5. Svc 인터페이스

`handler.go`의 구성된 `Svc` 인터페이스에 메서드를 추가하세요:

```go
type Svc interface {
    // ... 기존 메서드
    MyNewTool(ctx context.Context, req service.MyNewToolRequest) (service.MyNewToolResponse, error)
}
```

## 6. 핸들러

`handler.go`의 `handler`에 핸들러 메서드를 추가하세요. 핸들러는 서비스에 위임하고 결과를 `StructuredContent`로 래핑합니다:

```go
func (h *handler) handleMyNewTool(
    ctx context.Context,
    _ *sdkmcp.CallToolRequest,
    req service.MyNewToolRequest,
) (*sdkmcp.CallToolResult, any, error) {
    resp, err := h.service.MyNewTool(ctx, req)
    if err != nil {
        return nil, nil, err
    }
    return &sdkmcp.CallToolResult{
        StructuredContent: resp,
    }, nil, nil
}
```

## 7. 등록

`mcp.go`의 `registerTools` 함수에서 도구를 등록합니다. `toolRegistrations` 맵에 항목을 추가하세요:

```go
service.MyNewTool: {
    addTool[service.MyNewToolRequest](mcpServer, h.handleMyNewTool),
    true, // 도구가 변경 가능한 경우 false (invoke 또는 auth와 같은)
},
```

`registerTools` 함수 시그니처는:

```go
func registerTools(mcpServer *sdkmcp.Server, tools []service.Tool, h handler) {
```

`MakeToolDefinitions()`가 반환한 도구 정의를 반복하고 각각을 타입화된 핸들러와 함께 등록합니다. `toolRegistrations` 맵은 도구 이름 상수를 해당 핸들러에 연결합니다.
