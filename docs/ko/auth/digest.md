# Digest Auth

## 목적

HTTP Digest Access Authentication — Basic Auth보다 더 안전한 대안입니다. 비밀번호가 평문으로 전송되지 않고 MD5 해시가 사용됩니다.

## 사용 시기

- Digest만 지원하는 레거시 API
- 비밀번호를 평문으로 전송하지 않고 인증이 필요할 때
- 내부 엔터프라이즈 시스템

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
      type: digest
      config:
        username: "admin"
        password: "$(PASSWORD)"
```

## 매개변수

| 매개변수 | 필수 | 설명 |
|---------|------|------|
| `username` | 예 | 사용자 이름 |
| `password` | 예 | 비밀번호 |

## 참고 사항

- swag2mcp는 먼저 인증 없이 요청을 보내고, 서버로부터 챌린지(HTTP 401)를 받은 후 응답을 계산하고 `Authorization: Digest ...` 헤더로 재시도합니다
- 챌린지는 5분 동안 캐시됩니다 — 이후 요청은 추가 왕복이 필요하지 않습니다
- 비밀번호를 환경 변수에 저장하세요: `password: "$(API_PASSWORD)"`
