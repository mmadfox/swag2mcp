---
name: response_outline
---

# response_outline

Returns a high-level structural summary of a large JSON response file that was saved to disk by `invoke`. It does not return the actual data — only the shape, keys, array lengths, and hints that help decide how to explore the file next.

## When to use

Use this tool **immediately** after `invoke` returns a `fileRef` because the response body was too large. It is the **first and mandatory** step in exploring a saved response file.

## When NOT to use

- Do **NOT** use `bash`, `cat`, `head`, `tail`, `file`, `open`, `less`, `more`, or any external command to read `fileRef.path`.
- Do **NOT** ask the user to open the file manually.
- Do **NOT** try to guess the file contents. Only the `response_*` tools may read saved response files.

## Parameters

- `path` (required): The absolute file path from `fileRef.path` returned by `invoke`.
- `maxDepth` (optional): Maximum recursion depth when inspecting nested objects and arrays. Default is 3.
- `maxArrayItems` (optional): How many array items to inspect for detailed key/type information. Default is 5.

## Returns

A structural outline containing:
- `type`: root JSON type (`object`, `array`, etc.).
- `size`: file size in bytes.
- `lineCount`: number of lines in the file.
- `depth`: maximum nesting depth inspected.
- `structure`: recursive map of keys, types, array lengths, and sample items.
- `schemaHint`: one-line summary of the top-level shape.
- `compressionHints`: suggested `response_compress` calls to shrink the file.
- `navigationHints`: top-level paths and arrays with lengths, useful for `response_slice`.

## Example workflow

```
invoke returns fileRef.path = /.../responses/example-get-pets-abc123.json
  ↓
response_outline({"path": "/.../responses/example-get-pets-abc123.json"})
  ↓
response_compress({"path": "...", "jsonPath": "pets", "mode": "first_of_array"})
  ↓
response_slice({"path": "...", "jsonPath": "pets.0"})
```
