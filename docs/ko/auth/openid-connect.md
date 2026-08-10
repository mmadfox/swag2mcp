# OpenID Connect (OIDC Discovery)

## 목적

OpenID Connect Discovery 기반 인증. swag2mcp는 issuer에서 OIDC 디스커버리 문서를 가져와 토큰 엔드포인트를 찾고 **client credentials** 그랜트로 Bearer 토큰을 얻습니다. 토큰은 만료될 때까지 캐시됩니다.

## 사용 시기

- OIDC 디스커버리 문서(`.well-known/openid-configuration`)를 제공하는 API
- `client_id` + `client_secret`이 있는 서버 간 통합
- API가 OIDC를 사용하고 자동 토큰 획득이 필요한 경우

## 구성

```yaml
specs:
  - domain: my-api
    llm_title: My API
    base_url: https://api.example.com
    collections:
      - llm_title: Main
        location: https://example.com/spec.yaml
    auth:
      type: openid-connect
      config:
        issuer: "https://auth.example.com"
        client_id: "$(CLIENT_ID)"
        client_secret: "$(CLIENT_SECRET)"
        scopes:
          - openid
          - profile
```

## 매개변수

| 매개변수 | 필수 | 설명 |
|-----------|----------|-------------|
| `issuer` | 예 | OIDC issuer URL. swag2mcp는 `/.well-known/openid-configuration`을 추가하여 토큰 엔드포인트를 찾습니다 |
| `client_id` | 예 | 클라이언트 식별자 |
| `client_secret` | 예 | 클라이언트 시크릿 |
| `scopes` | 아니요 | 권한 목록(선택 사항) |

## 참고 사항

- 현재 토큰이 만료되면 swag2mcp가 자동으로 새 토큰을 요청합니다
- 토큰은 만료(`expires_in`)까지 캐시됩니다
- 서버가 `expires_in`을 제공하지 않으면 토큰은 1시간 동안 유효합니다
- `issuer`, `client_id`, `client_secret`은 환경 변수에 `$(VAR)` 구문을 지원합니다
- 토큰은 **client credentials** 그랜트로 얻습니다(서버 간, 사용자 없음)
- 디스커버리 문서에는 `token_endpoint` 필드가 있어야 합니다
