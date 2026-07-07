package service

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/id"
	"github.com/mmadfox/swag2mcp/internal/index"
	"github.com/mmadfox/swag2mcp/internal/spec"
	"github.com/mmadfox/swag2mcp/internal/types"
)

//go:embed testdata/invoke/*.yaml
var testDataFS embed.FS

// TestInvoke_GetRequest verifies that a simple GET request works correctly.
func TestInvoke_GetRequest(t *testing.T) {
	t.Parallel()

	specDoc := parseSpecFromFile(t, "users.yaml")
	testServer := newTestServer(t, testServerConfig{
		ExpectedMethod: http.MethodGet,
		ExpectedPath:   "/users",
		StatusCode:     http.StatusOK,
		ResponseBody:   map[string]any{"users": []any{}},
	})
	t.Cleanup(testServer.Close)

	serviceInstance := buildTestService(t, t.Name(), specDoc, nil, nil)

	// Override the base URL to point to our test server
	specInfo, _ := serviceInstance.index.SpecByID(specIDForTest(t))
	specInfo.BaseURL = testServer.URL

	endpointID := findEndpointID(t, serviceInstance, http.MethodGet, "/users")

	response, invokeError := serviceInstance.Invoke(context.Background(), InvokeRequest{
		EndpointID: endpointID,
	})
	if invokeError != nil {
		t.Fatalf("Invoke() returned error: %v", invokeError)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", response.StatusCode, http.StatusOK)
	}
}

// TestInvoke_GetRequestWithQuery verifies that query parameters are sent correctly.
func TestInvoke_GetRequestWithQuery(t *testing.T) {
	t.Parallel()

	specDoc := parseSpecFromFile(t, "users.yaml")
	testServer := newTestServer(t, testServerConfig{
		ExpectedMethod: http.MethodGet,
		ExpectedPath:   "/users",
		ExpectedQuery:  map[string]string{"limit": "10", "offset": "0"},
		StatusCode:     http.StatusOK,
		ResponseBody:   map[string]any{"users": []any{}},
	})
	t.Cleanup(testServer.Close)

	serviceInstance := buildTestService(t, t.Name(), specDoc, nil, nil)
	specInfo, _ := serviceInstance.index.SpecByID(specIDForTest(t))
	specInfo.BaseURL = testServer.URL

	endpointID := findEndpointID(t, serviceInstance, http.MethodGet, "/users")

	response, invokeError := serviceInstance.Invoke(context.Background(), InvokeRequest{
		EndpointID: endpointID,
		Parameters: map[string]any{"limit": "10", "offset": "0"},
	})
	if invokeError != nil {
		t.Fatalf("Invoke() returned error: %v", invokeError)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", response.StatusCode, http.StatusOK)
	}
}

// TestInvoke_PostRequestWithBody verifies that POST requests with a body work correctly.
func TestInvoke_PostRequestWithBody(t *testing.T) {
	t.Parallel()

	specDoc := parseSpecFromFile(t, "orders.yaml")
	testServer := newTestServer(t, testServerConfig{
		ExpectedMethod: http.MethodPost,
		ExpectedPath:   "/orders",
		ExpectedBody:   map[string]any{"productId": "prod-1", "quantity": float64(2)},
		StatusCode:     http.StatusCreated,
		ResponseBody:   map[string]any{"orderId": "ord-1"},
	})
	t.Cleanup(testServer.Close)

	serviceInstance := buildTestService(t, "test", specDoc, nil, nil)
	specInfo, _ := serviceInstance.index.SpecByID(specIDForTest(t))
	specInfo.BaseURL = testServer.URL

	endpointID := findEndpointID(t, serviceInstance, http.MethodPost, "/orders")

	response, invokeError := serviceInstance.Invoke(context.Background(), InvokeRequest{
		EndpointID: endpointID,
		RequestBody: map[string]any{
			"productId": "prod-1",
			"quantity":  2,
		},
	})
	if invokeError != nil {
		t.Fatalf("Invoke() returned error: %v", invokeError)
	}

	if response.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want %d", response.StatusCode, http.StatusCreated)
	}
}

// TestInvoke_DeleteRequest verifies that DELETE requests work correctly.
func TestInvoke_DeleteRequest(t *testing.T) {
	t.Parallel()

	specDoc := parseSpecFromFile(t, "users.yaml")
	testServer := newTestServer(t, testServerConfig{
		ExpectedMethod: http.MethodDelete,
		ExpectedPath:   "/users/user-123",
		StatusCode:     http.StatusNoContent,
	})
	t.Cleanup(testServer.Close)

	serviceInstance := buildTestService(t, "delete-request", specDoc, nil, nil)
	specInfo, _ := serviceInstance.index.SpecByID(specIDForTest(t))
	specInfo.BaseURL = testServer.URL

	endpointID := findEndpointID(t, serviceInstance, http.MethodDelete, "/users/{userId}")

	response, invokeError := serviceInstance.Invoke(context.Background(), InvokeRequest{
		EndpointID: endpointID,
		Parameters: map[string]any{"userId": "user-123"},
	})
	if invokeError != nil {
		t.Fatalf("Invoke() returned error: %v", invokeError)
	}

	if response.StatusCode != http.StatusNoContent {
		t.Errorf("StatusCode = %d, want %d", response.StatusCode, http.StatusNoContent)
	}
}

// TestInvoke_PatchRequest verifies that PATCH requests with headers work correctly.
func TestInvoke_PatchRequest(t *testing.T) {
	t.Parallel()

	specDoc := parseSpecFromFile(t, "orders.yaml")
	testServer := newTestServer(t, testServerConfig{
		ExpectedMethod: http.MethodPatch,
		ExpectedPath:   "/orders/ord-1",
		ExpectedHeaders: map[string]string{
			"X-Idempotency-Key": "idem-abc-123",
		},
		ExpectedBody: map[string]any{"status": "shipped"},
		StatusCode:   http.StatusOK,
		ResponseBody: map[string]any{"orderId": "ord-1", "status": "shipped"},
	})
	t.Cleanup(testServer.Close)

	serviceInstance := buildTestService(t, "patch-request", specDoc, nil, nil)
	specInfo, _ := serviceInstance.index.SpecByID(specIDForTest(t))
	specInfo.BaseURL = testServer.URL

	endpointID := findEndpointID(t, serviceInstance, http.MethodPatch, "/orders/{orderId}")

	response, invokeError := serviceInstance.Invoke(context.Background(), InvokeRequest{
		EndpointID: endpointID,
		Parameters: map[string]any{
			"orderId":           "ord-1",
			"X-Idempotency-Key": "idem-abc-123",
		},
		RequestBody: map[string]any{"status": "shipped"},
	})
	if invokeError != nil {
		t.Fatalf("Invoke() returned error: %v", invokeError)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", response.StatusCode, http.StatusOK)
	}
}

// TestInvoke_SpecHeaders verifies that spec-level headers are sent on every request.
func TestInvoke_SpecHeaders(t *testing.T) {
	t.Parallel()

	specDoc := parseSpecFromFile(t, "users.yaml")
	testServer := newTestServer(t, testServerConfig{
		ExpectedMethod: http.MethodGet,
		ExpectedPath:   "/users",
		ExpectedHeaders: map[string]string{
			"X-Source":      "swag2mcp",
			"X-Api-Version": "2024-01",
		},
		StatusCode:   http.StatusOK,
		ResponseBody: map[string]any{"users": []any{}},
	})
	t.Cleanup(testServer.Close)

	serviceInstance := buildTestService(t, "spec-headers", specDoc,
		map[string]string{
			"X-Source":      "swag2mcp",
			"X-Api-Version": "2024-01",
		},
		nil,
	)
	specInfo, _ := serviceInstance.index.SpecByID(specIDForTest(t))
	specInfo.BaseURL = testServer.URL

	endpointID := findEndpointID(t, serviceInstance, http.MethodGet, "/users")

	response, invokeError := serviceInstance.Invoke(context.Background(), InvokeRequest{
		EndpointID: endpointID,
	})
	if invokeError != nil {
		t.Fatalf("Invoke() returned error: %v", invokeError)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", response.StatusCode, http.StatusOK)
	}
}

// TestInvoke_CollectionHeaders verifies that collection-level headers override spec-level headers.
func TestInvoke_CollectionHeaders(t *testing.T) {
	t.Parallel()

	specDoc := parseSpecFromFile(t, "users.yaml")
	testServer := newTestServer(t, testServerConfig{
		ExpectedMethod: http.MethodGet,
		ExpectedPath:   "/users",
		ExpectedHeaders: map[string]string{
			"X-Region": "us-east-1",
		},
		StatusCode:   http.StatusOK,
		ResponseBody: map[string]any{"users": []any{}},
	})
	t.Cleanup(testServer.Close)

	serviceInstance := buildTestService(t, "coll-headers", specDoc,
		map[string]string{
			"X-Region": "eu-west-1",
		},
		nil,
	)
	specInfo, _ := serviceInstance.index.SpecByID(specIDForTest(t))
	specInfo.BaseURL = testServer.URL

	// Override collection headers
	collectionID := id.Collection(specInfo.ID, t.Name()+"/collection")
	collection, _ := serviceInstance.index.CollectionByID(collectionID)
	collection.Headers = map[string]string{
		"X-Region": "us-east-1",
	}

	endpointID := findEndpointID(t, serviceInstance, http.MethodGet, "/users")

	response, invokeError := serviceInstance.Invoke(context.Background(), InvokeRequest{
		EndpointID: endpointID,
	})
	if invokeError != nil {
		t.Fatalf("Invoke() returned error: %v", invokeError)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", response.StatusCode, http.StatusOK)
	}
}

// TestInvoke_BearerAuth verifies that Bearer token authentication works.
func TestInvoke_BearerAuth(t *testing.T) {
	t.Parallel()

	specDoc := parseSpecFromFile(t, "users.yaml")
	testServer := newTestServer(t, testServerConfig{
		ExpectedMethod: http.MethodGet,
		ExpectedPath:   "/users",
		AuthType:       "bearer",
		AuthCredentials: map[string]string{
			"token": "my-bearer-token",
		},
		StatusCode:   http.StatusOK,
		ResponseBody: map[string]any{"users": []any{}},
	})
	t.Cleanup(testServer.Close)

	authenticator := &auth.BearerTokenAuthClient{Token: "my-bearer-token"}
	if newError := authenticator.New(); newError != nil {
		t.Fatalf("failed to init authenticator: %v", newError)
	}

	serviceInstance := buildTestService(t, "bearer-auth", specDoc, nil, authenticator)
	specInfo, _ := serviceInstance.index.SpecByID(specIDForTest(t))
	specInfo.BaseURL = testServer.URL

	endpointID := findEndpointID(t, serviceInstance, http.MethodGet, "/users")

	response, invokeError := serviceInstance.Invoke(context.Background(), InvokeRequest{
		EndpointID: endpointID,
	})
	if invokeError != nil {
		t.Fatalf("Invoke() returned error: %v", invokeError)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", response.StatusCode, http.StatusOK)
	}
}

// TestInvoke_BasicAuth verifies that HTTP Basic authentication works.
func TestInvoke_BasicAuth(t *testing.T) {
	t.Parallel()

	specDoc := parseSpecFromFile(t, "users.yaml")
	testServer := newTestServer(t, testServerConfig{
		ExpectedMethod: http.MethodGet,
		ExpectedPath:   "/users",
		AuthType:       "basic",
		AuthCredentials: map[string]string{
			"username": "admin",
			"password": "secret",
		},
		StatusCode:   http.StatusOK,
		ResponseBody: map[string]any{"users": []any{}},
	})
	t.Cleanup(testServer.Close)

	authenticator := &auth.BasicAuthClient{Username: "admin", Password: "secret"}
	if newError := authenticator.New(); newError != nil {
		t.Fatalf("failed to init authenticator: %v", newError)
	}

	serviceInstance := buildTestService(t, "basic-auth", specDoc, nil, authenticator)
	specInfo, _ := serviceInstance.index.SpecByID(specIDForTest(t))
	specInfo.BaseURL = testServer.URL

	endpointID := findEndpointID(t, serviceInstance, http.MethodGet, "/users")

	response, invokeError := serviceInstance.Invoke(context.Background(), InvokeRequest{
		EndpointID: endpointID,
	})
	if invokeError != nil {
		t.Fatalf("Invoke() returned error: %v", invokeError)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", response.StatusCode, http.StatusOK)
	}
}

// TestInvoke_OAuth2ClientCredentials verifies that OAuth2 Client Credentials auth works.
func TestInvoke_OAuth2ClientCredentials(t *testing.T) {
	t.Parallel()

	authServer := newTestAuthServer(t, "oauth2-token-123")
	t.Cleanup(authServer.Close)

	specDoc := parseSpecFromFile(t, "users.yaml")
	testServer := newTestServer(t, testServerConfig{
		ExpectedMethod: http.MethodGet,
		ExpectedPath:   "/users",
		AuthType:       "oauth2-cc",
		AuthCredentials: map[string]string{
			"token": "oauth2-token-123",
		},
		StatusCode:   http.StatusOK,
		ResponseBody: map[string]any{"users": []any{}},
	})
	t.Cleanup(testServer.Close)

	authenticator := &auth.OAuth2ClientCredentialsAuthClient{
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		TokenURL:     authServer.URL + "/token",
	}
	if newError := authenticator.New(); newError != nil {
		t.Fatalf("failed to init authenticator: %v", newError)
	}

	serviceInstance := buildTestService(t, "oauth2-cc", specDoc, nil, authenticator)
	specInfo, _ := serviceInstance.index.SpecByID(specIDForTest(t))
	specInfo.BaseURL = testServer.URL

	endpointID := findEndpointID(t, serviceInstance, http.MethodGet, "/users")

	response, invokeError := serviceInstance.Invoke(context.Background(), InvokeRequest{
		EndpointID: endpointID,
	})
	if invokeError != nil {
		t.Fatalf("Invoke() returned error: %v", invokeError)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", response.StatusCode, http.StatusOK)
	}
}

// TestInvoke_DigestAuth verifies that HTTP Digest authentication works.
func TestInvoke_DigestAuth(t *testing.T) {
	t.Parallel()

	specDoc := parseSpecFromFile(t, "users.yaml")
	digestServer := newTestDigestServer(t, "digest-user", "digest-pass")
	t.Cleanup(digestServer.Close)

	serviceInstance := buildTestService(t, "digest-auth", specDoc, nil, nil)
	specInfo, _ := serviceInstance.index.SpecByID(specIDForTest(t))
	specInfo.BaseURL = digestServer.URL

	// Digest auth is applied via the auth transport, so we need to set it up
	authenticator := &auth.DigestAuthClient{
		Username: "digest-user",
		Password: "digest-pass",
	}
	if newError := authenticator.New(); newError != nil {
		t.Fatalf("failed to init authenticator: %v", newError)
	}
	specInfo.Auth = authenticator

	endpointID := findEndpointID(t, serviceInstance, http.MethodGet, "/users")

	response, invokeError := serviceInstance.Invoke(context.Background(), InvokeRequest{
		EndpointID: endpointID,
	})
	if invokeError != nil {
		t.Fatalf("Invoke() returned error: %v", invokeError)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", response.StatusCode, http.StatusOK)
	}
}

// TestInvoke_NotFoundError verifies that invoking a non-existent endpoint returns an error.
func TestInvoke_NotFoundError(t *testing.T) {
	t.Parallel()

	serviceInstance, serviceError := New()
	if serviceError != nil {
		t.Fatalf("failed to create service: %v", serviceError)
	}

	_, invokeError := serviceInstance.Invoke(context.Background(), InvokeRequest{
		EndpointID: "00000000000000000000000000000000",
	})
	if invokeError == nil {
		t.Fatal("expected error for non-existent endpoint, got nil")
	}
}

// TestInvoke_ValidationError verifies that an invalid endpoint ID returns a validation error.
func TestInvoke_ValidationError(t *testing.T) {
	t.Parallel()

	serviceInstance, serviceError := New()
	if serviceError != nil {
		t.Fatalf("failed to create service: %v", serviceError)
	}

	_, invokeError := serviceInstance.Invoke(context.Background(), InvokeRequest{
		EndpointID: "invalid-id",
	})
	if invokeError == nil {
		t.Fatal("expected validation error for invalid endpoint ID, got nil")
	}
}

// TestInvoke_ContentTypeHeader verifies that Content-Type is set for requests with a body.
func TestInvoke_ContentTypeHeader(t *testing.T) {
	t.Parallel()

	specDoc := parseSpecFromFile(t, "orders.yaml")
	testServer := newTestServer(t, testServerConfig{
		ExpectedMethod: http.MethodPost,
		ExpectedPath:   "/orders",
		ExpectedHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode:   http.StatusCreated,
		ResponseBody: map[string]any{"orderId": "ord-1"},
	})
	t.Cleanup(testServer.Close)

	serviceInstance := buildTestService(t, "content-type", specDoc, nil, nil)
	specInfo, _ := serviceInstance.index.SpecByID(specIDForTest(t))
	specInfo.BaseURL = testServer.URL

	endpointID := findEndpointID(t, serviceInstance, http.MethodPost, "/orders")

	response, invokeError := serviceInstance.Invoke(context.Background(), InvokeRequest{
		EndpointID: endpointID,
		RequestBody: map[string]any{
			"productId": "prod-1",
			"quantity":  2,
		},
	})
	if invokeError != nil {
		t.Fatalf("Invoke() returned error: %v", invokeError)
	}

	if response.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want %d", response.StatusCode, http.StatusCreated)
	}
}

// TestInvoke_TableDriven runs multiple test cases in a table-driven style.
func TestInvoke_TableDriven(t *testing.T) {
	t.Parallel()

	type invokeTestCase struct {
		Name           string
		SpecFile       string
		Method         string
		Path           string
		Parameters     map[string]any
		RequestBody    map[string]any
		SpecHeaders    map[string]string
		AuthType       string
		AuthToken      string
		ExpectedStatus int
	}

	testCases := []invokeTestCase{
		{
			Name:           "GET list users without parameters",
			SpecFile:       "users.yaml",
			Method:         http.MethodGet,
			Path:           "/users",
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:     "GET list users with query parameters",
			SpecFile: "users.yaml",
			Method:   http.MethodGet,
			Path:     "/users",
			Parameters: map[string]any{
				"limit":  "20",
				"offset": "5",
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:     "GET user by ID with path parameter",
			SpecFile: "users.yaml",
			Method:   http.MethodGet,
			Path:     "/users/{userId}",
			Parameters: map[string]any{
				"userId": "user-42",
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:     "DELETE user by ID",
			SpecFile: "users.yaml",
			Method:   http.MethodDelete,
			Path:     "/users/{userId}",
			Parameters: map[string]any{
				"userId": "user-to-delete",
			},
			ExpectedStatus: http.StatusNoContent,
		},
		{
			Name:     "POST create order with body",
			SpecFile: "orders.yaml",
			Method:   http.MethodPost,
			Path:     "/orders",
			RequestBody: map[string]any{
				"productId": "prod-99",
				"quantity":  1,
			},
			ExpectedStatus: http.StatusCreated,
		},
		{
			Name:     "PATCH update order with headers and body",
			SpecFile: "orders.yaml",
			Method:   http.MethodPatch,
			Path:     "/orders/{orderId}",
			Parameters: map[string]any{
				"orderId":           "ord-42",
				"X-Idempotency-Key": "idem-xyz-789",
			},
			RequestBody: map[string]any{
				"status": "cancelled",
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:     "GET with spec-level custom headers",
			SpecFile: "users.yaml",
			Method:   http.MethodGet,
			Path:     "/users",
			SpecHeaders: map[string]string{
				"X-Custom-Header": "custom-value",
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:           "GET with Bearer authentication",
			SpecFile:       "users.yaml",
			Method:         http.MethodGet,
			Path:           "/users",
			AuthType:       "bearer",
			AuthToken:      "table-driven-token",
			ExpectedStatus: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			specDoc := parseSpecFromFile(t, testCase.SpecFile)

			serverConfig := testServerConfig{
				ExpectedMethod: testCase.Method,
				StatusCode:     testCase.ExpectedStatus,
				ResponseBody:   map[string]any{"result": "ok"},
			}

			// Only set ExpectedPath for endpoints without path parameters
			if !strings.Contains(testCase.Path, "{") {
				serverConfig.ExpectedPath = testCase.Path
			}

			if testCase.AuthType != "" {
				serverConfig.AuthType = testCase.AuthType
				serverConfig.AuthCredentials = map[string]string{"token": testCase.AuthToken}
			}

			testServer := newTestServer(t, serverConfig)
			t.Cleanup(testServer.Close)

			var authenticator auth.Authenticator
			if testCase.AuthType == "bearer" {
				authenticator = &auth.BearerTokenAuthClient{Token: testCase.AuthToken}
				if newError := authenticator.New(); newError != nil {
					t.Fatalf("failed to init authenticator: %v", newError)
				}
			}

			serviceInstance := buildTestService(t, testCase.Name, specDoc, testCase.SpecHeaders, authenticator)
			specInfo, _ := serviceInstance.index.SpecByID(specIDForTest(t))
			specInfo.BaseURL = testServer.URL

			endpointID := findEndpointID(t, serviceInstance, testCase.Method, testCase.Path)

			response, invokeError := serviceInstance.Invoke(context.Background(), InvokeRequest{
				EndpointID:  endpointID,
				Parameters:  testCase.Parameters,
				RequestBody: testCase.RequestBody,
			})
			if invokeError != nil {
				t.Fatalf("Invoke() returned error: %v", invokeError)
			}

			if response.StatusCode != testCase.ExpectedStatus {
				t.Errorf("StatusCode = %d, want %d", response.StatusCode, testCase.ExpectedStatus)
			}
		})
	}
}

// testServerConfig holds configuration for a test HTTP server.
type testServerConfig struct {
	ExpectedMethod  string
	ExpectedPath    string
	ExpectedHeaders map[string]string
	ExpectedQuery   map[string]string
	ExpectedBody    map[string]any
	StatusCode      int
	ResponseBody    any
	ResponseHeaders map[string]string
	AuthType        string // "bearer", "basic", "oauth2-cc", "digest", ""
	AuthCredentials map[string]string
}

// newTestServer creates an [httptest.Server] that validates incoming requests against the config.
//
//nolint:gocognit
func newTestServer(t *testing.T, config testServerConfig) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if config.ExpectedMethod != "" && request.Method != config.ExpectedMethod {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = writer.Write([]byte(`{"error":"method not allowed"}`))
			return
		}

		if config.ExpectedPath != "" && !strings.HasSuffix(request.URL.Path, config.ExpectedPath) {
			writer.WriteHeader(http.StatusNotFound)
			_, _ = fmt.Fprintf(writer, `{"error":"path not found"}`)
			return
		}

		for headerName, headerValue := range config.ExpectedHeaders {
			actualValue := request.Header.Get(headerName)
			if actualValue != headerValue {
				writer.WriteHeader(http.StatusBadRequest)
				_, _ = fmt.Fprintf(writer,
					`{"error":"expected header %s: %s, got: %s"}`, headerName, headerValue, actualValue,
				)
				return
			}
		}

		for queryName, queryValue := range config.ExpectedQuery {
			actualValue := request.URL.Query().Get(queryName)
			if actualValue != queryValue {
				writer.WriteHeader(http.StatusBadRequest)
				_, _ = fmt.Fprintf(writer,
					`{"error":"expected query %s: %s, got: %s"}`, queryName, queryValue, actualValue,
				)
				return
			}
		}

		if config.ExpectedBody != nil {
			var receivedBody map[string]any
			if decodeError := json.NewDecoder(request.Body).Decode(&receivedBody); decodeError != nil {
				writer.WriteHeader(http.StatusBadRequest)
				_, _ = fmt.Fprintf(writer, `{"error":"invalid body: %s"}`, decodeError)
				return
			}
			for key, expectedValue := range config.ExpectedBody {
				actualValue, exists := receivedBody[key]
				if !exists {
					writer.WriteHeader(http.StatusBadRequest)
					_, _ = fmt.Fprintf(writer, `{"error":"missing body field: %s"}`, key)
					return
				}
				if fmt.Sprintf("%v", actualValue) != fmt.Sprintf("%v", expectedValue) {
					writer.WriteHeader(http.StatusBadRequest)
					_, _ = fmt.Fprintf(writer,
						`{"error":"body field %s: expected %v, got %v"}`, key, expectedValue, actualValue,
					)
					return
				}
			}
		}

		if config.AuthType != "" {
			authError := validateAuthHeader(t, request, config)
			if authError != "" {
				writer.WriteHeader(http.StatusUnauthorized)
				_, _ = fmt.Fprintf(writer, `{"error":"%s"}`, authError)
				return
			}
		}

		for headerName, headerValue := range config.ResponseHeaders {
			writer.Header().Set(headerName, headerValue)
		}

		writer.WriteHeader(config.StatusCode)
		if config.ResponseBody != nil {
			_ = json.NewEncoder(writer).Encode(config.ResponseBody)
		}
	}))
}

// validateAuthHeader checks the Authorization header against the expected auth type.
func validateAuthHeader(t *testing.T, request *http.Request, config testServerConfig) string {
	t.Helper()

	switch config.AuthType {
	case "bearer":
		expectedToken := config.AuthCredentials["token"]
		authHeader := request.Header.Get("Authorization")
		if authHeader != "Bearer "+expectedToken {
			return fmt.Sprintf("expected Bearer %s, got %s", expectedToken, authHeader)
		}
	case "basic":
		username, password, ok := request.BasicAuth()
		if !ok {
			return "missing basic auth"
		}
		expectedUser := config.AuthCredentials["username"]
		expectedPass := config.AuthCredentials["password"]
		if username != expectedUser || password != expectedPass {
			return fmt.Sprintf("expected %s:%s, got %s:%s", expectedUser, expectedPass, username, password)
		}
	case "oauth2-cc":
		expectedToken := config.AuthCredentials["token"]
		authHeader := request.Header.Get("Authorization")
		if authHeader != "Bearer "+expectedToken {
			return fmt.Sprintf("expected Bearer %s, got %s", expectedToken, authHeader)
		}
	case "digest":
		authHeader := request.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Digest ") {
			return fmt.Sprintf("expected Digest auth, got %s", authHeader)
		}
	}

	return ""
}

// newTestAuthServer creates an [httptest.Server] that simulates an OAuth2 token endpoint.
func newTestAuthServer(t *testing.T, token string) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		_ = request.ParseForm()
		grantType := request.Form.Get("grant_type")

		response := map[string]any{
			"access_token": token,
			"token_type":   "Bearer",
			"expires_in":   3600,
		}

		if grantType == "client_credentials" || grantType == "password" {
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		writer.WriteHeader(http.StatusBadRequest)
	}))
}

// newTestDigestServer creates an [httptest.Server] that challenges with Digest auth.
func newTestDigestServer(t *testing.T, _, _ string) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")
		if authHeader == "" {
			writer.Header().Set("WWW-Authenticate",
				`Digest realm="test-realm", nonce="test-nonce", opaque="test-opaque", qop="auth", algorithm="MD5"`,
			)
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Digest ") {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(`{"ok":true}`))
	}))
}

// parseSpecFromFile parses a YAML spec file and returns the parsed Doc.
func parseSpecFromFile(t *testing.T, filePath string) *spec.Doc {
	t.Helper()

	data := readTestFile(t, filePath)
	doc, parseError := spec.Parse(data)
	if parseError != nil {
		t.Fatalf("failed to parse spec %s: %v", filePath, parseError)
	}
	return doc
}

// readTestFile reads a test data file from the testdata directory.
func readTestFile(t *testing.T, relativePath string) []byte {
	t.Helper()

	data, readError := testDataFS.ReadFile("testdata/invoke/" + relativePath)
	if readError != nil {
		t.Fatalf("failed to read test file %s: %v", relativePath, readError)
	}
	return data
}

// buildTestService creates a Service with a pre-populated index for testing.
// Uses t.Name() as the unique domain to guarantee isolation across parallel tests.
func buildTestService(t *testing.T, _ string, specDoc *spec.Doc, specHeaders map[string]string, authenticator auth.Authenticator) *Service {
	t.Helper()

	uniqueDomain := t.Name()

	serviceInstance, serviceError := New()
	if serviceError != nil {
		t.Fatalf("failed to create service: %v", serviceError)
	}

	newIndex, newIndexError := index.New()
	if newIndexError != nil {
		t.Fatalf("failed to create index: %v", newIndexError)
	}
	serviceInstance.index = newIndex

	specID := id.Domain(uniqueDomain)
	specInfo := &types.Spec{
		ID:      specID,
		Domain:  uniqueDomain,
		BaseURL: "http://test-server",
		Headers: specHeaders,
		Auth:    authenticator,
	}

	collectionID := id.Collection(specID, uniqueDomain+"/collection")
	collectionInfo := &types.Collection{
		ID:     collectionID,
		SpecID: specID,
	}

	var allTags []*types.Tag
	var allEndpoints []*types.Endpoint

	for index, pathItem := range specDoc.PathItems {
		operation := pathItem.Operation
		if operation == nil {
			continue
		}

		tagName := fmt.Sprintf("%s-tag-%d", uniqueDomain, index)
		tagID := id.Tag(specID, collectionID, tagName)
		tagInfo := &types.Tag{
			ID:           tagID,
			SpecID:       specID,
			CollectionID: collectionID,
			Name:         tagName,
		}
		allTags = append(allTags, tagInfo)

		endpoint := &types.Endpoint{
			ID: id.Method(
				specID,
				collectionID,
				tagID,
				pathItem.Method,
				pathItem.Path,
				operation.ID,
			),
			SpecID:       specID,
			CollectionID: collectionID,
			TagID:        tagID,
			Tag:          tagName,
			Name:         pathItem.Method,
			Path:         pathItem.Path,
			Operation:    operation,
		}
		allEndpoints = append(allEndpoints, endpoint)
	}

	if ensureError := serviceInstance.index.EnsureIndex(specInfo, []*types.Collection{collectionInfo}, allTags, allEndpoints); ensureError != nil {
		t.Fatalf("failed to index: %v", ensureError)
	}

	return serviceInstance
}

// findEndpointID finds the first endpoint ID matching the given method and path pattern.
func findEndpointID(t *testing.T, serviceInstance *Service, method string, pathSuffix string) string {
	t.Helper()

	for cursor := range serviceInstance.index.IterateByEndpoints() {
		if cursor.Endpoint.Name == method && strings.HasSuffix(cursor.Endpoint.Path, pathSuffix) {
			return cursor.Endpoint.ID
		}
	}
	t.Fatalf("endpoint not found: %s %s", method, pathSuffix)
	return ""
}

// specIDForTest returns the domain ID for the current test.
func specIDForTest(t *testing.T) string {
	t.Helper()
	return id.Domain(t.Name())
}
