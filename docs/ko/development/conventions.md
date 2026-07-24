# 코드 규칙

## Go

- **Go 1.26+**
- **gofmt** / **gofumpt** / **goimports** / **gci**
- 줄당 **120자**
- 중첩 if 대신 **Guard clauses**
- **이름**: 비공개는 `camelCase`, 내보내기는 `PascalCase`

## 오류

LLM에 표시되는 오류에는 `LLMError`를 사용하세요:

```go
type LLMError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}
```

오류 코드:
- `validation_failed` — 잘못된 매개변수
- `not_found` — 리소스를 찾을 수 없음
- `rate_limit` — 속도 제한 초과
- `invoke_error` — API 호출 오류

## 인터페이스

- 작은 인터페이스 (1-3개 메서드)
- 인터페이스 구성
- 설정을 위한 함수형 옵션

## 테스트

- 테이블 기반 테스트
- 테스트 헬퍼 (`newTestService()`, `seedTestData()`)
- `go.uber.org/mock`을 통한 모의
- 핵심 패키지 80%+ 커버리지

## 설정

- YAML 형식
- 계단식: 전역 → spec → collection
- `go-playground/validator`를 통한 검증
- `$(VAR)`을 통한 환경 변수
