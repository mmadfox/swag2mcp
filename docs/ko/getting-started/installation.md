# 설치

## 요구 사항

- **macOS, Linux, Windows** (amd64 / arm64)
- **Go 1.26+** (`go install` 또는 소스에서 빌드하는 경우에만)

## 호환성

| 방법 | macOS | Linux | Windows |
|------|-------|-------|---------|
| 원라인 (curl) | ✅ | ✅ | ❌ |
| Homebrew | ✅ | ❌ | ❌ |
| Scoop | ❌ | ❌ | ✅ |
| APT (deb) | ❌ | ✅ | ❌ |
| RPM | ❌ | ✅ | ❌ |
| Docker | ✅ | ✅ | ✅ |
| GitHub Release | ✅ | ✅ | ✅ |
| go install | ✅ | ✅ | ✅ |
| 소스에서 빌드 | ✅ | ✅ | ✅ |

---

## macOS

### 원라인 (권장)

두 가지 설치 모드 — 적합한 방식을 선택하세요:

```bash
# sudo 사용 (/usr/local/bin에 설치 — 권장)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash

# sudo 없이 (~/.local/bin에 설치)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash -s -- --local
```

**`--sudo` (기본값):** `sudo`를 사용하여 `/usr/local/bin/swag2mcp`에 설치합니다. 비밀번호를 입력하라는 메시지가 표시됩니다. `sudo`가 실패하면 `~/.local/bin/swag2mcp`로 대체됩니다.

**`--local`:** `sudo` 없이 `~/.local/bin/swag2mcp`에 설치합니다. 설치 후 셸 설정에 추가하세요:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

**확인:**

```bash
swag2mcp --version
```

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp
```

**확인:**

```bash
swag2mcp --version
```

### GitHub Release

1. [최신 릴리스 페이지](https://github.com/mmadfox/swag2mcp/releases/latest)를 엽니다
2. Mac에 맞는 아카이브를 다운로드하세요:
   - **Apple Silicon**: `swag2mcp_darwin_arm64.tar.gz`
   - **Intel**: `swag2mcp_darwin_amd64.tar.gz`
3. 다운로드 폴더에서 터미널을 열고 실행:

```bash
tar -xzf swag2mcp_darwin_*.tar.gz
sudo mv swag2mcp /usr/local/bin/
```

**확인:**

```bash
swag2mcp --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

`$GOPATH/bin`이 `$PATH`에 있는지 확인하세요:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

**확인:**

```bash
swag2mcp --version
```

### 소스에서 빌드

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

**확인:**

```bash
swag2mcp --version
```

### Docker

`~/.swag2mcp`(또는 사용자 정의 경로)를 `/home/nonroot/.swag2mcp`에 마운트하세요.
엔트리포인트가 자동으로 파일 권한을 조정하여 컨테이너가 작업 공간을 읽을 수 있도록 합니다.

> **Apple Silicon (M1/M2/M3/M4) 사용자:** `--platform linux/amd64`를 추가하여 이미지를 실행하세요:
> `docker run --rm --platform linux/amd64 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ...`

```bash
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

stdio 전송으로 실행:

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp
```

HTTP 전송으로 실행:

```bash
docker run --rm -p 8080:8080 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp --transport sse --http-addr :8080
```

**확인:**

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest --version
```

---

## Linux

### 원라인 (권장)

두 가지 설치 모드 — 적합한 방식을 선택하세요:

```bash
# sudo 사용 (/usr/local/bin에 설치 — 권장)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash

# sudo 없이 (~/.local/bin에 설치)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash -s -- --local
```

**`--sudo` (기본값):** `sudo`를 사용하여 `/usr/local/bin/swag2mcp`에 설치합니다. 비밀번호를 입력하라는 메시지가 표시됩니다. `sudo`가 실패하면 `~/.local/bin/swag2mcp`로 대체됩니다.

**`--local`:** `sudo` 없이 `~/.local/bin/swag2mcp`에 설치합니다. 설치 후 셸 설정에 추가하세요:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

**확인:**

```bash
swag2mcp --version
```

### APT (Debian / Ubuntu)

1. [최신 릴리스 페이지](https://github.com/mmadfox/swag2mcp/releases/latest)를 엽니다
2. `swag2mcp_linux_amd64.deb`를 다운로드하세요
3. 설치:

```bash
# 다운로드 폴더에 있는지 확인하세요
pwd
sudo dpkg -i swag2mcp_linux_amd64.deb
```

**확인:**

```bash
swag2mcp --version
```

### RPM (Fedora / RHEL)

1. [최신 릴리스 페이지](https://github.com/mmadfox/swag2mcp/releases/latest)를 엽니다
2. `swag2mcp_linux_amd64.rpm`를 다운로드하세요
3. 설치:

```bash
# 다운로드 폴더에 있는지 확인하세요
pwd
sudo rpm -i swag2mcp_linux_amd64.rpm
```

**확인:**

```bash
swag2mcp --version
```

### Docker

`~/.swag2mcp`(또는 사용자 정의 경로)를 `/home/nonroot/.swag2mcp`에 마운트하세요.
엔트리포인트가 자동으로 파일 권한을 조정하여 컨테이너가 작업 공간을 읽을 수 있도록 합니다.

> **Apple Silicon (M1/M2/M3/M4) 사용자:** `--platform linux/amd64`를 추가하여 이미지를 실행하세요:
> `docker run --rm --platform linux/amd64 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ...`

```bash
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

stdio 전송으로 실행:

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp
```

HTTP 전송으로 실행:

```bash
docker run --rm -p 8080:8080 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp --transport sse --http-addr :8080
```

**확인:**

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest --version
```

### GitHub Release

1. [최신 릴리스 페이지](https://github.com/mmadfox/swag2mcp/releases/latest)를 엽니다
2. 아키텍처에 맞는 아카이브를 다운로드하세요:
   - **amd64**: `swag2mcp_linux_amd64.tar.gz`
   - **arm64**: `swag2mcp_linux_arm64.tar.gz`
3. 다운로드 폴더에서 터미널을 열고 실행:

```bash
tar -xzf swag2mcp_linux_*.tar.gz
sudo mv swag2mcp /usr/local/bin/
```

**확인:**

```bash
swag2mcp --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

`$GOPATH/bin`이 `$PATH`에 있는지 확인하세요:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

**확인:**

```bash
swag2mcp --version
```

### 소스에서 빌드

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

**확인:**

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

**확인:**

```powershell
swag2mcp --version
```

### GitHub Release

1. [최신 릴리스 페이지](https://github.com/mmadfox/swag2mcp/releases/latest)를 엽니다
2. `swag2mcp_windows_amd64.zip`을 다운로드하세요
3. ZIP 파일을 압축 해제합니다 (마우스 오른쪽 버튼 → 모두 추출, 또는 PowerShell 사용)
4. `swag2mcp.exe`를 `C:\Windows\System32\`로 이동

**확인:**

```powershell
swag2mcp --version
```

### go install

```powershell
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

**확인:**

```powershell
swag2mcp --version
```

### 소스에서 빌드

```powershell
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp.exe ./cmd/swag2mcp
```

**확인:**

```powershell
swag2mcp --version
```

---

> 모크 서버가 필요하신가요? [Mock Server 설치](install-mock.md)를 참조하세요.

## LLM 에이전트를 통한 설치

AI 기반 IDE(OpenCode, Cursor, Claude Desktop, VS Code 등)를 사용하는 경우 [LLM 스킬](/ko/llm-skills/overview)을 참조하여 에이전트를 통해 또는 수동으로 스킬을 설치하는 방법을 확인하세요.

> 일부 IDE 및 LLM 클라이언트는 스킬 추가 후 재시작이 필요합니다.

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
