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
```

## 매개변수

| 매개변수 | 필수 | 설명 |
|---------|------|------|
| `client_id` | 예 | 클라이언트 식별자 |
| `client_secret` | 예 | 클라이언트 시크릿 |
| `token_url` | 예 | 토큰 엔드포인트 URL |
| `scopes` | 아니요 | 권한 목록 (선택 사항) |

## 참고 사항

- swag2mcp는 현재 토큰이 만료되면 자동으로 새 토큰을 요청합니다
- 토큰은 만료 시간(`expires_in`)까지 캐시됩니다
- 서버가 `expires_in`을 제공하지 않으면 토큰은 1시간 동안 유효한 것으로 간주됩니다
- 모든 매개변수를 환경 변수에 저장할 수 있습니다
