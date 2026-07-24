# CLI 工作流程

此页面展示了从初始化到日常操作的使用 swag2mcp 的真实终端示例。

## 快速开始

```bash
# 1. 初始化工作区
mkdir -p .swag2mcp && swag2mcp init ./.swag2mcp

# 2. 列出你的 spec
swag2mcp ls
```

## 使用 YAML 添加 spec

### 简单 spec（公共 API）

```bash
swag2mcp add spec --yaml - <<EOF
domain: meteo
llm_title: Open-Meteo Weather API
base_url: https://api.open-meteo.com
collections:
  - llm_title: Weather Forecast
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
EOF
```

### 带认证的 spec（来自环境变量的 bearer 令牌）

```bash
swag2mcp add spec --yaml - <<EOF
domain: my-api
llm_title: My Protected API
base_url: https://api.example.com/v1
auth:
  type: bearer
  config:
    token: \$(MY_TOKEN)
collections:
  - llm_title: Users
    location: https://raw.githubusercontent.com/my-org/my-api/main/users.yaml
EOF
```

### 带多个 collection 的 spec

```bash
swag2mcp add spec --yaml - <<EOF
domain: meteo
llm_title: Open-Meteo APIs
base_url: https://api.open-meteo.com
collections:
  - llm_title: Forecast
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
  - llm_title: Air Quality
    base_url: https://air-quality-api.open-meteo.com
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
EOF
```

## 向现有 spec 添加 collection

```bash
swag2mcp add collection --yaml - <<EOF
spec_domain: meteo
llm_title: Marine Weather
location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
EOF
```

## 列出 spec

```bash
$ swag2mcp ls
Specifications:
  dadjoke (https://icanhazdadjoke.com)
    jokes (3 endpoints)
  meteo (https://api.open-meteo.com)
    forecast (5 endpoints)
    air-quality (8 endpoints)
    marine (4 endpoints)
```

### 按标签过滤

```bash
swag2mcp ls --tags=public
```

## 查看运行时信息

```bash
$ swag2mcp info
{
  "version": "v1.2.0",
  "workspace": "/home/user/.swag2mcp",
  "specs": {
    "total": 2,
    "active": 2,
    "disabled": 0,
    "collections": 4,
    "endpoints": 20
  },
  "http_client": {
    "timeout": "30s",
    "max_response_size": "2 KB",
    "follow_redirects": true
  },
  "mcp": {
    "transport": "stdio"
  },
  "auth": {
    "methods": ["bearer"]
  }
}
```

## 验证配置

```bash
$ swag2mcp validate
✅ Configuration is valid.
✓ Spec dadjoke: OK
✓ Spec meteo: OK
```

## 启动 MCP 服务器

### stdio（用于 IDE 集成）

```bash
swag2mcp mcp
```

### HTTP（用于远程访问）

```bash
swag2mcp mcp --transport sse --http-addr :8080
```

### 带标签过滤

```bash
swag2mcp mcp --tags=public
```

## 更新规范

刷新所有缓存的规范文件：

```bash
swag2mcp update
```

## 清理缓存

```bash
swag2mcp clean
```

## 导出和导入

### 备份你的工作区

```bash
swag2mcp export --output ~/backups/swag2mcp-2026-07-24.zip
```

### 在另一台机器上恢复

```bash
# 在新机器上
swag2mcp import --from-zip swag2mcp-2026-07-24.zip
```

## 交互式 TUI 浏览器

```bash
swag2mcp run
```

打开一个全屏终端 UI，用于搜索、浏览和调用 API。

## 模拟服务器

```bash
# 安装模拟二进制文件
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest

# 启动模拟服务器
swag2mcp-mock mockserver
```
