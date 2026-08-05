# VS Code 통합

## Via .vscode/mcp.json

1. VS Code용 MCP 확장을 설치하세요 (예: org.mcp의 MCP Client 등).
2. 프로젝트 루트에 `.vscode/mcp.json`을 생성:

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "${workspaceFolder}"]
    }
  }
}
```

> // "${{workspaceFolder}}"는 워크스페이스 경로로 전달됩니다

3. VS Code 창을 다시 로드하세요 (Ctrl+Shift+P → "Reload Window").
4. AI 어시스턴트를 사용하세요 — 이제 API를 인식합니다.

## 대안: VS Code 설정을 통해

`.vscode/settings.json`에서도 설정할 수 있습니다:

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

## 사용법

설정 후 VS Code AI 어시스턴트가 swag2mcp를 통해 API를 사용할 수 있습니다.

## 기타

클라이언트가 보이지 않나요? 모든 MCP 통합은 동일한 패턴을 따릅니다:
- 명령어를 `swag2mcp`로, 인수를 `mcp`로 설정
- 선택적으로 워크스페이스 경로 추가: `mcp /path/to/workspace`
- 정확한 설정 파일 위치와 형식은 클라이언트 문서를 확인

대부분의 MCP 클라이언트는 stdio 전송을 지원하며, 일부는 HTTP(SSE / Streamable HTTP)를 지원합니다.
