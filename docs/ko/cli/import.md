# import

## 목적

명세 파일을 워크스페이스의 `specs/` 디렉토리에 로컬로 저장하거나 ZIP 백업에서 전체 워크스페이스를 복원합니다. 세 가지 모드가 다양한 시나리오를 다룹니다: 단일 명세 추가, 기존 설정에서 대량 가져오기, 또는 전체 워크스페이스 복원.

## 사용 시기

- 명세 URL이나 파일이 있고 워크스페이스에 로컬로 저장하려고 할 때
- 설정의 모든 collection 명세 파일을 다운로드하여 워크스페이스를 자체 포함형으로 만들고 싶을 때
- `export`로 생성된 ZIP 백업에서 워크스페이스를 복원해야 할 때
- swag2mcp를 다른 머신으로 마이그레이션할 때

## 구문

```bash
swag2mcp import [path] [source] [name] [flags]
```

## 인수

| 인수 | 위치 | 필수 | 설명 |
|------|------|------|------|
| `path` | 1 | 아니요 | 워크스페이스 디렉토리. 생략 시 경로 해결 규칙에 따라 결정됩니다. |
| `source` | 2 | 다양 | 명세 파일의 URL 또는 로컬 경로, 또는 ZIP 아카이브 경로 |
| `name` | 3 | 다양 | 저장할 파일 이름 (예: `example-api.yaml`). 생략 시 URL에서 자동 생성됩니다. |

## 플래그

| 플래그 | 약어 | 타입 | 기본값 | 설명 |
|-------|------|------|--------|------|
| `--spec` | `-s` | `string` | `""` | 설정에서 collection 명세 파일을 다운로드. 값 없이 모든 spec, 또는 `--spec meteo,github`처럼 도메인 지정 |
| `--force` | `-f` | `bool` | `false` | 기존 명세 파일을 오류 없이 덮어쓰기 |
| `--from-zip` | | `string` | `""` | swag2mcp 백업 ZIP에서 워크스페이스 복원 |

## 작동 방식

### 모드 1 — URL 또는 파일에서 단일 가져오기

명세 파일을 다운로드하여 `specs/`에 저장:

```bash
swag2mcp import https://example.com/spec.yaml example-api.yaml
swag2mcp import /path/to/workspace https://example.com/spec.yaml example-api.yaml
swag2mcp import ./local-spec.yaml example-api.yaml
```

`name`을 생략하면 URL의 파일 이름에서 자동 생성됩니다:
```bash
swag2mcp import https://example.com/specs/petstore.yaml
# → petstore.yaml로 저장
```

기존 파일을 `--force`로 덮어쓰기:
```bash
swag2mcp import --force https://example.com/spec.yaml example-api.yaml
```

가져오기 후 출력에는 워크스페이스 경로, 저장된 파일 및 `swag2mcp.yaml`에 추가할 YAML 템플릿이 표시됩니다:

```
✅ Imported to /path/to/workspace
   specs/example-api.yaml

   Add to swag2mcp.yaml:
     specs:
       - domain: <your-domain>
         collections:
           - location: specs/example-api.yaml
```

### 모드 2 — 기존 설정에서 대량 가져오기 (`--spec`)

지정된 도메인의 모든 collection 명세 파일을 설정된 `location` URL에서 다운로드하여 `specs/`에 저장하고 설정을 로컬 복사본을 가리키도록 업데이트:

```bash
swag2mcp import --spec                # 모든 spec
swag2mcp import --spec meteo           # 특정 spec
swag2mcp import --spec meteo,github    # 여러 spec
swag2mcp import /path/to/workspace --spec meteo
```

지정된 도메인이 설정에 존재하지 않으면 명령이 오류를 반환합니다:
```
Error: import_no_match
  Spec "nonexistent" not found in config.
```

이렇게 하면 워크스페이스가 자체 포함형이 됩니다 — 가져오기 후에는 원격 명세 URL이 필요하지 않습니다.

### 모드 3 — ZIP 백업에서 복원

`swag2mcp export`로 생성된 ZIP 아카이브에서 전체 워크스페이스 복원:

```bash
swag2mcp import --from-zip /path/to/backup.zip
swag2mcp import /path/to/workspace /path/to/backup.zip
```

> **ZIP은 `swag2mcp export`로 생성된 것이어야 합니다.** 임의의 ZIP 파일은 작동하지 않습니다 — 아카이브에는 특정 내부 구조(`swag2mcp.yaml`, `specs/`, `auth_scripts/`)가 있습니다.

## 명령 후 검증

```bash
# 단일 또는 대량 가져오기
swag2mcp ls [path]
# 새 spec이 목록에 나타나야 함

# ZIP 복원
swag2mcp ls [path]
# 백업의 모든 spec이 나타나야 함
```

## 세부 사항

- **대량 모드는 설정 필요:** `--spec`을 사용할 때 설정 파일이 존재해야 합니다. 필요하면 먼저 `init`을 실행하세요.
- **단일 가져오기는 워크스페이스 생성:** 워크스페이스가 없으면 자동으로 생성됩니다.
- **ZIP 감지:** `.zip`으로 끝나는 위치 인수는 ZIP 소스로 처리됩니다. `--from-zip` 플래그가 위치 감지보다 우선합니다.
- **HTTP 클라이언트:** 설정의 전역 HTTP 클라이언트 설정이 가져오기 중에 적용됩니다(타임아웃, 프록시, 헤더 등).
