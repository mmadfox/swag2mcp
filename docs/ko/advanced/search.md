# 전문 검색

## 개요

swag2mcp에는 모든 spec의 모든 엔드포인트를 인덱싱하는 내장 전문 검색 엔진(bluge)이 포함되어 있습니다. LLM은 엔드포인트 ID를 몰라도 메서드, 경로, 요약 또는 태그로 엔드포인트를 검색할 수 있습니다.

## 인덱싱 작동 방식

spec이 추가되거나 업데이트되면 모든 엔드포인트가 인덱싱됩니다. 다음 필드를 검색할 수 있습니다:

| 필드 | 설명 | 예시 |
|------|------|------|
| `method` | HTTP 메서드 | `GET`, `POST`, `PUT` |
| `path` | API 엔드포인트 경로 | `/api/v1/users/{id}` |
| `summary` | OpenAPI 요약 | "Find pet by ID" |
| `tag` | 엔드포인트 카테고리 | "pets", "users" |
| `_all` | 모든 필드 결합 | method + path + tag + summary |

인덱스는 MCP 서버가 시작될 때마다 재구축됩니다. 빠른 검색을 위해 메모리에 저장됩니다.

## 쿼리 구문

검색은 정밀한 필터링을 위한 풍부한 쿼리 구문을 지원합니다:

| 예시 | 설명 |
|------|------|
| `pet` | 모든 필드에 대한 단순 텍스트 검색 |
| `method:GET` | 모든 GET 엔드포인트 찾기 |
| `tag:pets` | "pets" 태그의 엔드포인트 찾기 |
| `path:"/api/v1/users"` | 정확한 경로 일치 |
| `+method:POST +tag:pet` | 두 조건 모두 일치해야 함 |
| `-method:DELETE` | DELETE 메서드 제외 |
| `create~` | 퍼지 검색 (오타 허용) |
| `cr*` | 와일드카드 검색 |
| `"find pet"` | 구문 검색 |
| `+summary:pet -method:DELETE` | 요약에 "pet" 포함, DELETE 제외 |

### 필드별 검색

`field:value` 구문을 사용하여 특정 필드 내에서 검색할 수 있습니다:

```
method:GET
tag:pets
path:"/pet/findByStatus"
summary:"find pet by status"
```

### 부울 연산자

- `+` — 용어가 반드시 일치해야 함 (AND)
- `-` — 용어가 일치하지 않아야 함 (NOT)
- 용어 사이의 공백 — OR (모든 용어가 일치 가능)

### 퍼지 및 와일드카드

- `term~` — 퍼지 검색 (유사한 단어 일치, 오타 처리)
- `te*` — 와일드카드 (모든 문자 일치)
- `te?t` — 단일 문자 와일드카드

## 예시

```
# 모든 GET 요청 찾기
method:GET

# pet 태그의 POST 요청 찾기
+method:POST +tag:pet

# 정확한 경로로 엔드포인트 찾기
path:"/pet/findByStatus"

# 설명으로 찾기
"find pet by status"

# DELETE를 제외한 모든 항목 찾기
+summary:pet -method:DELETE

# "create" 퍼지 검색 (오타 처리)
create~
```

## MCP 도구

`search` MCP 도구는 LLM에 검색 엔진을 노출합니다:

```
→ search(query: "find pet by status", limit: 5)
← GET /pet/findByStatus — 상태별 애완동물 찾기
   GET /pet/{petId} — ID로 애완동물 찾기
```

### 매개변수

| 매개변수 | 필수 | 설명 |
|---------|------|------|
| `query` | 예 | 검색 쿼리 (구조화된 구문 지원) |
| `limit` | 예 | 최대 결과 수 (1-50) |

## 중요 참고 사항

- **인덱스는 메모리 내** — MCP 서버가 시작될 때마다 재구축됩니다. 영구 인덱스 파일이 없습니다.
- **모든 필드는 소문자로 변환됨** — 검색은 대소문자를 구분하지 않습니다
- **제한은 최대 50개** — 50개 이상의 결과를 요청할 수 없습니다
- **잘못된 쿼리 구문**은 예제와 함께 도움이 되는 오류 메시지를 반환합니다
- **`_all` 필드**는 method, path, tag, summary를 결합하여 단순 텍스트 검색에 사용됩니다
