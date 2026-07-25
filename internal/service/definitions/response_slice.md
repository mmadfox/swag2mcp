---
name: response_slice
---

# response_slice

Extracts a specific fragment of a saved JSON response file by logical jsonPath or by line range.

## When to use

Use this tool when you know which object, array, or field you want to inspect inside a large response. Prefer `jsonPath` over line numbers because it is stable and descriptive.

## When NOT to use

- Do **NOT** use `bash`, `cat`, `head`, `tail`, `file`, `open`, `less`, `more`, or any external command to read `fileRef.path`.
- Do **NOT** read the file manually. This tool is the only allowed way to extract fragments from saved response files.

## Parameters

- `path` (required): The absolute file path from `fileRef.path` returned by `invoke`.
- `jsonPath` (optional): Logical path such as `data.0`, `users.3.name`, or `company.departments.engineering.employees.0`. Use gjson dotted syntax.
- `line` (optional): 1-based line number to center the fragment on. The tool returns `around` lines above and below.
- `range` (optional): Exact line range as `start-end` (for example `120-240`).
- `around` (optional): Number of lines to include around `line`. Default is 20.

## Returns

- `slice.lines`: 1-based line range of the returned fragment.
- `slice.value`: the extracted JSON value parsed into a structured object.
- `slice.fragment`: raw JSON text when the fragment is small enough to include.
- `slice.context`: `object`, `array`, or `value` describing what was extracted.
- `slice.isComplete`: true when `value` is a valid JSON fragment.
- `slice.nextPath` / `slice.prevPath`: suggested adjacent paths for array navigation.
- `slice.nextLine` / `slice.prevLine`: suggested line numbers for line-based navigation.
- `fileRef`: only present when the extracted fragment exceeded the size limit and was saved to disk.

## Example

```
response_slice({
  "path": "/.../responses/...json",
  "jsonPath": "pets.0"
})
```

Then continue with `pets.1`, `pets.2`, etc. using `slice.nextPath`.
