# swag2mcp

<div style="background: #dc2626; color: white; padding: 20px 24px; border-radius: 12px; text-align: center; font-size: 1.4em; font-weight: 700; margin: 24px 0;">
  🚧 작업 진행 중 — 출시 예정!
</div>

OpenAPI/Swagger/Postman API 명세를 Model Context Protocol(MCP)을 통해 LLM 에이전트와 연결합니다.

<a href="https://www.youtube.com/watch?v=1Da4UmE2f9U" target="_blank">
  <img src="https://raw.githubusercontent.com/mmadfox/swag2mcp/main/docs/cover.png" alt="Preview">
</a>

## 당신의 API가 LLM과 대화합니다

한 줄의 설정으로 모든 OpenAPI/Swagger/Postman 파일을 MCP 서버로 변환합니다. LLM 에이전트가 API를 발견, 검사, 호출합니다 — 통합 코드가 전혀 필요 없습니다.

<img src="/architecture.svg" width="700" alt="swag2mcp architecture">

## 래퍼를 작성하지 마세요

새 API를 LLM에 연결할 때마다 동일한 보일러플레이트를 작성합니다: 명세 파싱, 인증, 오류 처리, 속도 제한. swag2mcp가 대신 처리합니다 — 19개의 즉시 사용 가능한 MCP 도구.

## 이런 분들에게 필요합니다

| 역할 | 이유 |
|------|------|
| **AI 에이전트 개발자** | 2일이 아닌 2분 만에 API 연결 |
| **MCP 엔지니어** | 핸들러 코드 불필요 — 명세만 가리키면 됩니다 |
| **아키텍트** | 회사 내 모든 LLM을 위한 단일 API 통합 레이어 |
| **데이터 분석가** | 코딩 없이 자연어로 API 접근 |
| **DevOps / SRE** | 추가 서비스 없이 LLM을 통한 모니터링 및 자동화 |
| **통합 엔지니어** | Basic부터 OAuth2, HMAC까지 9가지 인증 방식 기본 제공 |
| **QA 엔지니어** | 실제 API 없이 격리된 테스트를 위한 모의 서버 |
| **제품 관리자** | 백엔드 작업 없이 빠른 AI 기능 프로토타입 |
| **그 외 많은 분들** | |

---

## 라이선스

**GNU Affero General Public License v3.0**(AGPL v3)에 따라 라이선스가 부여됩니다.

전체 라이선스 텍스트는 [LICENSE](https://github.com/mmadfox/swag2mcp/blob/main/LICENSE)를 참조하세요.

```
SPDX-License-Identifier: AGPL-3.0-only
```
