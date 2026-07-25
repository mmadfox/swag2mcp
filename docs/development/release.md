# Release Process

## Preparation

1. Ensure all tests pass:

```bash
make testall
```

2. Check linting:

```bash
make lint
```

3. Update version if needed

## GoReleaser

Releases are automated via GoReleaser (`.goreleaser.yaml`):

```yaml
# Two binaries
- swag2mcp
- swag2mcp-mock

# Platforms
- linux/amd64, linux/arm64
- darwin/amd64, darwin/arm64
- windows/amd64

# Formats
- tar.gz (linux, macOS)
- zip (Windows)
```

## GitHub Actions

Release is triggered via `.github/workflows/release.yaml` on tag push:

```bash
git tag v1.0.0
git push origin v1.0.0
```

## Changelog

Auto-generated from commits. Excludes:

- `docs/`
- `test/`
- `ci/`
- `chore`
- `README`

## Publishing

After GitHub release creation:

1. Binaries are uploaded to GitHub Releases
2. Users can install via `go install` or download archives
