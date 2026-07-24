# version

## 用途

打印 swag2mcp 版本。用于验证安装的版本、报告错误或检查兼容性。

## 何时使用

- 你想检查安装的 swag2mcp 版本
- 你正在报告错误，需要包含版本信息
- 你想验证安装是否成功

## 语法

```bash
swag2mcp version
swag2mcp --version
```

## 参数

无。

## 标志

无。

## 工作原理

```bash
swag2mcp version
# swag2mcp v1.2.0

swag2mcp --version
# swag2mcp v1.2.0
```

## 输出格式

```
swag2mcp <version>
```

版本在构建时通过 `ldflags` 设置。如果未设置，默认为 `"dev"`。

## 细节

- **两种形式：** `swag2mcp version`（子命令）和 `swag2mcp --version`（全局标志）产生相同的输出。
- **无需配置：** 此命令无需工作区或配置文件即可工作。
