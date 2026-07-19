---
name: godeveloper
description: "Go development conventions, patterns, and best practices for the swag2mcp project. Covers naming, code organization, error handling, interface design, concurrency, testing, configuration, and project structure. Use when writing or reviewing Go code, designing packages, setting up tests, or structuring services."
license: MIT
metadata:
  author: mmadfox
  version: "1.0.0"
---

# Go Code Generation Skill — swag2mcp

This document defines how to write Go code for the **swag2mcp** project. It synthesizes the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md), [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments), [Go Names](https://talks.golang.org/2014/names.slide#1), and [Effective Go](https://go.dev/doc/effective_go), adapted to this project's conventions.

---

## 1. Naming Conventions

### 1.1 General Rules

- **MixedCase** — never use underscores in names (`ServeHTTP`, not `Serve_HTTP`; `maxLength`, not `max_length`).
- **Acronyms/Initialisms** — keep consistent case: `ServeHTTP`, `IDProcessor`, `urlPony`, `appID`. Never `ServeHttp`, `UrlPony`, `AppId`. Exception: generated protobuf code.
- **Distance rule** — the further a name is from its declaration, the more descriptive it must be. Local variables with small scope get short names; exported package-level names get descriptive names.

### 1.2 Packages

- Short, lowercase, single-word names: `auth`, `cache`, `config`, `index`, `spec`, `types`.
- No underscores or mixed caps.
- Avoid meaningless names: `util`, `common`, `misc`, `api`, `types`, `interfaces`.
- Package name must match the last component of its import path.
- Package name should lend meaning to the names it exports. `auth.Basic` not `auth.BasicAuth`.

### 1.3 Files

- Lowercase with underscores for compound names: `oauth2_cc.go`, `api_key.go`, `parse_v3.go`.
- Test files: `*_test.go` suffix.
- Generated files: `mock_*_test.go` prefix (e.g., `mock_svc_test.go`).
- `doc.go` for package documentation.

### 1.4 Types and Structs

- PascalCase for exported types: `LLMError`, `InvokeRequest`, `EndpointSearchItem`.
- camelCase for unexported types.
- Error types: PascalCase with `Error` suffix: `LLMError`, `ValidationError`.

### 1.5 Interfaces

- Single-method interfaces: name after the method + `er` suffix: `Authenticator`, `TokenURLSetter`, `MockBaseURLSetter`.
- Multi-method interfaces: descriptive name: `svc` (service interface for MCP handlers).
- Interfaces belong in the **consumer** package, not the implementor package.
- Do not define interfaces "for mocking" on the producer side. Return concrete types; let consumers define their own mock interfaces.

### 1.6 Functions and Methods

- PascalCase for exported, camelCase for unexported.
- Constructor functions: `New` prefix: `NewService`, `NewLLMError`, `NewConcurrentChecker`.
- Getters: omit `Get` prefix. `node.Parent()` not `node.GetParent()`.

### 1.7 Variables

- camelCase.
- Short for local scope: `i` over `index`, `r` over `reader`, `b` over `buffer`.
- Descriptive for wider scope.
- Unexported package-level globals: prefix with `_`: `_globalLogger`.

### 1.8 Constants

- PascalCase for exported, camelCase for unexported.
- Enums use `iota` with descriptive names: `validationFailedErrCode`, `notFoundErrCode`.
- **String and duration constants** — every string literal or duration expression (`30*time.Second`, `5*time.Minute`) used in more than one file within a package must be defined as a named constant in the package's main file (e.g. `auth.go`, `config.go`, `service.go`). Each constant must have a doc comment explaining where and how it is used:

```go
// headerAuthorization is the HTTP Authorization header name used by
// bearer, basic, digest, oauth2-cc, oauth2-pwd, and script auth clients.
const headerAuthorization = "Authorization"

// headerValueBearer is the Bearer token prefix used by bearer, oauth2-cc,
// oauth2-pwd, and script auth clients.
const headerValueBearer = "Bearer "

// paramInQuery is the value of the In field for API key auth placed in the
// query string. Used by APIKeyAuthClient.Apply.
const paramInQuery = "query"

// tokenRequestTimeout is the timeout for external HTTP requests
// (token endpoints, digest challenges) and script execution.
const tokenRequestTimeout = 30 * time.Second
```

Rationale: single source of truth prevents typos, enables grep-based auditing, and makes future changes (e.g. renaming a header or adjusting a timeout) safe and local.

### 1.9 Receivers

- One or two characters reflecting the type: `s *Service`, `h *handler`, `c *cache.Cache`.
- Consistent across all methods of a type.
- Never use `me`, `this`, `self`.

### 1.10 Parameters

- Short when types are descriptive: `func Escape(w io.Writer, s []byte)`.
- Descriptive when types are ambiguous: `func Unix(sec, nsec int64) Time`.

### 1.11 Return Values

- Name only for documentation purposes: `func Copy(dst Writer, src Reader) (written int64, err error)`.
- Avoid naming just to enable naked returns. Naked returns are acceptable only in very short functions (a handful of lines).

### 1.12 Errors

- Error types: `FooError` — `type LLMError struct { ... }`.
- Error values: `ErrFoo` — `var ErrNotFound = errors.New("not found")`.
- Error strings: lowercase, no trailing punctuation: `fmt.Errorf("something bad: %w", err)`.

---

## 2. Formatting and Style

### 2.1 Automatic Formatting

All code must pass these formatters (enforced by `.golangci.yml`):
- `gofmt` — standard Go formatting
- `gofumpt` — stricter Go formatting
- `goimports` — import management
- `gci` — import ordering

### 2.2 Import Ordering

Three groups separated by blank lines:
1. **Standard library** (`fmt`, `os`, `net/http`, etc.)
2. **Third-party** (`github.com/...`, `go.uber.org/...`, etc.)
3. **Local module** (`github.com/mmadfox/go/swag2mcp/...`)

```go
import (
    "context"
    "fmt"
    "time"

    "github.com/go-playground/validator/v10"
    "github.com/rs/zerolog"
    "github.com/stretchr/testify/require"

    "github.com/mmadfox/go/swag2mcp/internal/auth"
    "github.com/mmadfox/go/swag2mcp/internal/service"
)
```

### 2.3 Line Length

120-character limit (enforced by `.golangci.yml`). Break lines based on semantics, not length.

### 2.4 Grouping and Ordering

- `type` declarations before `const` before `var`.
- Related types, constants, and functions grouped together.
- Methods grouped by receiver type.

### 2.5 Reduce Nesting

Return early. Keep the normal path at minimal indentation:

```go
// Good
if err != nil {
    return err
}
// normal code

// Bad
if err != nil {
    return err
} else {
    // normal code
}
```

**Guard clause for empty/zero values** — check the negative condition first and return early, so the happy path stays flat:

```go
// Good
func (c *BearerTokenAuthClient) Apply(req *http.Request, out *Info) error {
    if c.Token == "" {
        return nil
    }
    setAuthHeader(req, out, "Authorization", "Bearer "+c.Token)
    return nil
}

// Bad — wraps the entire logic in an if block
func (c *BearerTokenAuthClient) Apply(req *http.Request, out *Info) error {
    if c.Token != "" {
        setAuthHeader(req, out, "Authorization", "Bearer "+c.Token)
    }
    return nil
}
```

**Guard clause for `out == nil`** — when a function accepts an optional output parameter, check `if out == nil { return nil }` early instead of wrapping the entire output logic in `if out != nil { ... }`:

```go
// Good
func (c *BasicAuthClient) Apply(req *http.Request, out *Info) error {
    if c.Username == "" || c.Password == "" {
        return nil
    }
    req.SetBasicAuth(c.Username, c.Password)
    if out == nil {
        return nil
    }
    val := req.Header.Get(headerAuthorization)
    if out.Headers == nil {
        out.Headers = make(map[string]string)
    }
    out.Headers[headerAuthorization] = val
    return nil
}

// Bad — wraps the entire output logic in an if block
func (c *BasicAuthClient) Apply(req *http.Request, out *Info) error {
    req.SetBasicAuth(c.Username, c.Password)
    if out != nil {
        val := req.Header.Get(headerAuthorization)
        if out.Headers == nil {
            out.Headers = make(map[string]string)
        }
        out.Headers[headerAuthorization] = val
    }
    return nil
}
```

### 2.6 Unnecessary Else

If the `if` branch returns/breaks/continues, omit the `else`.

### 2.7 Mutex Granularity

Keep `Lock`/`Unlock` pairs in small, focused methods. Never spread a single lock across a long function with multiple unlock points — it's error-prone and hard to review.

```go
// Good — Lock/Unlock encapsulated in tiny methods, Apply is flat
func (c *ScriptAuthClient) Apply(req *http.Request, out *Info) error {
    if token, ok := c.readCachedToken(); ok {
        setAuthHeader(req, out, headerAuthorization, bearerToken(token))
        return nil
    }
    token, expiresIn, err := c.execute()
    if err != nil {
        return fmt.Errorf("script auth: %w", err)
    }
    c.writeToken(token, expiresIn)
    setAuthHeader(req, out, headerAuthorization, bearerToken(token))
    return nil
}

func (c *ScriptAuthClient) readCachedToken() (string, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    if c.token != "" && time.Now().Before(c.expiresAt) {
        return c.token, true
    }
    return "", false
}

func (c *ScriptAuthClient) writeToken(token string, expiresIn int) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.token = token
    c.expiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
}

// Bad — Lock/Unlock spread across the function, manual unlock in multiple places
func (c *ScriptAuthClient) Apply(req *http.Request, out *Info) error {
    c.mu.Lock()
    if c.token != "" && time.Now().Before(c.expiresAt) {
        setAuthHeader(req, out, headerAuthorization, bearerToken(c.token))
        c.mu.Unlock()
        return nil
    }
    c.mu.Unlock()
    // ... fetch token ...
    c.mu.Lock()
    c.token = token
    c.expiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
    c.mu.Unlock()
    return nil
}
```

### 2.7 Variable Declarations

- `var s string` (zero value) — prefer `var` declaration.
- `s := ""` — use short declaration only when initializing with a non-zero value.
- Prefix unexported globals with `_`: `var _counter int`.

### 2.8 nil is a Valid Slice

```go
var s []string          // preferred (nil slice)
s := []string{}         // avoid (non-nil, zero-length)
```

`nil` slice and empty slice are functionally equivalent (`len` and `cap` are both zero). Use `nil` slice unless JSON encoding requires `[]` over `null`.

### 2.9 Struct and Map Initialization

Use multi-line initialization for complex structs:

```go
obj := Type{
    Field1: value1,
    Field2: value2,
}
```

### 2.10 Raw String Literals

Use backticks for multi-line strings and regex patterns to avoid escaping.

---

## 3. Comments and Documentation

### 3.1 Doc Comments

- Every exported declaration must have a doc comment.
- Format: `// Name sentence.` — starts with the name, ends with a period.
- Package comments in `doc.go` or adjacent to `package` clause with no blank line.

```go
// Package auth provides authentication methods for API specifications.
package auth

// Request represents a request to run a command.
type Request struct { ... }

// Encode writes the JSON encoding of req to w.
func Encode(w io.Writer, req *Request) error { ... }
```

### 3.2 Inline Comments

- `//` style only. No `/* */` block comments except in generated files.
- Complete sentences with proper punctuation.
- `TODO(username): description` for known issues.

### 3.3 Generated Files

Header: `// Code generated by ... DO NOT EDIT.`

---

## 4. Error Handling

### 4.1 Always Handle Errors

Never discard errors with `_`. Check every error return.

### 4.2 Error Wrapping

Use `fmt.Errorf("context: %w", err)` to wrap errors with context. Use `errors.Is` and `errors.As` for unwrapping.

### 4.3 LLMError — Project-Specific Error Type

Use `LLMError` for errors returned to the LLM. It has 4 codes:
- `validation_failed` — invalid input (wrong ID format, missing required fields)
- `not_found` — entity not found in index
- `rate_limit` — per-endpoint 10s cooldown exceeded
- `invoke_error` — HTTP request/response failures

**Messages must explain what to do next:**

```go
func (s *Service) EndpointByID(ctx context.Context, req EndpointByIDRequest) (EndpointByIDResponse, error) {
    if err := s.validateRequest(req); err != nil {
        return EndpointByIDResponse{}, NewLLMError(validationFailedErrCode,
            "The endpoint ID is invalid — it must be a 32-character hex string. "+
            "Use the search tool to find the correct endpoint ID.")
    }
    ep, ok := s.index.EndpointByID(req.ID)
    if !ok {
        return EndpointByIDResponse{}, NewLLMError(notFoundErrCode,
            fmt.Sprintf("No endpoint found with ID %q. Use the search tool to find the correct endpoint.", req.ID))
    }
    return EndpointByIDResponse{...}, nil
}
```

### 4.4 Sentinel Errors

```go
var ErrNotFound = errors.New("not found")
```

### 4.5 Error Naming

- Error types: `FooError`
- Error values: `ErrFoo`
- Error strings: lowercase, no trailing punctuation

### 4.6 Indent Error Flow

```go
x, err := f()
if err != nil {
    return err
}
// use x
```

### 4.7 Don't Panic

Use `error` return values for normal error handling. `panic` only for truly exceptional situations (e.g., programmer bugs, unrecoverable state).

### 4.8 In-Band Errors

Return additional values (`error`, `bool`) instead of using sentinel values like `-1` or `""` to signal errors.

---

## 5. Types and Interfaces

### 5.1 Interface Design

- Interfaces belong in the **consumer** package.
- Prefer small interfaces (1-3 methods).
- Compose larger interfaces from smaller ones:

```go
type svc interface {
    Specs(context.Context) (service.SpecsResponse, error)
    SpecByID(context.Context, service.SpecByIDRequest) (service.SpecByIDResponse, error)
    Search(context.Context, service.SearchRequest) (service.SearchResponse, error)
    Invoke(context.Context, service.InvokeRequest) (service.InvokeResponse, error)
    // ...
}
```

### 5.2 Verify Interface Compliance

Use compile-time checks:

```go
var _ svc = (*Service)(nil)
```

### 5.3 Receiver Type

- If in doubt, use pointer receiver.
- Value receiver for: small immutable structs, basic types, map/func/chan receivers.
- Pointer receiver for: mutation, `sync.Mutex` fields, large structs.
- Never mix receiver types on the same type.

### 5.4 Embedding

- Embed interfaces for composition.
- Embed structs for behavior reuse.
- Do not embed types that don't belong to the public API of the struct.

### 5.5 Generics

Use Go generics for type-safe containers and algorithms:

```go
type Cache[K KeyString, V any] interface {
    Get(key K) (V, bool)
    Set(key K, value V, cost int64) bool
}
```

---

## 6. Concurrency

### 6.1 Goroutine Lifetimes

Always make it clear when and why goroutines exit. Document goroutine lifetimes. Never fire-and-forget goroutines without a clear shutdown mechanism.

### 6.2 Channel Size

- Unbuffered channels by default.
- Buffer only when you have measured or clearly understand the buffering need.

### 6.3 Context

- `context.Context` is the first parameter of any function that makes RPCs, database calls, or does I/O.
- Never store Context in a struct. Pass it explicitly to each method.
- Use `context.Background()` only at the top level (main, init, tests).

### 6.4 Synchronous Functions

Prefer synchronous functions over asynchronous ones. Let callers add concurrency via goroutines if needed.

### 6.5 Zero-Value Mutex

`sync.Mutex` and `sync.RWMutex` are valid with zero values — no need to initialize with a constructor.

### 6.6 Avoid Mutable Globals

Use explicit dependency injection instead of package-level mutable state.

---

## 7. Testing

### 7.1 Test Framework

- Standard `testing` package + `github.com/stretchr/testify/require`.
- Mock framework: `go.uber.org/mock` (driven by `//go:generate`).

### 7.2 Table-Driven Tests

```go
func TestSomething(t *testing.T) {
    t.Parallel()
    tests := []struct {
        name          string
        input         string
        expected      string
        expectedError string
    }{
        {name: "valid input", input: "foo", expected: "bar"},
        {name: "empty input", input: "", expectedError: "input required"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            got, err := SomeFunc(tt.input)
            if tt.expectedError != "" {
                require.ErrorContains(t, err, tt.expectedError)
                return
            }
            require.NoError(t, err)
            require.Equal(t, tt.expected, got)
        })
    }
}
```

### 7.3 Test Helpers

- Use `t.Helper()` in helper functions.
- Produce useful failure messages: `t.Errorf("Foo(%q) = %d; want %d", in, got, want)`.

### 7.4 Coverage Targets

- Aim for **90% or higher package-level test coverage** for core packages whenever practical.
- Core packages include `auth`, `cache`, `config`, `env`, `httpclient`, `id`, `index`, `reader`, `server/mcp`, `service`, `spec`, `types`, and `workspace`.
- Treat coverage as an informational guide, not a hard gate; some error paths and generated code may be impractical to cover.
- Use `go test -coverprofile=coverage.out ./...` and inspect per-package coverage to identify gaps.

### 7.5 Mock Generation

```go
//go:generate go run go.uber.org/mock/mockgen -source=internal/server/mcp/handler.go -destination=internal/server/mcp/mock_svc_test.go -package=mcp
```

Run `go generate ./...` after changing the `svc` interface.

### 7.6 Build Tags

- `ci` — enables debug assertions
- `integration` — integration tests

### 7.7 Test File Naming

- `*_test.go` alongside source files (white-box testing).
- `*_test.go` in `package foo_test` for external tests.
- `*_fuzz_test.go` for fuzz tests.
- `*_benchmark_test.go` for benchmarks.

### 7.8 Service Test Pattern

Use `newTestService()` and `seedTestData()` helpers:

```go
func TestEndpointByID(t *testing.T) {
    t.Parallel()
    s := newTestService(t)
    seedTestData(t, s, "test-domain")

    resp, err := s.EndpointByID(context.Background(), service.EndpointByIDRequest{
        ID: "test-endpoint-id-32chars...",
    })
    require.NoError(t, err)
    require.Equal(t, "GET", resp.Method)
}
```

### 7.8 MCP Handler Test Pattern

Use `go.uber.org/mock/gomock`:

```go
func TestHandleSpecList(t *testing.T) {
    t.Parallel()
    ctrl := gomock.NewController(t)
    mockSvc := NewMocksvc(ctrl)
    mockSvc.EXPECT().Specs(gomock.Any()).Return(service.SpecsResponse{...}, nil)

    h := &handler{service: mockSvc}
    result, _, err := h.handleSpecList(context.Background(), nil, service.SpecsRequest{})
    require.NoError(t, err)
    require.NotNil(t, result)
}
```

---

## 8. Logging

### 8.1 Library

Use `log/slog` (standard library) for structured logging.

### 8.2 Structured Logging

```go
slog.Debug("processing request", "endpoint", endpointID, "duration", dur)
```

### 8.3 Log Levels

Debug, Info, Warn, Error.

---

## 9. Configuration

### 9.1 YAML Config Hierarchy

Config cascades: **global → spec → collection** (each level overrides the previous).

```yaml
mock_enabled: false
http_client:
  timeout: 30s
  follow_redirects: true
  max_response_size: 2048
mcp:
  transport: stdio
  addr: ":8080"
specs:
  - domain: "petstore"
    base_url: "https://petstore.swagger.io/v2"
    auth:
      type: api-key
      config:
        key: "X-API-Key"
        value: "abc123"
        in: header
    collections:
      - title: "Pet Operations"
        location: "https://petstore.swagger.io/v2/swagger.json"
```

### 9.2 Config Struct Tags

```go
type Config struct {
    MockEnabled  bool          `yaml:"mock_enabled"`
    HTTPClient   HTTPClientConfig `yaml:"http_client"`
    MCP          MCPConfig     `yaml:"mcp"`
    Specs        []Spec        `yaml:"specs" validate:"dive"`
}
```

### 9.3 Auth Dispatch

Config's `Auth.UnmarshalYAML` reads `type` field and dispatches to the correct auth package client. Supported types: `none`, `basic`, `bearer`, `digest`, `oauth2-cc`, `oauth2-pwd`, `api-key`, `script`.

---

## 10. Code Generation

### 10.1 Tools

- `mockgen` — mock implementations for interfaces (`go.uber.org/mock/mockgen`)
- `go generate ./...` — run all generators

### 10.2 Generated File Naming

`mock_*_test.go` for mockgen output.

### 10.3 go:generate Directives

Place in `generate.go` at project root, before the `package` declaration.

---

## 11. Project-Specific Conventions

### 11.1 Service Layer Pattern

Every service method follows this pattern:
1. Validate request with `s.validateRequest(req)` (uses `go-playground/validator`)
2. Look up entities from `s.index` (returns `LLMError` with `not_found` code)
3. Perform business logic
4. Return typed response or `LLMError`

```go
func (s *Service) Search(ctx context.Context, req SearchRequest) (SearchResponse, error) {
    if err := s.validateRequest(req); err != nil {
        return SearchResponse{}, NewLLMError(validationFailedErrCode, err.Error())
    }
    results, err := s.index.Search(req.Query, req.Limit)
    if err != nil {
        return SearchResponse{}, NewLLMError(invokeErrorCode, err.Error())
    }
    return SearchResponse{Results: results}, nil
}
```

### 11.2 Request/Response Pattern

Each method has a dedicated `{Method}Request` and `{Method}Response` struct. Request structs use `validate` tags for validation and `jsonschema` tags for documentation.

```go
type SearchRequest struct {
    Query string `json:"query" validate:"required,min=1" jsonschema:"description=Search query supporting field filters"`
    Limit int    `json:"limit" validate:"required,min=1,max=50" jsonschema:"description=Maximum results"`
}

type SearchResponse struct {
    Results []EndpointSearchItem `json:"results"`
}
```

### 11.3 Functional Options Pattern

```go
type NewOption func(*Service)

func New(opts ...NewOption) (*Service, error)

func WithDisableLLMAuth(disable bool) NewOption {
    return func(s *Service) {
        s.disableLLMAuth.Store(disable)
    }
}

func WithVersion(version string) NewOption {
    return func(s *Service) {
        s.version = version
    }
}
```

### 11.4 MCP Handler Pattern

```go
type svc interface {
    Specs(context.Context) (service.SpecsResponse, error)
    Search(context.Context, service.SearchRequest) (service.SearchResponse, error)
    Invoke(context.Context, service.InvokeRequest) (service.InvokeResponse, error)
    // ...
}

type handler struct {
    service svc
}

func (h *handler) handleSearch(ctx context.Context, _ *sdkmcp.CallToolRequest, req service.SearchRequest) (*sdkmcp.CallToolResult, any, error) {
    resp, err := h.service.Search(ctx, req)
    if err != nil {
        return nil, nil, err
    }
    return &sdkmcp.CallToolResult{StructuredContent: resp}, nil, nil
}
```

### 11.5 Rate Limiting

Per-endpoint 10-second cooldown for `invoke`:

```go
type invokeRateLimiter struct {
    mu    sync.Mutex
    calls map[string]time.Time
}

func (r *invokeRateLimiter) Allow(endpointID string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if last, ok := r.calls[endpointID]; ok && time.Since(last) < 10*time.Second {
        return NewLLMError(rateLimitErrCode,
            fmt.Sprintf("Endpoint %q was called %v ago. Wait %v before calling again.",
                endpointID, time.Since(last), 10*time.Second-time.Since(last)))
    }
    r.calls[endpointID] = time.Now()
    return nil
}
```

### 11.6 Response Size Management

Default max response size: 1 KB. When exceeded, save to `{workspace}/responses/` and return `FileReference`:

```go
type FileReference struct {
    Path        string `json:"path"`
    Size        int64  `json:"size"`
    SizeHint    string `json:"size_hint"`
    MaxSizeHint string `json:"max_size_hint"`
    Message     string `json:"message"`
    OpenCommand string `json:"open_command"`
}
```

### 11.7 ID Generation

MD5-based deterministic IDs:

```go
id.Domain("petstore")      // 32-char hex
id.Collection("petstore", "Pet Operations")
id.Tag("petstore", "Pet Operations", "pets")
id.Method("petstore", "Pet Operations", "pets", "GET", "/v2/pet/{petId}")
```

### 11.8 Auth Package

8 authentication methods in `internal/auth/`:

| Type | File | Config |
|------|------|--------|
| `none` | `noauth.go` | — |
| `basic` | `basic.go` | username, password |
| `bearer` | `bearer.go` | token |
| `digest` | `digest.go` | username, password |
| `oauth2-cc` | `oauth2_cc.go` | token_url, client_id, client_secret, scopes |
| `oauth2-pwd` | `oauth2_pwd.go` | token_url, client_id, client_secret, username, password, scopes |
| `api-key` | `api_key.go` | key, value, in (header/query/cookie) |
| `script` | `script.go` | path, args, env |

Each implements the `Authenticator` interface:

```go
type Authenticator interface {
    New() error
    Type() Type
    Apply(req *http.Request, out *Info) error
    Validate() error
}
```

### 11.9 Workspace Structure

`~/.swag2mcp` with subdirectories:
- `cache/` — downloaded spec files
- `specs/` — user spec files
- `responses/` — large invoke responses (cleaned after 48h)
- `auth_scripts/` — custom auth scripts

### 11.10 Linting

All code must pass linters in `.golangci.yml` (120-line limit, strict linters):
- `errcheck`, `errorlint`, `gocritic`, `gosec`, `govet`, `revive`, `staticcheck`, `testifylint`, `unparam`, `whitespace`, `bodyclose`, `depguard`, `perfsprint`, `prealloc`, `spancheck`, `usetesting`
- `depguard` denies: `gopkg.in/yaml.v2`/`v3` (use `go.yaml.in/yaml/v3`)

### 11.11 Go Version

This project uses **Go 1.23+**. Use Go 1.23+ features: `iter.Seq2`, generics, `t.Context()`, `slices`, `maps`, `cmp`.

### 11.12 Zip Slip Protection

Always use `filepath.Rel` to validate extracted paths, never `strings.HasPrefix`:

```go
// BAD — vulnerable to path traversal
if !strings.HasPrefix(filepath.Clean(fpath), filepath.Clean(destDir)+string(filepath.Separator)) {
    return fmt.Errorf("illegal file path: %s", f.Name)
}

// GOOD — safe
destDir := filepath.Clean(destDir)
fpath := filepath.Join(destDir, f.Name)
rel, err := filepath.Rel(destDir, fpath)
if err != nil || strings.HasPrefix(rel, "..") {
    return fmt.Errorf("zip slip detected: %s", f.Name)
}
```

`filepath.Rel` returns a relative path; if it starts with `..` the file would escape the destination directory. This is the standard Go idiom for zip slip prevention.

### 11.13 Coverage Strategy

- **Core packages** (`auth`, `cache`, `config`, `env`, `httpclient`, `id`, `index`, `server/mcp`, `service`, `spec`, `types`, `workspace`) — target 80%+
- **Integration packages** (`commands`, `tui`, `mockserver`) — informational only
- Commands: `make cover` (all), `make cover-core` (core only)

---

## 12. Quick Reference

| What | Convention |
|------|-----------|
| Package name | `auth`, `cache`, `config`, `index` |
| File name | `oauth2_cc.go`, `parse_v3.go` |
| Exported type | `LLMError`, `InvokeRequest` |
| Unexported type | `invokeRateLimiter` |
| Interface (1 method) | `Authenticator`, `TokenURLSetter` |
| Interface (multi) | `svc` |
| Error type | `LLMError` |
| Error value | `ErrNotFound` |
| Constructor | `NewService`, `NewLLMError` |
| Functional option | `WithDisableLLMAuth`, `WithVersion` |
| Receiver | `s *Service`, `h *handler` |
| Test file | `foo_test.go` |
| Generated file | `mock_svc_test.go` |
| Context param | First parameter |
| Error handling | `if err != nil { return err }` |
| LLM error | `NewLLMError(code, "what to do next")` |
| Slice declaration | `var s []T` (nil slice) |
| Import groups | std → third-party → local module |
| Config cascade | global → spec → collection |
| Auth dispatch | `UnmarshalYAML` reads `type` field |
| ID generation | `id.Domain()`, `id.Collection()`, `id.Tag()`, `id.Method()` |
| Rate limit | 10s per endpoint for `invoke` |
| Response size | 1 KB default, saved to file when exceeded |
| Mock framework | `go.uber.org/mock` |
| Test helpers | `newTestService()`, `seedTestData()` |
| Lint command | `make lint` |
| Test command | `go test ./...` |
| Generate command | `go generate ./...` |

### 11.13 Variable Groups

Group consecutive `var` declarations with `var ( ... )`. Single declarations are fine.

```go
// Good
var (
    workspaceDir string
    configPath   string
)

// Good (single)
var version string

// Bad
var workspaceDir string
var configPath string
```
