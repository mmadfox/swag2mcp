# export

## 목적

워크스페이스의 이식 가능한 ZIP 백업을 생성합니다. 아카이브에는 설정 파일, 모든 명세 파일, 인증 스크립트가 포함됩니다 — 다른 머신에서 워크스페이스를 복원하는 데 필요한 모든 것입니다.

## 사용 시기

- 변경 전에 워크스페이스를 백업하려고 할 때
- swag2mcp를 다른 머신으로 마이그레이션할 때
- API 설정을 동료와 공유하려고 할 때
- 재현 가능한 환경을 준비할 때

## 구문

```bash
swag2mcp export [path] [output] [flags]
```

## 인수

| 인수 | 위치 | 필수 | 설명 |
|------|------|------|------|
| `path` | 1 | 아니요 | 워크스페이스 디렉토리. 생략 시 경로 해결 규칙에 따라 결정됩니다. |
| `output` | 2 | 아니요 | 출력 ZIP 파일의 전체 경로. 생략 시 `./swag2mcp-backup-&lt;timestamp&gt;.zip`으로 기본 설정됩니다. |

## 플래그

| 플래그 | 약어 | 타입 | 기본값 | 설명 |
|-------|------|------|--------|------|
| `--spec` | `-s` | `stringSlice` | `nil` | 지정된 spec만 내보내기 (쉼표로 구분) |

## 작동 방식

### 기본 내보내기

타임스탬프가 있는 이름으로 현재 디렉토리에 ZIP을 생성합니다:

```bash
swag2mcp export
# ./swag2mcp-backup-2026-07-22-143022.zip 생성
```

### 커스텀 출력 경로

```bash
swag2mcp export /path/to/workspace /path/to/backup.zip
```

### 특정 spec 내보내기

```bash
swag2mcp export --spec meteo
swag2mcp export --spec meteo,store
```

## ZIP 내용

| 항목 | 설명 |
|------|------|
| `swag2mcp.meta` | 내보내기에 대한 메타데이터 |
| `swag2mcp.yaml` | 설정 파일 |
| `specs/` | 모든 명세 파일 (OpenAPI/Swagger/Postman) |
| `auth_scripts/` | 인증 스크립트 |
| `cache/` | 비어 있음 (캐시는 내보내지 않음) |
| `responses/` | 비어 있음 (응답은 내보내지 않음) |

## 복원

백업에서 복원하려면 `import`를 사용하세요:

```bash
swag2mcp import --from-zip /path/to/backup.zip
```

## 명령 후 검증

항상 ZIP 파일이 생성되었는지 확인하세요:

```bash
ls -la swag2mcp-backup-*.zip
# 또는 커스텀 출력 경로:
ls -la /path/to/backup.zip
```

## 세부 사항

- **출력은 파일 경로여야 함:** `[output]` 인수는 `.zip`으로 끝나는 전체 파일 경로여야 합니다. 디렉토리를 **전달하지 마세요** — 디렉토리 경로가 주어지면 명령어가 ZIP을 생성하지 않습니다.
- **기본 파일 이름:** UTC 타임스탬프를 사용한 `swag2mcp-backup-<YYYY-MM-DD-HHMMSS>.zip`.
- **`--spec` 필터:** 설정되면 지정된 spec만 포함됩니다. 다른 spec은 아카이브에서 제외됩니다.
- **설정 불필요:** `export`는 유효한 설정 파일 없이도 작동합니다. 워크스페이스에 있는 모든 것을 내보냅니다.
- **캐시와 응답은 제외됨:** 복원 시 오래된 상태가 되는 일시적인 데이터입니다. 설정, 명세, 인증 스크립트만 보존됩니다.
