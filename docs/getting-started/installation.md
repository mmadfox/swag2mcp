# Installation

## Requirements

- **macOS, Linux, or Windows** (amd64 / arm64)
- **Go 1.26+** (only for `go install` or building from source)

## Compatibility

| Method | macOS | Linux | Windows |
|--------|-------|-------|---------|
| One-liner (curl) | ✅ | ✅ | ❌ |
| Homebrew | ✅ | ❌ | ❌ |
| Scoop | ❌ | ❌ | ✅ |
| APT (deb) | ❌ | ✅ | ❌ |
| RPM | ❌ | ✅ | ❌ |
| Docker | ✅ | ✅ | ✅ |
| GitHub Release | ✅ | ✅ | ✅ |
| go install | ✅ | ✅ | ✅ |
| Build from source | ✅ | ✅ | ✅ |

---

## macOS

### One-liner (recommended)

Two install modes — choose the one that fits:

```bash
# With sudo (installs to /usr/local/bin — recommended)
curl -fsSL https://swag2mcp.io/install.sh | bash

# Without sudo (installs to ~/.local/bin)
curl -fsSL https://swag2mcp.io/install.sh | bash -s -- --local
```

**`--sudo` (default):** Installs to `/usr/local/bin/swag2mcp` using `sudo`. You will be prompted for your password. If `sudo` fails, falls back to `~/.local/bin/swag2mcp`.

**`--local`:** Installs to `~/.local/bin/swag2mcp` without `sudo`. After install, add to your shell config:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Verify:

```bash
swag2mcp --version
```

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp
```

Verify:

```bash
swag2mcp --version
```

### GitHub Release

1. Open the [latest release page](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Download the archive for your Mac:
   - **Apple Silicon**: `swag2mcp_darwin_arm64.tar.gz`
   - **Intel**: `swag2mcp_darwin_amd64.tar.gz`
3. Open Terminal in the download folder and run:

```bash
tar -xzf swag2mcp_darwin_*.tar.gz
sudo mv swag2mcp /usr/local/bin/
```

Verify:

```bash
swag2mcp --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

Ensure `$GOPATH/bin` is in your `$PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Verify:

```bash
swag2mcp --version
```

### Build from source

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

Verify:

```bash
swag2mcp --version
```

### Docker

Mount `~/.swag2mcp` (or your custom workspace path) to `/home/nonroot/.swag2mcp`.
The entrypoint automatically adjusts file permissions so the container can read your workspace.

> **Apple Silicon (M1/M2/M3/M4) users:** Add `--platform linux/amd64` to run the image:
> `docker run --rm --platform linux/amd64 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ...`

```bash
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

Run with stdio transport:

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp
```

Run with HTTP transport:

```bash
docker run --rm -p 8080:8080 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp --transport sse --http-addr :8080
```

Verify:

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest --version
```

---

## Linux

### One-liner (recommended)

Two install modes — choose the one that fits:

```bash
# With sudo (installs to /usr/local/bin — recommended)
curl -fsSL https://swag2mcp.io/install.sh | bash

# Without sudo (installs to ~/.local/bin)
curl -fsSL https://swag2mcp.io/install.sh | bash -s -- --local
```

**`--sudo` (default):** Installs to `/usr/local/bin/swag2mcp` using `sudo`. You will be prompted for your password. If `sudo` fails, falls back to `~/.local/bin/swag2mcp`.

**`--local`:** Installs to `~/.local/bin/swag2mcp` without `sudo`. After install, add to your shell config:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Verify:

```bash
swag2mcp --version
```

### APT (Debian / Ubuntu)

1. Open the [latest release page](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Download `swag2mcp_linux_amd64.deb`
3. Install:

```bash
# Make sure you are in the download folder
pwd
sudo dpkg -i swag2mcp_linux_amd64.deb
```

Verify:

```bash
swag2mcp --version
```

### RPM (Fedora / RHEL)

1. Open the [latest release page](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Download `swag2mcp_linux_amd64.rpm`
3. Install:

```bash
# Make sure you are in the download folder
pwd
sudo rpm -i swag2mcp_linux_amd64.rpm
```

Verify:

```bash
swag2mcp --version
```

### Docker

Mount `~/.swag2mcp` (or your custom workspace path) to `/home/nonroot/.swag2mcp`.
The entrypoint automatically adjusts file permissions so the container can read your workspace.

> **Apple Silicon (M1/M2/M3/M4) users:** Add `--platform linux/amd64` to run the image:
> `docker run --rm --platform linux/amd64 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ...`

```bash
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

Run with stdio transport:

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp
```

Run with HTTP transport:

```bash
docker run --rm -p 8080:8080 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp --transport sse --http-addr :8080
```

Verify:

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest --version
```

### GitHub Release

1. Open the [latest release page](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Download the archive for your architecture:
   - **amd64**: `swag2mcp_linux_amd64.tar.gz`
   - **arm64**: `swag2mcp_linux_arm64.tar.gz`
3. Open Terminal in the download folder and run:

```bash
tar -xzf swag2mcp_linux_*.tar.gz
sudo mv swag2mcp /usr/local/bin/
```

Verify:

```bash
swag2mcp --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

Ensure `$GOPATH/bin` is in your `$PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Verify:

```bash
swag2mcp --version
```

### Build from source

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

Verify:

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

Verify:

```powershell
swag2mcp --version
```

### GitHub Release

1. Open the [latest release page](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Download `swag2mcp_windows_amd64.zip`
3. Extract the ZIP file (right-click → Extract All, or use PowerShell)
4. Move `swag2mcp.exe` to `C:\Windows\System32\`

Verify:

```powershell
swag2mcp --version
```

### go install

```powershell
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

Verify:

```powershell
swag2mcp --version
```

### Build from source

```powershell
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp.exe ./cmd/swag2mcp
```

Verify:

```powershell
swag2mcp --version
```

### Docker

Mount `~/.swag2mcp` (or your custom workspace path) to `/home/nonroot/.swag2mcp`.
The entrypoint automatically adjusts file permissions so the container can read your workspace.

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

> Need the mock server? See [Install Mock Server](install-mock.md).

## Install via LLM Agent

If you use an AI-powered IDE (OpenCode, Cursor, Claude Desktop, VS Code, etc.), see [Skills for LLM](/llm-skills/overview) for how to install swag2mcp skills through your agent or manually.

> Some IDEs and LLM clients require a restart after adding skills.

---

## Next Steps

- [Quick Start](quickstart.md) — get running in 2 minutes
