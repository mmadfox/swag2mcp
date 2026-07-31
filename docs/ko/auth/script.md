# Script Auth

## 목적

외부 스크립트를 통한 인증. swag2mcp는 두 가지 스크립트 유형을 지원합니다:
- **`.sh`** — Linux 및 macOS의 셸 스크립트
- **`.bat`** — Windows의 배치 스크립트

스크립트는 작업 공간의 `auth_scripts` 디렉토리에 배치되어야 하며 JSON 토큰을 stdout에 출력해야 합니다.

## 사용 시기

- 커스텀 또는 비표준 인증 체계
- 복잡한 토큰 획득 로직 (다단계, 추가 검사 포함)
- 표준 방법이 요구 사항에 맞지 않을 때

## 설정

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: script
      config:
        domain: "jokes"
```

## 매개변수

| 매개변수 | 필수 | 설명 |
|---------|------|------|
| `domain` | 예 | 스크립트 파일 이름 (확장자 제외). spec의 `domain`과 일치해야 합니다 — 예: `domain: jokes`인 경우 스크립트 파일은 `jokes.sh`(Unix) 또는 `jokes.bat`(Windows)여야 합니다. |

## 스크립트 위치

스크립트는 워크스페이스의 `auth_scripts` 디렉토리에 있어야 합니다:

- **Linux / macOS:** `{workspace}/auth_scripts/{domain}.sh`
- **Windows:** `{workspace}/auth_scripts/{domain}.bat`

## 스크립트 출력 형식

스크립트는 stdout에 토큰과 만료 시간이 포함된 JSON을 출력해야 합니다:

```bash
#!/bin/bash
# auth_scripts/jokes.sh

TOKEN=$(curl -s -X POST https://auth.example.com/token \
  -d "grant_type=client_credentials" \
  -d "client_id=$CLIENT_ID" \
  -d "client_secret=$CLIENT_SECRET" | jq -r '.access_token')

echo "{\"token\": \"$TOKEN\", \"expires_in\": 3600}"
```

### JSON 필드

| 필드 | 필수 | 설명 |
|------|------|------|
| `token` | 예 | 인증 토큰 |
| `expires_in` | 아니요 | 토큰 수명(초) (기본값: 3600) |

## 참고 사항

- swag2mcp는 캐시된 토큰이 만료된 경우 모든 요청에서 스크립트를 실행합니다
- 스크립트는 30초 이내에 완료되어야 합니다
- 토큰은 만료 시간까지 캐시됩니다
- 스크립트 파일 이름 = `{domain}.sh` (Unix) 또는 `{domain}.bat` (Windows)
- `domain`에는 `/` 또는 `\`를 포함할 수 없습니다
