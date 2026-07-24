# MCP 도구

## 개요

swag2mcp는 LLM 에이전트가 Model Context Protocol을 통해 API에 완전히 접근할 수 있도록 **19개의 MCP 도구**를 제공합니다. 이 도구들은 전체 워크플로우를 다룹니다: 사용 가능한 API 발견, spec 계층 구조 탐색, 엔드포인트 검색 및 검사, API 호출 실행, 큰 응답 작업.

### 도구가 해결하는 문제

- **발견** — LLM이 사전에 ID를 알지 못해도 spec, collection, 태그를 찾을 수 있음
- **탐색** — spec → collection → tag → endpoint의 구조화된 계층 구조로 드릴다운
- **검색** — ID가 없을 때 모든 엔드포인트에 대한 전문 검색
- **검사** — 호출 전에 전체 OpenAPI 작업 객체 얻기
- **실행** — 자동 인증으로 실제 API 호출 실행
- **큰 응답 처리** — 인라인에 맞지 않는 초과된 응답의 개요, 압축, 조각 추출

### 읽기 전용 vs 변경 가능

| 타입 | 개수 | 도구 |
|------|------|------|
| **읽기 전용** | 17 | 모든 발견, 엔드포인트, 검색, 검사, 정보, 응답 도구 |
| **변경 가능** | 2 | `invoke` (실제 HTTP 호출), `auth` (토큰 검색) |

읽기 전용 도구는 MCP 프로토콜에서 `ReadOnlyHint=true`와 `IdempotentHint=true`로 표시되어 LLM에 부작용 없이 호출해도 안전함을 알립니다.

### 오류 처리

모든 도구는 머신 판독 가능 코드와 사람이 읽을 수 있는 메시지가 있는 구조화된 `LLMError` 객체로 오류를 반환합니다:

| 오류 코드 | 의미 |
|-----------|------|
| `validation_failed` | 잘못된 입력 (잘못된 ID 형식, 필수 필드 누락) |
| `not_found` | 인덱스 또는 워크스페이스에서 엔티티를 찾을 수 없음 |
| `rate_limit` | 동일한 엔드포인트에 10초 내 두 번째 `invoke` 호출 |
| `invoke_error` | HTTP 호출 실패, 다운로드 실패 |
| `auth_error` | 인증 토큰 검색 실패 |
| `config_error` | 설정 파일 로드 또는 저장 실패 |
| `parse_error` | 명세 파일 파싱 실패 |

## 카테고리

| 카테고리 | 도구 | 설명 |
|----------|------|------|
| **발견** | `spec_list`, `spec_by_id`, `collection_by_spec`, `collection_by_id`, `tag_by_spec`, `tag_by_collection`, `tag_by_id` | Spec 계층 구조 탐색: spec, collection, 태그 찾기 |
| **엔드포인트** | `endpoint_by_spec`, `endpoint_by_collection`, `endpoint_by_tag`, `endpoint_by_id` | 계층 구조의 다양한 수준에서 엔드포인트 보기 |
| **실행** | `search`, `inspect`, `invoke` | 검색, 전체 계약 검사, API 호출 |
| **유틸리티** | `auth`, `info`, `response_outline`, `response_compress`, `response_slice` | 인증 토큰, 런타임 정보, 큰 응답 처리 |
| **스킬** | [포맷팅 가이드](/mcp-tools/skills) | 도구 응답 표시 방식 사용자 정의 |

## 전체 목록

| 도구 | 설명 |
|------|------|
| `spec_list` | 워크스페이스의 모든 API 명세 나열 |
| `spec_by_id` | collection이 포함된 상세 spec 정보 가져오기 |
| `collection_by_spec` | spec 내의 collection 나열 |
| `collection_by_id` | 태그가 포함된 collection 세부 정보 가져오기 |
| `tag_by_spec` | spec 전체의 모든 태그 나열 |
| `tag_by_collection` | collection 내의 태그 나열 |
| `tag_by_id` | 태그 세부 정보 가져오기 (ID, 제목, 메서드 수) |
| `endpoint_by_spec` | spec의 모든 엔드포인트 나열 |
| `endpoint_by_collection` | collection의 엔드포인트 나열 |
| `endpoint_by_tag` | 태그의 엔드포인트 나열 |
| `endpoint_by_id` | 빠른 엔드포인트 요약 (메서드, 경로, 요약) |
| `search` | 모든 엔드포인트에 대한 전문 검색 |
| `inspect` | 전체 OpenAPI 작업 세부 정보 (매개변수, 스키마) |
| `invoke` | 실제 API 호출 실행 |
| `auth` | spec의 인증 토큰 또는 헤더 가져오기 |
| `info` | 런타임 정보 (버전, spec, 설정) |
| `response_outline` | 큰 응답 파일의 구조적 요약 |
| `response_compress` | 큰 응답을 압축하여 인라인에 맞추기 |
| `response_slice` | 큰 응답의 조각 추출 |

## 탐색 계층 구조

```
spec_list
  └── spec_by_id(id)
        └── collection_by_spec(specId)
              └── collection_by_id(id)
                    └── tag_by_collection(collectionId)
                          └── tag_by_id(id)
                                └── endpoint_by_tag(tagId)
                                      └── endpoint_by_id(id)
                                            └── inspect(endpointId)
                                                  └── invoke(endpointId)
```

ID가 없으면 `search`를 사용하여 쿼리로 엔드포인트를 찾으세요. `invoke`가 `fileRef`를 반환하면(응답이 너무 큼) `response_outline` → `response_compress` 또는 `response_slice`를 사용하여 데이터를 탐색하세요.
