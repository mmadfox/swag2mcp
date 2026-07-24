# 安装

## 要求

- **macOS、Linux 或 Windows**（amd64 / arm64）
- **Go 1.26+**（仅用于 `go install` 或从源码构建）

## 兼容性

| 方法 | macOS | Linux | Windows |
|------|-------|-------|---------|
| 一行命令（curl） | ✅ | ✅ | ❌ |
| Homebrew | ✅ | ✅ | ❌ |
| Scoop | ❌ | ❌ | ✅ |
| APT（deb） | ❌ | ✅ | ❌ |
| RPM | ❌ | ✅ | ❌ |
| Docker | ❌ | ✅ | ❌ |
| GitHub Release | ✅ | ✅ | ✅ |
| go install | ✅ | ✅ | ✅ |
| 从源码构建 | ✅ | ✅ | ✅ |

---

## macOS

### 一行命令（推荐）

```bash
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash
```

安装到 `/usr/local/bin/swag2mcp`（如果 `/usr/local/bin` 不可写，则安装到 `~/.local/bin/swag2mcp`）。

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp
```

### GitHub Release

::: code-group

```bash [Apple Silicon]
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_darwin_arm64.tar.gz
tar -xzf swag2mcp_darwin_arm64.tar.gz
sudo mv swag2mcp /usr/local/bin/
```

```bash [Intel]
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_darwin_amd64.tar.gz
tar -xzf swag2mcp_darwin_amd64.tar.gz
sudo mv swag2mcp /usr/local/bin/
```

:::

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

确保 `$GOPATH/bin` 在你的 `$PATH` 中：

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### 从源码构建

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

---

## Linux

### 一行命令（推荐）

```bash
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash
```

安装到 `/usr/local/bin/swag2mcp`（如果 `/usr/local/bin` 不可写，则安装到 `~/.local/bin/swag2mcp`）。

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp
```

### APT（Debian / Ubuntu）

```bash
# 从最新版本下载 .deb
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_linux_amd64.deb
sudo dpkg -i swag2mcp_linux_amd64.deb
```

### RPM（Fedora / RHEL）

```bash
# 从最新版本下载 .rpm
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_linux_amd64.rpm
sudo rpm -i swag2mcp_linux_amd64.rpm
```

### Docker

```bash
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

使用 stdio 传输运行：

```bash
docker run --rm -i ghcr.io/mmadfox/swag2mcp:latest swag2mcp mcp
```

使用 HTTP 传输运行：

```bash
docker run --rm -p 8080:8080 ghcr.io/mmadfox/swag2mcp:latest swag2mcp mcp --transport sse --http-addr :8080
```

### GitHub Release

::: code-group

```bash [amd64]
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_linux_amd64.tar.gz
tar -xzf swag2mcp_linux_amd64.tar.gz
sudo mv swag2mcp /usr/local/bin/
```

```bash [arm64]
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_linux_arm64.tar.gz
tar -xzf swag2mcp_linux_arm64.tar.gz
sudo mv swag2mcp /usr/local/bin/
```

:::

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

确保 `$GOPATH/bin` 在你的 `$PATH` 中：

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### 从源码构建

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

---

## Windows

### Scoop

```powershell
scoop bucket add mmadfox https://github.com/mmadfox/scoop-bucket
scoop install mmadfox/swag2mcp
```

### GitHub Release

```powershell
# 下载最新版本
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_windows_amd64.zip
Expand-Archive swag2mcp_windows_amd64.zip -DestinationPath .
move swag2mcp.exe C:\Windows\System32\
```

### go install

```powershell
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

### 从源码构建

```powershell
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp.exe ./cmd/swag2mcp
```

---

## 模拟服务器

`swag2mcp-mock` 二进制文件作为单独的下载提供。使用与主二进制文件相同的方法安装：

::: code-group

```bash [macOS / Linux]
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

```powershell [Windows]
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

:::

或从 [GitHub Releases](https://github.com/mmadfox/swag2mcp/releases) 下载 — 查找 `swag2mcp-mock_&lt;version&gt;_&lt;os&gt;_&lt;arch&gt;.tar.gz`。

---

## 通过 LLM 智能体安装

如果你使用 AI 驱动的 IDE（OpenCode、Cursor、Claude Desktop、VS Code 等），你可以通过智能体安装 swag2mcp：

1. 让智能体添加 swag2mcp 技能：

   ```
   "Create the .agents/skills/swag2mcp-cli directory and add the skill from https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md to .agents/skills/swag2mcp-cli/SKILL.md"
   "Create the .agents/skills/swag2mcp-format directory and add the skill from https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md to .agents/skills/swag2mcp-format/SKILL.md"
   ```

2. 然后告诉你的智能体：

   ```
   "Set up swag2mcp"
   ```

   智能体将下载并安装 swag2mcp，然后创建带有即用型 spec 的工作区。

> 某些 IDE 在添加技能后需要重启。

---

## 验证

```bash
swag2mcp --version
```

预期输出（版本可能不同）：

```
swag2mcp v*.*.*
```

---

## 下一步

- [快速开始](quickstart.md) — 2 分钟内运行起来
