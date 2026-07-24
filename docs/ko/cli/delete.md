# delete

## 목적

설정에서 **spec**(API 서비스) 또는 **collection**(명세 파일)을 제거합니다. `add`의 반대 작업입니다.

## 사용 시기

- API가 더 이상 필요하지 않을 때
- spec에서 특정 명세 파일을 제거하려고 할 때
- 워크스페이스를 정리할 때

## 구문

```bash
swag2mcp delete spec [path]
swag2mcp delete collection [path]
```

## 인수

| 인수 | 위치 | 필수 | 설명 |
|------|------|------|------|
| `path` | 1 | 아니요 | 워크스페이스 디렉토리. 생략 시 경로 해결 규칙에 따라 결정됩니다. |

## 플래그

없음. 두 하위 명령어 모두 순수하게 대화형입니다.

## 작동 방식

### Spec 삭제

목록에서 spec을 선택하라는 프롬프트를 표시한 후 삭제 전 확인을 요청합니다.

```bash
swag2mcp delete spec
```

### Collection 삭제

spec을 선택한 후 해당 spec 내의 collection을 선택하라는 프롬프트를 표시한 후 확인을 요청합니다.

```bash
swag2mcp delete collection
```

## ID 찾기

대화형 프롬프트는 ID 대신 사람이 읽을 수 있는 이름을 표시합니다. 참조용 ID가 필요하면:

```bash
# 모든 spec을 ID와 함께 나열
swag2mcp ls

# 특정 spec의 collection 나열
swag2mcp ls --tags
```

## 명령 후 검증

```bash
swag2mcp ls [path]
# 삭제된 spec 또는 collection이 더 이상 표시되지 않아야 함
```

## 세부 사항

- **TTY 필요:** 두 명령어 모두 대화형 터미널이 필요합니다. CI/CD 파이프라인, cron 작업 또는 비대화형 스크립트에서는 **작동하지 않습니다**.
- **`--force` 또는 `--yes` 없음:** 확인 프롬프트를 건너뛸 방법이 없습니다. 이는 실수로 인한 삭제를 방지하기 위한 의도적인 설계입니다.
- **자동 초기화:** 설정 파일이 없으면 `delete`가 자동으로 init 마법사를 먼저 실행합니다.
- **YAML 모드 없음:** `add`와 달리 `--yaml` 플래그가 없습니다. 삭제는 항상 대화형입니다.
