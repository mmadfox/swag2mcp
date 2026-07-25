# None

## 用途

无需认证。API 无需令牌或密钥即可访问。

## 何时使用

- 公共 API（Open-Meteo、icanhazdadjoke、PokéAPI）
- 测试和演示环境
- 当 API 不需要授权时

## 配置

设置 `type: none` 或直接省略 `auth` 部分：

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: none
```

## 参数

无。

## 说明

- 如果配置中完全缺少 `auth` 部分，等同于 `type: none`
- 不会向请求添加任何授权头
