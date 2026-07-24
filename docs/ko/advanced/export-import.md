# 내보내기 및 가져오기

## 개요

swag2mcp는 ZIP 아카이브를 통한 전체 워크스페이스 왕복을 지원합니다. 전체 워크스페이스(설정, 명세 파일, 인증 스크립트)를 ZIP 파일로 내보내고 다른 머신에서 복원할 수 있습니다.

## 내보내기

워크스페이스의 이식 가능한 ZIP 백업을 생성합니다.

```bash
# 기본 파일로 내보내기 (swag2mcp-backup-<timestamp>.zip)
swag2mcp export

# 커스텀 경로로 내보내기
swag2mcp export --output ~/backups/swag2mcp-backup.zip

# 특정 spec만 내보내기
swag2mcp export --spec meteo
swag2mcp export --spec meteo,store
```

### 내보내기에 포함되는 항목

| 항목 | 설명 |
|------|------|
| `swag2mcp.yaml` | 설정 파일 |
| `specs/` | 모든 명세 파일 (OpenAPI/Swagger/Postman) |
| `auth_scripts/` | 인증 스크립트 |
| `swag2mcp.meta` | 메타데이터 (호환성 버전 정보) |

캐시와 응답은 **내보내지지 않습니다** — 일시적인 데이터이며 복원 시 오래된 상태가 됩니다.

### 기본 파일 이름

출력 경로를 지정하지 않으면 현재 디렉토리에 `swag2mcp-backup-<YYYY-MM-DD-HHMMSS>.zip`으로 저장됩니다(UTC 타임스탬프).

## 가져오기

ZIP 백업에서 워크스페이스를 복원하거나 명세 파일을 가져옵니다.

### ZIP에서 복원

```bash
# 전체 워크스페이스 복원
swag2mcp import --from-zip /path/to/backup.zip

# 덮어쓰기로 복원
swag2mcp import --from-zip /path/to/backup.zip -f
```

ZIP은 `swag2mcp export`로 생성된 것이어야 합니다 — 임의의 ZIP 파일은 작동하지 않습니다.

### 단일 명세 파일 가져오기

명세 파일을 다운로드하여 워크스페이스에 추가:

```bash
swag2mcp import https://example.com/spec.yaml myspec
swag2mcp import /path/to/workspace https://example.com/spec.yaml myspec
```

### 기존 설정에서 대량 가져오기

지정된 spec(도메인)의 모든 collection 명세 파일을 다운로드:

```bash
swag2mcp import --spec meteo
swag2mcp import /path/to/workspace --spec meteo,store
```

각 collection의 명세 파일을 다운로드하여 `specs/`에 저장하고 설정을 로컬 복사본을 가리키도록 업데이트합니다.

## 사용 사례

### 백업

```bash
swag2mcp export --output swag2mcp-$(date +%Y-%m-%d).zip
```

### 다른 머신으로 전송

```bash
# 이전 머신에서
swag2mcp export --output swag2mcp.zip

# ZIP을 새 머신으로 복사한 후:
swag2mcp import --from-zip swag2mcp.zip
```

### 설정 공유

```bash
swag2mcp init
swag2mcp export --output template.zip
# template.zip을 동료와 공유
```

## 내보내기 후 검증

항상 ZIP 파일이 생성되었는지 확인하세요:

```bash
ls -la swag2mcp-backup-*.zip
```

## 중요 참고 사항

- **출력은 `.zip`으로 끝나는 파일 경로여야 함** — 디렉토리를 전달하지 마세요
- **캐시와 응답은 제외됨** — 설정, 명세, 인증 스크립트만 보존됩니다
- **ZIP은 자체 포함됨** — swag2mcp가 설치된 모든 머신에서 복원할 수 있습니다
- **Spec 필터** — `--spec`을 사용하여 특정 spec만 내보내거나 가져올 수 있습니다
