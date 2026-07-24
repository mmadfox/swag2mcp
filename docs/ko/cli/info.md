# info

## 목적

swag2mcp 런타임의 포괄적인 요약을 **JSON**으로 표시합니다. 버전, 워크스페이스 경로, spec 요약, HTTP 클라이언트 설정, MCP 전송 설정, 인증 방법, 모의 모드 상태가 포함됩니다.

## 사용 시기

- 워크스페이스의 머신 판독 가능 개요를 원할 때
- 디버깅을 위해 런타임 설정을 확인해야 할 때
- 활성화된 spec과 엔드포인트 수를 확인하려고 할 때
- HTTP 클라이언트 또는 MCP 전송 설정을 확인해야 할 때

## 구문

```bash
swag2mcp info [path]
```

## 인수

| 인수 | 위치 | 필수 | 설명 |
|------|------|------|------|
| `path` | 1 | 아니요 | 워크스페이스 디렉토리. 생략 시 경로 해결 규칙에 따라 결정됩니다. |

## 플래그

없음.

## 작동 방식

```bash
swag2mcp info
swag2mcp info ./my-workspace
```

## 출력

출력은 다음 구조의 JSON 객체입니다:

```json
{
  "version": "v1.2.0",
  "workspace": "~/.swag2mcp",
  "uptime": "2h 15m",
  "specs": {
    "total": 4,
    "active": 3,
    "disabled": 1,
    "collections": 6,
    "endpoints": 42
  },
  "http_client": {
    "timeout": "30s",
    "max_response_size": "2 KB",
    "proxy": "none",
    "follow_redirects": true,
    "max_redirects": 10,
    "randomize": false
  },
  "mcp": {
    "transport": "stdio",
    "addr": ":8080",
    "path": "/mcp"
  },
  "auth_methods": ["bearer", "api-key"],
  "mock_enabled": false
}
```

## 명령 후 검증

MCP 서버를 시작하기 전에 `info`를 사용하여 워크스페이스가 올바르게 로드되었고 모든 spec이 활성화되었는지 확인하세요.

## 세부 사항

- **자동 초기화:** 설정 파일이 없으면 `info`가 자동으로 init 마법사를 먼저 실행합니다.
- **JSON 전용:** 출력은 항상 JSON입니다. 사람이 읽을 수 있는 출력은 `ls`를 사용하세요.
- **`max_response_size`:** 사람이 읽을 수 있는 형식으로 표시됩니다 (예: `"1 KB"`, `"2 MB"`).
- **전문 검색 인덱스 없음:** `info`는 설정과 spec 메타데이터만 필요하므로 전문 검색 인덱싱을 비활성화합니다.
