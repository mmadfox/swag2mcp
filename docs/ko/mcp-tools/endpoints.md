# 엔드포인트 도구

엔드포인트 도구는 LLM이 계층 구조의 다양한 수준에서 API 엔드포인트를 볼 수 있게 합니다: spec의 모든 엔드포인트, collection의 엔드포인트, 태그의 엔드포인트, 또는 단일 엔드포인트 요약. 검사하거나 호출하기 전에 사용 가능한 작업을 발견하는 데 사용하세요.

---

## endpoint_by_spec

### 목적

모든 collection과 태그에 걸쳐 전체 spec의 모든 엔드포인트를 나열합니다. 가장 포괄적인 보기를 반환합니다 — 전체 컨텍스트(태그, collection, spec)가 있는 spec의 모든 엔드포인트.

### 사용 시기

- spec에서 사용 가능한 모든 엔드포인트를 보려고 할 때
- 필요한 엔드포인트가 어떤 collection이나 태그에 있는지 모를 때
- `spec_by_id` 후 전체 엔드포인트 목록을 얻으려고 할 때

### 매개변수

| 매개변수 | 타입 | 필수 | 설명 |
|---------|------|------|------|
| `specId` | string | 예 | spec의 32자 MD5 해시 |

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

| 필드 | 타입 | 설명 |
|------|------|------|
| `id` | string | 엔드포인트 식별자 |
| `tagId` | string | 상위 태그 식별자 |
| `tagName` | string | 사람이 읽을 수 있는 태그 이름 |
| `collectionId` | string | 상위 collection 식별자 |
| `collectionTitle` | string | 사람이 읽을 수 있는 collection 제목 |
| `specId` | string | 상위 spec 식별자 |
| `specDomain` | string | Spec 도메인 이름 |
| `method` | string | HTTP 메서드 (GET, POST, PUT, DELETE 등) |
| `path` | string | API 경로 (예: /v1/forecast) |
| `summary` | string | 엔드포인트 기능에 대한 사람이 읽을 수 있는 요약 |

### 세부 사항

- spec이 존재하지 않으면 `not_found` 반환
- 각 엔드포인트는 컨텍스트를 위해 전체 계보(spec → collection → tag)를 포함
- 단일 엔드포인트의 빠른 요약은 `endpoint_by_id` 사용

---

## endpoint_by_collection

### 목적

태그에 관계없이 특정 collection 내의 모든 엔드포인트를 나열합니다. spec 및 collection 메타데이터와 함께 collection별로 그룹화된 엔드포인트를 반환합니다.

### 사용 시기

- `collection_by_id` 후 collection의 모든 엔드포인트를 보려고 할 때
- collection의 전체 API 표면을 탐색하려고 할 때

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
  "endpoints": [
    {
      "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
      "tagId": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "tagName": "forecast",
      "method": "GET",
      "path": "/v1/forecast",
      "summary": "위치에 대한 날씨 예보 가져오기"
    }
  ]
}
```

### 세부 사항

- collection이 존재하지 않으면 `not_found` 반환
- 컨텍스트를 위해 spec 및 collection 메타데이터 포함
- collection 내의 모든 태그에서 엔드포인트가 함께 반환됨

---

## endpoint_by_tag

### 목적

특정 태그 아래 그룹화된 모든 엔드포인트를 나열합니다. 가장 집중된 보기입니다 — 하나의 collection 내의 하나의 태그에 있는 엔드포인트.

### 사용 시기

- `tag_by_id` 후 태그의 실제 엔드포인트를 보려고 할 때
- 태그를 알고 있고 해당 작업을 보려고 할 때

### 매개변수

| 매개변수 | 타입 | 필수 | 설명 |
|---------|------|------|------|
| `tagId` | string | 예 | 태그의 32자 MD5 해시 |

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
  "tag": {
    "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
    "title": "forecast",
    "countMethods": 5
  },
  "endpoints": [
    {
      "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
      "method": "GET",
      "path": "/v1/forecast",
      "summary": "위치에 대한 날씨 예보 가져오기"
    }
  ]
}
```

### 세부 사항

- 태그가 존재하지 않으면 `not_found` 반환
- 전체 컨텍스트 포함: spec, collection, 태그 메타데이터
- 엔드포인트는 단일 collection 내의 단일 태그로 범위가 제한됨

---

## endpoint_by_id

### 목적

단일 엔드포인트의 빠른 요약을 가져옵니다: 메서드, 경로, 요약, 폐기 상태. 가벼운 도구입니다 — 전체 OpenAPI 작업 객체(매개변수, 요청 본문, 응답 스키마)는 `inspect`를 사용하세요.

### 사용 시기

- 엔드포인트 ID가 있고 무엇을 하는지 빠르게 확인하려고 할 때
- 전체 세부 정보를 위해 `inspect`를 호출할지 결정하기 전에
- 호출 전에 메서드와 경로를 확인해야 할 때

### 매개변수

| 매개변수 | 타입 | 필수 | 설명 |
|---------|------|------|------|
| `id` | string | 예 | 엔드포인트의 32자 MD5 해시 |

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
  "tag": {
    "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
    "title": "forecast",
    "countMethods": 5
  },
  "endpoint": {
    "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
    "method": "GET",
    "path": "/v1/forecast",
    "summary": "위치에 대한 날씨 예보 가져오기"
  }
}
```

| 필드 | 타입 | 설명 |
|------|------|------|
| `endpoint.id` | string | 엔드포인트 식별자 |
| `endpoint.method` | string | HTTP 메서드 |
| `endpoint.path` | string | API 경로 |
| `endpoint.summary` | string | 사람이 읽을 수 있는 요약 |

### 세부 사항

- 엔드포인트가 존재하지 않으면 `not_found` 반환
- **빠른 요약**입니다 — 매개변수, 요청 본문 또는 응답 스키마를 반환하지 않음
- 전체 기술 세부 정보(`invoke` 전에 필요)는 `inspect` 사용
