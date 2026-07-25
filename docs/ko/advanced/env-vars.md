# 환경 변수

## 개요

swag2mcp는 `$(VAR_NAME)` 구문을 사용하여 설정 파일에서 환경 변수 치환을 지원합니다. 이를 통해 민감한 데이터(토큰, 비밀번호, 키)를 YAML 파일 밖에 보관할 수 있습니다.

## 작동 방식

swag2mcp가 시작되면 설정에서 `$(VAR_NAME)` 패턴을 스캔하고 해당 환경 변수의 값으로 대체합니다.

```yaml
specs:
  - domain: meteo
    auth:
      type: bearer
      config:
        token: "$(API_TOKEN)"
```

환경 변수 `API_TOKEN`이 설정되어 있으면 치환됩니다. 설정되어 있지 않으면 값이 비어 있게 됩니다.

## `$(VAR)`이 해결되는 위치

| 필드 | 예시 |
|------|------|
| 인증 `token` (bearer) | `token: "$(API_TOKEN)"` |
| 인증 `username` / `password` (basic, digest) | `password: "$(API_PASSWORD)"` |
| 인증 `client_id` / `client_secret` (oauth2-cc, oauth2-pwd) | `client_secret: "$(OAUTH_SECRET)"` |
| 인증 `api_key` / `secret_key` (hmac) | `api_key: "$(BINANCE_API_KEY)"` |
| 인증 `domain` (script) | `domain: "$(AUTH_DOMAIN)"` |
| MCP 서버 토큰 | `token: "$(MCP_TOKEN)"` |
| HTTP 클라이언트 헤더 | `"X-API-Key": "$(API_KEY)"` |
| HTTP 클라이언트 쿠키 값 | `value: "$(SESSION_TOKEN)"` |

## `$(VAR)`이 해결되지 않는 위치

- Base URL (`base_url`)
- Collection 위치 (`location`)
- Spec 도메인 이름 (`domain`)

## 예시

```bash
export API_TOKEN="eyJhbGciOiJIUzI1NiIs..."
export MCP_TOKEN="my-secret-token"

swag2mcp mcp
```

## 보안 모범 사례

- **절대** 시크릿을 YAML 파일에 직접 저장하지 마세요
- 환경 변수 또는 외부 시크릿 관리자를 사용하세요
- 하드코딩된 시크릿이 포함된 경우 YAML 파일을 `.gitignore`에 추가하세요
- 셸 프로필, IDE 설정 또는 배포 파이프라인에서 환경 변수를 설정하세요

## 구문 세부 사항

- `$(VAR_NAME)` — 표준 구문
- `$( VAR_NAME )` — 괄호 안의 공백은 허용되며 제거됩니다
- `$()` — 빈 변수 이름은 원래 문자열을 변경하지 않고 반환합니다
- 중첩된 `$(...)` 패턴은 해결되지 않습니다
