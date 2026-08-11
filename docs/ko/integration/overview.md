# 모든 MCP 클라이언트

swag2mcp는 **MCP 서버**(Model Context Protocol)입니다. 즉, 이 섹션에 나열된 것뿐만 아니라 **모든 MCP 클라이언트**와 함께 작동합니다. 에디터, IDE 또는 에이전트가 MCP 프로토콜을 지원한다면 swag2mcp를 연결할 수 있습니다.

## 범용 패턴

모든 MCP 클라이언트는 동일한 기본 설정을 사용합니다. swag2mcp를 MCP 서버로 추가합니다:

- **명령:** `swag2mcp`
- **인수:** `mcp`(선택적으로 워크스페이스 경로: `mcp /path/to/workspace`)

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

설정 파일의 정확한 위치와 형식(JSON, TOML, GUI 설정)은 클라이언트마다 다릅니다. **클라이언트의 MCP 문서**에서 설정 위치를 확인하세요.

## 전송

- **stdio** — 어디서나 작동합니다. 대부분의 MCP 클라이언트가 지원
- **HTTP(SSE / Streamable HTTP)** — HTTP 전송 옵션이 있는 클라이언트에서 지원

전송 플래그는 [`mcp` 명령](/ko/cli/mcp) 참조를 확인하세요.

## 테스트된 통합

| 클라이언트 | 가이드 |
|------------|--------|
| OpenCode | [OpenCode](/ko/integration/opencode) |
| Cursor | [Cursor](/ko/integration/cursor) |
| Claude Desktop | [Claude Desktop](/ko/integration/claude) |
| VS Code | [VS Code](/ko/integration/vscode) |
| Crush | [Crush](/ko/integration/crush) |

> 목록에 없는 클라이언트라도 **지원되지 않는다는 의미는 아닙니다.** MCP 프로토콜을 사용하는 한 위의 범용 패턴을 사용하고 클라이언트 매뉴얼에 따라 설정 위치를 확인하세요.
