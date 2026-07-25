# 配置级联

swag2mcp 使用三级配置级联。每个级别覆盖前一个级别。这让你可以在全局设置合理的默认值，并为特定 spec 或 collection 微调设置。

## 级别

```
全局 (http_client, mcp, mock_enabled, disable_ratelimiter, rate_limit_interval)
    ↓ 覆盖
Spec (specs[].http_client, specs[].auth, specs[].base_url, specs[].disable, specs[].tags)
    ↓ 覆盖
Collection (specs[].collections[].http_client, specs[].collections[].base_url, specs[].collections[].disable)
```

## 覆盖关系

| 参数 | 全局 | Spec | Collection |
|------|------|------|------------|
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

所有 `http_client` 设置可以在每个级别覆盖。Collection 级别的设置完全优先于 spec 和全局。

## 级联示例

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
      timeout: 60s  # 覆盖全局超时
      headers:
        "X-API-Version": "2"  # 添加到全局头
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          timeout: 120s  # 覆盖 spec 超时
          headers:
            "X-Custom": "value"  # 添加到 spec + 全局头
```

## "Forecast" Collection 的有效设置

```
timeout: 120s（来自 collection，覆盖 spec 的 60s 和全局的 30s）
max_response_size: 1048576（来自全局）
headers:
  - User-Agent: swag2mcp/1.0（来自全局）
  - X-API-Version: 2（来自 spec）
  - X-Custom: value（来自 collection）
```

## 合并方式

### HTTP 客户端设置

简单值（`timeout`、`max_response_size`、`user_agent`、`follow_redirects`、`max_redirects`、`random`）在每个级别**替换**。如果 spec 设置 `timeout: 60s`，它完全替换全局的 `30s`。

### 头

头在级别之间**合并**。所有三个级别的头被组合。如果相同的头键出现在多个级别，最低级别获胜。

### Cookie

Cookie 在级别之间**合并**。如果相同的 cookie 名称出现在多个级别，最低级别获胜。

### 代理

代理在每个级别**替换**。如果 spec 设置了代理，它完全替换该 spec 的全局代理。
