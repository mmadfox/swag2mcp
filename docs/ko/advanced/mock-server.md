# 모의 서버

## 개요

모의 서버는 OpenAPI 스키마를 기반으로 가짜 API 응답을 생성합니다. 실제 HTTP 호출 없이 API 통합을 테스트할 수 있습니다. 개발, LLM 에이전트 테스트, 데모에 유용합니다.

모의 서버는 **별도의 바이너리** — `swag2mcp-mock`입니다. 메인 `swag2mcp` 바이너리에 포함되어 있지 않으며 별도로 설치해야 합니다.

## 설치

```bash
# 옵션 1: GitHub Releases에서 다운로드
# swag2mcp-mock_<version>_<os>_<arch>.tar.gz 찾기

# 옵션 2: Go로 설치
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

## 설정

설정에서 모의 서버를 활성화하세요:

```yaml
mock_enabled: true
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092

specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
        base_mock_url: "127.0.0.1:9090"
```

## 매개변수

### mock_enabled

- **타입:** `bool`
- **기본값:** `false`
- **효과:** `true`일 때 모든 활성 collection에 `base_mock_url`이 설정되어야 합니다. 모의 서버는 각 collection에 대해 HTTP 서버를 시작합니다.

### mock_auth

모의 인증 서버의 포트입니다. 실제 자격 증명 없이 인증된 API를 테스트할 수 있도록 OAuth2, Digest, HMAC 인증 엔드포인트를 시뮬레이션합니다.

| 필드 | 기본값 | 설명 |
|------|--------|------|
| `oauth2_port` | `9090` | 모의 OAuth2 토큰 서버 포트 |
| `digest_port` | `9091` | 모의 Digest 인증 서버 포트 |
| `hmac_port` | `9092` | 모의 HMAC 인증 서버 포트 |

### base_mock_url (collection별)

- **타입:** `string`
- **필수:** 예 (`mock_enabled: true`일 때)
- **형식:** `host:port` (예: `localhost:8080`, `127.0.0.1:9000`)
- **효과:** 각 collection은 이 주소에서 자체 HTTP 서버를 얻습니다. 서버는 명세에 정의된 모든 엔드포인트에 무작위로 생성된 데이터로 응답합니다.

## 모의 서버 시작

```bash
# 기본 설정으로 시작
swag2mcp-mock mockserver

# TLS로 시작
swag2mcp-mock mockserver --tls

# 커스텀 TLS 인증서로 시작
swag2mcp-mock mockserver --tls --tls-cert cert.pem --tls-key key.pem
```

### TLS 플래그

| 플래그 | 설명 |
|-------|------|
| `--tls` | 자체 서명 인증서로 TLS 활성화 |
| `--tls-cert` | TLS 인증서 파일 경로 |
| `--tls-key` | TLS 키 파일 경로 |

`--tls`가 `--tls-cert`와 `--tls-key` 없이 설정되면 `localhost`용 자체 서명 인증서가 자동으로 생성됩니다.

## 모의 서버의 기능

모의 서버를 시작하면:

1. **모든 명세 파일 파싱** — 각 collection의 OpenAPI/Swagger 명세를 읽습니다
2. **핸들러 등록** — 명세에 정의된 모든 경로와 메서드에 대한 HTTP 핸들러를 생성합니다
3. **가짜 데이터 생성** — 응답 스키마와 일치하는 무작위 생성 데이터(올바른 타입, 형식, 구조)로 응답합니다
4. **인증 서버 시작** — 테스트를 위해 OAuth2, Digest, HMAC 인증 엔드포인트를 시뮬레이션합니다

### 모의 테스트

```bash
# 한 터미널에서:
swag2mcp-mock mockserver

# 다른 터미널에서:
curl http://localhost:8080/pets
# → [{"id":1,"name":"Pet_name","status":"available"}]
```

## 가짜 데이터 생성 방식

모의 서버는 OpenAPI 스키마를 기반으로 현실적인 가짜 데이터를 생성합니다:

- **문자열** — 무작위 단어, 문장 또는 형식별 값 (이메일, URL, UUID, 날짜, 전화번호 등)
- **숫자** — 지정된 범위 내의 무작위 정수 및 실수
- **부울** — 무작위 true/false
- **배열** — 1~3개의 무작위 항목
- **객체** — 모든 속성이 무작위 값으로 채워짐
- **열거형** — 열거형 목록에서 무작위 값
- **Nullable 필드** — 때때로 `null` 반환 (~10% 확률)

## 사용 사례

- **개발** — 실제 API 접근 없이 통합 테스트
- **LLM 에이전트 테스트** — LLM이 엔드포인트를 발견, 검사, 호출할 수 있는지 확인
- **데모** — 실제 API를 설정하지 않고 swag2mcp 작동 시연
- **부하 테스트** — 실제 API에 접근하지 않고 MCP 서버 부하 테스트

## 중요 참고 사항

- **별도 바이너리** — `swag2mcp-mock`은 메인 `swag2mcp` 바이너리에 포함되어 있지 않습니다. 별도로 설치하세요.
- **각 collection은 자체 포트를 가짐** — collection별로 `base_mock_url`을 설정하세요
- **인증 모의 서버는 전역** — OAuth2, Digest, HMAC 서버는 collection 수와 관계없이 설정된 포트에서 실행됩니다
- **명세 파싱 실패는 치명적이지 않음** — collection의 명세를 파싱할 수 없으면 경고와 함께 건너뜁니다
- **자체 서명 TLS** — 인증서 없이 `--tls`를 사용하면 localhost용 자체 서명 인증서가 생성됩니다
