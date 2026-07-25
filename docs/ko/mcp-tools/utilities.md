# 유틸리티 도구

유틸리티 도구는 지원 기능을 제공합니다: 인증 토큰 검색, 런타임 정보 가져오기, 인라인에 맞지 않는 큰 API 응답 작업.

---

## auth

### 목적

특정 spec의 인증 토큰, 헤더 또는 쿼리 매개변수를 검색합니다. LLM이 swag2mcp 외부에서 사용할 수 있는 자격 증명에 접근할 수 있게 합니다(예: curl 명령어 생성).

### 사용 시기

- 사용자가 명시적으로 원시 토큰이나 자격 증명을 요청할 때만
- 인증이 필요한 curl 명령어나 코드 스니펫을 생성할 때
- 사용자가 어떤 인증 방법이 설정되어 있는지 보려고 할 때

### 사용하지 말아야 할 때

- `inspect` 또는 `invoke` 전에 `auth`를 **호출하지 마세요** — `invoke`가 자동으로 인증을 획득하고 적용합니다
- 인증이 설정되어 있는지 확인하기 위해 `auth`를 **호출하지 마세요** — 대신 `info`를 사용하세요

### 작동 방식

spec의 인증 설정을 조회하고 인증 흐름(토큰 교환, 스크립트 실행 등)을 실행하여 현재 자격 증명을 획득합니다.

### 매개변수

| 매개변수 | 타입 | 필수 | 설명 |
|---------|------|------|------|
| `specId` | string | 예 | spec의 32자 MD5 해시 |

### 응답

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "headers": {
    "Authorization": "Bearer eyJhbGciOiJIUzI1NiIs...",
    "X-API-Key": "my-api-key"
  },
  "queryParams": {
    "api_key": "my-api-key"
  }
}
```

| 필드 | 타입 | 설명 |
|------|------|------|
| `token` | string | 원시 토큰 값 (bearer 토큰, API 키 등) |
| `headers` | object | 요청에 포함할 HTTP 헤더 |
| `queryParams` | object | 요청에 포함할 쿼리 매개변수 |

### 세부 사항

- **프로덕션에서 기본적으로 비활성화:** `--disable-llm-auth` 플래그(기본값: `true`)는 MCP 도구 목록에서 `auth` 도구를 완전히 제거합니다. LLM은 토큰을 보거나 요청할 수 없습니다. 디버깅 또는 단기 토큰을 위해 `--disable-llm-auth=false`로 설정하여 활성화하세요.
- **`invoke`가 인증을 자동 처리:** `invoke` 전에 `auth`를 호출할 필요가 없습니다. invoke 서비스가 자동으로 올바른 인증을 획득하고 적용합니다.
- **9가지 인증 방법 지원:** `none`, `basic`, `bearer`, `digest`, `hmac`, `oauth2-cc`(클라이언트 자격 증명), `oauth2-pwd`(비밀번호), `api-key`, `script`.
- 인증 방법이 실패하면 `auth_error` 반환 (예: OAuth2 토큰 엔드포인트에 연결할 수 없음, 스크립트 실행 실패).

---

## info

### 목적

swag2mcp 런타임의 포괄적인 요약을 반환합니다: 버전, 워크스페이스 경로, 활성 spec, HTTP 클라이언트 설정, MCP 전송 설정, 인증 방법, 모의 모드 상태.

### 사용 시기

- 사용자가 시스템 설정에 대해 물을 때
- 런타임 설정(타임아웃, 응답 크기 제한, 전송)을 확인해야 할 때
- 어떤 인증 방법을 사용할 수 있는지 알아야 할 때
- 설정 문제를 해결할 때

### 작동 방식

런타임 상태의 사전 계산된 스냅샷을 반환합니다. 매개변수가 필요하지 않습니다.

### 매개변수

없음.

### 응답

```json
{
  "version": "v1.2.0",
  "workspace": "~/.swag2mcp",
  "uptime": "2h 15m",
  "specs": {
    "total": 4,
    "active": 3,
    "disabled": 1,
    "collections": 6,
    "endpoints": 42
  },
  "http_client": {
    "timeout": "30s",
    "max_response_size": "2 KB",
    "follow_redirects": true,
    "max_redirects": 10,
    "randomize": false,
    "proxy": null,
    "headers": {},
    "cookies": []
  },
  "mcp": {
    "transport": "stdio",
    "addr": ":8080",
    "path": "/mcp",
    "auth_enabled": false
  },
  "auth": {
    "methods": ["bearer", "api-key"]
  },
  "mock": {
    "enabled": false
  }
}
```

| 필드 | 타입 | 설명 |
|------|------|------|
| `version` | string | swag2mcp 버전 |
| `workspace` | string | 워크스페이스 디렉토리 경로 |
| `uptime` | string | 서버 가동 시간 (사람이 읽을 수 있는 형식) |
| `specs` | object | Spec 요약: total, active, disabled, collections, endpoints |
| `http_client` | object | HTTP 클라이언트 설정 |
| `http_client.max_response_size` | string | 사람이 읽을 수 있는 형식의 최대 응답 크기 (예: "2 KB") |
| `mcp` | object | MCP 서버 설정 |
| `auth` | object | 사용 가능한 인증 방법 |
| `mock` | object | 모의 서버 상태 |

### 세부 사항

- `max_response_size`는 사람이 읽을 수 있는 형식으로 표시됨 (예: `"1 KB"`, `"2 MB"`)
- `uptime`은 서버 시작 시간에서 계산됨
- 데이터는 부트스트랩 시점에 캡처된 스냅샷 — MCP 서버가 시작된 시점의 상태를 반영

---

## response_outline

### 목적

`invoke`에 의해 디스크에 저장된 큰 JSON 응답 파일의 상위 수준 구조적 요약을 가져옵니다. 실제 값을 반환하지 않고 데이터의 형태(키, 타입, 배열 길이, 탐색 힌트)를 반환합니다.

### 사용 시기

- `invoke`가 `fileRef`를 반환한 직후 (응답이 인라인에 비해 너무 큼)
- 큰 응답 워크플로우의 **필수 첫 단계**

### 작동 방식

저장된 응답 파일을 읽고 구조를 분석합니다: 최상위 타입, 키, 배열 길이, 중첩 깊이, 압축 힌트.

### 매개변수

| 매개변수 | 타입 | 필수 | 설명 |
|---------|------|------|------|
| `path` | string | 예 | `fileRef.path`의 절대 경로 |
| `maxDepth` | int | 아니요 | 최대 재귀 깊이 (기본값: 3) |
| `maxArrayItems` | int | 아니요 | 검사할 배열 항목 수 (기본값: 5) |

### 응답

```json
{
  "outline": {
    "type": "object",
    "size": 1572864,
    "lineCount": 12500,
    "depth": 3,
    "structure": {
      "type": "object",
      "keys": ["data", "meta", "error"],
      "data": {
        "type": "array",
        "length": 500,
        "items": {
          "type": "object",
          "keys": ["id", "name", "status", "createdAt"]
        }
      }
    },
    "schemaHint": "3개 키가 있는 객체: data (array[500]), meta (object), error (null)",
    "keys": ["data", "meta", "error"],
    "itemCount": 500,
    "itemType": "object",
    "compressionHints": [
      "response_compress(path, 'first_of_array', 'data')",
      "response_compress(path, 'sample_array', 'data', arrayHead=3, arrayTail=2)",
      "response_compress(path, 'keys_only', 'data')",
      "response_compress(path, 'select_keys', 'data', selectKeys=[id, name])"
    ],
    "navigationHints": {
      "paths": ["data", "meta", "error"],
      "arrays": [
        {"path": "data", "length": 500}
      ]
    }
  }
}
```

| 필드 | 타입 | 설명 |
|------|------|------|
| `type` | string | 최상위 타입: "object" 또는 "array" |
| `size` | int | 파일 크기 (바이트) |
| `lineCount` | int | 파일의 줄 수 |
| `depth` | int | 검사된 최대 중첩 깊이 |
| `structure` | object | 키, 타입, 배열 길이가 있는 재귀적 구조 |
| `schemaHint` | string | 최상위 형태의 한 줄 요약 |
| `keys` | array | 최상위 키 (객체의 경우) |
| `itemCount` | int | 배열 길이 (배열의 경우) |
| `compressionHints` | array | 매개변수가 포함된 제안된 `response_compress` 호출 |
| `navigationHints` | object | 길이가 있는 최상위 경로 및 배열 |

### 세부 사항

- 경로가 유효하지 않거나 응답 디렉토리 내부가 아니면 `validation_failed` 반환
- 파일이 존재하지 않으면 `not_found` 반환
- 파일이 유효한 JSON이 아니면 `validation_failed` 반환
- `compressionHints` 필드는 `response_compress` 호출에 대한 즉시 사용 가능한 제안을 제공

---

## response_compress

### 목적

저장된 응답 파일 내의 JSON 값을 줄여 응답 크기 제한 내에 맞추고 LLM에 인라인으로 반환할 수 있도록 합니다. 여러 압축 모드를 통해 크기와 정보 간의 적절한 트레이드오프를 선택할 수 있습니다.

### 사용 시기

- `response_outline` 후 구조를 이해한 다음
- 큰 응답에서 인라인으로 데이터를 가져와야 할 때
- `response_slice`가 너무 좁고 더 넓은 보기가 필요할 때

### 작동 방식

저장된 응답 파일을 읽고, 지정된 JSON 경로로 이동하며, 압축 모드를 적용하고, 압축된 결과를 반환합니다. 결과가 여전히 크기 제한을 초과하면 새 파일에 저장됩니다.

### 매개변수

| 매개변수 | 타입 | 필수 | 설명 |
|---------|------|------|------|
| `path` | string | 예 | `fileRef.path`의 절대 경로 |
| `jsonPath` | string | 아니요 | 압축할 값의 경로 (예: `data` 또는 `data.0`) |
| `mode` | string | 예 | 압축 모드 (아래 표 참조) |
| `arrayHead` | int | 아니요 | `sample_array` 모드에서 유지할 선행 항목 (기본값: 3) |
| `arrayTail` | int | 아니요 | `sample_array` 모드에서 유지할 후행 항목 (기본값: 2) |
| `stringLen` | int | 아니요 | `truncate_strings` 모드의 최대 문자열 길이 (기본값: 80) |
| `selectKeys` | array | 아니요 | `select_keys` 모드에서 유지할 키 |

### 압축 모드

| 모드 | 설명 | 최적 용도 |
|------|------|----------|
| `first_of_array` | 배열의 첫 번째 요소만 유지 | 모든 요소가 동일한 구조일 때 |
| `sample_array` | 배열의 처음과 끝 유지 | 값의 범위를 확인해야 할 때 |
| `truncate_strings` | 모든 문자열을 `stringLen`자로 단축 | 문자열이 매우 길지만 구조가 중요할 때 |
| `keys_only` | 객체 값을 타입 이름으로 대체 | 구조만 필요할 때 |
| `select_keys` | 모든 객체에서 지정된 키만 유지 | 많은 객체에서 특정 필드가 필요할 때 |

### 응답

```json
{
  "body": [
    { "id": 1, "name": "Rex", "status": "available" },
    { "id": 2, "name": "Max", "status": "pending" }
  ],
  "hint": "first_of_array 모드를 사용하여 배열을 500개에서 2개 항목으로 압축했습니다"
}
```

| 필드 | 타입 | 설명 |
|------|------|------|
| `body` | any | 압축된 JSON 값 (크기 제한 이내일 때 표시) |
| `fileRef` | object | 파일 참조 (여전히 너무 클 때 표시) |
| `hint` | string | 압축된 내용에 대한 설명 |

### 세부 사항

- 압축된 결과가 여전히 `max_response_size`를 초과하면 새 파일에 저장되고 `FileReference`가 반환됨
- 기본값: `arrayHead=3`, `arrayTail=2`, `stringLen=80`
- 잘못된 경로, 잘못된 JSONPath 또는 JSON이 아닌 파일에 대해 `validation_failed` 반환
- 파일이 존재하지 않거나 JSONPath가 일치하지 않으면 `not_found` 반환

---

## response_slice

### 목적

논리적 JSON 경로 또는 줄 범위로 저장된 JSON 응답 파일의 특정 조각을 추출합니다. `response_compress`와 달리 원시 수정되지 않은 데이터를 제공합니다.

### 사용 시기

- 큰 응답에서 특정 요소나 값이 필요할 때
- `response_compress`가 충분한 세부 정보를 제공하지 않을 때
- 응답을 단계별로 탐색하려고 할 때

### 작동 방식

저장된 응답 파일을 읽고 JSON 경로(예: `data.3.name`) 또는 줄 범위(예: `120-240`)로 조각을 추출합니다. 배열과 객체를 단계별로 탐색하기 위한 탐색 힌트를 반환합니다.

### 매개변수

| 매개변수 | 타입 | 필수 | 설명 |
|---------|------|------|------|
| `path` | string | 예 | `fileRef.path`의 절대 경로 |
| `jsonPath` | string | 아니요 | 값의 논리적 경로 (예: `data.3.name`) |
| `line` | int | 아니요 | 조각을 중심으로 할 1부터 시작하는 줄 번호 |
| `range` | string | 아니요 | `start-end` 형식의 줄 범위 (예: `120-240`) |
| `around` | int | 아니요 | `line` 주변에 포함할 줄 수 (기본값: 20) |

### 응답

```json
{
  "slice": {
    "lines": [120, 130],
    "fragment": "{\n  \"id\": 1,\n  \"name\": \"Rex\"\n}",
    "value": {
      "id": 1,
      "name": "Rex"
    },
    "jsonPath": "data.0",
    "context": "object",
    "isComplete": true,
    "nextLine": 131,
    "prevLine": 119,
    "nextPath": "data.1",
    "prevPath": null
  }
}
```

| 필드 | 타입 | 설명 |
|------|------|------|
| `lines` | array | 1부터 시작하는 줄 범위 [start, end] |
| `fragment` | string | 원시 JSON 텍스트 (충분히 작을 때) |
| `value` | any | 추출된 JSON 값 |
| `jsonPath` | string | 사용된 JSON 경로 |
| `context` | string | "object", "array" 또는 "value" |
| `isComplete` | bool | 값이 유효한 JSON 조각일 때 true |
| `nextLine` | int | 줄 기반 탐색을 위한 제안된 다음 줄 |
| `prevLine` | int | 제안된 이전 줄 |
| `nextPath` | string | 배열 탐색을 위한 제안된 다음 JSON 경로 |
| `prevPath` | string | 제안된 이전 JSON 경로 |

### 세부 사항

- **줄 번호보다 `jsonPath` 선호** — JSON 경로는 안정적이고 설명적이며, 줄 번호는 파일이 재생성되면 변경됨
- 추출된 조각이 `max_response_size`를 초과하면 새 파일에 저장되고 `FileReference`가 반환됨
- 기본 `around`는 20줄
- 응답에는 배열 탐색을 위한 `nextPath`/`prevPath`와 줄 기반 탐색을 위한 `nextLine`/`prevLine`이 포함됨
- 잘못된 경로, 잘못된 JSONPath, 잘못된 줄/범위 또는 JSON이 아닌 파일에 대해 `validation_failed` 반환
- 파일이 존재하지 않거나 JSONPath가 일치하지 않으면 `not_found` 반환
