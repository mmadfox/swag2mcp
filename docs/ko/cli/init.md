# init

## 목적

`init` 명령어는 **워크스페이스** — `swag2mcp.yaml` 설정 파일과 캐시, 명세, 응답, 인증 스크립트용 하위 디렉토리가 있는 디렉토리를 생성합니다. swag2mcp를 설정할 때 가장 먼저 실행하는 명령어입니다.

## 사용 시기

- swag2mcp를 처음 설정할 때
- 특정 디렉토리에 새 워크스페이스를 생성하려고 할 때
- 손상되었거나 누락된 워크스페이스를 다시 초기화해야 할 때

## 구문

```bash
swag2mcp init [path] [flags]
```

## 인수

| 인수 | 위치 | 필수 | 설명 |
|------|------|------|------|
| `path` | 1 | 아니요 | 워크스페이스 디렉토리. 생략 시 `~/.swag2mcp`로 기본 설정됩니다. |

## 플래그

| 플래그 | 약어 | 타입 | 기본값 | 설명 |
|-------|------|------|--------|------|
| `--interactive` | `-i` | `bool` | `false` | 대화형 TUI 마법사 실행 |
| `--force` | `-f` | `bool` | `false` | 비어 있지 않은 디렉토리의 기존 설정 덮어쓰기 |

## 작동 방식

### 비대화형 모드 (기본값)

spec이 없는 최소 `swag2mcp.yaml`을 생성합니다. 이후 수동으로 파일을 편집합니다.

```bash
swag2mcp init
# ~/.swag2mcp/swag2mcp.yaml 생성

swag2mcp init ./my-project
# ./my-project/swag2mcp.yaml 생성

swag2mcp init /absolute/path
# /absolute/path/swag2mcp.yaml 생성
```

### 대화형 모드 (`-i`)

18단계 TUI 마법사를 실행하여 다음을 안내합니다:

1. 워크스페이스 디렉토리 선택
2. 도메인, 제목, base URL로 spec 추가
3. 위치 URL로 collection 설정
4. 인증 설정 (9가지 방법 모두)
5. HTTP 클라이언트 설정 (타임아웃, 프록시, 헤더 등)

```bash
swag2mcp init -i
```

### 강제 모드 (`--force`)

기본적으로 `init`은 비어 있지 않은 디렉토리에서 실행을 거부합니다. `--force`를 사용하여 덮어쓰세요:

```bash
swag2mcp init -f
swag2mcp init ./existing-dir -f
```

## 생성되는 항목

```
~/.swag2mcp/
├── swag2mcp.yaml       # 설정 파일
├── cache/               # 다운로드된 원격 명세 파일
├── specs/               # 로컬 명세 파일
├── responses/           # 저장된 API 호출 응답
└── auth_scripts/        # 인증 스크립트 (ScriptAuth 유형용)
```

## 명령 후 검증

```bash
ls ~/.swag2mcp/swag2mcp.yaml
# 파일이 존재하면 init 성공
```

## 세부 사항

- **경로 해결:** `[path]`는 **워크스페이스 디렉토리**이지 파일 경로가 아닙니다. CLI가 자동으로 `swag2mcp.yaml`을 추가합니다. 해결 순서: 명시적 `[path]` → 현재 디렉토리(`./`) → `~/.swag2mcp/`.
- **비어 있지 않은 디렉토리 확인:** `--force` 없이 `init`은 대상 디렉토리가 존재하고 비어 있지 않으면 오류를 반환합니다. 이는 실수로 인한 덮어쓰기를 방지합니다.
- **인증 스크립트 스텁:** spec이 `ScriptAuth`를 사용하면 `init`이 `auth_scripts/`에 스텁 스크립트 파일(Unix는 `.sh`, Windows는 `.bat`)을 생성합니다.
- **출력:** 성공 시 설정 경로와 힌트를 출력합니다: `"다음 단계: swag2mcp.yaml을 편집하거나 'swag2mcp ls'를 실행하여 설정된 spec을 확인하세요"`.
