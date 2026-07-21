package swag2mcp

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

//go:generate go install go.uber.org/mock/mockgen@latest

//go:generate mockgen -source=internal/server/mcp/handler.go -destination=internal/server/mcp/mock_svc_test.go -package=mcp

//go:generate mockgen -source=internal/service/deps.go -destination=internal/service/mock_deps_test.go -package=service -typed
