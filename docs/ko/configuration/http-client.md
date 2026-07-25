# HTTP 클라이언트

swag2mcp는 모든 API 호출에 설정 가능한 HTTP 클라이언트를 사용합니다. 이러한 설정은 전역적으로 정의되며 spec 및 collection 수준에서 재정의할 수 있습니다.

## 설정

```yaml
http_client:
  timeout: 30s
  max_response_size: 1048576
  user_agent: "swag2mcp-global/1.0"
  follow_redirects: true
  max_redirects: 10
  random: false
  proxy:
    url: ""
    username: ""
    password: ""
    bypass: []
  headers:
    "Accept": "application/json"
  cookies:
    - name: "session"
      value: "abc123"
      domain: ".example.com"
      path: "/"
```

## Timeout

swag2mcp가 API 응답을 기다리는 시간을 제어합니다.

- **타입:** duration (Go 형식: `30s`, `60s`, `2m`)
- **기본값:** `30s`
- **범위:** 1초 ~ 5분
- **효과:** API가 이 시간 내에 응답하지 않으면 요청이 timeout 오류로 실패합니다.
- **증가 시기:** 느린 API, 큰 페이로드, 불안정한 네트워크.
- **감소 시기:** 내부 API, 상태 확인, 빠른 실패 시나리오.

```yaml
http_client:
  timeout: 60s
```

## Max Response Size

응답이 LLM에 인라인으로 반환되지 않고 디스크에 저장되기 전의 최대 크기를 제한합니다.

- **타입:** `int` (바이트)
- **기본값:** `1048576` (1 MB)
- **범위:** 256 ~ 10,485,760 바이트 (10 MB)
- **효과:** 응답이 이 제한을 초과하면 `{workspace}/responses/`에 JSON 파일로 저장됩니다. LLM은 파일 참조를 받고 `response_outline`, `response_compress`, `response_slice` 도구로 탐색할 수 있습니다.
- **증가 시기:** 대용량 데이터 세트를 반환하는 API (보고서, 로그, 분석).
- **감소 시기:** 제한된 LLM 컨텍스트 윈도우 또는 모든 응답에 파일 기반 접근을 선호할 때.

```yaml
http_client:
  max_response_size: 4194304  # 4 MB
```

## User-Agent

모든 요청과 함께 전송되는 `User-Agent` 헤더입니다. 일부 API는 특정 사용자 에이전트를 요구하거나 알려진 봇 사용자 에이전트를 차단합니다.

- **타입:** `string`
- **기본값:** `"swag2mcp-global/1.0"`
- **효과:** API 서버에 애플리케이션을 식별합니다.
- **변경 시기:** API가 특정 사용자 에이전트를 요구하거나 분석을 위해 애플리케이션을 식별하려고 할 때.

```yaml
http_client:
  user_agent: "MyApp/1.0"
```

## Follow Redirects

swag2mcp가 HTTP 리디렉션(3xx 상태 코드)을 자동으로 따라갈지 제어합니다.

- **타입:** `bool`
- **기본값:** `true`
- **효과:** `true`일 때 swag2mcp가 `max_redirects` 횟수까지 리디렉션을 따릅니다. `false`일 때 리디렉션 응답이 그대로 반환됩니다.
- **비활성화 시기:** 루프에서 리디렉션하는 API, 리디렉션 대상을 수동으로 검사해야 하는 보안에 민감한 엔드포인트.

```yaml
http_client:
  follow_redirects: false
```

## Max Redirects

swag2mcp가 중단하기 전에 따르는 리디렉션 횟수를 제한합니다.

- **타입:** `int`
- **기본값:** `10`
- **범위:** 0 ~ 50
- **효과:** API가 이 제한보다 더 많이 리디렉션하면 요청이 실패합니다.
- **변경 시기:** 긴 리디렉션 체인이 있는 API, 또는 리디렉션 루프에서 더 빠른 실패를 위해 감소.

```yaml
http_client:
  max_redirects: 5
```

## Randomizer

각 요청에 브라우저와 유사한 무작위 헤더를 추가하여 핑거프린팅 및 차단을 방지합니다.

- **타입:** `bool`
- **기본값:** `false`
- **효과:** `true`일 때 swag2mcp가 각 요청에 대해 무작위 헤더를 생성합니다: `User-Agent` (실제 브라우저 문자열 풀에서), `Accept`, `Accept-Language`, `Accept-Encoding`, `Cache-Control`. 이는 `user_agent` 설정을 재정의합니다.
- **활성화 시기:** User-Agent 또는 헤더 패턴을 기반으로 요청을 차단하는 API, 스크래핑 시나리오.

```yaml
http_client:
  random: true
```

## Proxy

프록시 서버는 swag2mcp와 대상 API 간의 중개자 역할을 합니다. 모든 HTTP 트래픽이 이를 통해 라우팅됩니다.

**프록시가 필요할 수 있는 경우:**
- **회사 네트워크** — 모든 아웃바운드 트래픽이 회사 프록시를 통과해야 함
- **지역 제한** — 일부 API는 지역이 제한되어 있으며, 올바른 지역의 프록시가 이를 우회함
- **고정 IP** — IP 허용 목록이 필요한 API
- **익명성** — 대상 API로부터 원본 IP 숨기기

### Proxy URL

- **타입:** `string`
- **기본값:** `""` (프록시 없음)
- **지원 스키마:** `http`, `https`, `socks5`, `socks5h`
- **`$(VAR)` 지원:** ✅ 런타임에 해결

| 스키마 | 설명 | 사용 사례 |
|--------|------|----------|
| `http` | HTTP 트래픽용 HTTP 프록시 | 회사 프록시, 기본 프록시 |
| `https` | HTTPS 프록시 (CONNECT 터널) | 보안 회사 프록시 |
| `socks5` | SOCKS5 프록시 (DNS 로컬 해결) | 범용, 모든 프로토콜 |
| `socks5h` | SOCKS5 프록시 (DNS 프록시에서 해결) | 프록시가 더 나은 DNS 해상도를 가질 때 |

### Proxy 인증

프록시가 인증을 요구하는 경우 `username`과 `password`를 제공하세요:

- **`$(VAR)` 지원:** ✅ 세 필드 모두(`url`, `username`, `password`) 런타임에 해결

```yaml
http_client:
  proxy:
    url: "http://proxy.example.com:8080"
    username: "proxyuser"
    password: "$(PROXY_PASSWORD)"
```

### Proxy Bypass

프록시를 **통하지 않아야** 하는 도메인 목록입니다. 내부 서비스, localhost 또는 직접 접근만 가능한 API에 유용합니다.

```yaml
http_client:
  proxy:
    url: "http://proxy.example.com:8080"
    bypass:
      - "localhost"
      - "127.0.0.1"
      - "*.internal.company.com"
      - "api.local"
```

Bypass는 와일드카드 패턴을 지원합니다(`*.example.com`은 모든 하위 도메인과 일치).

## Headers

모든 요청에 추가되는 커스텀 HTTP 헤더입니다. 헤더는 계단식 수준 간에 병합됩니다:

```
전역 헤더 → Spec 헤더 (병합) → Collection 헤더 (병합)
```

Collection 헤더는 동일한 키에 대해 spec 헤더를 재정의하고, spec 헤더는 전역 헤더를 재정의합니다.

```yaml
http_client:
  headers:
    "Accept": "application/json"
    "Accept-Language": "en-US"
```

헤더 값은 `$(ENV_VAR)` 해결을 지원합니다.

## Cookies

모든 요청과 함께 전송되는 쿠키입니다. 쿠키는 계단식 수준 간에 병합됩니다(낮은 수준이 동일한 쿠키 이름에 대해 전역을 재정의).

```yaml
http_client:
  cookies:
    - name: "session"
      value: "abc123"
      domain: ".example.com"
      path: "/"
      secure: false
      http_only: false
```

### Cookie 필드

| 필드 | 필수 | 설명 |
|------|------|------|
| `name` | 예 | 쿠키 이름 |
| `value` | 예 | 쿠키 값 (`$(ENV_VAR)` 해결 지원) |
| `domain` | 아니요 | 도메인 범위 (예: `.example.com`) |
| `path` | 아니요 | 경로 범위 (예: `/`) |
| `secure` | 아니요 | HTTPS를 통해서만 전송 |
| `http_only` | 아니요 | JavaScript로 접근 불가 |

## Spec 수준의 커스텀 헤더

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    http_client:
      headers:
        "Accept": "application/json"
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Spec 수준의 쿠키

```yaml
specs:
  - domain: example
    llm_title: Example API
    base_url: https://api.example.com
    http_client:
      cookies:
        - name: "session"
          value: "abc123"
        - name: "csrf"
          value: "$(CSRF_TOKEN)"
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## 계단식

HTTP 클라이언트 설정은 전역에서 spec으로, spec에서 collection으로 계단식으로 적용됩니다. 모든 설정은 모든 수준에서 재정의할 수 있습니다:

```
전역 (http_client)
    ↓ 재정의 (모든 설정)
Spec (specs[].http_client)
    ↓ 재정의 (모든 설정)
Collection (specs[].collections[].http_client)
```

**모든 HTTP 클라이언트 설정**(timeout, proxy, user-agent, redirects, response size, randomizer, headers, cookies)은 spec 및 collection 수준에서 재정의할 수 있습니다.

자세한 내용은 [설정 계단식](./cascade)을 참조하세요.
