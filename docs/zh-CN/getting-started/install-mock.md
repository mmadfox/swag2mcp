# 安装 Mock 服务器

`swag2mcp-mock` 二进制文件是一个独立的工具，用于使用模拟响应测试 API。它支持与主二进制文件相同的安装方法。

## 兼容性

| 方法 | macOS | Linux | Windows |
|--------|-------|-------|---------|
| 一行命令 (curl) | ✅ | ✅ | ❌ |
| Homebrew | ✅ | ❌ | ❌ |
| Scoop | ❌ | ❌ | ✅ |
| GitHub Release | ✅ | ✅ | ✅ |
| go install | ✅ | ✅ | ✅ |
| 从源码构建 | ✅ | ✅ | ✅ |
| Docker | ✅ | ✅ | ✅ |

---

## macOS

### 一行命令（推荐）

两种安装模式 — 选择适合您的方式：

```bash
# 使用 sudo（安装到 /usr/local/bin — 推荐）
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash

# 不使用 sudo（安装到 ~/.local/bin）
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash -s -- --local
```

**`--sudo`（默认）：** 使用 `sudo` 安装到 `/usr/local/bin/swag2mcp-mock`。系统会提示您输入密码。如果 `sudo` 失败，则回退到 `~/.local/bin/swag2mcp-mock`。

**`--local`：** 不使用 `sudo` 安装到 `~/.local/bin/swag2mcp-mock`。安装后，添加到您的 shell 配置中：

```bash
export PATH="$HOME/.local/bin:$PATH"
```

验证：

```bash
swag2mcp-mock --version
```

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp-mock
```

验证：

```bash
swag2mcp-mock --version
```

### GitHub Release

1. 打开[最新发布页面](https://github.com/mmadfox/swag2mcp/releases/latest)
2. 下载适用于您 Mac 的压缩包：
   - **Apple Silicon**：`swag2mcp-mock_darwin_arm64.tar.gz`
   - **Intel**：`swag2mcp-mock_darwin_amd64.tar.gz`
3. 在下载文件夹中打开终端并运行：

```bash
tar -xzf swag2mcp-mock_darwin_*.tar.gz
sudo mv swag2mcp-mock /usr/local/bin/
```

验证：

```bash
swag2mcp-mock --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

确保 `$GOPATH/bin` 在您的 `$PATH` 中：

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

验证：

```bash
swag2mcp-mock --version
```

### 从源码构建

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp-mock ./cmd/swag2mcp-mock
sudo mv swag2mcp-mock /usr/local/bin/
```

验证：

```bash
swag2mcp-mock --version
```

### Docker

容器以非 root 用户身份运行，需要访问您的工作区目录。
将 `~/.swag2mcp`（或您的自定义工作区路径）挂载到 `/home/nonroot/.swag2mcp`：

> **Apple Silicon (M1/M2/M3/M4) 用户：** 添加 `--platform linux/amd64` 来运行镜像：
> `docker run --rm --platform linux/amd64 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ...`

```bash
docker pull ghcr.io/mmadfox/swag2mcp-mock:latest
```

运行 mock 服务器：

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest
```

验证：

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest --version
```

---

## Linux

### 一行命令（推荐）

两种安装模式 — 选择适合您的方式：

```bash
# 使用 sudo（安装到 /usr/local/bin — 推荐）
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash

# 不使用 sudo（安装到 ~/.local/bin）
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash -s -- --local
```

**`--sudo`（默认）：** 使用 `sudo` 安装到 `/usr/local/bin/swag2mcp-mock`。系统会提示您输入密码。如果 `sudo` 失败，则回退到 `~/.local/bin/swag2mcp-mock`。

**`--local`：** 不使用 `sudo` 安装到 `~/.local/bin/swag2mcp-mock`。安装后，添加到您的 shell 配置中：

```bash
export PATH="$HOME/.local/bin:$PATH"
```

验证：

```bash
swag2mcp-mock --version
```

### GitHub Release

1. 打开[最新发布页面](https://github.com/mmadfox/swag2mcp/releases/latest)
2. 下载适用于您架构的压缩包：
   - **amd64**：`swag2mcp-mock_linux_amd64.tar.gz`
   - **arm64**：`swag2mcp-mock_linux_arm64.tar.gz`
3. 在下载文件夹中打开终端并运行：

```bash
tar -xzf swag2mcp-mock_linux_*.tar.gz
sudo mv swag2mcp-mock /usr/local/bin/
```

验证：

```bash
swag2mcp-mock --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

确保 `$GOPATH/bin` 在您的 `$PATH` 中：

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

验证：

```bash
swag2mcp-mock --version
```

### 从源码构建

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp-mock ./cmd/swag2mcp-mock
sudo mv swag2mcp-mock /usr/local/bin/
```

验证：

```bash
swag2mcp-mock --version
```

### Docker

容器以非 root 用户身份运行，需要访问您的工作区目录。
将 `~/.swag2mcp`（或您的自定义工作区路径）挂载到 `/home/nonroot/.swag2mcp`：

```bash
docker pull ghcr.io/mmadfox/swag2mcp-mock:latest
```

运行 mock 服务器：

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest
```

验证：

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest --version
```

---

## Windows

### Scoop

```powershell
scoop bucket add mmadfox https://github.com/mmadfox/scoop-bucket
scoop install mmadfox/swag2mcp-mock
```

验证：

```powershell
swag2mcp-mock --version
```

### GitHub Release

1. 打开[最新发布页面](https://github.com/mmadfox/swag2mcp/releases/latest)
2. 下载 `swag2mcp-mock_windows_amd64.zip`
3. 解压 ZIP 文件（右键单击 → 全部提取，或使用 PowerShell）
4. 将 `swag2mcp-mock.exe` 移动到 `C:\Windows\System32\`

验证：

```powershell
swag2mcp-mock --version
```

### go install

```powershell
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

验证：

```powershell
swag2mcp-mock --version
```

### 从源码构建

```powershell
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp-mock.exe ./cmd/swag2mcp-mock
```

验证：

```powershell
swag2mcp-mock --version
```

---

## 下一步

- [快速入门](quickstart.md) — 2 分钟开始使用
- [Mock 服务器](../advanced/mock-server.md) — 配置和运行 mock 服务器
