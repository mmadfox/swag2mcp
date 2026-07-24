# Installation

## Requirements

- **macOS, Linux, or Windows** (amd64 / arm64)
- **Go 1.26+** (only for `go install` or building from source)

## Compatibility

| Method | macOS | Linux | Windows |
|--------|-------|-------|---------|
| One-liner (curl) | ✅ | ✅ | ❌ |
| Homebrew | ✅ | ✅ | ❌ |
| Scoop | ❌ | ❌ | ✅ |
| APT (deb) | ❌ | ✅ | ❌ |
| RPM | ❌ | ✅ | ❌ |
| Docker | ❌ | ✅ | ❌ |
| GitHub Release | ✅ | ✅ | ✅ |
| go install | ✅ | ✅ | ✅ |
| Build from source | ✅ | ✅ | ✅ |

---

## macOS

### One-liner (recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash
```

Installs to `/usr/local/bin/swag2mcp` (or `~/.local/bin/swag2mcp` if `/usr/local/bin` is not writable).

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

Ensure `$GOPATH/bin` is in your `$PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Build from source

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

---

## Linux

### One-liner (recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash
```

Installs to `/usr/local/bin/swag2mcp` (or `~/.local/bin/swag2mcp` if `/usr/local/bin` is not writable).

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp
```

### APT (Debian / Ubuntu)

```bash
# Download the .deb from the latest release
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_linux_amd64.deb
sudo dpkg -i swag2mcp_linux_amd64.deb
```

### RPM (Fedora / RHEL)

```bash
# Download the .rpm from the latest release
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_linux_amd64.rpm
sudo rpm -i swag2mcp_linux_amd64.rpm
```

### Docker

```bash
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

Run with stdio transport:

```bash
docker run --rm -i ghcr.io/mmadfox/swag2mcp:latest swag2mcp mcp
```

Run with HTTP transport:

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

Ensure `$GOPATH/bin` is in your `$PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Build from source

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
# Download the latest release
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_windows_amd64.zip
Expand-Archive swag2mcp_windows_amd64.zip -DestinationPath .
move swag2mcp.exe C:\Windows\System32\
```

### go install

```powershell
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

### Build from source

```powershell
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp.exe ./cmd/swag2mcp
```

---

## Mock Server

The `swag2mcp-mock` binary is available as a separate download. Install it using the same method as the main binary:

::: code-group

```bash [macOS / Linux]
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

```powershell [Windows]
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

:::

Or download from [GitHub Releases](https://github.com/mmadfox/swag2mcp/releases) — look for `swag2mcp-mock_<version>_<os>_<arch>.tar.gz`.

---

## Install via LLM Agent

If you use an AI-powered IDE (OpenCode, Cursor, Claude Desktop, VS Code, etc.), you can install swag2mcp through your agent:

1. Ask your agent to add the swag2mcp skills:

   ```
   "Create the .agents/skills/swag2mcp-cli directory and add the skill from https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md to .agents/skills/swag2mcp-cli/SKILL.md"
   "Create the .agents/skills/swag2mcp-format directory and add the skill from https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md to .agents/skills/swag2mcp-format/SKILL.md"
   ```

2. Then tell your agent:

   ```
   "Set up swag2mcp"
   ```

   The agent will download and install swag2mcp, then create a workspace with ready-to-use specs.

> Some IDEs require a restart after adding skills.

---

## Verify

```bash
swag2mcp --version
```

Expected output (version may vary):

```
swag2mcp v*.*.*
```

---

## Next Steps

- [Quick Start](quickstart.md) — get running in 2 minutes
