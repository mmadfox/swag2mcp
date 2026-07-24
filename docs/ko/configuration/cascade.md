# 설정 계단식

swag2mcp는 3단계 설정 계단식을 사용합니다. 각 수준이 이전 수준을 재정의합니다. 이를 통해 전역적으로 합리적인 기본값을 설정하고 특정 spec 또는 collection에 대해 세부 조정할 수 있습니다.

## 수준

```
전역 (http_client, mcp, mock_enabled, disable_ratelimiter, rate_limit_interval)
    ↓ 재정의
Spec (specs[].http_client, specs[].auth, specs[].base_url, specs[].disable, specs[].tags)
    ↓ 재정의
Collection (specs[].collections[].http_client, specs[].collections[].base_url, specs[].collections[].disable)
```

## 재정의되는 항목

| 매개변수 | 전역 | Spec | Collection |
|---------|------|------|------------|
| `http_client.timeout` | ✅ | ✅ | ✅ |
| `http_client.max_response_size` | ✅ | ✅ | ✅ |
| `http_client.user_agent` | ✅ | ✅ | ✅ |
| `http_client.follow_redirects` | ✅ | ✅ | ✅ |
| `http_client.max_redirects` | ✅ | ✅ | ✅ |
| `http_client.proxy` | ✅ | ✅ | ✅ |
| `http_client.random` | ✅ | ✅ | ✅ |
| `http_client.headers` | ✅ | ✅ | ✅ |
| `http_client.cookies` | ✅ | ✅ | ✅ |
| `base_url` | ❌ | ✅ | ✅ |
| `auth` | ❌ | ✅ | ❌ |
| `disable` | ❌ | ✅ | ✅ |
| `tags` | ❌ | ✅ | ❌ |
| `mock_enabled` | ✅ | ❌ | ❌ |
| `disable_ratelimiter` | ✅ | ❌ | ❌ |
| `rate_limit_interval` | ✅ | ❌ | ❌ |

모든 `http_client` 설정은 모든 수준에서 재정의할 수 있습니다. Collection 수준 설정은 spec 및 전역보다 완전히 우선합니다.

## 계단식 예시

```yaml
http_client:
  timeout: 30s
  max_response_size: 1048576
  headers:
    "User-Agent": "swag2mcp/1.0"

specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    http_client:
      timeout: 60s  # 전역 timeout 재정의
      headers:
        "X-API-Version": "2"  # 전역 헤더에 추가
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          timeout: 120s  # spec timeout 재정의
          headers:
            "X-Custom": "value"  # spec + 전역 헤더에 추가
```

## "Forecast" Collection의 유효 설정

```
timeout: 120s (collection에서, spec 60s 및 전역 30s 재정의)
max_response_size: 1048576 (전역에서)
headers:
  - User-Agent: swag2mcp/1.0 (전역에서)
  - X-API-Version: 2 (spec에서)
  - X-Custom: value (collection에서)
```

## 병합 방식

### HTTP 클라이언트 설정

단순 값(`timeout`, `max_response_size`, `user_agent`, `follow_redirects`, `max_redirects`, `random`)은 각 수준에서 **대체**됩니다. spec이 `timeout: 60s`로 설정하면 전역 `30s`를 완전히 대체합니다.

### 헤더

헤더는 수준 간에 **병합**됩니다. 세 수준 모두의 헤더가 결합됩니다. 동일한 헤더 키가 여러 수준에 나타나면 가장 낮은 수준이 우선합니다.

### 쿠키

쿠키는 수준 간에 **병합**됩니다. 동일한 쿠키 이름이 여러 수준에 나타나면 가장 낮은 수준이 우선합니다.

### 프록시

프록시는 각 수준에서 **대체**됩니다. spec이 프록시를 설정하면 해당 spec의 전역 프록시를 완전히 대체합니다.
