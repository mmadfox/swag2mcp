package mockserver

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mmadfox/swag2mcp/internal/cache"
	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/spec"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

const (
	apiShutdownTimeout = 5
	contentTypeJSON    = "application/json"
	responseStatus200  = "200"
	responseStatus201  = "201"
	responseStatusDef  = "default"
)

// apiMockServer serves mock responses for a single OpenAPI/Swagger collection.
// It parses the spec document and registers HTTP handlers for each endpoint.
type apiMockServer struct {
	specDomain      string
	collectionTitle string
	addr            string
	server          *http.Server
	logger          *slog.Logger
	tlsConfig       *tls.Config
	doc             *spec.Doc
	workspace       *workspace.Workspace
}

// newAPIMockServer creates a new API mock server for the given spec and collection.
// It parses the OpenAPI document at the collection's location.
func newAPIMockServer(
	specConfig *config.Spec,
	collectionConfig *config.Collection,
	addr string,
	tlsConfig *tls.Config,
	logger *slog.Logger,
	ws *workspace.Workspace,
) *apiMockServer {
	specDocument, parseError := parseSpecDocument(collectionConfig.Location, ws)
	if parseError != nil {
		logger.Error("failed to parse spec, skipping collection",
			"spec", specConfig.Domain,
			"collection", collectionConfig.Location,
			"error", parseError,
		)
		return nil
	}

	return &apiMockServer{
		specDomain:      specConfig.Domain,
		collectionTitle: collectionConfig.LLMTitle,
		addr:            addr,
		tlsConfig:       tlsConfig,
		logger:          logger,
		doc:             specDocument,
		workspace:       ws,
	}
}

// parseSpecDocument reads and parses an OpenAPI/Swagger spec from a local path or URL.
// It uses the cache to resolve remote URLs to local files.
func parseSpecDocument(location string, ws *workspace.Workspace) (*spec.Doc, error) {
	localPath := location

	workspaceDir := ""
	if ws != nil {
		workspaceDir = ws.Root()
	}

	cacheDir := cache.New(workspaceDir)
	if resolvedPath, resolveError := cacheDir.Resolve(context.Background(), location); resolveError == nil {
		localPath = resolvedPath
	}

	data, readError := os.ReadFile(localPath)
	if readError != nil {
		return nil, fmt.Errorf("read spec file: %w", readError)
	}

	specDocument, parseError := spec.Parse(data)
	if parseError != nil {
		return nil, fmt.Errorf("parse spec: %w", parseError)
	}

	return specDocument, nil
}

// start begins listening for HTTP requests on the configured port.
// It registers a handler for each path+method combination in the spec document.
func (m *apiMockServer) start(ctx context.Context) {
	mux := http.NewServeMux()

	for _, pathItem := range m.doc.PathItems {
		if pathItem == nil || pathItem.Operation == nil {
			continue
		}

		path := pathItem.Path
		method := strings.ToUpper(pathItem.Method)

		operation := pathItem.Operation
		handler := m.createEndpointHandler(operation)

		pattern := method + " " + path
		mux.HandleFunc(pattern, handler)
	}

	address := extractHostPort(m.addr)
	if !strings.Contains(address, ":") {
		address = ":" + address
	}
	m.server = &http.Server{
		Addr:              address,
		Handler:           mux,
		ReadHeaderTimeout: apiShutdownTimeout * time.Second,
	}

	if m.tlsConfig != nil {
		m.server.TLSConfig = m.tlsConfig
	}

	go func() {
		serveError := m.server.ListenAndServe()
		if serveError != nil && serveError != http.ErrServerClosed {
			m.logger.ErrorContext(ctx, "API mock server error",
				"spec", m.specDomain,
				"addr", m.addr,
				"error", serveError,
			)
		}
	}()

	m.logger.InfoContext(ctx, "API mock server started",
		"spec", m.specDomain,
		"addr", m.addr,
		"endpoints", len(m.doc.PathItems),
	)
}

// shutdown gracefully stops the HTTP server.
func (m *apiMockServer) shutdown() {
	if m.server != nil {
		shutdownContext, shutdownCancel := context.WithTimeout(
			context.Background(),
			apiShutdownTimeout*time.Second,
		)
		defer shutdownCancel()
		if err := m.server.Shutdown(shutdownContext); err != nil {
			m.logger.Warn("mock api server shutdown error", "error", err)
		}
	}
}

// createEndpointHandler returns an HTTP handler that generates a mock response
// for the given operation. It finds the response schema and generates random
// data that conforms to it.
func (m *apiMockServer) createEndpointHandler(operation *spec.Operation) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, _ *http.Request) {
		responseSchema := m.findResponseSchema(operation)
		if responseSchema == nil {
			responseWriter.Header().Set("Content-Type", "application/json")
			responseWriter.WriteHeader(http.StatusOK)
			_, _ = responseWriter.Write([]byte(`{}`))
			return
		}

		generatedData := GenerateFromSchema(responseSchema)

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)

		encodeError := json.NewEncoder(responseWriter).Encode(generatedData)
		if encodeError != nil {
			m.logger.ErrorContext(context.Background(), "failed to encode response",
				"error", encodeError,
			)
			http.Error(responseWriter, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// findResponseSchema finds the first response schema from the operation's
// responses map. It prefers 200, 201, and default status codes, then falls
// back to any available response.
func (m *apiMockServer) findResponseSchema(operation *spec.Operation) *spec.Schema {
	if operation == nil || operation.Responses == nil {
		return nil
	}

	preferredStatusCodes := []string{responseStatus200, responseStatus201, responseStatusDef}

	for _, statusCode := range preferredStatusCodes {
		response, exists := operation.Responses[statusCode]
		if !exists || response == nil {
			continue
		}

		schema := schemaForContentType(response.Content)
		if schema != nil {
			return schema
		}
	}

	for _, response := range operation.Responses {
		if response == nil {
			continue
		}
		schema := schemaForContentType(response.Content)
		if schema != nil {
			return schema
		}
	}

	return nil
}

// schemaForContentType extracts the JSON schema from a content map,
// preferring application/json and */* content types.
func schemaForContentType(content map[string]*spec.MediaType) *spec.Schema {
	if content == nil {
		return nil
	}

	preferredTypes := []string{contentTypeJSON, "*/*"}

	for _, contentType := range preferredTypes {
		mediaType, exists := content[contentType]
		if !exists || mediaType == nil {
			continue
		}
		return mediaType.Schema
	}

	for _, mediaType := range content {
		if mediaType == nil {
			continue
		}
		return mediaType.Schema
	}

	return nil
}
