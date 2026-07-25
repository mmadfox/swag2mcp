# Bearer Auth

## 목적

Bearer Token 인증 — 현대 REST API에서 가장 일반적인 방법입니다. 토큰은 `Authorization: Bearer &lt;token&gt;` 헤더로 전송됩니다.

## 사용 시기

- 현대 REST API
- JWT (JSON Web Tokens)
- OAuth2 액세스 토큰 (토큰이 이미 획득된 경우)
- Bearer Token을 허용하는 모든 API

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
      type: bearer
      config:
        token: "eyJhbGciOiJIUzI1NiIs..."
```

## 매개변수

| 매개변수 | 필수 | 설명 |
|---------|------|------|
| `token` | 예 | Bearer 토큰 (JWT, OAuth2 토큰 등) |

## 참고 사항

- 토큰은 정적입니다 — 만료되면 설정에서 수동으로 업데이트해야 합니다
- 자동 토큰 갱신을 위해서는 `oauth2-cc` 또는 `oauth2-pwd`를 사용하세요
- 토큰을 환경 변수에 저장하세요: `token: "$(API_TOKEN)"`
