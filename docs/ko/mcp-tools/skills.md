# 스킬

## 출력 형식 사용자 정의

모든 swag2mcp MCP 도구는 구조화된 JSON 데이터를 반환합니다. 이 데이터가 사용자에게 **표시되는 방식**은 LLM의 포맷팅 스킬에 따라 달라지며 — 완전히 제어할 수 있습니다.

### 기본 포맷 스킬

swag2mcp는 모든 도구 응답에 대해 간결하고 사람이 읽을 수 있는 마크다운을 정의하는 내장 포맷팅 스킬과 함께 제공됩니다:

[swag2mcp-format SKILL.md](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md)

이 스킬은 모든 19개 MCP 도구를 다음으로 다룹니다:
- 목록을 위한 간결한 테이블 (spec, collection, 태그, 엔드포인트)
- 상세 보기를 위한 인라인 헤더
- `inspect`를 위한 간결한 스키마 표현
- 모든 응답에서 일관된 스타일

### 스킬이 중요한 이유

동일한 데이터가 스킬에 따라 완전히 다른 방식으로 표시될 수 있습니다:

| 스타일 | 출력 예시 |
|-------|----------|
| **간결한 테이블** (기본값) | `GET /pet/{petId}` — ID로 애완동물 찾기 |
| **상세** | `Method: GET, Path: /pet/{petId}, Summary: ID로 애완동물 찾기, Deprecated: false` |
| **최소** | `GET /pet/{petId}` |
| **기술적** | `GET /pet/{petId} → 200: Pet object, 404: Not found` |
| **커스텀** | 설명할 수 있는 모든 형식 |

### 자신만의 스킬 만들기

원하는 정확한 출력 형식을 설명하여 자신만의 포맷팅 스킬을 작성할 수 있습니다. 스킬은 각 도구에 대한 포맷팅 규칙이 있는 마크다운 파일입니다. 몇 가지 아이디어:

- **JSON 출력** — 머신 소비를 위한 원시 JSON 반환
- **CSV 스타일** — 스프레드시트 가져오기용 테이블 형식 데이터
- **다이어그램 친화적** — API 구조의 Mermaid 또는 ASCII 다이어그램
- **최소** — 메서드와 경로만, 그 외에는 없음
- **문서 스타일** — 전체 설명, 예시, 참고 사항

### 유일한 한계는 모델입니다

포맷된 출력의 품질은 전적으로 LLM이 포맷팅 규칙을 따르는 능력에 달려 있습니다. 명확한 예시가 있는 잘 작성된 스킬은 일관되고 신뢰할 수 있는 출력을 생성합니다. 모호한 스킬은 일관성 없는 결과를 만듭니다.

다음을 할 수 있습니다:
- 기본 스킬을 있는 그대로 사용
- 포크하여 취향에 맞게 포맷팅 조정
- 처음부터 직접 작성
- 작업에 따라 스킬 간 전환

### 스킬 사용 방법

스킬은 시스템 프롬프트 또는 에이전트 설정의 일부로 LLM 클라이언트(OpenCode, Cursor, Claude Desktop 등)에 의해 로드됩니다. 스킬 파일을 연결하는 방법은 클라이언트의 문서를 참조하세요.

OpenCode의 경우 스킬이 `opencode.json`에서 설정됩니다:

```json
{
  "skills": [
    {
      "name": "swag2mcp-format",
      "sourceURL": "https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md"
    }
  ]
}
```
