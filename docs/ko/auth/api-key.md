# API Key

## 목적

API 키를 통한 인증입니다. 키는 HTTP 헤더 또는 URL 쿼리 매개변수로 전송할 수 있습니다.

## 사용 시기

- API 키를 사용하는 서비스
- 날씨 서비스, 지리 데이터, 번역 API
- API가 헤더(`X-API-Key`) 또는 쿼리 매개변수(`?api_key=...`)에서 키를 기대할 때

## 설정

### 헤더의 키

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: api-key
      config:
        key: "X-API-Key"
        in: header
        value: "$(API_KEY)"
```

### 쿼리 매개변수의 키

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: api-key
      config:
        key: "api_key"
        in: query
        value: "$(API_KEY)"
```

## 매개변수

| 매개변수 | 필수 | 설명 |
|---------|------|------|
| `key` | 예 | 헤더 또는 쿼리 매개변수의 이름 |
| `in` | 예 | 키 위치: `header` 또는 `query` |
| `value` | 예 | 키 값 |

## 참고 사항

- `header` 모드에서 키는 HTTP 헤더로 추가됩니다
- `query` 모드에서 키는 URL 매개변수로 추가됩니다
- 값을 환경 변수에 저장하세요: `value: "$(MY_API_KEY)"`
