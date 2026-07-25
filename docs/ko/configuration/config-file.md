# 설정 파일

swag2mcp는 YAML 설정 파일을 사용합니다. `swag2mcp init`으로 생성됩니다.

## 위치

- **Linux/macOS**: `~/.swag2mcp/swag2mcp.yaml`
- **Windows**: `%USERPROFILE%\.swag2mcp\swag2mcp.yaml`

## 기본 구조

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

## 전체 예시

```yaml
# ── 전역 HTTP 클라이언트 ──────────────────────────────────
http_client:
  timeout: 30s
  max_response_size: 1048576
  user_agent: "swag2mcp-global/1.0"
  follow_redirects: true
  max_redirects: 10
  random: false
  proxy:
    url: ""
    username: ""
    password: ""
    bypass: []
  headers:
    "Accept": "application/json"
  cookies:
    - name: "session"
      value: "abc123"

# ── MCP 서버 ──────────────────────────────────────────
mcp:
  transport: stdio
  addr: ":8080"
  path: "/mcp"
  auth:
    token: ""

# ── 모의 서버 ─────────────────────────────────────────
mock_enabled: false
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092

# ── 속도 제한기 ────────────────────────────────────────
disable_ratelimiter: false
rate_limit_interval: 10s

# ── Specs ───────────────────────────────────────────────
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    llm_instruction: "날씨 예보 및 기후 데이터에 이 API를 사용하세요"
    base_url: https://api.open-meteo.com
    disable: false
    tags: ["weather", "climate"]
    http_client:
      timeout: 10s
    auth:
      type: bearer
      config:
        token: "$(TOKEN)"
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
        disable: false
        http_client:
          timeout: 5s

  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## 환경 변수

`$(VAR_NAME)` 구문을 사용하여 환경 변수를 참조하세요. swag2mcp는 시작 시 이를 해결합니다.

```yaml
specs:
  - domain: meteo
    auth:
      type: bearer
      config:
        token: "$(API_TOKEN)"

mcp:
  auth:
    token: "$(MCP_TOKEN)"
```

`$(VAR)`은 다음에서 해결됩니다:
- 인증 설정 필드: `token`, `username`, `password`, `client_id`, `client_secret`, `api_key`, `secret_key`, `domain`
- MCP 서버 인증 토큰: `mcp.auth.token`
- HTTP 클라이언트 헤더 및 쿠키 값

`$(VAR)`은 base URL 또는 collection location에서 **해결되지 않습니다**.

## 검증

```bash
# 기본 워크스페이스 검증 (~/.swag2mcp)
swag2mcp validate

# 커스텀 프로젝트 워크스페이스 검증
swag2mcp validate ./my-project
```

워크스페이스가 홈 디렉토리에 없는 경우(예: 프로젝트 저장소 내부) `validate`, `update`, `mcp` 또는 다른 명령어를 실행할 때 항상 경로를 지정하세요. 그렇지 않으면 swag2mcp가 기본 `~/.swag2mcp` 워크스페이스를 사용합니다.
