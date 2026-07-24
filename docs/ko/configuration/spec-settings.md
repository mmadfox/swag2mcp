# Spec 설정

Spec 설정은 API 서비스를 정의하고 해당 API에 대한 전역 설정을 재정의합니다. 각 spec은 하나의 논리적 API(예: "Open-Meteo Weather APIs")를 나타내며 여러 collection(명세 파일)을 포함할 수 있습니다.

## Spec 섹션

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    llm_instruction: "날씨 예보 및 기후 데이터에 이 API를 사용하세요"
    base_url: https://api.open-meteo.com
    disable: false
    tags: ["weather", "climate"]
    http_client:
      timeout: 10s
      max_response_size: 1024
    auth:
      type: bearer
      config:
        token: "$(TOKEN)"
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

## 매개변수

### domain

- **타입:** `string`
- **필수:** 예
- **설명:** 이 API spec의 고유 식별자. 내부적으로 spec을 참조하는 데 사용됩니다.
- **규칙:** 1-60자. 소문자(`a-z`), 숫자(`0-9`), 하이픈(`-`), 밑줄(`_`)만 허용.
- **예시:** `meteo`, `binance`, `my-api`

### llm_title

- **타입:** `string`
- **필수:** 예
- **설명:** LLM이 이 API를 참조할 때 사용하는 사람이 읽을 수 있는 이름. MCP 도구 응답에 표시됩니다.
- **규칙:** 5-120자. 문자, 숫자, 공백, 기본 구두점만 허용.
- **예시:** `Open-Meteo Weather APIs`, `Binance Market Data`

### llm_instruction

- **타입:** `string`
- **기본값:** `""`
- **설명:** 이 API 사용 방법에 대한 LLM 지침. API가 무엇을 하는지와 사용 시기를 설명합니다.
- **규칙:** 최대 500자. 문자, 숫자, 공백, 기본 구두점만 허용.
- **예시:** `"날씨 예보, 현재 상태 및 기후 데이터에 이 API를 사용하세요."`

### base_url

- **타입:** `string`
- **필수:** 예
- **설명:** 이 spec의 모든 API 요청에 대한 기본 URL. OpenAPI 명세의 엔드포인트 경로가 이 URL에 추가됩니다.
- **예시:** `https://api.open-meteo.com`, `https://api.binance.com`
- **참고:** 다른 collection이 다른 base URL을 사용하는 경우 collection 수준에서 재정의할 수 있습니다.

### disable

- **타입:** `bool`
- **기본값:** `false`
- **설명:** `true`일 때 이 spec은 MCP 도구에서 제외됩니다. 로드, 인덱싱되지 않으며 LLM이 사용할 수 없습니다.
- **사용 시기:** 설정에서 제거하지 않고 API를 임시로 비활성화합니다. 다운되었거나, 폐기되었거나, 유지보수 중인 API에 유용합니다.

### tags

- **타입:** `[]string` (문자열 배열)
- **기본값:** `[]`
- **설명:** spec 필터링용 태그. CLI 명령어(`ls`, `validate`, `mcp`, `update`)에서 `--tags` 플래그와 함께 사용됩니다.
- **예시:** `["public", "weather"]`, `["internal", "production"]`
- **효과:** `swag2mcp mcp --tags=public`을 실행하면 `public` 태그가 있는 spec만 로드됩니다.

### http_client

- **타입:** `object`
- **기본값:** 전역에서 상속
- **설명:** 이 spec의 전역 HTTP 클라이언트 설정을 재정의합니다. 전역 `http_client`의 모든 설정을 재정의할 수 있습니다: `timeout`, `max_response_size`, `user_agent`, `follow_redirects`, `max_redirects`, `random`, `proxy`, `headers`, `cookies`.
- **예시:**
  ```yaml
  http_client:
    timeout: 60s
    max_response_size: 4194304
    headers:
      "X-DC": "us-east-1"
  ```

### auth

- **타입:** `object`
- **기본값:** `none` (인증 없음)
- **설명:** 이 spec의 인증 설정입니다. 모든 9가지 방법과 해당 매개변수는 [인증](/auth/overview) 섹션을 참조하세요.
- **예시:**
  ```yaml
  auth:
    type: bearer
    config:
      token: "$(API_TOKEN)"
  ```

### collections

- **타입:** `[]object` (collection 배열)
- **필수:** 예 (최소 1개)
- **설명:** 이 spec에 속하는 OpenAPI/Swagger/Postman 명세 파일 목록입니다. 각 collection은 하나의 명세 파일입니다.
- **규칙:** spec당 1-30개 collection.
- **참고:** 모든 collection 매개변수는 [Collection 설정](./collection-settings)을 참조하세요.

## Spec 비활성화

비활성화된 spec은 로드되거나 인덱싱되지 않습니다. LLM이 보거나 사용할 수 없습니다.

```yaml
specs:
  - domain: old-api
    llm_title: Old API
    base_url: https://old-api.example.com
    disable: true
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## HTTP Client 재정의

전역 수준의 모든 `http_client` 설정은 spec 수준에서 재정의할 수 있습니다. Spec 값은 이 spec에 대해서만 전역 값보다 우선합니다.

```yaml
specs:
  - domain: slow-api
    llm_title: Slow API
    base_url: https://slow-api.example.com
    http_client:
      timeout: 120s
      max_response_size: 8388608
      headers:
        "X-DC": "us-east-1"
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Proxy 재정의

이 spec이 전역과 다른 프록시가 필요한 경우 spec 수준에서 설정하세요:

```yaml
specs:
  - domain: proxied-api
    llm_title: Proxied API
    base_url: https://api.example.com
    http_client:
      proxy:
        url: http://proxy.company.com:8080
        username: $(PROXY_USER)
        password: $(PROXY_PASS)
        bypass:
          - "*.local"
          - "10.0.0.0/8"
    collections:
      - llm_title: Main
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo.json
```
