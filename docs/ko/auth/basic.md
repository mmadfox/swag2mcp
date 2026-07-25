# Basic Auth

## 목적

HTTP Basic Authentication — 사용자 이름과 비밀번호로 인증하는 가장 간단한 방법입니다.

## 사용 시기

- Basic Auth만 지원하는 레거시 API
- 복잡한 토큰이 필요 없는 간단한 인증
- 내부 서비스

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
      type: basic
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

- 비밀번호는 `Authorization: Basic ...` 헤더에 Base64로 인코딩되어 전송됩니다 — 이는 **암호화가 아닙니다**. 항상 HTTPS를 사용하세요.
- 비밀번호를 환경 변수에 저장하세요: `password: "$(MY_PASSWORD)"`
