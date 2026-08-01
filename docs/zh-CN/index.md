# swag2mcp

通过模型上下文协议（MCP）将 OpenAPI/Swagger/Postman API 规范与 LLM 智能体连接起来。

并非所有 API 都支持 MCP — 私有端点、内部服务、遗留系统和第三方 API 很少支持。swag2mcp 将任何 REST API 封装到 MCP 接口中，让 LLM 智能体无需修改一行服务器代码即可即时访问您的整个 API 表面。通过实时 API 调用，LLM 获得真实世界的知识，从而做出明智的决策、自动化工作流程并对您的数据采取行动 — 而不仅仅是猜测。

<a href="https://www.youtube.com/watch?v=1Da4UmE2f9U" target="_blank">
  <img src="https://raw.githubusercontent.com/mmadfox/swag2mcp/main/docs/cover.jpg" alt="预览">
</a>

## 你的 API 会说 LLM 语言

一行配置即可将任何 OpenAPI/Swagger/Postman 文件转换为 MCP 服务器。LLM 智能体可以发现、检查和调用你的 API — 无需编写集成代码。

<img src="/architecture.svg" width="700" alt="swag2mcp 架构">

## 告别重复的包装代码

每次将新 API 连接到 LLM 时，你都要编写相同的样板代码：规范解析、认证、错误处理、速率限制。swag2mcp 为你完成这一切 — 19 个现成的 MCP 工具。

## 谁需要它

| 角色 | 原因 |
|------|------|
| **AI 智能体开发者** | 2 分钟连接任何 API，而不是 2 天 |
| **MCP 工程师** | 无需处理代码 — 只需指向规范即可 |
| **架构师** | 为公司所有 LLM 提供统一的 API 集成层 |
| **数据分析师** | 通过自然语言访问 API，无需编码 |
| **DevOps / SRE** | 通过 LLM 进行监控和自动化，无需额外服务 |
| **集成工程师** | 9 种开箱即用的认证方法 — 从 Basic 到 OAuth2 到 HMAC |
| **QA 工程师** | 无需真实 API 即可进行隔离测试的模拟服务器 |
| **产品经理** | 无需后端工作即可快速构建 AI 功能原型 |
| **以及其他许多人** | |

---

## 许可证

基于 **Apache License, Version 2.0** 许可。

完整许可文本请参见 [LICENSE](https://github.com/mmadfox/swag2mcp/blob/main/LICENSE)。
