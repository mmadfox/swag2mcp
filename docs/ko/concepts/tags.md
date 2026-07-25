# Tags

태그는 collection 내에서 관련 엔드포인트를 그룹화하는 카테고리입니다. 태그가 있을 수도 있고 없을 수도 있습니다 — 모든 collection에 태그가 있는 것은 아니며, collection은 여러 개의 태그를 가질 수 있습니다.

태그는 OpenAPI/Swagger/Postman 파일 자체에서 가져옵니다. 태그에 대한 **YAML 설정이 없습니다** — `swag2mcp.yaml`에서 태그를 생성, 이름 변경 또는 삭제할 수 없습니다. 태그를 변경하는 유일한 방법은 원본 명세 파일을 편집하는 것입니다.

## 계층 구조

```
Spec (domain, e.g. "meteo")
  └── Collection (spec file, e.g. forecast.yml)
        └── Tag "weather"
              └── GET /forecast
              └── GET /forecast/hourly
        └── Tag "alerts"
              └── GET /alerts
```

## 태그 생성 방식

태그는 파싱 중에 명세 문서에서 추출됩니다:

**OpenAPI 3.x / Swagger 2.0** — 각 작업의 `tags` 목록이 태그가 됩니다:

```yaml
paths:
  /pet:
    get:
      tags: ["pets"]
      summary: "ID로 애완동물 찾기"
    post:
      tags: ["pets"]
      summary: "새 애완동물 추가"
  /pet/{petId}/uploadImage:
    post:
      tags: ["pet_images"]
      summary: "이미지 업로드"
```

**Postman** — 각 최상위 폴더가 태그가 됩니다. 중첩된 폴더는 마지막 폴더 이름을 사용합니다.

엔드포인트에 태그가 없으면 `"default"` 태그 아래에 배치됩니다.

## 목적

태그는 LLM이 관련 엔드포인트 그룹을 찾는 데 도움을 줍니다. LLM은 collection의 모든 엔드포인트를 검색하는 대신 먼저 올바른 태그를 찾은 다음 해당 태그 내의 엔드포인트만 나열할 수 있습니다.

## 태그용 MCP 도구

| 도구 | 설명 |
|------|------|
| `tag_by_spec` | 전체 spec의 모든 태그 |
| `tag_by_collection` | 특정 collection 내의 태그 |
| `tag_by_id` | 태그 세부 정보 (제목, 메서드 수) |
| `endpoint_by_tag` | 태그 아래 그룹화된 엔드포인트 |

## 예시

```
쿼리: "애완동물 collection의 모든 태그 표시"
→ tag_by_collection(collectionId: "...")
→ 결과: pets (5 methods), pet_images (1 method)
```

## 제한 사항

- 태그는 설정 관점에서 읽기 전용입니다. 태그를 추가, 이름 변경 또는 제거하려면 원본 OpenAPI/Swagger/Postman 파일을 편집하고 `swag2mcp update`를 실행하세요.
- YAML 설정에서 collection별로 태그를 필터링하거나 비활성화할 수 없습니다.
