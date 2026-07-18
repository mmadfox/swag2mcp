---
name: response_compress
---

# response_compress

Reduces a JSON value inside a saved response file so it fits within the response size limit and can be returned to the LLM inline.

## When to use

Use this tool after `response_outline` when you want to see a representative sample of a large array or reduce verbose string/object content before reading specific items.

## When NOT to use

- Do **NOT** use `bash`, `cat`, `head`, `tail`, `file`, `open`, `less`, `more`, or any external command to read `fileRef.path`.
- Do **NOT** read the file manually. Always use `response_compress` or `response_slice` to access data inside saved response files.

## Parameters

- `path` (required): The absolute file path from `fileRef.path` returned by `invoke`.
- `jsonPath` (optional): Path to the value to compress. Default is the root of the file.
- `mode` (required): Compression strategy.
  - `first_of_array`: keep only the first element of an array.
  - `sample_array`: keep a head and tail sample of an array.
  - `truncate_strings`: shorten every string to `stringLen` characters.
  - `keys_only`: replace object values with type names.
  - `select_keys`: keep only the keys listed in `selectKeys` for every object in an array.
- `arrayHead` (optional): Number of leading array items for `sample_array`. Default is 3.
- `arrayTail` (optional): Number of trailing array items for `sample_array`. Default is 2.
- `stringLen` (optional): Maximum string length for `truncate_strings`. Default is 80.
- `selectKeys` (optional): Keys to keep for `select_keys` mode.

## Returns

Either:
- `body`: the compressed JSON value inline, or
- `fileRef`: if the compressed result is still too large, a new saved file path and metadata.
- `hint`: a short explanation of what was compressed and how to continue exploring.

## Recommended workflow

After `response_outline` shows a large array such as `pets` with 5000 items:

```
response_compress({
  "path": "/.../responses/...json",
  "jsonPath": "pets",
  "mode": "first_of_array"
})
```

Then use `response_slice` with `jsonPath` like `pets.0`, `pets.1`, etc.
