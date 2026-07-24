# Collection 설정

Collection 설정은 단일 OpenAPI/Swagger/Postman 명세 파일을 정의하고 해당 파일에 대한 spec 설정을 재정의합니다. 각 collection은 spec에 속하며 하나의 API 명세 문서를 나타냅니다.

## Collection 섹션

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        llm_instruction: "현재 및 예보 날씨 데이터에 사용"
        disable: false
        base_url: https://forecast-api.open-meteo.com
        base_mock_url: localhost:8081
        http_client:
          timeout: 5s
```

## 매개변수

### llm_title

- **타입:** `string`
- **필수:** 아니요
- **설명:** 이 collection의 사람이 읽을 수 있는 이름. MCP 도구 응답에 표시됩니다.
- **규칙:** 최대 120자. 문자, 숫자, 공백, 기본 구두점만 허용.
- **예시:** `Forecast`, `Air Quality`, `Market Data`

### llm_instruction

- **타입:** `string`
- **기본값:** `""`
- **설명:** 이 특정 collection에 대한 LLM 지침. 이 collection이 제공하는 엔드포인트를 설명합니다.
- **규칙:** 최대 360자. 문자, 숫자, 공백, 기본 구두점만 허용.
- **예시:** `"현재 및 예보 날씨 데이터에 사용하세요."`

### title

- **타입:** `string`
- **기본값:** `""`
- **설명:** 명세 파일의 원시 제목. 런타임에 자동으로 채워집니다. 일반적으로 YAML에서 설정할 필요가 없습니다.

### location

- **타입:** `string`
- **필수:** 예
- **설명:** OpenAPI 3.x, Swagger 2.0 또는 Postman collection 명세 파일의 URL 또는 로컬 파일 경로.
- **규칙:** 5-250자.
- **예시:**
  - URL: `https://raw.githubusercontent.com/org/repo/main/spec.yaml`
  - 로컬: `./specs/my-api.json`
  - 로컬 (절대 경로): `/home/user/.swag2mcp/specs/my-api.yaml`

### disable

- **타입:** `bool`
- **기본값:** `false`
- **설명:** `true`일 때 이 collection은 MCP 도구에서 제외됩니다. 로드되거나 인덱싱되지 않습니다.
- **사용 시기:** 설정에서 제거하지 않고 collection을 임시로 비활성화합니다. 명세 파일이 업데이트 중이거나 API 버전이 폐기된 경우 유용합니다.

### http_client

- **타입:** `object`
- **기본값:** spec에서 상속 (또는 전역)
- **설명:** 이 collection의 HTTP 클라이언트 설정을 재정의합니다. 전역 `http_client`의 모든 설정을 재정의할 수 있습니다: `timeout`, `max_response_size`, `user_agent`, `follow_redirects`, `max_redirects`, `random`, `proxy`, `headers`, `cookies`.
- **예시:**
  ```yaml
  http_client:
    timeout: 120s
    headers:
      "X-Custom": "value"
    cookies:
      - name: "session"
        value: "abc123"
  ```

### base_url

- **타입:** `string`
- **기본값:** `""` (spec에서 상속)
- **설명:** 이 collection의 spec 수준 `base_url`을 재정의합니다. 동일한 spec 내의 다른 collection이 다른 base URL을 사용할 때 사용합니다.
- **예시:** spec에 `base_url: https://api.open-meteo.com`이 있지만 한 collection이 `https://air-quality-api.open-meteo.com`을 사용하는 경우 collection 수준에서 `base_url`을 설정하세요.

### base_mock_url

- **타입:** `string`
- **기본값:** `""`
- **설명:** `host:port` 형식의 모의 서버 주소. 전역 설정에서 `mock_enabled: true`일 때 필수입니다.
- **규칙:** Host는 `localhost`, `127.0.0.1` 또는 `0.0.0.0`이어야 합니다. Port는 유효한 포트 번호여야 합니다.
- **예시:** `localhost:8081`, `127.0.0.1:9000`
- **사용 시기:** `mock_enabled: true`이고 이 collection을 가짜 응답으로 테스트하려고 할 때.

## 하나의 Spec에서 여러 Collection

spec은 여러 collection을 가질 수 있습니다 — 예를 들어, API가 다른 서비스에 대해 별도의 명세 파일을 가질 때:

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
      - llm_title: Marine
        base_url: https://marine-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

## Collection 비활성화

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
        disable: true
```

## HTTP Client 재정의

모든 `http_client` 설정은 collection 수준에서 재정의할 수 있습니다. Collection 값은 이 collection에 대해서만 spec 및 전역 값보다 우선합니다.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          timeout: 120s
          headers:
            "X-Custom": "value"
          cookies:
            - name: "session"
              value: "abc123"
```
