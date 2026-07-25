# Tags Filtering

Demonstrates how to use `tags` to organize specs and filter them at startup
with the `--tags` flag. Tags allow you to run only a subset of your specs.

## What it demonstrates

- `tags` at spec level — arbitrary labels for categorization
- `--tags` CLI flag — filter which specs are loaded
- Multiple tags per spec
- Tags are comma-separated in the CLI flag
- Specs without matching tags are skipped entirely

## Expected behavior

- With `--tags=production`: only "payments-api" and "analytics-api" are loaded
- With `--tags=staging`: only "staging-api" is loaded
- With no filter: all 4 specs are loaded
- Filtered-out specs do not appear in any tool responses
