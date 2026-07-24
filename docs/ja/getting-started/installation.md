# インストール

## システム要件

- **macOS、Linux、Windows**（amd64 / arm64）
- **Go 1.26+**（`go install` またはソースからのビルドの場合のみ）

## 互換性

| 方法 | macOS | Linux | Windows |
|------|-------|-------|---------|
| One-liner (curl) | ✅ | ✅ | ❌ |
| Homebrew | ✅ | ✅ | ❌ |
| Scoop | ❌ | ❌ | ✅ |
| APT (deb) | ❌ | ✅ | ❌ |
| RPM | ❌ | ✅ | ❌ |
| Docker | ❌ | ✅ | ❌ |
| GitHub Release | ✅ | ✅ | ✅ |
| go install | ✅ | ✅ | ✅ |
| ソースからビルド | ✅ | ✅ | ✅ |

---

## macOS

### One-liner（推奨）

```bash
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash
```

`/usr/local/bin/swag2mcp` にインストールされます（`/usr/local/bin` が書き込み不可の場合は `~/.local/bin/swag2mcp`）。

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

`$GOPATH/bin` が `$PATH` に含まれていることを確認してください：

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### ソースからビルド

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

---

## Linux

### One-liner（推奨）

```bash
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash
```

`/usr/local/bin/swag2mcp` にインストールされます（`/usr/local/bin` が書き込み不可の場合は `~/.local/bin/swag2mcp`）。

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp
```

### APT（Debian / Ubuntu）

```bash
# 最新リリースから .deb をダウンロード
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_linux_amd64.deb
sudo dpkg -i swag2mcp_linux_amd64.deb
```

### RPM（Fedora / RHEL）

```bash
# 最新リリースから .rpm をダウンロード
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_linux_amd64.rpm
sudo rpm -i swag2mcp_linux_amd64.rpm
```

### Docker

```bash
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

stdio トランスポートで実行：

```bash
docker run --rm -i ghcr.io/mmadfox/swag2mcp:latest swag2mcp mcp
```

HTTP トランスポートで実行：

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

`$GOPATH/bin` が `$PATH` に含まれていることを確認してください：

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### ソースからビルド

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
# 最新リリースをダウンロード
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_windows_amd64.zip
Expand-Archive swag2mcp_windows_amd64.zip -DestinationPath .
move swag2mcp.exe C:\Windows\System32\
```

### go install

```powershell
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

### ソースからビルド

```powershell
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp.exe ./cmd/swag2mcp
```

---

## モックサーバー

`swag2mcp-mock` バイナリは別途ダウンロード可能です。メインのバイナリと同じ方法でインストールしてください：

::: code-group

```bash [macOS / Linux]
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

```powershell [Windows]
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

:::

または [GitHub Releases](https://github.com/mmadfox/swag2mcp/releases) からダウンロード — `swag2mcp-mock_<version>_<os>_<arch>.tar.gz` を探してください。

---

## LLM エージェント経由のインストール

AI 搭載 IDE（OpenCode、Cursor、Claude Desktop、VS Code など）を使用している場合、エージェントを通じて swag2mcp をインストールできます：

1. エージェントに swag2mcp スキルを追加するよう依頼します：

   ```
   ".agents/skills/swag2mcp-cli ディレクトリを作成し、https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md のスキルを .agents/skills/swag2mcp-cli/SKILL.md に追加してください"
   ".agents/skills/swag2mcp-format ディレクトリを作成し、https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md のスキルを .agents/skills/swag2mcp-format/SKILL.md に追加してください"
   ```

2. 次にエージェントに指示します：

   ```
   "swag2mcp をセットアップしてください"
   ```

   エージェントが swag2mcp をダウンロードしてインストールし、使用可能なスペックを含むワークスペースを作成します。

> 一部の IDE ではスキル追加後に再起動が必要です。

---

## 確認

```bash
swag2mcp --version
```

期待される出力（バージョンは異なる場合があります）：

```
swag2mcp v*.*.*
```

---

## 次のステップ

- [クイックスタート](quickstart.md) — 2分で使い始める
