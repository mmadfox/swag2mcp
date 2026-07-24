# Caching

## Overview

swag2mcp caches downloaded spec files so the MCP server starts faster on subsequent runs. Instead of downloading the same spec file every time, it reuses the cached copy.

## How caching works

When you add a spec with a remote URL, swag2mcp downloads it and saves it to the `cache/` directory. On the next startup, it checks if the cached copy is still fresh. If it is, the download is skipped.

### What gets cached

| Source | Behavior |
|--------|----------|
| **Remote URL** (http/https) | Always cached. Downloaded once, reused until the cache expires. |
| **Local file in `specs/`** | Used directly from the `specs/` directory. Never cached — changes are immediately visible. |
| **Local file outside `specs/`** | Copied to the cache. If the source file changes (modification time), the cache is invalidated. |

### Cache expiration (TTL)

Each cached file gets a random expiration time between **1 hour and 48 hours**. The randomness prevents all cached files from expiring at the same time (which would cause a thundering herd of downloads).

- The TTL is reset every time the MCP server starts
- If a cached file is still within its TTL, it is reused
- If the TTL has expired, the file is downloaded again

### Cache structure

```
~/.swag2mcp/cache/
├── a1b2c3d4e5f6a7b8.spec    # Cached spec file
├── a1b2c3d4e5f6a7b8.meta    # Metadata (source, TTL, cached at)
├── b2c3d4e5f6a7b8c9.spec
├── b2c3d4e5f6a7b8c9.meta
└── ...
```

The cache key is derived from the spec file URL or path. Each cached file has a companion `.meta` file that stores when it was cached and when it expires.

## Managing the cache

### Force a refresh

Run `swag2mcp update` to clear the entire cache and re-download all spec files:

```bash
swag2mcp update
```

This validates the config, clears the cache, and downloads everything fresh.

### Clear the cache manually

```bash
swag2mcp clean
```

This removes all cached spec files and saved API responses. The next time you start the MCP server, all specs will be downloaded again.

### Automatic cleanup

When the MCP server starts (`swag2mcp mcp`), saved API responses older than 48 hours are automatically removed. This prevents the `responses/` directory from growing indefinitely.

## Important notes

- **Local files in `specs/` are never cached** — if you edit a spec file directly in the `specs/` directory, the changes are immediately visible without clearing the cache
- **Remote URLs are always cached** — there is no way to bypass the cache for remote URLs except by running `swag2mcp update` or `swag2mcp clean`
- **The cache is local** — it is stored on disk and does not sync between machines. Use `swag2mcp export` and `swag2mcp import` to transfer specs between machines
