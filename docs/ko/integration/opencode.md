# OpenCode 통합

## stdio

```json
{
  "mcp": {
    "swag2mcp": {
      "type": "local",
      "command": ["swag2mcp", "mcp"],
      "enabled": true
    }
  }
}
```

## HTTP

```json
{
  "mcp": {
    "swag2mcp": {
      "type": "local",
      "command": ["swag2mcp", "mcp", "--transport", "sse", "--http-addr", "127.0.0.1:8080"],
      "enabled": true
    }
  }
}
```

## 예시 쿼리

연결 후 다음을 물어볼 수 있습니다:

- "어떤 API가 있나요?"
- "petstore의 모든 엔드포인트 보여주세요"
- "사용자 생성 API 찾아주세요"
- "GET /pet/1을 호출하고 결과를 보여주세요"

## 기타

클라이언트가 보이지 않나요? 모든 MCP 통합은 동일한 패턴을 따릅니다:
- 명령어를 `swag2mcp`로, 인수를 `mcp`로 설정
- 선택적으로 워크스페이스 경로 추가: `mcp /path/to/workspace`
- 정확한 설정 파일 위치와 형식은 클라이언트 문서를 확인

대부분의 MCP 클라이언트는 stdio 전송을 지원하며, 일부는 HTTP(SSE / Streamable HTTP)를 지원합니다.
