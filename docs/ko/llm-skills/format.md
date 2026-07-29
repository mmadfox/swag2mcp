# swag2mcp-format

**swag2mcp-format** 스킬은 LLM이 swag2mcp MCP 도구 응답을 깔끔하고 간결하며 읽기 쉬운 Markdown 형식으로 표시하는 방법을 가르칩니다. 이 스킬이 없으면 LLM이 응답 형식을 자체적으로 결정합니다 — 종종 장황하고 일관성이 없습니다.

## 적용 범위

모든 swag2mcp MCP 도구:

- `spec_list`, `spec_by_id` — 사양 개요 및 세부 정보
- `collection_by_spec`, `collection_by_id` — 태그가 있는 컬렉션
- `tag_by_spec`, `tag_by_collection`, `tag_by_id` — 태그 목록
- `endpoint_by_spec`, `endpoint_by_collection`, `endpoint_by_tag`, `endpoint_by_id` — 엔드포인트 목록
- `search` — 검색 결과
- `inspect` — 간결한 스키마가 포함된 전체 작업 세부 정보
- `invoke` — API 호출 결과
- `auth` — 인증 정보
- `info` — 런타임 정보

## LLM 에이전트를 통한 설치

다음 요청을 AI 기반 IDE에 복사하세요:

```
.agents/skills/swag2mcp-format/ 디렉토리를 만들고
https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md 에서 스킬을 추가해
```

## 수동 설치

```bash
mkdir -p .agents/skills/swag2mcp-format
curl -o .agents/skills/swag2mcp-format/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md
```

## 재시작 필요

스킬을 추가한 후 LLM 클라이언트 또는 IDE를 재시작하세요（[개요](overview.md#재시작-필요) 참조）。
