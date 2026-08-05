# Cursor 통합

## stdio

### Via Settings UI

1. Cursor 설정 열기 (Cmd+, / Ctrl+,)
2. Go to **MCP 서버**
3. Click **새 서버 추가**
4. Fill in:
   - **이름:** `swag2mcp`
   - **유형:** `command`
   - **명령어:** `swag2mcp mcp`
5. Click **저장**

### 설정 파일을 통해

In `~/.cursor/mcp.json`:

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

After connecting, Cursor AI Agent can:

- API 탐색
- 관련 엔드포인트 찾기
- API 호출 및 결과 표시
- 요청 디버깅 도움

## 기타

클라이언트가 보이지 않나요? 모든 MCP 통합은 동일한 패턴을 따릅니다:
- 명령어를 `swag2mcp`로, 인수를 `mcp`로 설정
- 선택적으로 워크스페이스 경로 추가: `mcp /path/to/workspace`
- 정확한 설정 파일 위치와 형식은 클라이언트 문서를 확인

대부분의 MCP 클라이언트는 stdio 전송을 지원하며, 일부는 HTTP(SSE / Streamable HTTP)를 지원합니다.
