# Specs

Spec은 API 도메인 또는 서비스를 나타내는 논리적 컨테이너입니다(예: YouTube, Binance, Open-Meteo). 각 spec은 고유한 `domain`, `base_url`, 선택적 `auth`를 가지며 하나 이상의 collection을 포함합니다.

[Collections](./collections)는 OpenAPI/Swagger/Postman 파일을 가리킵니다 — spec 자체는 파일이 아니라 이를 둘러싼 그룹입니다.

## Domain — 이름 규칙

`domain`은 spec의 고유 식별자입니다. 시스템 전체에서 기본 키로 사용됩니다.

| 규칙 | 제약 |
|------|------|
| 문자 | `a-z`, `0-9`, `_`, `-`만 허용 |
| 길이 | 1–60자 |
| 고유성 | **중복 불가** — 두 개의 활성 spec이 동일한 domain을 공유할 수 없음 |

**유효한 예:** `meteo`, `binance`, `github-api`, `my_service`, `openai-v1`

**유효하지 않은 예:** `Meteo`(대문자), `my api`(공백), `my.api`(점), `a-very-long-domain-name-that-exceeds-sixty-characters`(너무 김)

## Spec 필드

| 필드 | YAML 키 | 필수 | 설명 |
|------|---------|------|------|
| [Domain](#domain--naming-rules) | `domain` | ✅ | 고유 API 식별자 (1–60자, `a-z0-9_-`) |
| LLM Title | `llm_title` | ✅ | LLM이 이 API를 참조할 때 사용하는 사람이 읽을 수 있는 이름 (5–120자) |
| [LLM Instruction](#llm-instruction) | `llm_instruction` | ❌ | swag2mcp 시스템 프롬프트에 주입되는 짧은 힌트 (최대 500자) |
| Base URL | `base_url` | ✅ | 모든 API 요청의 기본 URL (유효한 URL) |
| [Disable](#disable) | `disable` | ❌ | 로딩 및 인덱싱 중 이 spec 건너뛰기 |
| [Tags](#tags) | `tags` | ❌ | 필터링용 태그 (예: `["public", "demo"]`) |
| [Auth](#auth) | `auth` | ❌ | 인증 설정 |
| [HTTP Client](#http-client) | `http_client` | ❌ | Spec별 HTTP 설정 (헤더, 쿠키) |
| [Collections](./collections) | `collections` | ✅ | 1–30개 collection 목록 |

## 검증

swag2mcp가 설정을 검증할 때 모든 spec에 대해 다음 규칙이 확인됩니다:

| 확인 | 규칙 |
|------|------|
| **중복 도메인** | 두 개의 활성 spec이 동일한 `domain`을 공유할 수 없음 |
| **도메인 형식** | `^[a-z0-9_-]{1,60}$`와 일치해야 함 |
| **LLM Title** | 필수, 5–120자, 문자/숫자/공백/기본 구두점 |
| **LLM Instruction** | 최대 500자, title과 동일한 문자 세트 |
| **Base URL** | 필수, 유효한 URL이어야 함 |
| **Collections** | 필수, 1–30개 항목 |
| **Auth** | 인증 유형별로 검증 (예: bearer는 `token` 필요, basic은 `username` + `password` 필요) |
| **Location** | 각 collection의 `location`은 유효한 URL 또는 파일 경로여야 함 (5–250자) |

검증은 모든 `swag2mcp mcp` 시작 시 실행됩니다. 실패하면 MCP 서버가 시작되지 않습니다 — 일부 IDE에서는 서버가 연결되지 않고 LLM이 수정해야 할 사항을 설명하는 명확한 오류 메시지를 받게 됩니다.

서버를 시작하기 전에 문제를 진단하려면 [`validate`](../cli/validate.md) 명령어를 사용하세요:

```bash
# 기본 워크스페이스 검증 (~/.swag2mcp)
swag2mcp validate

# 커스텀 프로젝트 워크스페이스 검증
swag2mcp validate ./my-project
```

## LLM Instruction

각 spec에 `llm_instruction`을 설정하는 것이 좋습니다 — 이 API의 용도와 사용 시기를 LLM에 알려주는 짧은 힌트(최대 500자)입니다. 이 지침은 swag2mcp 시스템 프롬프트에 주입되어 LLM이 추가 컨텍스트 없이 spec의 목적을 이해하는 데 도움을 줍니다.

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    llm_instruction: "이 API를 사용하여 무작위 아재개그를 얻거나 키워드로 특정 농담을 검색하세요."
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

Collection은 더 구체적인 지침을 위해 자체 `llm_instruction`(최대 360자)를 가질 수도 있습니다.

## Auth

인증은 spec 수준에서 설정되며 모든 collection에 적용됩니다. swag2mcp는 9가지 인증 방법을 지원합니다:

| 방법 | YAML 타입 | 주요 필드 |
|------|-----------|----------|
| [None](../auth/none.md) | `none` | — |
| [Basic](../auth/basic.md) | `basic` | `username`, `password` |
| [Bearer](../auth/bearer.md) | `bearer` | `token` |
| [Digest](../auth/digest.md) | `digest` | `username`, `password` |
| [OAuth2 Client Credentials](../auth/oauth2-cc.md) | `oauth2-cc` | `client_id`, `client_secret`, `token_url` |
| [OAuth2 Password](../auth/oauth2-pwd.md) | `oauth2-pwd` | `username`, `password`, `client_id`, `token_url` |
| [API Key](../auth/api-key.md) | `api-key` | `key`, `value`, `in` (`header` 또는 `query`) |
| [HMAC](../auth/hmac.md) | `hmac` | `api_key`, `secret_key` |
| [Script](../auth/script.md) | `script` | `domain` |

각 방법에 대한 자세한 내용은 [인증 개요](../auth/overview.md)를 참조하세요.

## HTTP Client

spec 수준에서 HTTP 설정을 재정의할 수 있습니다. 이는 이 spec의 collection이 만드는 모든 요청에 적용됩니다.

```yaml
specs:
  - domain: slow-api
    llm_title: Slow API
    base_url: https://slow-api.example.com
    http_client:
      headers:
        X-API-Version: "2"
      cookies:
        - name: session
          value: abc123
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

설정은 전역 → spec → collection 순으로 계단식으로 적용됩니다. 자세한 내용은 [설정 계단식](../configuration/cascade.md)을 참조하세요.

## Tags

태그를 사용하면 카테고리별로 spec을 필터링할 수 있습니다. `swag2mcp ls` 또는 부트스트랩 중에 `--tags` 플래그와 함께 사용하세요.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    tags: ["weather", "public"]
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

```bash
# "weather" 태그가 있는 spec만 나열
swag2mcp ls --tags weather
```

## Disable

`disable: true`로 설정하면 spec을 완전히 건너뜁니다. 로드, 인덱싱되지 않으며 LLM이 사용할 수 없습니다.

```yaml
specs:
  - domain: old-api
    llm_title: Old API (Deprecated)
    base_url: https://old-api.example.com
    disable: true
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## 예시

### 최소 Spec

```yaml
specs:
  - domain: dadjokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

### 인증이 있는 Spec

```yaml
specs:
  - domain: binance
    llm_title: Binance Market Data API
    base_url: https://api.binance.com
    auth:
      type: hmac
      config:
        api_key: $(BINANCE_API_KEY)
        secret_key: $(BINANCE_SECRET_KEY)
    collections:
      - llm_title: Market Data
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/binance.yaml
```

### 여러 Collection이 있는 Spec

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
      - llm_title: Marine
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

### LLM Instruction과 Tags가 있는 Spec

```yaml
specs:
  - domain: rickandmorty
    llm_title: Rick and Morty API
    llm_instruction: "이 API를 사용하여 Rick and Morty 쇼의 캐릭터, 에피소드, 위치 정보를 가져오세요."
    base_url: https://rickandmortyapi.com/api
    tags: ["entertainment", "public"]
    collections:
      - llm_title: Characters
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/rick-and-morty.json
```

## 관련 항목

- [Spec 설정 (config)](../configuration/spec-settings.md) — 전체 YAML 참조
- [설정 계단식](../configuration/cascade.md) — 설정이 서로를 재정의하는 방식
- [인증 개요](../auth/overview.md) — 모든 9가지 인증 방법
- [HTTP Client](../configuration/http-client.md) — HTTP 클라이언트 설정
- [Collections](./collections) — spec 내의 명세 파일
