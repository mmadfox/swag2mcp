# 응답 크기 관리

## 개요

API 응답은 매우 클 수 있습니다 — 때로는 LLM의 컨텍스트 윈도우에 맞지 않을 정도로 큽니다. swag2mcp는 초과된 응답을 디스크에 저장하고 탐색 도구를 제공하여 자동으로 응답 크기를 관리합니다.

## 작동 방식

1. **`invoke` 호출** — swag2mcp가 API 요청을 실행합니다
2. **응답이 작은 경우** (제한 이내) — LLM에 인라인으로 반환됩니다
3. **응답이 너무 큰 경우** (제한 초과) — `{workspace}/responses/`에 JSON 파일로 저장됩니다. LLM은 전체 응답 대신 파일 참조를 받습니다

### 예시: 작은 응답 (인라인)

```json
{
  "statusCode": 200,
  "body": {
    "id": 1,
    "name": "Rex",
    "status": "available"
  }
}
```

### 예시: 큰 응답 (파일 참조)

```json
{
  "statusCode": 200,
  "fileRef": {
    "path": "/Users/user/.swag2mcp/responses/response_a1b2c3d4.json",
    "size": 1572864,
    "sizeHint": "1.5 MB",
    "maxSizeHint": "2 KB",
    "message": "응답이 2 KB 제한을 초과하여 디스크에 저장되었습니다.",
    "openCmd": "open /Users/user/.swag2mcp/responses/response_a1b2c3d4.json"
  }
}
```

## 설정

```yaml
http_client:
  max_response_size: 1048576  # 1 MB (바이트)
```

### max_response_size

- **타입:** `int` (바이트)
- **기본값:** `1048576` (1 MB)
- **범위:** 256 ~ 10,485,760 바이트 (10 MB)
- **효과:** 이 크기보다 큰 응답은 인라인으로 반환되지 않고 디스크에 저장됩니다
- **증가 시기:** 대용량 데이터 세트를 반환하는 API (보고서, 로그, 분석)
- **감소 시기:** 제한된 LLM 컨텍스트 윈도우 또는 모든 응답에 파일 기반 접근을 선호할 때

## 큰 응답 작업

`invoke`가 `fileRef`를 반환하면 다음 세 가지 도구를 사용하여 데이터를 탐색하세요:

### 1. response_outline — 구조 이해

응답의 구조적 요약을 얻습니다: 키, 타입, 배열 길이, 탐색 힌트.

```json
→ response_outline(path: "/path/to/file.json")
← {
    "type": "object",
    "size": 1572864,
    "keys": ["data", "meta"],
    "itemCount": 500,
    "compressionHints": [
      "response_compress(path, 'first_of_array', 'data')",
      "response_compress(path, 'sample_array', 'data', arrayHead=3, arrayTail=2)"
    ]
  }
```

### 2. response_compress — 더 작은 버전 얻기

데이터를 압축하여 인라인에 맞춥니다. 여러 압축 모드로 적절한 트레이드오프를 선택할 수 있습니다.

| 모드 | 설명 | 최적 용도 |
|------|------|----------|
| `first_of_array` | 배열의 첫 번째 요소만 유지 | 모든 요소가 동일한 구조일 때 |
| `sample_array` | 배열의 처음(3)과 끝(2) 유지 | 값의 범위를 확인해야 할 때 |
| `truncate_strings` | 모든 문자열을 N자로 단축 | 문자열이 매우 길 때 |
| `keys_only` | 값을 타입 이름으로 대체 | 구조만 필요할 때 |
| `select_keys` | 지정된 키만 유지 | 특정 필드가 필요할 때 |

```json
→ response_compress(path: "/path/to/file.json", mode: "first_of_array", jsonPath: "data")
← {
    "body": [{ "id": 1, "name": "Rex" }],
    "hint": "first_of_array 모드를 사용하여 배열을 500개에서 1개 항목으로 압축했습니다"
  }
```

### 3. response_slice — 특정 조각 추출

JSON 경로 또는 줄 범위로 특정 요소나 값을 가져옵니다.

```json
→ response_slice(path: "/path/to/file.json", jsonPath: "data.0")
← {
    "slice": {
      "value": { "id": 1, "name": "Rex" },
      "jsonPath": "data.0",
      "nextPath": "data.1",
      "prevPath": null
    }
  }
```

## 전체 워크플로우

```
1. invoke(endpoint) → fileRef (응답 1.5 MB)
2. response_outline(path) → 구조: { data: Array(500) }
3. response_compress(path, mode: "first_of_array", jsonPath: "data") → 첫 번째 항목
4. response_slice(path, jsonPath: "data.0") → 첫 번째 항목 상세
5. response_slice(path, jsonPath: "data.1") → 두 번째 항목
```

## 자동 정리

MCP 서버가 시작될 때(`swag2mcp mcp`) 48시간보다 오래된 응답 파일이 자동으로 제거됩니다. 수동으로 정리할 수도 있습니다:

```bash
swag2mcp clean
```

## 중요 참고 사항

- **제한은 바이트 단위** — `1048576` = 1 MB, `2097152` = 2 MB 등
- **파일 참조에는 열기 명령이 포함됨** — macOS에서는 `open`, Linux에서는 `xdg-open`
- **응답 파일은 무작위 접미사로 이름이 지정됨** — 동시 호출 간 충돌 없음
- **응답 디렉토리는 자동으로 생성됨** — 수동 설정 불필요
