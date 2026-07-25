package swag2mcp

//go:generate go install go.uber.org/mock/mockgen@latest

//go:generate mockgen -source=internal/server/mcp/handler.go -destination=internal/server/mcp/mock_svc_test.go -package=mcp
