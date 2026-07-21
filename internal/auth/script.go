package auth

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	windowsOS = "windows"
)

// ScriptAuthClient holds a domain name used to locate the auth script.
// The script must be located at {workspaceDir}/auth_scripts/{domain}.sh (Unix)
// or {workspaceDir}/auth_scripts/{domain}.bat (Windows).
// The script must output a JSON object: {"token": "...", "expires_in": N}.
type ScriptAuthClient struct {
	Domain string `yaml:"domain" validate:"required,script_domain_format"`

	mu           sync.Mutex
	token        string
	expiresAt    time.Time
	workspaceDir string
}

// New resolves environment variables in Domain, trims whitespace, and validates that the domain contains no path separators.
func (c *ScriptAuthClient) New() error {
	c.Domain = resolveEnv(c.Domain)
	c.Domain = strings.TrimSpace(c.Domain)
	if strings.ContainsAny(c.Domain, "/\\") {
		return fmt.Errorf("script domain must be a name without path separators, got %q", c.Domain)
	}
	return nil
}

// SetWorkspaceDir sets the workspace directory for script resolution.
func (c *ScriptAuthClient) SetWorkspaceDir(dir string) {
	c.workspaceDir = dir
}

// Type returns the authentication type for script-based auth.
func (c *ScriptAuthClient) Type() Type {
	return ScriptAuth
}

// Apply executes the external auth script to obtain a Bearer token and sets it on the request, caching the token until expiry.
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

func (c *ScriptAuthClient) scriptPath() string {
	ext := ".sh"
	if runtime.GOOS == windowsOS {
		ext = ".bat"
	}
	return filepath.Join(c.workspaceDir, "auth_scripts", c.Domain+ext)
}

func (c *ScriptAuthClient) execute() (string, int, error) {
	scriptPath := c.scriptPath()

	ctx, cancel := context.WithTimeout(context.Background(), tokenRequestTimeout)
	defer cancel()

	var cmd *exec.Cmd
	if runtime.GOOS == windowsOS {
		cmd = exec.CommandContext(ctx, "cmd", "/c", scriptPath)
	} else {
		cmd = exec.CommandContext(ctx, "sh", scriptPath)
	}

	output, err := cmd.Output()
	if err != nil {
		return "", 0, fmt.Errorf("execute script %s: %w", scriptPath, err)
	}

	return decodeTokenResponse(bytes.NewReader(output))
}

// Validate checks that the Domain field is present and has a valid script domain format.
func (c *ScriptAuthClient) Validate() error {
	return authValidator.Struct(c)
}
