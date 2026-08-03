# OAuth2 Client Credentials

## 목적

OAuth2 Client Credentials Grant — 서버 간 통신을 위한 인증입니다. 애플리케이션이 사용자 개입 없이 client_id와 client_secret을 사용하여 토큰을 획득합니다.

## 사용 시기

- 마이크로서비스 및 서버 간 통합
- 머신 간 통신
- API가 OAuth2를 사용하고 client_id + client_secret이 있을 때

## 설정

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: oauth2-cc
      config:
        client_id: "$(CLIENT_ID)"
        client_secret: "$(CLIENT_SECRET)"
        token_url: "https://auth.example.com/oauth/token"
        scopes:
          - read
          - write
        request_format: form
```

## 매개변수

| 매개변수 | 필수 | 설명 |
|---------|------|------|
| `client_id` | 예 | 클라이언트 식별자 |
| `client_secret` | 예 | 클라이언트 시크릿 |
| `token_url` | 예 | 토큰 엔드포인트 URL |
| `scopes` | 아니요 | 권한 목록 (선택 사항) |
| `request_format` | 아니요 | 요청 형식: `form` (기본값, `application/x-www-form-urlencoded`) 또는 `json` (`application/json`) |

## 참고 사항

- swag2mcp는 현재 토큰이 만료되면 자동으로 새 토큰을 요청합니다
- 토큰은 만료 시간(`expires_in`)까지 캐시됩니다
- 서버가 `expires_in`을 제공하지 않으면 토큰은 1시간 동안 유효한 것으로 간주됩니다
- `client_id`와 `client_secret`은 `$(VAR)` 구문을 통한 환경 변수를 지원합니다
- `token_url`과 `scopes`는 그대로 사용됩니다 (환경 변수 해결 없음)
- `request_format: json`은 토큰 요청을 JSON 본문으로 전송합니다 — 토큰 엔드포인트가 `Content-Type: application/json`을 필요로 할 때 사용하세요
