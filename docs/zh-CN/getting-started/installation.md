# 安装

## 要求

- **macOS、Linux 或 Windows**（amd64 / arm64）
- **Go 1.26+**（仅用于 `go install` 或从源码构建）

## 兼容性

| 方法 | macOS | Linux | Windows |
|------|-------|-------|---------|
| 一行命令（curl） | ✅ | ✅ | ❌ |
| Homebrew | ✅ | ❌ | ❌ |
| Scoop | ❌ | ❌ | ✅ |
| APT（deb） | ❌ | ✅ | ❌ |
| RPM | ❌ | ✅ | ❌ |
| Docker | ✅ | ✅ | ✅ |
| GitHub Release | ✅ | ✅ | ✅ |
| go install | ✅ | ✅ | ✅ |
| 从源码构建 | ✅ | ✅ | ✅ |

---

## macOS

### 一行命令（推荐）

两种安装模式 — 选择适合您的一种：

```bash
# 使用 sudo（安装到 /usr/local/bin — 推荐）
curl -fsSL https://swag2mcp.io/install.sh | bash

# 不使用 sudo（安装到 ~/.local/bin）
curl -fsSL https://swag2mcp.io/install.sh | bash -s -- --local
```

**`--sudo`（默认）：** 使用 `sudo` 安装到 `/usr/local/bin/swag2mcp`。系统会提示您输入密码。如果 `sudo` 失败，则回退到 `~/.local/bin/swag2mcp`。

**`--local`：** 无需 `sudo` 安装到 `~/.local/bin/swag2mcp`。安装后，添加到您的 shell 配置中：

```bash
export PATH="$HOME/.local/bin:$PATH"
```

**验证：**

```bash
swag2mcp --version
```

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp
```

**验证：**

```bash
swag2mcp --version
```

### GitHub Release

1. 打开[最新发布页面](https://github.com/mmadfox/swag2mcp/releases/latest)
2. 下载适合您 Mac 的压缩包：
   - **Apple Silicon**：`swag2mcp_darwin_arm64.tar.gz`
   - **Intel**：`swag2mcp_darwin_amd64.tar.gz`
3. 在下载文件夹中打开终端并运行：

```bash
tar -xzf swag2mcp_darwin_*.tar.gz
sudo mv swag2mcp /usr/local/bin/
```

**验证：**

```bash
swag2mcp --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

确保 `$GOPATH/bin` 在你的 `$PATH` 中：

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

**验证：**

```bash
swag2mcp --version
```

### 从源码构建

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

**验证：**

```bash
swag2mcp --version
```

### Docker

将 `~/.swag2mcp`（或您的自定义路径）挂载到 `/home/nonroot/.swag2mcp`。
入口点会自动调整文件权限，使容器能够读取您的工作区。

> **Apple Silicon (M1/M2/M3/M4) 用户：** 添加 `--platform linux/amd64` 来运行镜像：
> `docker run --rm --platform linux/amd64 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ...`

```bash
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

使用 stdio 传输运行：

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp
```

使用 HTTP 传输运行：

```bash
docker run --rm -p 8080:8080 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp --transport sse --http-addr :8080
```

**验证：**

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest --version
```

---

## Linux

### 一行命令（推荐）

两种安装模式 — 选择适合您的一种：

```bash
# 使用 sudo（安装到 /usr/local/bin — 推荐）
curl -fsSL https://swag2mcp.io/install.sh | bash

# 不使用 sudo（安装到 ~/.local/bin）
curl -fsSL https://swag2mcp.io/install.sh | bash -s -- --local
```

**`--sudo`（默认）：** 使用 `sudo` 安装到 `/usr/local/bin/swag2mcp`。系统会提示您输入密码。如果 `sudo` 失败，则回退到 `~/.local/bin/swag2mcp`。

**`--local`：** 无需 `sudo` 安装到 `~/.local/bin/swag2mcp`。安装后，添加到您的 shell 配置中：

```bash
export PATH="$HOME/.local/bin:$PATH"
```

**验证：**

```bash
swag2mcp --version
```

### APT（Debian / Ubuntu）

1. 打开[最新发布页面](https://github.com/mmadfox/swag2mcp/releases/latest)
2. 下载 `swag2mcp_linux_amd64.deb`
3. 安装：

```bash
# 确保您在下载文件夹中
pwd
sudo dpkg -i swag2mcp_linux_amd64.deb
```

**验证：**

```bash
swag2mcp --version
```

### RPM（Fedora / RHEL）

1. 打开[最新发布页面](https://github.com/mmadfox/swag2mcp/releases/latest)
2. 下载 `swag2mcp_linux_amd64.rpm`
3. 安装：

```bash
# 确保您在下载文件夹中
pwd
sudo rpm -i swag2mcp_linux_amd64.rpm
```

**验证：**

```bash
swag2mcp --version
```

### Docker

将 `~/.swag2mcp`（或您的自定义路径）挂载到 `/home/nonroot/.swag2mcp`。
入口点会自动调整文件权限，使容器能够读取您的工作区。

> **Apple Silicon (M1/M2/M3/M4) 用户：** 添加 `--platform linux/amd64` 来运行镜像：
> `docker run --rm --platform linux/amd64 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ...`

```bash
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

使用 stdio 传输运行：

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp
```

使用 HTTP 传输运行：

```bash
docker run --rm -p 8080:8080 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp --transport sse --http-addr :8080
```

**验证：**

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest --version
```

### GitHub Release

1. 打开[最新发布页面](https://github.com/mmadfox/swag2mcp/releases/latest)
2. 下载适合您架构的压缩包：
   - **amd64**：`swag2mcp_linux_amd64.tar.gz`
   - **arm64**：`swag2mcp_linux_arm64.tar.gz`
3. 在下载文件夹中打开终端并运行：

```bash
tar -xzf swag2mcp_linux_*.tar.gz
sudo mv swag2mcp /usr/local/bin/
```

**验证：**

```bash
swag2mcp --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

确保 `$GOPATH/bin` 在你的 `$PATH` 中：

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

**验证：**

```bash
swag2mcp --version
```

### 从源码构建

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

**验证：**

```bash
swag2mcp --version
```

---

## Windows

### Scoop

```powershell
scoop bucket add mmadfox https://github.com/mmadfox/scoop-bucket
scoop install mmadfox/swag2mcp
```

**验证：**

```powershell
swag2mcp --version
```

### GitHub Release

1. 打开[最新发布页面](https://github.com/mmadfox/swag2mcp/releases/latest)
2. 下载 `swag2mcp_windows_amd64.zip`
3. 解压 ZIP 文件（右键 → 全部提取，或使用 PowerShell）
4. 将 `swag2mcp.exe` 移动到 `C:\Windows\System32\`

**验证：**

```powershell
swag2mcp --version
```

### go install

```powershell
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

**验证：**

```powershell
swag2mcp --version
```

### 从源码构建

```powershell
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp.exe ./cmd/swag2mcp
```

**验证：**

```powershell
swag2mcp --version
```

### Docker

将 `~/.swag2mcp`（或您的自定义路径）挂载到 `/home/nonroot/.swag2mcp`。入口点会自动调整文件权限，使容器能够读取您的工作区。

```powershell
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

Run with stdio transport:

```powershell
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp
```

Run with HTTP transport:

```powershell
docker run --rm -p 8080:8080 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp --transport sse --http-addr :8080
```

Verify:

```powershell
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest --version
```

---

> 需要 mock 服务器？请查看[安装 Mock Server](install-mock.md)。

## 通过 LLM 智能体安装

如果您使用 AI 驱动的 IDE（OpenCode、Cursor、Claude Desktop、VS Code 等），请参阅 [LLM 技能](/zh-CN/llm-skills/overview) 了解如何通过智能体或手动安装技能。

> 某些 IDE 和 LLM 客户端在添加技能后需要重启。

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
