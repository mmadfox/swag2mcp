# 엔드포인트

엔드포인트는 호출할 수 있는 특정 HTTP 메서드 + 경로입니다(예: `GET /api/users/{id}`). 엔드포인트는 LLM이 발견, 검사, 호출하는 실제 API 작업입니다.

## 구조

각 엔드포인트에는 다음이 포함됩니다:

- **HTTP 메서드**: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS
- **경로**: `/api/v1/users/{id}`
- **요약**: 엔드포인트가 무엇을 하는지에 대한 짧은 설명 — LLM이 한눈에 목적을 이해하는 데 매우 유용
- **설명**: 엔드포인트의 동작, 매개변수, 사용 사례에 대한 상세 설명
- **매개변수**: path, query, header, cookie
- **요청 본문**: POST/PUT/PATCH용
- **응답**: 상태 코드 및 응답 스키마

`summary`와 `description` 필드는 OpenAPI/Swagger/Postman 파일에서 가져옵니다. 이는 LLM이 엔드포인트의 기능을 이해하는 기본 방법입니다. 잘 작성된 요약은 엔드포인트 발견을 훨씬 더 효과적으로 만듭니다.

## 엔드포인트용 MCP 도구

| 도구 | 설명 |
|------|------|
| `endpoint_by_spec` | spec의 모든 엔드포인트 |
| `endpoint_by_collection` | collection의 엔드포인트 |
| `endpoint_by_tag` | 태그의 엔드포인트 |
| `endpoint_by_id` | 빠른 엔드포인트 요약 |
| `inspect` | 전체 엔드포인트 상세 (스키마, 매개변수) |
| `invoke` | 엔드포인트 호출 |
| `search` | 텍스트로 엔드포인트 검색 |

## 폐기된 엔드포인트

명세에서 `deprecated`로 표시된 엔드포인트는 검사 시 알림과 함께 표시됩니다.

## 설정

엔드포인트는 swag2mcp 관점에서 **읽기 전용**입니다. 엔드포인트에 대한 YAML 설정이 없습니다 — `swag2mcp.yaml`에서 엔드포인트를 추가, 제거, 이름 변경 또는 수정할 수 없습니다.

엔드포인트를 변경하려면(새로 추가, 요약 업데이트, 매개변수 수정, 폐기 표시) 원본 OpenAPI/Swagger/Postman 파일을 편집하고 `swag2mcp update`를 실행하여 다시 파싱하고 재인덱싱하세요.

## 예시

```
쿼리: "GET /pet/{petId}의 세부 정보 표시"
→ inspect(endpointId: "abc123...")
→ 결과:
  GET /pet/{petId}
  요약: ID로 애완동물 찾기
  설명: ID로 단일 애완동물 반환
  매개변수:
    - petId (path, integer, 필수)
  응답:
    - 200: Pet 객체
    - 400: 오류
    - 404: 찾을 수 없음
```
