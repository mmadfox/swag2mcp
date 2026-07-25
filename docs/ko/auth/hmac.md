# HMAC Auth

## 목적

HMAC-SHA256 요청 서명 — 암호화폐 거래소(Binance, Bybit 등)에서 사용하는 인증 방식입니다. 각 요청이 비밀 키로 서명됩니다.

## 사용 시기

- Binance API 및 Binance 호환 거래소
- 암호화폐 거래 플랫폼
- 요청 서명이 필요한 API

## 설정

```yaml
specs:
  - domain: binance
    llm_title: Binance Market Data
    base_url: https://api.binance.com
    collections:
      - llm_title: Market Data
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/binance.yaml
    auth:
      type: hmac
      config:
        api_key: "$(BINANCE_API_KEY)"
        secret_key: "$(BINANCE_SECRET_KEY)"
```

## 매개변수

| 매개변수 | 필수 | 설명 |
|---------|------|------|
| `api_key` | 예 | 공개 API 키 |
| `secret_key` | 예 | 서명용 비밀 키 |

## 참고 사항

- swag2mcp는 모든 요청에 자동으로 타임스탬프(Unix 밀리초)를 추가합니다
- 서명은 모든 요청 매개변수에서 계산됩니다
- 키를 환경 변수에 저장하세요: `api_key: "$(BINANCE_API_KEY)"`
- 이 방식은 Binance API 및 유사한 거래소와 호환됩니다
