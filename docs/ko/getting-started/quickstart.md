# 빠른 시작

2분 만에 swag2mcp를 실행하세요.

## 1. 초기화

### 홈 디렉토리 (권장)

전체 시스템을 위한 일회성 설정입니다. 설정이 홈 폴더에 저장됩니다.

::: code-group

```bash [macOS / Linux]
swag2mcp init
# ~/.swag2mcp/swag2mcp.yaml 생성
```

```powershell [Windows]
swag2mcp.exe init
# %USERPROFILE%\.swag2mcp\swag2mcp.yaml 생성
```

:::

### 프로젝트 디렉토리

프로젝트 내부의 격리된 워크스페이스용입니다.

::: code-group

```bash [macOS / Linux]
mkdir -p ./swag2mcp && swag2mcp init ./swag2mcp
```

```powershell [Windows]
mkdir ./swag2mcp; swag2mcp.exe init ./swag2mcp
```

:::

### ZIP에서

준비된 워크스페이스가 있는 경우(예: 동료로부터):

```bash
swag2mcp import --from-zip workspace.zip
```

## 2. 에이전트 스킬 설치 (권장)

swag2mcp 스킬을 설치하여 AI 에이전트에게 모든 명령어, 플래그, 설정 형식, 실제 예시를 가르치세요.

에이전트에게 요청:

```bash
"https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md에서 swag2mcp-cli 스킬을 추가하세요"
"https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md에서 swag2mcp-format 스킬을 추가하세요"
```

> 일부 IDE는 스킬 추가 후 재시작이 필요합니다.

## 3. LLM 클라이언트 / IDE 설정

swag2mcp에 연결하도록 IDE를 설정하세요. IDE는 필요할 때 MCP 서버를 자동으로 시작합니다.

::: code-group

```json [OpenCode]
{
  "mcp": {
    "swag2mcp": {
      "type": "local",
      "command": ["swag2mcp", "mcp"],
      "enabled": true
    }
  }
}
```

```json [Claude Desktop]
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp"]
    }
  }
}
```

```json [Crush]
{
  "mcp": {
    "swag2mcp": {
      "type": "stdio",
      "command": "swag2mcp",
      "args": ["mcp"]
    }
  }
}
```

:::

다른 IDE(Cursor, VS Code, JetBrains)는 [통합 가이드](../integration/opencode.md)를 참조하세요.

> 커스텀 경로(예: `./swag2mcp`)에 워크스페이스를 초기화한 경우 명령어에 전체 경로를 사용하세요:
> `"command": ["swag2mcp", "mcp", "/absolute/path/to/swag2mcp"]`

> **설정 변경 후에는 변경 사항을 적용하려면 MCP 서버를 다시 시작하세요.**

## 4. MCP 서버 시작

### stdio (기본값) — 로컬 IDE용

설정할 것이 없습니다. IDE가 위 설정을 통해 swag2mcp를 자동으로 시작합니다.

```bash
swag2mcp mcp
```

### SSE / Streamable HTTP — 원격 접근용

```bash
swag2mcp mcp --transport sse --http-addr :8080
```

또는 `swag2mcp.yaml`에서 설정:

```yaml
mcp:
  transport: sse
  addr: ":8080"
  path: "/mcp"
```

모든 플래그는 [MCP 서버 참조](../configuration/mcp-server.md)를 참조하세요.

### 태그로 spec 필터링

```bash
swag2mcp mcp --tags weather,public
```

일치하는 태그가 있는 spec만 LLM이 사용할 수 있습니다.

### 작동 확인

연결 후 LLM 에이전트에게 물어보세요:

```bash
"어떤 MCP 도구를 지원하나요?"
```

에이전트가 swag2mcp 도구(`spec_list`, `search`, `invoke` 등)를 나열하면 모든 것이 정상입니다.

### 시도해 볼 예시 쿼리

| 에이전트에게 물어보기 | 결과 |
|-----------------|------|
| "뉴욕 날씨는 어떤가요?" | `invoke` — Open-Meteo 예보 API 호출 |
| "현재 BTC 가격은 얼마인가요?" | `invoke` — Binance 티커 API 호출 |
| "아재개그 해주세요" | `invoke` — icanhazdadjoke API 호출 |
| "피카츄 보여주세요" | `invoke` — 이름으로 PokéAPI 호출 |
| "릭 산체스는 누구인가요?" | `invoke` — Rick and Morty 캐릭터 API 호출 |
| "베이징의 대기질은 어떤가요?" | `invoke` — Open-Meteo 대기질 API 호출 |
| "포르투갈 근처 파도는 얼마나 높나요?" | `invoke` — Open-Meteo 해양 API 호출 |
| "개에 관한 농담 검색" | `invoke` — dadjoke 검색 엔드포인트 호출 |
| "모든 포켓몬 나열" | `invoke` — PokéAPI 목록 엔드포인트 호출 |
| "에베레스트 산의 고도는?" | `invoke` — Open-Meteo 고도 API 호출 |

## 5. 다음 단계

- [개념](../concepts/overview.md) — 아키텍처 이해
- [설정](../configuration/config-file.md) — 설정 사용자 정의
- [CLI 명령어](../cli/overview.md) — 전체 명령어 참조
