# MCP 서버

MCP 서버는 LLM 에이전트의 주요 상호 작용 지점입니다. 설정된 모든 API를 LLM이 호출할 수 있는 MCP 도구로 노출합니다.

## 설정

```yaml
mcp:
  transport: stdio
```

## 전송 방식

세 가지 전송 유형을 사용할 수 있습니다:

| 전송 | 설명 | 사용 시기 |
|------|------|----------|
| `stdio` | 표준 입출력 | 로컬 LLM 클라이언트 (VS Code, Cursor, Claude Desktop) |
| `sse` | Server-Sent Events | 원격 클라이언트, HTTP 기반 통신 |
| `streamable-http` | HTTP 스트리밍 | 웹 클라이언트, 최신 MCP 클라이언트 |

### stdio (기본값)

LLM 클라이언트가 swag2mcp를 하위 프로세스로 실행합니다. 통신은 표준 입력과 출력을 통해 이루어집니다. 네트워크 포트가 필요하지 않습니다.

```yaml
mcp:
  transport: stdio
```

```bash
swag2mcp mcp
```

### SSE

HTTP 기반 통신을 위한 Server-Sent Events 전송입니다. MCP 서버가 HTTP 포트에서 수신하고 LLM 클라이언트가 원격으로 연결합니다.

```yaml
mcp:
  transport: sse
  addr: "127.0.0.1:8080"
  path: "/mcp"
```

```bash
swag2mcp mcp --transport sse --http-addr 127.0.0.1:8080
```

### Streamable HTTP

스트리밍 응답을 지원하는 최신 HTTP 전송입니다. SSE와 유사하지만 다른 프로토콜을 사용합니다.

```yaml
mcp:
  transport: streamable-http
  addr: "127.0.0.1:8080"
  path: "/mcp"
```

```bash
swag2mcp mcp --transport streamable-http --http-addr 0.0.0.0:8080
```

## 매개변수

### transport

- **타입:** `string`
- **기본값:** `"stdio"`
- **옵션:** `stdio`, `sse`, `streamable-http`
- **효과:** MCP 서버가 LLM 클라이언트와 통신하는 방식을 결정합니다.

### addr

- **타입:** `string`
- **기본값:** `":8080"`
- **설명:** SSE 및 Streamable HTTP 전송의 수신 주소입니다. 형식: `host:port`.
- **예시:** `":8080"`, `"127.0.0.1:8080"`, `"0.0.0.0:9000"`

### path

- **타입:** `string`
- **기본값:** `"/mcp"`
- **설명:** MCP 엔드포인트의 URL 경로입니다. LLM 클라이언트가 `http://&lt;addr&gt;&lt;path&gt;`로 요청을 보냅니다.
- **예시:** `"/mcp"`, `"/api/mcp"`, `"/v1/mcp"`

### auth.token

- **타입:** `string`
- **기본값:** `""` (인증 없음)
- **설명:** HTTP 전송 인증용 Bearer 토큰입니다. 설정되면 LLM 클라이언트가 모든 요청에 `Authorization: Bearer &lt;token&gt;`을 포함해야 합니다.
- **참고:** `$(ENV_VAR)` 해결을 지원합니다.

### auth.type

- **Type:** `string`
- **Default:** `""` (no JWT auth)
- **Options:** `jwks`, `oidc`, `introspection`
- **Description:** JWT authentication type for HTTP transport. When set, enables dynamic token verification using JWKS, OIDC Discovery, or token introspection.

### auth.jwks_url

- **Type:** `string`
- **Default:** `""`
- **Description:** URL of the JWKS (JSON Web Key Set) endpoint. Required when `auth.type` is `jwks` or resolved via OIDC discovery.

### auth.issuer

- **Type:** `string`
- **Default:** `""`
- **Description:** Expected JWT issuer (`iss` claim). If set, tokens with a different issuer are rejected.

### auth.audience

- **Type:** `string`
- **Default:** `""`
- **Description:** Expected JWT audience (`aud` claim). If set, tokens without this audience are rejected.

### auth.introspection_url

- **Type:** `string`
- **Default:** `""`
- **Description:** Token introspection endpoint URL. Required when `auth.type` is `introspection`.

### auth.client_id

- **Type:** `string`
- **Default:** `""`
- **Description:** Client ID for introspection auth. Required when `auth.type` is `introspection`.

### auth.client_secret

- **Type:** `string`
- **Default:** `""`
- **Description:** Client secret for introspection auth. Supports `$(ENV_VAR)` resolution.

## HTTP 인증

Bearer 토큰으로 MCP HTTP 엔드포인트를 보호하세요:

```yaml
mcp:
  auth:
    token: "my-secret-token"
```

또는 CLI 플래그를 통해:

```bash
swag2mcp mcp --auth-token "my-secret-token"
```

### With JWT authentication (JWKS)

Protect the MCP HTTP endpoint with JWT verification via a JWKS endpoint:

```yaml
mcp:
  auth:
    type: jwks
    jwks_url: "https://auth.example.com/.well-known/jwks.json"
    issuer: "https://auth.example.com/"
    audience: "swag2mcp"
```

```bash
swag2mcp mcp --transport sse --http-addr 0.0.0.0:8080 \
  --auth-type jwks \
  --auth-jwks-url "https://auth.example.com/.well-known/jwks.json" \
  --auth-issuer "https://auth.example.com/" \
  --auth-audience "swag2mcp"
```

### With JWT authentication (OIDC Discovery)

```yaml
mcp:
  auth:
    type: oidc
    issuer: "https://auth.example.com/"
    audience: "swag2mcp"
```

### With JWT authentication (Token Introspection)

```yaml
mcp:
  auth:
    type: introspection
    introspection_url: "https://auth.example.com/introspect"
    client_id: "my-client"
    client_secret: "$(MCP_CLIENT_SECRET)"
```

```bash
swag2mcp mcp --transport sse --http-addr 0.0.0.0:8080 \
  --auth-type introspection \
  --auth-introspection-url "https://auth.example.com/introspect" \
  --auth-client-id "my-client" \
  --auth-client-secret "$(MCP_CLIENT_SECRET)"
```

## Health Check

MCP 서버는 MCP 초기화 없이 작동하는 health check 엔드포인트를 제공합니다:

```bash
curl http://127.0.0.1:8080/health
# {"status":"ok","version":"v1.2.0"}
```

## 시작 플래그

CLI 플래그는 YAML 설정을 재정의합니다. 플래그가 설정되지 않으면 YAML의 `mcp` 섹션 값이 폴백으로 사용됩니다.

| 플래그 | 타입 | 기본값 | 설명 |
|-------|------|--------|------|
| `--transport` | string | `"stdio"` | 전송 유형: `stdio`, `sse`, `streamable-http` |
| `--http-addr` | string | `":8080"` | HTTP 서버 주소 (SSE 및 Streamable HTTP용) |
| `--http-path` | string | `"/mcp"` | MCP 핸들러의 URL 경로 |
| `--auth-token` | string | `""` | HTTP 전송 인증용 Bearer 토큰 |
| `--logfile` | string | `""` | 로그 파일 경로 (설정하지 않으면 stderr로 로깅) |
| `--disable-llm-auth` | bool | `true` | MCP 도구 목록에서 `auth` 도구 제거 |
| `--dump-dir` | string | `""` | 디버깅용 HTTP 요청 덤프 디렉토리 |
| `--tags` | string | `""` | 태그로 spec 필터링 (쉼표로 구분) |
| `--auth-type` | string | `""` | JWT auth type: `jwks`, `oidc`, `introspection` |
| `--auth-jwks-url` | string | `""` | JWKS URL for JWT auth |
| `--auth-issuer` | string | `""` | JWT issuer for token validation |
| `--auth-audience` | string | `""` | JWT audience for token validation |
| `--auth-introspection-url` | string | `""` | Token introspection URL |
| `--auth-client-id` | string | `""` | Client ID for introspection auth |
| `--auth-client-secret` | string | `""` | Client secret for introspection auth |
