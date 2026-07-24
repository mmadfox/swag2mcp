# CLI 워크플로우

이 페이지는 초기화부터 일상적인 작업까지 터미널에서 swag2mcp를 사용하는 실제 예시를 보여줍니다.

## 빠른 시작

```bash
# 1. 워크스페이스 초기화
mkdir -p .swag2mcp && swag2mcp init ./.swag2mcp

# 2. spec 나열
swag2mcp ls
```

## YAML로 spec 추가

### 간단한 spec (공개 API)

```bash
swag2mcp add spec --yaml - <<EOF
domain: meteo
llm_title: Open-Meteo Weather API
base_url: https://api.open-meteo.com
collections:
  - llm_title: Weather Forecast
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
EOF
```

### 인증이 있는 spec (env의 bearer 토큰)

```bash
swag2mcp add spec --yaml - <<EOF
domain: my-api
llm_title: My Protected API
base_url: https://api.example.com/v1
auth:
  type: bearer
  config:
    token: \$(MY_TOKEN)
collections:
  - llm_title: Users
    location: https://raw.githubusercontent.com/my-org/my-api/main/users.yaml
EOF
```

### 여러 collection이 있는 spec

```bash
swag2mcp add spec --yaml - <<EOF
domain: meteo
llm_title: Open-Meteo APIs
base_url: https://api.open-meteo.com
collections:
  - llm_title: Forecast
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
  - llm_title: Air Quality
    base_url: https://air-quality-api.open-meteo.com
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
EOF
```

## 기존 spec에 collection 추가

```bash
swag2mcp add collection --yaml - <<EOF
spec_domain: meteo
llm_title: Marine Weather
location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
EOF
```

## Spec 나열

```bash
$ swag2mcp ls
Specifications:
  dadjoke (https://icanhazdadjoke.com)
    jokes (3 endpoints)
  meteo (https://api.open-meteo.com)
    forecast (5 endpoints)
    air-quality (8 endpoints)
    marine (4 endpoints)
```

### 태그로 필터링

```bash
swag2mcp ls --tags=public
```

## 런타임 정보 보기

```bash
$ swag2mcp info
{
  "version": "v1.2.0",
  "workspace": "/home/user/.swag2mcp",
  "specs": {
    "total": 2,
    "active": 2,
    "disabled": 0,
    "collections": 4,
    "endpoints": 20
  },
  "http_client": {
    "timeout": "30s",
    "max_response_size": "2 KB",
    "follow_redirects": true
  },
  "mcp": {
    "transport": "stdio"
  },
  "auth": {
    "methods": ["bearer"]
  }
}
```

## 설정 검증

```bash
$ swag2mcp validate
✅ Configuration is valid.
✓ Spec dadjoke: OK
✓ Spec meteo: OK
```

## MCP 서버 시작

### stdio (IDE 통합용)

```bash
swag2mcp mcp
```

### HTTP (원격 접근용)

```bash
swag2mcp mcp --transport sse --http-addr :8080
```

### 태그 필터 사용

```bash
swag2mcp mcp --tags=public
```

## Spec 업데이트

모든 캐시된 명세 파일 새로고침:

```bash
swag2mcp update
```

## 캐시 정리

```bash
swag2mcp clean
```

## 내보내기 및 가져오기

### 워크스페이스 백업

```bash
swag2mcp export --output ~/backups/swag2mcp-2026-07-24.zip
```

### 다른 머신에서 복원

```bash
# 새 머신에서
swag2mcp import --from-zip swag2mcp-2026-07-24.zip
```

## 대화형 TUI 탐색기

```bash
swag2mcp run
```

전체 화면 터미널 UI가 열려 API를 검색, 탐색, 호출할 수 있습니다.

## 모의 서버

```bash
# 모의 바이너리 설치
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest

# 모의 서버 시작
swag2mcp-mock mockserver
```
