# 발견 도구

발견 도구는 LLM이 spec 계층 구조를 탐색할 수 있게 합니다: 모든 spec 찾기, spec으로 드릴다운하여 collection 보기, collection 내의 태그 탐색. `spec_list`로 시작하여 사용 가능한 API를 확인한 후 ID를 사용하여 더 깊이 들어가세요.

---

## spec_list

### 목적

워크스페이스에 등록된 모든 API 명세를 나열합니다. 모든 세션의 시작점입니다 — LLM이 먼저 호출하여 사용 가능한 API를 발견합니다.

### 사용 시기

- 세션 시작 시 어떤 API가 설정되어 있는지 확인
- spec을 추가하거나 제거한 후 목록 새로고침
- 다른 도구에 필요한 spec ID가 필요할 때

### 작동 방식

고유 ID와 도메인 이름이 있는 모든 spec 목록을 반환합니다. 매개변수가 필요하지 않습니다.

### 매개변수

없음.

### 응답

```json
{
  "specs": [
    {
      "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
      "domain": "meteo"
    },
    {
      "id": "b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7",
      "domain": "dadjoke"
    }
  ]
}
```

| 필드 | 타입 | 설명 |
|------|------|------|
| `id` | string | 32자 MD5 해시, spec의 고유 식별자 |
| `domain` | string | spec의 도메인 이름 (예: "meteo", "dadjoke") |

### 세부 사항

- `id`와 `domain`만 반환 — 전체 세부 정보(collection, 태그)는 `spec_by_id` 사용
- 모든 ID는 32자 MD5 16진수 문자열 (`^[0-9a-f]{32}$`)
- 설정된 spec이 없으면 빈 배열 반환

---

## spec_by_id

### 목적

특정 spec의 상세 정보를 가져옵니다: 도메인, 모든 collection, 통계(태그 수, 메서드 수).

### 사용 시기

- `spec_list` 후 spec 내부의 collection을 보려고 할 때
- 추가 탐색을 위해 collection ID가 필요할 때

### 작동 방식

spec ID를 받아 spec 메타데이터와 모든 collection을 개수와 함께 반환합니다.

### 매개변수

| 매개변수 | 타입 | 필수 | 설명 |
|---------|------|------|------|
| `id` | string | 예 | spec의 32자 MD5 해시 |

### 응답

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collections": [
    {
      "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
      "title": "Weather Forecast",
      "llmTitle": "Forecast API",
      "countTags": 3,
      "countMethods": 12
    }
  ]
}
```

| 필드 | 타입 | 설명 |
|------|------|------|
| `spec.id` | string | Spec 식별자 |
| `spec.domain` | string | Spec 도메인 이름 |
| `collections[].id` | string | Collection 식별자 |
| `collections[].title` | string | 사람이 읽을 수 있는 제목 |
| `collections[].llmTitle` | string | LLM 친화적 제목 (선택 사항) |
| `collections[].countTags` | int | collection의 태그 수 |
| `collections[].countMethods` | int | collection의 HTTP 메서드 수 |

### 세부 사항

- spec ID가 존재하지 않으면 `not_found` 오류 반환
- `id`는 유효한 32자 MD5 16진수 문자열이어야 함

---

## collection_by_spec

### 목적

특정 spec 내의 모든 collection을 나열합니다. `spec_by_id`와 유사하지만 추가 spec 메타데이터 없이 collection 목록만 반환합니다.

### 사용 시기

- 이미 spec ID가 있고 collection 목록만 필요할 때
- `spec_by_id`의 더 가벼운 대안으로

### 매개변수

| 매개변수 | 타입 | 필수 | 설명 |
|---------|------|------|------|
| `specId` | string | 예 | spec의 32자 MD5 해시 |

### 응답

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collections": [
    {
      "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
      "title": "Weather Forecast",
      "llmTitle": "Forecast API",
      "countTags": 3,
      "countMethods": 12
    }
  ]
}
```

### 세부 사항

- spec이 존재하지 않으면 `not_found` 반환
- `spec_by_id`와 동일한 데이터지만 추가 spec 래퍼 없음

---

## collection_by_id

### 목적

특정 collection의 상세 정보를 가져옵니다: 메타데이터, 상위 spec, collection 내의 모든 태그.

### 사용 시기

- `collection_by_spec` 후 collection 내부의 태그를 보려고 할 때
- `tag_by_id` 또는 `endpoint_by_tag`에 필요한 태그 ID가 필요할 때

### 매개변수

| 매개변수 | 타입 | 필수 | 설명 |
|---------|------|------|------|
| `id` | string | 예 | collection의 32자 MD5 해시 |

### 응답

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collection": {
    "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
    "title": "Weather Forecast",
    "countMethods": 12
  },
  "tags": [
    {
      "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "title": "forecast",
      "countMethods": 5
    },
    {
      "id": "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7",
      "title": "current",
      "countMethods": 7
    }
  ]
}
```

| 필드 | 타입 | 설명 |
|------|------|------|
| `spec` | object | 상위 spec (id, domain) |
| `collection` | object | Collection 메타데이터 (id, title, countMethods) |
| `tags[]` | array | id, title, countMethods가 있는 태그 목록 |

### 세부 사항

- collection ID가 존재하지 않으면 `not_found` 반환
- 태그는 ID와 함께 반환 — 실제 엔드포인트를 보려면 `endpoint_by_tag(tagId)` 사용

---

## tag_by_spec

### 목적

모든 collection에 걸쳐 전체 spec의 모든 태그를 나열합니다. 사용 가능한 모든 태그의 조감도를 보는 데 유용합니다.

### 사용 시기

- 각 collection으로 드릴다운하지 않고 spec의 모든 태그를 보려고 할 때
- 필요한 태그가 어떤 collection에 있는지 모를 때

### 매개변수

| 매개변수 | 타입 | 필수 | 설명 |
|---------|------|------|------|
| `specId` | string | 예 | spec의 32자 MD5 해시 |

### 응답

```json
{
  "tags": [
    {
      "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "title": "forecast",
      "countMethods": 5
    },
    {
      "id": "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7",
      "title": "current",
      "countMethods": 7
    }
  ]
}
```

### 세부 사항

- spec이 존재하지 않으면 `not_found` 반환
- 태그는 spec의 모든 collection에서 집계됨

---

## tag_by_collection

### 목적

특정 collection 내의 모든 태그를 나열합니다. `tag_by_spec`과 달리 상위 spec 및 collection 메타데이터도 반환합니다.

### 사용 시기

- `collection_by_id` 후 태그 목록을 확인하려고 할 때
- 전체 컨텍스트(spec + collection + 태그)가 필요할 때

### 매개변수

| 매개변수 | 타입 | 필수 | 설명 |
|---------|------|------|------|
| `collectionId` | string | 예 | collection의 32자 MD5 해시 |

### 응답

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collection": {
    "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
    "title": "Weather Forecast",
    "countMethods": 12
  },
  "tags": [
    {
      "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "title": "forecast",
      "countMethods": 5
    }
  ]
}
```

### 세부 사항

- collection이 존재하지 않으면 `not_found` 반환
- `tag_by_spec`과 동일한 태그 데이터지만 하나의 collection으로 범위가 제한됨

---

## tag_by_id

### 목적

단일 태그에 대한 정보를 가져옵니다: ID, 제목, 포함된 메서드 수. 태그 자체에 대한 정보를 알려줍니다 — 실제 엔드포인트를 보려면 `endpoint_by_tag`를 사용하세요.

### 사용 시기

- 태그 ID가 있고 이름과 크기를 확인하려고 할 때
- `endpoint_by_tag`를 호출하기 전에 예상되는 엔드포인트 수를 이해하려고 할 때

### 매개변수

| 매개변수 | 타입 | 필수 | 설명 |
|---------|------|------|------|
| `id` | string | 예 | 태그의 32자 MD5 해시 |

### 응답

```json
{
  "tag": {
    "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
    "title": "forecast",
    "countMethods": 5
  }
}
```

| 필드 | 타입 | 설명 |
|------|------|------|
| `tag.id` | string | 태그 식별자 |
| `tag.title` | string | 사람이 읽을 수 있는 태그 이름 |
| `tag.countMethods` | int | 이 태그의 HTTP 메서드 수 |

### 세부 사항

- 태그가 존재하지 않으면 `not_found` 반환
- 이 도구는 태그 메타데이터만 반환 — 실제 엔드포인트 목록은 `endpoint_by_tag` 사용
