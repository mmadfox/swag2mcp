# mcp

## 목적

**MCP(Model Context Protocol) 서버**를 시작합니다 — LLM 통합의 기본 모드입니다. AI 에이전트(Claude, Cursor, OpenCode 등)가 16개의 MCP 도구를 통해 API에 접근할 수 있도록 실행하는 명령어입니다.

## 사용 시기

- LLM 에이전트를 API에 연결하려고 할 때
- IDE(VS Code, Cursor, JetBrains) 또는 데스크톱 앱(Claude Desktop)을 설정할 때
- MCP 프로토콜을 통해 API를 노출해야 할 때
- 통합 전에 MCP 서버를 테스트할 때

## 구문

```bash
swag2mcp mcp [path] [flags]
```

## 인수

| 인수 | 위치 | 필수 | 설명 |
|------|------|------|------|
| `path` | 1 | 아니요 | 워크스페이스 디렉토리. 생략 시 경로 해결 규칙에 따라 결정됩니다. |

## 플래그

| 플래그 | 약어 | 타입 | 기본값 | 설명 |
|-------|------|------|--------|------|
| `--transport` | | `string` | `"stdio"` | MCP 전송: `stdio`, `sse`, `streamable-http` |
| `--http-addr` | | `string` | `":8080"` | HTTP 서버 주소 (`sse` 및 `streamable-http`용) |
| `--http-path` | | `string` | `"/mcp"` | MCP 핸들러의 HTTP 경로 |
| `--auth-token` | | `string` | `""` | HTTP 전송 인증용 Bearer 토큰 |
| `--logfile` | `-f` | `string` | `""` | 로그 파일 경로. 설정하지 않으면 stderr로 로깅 |
| `--disable-llm-auth` | | `bool` | `true` | MCP 도구 목록에서 `auth` 도구 제거 |
| `--dump-dir` | | `string` | `""` | 디버깅용 HTTP 요청 덤프 디렉토리 |
| `--tags` | `-t` | `string` | `""` | 태그로 spec 필터링 (쉼표로 구분) |

## 작동 방식

### stdio 전송 (기본값)

LLM 클라이언트(IDE, Claude Desktop 등)가 MCP 서버를 하위 프로세스로 실행할 때 사용됩니다. 서버는 표준 입출력으로 통신합니다.

```bash
swag2mcp mcp
```

### SSE 전송

HTTP 기반 통신을 위한 Server-Sent Events 전송입니다. MCP 핸드셰이크 시퀀스가 필요합니다.

```bash
swag2mcp mcp --transport sse --http-addr :8080
```

### Streamable HTTP 전송

스트리밍 응답을 지원하는 최신 HTTP 전송입니다.

```bash
swag2mcp mcp --transport streamable-http --http-addr 0.0.0.0:8080
```

### 인증 사용

Bearer 토큰으로 HTTP 엔드포인트 보호:

```bash
swag2mcp mcp --transport sse --http-addr :8080 --auth-token "my-secret"
```

### 태그 필터링 사용

특정 태그가 있는 spec만 로드:

```bash
swag2mcp mcp --tags=public
```

### auth 도구 활성화 (디버그 모드)

LLM이 `auth` 도구를 통해 새 토큰을 요청하도록 허용:

```bash
swag2mcp mcp --disable-llm-auth=false
```

### 요청 덤프 디렉토리 사용

디버깅을 위해 모든 HTTP 요청 저장:

```bash
swag2mcp mcp --dump-dir ./dumps
```

## MCP HTTP 전송 — 핸드셰이크 프로토콜

`sse` 또는 `streamable-http`를 사용할 때 MCP 프로토콜은 특정 핸드셰이크가 필요합니다. 초기화 전에는 도구 호출이 실패합니다:

```
1단계: POST /mcp → {"method":"initialize", ...}
2단계: POST /mcp → {"method":"notifications/initialized"}
3단계: POST /mcp → {"method":"tools/list", ...}   ← 이제 작동
```

### Health check

초기화 없이 작동:

```bash
curl http://localhost:8080/health
# → {"status":"ok","version":"v1.2.0"}
```

## IDE 설정 예시

### VS Code (`.vscode/settings.json` 또는 전역 설정)

```json
{
  "mcp": {
    "servers": {
      "swag2mcp": {
        "command": "swag2mcp",
        "args": ["mcp", "/absolute/path/to/.swag2mcp"]
      }
    }
  }
}
```

### Cursor / Windsurf (`~/.cursor/mcp.json` 또는 프로젝트 `.cursor/mcp.json`)

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/absolute/path/to/.swag2mcp"]
    }
  }
}
```

### Claude Desktop (macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`)

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/absolute/path/to/.swag2mcp"]
    }
  }
}
```

### JetBrains IDE (설정 → 도구 → MCP)

- 이름: `swag2mcp`
- 명령어: `swag2mcp`
- 인수: `mcp /absolute/path/to/.swag2mcp`

> **IDE 설정에서 워크스페이스 디렉토리에 항상 절대 경로를 사용하세요.** 상대 경로는 IDE의 작업 디렉토리에 따라 실패할 수 있습니다.

## 출력

성공 시 서버가 출력:

```
MCP server listening on http://127.0.0.1:8080/mcp
```

## 세부 사항

- **자동 초기화 없음:** 설정 파일이 없으면 `mcp`가 오류를 반환합니다: `"configuration not found at <path>"`. 먼저 `init`을 실행하세요.
- **`--disable-llm-auth` (기본값: `true`):** 활성화되면 `auth` 도구가 MCP 도구 목록에서 완전히 제거됩니다. LLM은 토큰을 보거나 요청할 수 없습니다. 인증은 여전히 작동합니다 — 토큰은 LLM을 통하지 않고 표준 설정 메커니즘을 통해 획득됩니다. 이 모드는 **프로덕션**에 권장됩니다. **디버깅** 또는 단기 토큰 사용 시 `--disable-llm-auth=false`로 설정하여 LLM이 `auth` 도구를 통해 새 토큰을 요청하도록 허용하세요.
- **YAML 설정 폴백:** CLI 플래그가 명시적으로 설정되지 않은 경우 `swag2mcp.yaml`의 `mcp` 섹션에서 값을 가져옵니다(있는 경우). 이를 통해 매번 플래그를 전달하지 않고 설정 파일에서 서버를 구성할 수 있습니다.
- **응답 정리:** 시작 시 48시간보다 오래된 응답이 `responses/` 디렉토리에서 자동으로 제거됩니다.
- **경로 해결 경고:** `[path]`가 생략되면 `mcp`는 현재 디렉토리에서 `swag2mcp.yaml`을 먼저 검색한 후 `~/.swag2mcp/`로 폴백됩니다. 잘못된 디렉토리에서 명령어를 실행하면 의도한 것과 다른 워크스페이스가 로드될 수 있습니다. **서비스로 실행하거나 IDE 설정에서 항상 `[path]`를 명시적으로 지정하세요.**
