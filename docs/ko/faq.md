# FAQ

## 일반

### swag2mcp란 무엇이며 어떤 문제를 해결하나요?

swag2mcp는 OpenAPI/Swagger/Postman API 명세를 Model Context Protocol(MCP)을 통해 LLM 에이전트와 연결합니다. 각 API를 AI 에이전트에 연결하기 위해 커스텀 코드를 작성하는 대신, YAML 파일에 한 번 구성하면 LLM이 API를 발견, 검사, 호출할 수 있는 19개의 도구를 사용할 수 있습니다.

### 다른 API-to-LLM 도구와 어떻게 다른가요?

- **코딩 불필요** — YAML로 API를 구성하며 통합 코드가 필요 없습니다
- **19개의 MCP 도구** — 발견부터 호출, 대용량 응답 처리까지 완벽한 도구 모음
- **9가지 인증 방식** — 모든 API 인증 체계와 호환
- **전문 검색** — bluge 기반의 모든 엔드포인트 검색
- **TUI 탐색기** — 대화형 터미널 인터페이스로 탐색 및 테스트
- **모의 서버** — 실제 API 호출 없이 테스트

### 어떤 API 명세 형식을 지원하나요?

OpenAPI 3.x, Swagger 2.0, Postman Collections v2.1을 지원합니다.

### spec과 collection의 차이는 무엇인가요?

**Spec**은 논리적 API 서비스(예: "Open-Meteo Weather APIs")를 나타냅니다. **Collection**은 하나의 OpenAPI/Swagger/Postman 파일입니다. 하나의 spec은 여러 collection을 가질 수 있습니다 — 예를 들어, API가 예보, 대기질, 해양 등 서로 다른 서비스에 대해 별도의 명세 파일을 가질 때 사용합니다.

### 어떤 MCP 전송 방식을 지원하나요?

세 가지 전송 방식: `stdio`(기본값, 로컬 LLM 클라이언트용), `sse`(원격 클라이언트용 Server-Sent Events), `streamable-http`(최신 HTTP 스트리밍).

### swag2mcp를 모든 LLM과 함께 사용할 수 있나요?

네, MCP 프로토콜을 지원하는 모든 LLM 클라이언트(Claude Desktop, VS Code, Cursor, Windsurf, JetBrains IDE, OpenCode 등)와 함께 사용할 수 있습니다.

## 설치

### swag2mcp를 어떻게 설치하나요?

```bash
# 옵션 1: GitHub Releases에서 다운로드
# https://github.com/mmadfox/swag2mcp/releases/latest 로 이동
# OS와 아키텍처에 맞는 아카이브 다운로드

# 옵션 2: Go로 설치
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

### Go가 설치되어 있어야 하나요?

아니요. Linux(amd64, arm64), macOS(amd64, arm64), Windows(amd64)용 사전 빌드된 바이너리가 [GitHub Releases 페이지](https://github.com/mmadfox/swag2mcp/releases)에서 제공됩니다.

### 모의 서버는 어떻게 설치하나요?

모의 서버는 별도의 바이너리입니다:

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

또는 GitHub Releases에서 `swag2mcp-mock_&lt;version&gt;_&lt;os&gt;_&lt;arch&gt;.tar.gz`를 다운로드하세요.

## 시작하기

### 빠르게 시작하려면 어떻게 해야 하나요?

```bash
# 1. 워크스페이스 초기화
mkdir -p .swag2mcp && swag2mcp init ./.swag2mcp

# 2. MCP 서버 시작 (init 후 공개 예제 명세가 포함됩니다)
swag2mcp mcp
```

`init` 후 워크스페이스에는 이미 여러 공개 예제 명세(icanhazdadjoke, Open-Meteo, Binance, PokéAPI)가 포함됩니다. 바로 MCP 서버를 시작할 수 있습니다 — 수동으로 명세를 추가할 필요가 없습니다.

자체 API를 추가하려면:

```bash
swag2mcp add spec --yaml - <<EOF
domain: dadjoke
llm_title: icanhazdadjoke API
base_url: https://icanhazdadjoke.com
collections:
  - llm_title: Jokes
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
EOF
```

### IDE에 swag2mcp를 어떻게 연결하나요?

**VS Code** (`.vscode/settings.json`):
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

**Cursor** (`~/.cursor/mcp.json`):
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

**Claude Desktop** (`claude_desktop_config.json`):
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

항상 워크스페이스 디렉토리의 절대 경로를 사용하세요.

## 설정

### 설정 파일은 어디에 있나요?

기본 위치: `~/.swag2mcp/swag2mcp.yaml`. 다른 디렉토리에 생성하여 명령어에 경로를 전달할 수도 있습니다.

### API를 어떻게 추가하나요?

```bash
# 대화형 모드
swag2mcp add spec

# YAML 사용 (스크립팅에 권장)
swag2mcp add spec --yaml - <<EOF
domain: my-api
llm_title: My API
base_url: https://api.example.com/v1
collections:
  - llm_title: Main
    location: https://example.com/spec.yaml
EOF
```

### 기존 spec에 collection을 어떻게 추가하나요?

```bash
swag2mcp add collection --yaml - <<EOF
spec_domain: meteo
llm_title: Air Quality
location: https://example.com/air-quality.yaml
EOF
```

### spec을 임시로 비활성화하려면 어떻게 하나요?

spec 설정에서 `disable: true`로 설정하세요. 해당 spec은 로드되거나 인덱싱되지 않습니다.

### 로드할 spec을 필터링할 수 있나요?

네, `--tags` 플래그를 사용하세요: `swag2mcp mcp --tags=public`. 일치하는 태그가 있는 spec만 로드됩니다.

### 시크릿에 환경 변수를 어떻게 사용하나요?

인증 필드에 `$(VAR_NAME)` 구문을 사용하세요:

```yaml
auth:
  type: bearer
  config:
    token: "$(MY_API_TOKEN)"
```

시작 전에 변수를 설정하세요: `export MY_API_TOKEN="eyJhbGci..."`

## 인증

### 어떤 인증 방식을 지원하나요?

9가지 방식: `none`, `basic`, `bearer`, `digest`, `hmac`, `oauth2-cc`(클라이언트 자격 증명), `oauth2-pwd`(비밀번호 그랜트), `api-key`, `script`.

### 토큰을 어떻게 전달하나요?

설정 파일 또는 환경 변수를 통해:

```yaml
auth:
  type: bearer
  config:
    token: "$(MY_TOKEN)"
```

### invoke 전에 auth를 호출해야 하나요?

아니요. `invoke` 도구는 spec 설정에서 자동으로 인증을 적용합니다. `auth` MCP 도구는 사용자에게 토큰을 보여주려는 경우(예: curl 명령어 생성)에만 필요합니다.

### auth 도구가 보이지 않는 이유는 무엇인가요?

`auth` 도구는 기본적으로 비활성화되어 있습니다(`--disable-llm-auth=true`). 이는 프로덕션을 위한 보안 조치입니다. 활성화하려면: `swag2mcp mcp --disable-llm-auth=false`.

### OAuth2 토큰은 어떻게 갱신되나요?

OAuth2 Client Credentials 및 Password Grant 토큰은 만료 시 자동으로 갱신됩니다. Bearer 토큰은 정적이며 수동으로 업데이트해야 합니다.

## MCP 서버

### MCP 서버를 어떻게 시작하나요?

```bash
# 기본 (stdio 전송)
swag2mcp mcp

# HTTP 전송 사용
swag2mcp mcp --transport sse --http-addr :8080
```

### 포트를 어떻게 변경하나요?

```bash
swag2mcp mcp --transport sse --http-addr 0.0.0.0:9090
```

### MCP HTTP 엔드포인트를 어떻게 보호하나요?

Bearer 토큰을 설정하세요:

```bash
swag2mcp mcp --transport sse --http-addr :8080 --auth-token "my-secret"
```

LLM 클라이언트는 모든 요청에 `Authorization: Bearer my-secret`을 포함해야 합니다.

### HTTP 전송의 MCP 핸드셰이크는 무엇인가요?

SSE 및 Streamable HTTP 전송의 경우 MCP 프로토콜은 3단계 핸드셰이크가 필요합니다:

```
1단계: POST /mcp → {"method":"initialize", ...}
2단계: POST /mcp → {"method":"notifications/initialized"}
3단계: POST /mcp → {"method":"tools/list", ...}  ← 이제 작동
```

초기화 전에는 도구 호출이 실패합니다.

## 사용법

### 엔드포인트를 어떻게 검색하나요?

`search` MCP 도구 또는 TUI(`swag2mcp run`)를 사용하세요. 검색은 필드 필터(`method:GET`, `tag:pets`), 퍼지 검색, 와일드카드, 부울 연산자를 지원합니다.

### API를 어떻게 호출하나요?

LLM이 `invoke` MCP 도구를 사용합니다. 항상 먼저 엔드포인트를 검사하여 필요한 매개변수를 이해하세요:

```
inspect(endpointId: "...")  → 계약 이해
invoke(endpointId: "...", parameters: {...})  → 호출 실행
```

### 응답이 너무 크면 어떻게 되나요?

`max_response_size`(기본값 1MB)를 초과하는 응답은 디스크에 저장됩니다. LLM은 파일 참조를 받고 `response_outline`, `response_compress`, `response_slice` 도구로 탐색할 수 있습니다.

### 속도 제한기는 어떻게 작동하나요?

각 엔드포인트에는 10초의 쿨다운이 있습니다. LLM이 10초 이내에 동일한 엔드포인트를 두 번 호출하면 두 번째 호출은 자동으로 차단됩니다. 설정에서 비활성화하거나 조정할 수 있습니다.

### 실제 API 호출 없이 테스트할 수 있나요?

네, 모의 서버를 사용하세요:

```bash
swag2mcp-mock mockserver
```

OpenAPI 스키마를 기반으로 가짜 응답을 생성합니다.

## 워크스페이스 관리

### 설정을 어떻게 백업하나요?

```bash
swag2mcp export --output ~/backups/swag2mcp-2026-07-24.zip
```

### 다른 머신으로 전송하려면 어떻게 하나요?

```bash
# 이전 머신에서
swag2mcp export --output swag2mcp.zip

# ZIP을 복사한 후 새 머신에서
swag2mcp import --from-zip swag2mcp.zip
```

### 명세 파일을 어떻게 업데이트하나요?

```bash
swag2mcp update
```

설정을 재검증하고, 캐시를 지우고, 모든 명세 파일을 다시 다운로드합니다.

### 디스크 공간을 어떻게 정리하나요?

```bash
swag2mcp clean
```

캐시된 명세 파일과 저장된 API 응답을 제거합니다. 오래된 응답(48시간 초과)은 MCP 서버 시작 시 자동으로 정리됩니다.

## TUI

### TUI란 무엇이며 어떻게 사용하나요?

TUI(터미널 사용자 인터페이스)는 대화형 API 탐색기입니다. `swag2mcp run`으로 실행합니다. 세 가지 모드가 있습니다: 검색(전문 검색), 탐색(트리 탐색: Spec → Collection → Tag → Endpoint), 인증(토큰 보기).

### 키보드 단축키는 무엇인가요?

| 키 | 동작 |
|-----|------|
| `↑/↓` | 이동 |
| `Enter` | 선택 |
| `Esc` | 뒤로 |
| `Tab` | 모드 전환 |
| `/` | 검색 |
| `N/P` | 다음/이전 페이지 |
| `q` | 종료 |

## 고급

### 프록시를 사용할 수 있나요?

네, `http_client.proxy`에서 설정하세요:

```yaml
http_client:
  proxy:
    url: "http://proxy.company.com:8080"
    username: "$(PROXY_USER)"
    password: "$(PROXY_PASS)"
    bypass:
      - "localhost"
      - "*.internal.com"
```

### 커스텀 인증 방식을 추가할 수 있나요?

네, `internal/auth/`에서 `Authenticator` 인터페이스를 구현하고 설정 파서에 등록하세요. 자세한 내용은 개발 섹션을 참조하세요.

### 커스텀 MCP 도구를 추가할 수 있나요?

네, `Svc` 인터페이스에 메서드를 추가하고, 서비스 레이어에서 구현하고, 핸들러를 추가하고, 등록하세요. 자세한 내용은 개발 섹션을 참조하세요.

### `swag2mcp`와 `swag2mcp-mock`의 차이는 무엇인가요?

`swag2mcp`는 CLI 명령어와 MCP 서버가 포함된 메인 바이너리입니다. `swag2mcp-mock`은 실제 API 호출 없이 테스트하기 위한 모의 서버를 시작하는 별도의 바이너리입니다.
