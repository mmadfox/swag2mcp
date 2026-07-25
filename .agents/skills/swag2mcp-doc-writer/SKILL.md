---
name: swag2mcp-doc-writer
description: |
  Documentation writing and maintenance for swag2mcp.
  Use when the user asks to write, edit, add, remove, or restructure
  documentation files. Enforces multi-language consistency across all
  8 locales and follows project writing conventions.
license: MIT
metadata:
  author: mmadfox
  version: "1.0.0"
---

# swag2mcp-doc-writer — Documentation Writing Skill

## When this skill activates

Activate when the user asks to write, edit, add, remove, or restructure any documentation file in the `docs/` directory. This includes:

- Creating new pages or sections
- Updating existing content (commands, flags, config fields, examples)
- Adding or removing files
- Renaming or moving files
- Updating the VitePress config (`config.mjs`)

## Critical Rule: Multi-Language Consistency

The documentation is maintained in **8 locales**. The root `docs/` directory is English. Every change to an English file **MUST** be replicated identically to all 7 translated locales:

| Locale | Directory | Prefix |
|--------|-----------|--------|
| English (default) | `docs/` | `/` |
| Russian | `docs/ru/` | `/ru/` |
| German | `docs/de/` | `/de/` |
| French | `docs/fr/` | `/fr/` |
| Spanish | `docs/es/` | `/es/` |
| Chinese | `docs/zh-CN/` | `/zh-CN/` |
| Japanese | `docs/ja/` | `/ja/` |
| Korean | `docs/ko/` | `/ko/` |

### Synchronization rules

1. **File structure is identical** across all locales — same subdirectories, same filenames, same number of files
2. **Content is copied as-is** from English to all locales. Translation is done separately — the skill only ensures structural consistency
3. **When adding a file** — create it in English first, then copy to all 7 locale directories
4. **When deleting a file** — remove from English and all 7 locale directories
5. **When renaming a file** — rename in English and all 7 locale directories, then update all links in `config.mjs` and all `.md` files that reference the old path
6. **When updating content** — update the English file, then copy the exact same content to the corresponding file in each locale directory
7. **`config.mjs`** — only one file exists at `docs/.vitepress/config.mjs`. It contains sidebar and nav definitions for all 8 locales. Update it once.

### Verification command

After any change, run:

```bash
npm run docs:build
```

The build must complete without errors. If it fails, fix the issues before reporting success.

---

## Documentation Structure

### Root files (3)

| File | Purpose |
|------|---------|
| `index.md` | Landing page — hero, architecture SVG, "Who needs this" table, license |
| `faq.md` | Frequently asked questions |
| `troubleshooting.md` | Common issues and solutions |

### Section directories (10 sections, 67 files)

```
getting-started/   (2 files)  — installation.md, quickstart.md
concepts/          (6 files)  — overview.md, specs.md, collections.md, tags.md, endpoints.md, workspace.md
configuration/     (7 files)  — config-file.md, global-settings.md, spec-settings.md, collection-settings.md, http-client.md, mcp-server.md, cascade.md
cli/               (14 files) — overview.md, add.md, clean.md, delete.md, export.md, import.md, info.md, init.md, ls.md, mcp.md, run.md, update.md, validate.md, version.md
mcp-tools/         (6 files)  — overview.md, discovery.md, endpoints.md, execution.md, utilities.md, skills.md
auth/              (10 files) — overview.md, none.md, basic.md, bearer.md, api-key.md, digest.md, hmac.md, oauth2-cc.md, oauth2-pwd.md, script.md
advanced/          (8 files)  — search.md, ratelimit.md, response-size.md, caching.md, mock-server.md, tui.md, export-import.md, env-vars.md
integration/       (5 files)  — opencode.md, cursor.md, claude.md, vscode.md, crush.md
examples/          (2 files)  — cli-workflow.md, llm-session.md
development/       (7 files)  — overview.md, project-structure.md, building.md, testing.md, conventions.md, new-auth.md, new-tool.md
```

**Total: 70 `.md` files per locale × 8 locales = 560 files.**

---

## Writing Conventions

### Headings

- Start every page with `# Title` (H1) — no YAML frontmatter
- Use `## Section` (H2) for major sections
- Use `### Subsection` (H3) for subsections
- Do NOT use `####` or deeper unless absolutely necessary
- Keep heading text concise and descriptive

### Code blocks

Use the correct language tag for every code block:

```markdown
```bash
swag2mcp init
```

```yaml
domain: meteo
base_url: https://api.example.com
```

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp"]
    }
  }
}
```

```go
func main() {
    fmt.Println("hello")
}
```

```powershell
swag2mcp.exe init
```
```

### Platform-specific code (code groups)

Use `::: code-group` when showing the same command for different platforms:

```markdown
::: code-group

```bash [macOS / Linux]
swag2mcp init
```

```powershell [Windows]
swag2mcp.exe init
```

:::
```

### Notes and callouts

Use simple blockquotes. Do NOT use `::: warning`, `::: tip`, `[!NOTE]`, or other admonition syntax:

```markdown
> Some IDEs require a restart after adding skills.

> **Always use an absolute path** to the workspace directory in IDE config.
```

### Tables

Use standard GFM tables with `|` separators. Align columns for readability:

```markdown
| Method | macOS | Linux | Windows |
|--------|-------|-------|---------|
| One-liner | ✅ | ✅ | ❌ |
| Homebrew | ✅ | ✅ | ❌ |
```

### Lists

Use `-` for unordered lists. Bold the term, then `—` for the description:

```markdown
- **Cooldown:** 10 seconds per endpoint
- **Scope:** Per-endpoint — calling endpoint A does not affect endpoint B
```

### Links

- **Internal (same section):** `./specs` or `./quickstart`
- **Internal (cross-section):** `/cli/overview` or `/getting-started/installation`
- **External:** full URL `https://github.com/mmadfox/swag2mcp`
- **Locale links in config.mjs:** always include the locale prefix: `/ru/getting-started/installation`

### Emoji

Use sparingly. Only these are approved:

- ✅ for supported/compatible
- ❌ for not supported
- 🚧 for work in progress

### Horizontal rules

Use `---` to separate major sections in long pages (installation, FAQ, troubleshooting).

### ASCII hierarchy diagrams

For showing tree structures (spec → collection → tag → endpoint):

```
Spec (domain, e.g. "meteo")
  └── Collection 1 (spec file, e.g. forecast.yml)
        └── Tag 1 (category)
              └── Endpoint (GET /api/forecast)
```

### Inline HTML

Only use in `index.md` for:
- The WIP banner (`<div style="...">`)
- Image embedding with links (`<a href="..."><img src="..."></a>`)
- The architecture SVG (`<img src="/architecture.svg">`)

Do NOT use inline HTML in any other file.

### Inline formatting

- **Bold** for emphasis, UI labels, key terms
- `Backticks` for commands, file paths, field names, code, flags (`--flag`)
- Do NOT use `_italic_` — prefer bold

---

## Section-Specific Writing Rules

### `getting-started/installation.md`

- Start with a compatibility table showing all methods × platforms
- Group methods by platform (macOS, Linux, Windows)
- Each method gets a `###` heading
- Use `::: code-group` for platform variants within a method
- Include a "Mock Server" subsection
- Include an "Install via LLM Agent" subsection with exact prompts
- End with "Verify" and "Next Steps"

### `getting-started/quickstart.md`

- Step-by-step: init → skills → IDE config → MCP start → verify
- Use `::: code-group` for platform-specific commands
- Include example LLM queries table
- End with "What's Next" links

### `cli/*.md`

- Each CLI command gets its own file
- Start with `# Command Name` (e.g., `# swag2mcp init`)
- Show usage: `swag2mcp <command> [path] [flags]`
- List all flags in a table: Flag, Shorthand, Type, Default, Description
- Show examples with ` ```bash `
- Describe behavior in plain terms

### `configuration/*.md`

- Document every YAML field with: Field, Type, Required, Default, Description
- Show complete YAML examples
- Explain cascade behavior (global → spec → collection)
- Note which fields support `$(VAR)` environment variable resolution

### `auth/*.md`

- Each auth method gets its own file
- Show YAML config example
- Show environment variable usage with `$(VAR)`
- Include any method-specific notes

### `integration/*.md`

- Each client gets its own file
- Show stdio config first, then HTTP if applicable
- End with `## Others` section explaining the universal pattern
- Config JSON must match real MCP client formats

### `advanced/*.md`

- Each advanced feature gets its own file
- Start with overview, then how it works, then configuration
- Use diagrams (ASCII) where helpful
- Include important notes at the end

### `index.md`

- Hero section with WIP banner
- One-liner description
- YouTube preview image
- Architecture SVG
- "Who needs this" table
- License

### `faq.md`

- Q&A format with `###` headings for each question
- Group related questions under `##` section headings
- Link to relevant detailed docs

### `troubleshooting.md`

- Problem → Cause → Solution format
- Group by category under `##` headings
- Include specific commands to diagnose and fix

---

## VitePress Config (`config.mjs`)

The config file at `docs/.vitepress/config.mjs` defines:

- **`nav`** — top navigation bar (Home, Quick Start, Installation, Integration, GitHub)
- **`sidebar`** — left sidebar with all sections and page links
- **`i18n`** — array of 7 locale objects, each with translated `nav` and `sidebar`

When adding a new page:
1. Add the file to English `docs/` and all 7 locale directories
2. Add the link to the English `sidebar` array
3. Add the link to each locale's `sidebar` array in the `i18n` section

When removing a page:
1. Remove the file from English and all 7 locale directories
2. Remove the link from the English `sidebar` array
3. Remove the link from each locale's `sidebar` array in the `i18n` section

---

## Verification Checklist

Before reporting completion, always:

1. **Build the docs:**
   ```bash
   npm run docs:build
   ```
   Must complete with `build complete in X.XXs.` and no errors.

2. **Check file counts** (optional, for structural changes):
   ```bash
   for lang in ru de fr es zh-CN ja ko; do echo "$lang: $(find docs/$lang -name '*.md' | wc -l) files"; done
   echo "en: $(find docs -maxdepth 1 -name '*.md' | wc -l) + $(for d in getting-started concepts configuration cli mcp-tools auth advanced integration examples development; do find "docs/$d" -name '*.md'; done | wc -l) files"
   ```
   All locales should have the same count.

3. **Verify locale files exist** for any newly created page:
   ```bash
   for lang in ru de fr es zh-CN ja ko; do test -f "docs/$lang/<path>" && echo "$lang: OK" || echo "$lang: MISSING"; done
   ```
