# 인증

## 개요

swag2mcp는 인증이 필요한 API 작업을 위한 **9가지 인증 방법**을 지원합니다. 설정 파일에 한 번 구성하면 `invoke`를 통한 모든 API 호출에 올바른 토큰과 헤더가 자동으로 포함됩니다.

### 설정 위치

인증은 `swag2mcp.yaml`의 **spec** 수준에서 설정됩니다:

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: bearer
      config:
        token: "my-token"
```

### 작동 방식

- 설정에서 인증 유형과 매개변수를 지정합니다
- swag2mcp는 `invoke`를 호출할 때 모든 요청에 자동으로 적용합니다
- API를 호출하기 전에 토큰을 요청할 **필요가 없습니다** — 자동으로 처리됩니다
- 토큰이 만료되면(OAuth2, Script) swag2mcp가 자체적으로 갱신합니다

### 환경 변수

민감한 데이터(토큰, 비밀번호, 키)는 `$(VAR_NAME)` 구문을 사용하여 환경 변수에 저장할 수 있습니다:

```yaml
auth:
  type: bearer
  config:
    token: "$(MY_API_TOKEN)"
```

swag2mcp는 시작 시 `MY_API_TOKEN`의 값을 치환합니다.

### MCP auth 도구

LLM 에이전트는 `auth` MCP 도구를 통해 토큰이나 헤더를 검색할 수 있습니다 — 예를 들어, curl 명령어를 만들거나 사용자에게 보여주는 용도입니다.

**프로덕션**에서는 이 도구를 `--disable-llm-auth`로 비활성화해야 합니다(기본적으로 활성화됨). 이렇게 하면 LLM이 토큰에 접근할 수 없습니다.

### 방법

| 방법 | 설명 | 최적 용도 |
|------|------|----------|
| [`none`](/auth/none) | 인증 없음 | 공개 API |
| [`basic`](/auth/basic) | HTTP Basic (username + password) | 레거시 API, 간단한 인증 |
| [`bearer`](/auth/bearer) | Bearer Token (JWT, token) | 현대 REST API |
| [`api-key`](/auth/api-key) | 헤더 또는 쿼리 매개변수의 API 키 | API 키가 있는 서비스 |
| [`digest`](/auth/digest) | HTTP Digest (username + password) | 레거시 API, Basic보다 안전 |
| [`hmac`](/auth/hmac) | HMAC-SHA256 서명 (Binance 스타일) | 암호화폐 거래소 |
| [`oauth2-cc`](/auth/oauth2-cc) | OAuth2 Client Credentials | 서버 간, 마이크로서비스 |
| [`oauth2-pwd`](/auth/oauth2-pwd) | OAuth2 Password Grant | 사용자 로그인이 있는 앱 |
| [`script`](/auth/script) | 토큰 획득을 위한 외부 스크립트 | 모든 커스텀 인증 체계 |
