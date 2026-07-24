# 실행 도구

실행 도구는 swag2mcp의 핵심입니다: **search**는 ID가 없을 때 엔드포인트를 찾고, **inspect**는 전체 OpenAPI 계약을 보여주며, **invoke**는 실제 API 호출을 실행합니다. 항상 이 순서로 사용하세요: search → inspect → invoke.

---

## search

### 목적

엔드포인트 ID가 없을 때 엔드포인트를 찾는 유일한 도구입니다. bluge 검색 엔진을 사용하여 모든 spec의 모든 엔드포인트에 대해 전문 검색을 수행합니다.

### 사용 시기

- 엔드포인트 ID를 모를 때
- 키워드, 메서드, 태그 또는 경로로 엔드포인트를 찾으려고 할 때
- 특정 기능에 어떤 엔드포인트가 존재하는지 발견해야 할 때

### 작동 방식

모든 spec의 전문 검색 인덱스를 검색합니다. 필드 필터, 부울 연산자, 퍼지 매칭, 와일드카드 등이 있는 구조화된 쿼리를 지원합니다.

### 매개변수

| 매개변수 | 타입 | 필수 | 설명 |
|---------|------|------|------|
| `query` | string | 예 | 검색 쿼리 (구조화된 구문 지원) |
| `limit` | int | 예 | 반환할 최대 결과 수 (1-50) |

### 쿼리 구문

| 예시 | 설명 |
|------|------|
| `pet` | 모든 필드에 대한 단순 텍스트 검색 |
| `method:GET` | HTTP 메서드로 필터링 |
| `tag:pet` | 태그 이름으로 필터링 |
| `path:"/api/v1/users"` | 정확한 경로 검색 |
| `+method:POST +tag:pet` | 두 조건 모두 일치해야 함 |
| `-method:DELETE` | DELETE 메서드 제외 |
| `create~` | 퍼지 검색 (오타 허용) |
| `path:/api/v1/*` | 와일드카드 경로 검색 |
| `/pattern/` | 정규식 검색 |
| `term^3` | 용어 관련성 부스트 |

**검색 가능한 필드:** `method` (키워드), `tag` (키워드), `path` (텍스트), `summary` (텍스트), `_all` (기본 텍스트 필드).

**지원되지 않음:** 그룹화를 위한 괄호, 명시적 `AND`/`OR` 연산자, 필드 그룹화.

### 응답

```json
{
  "endpoints": [
    {
      "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
      "tagId": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "tagName": "forecast",
      "collectionId": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
      "collectionTitle": "Weather Forecast",
      "specId": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
      "specDomain": "meteo",
      "method": "GET",
      "path": "/v1/forecast",
      "summary": "위치에 대한 날씨 예보 가져오기"
    }
  ]
}
```

각 결과는 전체 계보(spec → collection → tag)를 포함하여 LLM이 관련 엔드포인트로 이동할 수 있도록 합니다.

### 세부 사항

- `limit`은 1에서 50 사이여야 함 (그렇지 않으면 `validation_failed` 반환)
- `query`는 필수 (비어 있으면 `validation_failed` 반환)
- 결과는 관련성 순서로 반환됨 (최적 일치 우선)
- 필드 필터(`method:GET`, `tag:pet`)를 사용하여 결과 좁히기
- 정확한 경로 매칭은 따옴표 사용: `path:"/v1/forecast"`

---

## inspect

### 목적

엔드포인트의 전체 OpenAPI 작업 객체를 검색합니다: 모든 매개변수, 요청 본문 스키마, 응답 스키마, base URL, 전체 URL. `invoke` **전에** 호출하여 엔드포인트의 계약을 이해하는 도구입니다.

### 사용 시기

- 항상 `invoke` 전에 — 올바른 호출을 위해 전체 계약이 필요함
- API의 기술적 세부 정보를 사용자에게 설명해야 할 때
- 필수 매개변수, 요청 본문 구조 또는 응답 형식을 알아야 할 때

### 작동 방식

인덱스에서 엔드포인트를 조회하고 모든 스키마가 해결된 완전한 OpenAPI 작업 객체를 반환합니다.

### 매개변수

| 매개변수 | 타입 | 필수 | 설명 |
|---------|------|------|------|
| `endpointId` | string | 예 | 엔드포인트의 32자 MD5 해시 |

### 응답

```json
{
  "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
  "tagId": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
  "collectionId": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
  "specId": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
  "specDomain": "meteo",
  "method": "POST",
  "path": "/pet",
  "baseUrl": "https://meteo.swagger.io/v2",
  "fullUrl": "https://meteo.swagger.io/v2/pet",
  "operation": {
    "id": "addPet",
    "tags": ["pet"],
    "summary": "새 애완동물 추가",
    "description": "스토어에 새 애완동물 추가",
    "deprecated": false,
    "parameters": [
      {
        "name": "petId",
        "in": "path",
        "description": "애완동물의 ID",
        "required": true,
        "schema": {
          "type": "integer",
          "format": "int64"
        }
      }
    ],
    "requestBody": {
      "description": "추가할 애완동물 객체",
      "required": true,
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "name": { "type": "string" },
              "status": { "type": "string", "enum": ["available", "pending", "sold"] }
            },
            "required": ["name"]
          }
        }
      }
    },
    "responses": {
      "200": {
        "description": "성공적인 작업",
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/Pet"
            }
          }
        }
      },
      "405": {
        "description": "잘못된 입력"
      }
    }
  }
}
```

| 필드 | 타입 | 설명 |
|------|------|------|
| `baseUrl` | string | API의 Base URL (설정에서) |
| `fullUrl` | string | 엔드포인트의 전체 URL (base + path) |
| `operation.parameters[]` | array | 이름, 위치(path/query/header/cookie), 설명, 필수 플래그, 스키마가 있는 매개변수 |
| `operation.requestBody` | object | 콘텐츠 유형 및 스키마가 있는 요청 본문 |
| `operation.responses` | map | 설명 및 스키마가 있는 응답 코드 |
| `operation.deprecated` | bool | 엔드포인트가 폐기되었는지 여부 |

### 세부 사항

- 엔드포인트가 존재하지 않으면 `not_found` 반환
- 전체 OpenAPI 작업을 반환하는 **유일한** 도구 — `endpoint_by_id`는 요약만 반환
- 필수 매개변수와 본문 구조를 이해하려면 `invoke` 전에 항상 `inspect` 호출
- `operation` 객체에는 전체 스키마 정의로 해결된 `$ref` 참조가 포함됨

---

## invoke

### 목적

엔드포인트에 대한 실제 API 호출을 실행합니다. 실제 HTTP 요청을 만드는 유일한 도구입니다. 인증이 자동으로 적용됩니다 — 먼저 `auth`를 호출할 필요가 없습니다.

### 사용 시기

- 엔드포인트의 계약을 이해하기 위해 `inspect`를 호출한 후에만
- 파괴적 작업(POST, PUT, PATCH, DELETE)의 경우 명시적 사용자 확인이 있는 경우에만
- 사용자가 API 호출을 요청하고 모든 필수 매개변수가 있을 때

### 작동 방식

1. 인덱스에서 엔드포인트 조회
2. URL에 경로 매개변수 대입
3. 쿼리 매개변수 추가
4. 헤더 및 쿠키 추가
5. 요청 본문을 JSON으로 직렬화
6. 자동으로 인증 획득 및 적용 (토큰, 헤더, 쿼리 매개변수)
7. HTTP 요청 실행
8. 응답 반환 또는 너무 크면 파일에 저장

### 매개변수

| 매개변수 | 타입 | 필수 | 설명 |
|---------|------|------|------|
| `endpointId` | string | 예 | 엔드포인트의 32자 MD5 해시 |
| `parameters` | object | 아니요 | 키-값 쌍의 경로, 쿼리, 헤더 매개변수 |
| `requestBody` | object | 아니요 | POST/PUT/PATCH 요청의 요청 본문 |
| `headers` | object | 아니요 | 전송할 추가 HTTP 헤더 |
| `cookies` | object | 아니요 | 전송할 추가 HTTP 쿠키 |

### 응답 (인라인)

```json
{
  "statusCode": 200,
  "headers": {
    "content-type": "application/json"
  },
  "body": {
    "id": 1,
    "name": "Rex",
    "status": "available"
  }
}
```

### 응답 (파일 참조 — 본문이 크기 제한 초과 시)

```json
{
  "statusCode": 200,
  "headers": {
    "content-type": "application/json"
  },
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

| 필드 | 타입 | 설명 |
|------|------|------|
| `statusCode` | int | HTTP 응답 상태 코드 |
| `headers` | object | HTTP 응답 헤더 |
| `body` | any | 응답 본문 (크기 제한 이내일 때 표시) |
| `fileRef` | object | 파일 참조 (본문이 크기 제한 초과 시 표시) |

### 큰 응답 작업

`invoke`가 `fileRef`를 반환하면 응답 도구를 사용하여 데이터를 탐색하세요:

1. **`response_outline(path)`** — 구조적 요약 얻기 (키, 타입, 배열 길이)
2. **`response_compress(path, mode)`** — 데이터를 압축하여 인라인에 맞추기
3. **`response_slice(path, jsonPath)`** — 특정 조각 추출

### 세부 사항

- **인증은 자동:** `invoke` 도구는 spec의 인증 설정에서 자동으로 인증을 획득하고 적용합니다. 먼저 `auth`를 호출할 **필요가 없습니다**.
- **속도 제한:** 각 엔드포인트에는 10초의 쿨다운이 있습니다. 동일한 엔드포인트에 대한 10초 내의 두 번째 호출은 자동으로 차단됩니다(`rate_limit` 오류 반환).
- **응답 크기 제한:** 기본값은 2 KB입니다(`max_response_size`로 설정 가능). 응답이 이 제한을 초과하면 `{workspace}/responses/`에 저장되고 인라인 `body` 대신 `FileReference`가 반환됩니다.
- **매개변수 처리:** 경로 매개변수는 URL에 대입됩니다. 쿼리 매개변수가 추가됩니다. 요청의 매개변수가 작업 명세 기본값을 재정의합니다.
- **요청 본문:** POST/PUT/PATCH의 경우 본문이 JSON으로 직렬화됩니다. `Content-Type`이 자동으로 `application/json`으로 설정됩니다.
- **오류 처리:** HTTP 오류(2xx 아님)는 힌트에 상태 코드와 응답 본문과 함께 `invoke_error`로 반환됩니다.
- **파괴적 작업:** 명시적 사용자 확인 없이 POST/PUT/PATCH/DELETE를 호출하지 마세요.
