# Adding a New Auth Method

## Steps

1. Create `internal/auth/<name>.go`
2. Implement the `Authenticator` interface
3. Add type to `internal/auth/auth.go`
4. Add parsing to `internal/config/auth.go`
5. Write tests

## Authenticator Interface

```go
type Authenticator interface {
    Type() Type
    SetHeaders(req *http.Request) error
    SetQueryParams(req *http.Request) error
    Info() Info
}
```

## Example

```go
package auth

type MyAuth struct {
    token string
}

func (a *MyAuth) Type() Type {
    return TypeMyAuth
}

func (a *MyAuth) SetHeaders(req *http.Request) error {
    req.Header.Set("X-My-Auth", a.token)
    return nil
}

func (a *MyAuth) SetQueryParams(req *http.Request) error {
    return nil
}

func (a *MyAuth) Info() Info {
    return Info{
        Type:   TypeMyAuth,
        Fields: map[string]string{"token": a.token},
    }
}
```

## Configuration

In `internal/config/auth.go`:

```go
type MyAuthConfig struct {
    Token string `yaml:"token" validate:"required"`
}
```
