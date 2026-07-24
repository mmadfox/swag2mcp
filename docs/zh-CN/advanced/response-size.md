# 响应大小管理

## 概述

API 响应可能非常大 — 有时大到无法放入 LLM 的上下文窗口。swag2mcp 通过将过大的响应保存到磁盘并提供工具来探索它们，自动管理响应大小。

## 工作原理

1. **你调用 `invoke`** — swag2mcp 发起 API 请求
2. **如果响应很小**（在限制内）— 直接内联返回给 LLM
3. **如果响应太大**（超过限制）— 保存到 `{workspace}/responses/` 作为 JSON 文件。LLM 收到文件引用而不是完整响应

### 示例：小响应（内联）

```json
{
  "statusCode": 200,
  "body": {
    "id": 1,
    "name": "Rex",
    "status": "available"
  }
}
```

### 示例：大响应（文件引用）

```json
{
  "statusCode": 200,
  "fileRef": {
    "path": "/Users/user/.swag2mcp/responses/response_a1b2c3d4.json",
    "size": 1572864,
    "sizeHint": "1.5 MB",
    "maxSizeHint": "2 KB",
    "message": "Response exceeds the 2 KB limit and has been saved to disk.",
    "openCmd": "open /Users/user/.swag2mcp/responses/response_a1b2c3d4.json"
  }
}
```

## 配置

```yaml
http_client:
  max_response_size: 1048576  # 1 MB（字节）
```

### max_response_size

- **类型：** `int`（字节）
- **默认值：** `1048576`（1 MB）
- **范围：** 256 到 10,485,760 字节（10 MB）
- **效果：** 超过此大小的响应将保存到磁盘，而不是内联返回
- **何时增加：** 返回大数据集的 API（报告、日志、分析）
- **何时减少：** LLM 上下文窗口有限，或你更倾向于基于文件的访问

## 处理大响应

当 `invoke` 返回 `fileRef` 时，使用这三个工具探索数据：

### 1. response_outline — 了解结构

获取响应的结构摘要：键、类型、数组长度和导航提示。

```json
→ response_outline(path: "/path/to/file.json")
← {
    "type": "object",
    "size": 1572864,
    "keys": ["data", "meta"],
    "itemCount": 500,
    "compressionHints": [
      "response_compress(path, 'first_of_array', 'data')",
      "response_compress(path, 'sample_array', 'data', arrayHead=3, arrayTail=2)"
    ]
  }
```

### 2. response_compress — 获取更小的版本

压缩数据以适应内联。多种压缩模式让你选择合适的权衡。

| 模式 | 描述 | 最适合 |
|------|------|--------|
| `first_of_array` | 只保留数组的第一个元素 | 所有元素结构相同时 |
| `sample_array` | 保留数组的头部（3）和尾部（2） | 需要查看值的范围时 |
| `truncate_strings` | 将每个字符串缩短到 N 个字符 | 字符串非常长时 |
| `keys_only` | 将值替换为类型名称 | 只需要结构时 |
| `select_keys` | 只保留指定的键 | 需要特定字段时 |

```json
→ response_compress(path: "/path/to/file.json", mode: "first_of_array", jsonPath: "data")
← {
    "body": [{ "id": 1, "name": "Rex" }],
    "hint": "Compressed array from 500 to 1 item using first_of_array mode"
  }
```

### 3. response_slice — 提取特定片段

通过 JSON 路径或行范围获取特定元素或值。

```json
→ response_slice(path: "/path/to/file.json", jsonPath: "data.0")
← {
    "slice": {
      "value": { "id": 1, "name": "Rex" },
      "jsonPath": "data.0",
      "nextPath": "data.1",
      "prevPath": null
    }
  }
```

## 完整工作流程

```
1. invoke(endpoint) → fileRef（响应为 1.5 MB）
2. response_outline(path) → 结构：{ data: Array(500) }
3. response_compress(path, mode: "first_of_array", jsonPath: "data") → 第一个项目
4. response_slice(path, jsonPath: "data.0") → 完整第一个项目详情
5. response_slice(path, jsonPath: "data.1") → 第二个项目
```

## 自动清理

当 MCP 服务器启动（`swag2mcp mcp`）时，超过 48 小时的响应文件会自动删除。你也可以手动清理：

```bash
swag2mcp clean
```

## 重要说明

- **限制以字节为单位** — `1048576` = 1 MB，`2097152` = 2 MB 等
- **文件引用包含打开命令** — macOS 上为 `open`，Linux 上为 `xdg-open`
- **响应文件使用随机后缀命名** — 并发调用之间不会冲突
- **响应目录自动创建** — 无需手动设置
