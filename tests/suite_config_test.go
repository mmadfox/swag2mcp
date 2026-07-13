package tests

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScript_Validate_ValidConfig(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: test-api
    llm_title: Test API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	writeConfig(t, ws, configContent)

	stdout, _, code := runCommandInWS(t, ws, "validate", ".")
	assertEqual(t, "exit code", code, 0)
	assertContains(t, "stdout", stdout, "valid")
}

func TestScript_Validate_DuplicateDomain(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
  - domain: petstore
    llm_title: Duplicate
    base_url: https://api.example.com
    collections:
      - title: Store
        location: ./testdata/petstore.yaml
`
	writeConfig(t, ws, configContent)

	_, stderr, code := runCommandInWS(t, ws, "validate", ".")
	assertNotEqual(t, "exit code", code, 0)
	assertContains(t, "stderr", stderr, "duplicate")
}

func TestScript_Validate_InvalidDomainFormat(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: "UPPERCASE INVALID"
    llm_title: Bad API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	writeConfig(t, ws, configContent)

	_, stderr, code := runCommandInWS(t, ws, "validate", ".")
	assertNotEqual(t, "exit code", code, 0)
	assertContains(t, "stderr", stderr, "Domain")
}

func TestScript_Validate_UnreachableLocation(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: test-api
    llm_title: Test API
    base_url: https://api.example.com
    collections:
      - title: Missing
        location: ./nonexistent.yaml
`
	writeConfig(t, ws, configContent)

	_, _, code := runCommandInWS(t, ws, "validate", ".")
	assertNotEqual(t, "exit code", code, 0)
}

func TestScript_Validate_TagFilter(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: public-api
    llm_title: Public API
    base_url: https://api.example.com
    tags: ["public"]
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
  - domain: internal-api
    llm_title: Internal API
    base_url: https://api.example.com
    tags: ["internal"]
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	writeConfig(t, ws, configContent)

	_, _, code := runCommandInWS(t, ws, "validate", "-t", "public", ".")
	assertEqual(t, "exit code", code, 0)
}

func TestScript_AddSpec_FromYAML(t *testing.T) {
	ws := newTestWorkspace(t)

	yamlData := `domain: added-spec
llm_title: Added Spec
base_url: https://api.example.com
collections:
  - title: Pets
    location: ./testdata/petstore.yaml
`
	stdout, stderr, code := runCommandInWS(t, ws, "add", "spec", "--yaml", yamlData, ".")
	assertEqual(t, "exit code", code, 0)
	assertContains(t, "stdout", stdout+stderr, "added")

	stdout2, _, _ := runCommandInWS(t, ws, "ls", ".")
	assertContains(t, "stdout", stdout2, "added-spec")
}

func TestScript_AddSpec_FromStdin(t *testing.T) {
	ws := newTestWorkspace(t)

	yamlData := `domain: stdin-spec
llm_title: Stdin Spec
base_url: https://api.example.com
collections:
  - title: Pets
    location: ./testdata/petstore.yaml
`
	stdout, stderr, code := runCommandWithStdinInWS(t, ws, yamlData, "add", "spec", "--yaml", "-", ".")
	assertEqual(t, "exit code", code, 0)
	assertContains(t, "stdout", stdout+stderr, "added")

	stdout2, _, _ := runCommandInWS(t, ws, "ls", ".")
	assertContains(t, "stdout", stdout2, "stdin-spec")
}

func TestScript_AddSpec_InvalidYAML(t *testing.T) {
	ws := newTestWorkspace(t)

	_, _, code := runCommandInWS(t, ws, "add", "spec", "--yaml", "invalid: [yaml: broken", ".")
	assertNotEqual(t, "exit code", code, 0)
}

func TestScript_AddCollection_FromYAML(t *testing.T) {
	ws := newTestWorkspace(t)

	specYAML := `domain: test-api
llm_title: Test API
base_url: https://api.example.com
collections:
  - title: Existing
    location: ./testdata/petstore.yaml
`
	runCommandInWS(t, ws, "add", "spec", "--yaml", specYAML, ".")

	collectionYAML := `spec_domain: test-api
llm_title: Added Collection
location: ./testdata/petstore.yaml
`
	stdout, stderr, code := runCommandInWS(t, ws, "add", "collection", "--yaml", collectionYAML, ".")
	assertEqual(t, "exit code", code, 0)
	assertContains(t, "stdout", stdout+stderr, "added")

	stdout2, _, _ := runCommandInWS(t, ws, "ls", ".")
	assertContains(t, "stdout", stdout2, "Added Collection")
}

func TestScript_DeleteSpec(t *testing.T) {
	ws := newTestWorkspace(t)

	specYAML := `domain: to-delete
llm_title: To Delete
base_url: https://api.example.com
collections:
  - title: Pets
    location: ./testdata/petstore.yaml
`
	runCommandInWS(t, ws, "add", "spec", "--yaml", specYAML, ".")

	stdout, _, _ := runCommandInWS(t, ws, "ls", ".")
	assertContains(t, "stdout", stdout, "to-delete")

	_, _, code := runCommandWithStdinInWS(t, ws, "5\ny\n", "delete", "spec", ".")
	assertEqual(t, "exit code", code, 0)

	stdout2, _, _ := runCommandInWS(t, ws, "ls", ".")
	assertNotContains(t, "stdout", stdout2, "to-delete")
}

func TestScript_DeleteSpec_Cancel(t *testing.T) {
	ws := newTestWorkspace(t)

	specYAML := `domain: keep-me
llm_title: Keep Me
base_url: https://api.example.com
collections:
  - title: Pets
    location: ./testdata/petstore.yaml
`
	runCommandInWS(t, ws, "add", "spec", "--yaml", specYAML, ".")

	runCommandWithStdinInWS(t, ws, "n\n", "delete", "spec", ".")

	stdout, _, _ := runCommandInWS(t, ws, "ls", ".")
	assertContains(t, "stdout", stdout, "keep-me")
}

func TestScript_ListSpecs(t *testing.T) {
	ws := newTestWorkspace(t)

	specYAML := `domain: list-test
llm_title: List Test
base_url: https://api.example.com
collections:
  - title: Pets
    location: ./testdata/petstore.yaml
`
	runCommandInWS(t, ws, "add", "spec", "--yaml", specYAML, ".")

	stdout, _, code := runCommandInWS(t, ws, "ls", ".")
	assertEqual(t, "exit code", code, 0)
	assertContains(t, "stdout", stdout, "list-test")
	assertContains(t, "stdout", stdout, "List Test")
	assertContains(t, "stdout", stdout, "petstore.yaml")
}

func TestScript_ListSpecs_Empty(t *testing.T) {
	ws := newTestWorkspace(t)

	_, _, code := runCommandInWS(t, ws, "ls", ".")
	assertEqual(t, "exit code", code, 0)
}

func TestScript_ListSpecs_TagFilter(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: public-api
    llm_title: Public API
    base_url: https://api.example.com
    tags: ["public"]
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
  - domain: internal-api
    llm_title: Internal API
    base_url: https://api.example.com
    tags: ["internal"]
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	writeConfig(t, ws, configContent)

	stdout, _, _ := runCommandInWS(t, ws, "ls", "-t", "public", ".")
	assertContains(t, "stdout", stdout, "public-api")
	assertNotContains(t, "stdout", stdout, "internal-api")
}

func TestScript_Update_ReCachesSpecs(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: test-api
    llm_title: Test API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	writeConfig(t, ws, configContent)

	_, stderr, code := runCommandInWS(t, ws, "update", ".")
	assertEqual(t, "exit code", code, 0)
	assertContains(t, "stderr", stderr, "updated")
}

func TestScript_Update_InvalidConfig(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: "INVALID DOMAIN"
    llm_title: Bad
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./nonexistent.yaml
`
	writeConfig(t, ws, configContent)

	_, _, code := runCommandInWS(t, ws, "update", ".")
	assertNotEqual(t, "exit code", code, 0)
}

func TestScript_Clean_RemovesCache(t *testing.T) {
	ws := newTestWorkspace(t)

	root := wsDir(ws)
	cacheDir := filepath.Join(root, "cache")
	_ = os.MkdirAll(cacheDir, 0755)
	dummyFile := filepath.Join(cacheDir, "test.cache")
	_ = os.WriteFile(dummyFile, []byte("data"), 0644)

	responsesDir := filepath.Join(root, "responses")
	_ = os.MkdirAll(responsesDir, 0755)
	dummyResp := filepath.Join(responsesDir, "test.json")
	_ = os.WriteFile(dummyResp, []byte("{}"), 0644)

	stdout, _, code := runCommandInWS(t, ws, "clean", ".")
	assertEqual(t, "exit code", code, 0)
	assertContains(t, "stdout", stdout, "Removed")

	if _, err := os.Stat(dummyFile); !os.IsNotExist(err) {
		t.Errorf("cache file was not removed")
	}
	if _, err := os.Stat(dummyResp); !os.IsNotExist(err) {
		t.Errorf("response file was not removed")
	}
}

func TestScript_Clean_PreservesSpecs(t *testing.T) {
	ws := newTestWorkspace(t)

	root := wsDir(ws)
	specsDir := filepath.Join(root, "specs")
	_ = os.MkdirAll(specsDir, 0755)
	specFile := filepath.Join(specsDir, "test.yaml")
	_ = os.WriteFile(specFile, []byte("spec: test"), 0644)

	runCommandInWS(t, ws, "clean", ".")

	if _, err := os.Stat(specFile); os.IsNotExist(err) {
		t.Errorf("specs directory was cleaned but should be preserved")
	}
}

func TestScript_EnvVarResolution(t *testing.T) {
	ws := newTestWorkspace(t)
	t.Setenv("TEST_BASE_URL", "https://env-test.example.com")

	configContent := `specs:
  - domain: env-test
    llm_title: Env Test
    base_url: https://env-test.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	writeConfig(t, ws, configContent)

	_, _, code := runCommandInWS(t, ws, "validate", ".")
	assertEqual(t, "exit code", code, 0)
}

func TestScript_ConfigCascade(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `http_client:
  timeout: 10s
  headers:
    X-Global: "true"
specs:
  - domain: cascade-test
    llm_title: Cascade Test
    base_url: https://spec.example.com
    http_client:
      timeout: 30s
      headers:
        X-Spec: "spec-only"
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
        http_client:
          headers:
            X-Collection: "collection-only"
`
	writeConfig(t, ws, configContent)

	_, _, code := runCommandInWS(t, ws, "validate", ".")
	assertEqual(t, "exit code", code, 0)
}
