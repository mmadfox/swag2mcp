# HMAC Auth

## 用途

HMAC-SHA256 请求签名 — 加密货币交易所（Binance、Bybit 等）使用的认证方法。每个请求使用密钥签名。

## 何时使用

- Binance API 和兼容 Binance 的交易所
- 加密货币交易平台
- 需要请求签名的 API

## 配置

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

## 参数

| 参数 | 必需 | 描述 |
|------|------|------|
| `api_key` | 是 | 公共 API 密钥 |
| `secret_key` | 是 | 用于签名的密钥 |

## 说明

- swag2mcp 自动为每个请求添加时间戳（Unix 毫秒）
- 签名从所有请求参数计算得出
- 将密钥存储在环境变量中：`api_key: "$(BINANCE_API_KEY)"`
- 此方法兼容 Binance API 和类似交易所
