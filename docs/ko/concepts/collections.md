# Collections

Collection은 특정 API를 설명하는 단일 OpenAPI/Swagger/Postman 파일입니다. `location`(URL 또는 로컬 파일 경로)을 가리키며 spec(도메인)에 속합니다.

하나의 spec은 여러 collection을 가질 수 있습니다 — 예를 들어, "meteo" spec에는 각각 다른 명세 파일을 가리키는 "Forecast", "Air Quality", "Marine" collection이 있을 수 있습니다.

## Collection 필드

| 필드 | YAML 키 | 필수 | 설명 |
|------|---------|------|------|
| [LLM Title](#llm-instruction) | `llm_title` | ❌ | LLM용 collection 표시 이름 (최대 120자). 설정되지 않으면 명세 문서에서 자동 채움 |
| [LLM Instruction](#llm-instruction) | `llm_instruction` | ❌ | LLM용 짧은 힌트 (최대 360자). 설정되지 않으면 명세 문서에서 자동 채움 |
| Title | `title` | ❌ | 원본 명세 제목 재정의 (파싱된 문서에서 자동 채움) |
| [Location](#location--how-spec-files-are-resolved) | `location` | ✅ | 명세 파일의 URL 또는 경로 (5–250자) |
| [Disable](#disable) | `disable` | ❌ | 로딩 중 이 collection 건너뛰기 |
| [HTTP Client](#http-client-override) | `http_client` | ❌ | Collection별 HTTP 설정 (헤더, 쿠키) |
| [Base URL](#base-url-override) | `base_url` | ❌ | 이 collection의 spec base URL 재정의 |
| [Mock Server](#mock-server) | `base_mock_url` | ❌ | `host:port` 형식의 모의 서버 주소. `mock_enabled: true`일 때 필수 |

## Location — 명세 파일 해결 방식

`location` 필드는 swag2mcp에 OpenAPI/Swagger/Postman 파일을 찾을 위치를 알려줍니다. 여러 소스 유형을 지원합니다:

| 소스 | 예시 | 설명 |
|------|------|------|
| **원격 URL** | `https://raw.githubusercontent.com/.../spec.yaml` | 다운로드 및 캐시 |
| **로컬 파일 (절대 경로)** | `/home/user/my-api.yaml` | 파일 시스템에서 읽기, 캐시 |
| **로컬 파일 (상대 경로)** | `./my-api.yaml` | 절대 경로로 해결, 캐시 |
| **워크스페이스 로컬 파일** | `specs/my-api.yaml` | `~/.swag2mcp/specs/`에 저장, 직접 사용 (캐시되지 않음) |
| **file:// URI** | `file:///home/user/spec.yaml` | 로컬 경로로 변환, 캐시 |

swag2mcp가 자동으로 소스 유형을 감지합니다:

- `https://` 또는 `http://` → 원격 URL (캐시됨)
- `file://` → 로컬 파일 (파일 시스템 경로로 변환)
- 그 외 → 로컬 파일 (`~`는 홈 디렉토리로 확장)

### 원격 URL

원격 URL을 사용하면 swag2mcp가 파일을 다운로드하여 로컬에 캐시합니다. 이후 시작 시 반복 다운로드를 피하기 위해 캐시가 재사용됩니다.

### 로컬 파일

로컬 파일은 파일 시스템에서 직접 읽습니다. 파일이 워크스페이스 `specs/` 디렉토리 외부에 있으면 일관성을 위해 캐시로 복사됩니다.

### 워크스페이스 로컬 파일

워크스페이스 내부의 `specs/` 디렉토리(`~/.swag2mcp/specs/`)는 로컬 명세 파일을 위한 권장 위치입니다. 여기에 저장된 파일은 캐싱 없이 직접 사용됩니다. 참조하려면 `specs/`로 시작하는 상대 경로를 사용하세요.

> **참고:** `specs/`는 디렉토리 이름(예: `cache/` 또는 `responses/`)일 뿐, "spec" 개념이 아닙니다. collection이 가리키는 실제 OpenAPI/Swagger/Postman 파일을 저장합니다.

```bash
# 명세 파일을 워크스페이스로 가져오기
swag2mcp import https://example.com/api.yaml myspec

# 가져오기 후 location이 다음과 같이 변경됨:
# specs/myspec.yaml
```

## 캐시 시스템

swag2mcp는 매번 시작할 때마다 다운로드하지 않도록 원격 명세 파일을 캐시합니다.

### 작동 방식

1. 원격 URL이 있는 collection이 로드되면 swag2mcp가 캐시를 확인합니다
2. 유효한(만료되지 않은) 캐시 항목이 있으면 직접 사용됩니다
3. 없으면 파일이 다운로드, 파싱되어 캐시에 저장됩니다

### 캐시 구조

```
~/.swag2mcp/
  cache/
    {sha256_hash}.spec    # 캐시된 명세 파일 내용
    {sha256_hash}.meta    # 캐시 메타데이터 (JSON)
```

각 캐시 파일에는 다음을 포함하는 메타데이터 파일이 있습니다:

```json
{
  "source": "https://example.com/api.yaml",
  "source_type": "url",
  "cached_at": "2024-01-01T00:00:00Z",
  "mod_time": "2024-01-01T00:00:00Z",
  "ttl_sec": 3600
}
```

### 캐시 TTL

각 캐시 파일은 **1시간에서 48시간** 사이의 **무작위 TTL**을 받습니다. 이는 모든 캐시 파일이 동시에 만료되는 것을 방지합니다(폭주 문제 방지).

### 캐시 키

캐시 키는 원시 location 문자열의 SHA-256 해시입니다(처음 16바이트 = 32자 16진수).

### 캐시 관리

```bash
# 캐시와 응답 지우기, 모든 명세 파일 다시 다운로드
swag2mcp update

# 캐시와 응답만 지우기
swag2mcp clean
```

- `swag2mcp update` — 설정 검증, `cache/`와 `responses/` 지우기, 모든 collection location 재캐싱
- `swag2mcp clean` — `cache/`와 `responses/`의 모든 내용 제거, 고아 인증 스크립트도 제거
- 오래된 응답은 MCP 서버 시작 후 48시간이 지나면 자동으로 정리됨

## 검증

모든 collection은 설정이 로드될 때 검증됩니다. 검증은 모든 `swag2mcp mcp` 시작 시 실행됩니다. 실패하면 MCP 서버가 시작되지 않습니다 — 일부 IDE에서는 서버가 연결되지 않고 LLM이 수정해야 할 사항을 설명하는 명확한 오류 메시지를 받게 됩니다.

| 확인 | 규칙 |
|------|------|
| **Location** | 필수, 5–250자 |
| **Location 접근성** | 접근 가능한 URL 또는 존재하는 파일이어야 함 |
| **Location 유효성** | 유효한 OpenAPI 3.x, Swagger 2.0 또는 Postman 파일이어야 함 |
| **LLM Title** | 최대 120자, 문자/숫자/기본 구두점 |
| **LLM Instruction** | 최대 360자, title과 동일한 문자 세트 |
| **Base URL** | 설정된 경우 유효한 URL이어야 함 |
| **Base Mock URL** | `host:port` 또는 `host:port/path` 형식이어야 하며 host는 `localhost`, `127.0.0.1`, 또는 `0.0.0.0` |
| **Mock 필수** | `mock_enabled: true`이면 모든 collection에 `base_mock_url`이 있어야 함 |
| **중복 모의 포트** | 두 collection이 동일한 모의 포트를 공유할 수 없음 |

서버를 시작하기 전에 문제를 진단하려면 [`validate`](../cli/validate.md) 명령어를 사용하세요:

```bash
# 기본 워크스페이스 검증 (~/.swag2mcp)
swag2mcp validate

# 커스텀 프로젝트 워크스페이스 검증
swag2mcp validate ./my-project
```

## Collection 추가

### YAML 설정을 통해

`~/.swag2mcp/swag2mcp.yaml`을 직접 편집:

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

편집 후 변경 사항을 적용하려면 MCP 서버를 다시 시작하세요(`swag2mcp mcp`).

### CLI를 통해

```bash
# 대화형 모드
swag2mcp add collection

# YAML로 비대화형
swag2mcp add collection --yaml 'spec_domain: meteo
llm_title: Forecast
location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml'

# stdin에서 파이프
cat collection.yaml | swag2mcp add collection --yaml -

# YAML 예시 보기
swag2mcp add collection --example
```

### 가져오기를 통해

```bash
# 명세 파일을 워크스페이스로 가져오기
swag2mcp import https://example.com/api.yaml
```

## LLM Instruction

Collection은 더 구체적인 지침을 위해 자체 `llm_instruction`(최대 360자)를 가질 수 있습니다. 이는 spec 수준 지침과 함께 swag2mcp 시스템 프롬프트에 주입됩니다.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        llm_instruction: "현재 날씨 및 일일 예보에 이 collection을 사용하세요."
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        llm_instruction: "대기질 지수 및 오염 데이터에 이 collection을 사용하세요."
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
```

`llm_title`이 설정되지 않으면 명세 문서의 `title` 필드에서 자동으로 채워집니다. `llm_instruction`이 설정되지 않으면 명세 문서의 `description` 필드에서 채워집니다.

## Disable

`disable: true`로 설정하면 collection을 건너뜁니다. 로드, 인덱싱되지 않으며 LLM이 사용할 수 없습니다.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
        disable: true
```

## Base URL 재정의

각 collection은 spec의 `base_url`을 재정의할 수 있습니다. 동일한 spec 내의 다른 collection이 다른 API 엔드포인트를 사용할 때 유용합니다.

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

## HTTP Client 재정의

Collection은 spec 및 전역 수준의 HTTP 설정(헤더, 쿠키)을 재정의할 수 있습니다.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          headers:
            X-API-Version: "2"
          cookies:
            - name: session
              value: abc123
```

설정은 전역 → spec → collection 순으로 계단식으로 적용됩니다. 자세한 내용은 [설정 계단식](../configuration/cascade.md)을 참조하세요.

## 모의 서버

설정 수준에서 `mock_enabled: true`가 설정되면 모든 collection에 `base_mock_url`이 설정되어야 합니다. 이는 swag2mcp에 이 collection에 대해 모의 서버가 실행 중인 위치를 알려줍니다.

```yaml
mock_enabled: true
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        base_mock_url: localhost:8080
```

자세한 내용은 [모의 서버](../advanced/mock-server.md)를 참조하세요.

## 예시

### 최소 Collection

```yaml
specs:
  - domain: dadjokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

### 모든 필드가 있는 전체 Collection

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        llm_instruction: "현재 날씨 및 일일 예보에 사용하세요."
        title: "Custom Title"
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        disable: false
        base_url: https://forecast-api.open-meteo.com
        base_mock_url: localhost:8080
        http_client:
          headers:
            X-Custom: value
```

### Spec당 여러 Collection

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

### 워크스페이스의 로컬 파일 (specs/ 디렉토리)

```yaml
specs:
  - domain: myapi
    llm_title: My Internal API
    base_url: https://api.mycompany.com
    collections:
      - llm_title: Users
        location: specs/users.openapi.json
      - llm_title: Orders
        location: specs/orders.openapi.json
```

### 비활성화된 Collection

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
        disable: true
```

## 관련 항목

- [Collection 설정 (config)](../configuration/collection-settings.md) — 전체 YAML 참조
- [설정 계단식](../configuration/cascade.md) — 설정이 서로를 재정의하는 방식
- [Specs](./specs) — collection의 논리적 컨테이너
- [HTTP Client](../configuration/http-client.md) — HTTP 클라이언트 설정
- [모의 서버](../advanced/mock-server.md) — 모의 서버 설정
- [CLI: validate](../cli/validate.md) — validate 명령어 참조
- [CLI: update](../cli/update.md) — update 명령어 참조
- [CLI: clean](../cli/clean.md) — clean 명령어 참조
