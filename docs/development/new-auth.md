# Adding a New Auth Method

## Steps

1. **Create the auth client** in `internal/auth/&lt;name&gt;.go`
2. **Implement the `Authenticator` interface**
3. **Add type constant** to `internal/auth/auth.go`
4. **Add YAML decoder** to `internal/config/auth.go`
5. **Register decoder** in the `authDecoders` map
6. **Write tests**

## 1. Auth client

Create `internal/auth/my_auth.go`:

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

## 2. Authenticator interface

Every auth client must implement:

```go
type Authenticator interface {
    New() error                    // Initialize, resolve env vars
    Type() Type                    // Return the auth type identifier
    Apply(req *http.Request, out *Info) error  // Apply auth to request
    Validate() error               // Validate required fields
}
```

## 3. Type constant

Add to `internal/auth/auth.go`:

```go
const MyAuth Type = "my-auth"
```

## 4. YAML decoder

Add a decoder function in `internal/config/auth.go`. The decoder receives a `*yaml.Node` and must decode it into your auth client struct:

```go
func decodeMyAuth(node *yaml.Node) (auth.Authenticator, error) {
    var client auth.MyAuthClient
    if err := decodeConfig(node, &client); err != nil {
        return nil, err
    }
    return &client, nil
}
```

The `decodeConfig` helper handles the common pattern: it checks that the node is not empty, decodes YAML into the struct, and returns a descriptive error on failure.

## 5. Register decoder

Add your decoder to the `authDecoders` map in `internal/config/auth.go`:

```go
var authDecoders = map[string]authDecoder{
    // ... existing decoders
    auth.MyAuth.String(): decodeMyAuth,
}
```

The `UnmarshalYAML` method on `Auth` reads the `type` field from the YAML, normalises underscores to hyphens, looks up the decoder in `authDecoders`, and calls it with the `config` node. This is how swag2mcp knows which auth client to instantiate for each spec.

## 6. Tests

Create `internal/auth/my_auth_test.go` with table-driven tests covering:

- `New()` resolves env vars correctly
- `Type()` returns the correct type
- `Apply()` sets the right headers/query params
- `Apply()` handles empty values gracefully
- `Validate()` passes for valid config
- `Validate()` fails for missing required fields
