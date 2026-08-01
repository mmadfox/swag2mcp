# Install Mock Server

The `swag2mcp-mock` binary is a separate tool for testing APIs with fake responses. It supports the same installation methods as the main binary.

## Compatibility

| Method | macOS | Linux | Windows |
|--------|-------|-------|---------|
| One-liner (curl) | ✅ | ✅ | ❌ |
| Homebrew | ✅ | ❌ | ❌ |
| Scoop | ❌ | ❌ | ✅ |
| GitHub Release | ✅ | ✅ | ✅ |
| go install | ✅ | ✅ | ✅ |
| Build from source | ✅ | ✅ | ✅ |
| Docker | ✅ | ✅ | ✅ |

---

## macOS

### One-liner (recommended)

Two install modes — choose the one that fits:

```bash
# With sudo (installs to /usr/local/bin — recommended)
curl -fsSL https://swag2mcp.io/install-mock.sh | bash

# Without sudo (installs to ~/.local/bin)
curl -fsSL https://swag2mcp.io/install-mock.sh | bash -s -- --local
```

**`--sudo` (default):** Installs to `/usr/local/bin/swag2mcp-mock` using `sudo`. You will be prompted for your password. If `sudo` fails, falls back to `~/.local/bin/swag2mcp-mock`.

**`--local`:** Installs to `~/.local/bin/swag2mcp-mock` without `sudo`. After install, add to your shell config:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Verify:

```bash
swag2mcp-mock --version
```

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp-mock
```

Verify:

```bash
swag2mcp-mock --version
```

### GitHub Release

1. Open the [latest release page](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Download the archive for your Mac:
   - **Apple Silicon**: `swag2mcp-mock_darwin_arm64.tar.gz`
   - **Intel**: `swag2mcp-mock_darwin_amd64.tar.gz`
3. Open Terminal in the download folder and run:

```bash
tar -xzf swag2mcp-mock_darwin_*.tar.gz
sudo mv swag2mcp-mock /usr/local/bin/
```

Verify:

```bash
swag2mcp-mock --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

Ensure `$GOPATH/bin` is in your `$PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Verify:

```bash
swag2mcp-mock --version
```

### Build from source

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp-mock ./cmd/swag2mcp-mock
sudo mv swag2mcp-mock /usr/local/bin/
```

Verify:

```bash
swag2mcp-mock --version
```

### Docker

Mount `~/.swag2mcp` (or your custom workspace path) to `/home/nonroot/.swag2mcp`.
The entrypoint automatically adjusts file permissions so the container can read your workspace.

> **Apple Silicon (M1/M2/M3/M4) users:** Add `--platform linux/amd64` to run the image:
> `docker run --rm --platform linux/amd64 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ...`

```bash
docker pull ghcr.io/mmadfox/swag2mcp-mock:latest
```

Run mock server:

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest
```

Verify:

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest --version
```

---

## Linux

### One-liner (recommended)

Two install modes — choose the one that fits:

```bash
# With sudo (installs to /usr/local/bin — recommended)
curl -fsSL https://swag2mcp.io/install-mock.sh | bash

# Without sudo (installs to ~/.local/bin)
curl -fsSL https://swag2mcp.io/install-mock.sh | bash -s -- --local
```

**`--sudo` (default):** Installs to `/usr/local/bin/swag2mcp-mock` using `sudo`. You will be prompted for your password. If `sudo` fails, falls back to `~/.local/bin/swag2mcp-mock`.

**`--local`:** Installs to `~/.local/bin/swag2mcp-mock` without `sudo`. After install, add to your shell config:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Verify:

```bash
swag2mcp-mock --version
```

### GitHub Release

1. Open the [latest release page](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Download the archive for your architecture:
   - **amd64**: `swag2mcp-mock_linux_amd64.tar.gz`
   - **arm64**: `swag2mcp-mock_linux_arm64.tar.gz`
3. Open Terminal in the download folder and run:

```bash
tar -xzf swag2mcp-mock_linux_*.tar.gz
sudo mv swag2mcp-mock /usr/local/bin/
```

Verify:

```bash
swag2mcp-mock --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

Ensure `$GOPATH/bin` is in your `$PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Verify:

```bash
swag2mcp-mock --version
```

### Build from source

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp-mock ./cmd/swag2mcp-mock
sudo mv swag2mcp-mock /usr/local/bin/
```

Verify:

```bash
swag2mcp-mock --version
```

### Docker

Mount `~/.swag2mcp` (or your custom workspace path) to `/home/nonroot/.swag2mcp`.
The entrypoint automatically adjusts file permissions so the container can read your workspace.

```bash
docker pull ghcr.io/mmadfox/swag2mcp-mock:latest
```

Run mock server:

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest
```

Verify:

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

Verify:

```powershell
swag2mcp-mock --version
```

### GitHub Release

1. Open the [latest release page](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Download `swag2mcp-mock_windows_amd64.zip`
3. Extract the ZIP file (right-click → Extract All, or use PowerShell)
4. Move `swag2mcp-mock.exe` to `C:\Windows\System32\`

Verify:

```powershell
swag2mcp-mock --version
```

### go install

```powershell
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

Verify:

```powershell
swag2mcp-mock --version
```

### Build from source

```powershell
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp-mock.exe ./cmd/swag2mcp-mock
```

Verify:

```powershell
swag2mcp-mock --version
```

---

## Next Steps

- [Quick Start](quickstart.md) — get running in 2 minutes
- [Mock Server](../advanced/mock-server.md) — configure and run the mock server
