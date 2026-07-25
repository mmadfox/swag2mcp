# 개념

## 아키텍처

swag2mcp는 API 명세와 LLM 에이전트 간의 브리지 역할을 합니다:

<img src="/architecture.svg" width="800" alt="swag2mcp architecture">

## 핵심 개념

**Spec** — API 도메인 또는 서비스를 나타내는 논리적 컨테이너(예: YouTube, Binance, Open-Meteo). 각 spec은 고유한 `domain`, `base_url`, 선택적 `auth`를 가지며 하나 이상의 collection을 포함합니다. 또한 `llm_instruction`을 설정할 수 있습니다 — swag2mcp 시스템 프롬프트에 주입되어 LLM에 이 spec의 용도와 사용 시기를 알려주는 짧은 힌트입니다. 자세히 알아보기: [Specs](./specs).

**Collection** — 특정 API를 설명하는 단일 OpenAPI/Swagger/Postman 파일입니다. `location`(URL 또는 로컬 파일 경로)을 가리킵니다. 하나의 spec은 여러 collection을 가질 수 있습니다 — 예를 들어, "meteo" spec에는 각각 다른 명세 파일을 가리키는 "Forecast", "Air Quality", "Marine" collection이 있을 수 있습니다. 자세히 알아보기: [Collections](./collections).

**Tag** — collection 내 엔드포인트의 카테고리입니다. LLM이 올바른 작업을 더 정확하게 찾는 데 도움을 줍니다. 자세히 알아보기: [Tags](./tags).

**Endpoint** — 특정 HTTP 메서드 + 경로(예: `GET /api/users`). LLM은 설명으로 엔드포인트를 찾고, 매개변수와 스키마를 검사한 후 호출할 수 있습니다. 자세히 알아보기: [Endpoints](./endpoints).

**Workspace** — swag2mcp가 설정, 명세 캐시, 저장된 응답, 인증 스크립트를 저장하는 디렉토리입니다. 자세히 알아보기: [Workspace](./workspace).

## 작동 방식

1. **spec 또는 collection 추가** — YAML 설정(`~/.swag2mcp/swag2mcp.yaml`)에 정의합니다. 예:

   ```yaml
   specs:
     - domain: jokes
       llm_title: Dad Joke API
       base_url: https://icanhazdadjoke.com
       collections:
         - llm_title: Jokes
           location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
   ```
2. **swag2mcp가 각 collection을 파싱** — Tag와 Endpoint를 생성하고 검색을 위해 인덱싱합니다.
3. **LLM이 올바른 엔드포인트를 찾음** — MCP 도구(`search`, `endpoint_by_tag`, `inspect`)를 통해 LLM이 설명과 일치하는 엔드포인트를 검색하고 매개변수와 요청 스키마를 검토합니다.
4. **LLM이 엔드포인트를 호출** — MCP 도구 `invoke`를 통해 LLM이 요청을 보냅니다. swag2mcp는 호출 전에 엔드포인트의 OpenAPI 스키마에 대해 모든 입력 매개변수(path params, query params, headers, request body)를 검증합니다. 스키마와 일치하지 않으면 LLM은 무엇이 잘못되었는지 설명하는 명확한 오류를 받습니다. 검증이 완료되면 swag2mcp가 실제 HTTP 호출을 실행하고 결과를 반환합니다.
5. **결과가 LLM으로 돌아감** — API 응답이 에이전트에 전달됩니다. 큰 응답은 워크스페이스에 저장되며 세 가지 전용 MCP 도구로 탐색할 수 있습니다: `response_outline`(구조 보기), `response_compress`(대표 샘플로 축소), `response_slice`(특정 조각 추출).

swag2mcp는 LLM과 API 세계 간의 브리지입니다. API 명세를 추가하면 LLM이 MCP 프로토콜을 통해 올바른 엔드포인트를 찾고, 문서를 검사하고, 호출합니다. 여러분이 해야 할 일은 spec을 추가하고 MCP 서버를 시작하는 것뿐입니다.

> **설정은 언제든지 편집 가능합니다.** YAML 설정 파일(`~/.swag2mcp/swag2mcp.yaml`)은 수동으로 편집할 수 있습니다 — spec 추가, 인증 변경, 설정 조정. 편집 후에는 변경 사항을 적용하려면 MCP 서버를 다시 시작하세요(`swag2mcp mcp`).

## 계층 구조

```
Spec (domain, e.g. "meteo")
  └── Collection 1 (spec file, e.g. forecast.yml)
        └── Tag 1 (category)
              └── Endpoint (GET /api/forecast)
              └── Endpoint (POST /api/forecast)
        └── Tag 2
              └── Endpoint (GET /api/forecast/{id})
  └── Collection 2 (spec file, e.g. air-quality.yml)
        └── Tag 3
              └── Endpoint (GET /api/air-quality)
```
