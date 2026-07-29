# Skills for LLM

스킬은 LLM 에이전트가 swag2mcp를 더 효과적으로 사용하는 방법을 가르치는 Markdown 파일입니다. 에이전트의 시스템 프롬프트의 일부로 로드되어 LLM에게 응답 형식 지정 및 CLI 명령 이해에 대한 정확한 지침을 제공합니다.

## 사용 가능한 스킬

| 스킬 | 설명 | 소스 |
|------|------|------|
| **swag2mcp-format** | 모든 MCP 도구 응답을 간결하고 읽기 쉬운 Markdown 테이블로 포맷 | `swag2mcp-format/SKILL.md` |
| **swag2mcp-cli** | 전체 CLI 참조 — LLM이 모든 명령, 플래그 및 설정 옵션을 정확히 파악 | `swag2mcp-cli/SKILL.md` |

## 스킬이 중요한 이유

포맷 스킬이 없으면 LLM이 도구 결과를 자체적으로 표시하는 방식을 결정합니다 — 종종 장황하고 일관성이 없습니다. 포맷 스킬은 모든 응답이 동일한 깔끔한 패턴을 따르도록 보장합니다: 목록에는 간결한 테이블, 세부 정보에는 인라인 헤더, 간결한 스키마.

CLI 스킬을 사용하면 LLM이 swag2mcp 명령에 대한 "어떻게..." 질문에 추측 없이 정확하게 답변할 수 있습니다.

## LLM 에이전트를 통한 설치

다음 요청을 AI 기반 IDE(OpenCode, Cursor, Claude Desktop, VS Code 등)에 복사하세요:

```
내 프로젝트에 swag2mcp 스킬을 추가해줘:

1. .agents/skills/swag2mcp-format/ 디렉토리를 만들고 https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md 에서 스킬을 추가해
2. .agents/skills/swag2mcp-cli/ 디렉토리를 만들고 https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md 에서 스킬을 추가해
```

에이전트가 두 스킬 파일을 다운로드하여 올바른 위치에 배치합니다.

## 수동 설치

LLM 클라이언트가 에이전트 기반 설정을 지원하지 않는 경우 파일을 수동으로 다운로드하세요:

```bash
mkdir -p .agents/skills/swag2mcp-format
mkdir -p .agents/skills/swag2mcp-cli

curl -o .agents/skills/swag2mcp-format/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md
curl -o .agents/skills/swag2mcp-cli/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md
```

## LLM 클라이언트 구성

OpenCode의 경우 `opencode.json`에 스킬을 추가하세요:

```json
{
  "skills": [
    {
      "name": "swag2mcp-format",
      "sourceURL": "https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md"
    },
    {
      "name": "swag2mcp-cli",
      "sourceURL": "https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md"
    }
  ]
}
```

## 재시작 필요

**스킬을 추가한 후 LLM 클라이언트 또는 IDE를 재시작하세요.** 일부 도구는 시작 시에만 스킬을 로드합니다. 스킬이 적용되지 않는 것 같으면 다음을 시도하세요:

- **OpenCode**: 애플리케이션을 다시 시작하거나 opencode 명령을 다시 실행
- **Cursor**: 창을 닫고 다시 열기 (`Cmd+Shift+W` / `Ctrl+Shift+W`)
- **Claude Desktop**: 앱을 종료하고 다시 실행
- **VS Code**: 창 다시 로드 (`Ctrl+Shift+P` → "Developer: Reload Window")
