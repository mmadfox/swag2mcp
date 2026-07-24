# None

## 목적

인증이 필요하지 않습니다. 토큰이나 키 없이 API에 접근할 수 있습니다.

## 사용 시기

- 공개 API (Open-Meteo, icanhazdadjoke, PokéAPI)
- 테스트 및 데모 환경
- API가 인증을 요구하지 않을 때

## 설정

`type: none`으로 설정하거나 `auth` 섹션을 생략하세요:

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: none
```

## 매개변수

없음.

## 참고 사항

- 설정에서 `auth` 섹션이 완전히 없으면 `type: none`과 동일합니다
- 요청에 인증 헤더가 추가되지 않습니다
