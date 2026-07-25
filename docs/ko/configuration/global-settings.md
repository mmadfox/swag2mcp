# 전역 설정

전역 설정은 `swag2mcp.yaml`의 최상위 설정 블록입니다. spec 또는 collection 수준에서 재정의되지 않는 한 모든 spec에 적용됩니다.

## 구조

```yaml
http_client:
  # 모든 API 호출의 HTTP 클라이언트 설정

mcp:
  # MCP 서버 설정

mock_enabled: false
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092

disable_ratelimiter: false
rate_limit_interval: 10s
```

## HTTP 클라이언트

swag2mcp가 API에 HTTP 요청을 하는 방식을 제어합니다: 타임아웃, 응답 크기 제한, 프록시, 헤더, 쿠키, 리디렉션, 사용자 에이전트. 이러한 설정은 spec과 collection으로 계단식으로 적용됩니다.

모든 매개변수와 예시는 [HTTP 클라이언트](./http-client)를 참조하세요.

## MCP 서버

MCP 서버가 LLM 에이전트와 통신하는 방식을 제어합니다: 전송 유형(stdio, SSE, Streamable HTTP), 주소, 경로, 선택적 bearer 토큰 인증.

모든 매개변수, 전송 방식, 시작 플래그는 [MCP 서버](./mcp-server)를 참조하세요.

## 모의 서버

모의 서버는 OpenAPI 스키마를 기반으로 가짜 API 응답을 생성합니다. 실제 API에 접근하지 않고 테스트할 때 유용합니다.

```yaml
mock_enabled: true
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092
```

### mock_enabled

- **타입:** `bool`
- **기본값:** `false`
- **효과:** `true`일 때 swag2mcp가 `base_mock_url`이 설정된 모든 spec에 대해 모의 서버를 시작합니다. 각 collection에 `base_mock_url`이 설정되어야 합니다.
- **활성화 시기:** 실제 HTTP 호출 없이 API 통합을 테스트하려고 할 때. 모의 서버는 OpenAPI 스키마를 기반으로 가짜 데이터를 반환합니다.

### mock_auth

모의 인증 서버의 포트 설정입니다. 모의 서버로 인증 방법(OAuth2, Digest, HMAC)을 테스트할 때 사용됩니다.

| 필드 | 타입 | 기본값 | 설명 |
|------|------|--------|------|
| `oauth2_port` | int | `9090` | 모의 OAuth2 토큰 서버 포트 (1024-65535) |
| `digest_port` | int | `9091` | 모의 Digest 인증 서버 포트 (1024-65535) |
| `hmac_port` | int | `9092` | 모의 HMAC 인증 서버 포트 (1024-65535) |

## 속도 제한기

속도 제한기는 LLM이 동일한 API 엔드포인트를 너무 자주 호출하는 것을 방지합니다. 기본적으로 각 엔드포인트는 10초에 한 번씩 호출할 수 있습니다.

```yaml
disable_ratelimiter: false
rate_limit_interval: 10s
```

### disable_ratelimiter

- **타입:** `bool`
- **기본값:** `false`
- **효과:** `true`일 때 엔드포인트별 속도 제한기가 완전히 비활성화됩니다. LLM이 대기 없이 동일한 엔드포인트를 반복해서 호출할 수 있습니다.
- **활성화 시기:** 테스트, 디버깅 또는 동일한 엔드포인트를 빠르게 연속해서 여러 번 호출해야 할 때.
- **비활성화 유지(권장):** 프로덕션. 속도 제한기는 실수로 인한 남용을 방지하고 API 속도 제한을 준수합니다.

### rate_limit_interval

- **타입:** duration (Go 형식: `10s`, `30s`, `1m`)
- **기본값:** `10s`
- **효과:** LLM이 동일한 엔드포인트 호출 간 대기해야 하는 시간을 설정합니다.
- **변경 시기:** 엄격한 속도 제한이 있는 API는 증가. 부하를 제어할 수 있는 내부 API는 감소.
- **범위:** 모든 유효한 duration (예: `5s`, `30s`, `1m`, `2m`).

## 계단식

전역 설정은 spec 및 collection 수준에서 재정의할 수 있습니다. 모든 `http_client` 설정(timeout, proxy, user-agent, redirects, response size, randomizer, headers, cookies)은 spec 및 collection 수준에서 재정의할 수 있습니다.

```
전역 (http_client, mock_enabled, disable_ratelimiter, rate_limit_interval)
    ↓ 재정의 (http_client만)
Spec (specs[].http_client)
    ↓ 재정의 (http_client만)
Collection (specs[].collections[].http_client)
```

자세한 내용은 [설정 계단식](./cascade)을 참조하세요.
