# swag2mcp-cli

**swag2mcp-cli** 스킬은 LLM에게 완전한 swag2mcp CLI 참조를 제공합니다 — 모든 명령, 플래그, 인수 및 설정 옵션. 이 스킬을 통해 LLM은 "어떻게..." 질문에 추측 없이 정확하게 답변할 수 있습니다.

## 적용 범위

13개 CLI 명령 모두:

| 명령 | 목적 |
|------|------|
| `init` | 작업 공간 및 설정 초기화 |
| `add` | 사양 또는 컬렉션 추가 |
| `delete` | 사양 또는 컬렉션 삭제 |
| `ls` | 구성된 사양 나열 |
| `run` | API Explorer TUI 시작 |
| `validate` | 설정 파일 검증 |
| `clean` | 캐시된 데이터 정리 |
| `update` | 설정에서 캐시 업데이트 |
| `mcp` | MCP 서버 시작 |
| `version` | 버전 정보 표시 |
| `info` | 런타임 정보 표시 |
| `import` | 명세 파일 가져오기 또는 백업에서 워크스페이스 복원 |
| `export` | 작업 공간을 ZIP 파일로 내보내기 |

모든 플래그, 설정 파일 구조, 인증 방법 및 고급 옵션 포함.

## 직접 링크

<https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md>

## LLM 에이전트를 통한 설치

다음 요청을 AI 기반 IDE에 복사하세요:

```
.agents/skills/swag2mcp-cli/ 디렉토리를 만들고
https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md 에서 스킬을 추가해
```

## 수동 설치

```bash
mkdir -p .agents/skills/swag2mcp-cli
curl -o .agents/skills/swag2mcp-cli/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md
```

## 재시작 필요

스킬을 추가한 후 LLM 클라이언트 또는 IDE를 재시작하세요（[개요](overview.md#재시작-필요) 참조）。
