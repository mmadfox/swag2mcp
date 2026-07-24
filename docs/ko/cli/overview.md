# CLI 명령어

## 개요

`swag2mcp` CLI는 워크스페이스 초기화 및 API 명세 관리부터 LLM 통합을 위한 MCP 서버 시작까지 모든 작업의 단일 진입점입니다. OpenAPI/Swagger/Postman 명세 작업의 전체 라이프사이클을 다루는 **13개의 명령어**를 제공합니다.

### CLI가 해결하는 문제

- **워크스페이스 라이프사이클** — 생성(`init`), 검사(`info`, `ls`), 정리(`clean`), 업데이트(`update`), 제거(`delete`)
- **Spec 및 Collection 관리** — 추가(`add`), 나열(`ls`), 삭제(`delete`)
- **실행 모드** — LLM 도구 접근을 위한 MCP 서버 시작(`mcp`) 또는 대화형 TUI 탐색기 실행(`run`)
- **진단** — 설정 검증(`validate`), 버전 표시(`version`), 런타임 정보 표시(`info`)
- **백업 및 복원** — ZIP을 통한 전체 워크스페이스 왕복(`export`, `import`)

### 주요 세부 사항

- **경로 해결** — `[path]`를 받는 명령어는 **워크스페이스 디렉토리**(파일 경로가 아님)를 기대합니다. 해결 순서: 명시적 `[path]` → 현재 디렉토리(`./`) → `~/.swag2mcp/`. CLI가 자동으로 `swag2mcp.yaml`을 추가합니다. 서비스로 실행하거나 IDE 설정에서 잘못된 워크스페이스 로드를 방지하기 위해 항상 명시적 경로를 전달하세요.
- **Spec vs Collection** — **spec**은 논리적 API 서비스(예: "Open-Meteo API")를 나타내고, **collection**은 하나의 OpenAPI/Swagger/Postman 파일입니다. 하나의 spec은 여러 collection을 가질 수 있습니다.
- **`--version`** 은 플래그(`swag2mcp --version`)와 하위 명령어(`swag2mcp version`) 모두로 지원됩니다.
- **`add spec` / `add collection`** 은 `--yaml`(인라인 문자열 또는 stdin용 `-`)을 통해 YAML 입력을 받습니다. 파일 또는 heredoc에서 파이핑하면 특수 문자로 인한 셸 따옴표 문제를 피할 수 있습니다.
- **`delete`** 는 TTY(대화형 터미널)가 필요합니다. `--force` 또는 `--yes` 플래그가 없습니다 — 항상 선택 및 확인을 위한 프롬프트를 표시합니다.
- **`mcp`** 는 LLM 통합을 위한 기본 명령어입니다. 세 가지 전송 방식을 지원합니다: `stdio`(기본값), `sse`, `streamable-http`. `--disable-llm-auth` 플래그(기본값: `true`)는 MCP 도구 목록에서 `auth` 도구를 제거하여 LLM이 토큰을 보거나 요청하지 못하게 합니다. 인증은 여전히 작동합니다 — 토큰은 LLM을 통하지 않고 표준 설정 메커니즘을 통해 획득됩니다. 이 모드는 **프로덕션**에 권장됩니다(LLM이 자격 증명에 접근할 수 없음). **디버깅** 또는 단기 토큰 사용 시 `--disable-llm-auth=false`로 설정하여 LLM이 `auth` 도구를 통해 새 토큰을 요청하도록 허용하세요.
- **`validate`** 는 YAML 구문, 설정 구조, 명세 파일 존재 여부, URL 접근성, 명세 형식(OpenAPI/Swagger/Postman), 인증 설정, HTTP 클라이언트 정확성을 확인합니다. 인증 엔드포인트나 API 엔드포인트 가용성은 **테스트하지 않습니다**.
- **`export` / `import`** 는 전체 워크스페이스 왕복을 제공합니다 — 설정 파일, 명세 파일, 캐시, 인증 스크립트가 모두 ZIP 아카이브에 포함됩니다.
- **`clean`** 은 `cache/`와 `responses/` 디렉토리를 제거하지만 `specs/`와 `auth_scripts/`는 보존합니다. 오래된 응답(48시간 초과)은 `mcp` 시작 시 자동으로 정리됩니다.

## 명령어

| 명령어 | 설명 |
|--------|------|
| [`init`](/cli/init) | 기본 설정으로 워크스페이스 디렉토리 초기화 |
| [`add`](/cli/add) | 설정에 spec 또는 collection 추가 |
| [`delete`](/cli/delete) | 대화형으로 spec 또는 collection 삭제 |
| [`ls`](/cli/ls) | 모든 spec과 해당 collection 나열 |
| [`run`](/cli/run) | 대화형 TUI API 탐색기 실행 |
| [`validate`](/cli/validate) | 설정 및 명세 파일 검증 |
| [`clean`](/cli/clean) | 캐시된 명세와 호출 응답 지우기 |
| [`update`](/cli/update) | 모든 spec 재검증, 재캐싱, 재인덱싱 |
| [`mcp`](/cli/mcp) | LLM 도구 접근을 위한 MCP 서버 시작 |
| [`version`](/cli/version) | swag2mcp 버전 출력 |
| [`info`](/cli/info) | 상세 설정 및 런타임 정보 표시 |
| [`import`](/cli/import) | 명세 파일 가져오기 또는 ZIP에서 워크스페이스 복원 |
| [`export`](/cli/export) | 워크스페이스를 이식 가능한 ZIP 백업으로 내보내기 |

## 전역 플래그

| 플래그 | 설명 |
|-------|------|
| `--version` | 버전 표시 (`version` 하위 명령어와 동일) |
| `--help` | 모든 명령어에 대한 도움말 표시 |
