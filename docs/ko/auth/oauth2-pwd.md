# OAuth2 Password Grant

## 목적

OAuth2 Resource Owner Password Grant — 사용자의 사용자 이름과 비밀번호를 사용한 인증입니다. 사용자가 자격 증명을 앱에 신뢰하는 퍼스트파티 애플리케이션에 적합합니다.

## 사용 시기

- 퍼스트파티 애플리케이션 (모바일, 웹)
- Keycloak 및 유사한 ID 제공자와의 통합
- API가 OAuth2 Password Grant를 지원할 때

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
      type: oauth2-pwd
      config:
        client_id: "$(CLIENT_ID)"
        client_secret: "$(CLIENT_SECRET)"
        username: "$(USERNAME)"
        password: "$(PASSWORD)"
        token_url: "https://auth.example.com/oauth/token"
        scopes:
          - openid
          - profile
        request_format: form
```

## 매개변수

| 매개변수 | 필수 | 설명 |
|---------|------|------|
| `client_id` | 예 | 클라이언트 식별자 |
| `username` | 예 | 사용자 이름 |
| `password` | 예 | 비밀번호 |
| `token_url` | 예 | 토큰 엔드포인트 URL |
| `client_secret` | 아니요 | 클라이언트 시크릿 (선택 사항, 공개 클라이언트용) |
| `scopes` | 아니요 | 권한 목록 (선택 사항) |
| `request_format` | 아니요 | 요청 형식: `form` (기본값, `application/x-www-form-urlencoded`) 또는 `json` (`application/json`) |

## 참고 사항

- `client_secret`은 선택 사항 — **공개 클라이언트**가 지원됩니다 (예: Keycloak)
- swag2mcp는 토큰이 만료되면 자동으로 갱신합니다
- 토큰은 만료까지 캐시됩니다
- `client_id`, `client_secret`, `username`, `password`는 `$(VAR)` 구문을 통한 환경 변수를 지원합니다
- `token_url`과 `scopes`는 그대로 사용됩니다 (환경 변수 해결 없음)
- `request_format: json`은 토큰 요청을 JSON 본문으로 전송합니다 — 토큰 엔드포인트가 `Content-Type: application/json`을 필요로 할 때 사용하세요
