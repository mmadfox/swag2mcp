# VS Code 통합

## VS Code 설정을 통해

`.vscode/settings.json`:

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

## 확장을 통해

VS Code용 MCP 확장을 설치하고 추가:

```json
{
  "mcp.servers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp"]
    }
  }
}
```

## 사용법

설정 후 VS Code AI Assistant가 swag2mcp를 통해 API와 작업할 수 있습니다.

## 기타

클라이언트가 보이지 않나요? 모든 MCP 통합은 동일한 패턴을 따릅니다:
- 명령어를 `swag2mcp`로, 인수를 `mcp`로 설정
- 선택적으로 워크스페이스 경로 추가: `mcp /path/to/workspace`
- 정확한 설정 파일 위치와 형식은 클라이언트 문서를 확인

대부분의 MCP 클라이언트는 stdio 전송을 지원하며, 일부는 HTTP(SSE / Streamable HTTP)를 지원합니다.
