# 설치

## 요구 사항

- **macOS, Linux, Windows** (amd64 / arm64)
- **Go 1.26+** (`go install` 또는 소스에서 빌드하는 경우에만)

## 호환성

| 방법 | macOS | Linux | Windows |
|------|-------|-------|---------|
| 원라인 (curl) | ✅ | ✅ | ❌ |
| Homebrew | ✅ | ✅ | ❌ |
| Scoop | ❌ | ❌ | ✅ |
| APT (deb) | ❌ | ✅ | ❌ |
| RPM | ❌ | ✅ | ❌ |
| Docker | ❌ | ✅ | ❌ |
| GitHub Release | ✅ | ✅ | ✅ |
| go install | ✅ | ✅ | ✅ |
| 소스에서 빌드 | ✅ | ✅ | ✅ |

---

## macOS

### 원라인 (권장)

```bash
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash
```

`/usr/local/bin/swag2mcp`에 설치됩니다(`/usr/local/bin`에 쓸 수 없는 경우 `~/.local/bin/swag2mcp`).

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

`$GOPATH/bin`이 `$PATH`에 있는지 확인하세요:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### 소스에서 빌드

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

---

## Linux

### 원라인 (권장)

```bash
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash
```

`/usr/local/bin/swag2mcp`에 설치됩니다(`/usr/local/bin`에 쓸 수 없는 경우 `~/.local/bin/swag2mcp`).

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp
```

### APT (Debian / Ubuntu)

```bash
# 최신 릴리스에서 .deb 다운로드
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_linux_amd64.deb
sudo dpkg -i swag2mcp_linux_amd64.deb
```

### RPM (Fedora / RHEL)

```bash
# 최신 릴리스에서 .rpm 다운로드
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_linux_amd64.rpm
sudo rpm -i swag2mcp_linux_amd64.rpm
```

### Docker

```bash
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

stdio 전송으로 실행:

```bash
docker run --rm -i ghcr.io/mmadfox/swag2mcp:latest swag2mcp mcp
```

HTTP 전송으로 실행:

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

`$GOPATH/bin`이 `$PATH`에 있는지 확인하세요:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### 소스에서 빌드

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
# 최신 릴리스 다운로드
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_windows_amd64.zip
Expand-Archive swag2mcp_windows_amd64.zip -DestinationPath .
move swag2mcp.exe C:\Windows\System32\
```

### go install

```powershell
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

### 소스에서 빌드

```powershell
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp.exe ./cmd/swag2mcp
```

---

## 모의 서버

`swag2mcp-mock` 바이너리는 별도 다운로드로 제공됩니다. 메인 바이너리와 동일한 방법으로 설치하세요:

::: code-group

```bash [macOS / Linux]
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

```powershell [Windows]
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

:::

또는 [GitHub Releases](https://github.com/mmadfox/swag2mcp/releases)에서 다운로드 — `swag2mcp-mock_<version>_<os>_<arch>.tar.gz`를 찾으세요.

---

## LLM 에이전트를 통한 설치

AI 기반 IDE(OpenCode, Cursor, Claude Desktop, VS Code 등)를 사용하는 경우 에이전트를 통해 swag2mcp를 설치할 수 있습니다:

1. 에이전트에게 swag2mcp 스킬을 추가하도록 요청하세요:

   ```
   ".agents/skills/swag2mcp-cli 디렉토리를 만들고 https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md의 스킬을 .agents/skills/swag2mcp-cli/SKILL.md에 추가하세요"
   ".agents/skills/swag2mcp-format 디렉토리를 만들고 https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md의 스킬을 .agents/skills/swag2mcp-format/SKILL.md에 추가하세요"
   ```

2. 그런 다음 에이전트에게 말하세요:

   ```
   "swag2mcp 설정"
   ```

   에이전트가 swag2mcp를 다운로드 및 설치하고 사용 준비가 된 명세가 있는 워크스페이스를 생성합니다.

> 일부 IDE는 스킬 추가 후 재시작이 필요합니다.

---

## 확인

```bash
swag2mcp --version
```

예상 출력 (버전은 다를 수 있음):

```
swag2mcp v*.*.*
```

---

## 다음 단계

- [빠른 시작](quickstart.md) — 2분 만에 실행
