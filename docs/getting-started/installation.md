# Installation

> [!WARNING]
> **Work in Progress** — This section is being updated. The installation methods below are functional but documentation is under active development.

## Requirements

- **Go 1.26+** (for building from source)
- **macOS, Linux, or Windows** (amd64 / arm64)

## Option 1: GitHub Releases (recommended)

Download from the [releases page](https://github.com/mmadfox/swag2mcp/releases):

| Platform | Architecture | Format |
|----------|-------------|--------|
| macOS | amd64 | tar.gz |
| macOS | arm64 (Apple Silicon) | tar.gz |
| Linux | amd64 | tar.gz |
| Linux | arm64 | tar.gz |
| Windows | amd64 | zip |

```bash
# Example for macOS ARM64
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_darwin_arm64.tar.gz
tar -xzf swag2mcp_darwin_arm64.tar.gz
sudo mv swag2mcp /usr/local/bin/
```

## Option 2: go install

```bash
go install github.com/mmadfox/swag2mcp@latest
```

Ensure `$GOPATH/bin` is in your `$PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## Option 3: Build from source

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

## Verify

```bash
swag2mcp version
```

Expected output:

```
swag2mcp version dev
```

## Mock Server

The `swag2mcp-mock` binary is included:

```bash
swag2mcp-mock mockserver
```

## Next Steps

- [Quick Start](quickstart.md) — get running in 2 minutes
- [First API](first-api.md) — add your first spec
