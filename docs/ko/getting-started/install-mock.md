# Mock 서버 설치

`swag2mcp-mock` 바이너리는 가짜 응답으로 API를 테스트하기 위한 별도의 도구입니다. 메인 바이너리와 동일한 설치 방법을 지원합니다.

## 호환성

| 방법 | macOS | Linux | Windows |
|--------|-------|-------|---------|
| 한 줄 명령 (curl) | ✅ | ✅ | ❌ |
| Homebrew | ✅ | ❌ | ❌ |
| Scoop | ❌ | ❌ | ✅ |
| GitHub Release | ✅ | ✅ | ✅ |
| go install | ✅ | ✅ | ✅ |
| 소스에서 빌드 | ✅ | ✅ | ✅ |
| Docker | ✅ | ✅ | ✅ |

---

## macOS

### 한 줄 명령 (권장)

두 가지 설치 모드 — 적합한 방식을 선택하세요:

```bash
# sudo 사용 (/usr/local/bin에 설치 — 권장)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash

# sudo 없이 (~/.local/bin에 설치)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash -s -- --local
```

**`--sudo` (기본값):** `sudo`를 사용하여 `/usr/local/bin/swag2mcp-mock`에 설치합니다. 비밀번호를 입력하라는 메시지가 표시됩니다. `sudo`가 실패하면 `~/.local/bin/swag2mcp-mock`으로 대체됩니다.

**`--local`:** `sudo` 없이 `~/.local/bin/swag2mcp-mock`에 설치합니다. 설치 후 셸 구성에 추가하세요:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

확인:

```bash
swag2mcp-mock --version
```

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp-mock
```

확인:

```bash
swag2mcp-mock --version
```

### GitHub Release

1. [최신 릴리스 페이지](https://github.com/mmadfox/swag2mcp/releases/latest) 열기
2. Mac에 맞는 아카이브 다운로드:
   - **Apple Silicon**: `swag2mcp-mock_darwin_arm64.tar.gz`
   - **Intel**: `swag2mcp-mock_darwin_amd64.tar.gz`
3. 다운로드 폴더에서 터미널을 열고 실행:

```bash
tar -xzf swag2mcp-mock_darwin_*.tar.gz
sudo mv swag2mcp-mock /usr/local/bin/
```

확인:

```bash
swag2mcp-mock --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

`$GOPATH/bin`이 `$PATH`에 있는지 확인하세요:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

확인:

```bash
swag2mcp-mock --version
```

### 소스에서 빌드

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp-mock ./cmd/swag2mcp-mock
sudo mv swag2mcp-mock /usr/local/bin/
```

확인:

```bash
swag2mcp-mock --version
```

### Docker

컨테이너는 비-root 사용자로 실행되며 워크스페이스 디렉토리에 대한 액세스가 필요합니다.
`~/.swag2mcp`(또는 사용자 정의 워크스페이스 경로)를 `/home/nonroot/.swag2mcp`에 마운트하세요:

> **Apple Silicon (M1/M2/M3/M4) 사용자:** 이미지를 실행하려면 `--platform linux/amd64`를 추가하세요:
> `docker run --rm --platform linux/amd64 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ...`

```bash
docker pull ghcr.io/mmadfox/swag2mcp-mock:latest
```

Mock 서버 실행:

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest mockserver
```

확인:

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest --version
```

---

## Linux

### 한 줄 명령 (권장)

두 가지 설치 모드 — 적합한 방식을 선택하세요:

```bash
# sudo 사용 (/usr/local/bin에 설치 — 권장)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash

# sudo 없이 (~/.local/bin에 설치)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash -s -- --local
```

**`--sudo` (기본값):** `sudo`를 사용하여 `/usr/local/bin/swag2mcp-mock`에 설치합니다. 비밀번호를 입력하라는 메시지가 표시됩니다. `sudo`가 실패하면 `~/.local/bin/swag2mcp-mock`으로 대체됩니다.

**`--local`:** `sudo` 없이 `~/.local/bin/swag2mcp-mock`에 설치합니다. 설치 후 셸 구성에 추가하세요:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

확인:

```bash
swag2mcp-mock --version
```

### GitHub Release

1. [최신 릴리스 페이지](https://github.com/mmadfox/swag2mcp/releases/latest) 열기
2. 아키텍처에 맞는 아카이브 다운로드:
   - **amd64**: `swag2mcp-mock_linux_amd64.tar.gz`
   - **arm64**: `swag2mcp-mock_linux_arm64.tar.gz`
3. 다운로드 폴더에서 터미널을 열고 실행:

```bash
tar -xzf swag2mcp-mock_linux_*.tar.gz
sudo mv swag2mcp-mock /usr/local/bin/
```

확인:

```bash
swag2mcp-mock --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

`$GOPATH/bin`이 `$PATH`에 있는지 확인하세요:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

확인:

```bash
swag2mcp-mock --version
```

### 소스에서 빌드

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp-mock ./cmd/swag2mcp-mock
sudo mv swag2mcp-mock /usr/local/bin/
```

확인:

```bash
swag2mcp-mock --version
```

### Docker

컨테이너는 비-root 사용자로 실행되며 워크스페이스 디렉토리에 대한 액세스가 필요합니다.
`~/.swag2mcp`(또는 사용자 정의 워크스페이스 경로)를 `/home/nonroot/.swag2mcp`에 마운트하세요:

```bash
docker pull ghcr.io/mmadfox/swag2mcp-mock:latest
```

Mock 서버 실행:

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest mockserver
```

확인:

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

확인:

```powershell
swag2mcp-mock --version
```

### GitHub Release

1. [최신 릴리스 페이지](https://github.com/mmadfox/swag2mcp/releases/latest) 열기
2. `swag2mcp-mock_windows_amd64.zip` 다운로드
3. ZIP 파일 압축 풀기 (마우스 오른쪽 버튼 클릭 → 모두 추출, 또는 PowerShell 사용)
4. `swag2mcp-mock.exe`를 `C:\Windows\System32\`로 이동

확인:

```powershell
swag2mcp-mock --version
```

### go install

```powershell
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

확인:

```powershell
swag2mcp-mock --version
```

### 소스에서 빌드

```powershell
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp-mock.exe ./cmd/swag2mcp-mock
```

확인:

```powershell
swag2mcp-mock --version
```

---

## 다음 단계

- [빠른 시작](quickstart.md) — 2분 만에 시작하기
- [Mock 서버](../advanced/mock-server.md) — Mock 서버 구성 및 실행
