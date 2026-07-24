# 새 인증 방법 추가

## 단계

1. **인증 클라이언트 생성** `internal/auth/&lt;name&gt;.go`
2. **`Authenticator` 인터페이스 구현**
3. **타입 상수 추가** `internal/auth/auth.go`
4. **YAML 디코더 추가** `internal/config/auth.go`
5. **`authDecoders` 맵에 디코더 등록**
6. **테스트 작성**

## 1. 인증 클라이언트

`internal/auth/my_auth.go` 생성:

```go
package auth

import "net/http"

type MyAuthClient struct {
    Token string `yaml:"token" validate:"required"`
}

func (c *MyAuthClient) New() error {
    c.Token = resolveEnv(c.Token)
    return nil
}

func (c *MyAuthClient) Type() Type {
    return MyAuth
}

func (c *MyAuthClient) Apply(req *http.Request, out *Info) error {
    if c.Token == "" {
        return nil
    }
    setAuthHeader(req, out, "X-My-Auth", c.Token)
    return nil
}

func (c *MyAuthClient) Validate() error {
    return authValidator.Struct(c)
}
```

## 2. Authenticator 인터페이스

모든 인증 클라이언트는 다음을 구현해야 합니다:

```go
type Authenticator interface {
    New() error                    // 초기화, 환경 변수 해결
    Type() Type                    // 인증 유형 식별자 반환
    Apply(req *http.Request, out *Info) error  // 요청에 인증 적용
    Validate() error               // 필수 필드 검증
}
```

## 3. 타입 상수

`internal/auth/auth.go`에 추가:

```go
const MyAuth Type = "my-auth"
```

## 4. YAML 디코더

`internal/config/auth.go`에 디코더 함수를 추가합니다. 디코더는 `*yaml.Node`를 받아 인증 클라이언트 구조체로 디코딩해야 합니다:

```go
func decodeMyAuth(node *yaml.Node) (auth.Authenticator, error) {
    var client auth.MyAuthClient
    if err := decodeConfig(node, &client); err != nil {
        return nil, err
    }
    return &client, nil
}
```

`decodeConfig` 헬퍼는 일반적인 패턴을 처리합니다: 노드가 비어 있지 않은지 확인하고, YAML을 구조체로 디코딩하며, 실패 시 설명적인 오류를 반환합니다.

## 5. 디코더 등록

`internal/config/auth.go`의 `authDecoders` 맵에 디코더를 추가하세요:

```go
var authDecoders = map[string]authDecoder{
    // ... 기존 디코더
    auth.MyAuth.String(): decodeMyAuth,
}
```

`Auth`의 `UnmarshalYAML` 메서드는 YAML에서 `type` 필드를 읽고, 밑줄을 하이픈으로 정규화하며, `authDecoders`에서 디코더를 조회하고, `config` 노드로 호출합니다. 이것이 swag2mcp가 각 spec에 대해 인스턴스화할 인증 클라이언트를 알 수 있는 방법입니다.

## 6. 테스트

`internal/auth/my_auth_test.go`를 테이블 기반 테스트로 생성하여 다음을 다룹니다:

- `New()`가 환경 변수를 올바르게 해결하는지
- `Type()`이 올바른 타입을 반환하는지
- `Apply()`가 올바른 헤더/쿼리 매개변수를 설정하는지
- `Apply()`가 빈 값을 적절히 처리하는지
- `Validate()`가 유효한 설정에 대해 통과하는지
- `Validate()`가 필수 필드 누락에 대해 실패하는지
