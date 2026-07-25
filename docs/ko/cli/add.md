# add

## 목적

기존 설정에 새 **spec**(API 서비스) 또는 **collection**(OpenAPI/Swagger/Postman 파일)을 추가합니다. 새 API로 워크스페이스를 확장하는 기본 방법입니다.

## 사용 시기

- LLM 에이전트에 연결할 새 API가 있을 때
- OpenAPI 명세 URL을 찾아 추가하려고 할 때
- 기존 spec에 추가 명세 파일(collection)을 추가하려고 할 때
- 대화형 마법사 대신 직접 YAML을 작성하는 것을 선호할 때

## 구문

```bash
swag2mcp add spec [path] [flags]
swag2mcp add collection [path] [flags]
```

## 인수

| 인수 | 위치 | 필수 | 설명 |
|------|------|------|------|
| `path` | 1 | 아니요 | 워크스페이스 디렉토리. 생략 시 경로 해결 규칙에 따라 결정됩니다. |

## 플래그

### `add spec`

| 플래그 | 약어 | 타입 | 기본값 | 설명 |
|-------|------|------|--------|------|
| `--yaml` | `-y` | `string` | `""` | YAML 입력 인라인 또는 `-` (stdin) |
| `--example` | `-e` | `bool` | `false` | YAML 템플릿 출력 후 종료 |

### `add collection`

| 플래그 | 약어 | 타입 | 기본값 | 설명 |
|-------|------|------|--------|------|
| `--yaml` | `-y` | `string` | `""` | YAML 입력 인라인 또는 `-` (stdin) |
| `--example` | `-e` | `bool` | `false` | YAML 템플릿 출력 후 종료 |

## 작동 방식

### 대화형 모드 (기본값)

spec 또는 collection 필드를 단계별로 입력할 수 있는 TUI 마법사를 실행합니다.

```bash
swag2mcp add spec
swag2mcp add collection
```

### YAML 인라인 모드

YAML을 문자열로 직접 전달합니다. **셸 따옴표에 주의하세요** — `:`, `#`, `&`, `{` 같은 특수 문자는 명령어를 깨뜨릴 수 있습니다.

```bash
swag2mcp add spec --yaml 'domain: meteo
llm_title: Open-Meteo API
base_url: https://meteo.swagger.io/v2
collections:
  - llm_title: Main
    location: https://example.com/spec.json'
```

### YAML from stdin (복잡한 YAML에 권장)

파일에서 파이프하거나 heredoc을 사용하여 셸 따옴표 문제를 완전히 피하세요:

```bash
# 파일에서 파이프
cat spec.yaml | swag2mcp add spec --yaml -

# Heredoc
swag2mcp add spec --yaml - <<EOF
domain: my-api
llm_title: My API
llm_instruction: "X & Y # important에 이 API를 사용하세요"
base_url: https://api.example.com/v1
collections:
  - llm_title: Main
    location: https://raw.githubusercontent.com/org/repo/main/spec.yaml
EOF
```

### YAML 템플릿

예상되는 YAML 구조를 출력하고 종료:

```bash
swag2mcp add spec --example
swag2mcp add collection --example
```

## YAML 형식

### Spec

```yaml
domain: meteo
llm_title: Open-Meteo API
llm_instruction: 이 API를 사용하여 애완동물을 관리하세요.
base_url: https://meteo.swagger.io/v2
tags: [public, demo]
auth:
  type: bearer
  config:
    token: $(TOKEN)
collections:
  - llm_title: Open-Meteo Swagger
    location: https://example.com/spec.json
```

### Collection

```yaml
spec_domain: meteo
llm_title: Orders Collection
location: https://example.com/orders.json
```

## 명령 후 검증

```bash
swag2mcp ls [path]
# 새 spec 또는 collection이 목록에 나타나야 함
```

## 세부 사항

- **자동 초기화:** 설정 파일이 없으면 `add`가 자동으로 init 마법사를 먼저 실행합니다. `init`을 별도로 실행할 필요가 없습니다.
- **셸 따옴표:** 인라인 YAML(`--yaml '...'`)은 특수 문자에 취약합니다. 단순한 값 이상의 경우 heredoc 또는 파이프와 함께 `--yaml -`을 선호하세요.
- **`--example`은 즉시 종료**되며 기존 설정을 확인하거나 수정하지 않습니다.
- **`add spec` vs `add collection`:** 새 API 서비스(새 도메인)에는 `add spec`을 사용하세요. 기존 spec에 다른 명세 파일을 추가하려면 `add collection`을 사용하세요.
