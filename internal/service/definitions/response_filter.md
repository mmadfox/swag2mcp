---
name: response_filter
---

# response_filter

Filters, searches, and paginates through arrays in saved JSON response files.

## When to use

Use this tool when:
- You need to find specific items in a large array (e.g. "find bitcoin in 18090 items")
- You want to filter an array by a condition (e.g. "status = active", "price > 100")
- You need to paginate through a large array page by page
- The response is too large to inspect manually with `response_slice`

## When NOT to use

- Do **NOT** use `bash`, `cat`, `head`, `tail`, `file`, `open`, `less`, `more`, or any external command to read `fileRef.path`.
- Do **NOT** read the file manually. This tool is the only allowed way to filter and paginate saved response files.

## Parameters

- `path` (required): The absolute file path from `fileRef.path` returned by `invoke`.
- `jsonPath` (required): Path to the array to filter (e.g. `pets`, `data.items`, `results`).
- `search` (optional): Full-text search across all fields of each item. Case-insensitive substring match (e.g. `"bitcoin"`).
- `filter` (optional): Structured filter condition. Format: `field operator value`. Supported operators: `=`, `!=`, `contains`, `>`, `<`, `>=`, `<=`. Examples: `status = active`, `price > 100`, `name contains bitcoin`.
- `page` (optional): Page number starting from 1. Default is 1.
- `pageSize` (optional): Items per page (max 50). Default is 10.

## Returns

- `page`: Current page number.
- `pageSize`: Items per page.
- `total`: Total number of matching items.
- `totalPages`: Total number of pages.
- `items`: Array of matching items for the current page.
- `strategy`: Whether the file was processed in memory (`memory`) or streamed (`streaming`).

## Examples

```
response_filter({
  "path": "/.../responses/...json",
  "jsonPath": "pets",
  "search": "fluffy",
  "page": 1,
  "pageSize": 5
})
```

```
response_filter({
  "path": "/.../responses/...json",
  "jsonPath": "data.items",
  "filter": "price > 50",
  "page": 2,
  "pageSize": 20
})
```
