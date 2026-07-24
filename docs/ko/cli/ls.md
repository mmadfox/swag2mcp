# ls

## 목적

설정된 모든 **spec**과 해당 **collection**을 사람이 읽을 수 있는 형식으로 나열합니다. 워크스페이스에서 사용 가능한 API를 확인하는 기본 방법입니다.

## 사용 시기

- 어떤 API가 설정되어 있는지 확인하려고 할 때
- spec 또는 collection ID를 찾아야 할 때
- 각 collection에 몇 개의 엔드포인트가 있는지 확인하려고 할 때
- 태그로 spec을 필터링하려고 할 때

## 구문

```bash
swag2mcp ls [path] [flags]
```

## 인수

| 인수 | 위치 | 필수 | 설명 |
|------|------|------|------|
| `path` | 1 | 아니요 | 워크스페이스 디렉토리. 생략 시 경로 해결 규칙에 따라 결정됩니다. |

## 플래그

| 플래그 | 약어 | 타입 | 기본값 | 설명 |
|-------|------|------|--------|------|
| `--tags` | `-t` | `string` | `""` | 태그로 spec 필터링 (쉼표로 구분) |

## 작동 방식

### 모든 spec 나열

도메인, collection, 엔드포인트 수와 함께 모든 spec을 표시:

```bash
swag2mcp ls
```

출력 예시:

```
Specifications:
  dadjoke (https://icanhazdadjoke.com)
    jokes (3 endpoints)
  meteo (https://meteo.swagger.io/v2)
    forecast (5 endpoints)
    current (8 endpoints)
  binance (https://api.binance.com)
    market-data (12 endpoints)
```

### 태그로 필터링

지정된 태그가 있는 spec만 표시:

```bash
swag2mcp ls --tags=public
swag2mcp ls --tags=public,internal
```

## 명령 후 검증

`add`, `delete`, `update` 또는 `import` 후 `ls`를 사용하여 워크스페이스 상태가 예상과 일치하는지 확인하세요.

## 세부 사항

- **자동 초기화:** 설정 파일이 없으면 `ls`가 자동으로 init 마법사를 먼저 실행합니다.
- **태그 필터링:** 태그는 쉼표로 구분됩니다. 지정된 태그 중 **하나라도** 일치하는 spec이 표시됩니다(OR 로직).
- **출력 형식:** 출력은 일반 텍스트이며 JSON이 아닙니다. 머신 판독 가능 출력은 `info`를 사용하세요.
