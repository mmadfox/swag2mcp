# 프로젝트 구조

```
swag2mcp/
├── cmd/
│   ├── swag2mcp/          # 메인 바이너리
│   │   └── main.go
│   └── swag2mcp-mock/     # 모의 서버
│       └── main.go
├── internal/
│   ├── auth/              # 9가지 인증 방법
│   ├── cache/             # 명세 캐싱
│   ├── commands/          # 13개 CLI 명령어 (cobra)
│   ├── config/            # YAML 설정
│   ├── env/               # 환경 변수
│   ├── httpclient/        # HTTP 클라이언트
│   ├── id/                # MD5 ID 생성
│   ├── index/             # 전문 검색 (bluge)
│   ├── model/             # 데이터 모델
│   ├── reader/            # 대용량 응답 읽기
│   ├── server/
│   │   ├── mcp/           # MCP 서버 (19개 도구)
│   │   └── mockserver/    # 모의 서버
│   ├── service/           # 비즈니스 로직
│   ├── spec/              # 명세 파서
│   ├── tui/               # TUI 인터페이스
│   └── workspace/         # 워크스페이스 관리
├── specs/                 # 샘플 명세
├── tests/                 # 통합 테스트
├── docs/                  # 문서
├── examples/              # 설정 예시
└── playground/            # 개발 샌드박스
```

## 주요 패키지

| 패키지 | 설명 |
|--------|------|
| `auth` | 9가지 인증 방법 |
| `cache` | TTL이 있는 디스크 기반 캐싱 |
| `commands` | Cobra CLI 명령어 |
| `config` | 계단식 YAML 설정 |
| `httpclient` | 설정 가능한 HTTP 클라이언트 |
| `index` | 전문 검색 (bluge) |
| `server/mcp` | MCP 서버 (3가지 전송 방식) |
| `service` | 비즈니스 로직 (핵심) |
| `spec` | OpenAPI/Swagger/Postman 파서 |
| `tui` | Bubbletea TUI |
| `workspace` | 파일 관리 |
