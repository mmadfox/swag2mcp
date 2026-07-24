# version

## 목적

swag2mcp 버전을 출력합니다. 설치된 버전 확인, 버그 신고, 호환성 확인에 유용합니다.

## 사용 시기

- 설치된 swag2mcp 버전을 확인하려고 할 때
- 버그를 신고하고 버전을 포함해야 할 때
- 설치가 성공했는지 확인하려고 할 때

## 구문

```bash
swag2mcp version
swag2mcp --version
```

## 인수

없음.

## 플래그

없음.

## 작동 방식

```bash
swag2mcp version
# swag2mcp v1.2.0

swag2mcp --version
# swag2mcp v1.2.0
```

## 출력 형식

```
swag2mcp <version>
```

버전은 빌드 시 `ldflags`를 통해 설정됩니다. 설정되지 않으면 기본값은 `"dev"`입니다.

## 세부 사항

- **두 가지 형식:** `swag2mcp version`(하위 명령어)과 `swag2mcp --version`(전역 플래그) 모두 동일한 출력을 생성합니다.
- **설정 불필요:** 이 명령어는 워크스페이스나 설정 파일 없이 작동합니다.
