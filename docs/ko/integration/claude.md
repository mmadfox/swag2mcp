# Claude Desktop 통합

## stdio

`claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp"]
    }
  }
}
```

## 커스텀 워크스페이스

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/path/to/workspace"]
    }
  }
}
```

## 사용법

Claude Desktop을 다시 시작한 후 다음을 할 수 있습니다:

- "모든 API 목록을 보여주세요"
- "주문 생성 엔드포인트를 찾아주세요"
- "모스크바의 날씨 API를 호출해주세요"

## 기타

클라이언트가 보이지 않나요? 모든 MCP 통합은 동일한 패턴을 따릅니다:
- 명령어를 `swag2mcp`로, 인수를 `mcp`로 설정
- 선택적으로 워크스페이스 경로 추가: `mcp /path/to/workspace`
- 정확한 설정 파일 위치와 형식은 클라이언트 문서를 확인

대부분의 MCP 클라이언트는 stdio 전송을 지원하며, 일부는 HTTP(SSE / Streamable HTTP)를 지원합니다.
