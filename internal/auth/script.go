package auth

import (
	"context"
	"encoding/json"
	"errors"
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
	scriptDefaultExpiresIn = 3600
	windowsOS              = "windows"
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

type scriptTokenResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

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

func (c *ScriptAuthClient) Type() Type {
	return ScriptAuth
}

func (c *ScriptAuthClient) Apply(req *http.Request, out *Info) error {
	c.mu.Lock()
	if c.token != "" && time.Now().Before(c.expiresAt) {
		setAuthHeader(req, out, "Authorization", "Bearer "+c.token)
		c.mu.Unlock()
		return nil
	}
	c.mu.Unlock()

	token, expiresIn, execErr := c.execute()
	if execErr != nil {
		return fmt.Errorf("script auth: %w", execErr)
	}

	c.mu.Lock()
	c.token = token
	c.expiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
	setAuthHeader(req, out, "Authorization", "Bearer "+c.token)
	c.mu.Unlock()
	return nil
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) //nolint:mnd // Script execution timeout.
	defer cancel()

	var cmd *exec.Cmd
	if runtime.GOOS == windowsOS {
		cmd = exec.CommandContext(ctx, "cmd", "/c", scriptPath)
	} else {
		cmd = exec.CommandContext(ctx, "sh", scriptPath)
	}

	output, runErr := cmd.Output()
	if runErr != nil {
		return "", 0, fmt.Errorf("execute script %s: %w", scriptPath, runErr)
	}

	outStr := strings.TrimSpace(string(output))

	var sr scriptTokenResponse
	if unmarshalErr := json.Unmarshal([]byte(outStr), &sr); unmarshalErr != nil {
		return "", 0, fmt.Errorf("script must output valid JSON with 'token' field, got: %s", outStr)
	}

	if sr.Token == "" {
		return "", 0, errors.New("script JSON response missing 'token' field")
	}

	expiresIn := sr.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = scriptDefaultExpiresIn
	}

	return sr.Token, expiresIn, nil
}

func (c *ScriptAuthClient) Validate() error {
	return authValidator.Struct(c)
}
