# 문제 해결

## 설치 문제

### swag2mcp: command not found

바이너리가 PATH에 없습니다.

```bash
# Go가 설치되어 있는지 확인
go version

# Go가 바이너리를 설치하는 위치 확인
go env GOPATH
# 일반적으로 ~/go 또는 ~/go/bin

# PATH에 추가 (~/.zshrc 또는 ~/.bashrc에 추가)
export PATH=$PATH:$(go env GOPATH)/bin

# 또는 전체 경로 사용
~/go/bin/swag2mcp --version
```

GitHub Releases에서 바이너리를 다운로드한 경우 PATH에 있는 디렉토리에 있는지 확인하세요:

```bash
# /usr/local/bin으로 이동 (macOS/Linux)
sudo mv swag2mcp /usr/local/bin/
```

### permission denied

바이너리에 실행 권한이 없습니다.

```bash
# go install의 경우 (소유권 수정)
sudo chown -R $(whoami) $(go env GOPATH)

# 다운로드한 바이너리의 경우
chmod +x /path/to/swag2mcp
```

### Go 버전이 너무 오래됨

swag2mcp는 Go 1.23+이 필요합니다.

```bash
go version
# 버전이 1.23 미만이면 Go 업데이트:
# https://go.dev/dl/
```

### 모의 서버를 찾을 수 없음

모의 서버는 별도의 바이너리입니다. 명시적으로 설치하세요:

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

## 설정 문제

### 설정 파일을 찾을 수 없음

swag2mcp가 `swag2mcp.yaml`을 찾을 수 없습니다.

```bash
# 새 설정 생성
swag2mcp init

# 또는 경로를 명시적으로 지정
swag2mcp mcp /path/to/workspace
swag2mcp ls /path/to/workspace
```

**일반적인 원인:** 임의의 디렉토리에서 `swag2mcp mcp`를 실행하여 프로젝트의 워크스페이스 대신 `~/.swag2mcp/`를 찾았습니다. 항상 경로를 명시적으로 전달하세요.

### 잘못된 워크스페이스가 로드됨

예상과 다른 워크스페이스가 로드되었습니다.

**해결 순서:** 명시적 `[path]` → 현재 디렉토리(`./`) → `~/.swag2mcp/`. 경로 없이 `swag2mcp mcp`를 `swag2mcp.yaml`이 없는 디렉토리에서 실행하면 `~/.swag2mcp/`로 폴백됩니다.

**해결 방법:** 항상 워크스페이스 경로를 전달하세요: `swag2mcp mcp /path/to/your/workspace`

### YAML 파싱 오류

설정 파일에 잘못된 YAML 구문이 있습니다.

```bash
# 설정 검증
swag2mcp validate

# 일반적인 실수:
# - 공백 대신 탭 사용 (YAML은 공백 필요)
# - 중첩 필드의 들여쓰기 누락
# - 특수 문자가 있는 따옴표 없는 문자열 (: # & {)
```

**팁:** YAML 린터 또는 YAML을 지원하는 편집기를 사용하여 구문 오류를 찾으세요.

### 검증 실패: "no specifications defined"

설정 파일이 존재하지만 spec이 없습니다.

```bash
# spec 추가
swag2mcp add spec

# 또는 swag2mcp.yaml을 편집하여 최소 하나의 spec 추가
```

### 검증 실패: "duplicate domain"

두 spec이 동일한 `domain` 값을 가지고 있습니다. 도메인은 고유해야 합니다.

```bash
# 현재 spec 목록 보기
swag2mcp ls

# swag2mcp.yaml에서 중복 도메인 확인
```

### 검증 실패: "invalid spec location"

`location` URL 또는 파일 경로에 접근할 수 없거나 유효한 spec 파일이 아닙니다.

```bash
# URL에 접근 가능한지 확인
curl -I https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml

# 로컬 파일이 존재하는지 확인
ls -la ./specs/my-api.yaml

# 파일이 유효한 OpenAPI/Swagger/Postman인지 확인
# (단순한 JSON이나 HTML 페이지가 아닌지)
```

**일반적인 원인:** `location` 필드가 API 엔드포인트 자체(예: `https://api.example.com/v1/users`)를 가리키고 있습니다. 위치는 OpenAPI/Swagger/Postman 파일을 가리켜야 합니다.

## MCP 서버 문제

### 포트가 이미 사용 중

다른 프로세스가 포트를 사용 중입니다.

```bash
# 프로세스 찾기
lsof -i :8080

# 종료
kill <PID>

# 또는 다른 포트 사용
swag2mcp mcp --transport sse --http-addr :9090
```

### 연결 거부됨

MCP 서버가 실행 중이 아니거나 접근할 수 없습니다.

```bash
# 서버가 실행 중인지 확인
swag2mcp mcp --transport sse --http-addr 127.0.0.1:8080

# 다른 터미널에서 health 엔드포인트 확인
curl http://127.0.0.1:8080/health

# 커스텀 경로 사용 시
curl http://127.0.0.1:8080/custom-path/health
```

### MCP 도구가 LLM 클라이언트에 표시되지 않음

LLM 클라이언트가 도구를 볼 수 없습니다.

```bash
# spec이 로드되었는지 확인
swag2mcp ls

# spec이 비활성화되지 않았는지 확인
swag2mcp validate

# 서버 로그 확인
swag2mcp mcp --logfile /tmp/swag2mcp.log
cat /tmp/swag2mcp.log

# IDE 설정의 워크스페이스 경로가 올바른지 확인
# (절대 경로여야 함)
```

**일반적인 원인:**
- IDE 설정의 잘못된 워크스페이스 경로
- 모든 spec에 `disable: true` 설정
- `--tags`로 인해 spec이 필터링됨
- 지정된 경로에 설정 파일이 없음

### MCP 핸드셰이크 실패 (HTTP 전송)

SSE 및 Streamable HTTP 전송의 경우 MCP 프로토콜은 도구 호출 전에 초기화가 필요합니다.

```
1단계: POST /mcp → {"method":"initialize", ...}
2단계: POST /mcp → {"method":"notifications/initialized"}
3단계: POST /mcp → {"method":"tools/list", ...}  ← 이제 작동
```

LLM 클라이언트가 도구를 호출하기 전에 핸드셰이크를 완료하는지 확인하세요.

### Health check가 404 반환

Health 엔드포인트 경로가 MCP 경로와 다를 수 있습니다.

```bash
# 기본 health 엔드포인트
curl http://127.0.0.1:8080/health

# MCP 경로를 변경해도 health는 여전히 /health에 있습니다
# (--http-path의 영향을 받지 않음)
```

### Auth 도구를 사용할 수 없음

`auth` MCP 도구가 표시되지 않습니다.

`auth` 도구는 **기본적으로 비활성화**되어 있습니다(`--disable-llm-auth=true`). 이는 프로덕션 보안을 위한 의도적인 설정입니다.

```bash
# auth 도구 활성화
swag2mcp mcp --disable-llm-auth=false
```

## 인증 문제

### 401 Unauthorized

자격 증명이 없거나 유효하지 않아 API가 요청을 거부했습니다.

```bash
# 인증이 설정되었는지 확인
swag2mcp info

# 설정 검증
swag2mcp validate

# 환경 변수가 설정되었는지 확인
echo $MY_TOKEN

# 토큰이 만료되지 않았는지 확인 (bearer 토큰은 정적)
```

**일반적인 원인:**
- 토큰이 없거나 비어 있음
- 환경 변수가 설정되지 않음
- 토큰이 만료됨 (bearer 토큰은 자동 갱신되지 않음)
- 잘못된 인증 유형 설정

### 403 Forbidden

권한이 부족하여 API가 요청을 거부했습니다.

- 토큰에 필요한 범위가 없을 수 있습니다
- API 키가 이 리소스에 접근 권한이 없을 수 있습니다
- 필요한 권한에 대한 API 문서를 확인하세요

### OAuth2 토큰 엔드포인트에 연결할 수 없음

swag2mcp가 OAuth2 토큰 URL에 연결할 수 없습니다.

```bash
# 설정의 token_url 확인
# URL이 올바르고 접근 가능한지 확인
curl -X POST https://auth.example.com/oauth/token \
  -d "grant_type=client_credentials" \
  -d "client_id=test" \
  -d "client_secret=test"

# 네트워크 연결 확인
# 회사 프록시 뒤에 있는 경우 프록시 설정 확인
```

### Digest 인증 실패

swag2mcp가 Digest 인증 핸드셰이크를 완료할 수 없습니다.

- 서버가 401 응답과 함께 `WWW-Authenticate: Digest ...` 헤더를 반환해야 합니다
- 챌린지는 5분 동안 캐시됩니다 — 서버가 nonce를 변경하면 캐시가 만료될 때까지 기다리세요
- 사용자 이름과 비밀번호가 올바른지 확인하세요

### HMAC 서명 불일치

API가 HMAC 서명된 요청을 거부했습니다.

- `api_key`와 `secret_key`가 올바른지 확인하세요
- API가 Binance 스타일 HMAC-SHA256 서명을 사용하는지 확인하세요
- 일부 거래소는 다른 서명 방식을 사용합니다 — HMAC 인증은 Binance 호환 API 전용입니다

### Script 인증 실패

외부 인증 스크립트가 실패했습니다.

```bash
# 스크립트가 존재하는지 확인
ls -la ~/.swag2mcp/auth_scripts/my-domain.sh

# 수동으로 스크립트 실행 테스트
sh ~/.swag2mcp/auth_scripts/my-domain.sh

# 스크립트 출력 형식 확인 (JSON이어야 함: {"token": "...", "expires_in": 3600})
# 스크립트가 30초 이내에 완료되는지 확인
# 스크립트에 실행 권한이 있는지 확인
chmod +x ~/.swag2mcp/auth_scripts/my-domain.sh
```

## 검색 문제

### 검색 결과 없음

검색 결과가 없습니다.

```bash
# spec이 로드되었는지 확인
swag2mcp ls

# spec이 비활성화되지 않았는지 확인
swag2mcp validate

# 더 간단한 쿼리 시도
# method로 검색: method:GET
# tag로 검색: tag:pets

# 인덱스는 MCP 서버 시작 시마다 재구축됩니다
# 방금 spec을 추가했다면 서버를 다시 시작하세요
```

### 검색 결과가 관련 없음

쿼리가 너무 광범위하거나 모호합니다.

- 필드 필터를 사용하여 좁히기: `method:GET +tag:pets`
- 정확한 구문 사용: `"find pet by status"`
- `limit` 매개변수를 사용하여 더 집중된 결과 얻기

## API 호출 문제

### invoke가 오류 반환

API 호출이 실패했습니다.

```bash
# 오류 메시지 확인 — HTTP 상태 코드 포함
# 4xx 오류: 매개변수, 인증 또는 권한 확인
# 5xx 오류: API 서버에 문제가 있음

# 호출 전에 항상 엔드포인트 검사
inspect(endpointId: "...")

# 모든 필수 매개변수가 제공되었는지 확인
# 매개변수 유형 확인 (문자열, 숫자, 부울)
```

### 속도 제한 오류

LLM이 동일한 엔드포인트를 너무 빨리 호출했습니다.

각 엔드포인트에는 10초의 쿨다운이 있습니다. 다시 호출하기 전에 기다리거나 속도 제한기를 비활성화하세요:

```yaml
disable_ratelimiter: true
```

### 응답이 너무 큼 (fileRef 반환)

응답이 `max_response_size`를 초과했습니다.

정상적인 현상입니다. 응답 도구를 사용하여 데이터를 탐색하세요:

```
1. response_outline(path) → 구조 이해
2. response_compress(path, mode: "first_of_array") → 샘플 얻기
3. response_slice(path, jsonPath: "data.0") → 특정 데이터 얻기
```

또는 제한을 늘리세요:

```yaml
http_client:
  max_response_size: 4194304  # 4 MB
```

### 느린 API 응답

API 응답 시간이 너무 깁니다.

```yaml
http_client:
  timeout: 120s  # 기본 30s에서 증가
```

## 워크스페이스 문제

### swag2mcp init 실패: "directory is not empty"

대상 디렉토리에 이미 파일이 있습니다.

```bash
# --force를 사용하여 덮어쓰기
swag2mcp init --force

# 또는 다른 디렉토리 사용
swag2mcp init ./new-workspace
```

### swag2mcp update 실패

하나 이상의 spec 파일을 다운로드할 수 없습니다.

```bash
# 어떤 URL이 실패했는지 오류 메시지 확인
# URL에 접근 가능한지 확인
curl -I <failed-url>

# 네트워크 연결 확인
# 프록시 설정 확인
```

### 내보내기에서 ZIP이 생성되지 않음

`[output]` 인수는 디렉토리가 아닌 `.zip`으로 끝나는 파일 경로여야 합니다.

```bash
# 올바름
swag2mcp export /path/to/workspace /path/to/backup.zip

# 잘못됨 (ZIP이 생성되지 않음)
swag2mcp export /path/to/workspace /some/directory
```

### 가져오기 실패: "not a valid swag2mcp backup"

ZIP 파일이 `swag2mcp export`로 생성되지 않았습니다.

`swag2mcp export`로 생성된 ZIP 아카이브만 가져올 수 있습니다. 아카이브에는 특정 내부 구조(`swag2mcp.yaml`, `specs/`, `auth_scripts/`)가 있습니다.

## TUI 문제

### TUI가 올바르게 렌더링되지 않음

터미널이 너무 작거나 필요한 기능을 지원하지 않습니다.

- 최소 터미널 크기: 80×24 문자
- TUI는 Bubbletea를 사용하며 대부분의 최신 터미널에서 작동합니다
- 터미널 창 크기를 조정해 보세요
- 다른 터미널 에뮬레이터를 시도해 보세요

### TUI에 "no specs found" 표시

워크스페이스에 설정된 spec이 없습니다.

```bash
# spec 확인
swag2mcp ls

# spec 추가
swag2mcp add spec
```

## 모의 서버 문제

### 모의 서버가 시작되지 않음

```bash
# 설정에 mock_enabled: true가 있는지 확인
# 모든 collection에 base_mock_url이 설정되어 있는지 확인
# 포트가 사용 중이 아닌지 확인
lsof -i :9090

# 모의 서버 로그 확인
swag2mcp-mock mockserver
```

### 모의 서버가 빈 응답 반환

명세 파일에 응답 스키마가 정의되지 않았을 수 있습니다.

- 모의 서버는 응답 스키마에서 데이터를 생성합니다
- 스키마가 없으면 `{}`를 반환합니다
- OpenAPI 명세에 `responses`와 `schema`가 정의되어 있는지 확인하세요

## 네트워크 문제

### 프록시 연결 실패

swag2mcp가 설정된 프록시를 통해 연결할 수 없습니다.

```bash
# 프록시 URL 형식 확인 (스키마 포함: http://, https://, socks5://)
# 프록시 자격 증명 확인
# 바이패스 목록 확인 — 대상이 바이패스 목록에 있을 수 있음
# curl로 프록시 테스트
curl -x http://proxy.company.com:8080 https://api.example.com
```

### TLS/SSL 오류

인증서 검증에 실패했습니다.

- MCP 서버에 자체 서명 인증서를 사용하는 경우 클라이언트가 이를 신뢰해야 합니다
- `--tls`와 함께 모의 서버를 사용하면 자체 서명 인증서가 자동으로 생성됩니다
- API 호출의 경우 swag2mcp는 시스템의 인증서 저장소를 사용합니다

## 기타 문제

### 높은 디스크 사용량

캐시 및 응답 디렉토리가 시간이 지남에 따라 커질 수 있습니다.

```bash
# 모든 항목 정리
swag2mcp clean

# 오래된 응답(48시간 초과)은 MCP 서버 시작 시 자동으로 정리됩니다
# 캐시 파일은 1-48시간 사이에 무작위로 만료됩니다
```

### "command not found" after go install

`go install` 디렉토리가 PATH에 없습니다.

```bash
# Go가 바이너리를 설치하는 위치 확인
go env GOPATH
# PATH에 추가
export PATH=$PATH:$(go env GOPATH)/bin
```

### LLM이 도구를 올바르게 사용하지 않음

LLM에 더 나은 지침이나 포맷팅 스킬이 필요할 수 있습니다.

- spec 설정에서 `llm_instruction`을 사용하여 API가 무엇을 하는지 설명하세요
- 일관된 출력 포맷팅을 위해 [swag2mcp-format 스킬](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md)을 고려하세요
- LLM 응답의 품질은 모델과 받는 지침에 따라 달라집니다

### 버그를 어떻게 신고하나요?

[GitHub](https://github.com/mmadfox/swag2mcp/issues)에 이슈를 열고 다음 정보를 포함하세요:
- swag2mcp 버전 (`swag2mcp --version`)
- 운영 체제 및 아키텍처
- 실행한 정확한 명령어
- 전체 오류 메시지
- 설정 파일 (시크릿 제거)
