# 워크스페이스

워크스페이스는 swag2mcp가 모든 데이터(설정, 캐시된 명세, 로컬 명세 파일, 저장된 응답, 인증 스크립트)를 저장하는 디렉토리입니다.

## 구조

```
~/.swag2mcp/                          # 워크스페이스 루트 (기본값)
├── swag2mcp.yaml                     # 설정 파일
├── cache/                            # 캐시된 원격 명세 파일
│   ├── a1b2c3d4e5f6...spec          # 캐시된 명세 내용
│   └── a1b2c3d4e5f6...meta          # 캐시 메타데이터 (JSON)
├── specs/                            # 로컬 명세 파일
│   └── my-api.yaml
├── responses/                        # 저장된 API 응답 (큰 응답)
│   ├── meteo-get-forecast-abc123.json
│   └── response-fragment-def456.json
└── auth_scripts/                     # 인증 스크립트
    ├── meteo.sh                      # Unix 셸 스크립트
    └── meteo.bat                     # Windows 배치 스크립트
```

## 기본 경로

- **Linux/macOS**: `~/.swag2mcp/`
- **Windows**: `%USERPROFILE%\.swag2mcp\`

## 커스텀 경로

```bash
swag2mcp mcp /path/to/workspace
swag2mcp mcp ./my-workspace
```

## 디렉토리

### cache/

다운로드된 원격 명세 파일을 저장합니다. 각 파일은 URL의 SHA-256 해시를 파일 이름으로 사용하여 캐시됩니다:

- `{hash}.spec` — 캐시된 명세 파일 내용
- `{hash}.meta` — JSON 메타데이터 (소스 URL, 캐시 시간, TTL)

각 캐시 파일은 1시간에서 48시간 사이의 무작위 TTL을 가집니다. 캐시는 매 시작 시 자동으로 확인됩니다 — 유효한(만료되지 않은) 항목이 있으면 다운로드 없이 재사용됩니다.

**명령어:**
- `swag2mcp update` — 캐시 지우기 및 모든 명세 다시 다운로드
- `swag2mcp clean` — 캐시 및 응답 지우기

### specs/

collection이 `location: specs/{name}`을 통해 가리키는 로컬 명세 파일을 저장합니다. 여기의 파일은 캐싱 없이 직접 사용됩니다.

이 디렉토리는 다음으로 채워집니다:
- `swag2mcp import <source> <name>` — 원격 명세를 다운로드하여 여기에 저장
- `swag2mcp export` — 명세를 내보내기 ZIP에 복사
- 수동 배치 — 직접 명세 파일을 여기에 복사 가능

### responses/

`max_response_size` 제한(기본값 1MB)을 초과하는 API 응답을 저장합니다. LLM이 엔드포인트를 호출하고 응답이 너무 크면 swag2mcp가 여기에 저장하고 파일 참조를 반환합니다.

이름 규칙: `{domain}-{method}-{path_with_underscores}-{6char_hex}.json`

오래된 응답은 MCP 서버 시작 후 48시간이 지나면 자동으로 정리됩니다.

### auth_scripts/

`script` 인증 유형의 인증 스크립트를 저장합니다. 각 스크립트는 spec의 도메인 이름을 따서 명명됩니다.

#### 이름 규칙

| 플랫폼 | 파일 이름 | 예시 |
|--------|---------|------|
| Unix (Linux, macOS) | `{domain}.sh` | `meteo.sh` |
| Windows | `{domain}.bat` | `meteo.bat` |

도메인에는 `/` 또는 `\` 문자를 포함할 수 없습니다.

#### 스크립트 작동 방식

1. swag2mcp가 30초 타임아웃으로 스크립트를 실행합니다
2. 스크립트는 stdout에 유효한 JSON을 출력해야 합니다
3. swag2mcp가 JSON을 파싱하고 API 요청에 토큰을 사용합니다

#### 예상 출력 형식

```json
{
  "token": "your-token-here",
  "expires_in": 3600
}
```

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| `token` | string | ✅ | 인증 토큰 |
| `access_token` | string | ❌ | `token`의 대안 (먼저 확인됨) |
| `token_type` | string | ❌ | 토큰 유형 (예: "Bearer") |
| `expires_in` | number | ❌ | 토큰 수명(초) (기본값: 3600) |

#### 실행

| 플랫폼 | 명령어 |
|--------|--------|
| Unix | `sh {domain}.sh` |
| Windows | `cmd /c {domain}.bat` |

#### 토큰 캐싱

토큰은 만료될 때까지 메모리에 캐시됩니다. 각 API 호출 시 swag2mcp가 먼저 캐시를 확인합니다 — 스크립트는 캐시된 토큰이 만료된 경우에만 실행됩니다.

#### 스텁 생성

`auth: { type: script, config: { domain: "myapi" } }`를 설정하면 swag2mcp가 자동으로 스텁 스크립트를 생성합니다:

**Unix (`auth_scripts/myapi.sh`):**
```bash
#!/bin/sh
echo '{"token": "your-token-here", "expires_in": 3600}'
```

**Windows (`auth_scripts/myapi.bat`):**
```bat
@echo off
echo {"token": "your-token-here", "expires_in": 3600}
```

플레이스홀더 토큰을 실제 인증 로직으로 교체하세요.

#### 고아 정리

spec을 삭제하면 해당 인증 스크립트가 고아가 됩니다. swag2mcp는 다음 경우에 고아 스크립트를 자동으로 제거합니다:
- `swag2mcp update`
- `swag2mcp clean`

## 명령어

### update

```bash
swag2mcp update [path]
```

설정을 검증하고, 캐시와 응답을 지운 후 모든 명세 파일을 다시 다운로드합니다. 또한 인증 스크립트가 존재하는지 확인하고 고아 스크립트를 제거합니다.

다음 후에 이 명령어를 사용하세요:
- collection 추가 또는 제거
- collection 위치 변경
- 재캐싱이 필요한 명세 파일 편집

### clean

```bash
swag2mcp clean [path]
```

`cache/`와 `responses/`의 모든 내용과 고아 인증 스크립트를 제거합니다. 명세를 재캐싱하지는 않습니다 — 그 용도는 `update`를 사용하세요.

### validate

```bash
swag2mcp validate [path]
```

모든 collection 위치를 포함한 설정을 검증합니다. [CLI: validate](../cli/validate.md) 참조.

## 내보내기 및 가져오기

```bash
# 워크스페이스를 ZIP으로 내보내기 (기본 이름: swag2mcp-backup-{date}.zip)
swag2mcp export

# 특정 경로로 내보내기
swag2mcp export /path/to/workspace /path/to/backup.zip

# 특정 spec만 내보내기
swag2mcp export --spec meteo

# 백업에서 복원
swag2mcp import --from-zip /path/to/backup.zip
swag2mcp import /path/to/workspace /path/to/backup.zip
```

내보내기 포함: `swag2mcp.yaml`, `specs/`, `auth_scripts/`. 캐시와 응답은 제외됩니다(로컬 데이터).

## .gitignore

워크스페이스가 Git 저장소 내에 있는 경우 `.gitignore`에 다음 항목을 추가하세요:

```gitignore
# swag2mcp — 로컬 데이터만
.swag2mcp/cache/
.swag2mcp/responses/
```

`cache/`와 `responses/` 디렉토리에는 커밋해서는 안 되는 로컬, 머신별 데이터가 포함됩니다. 그 외 모든 것(`swag2mcp.yaml`, `specs/`, `auth_scripts/`)은 설정이 팀 전체에 공유되도록 저장소에 있어야 합니다.
