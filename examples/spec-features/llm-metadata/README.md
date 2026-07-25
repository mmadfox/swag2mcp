# LLM Metadata

Demonstrates how to provide rich metadata for LLM agents using `llm_title`
and `llm_instruction` fields. These fields help the LLM understand when and
how to use each API.

## What it demonstrates

- `llm_title` at spec level — a short, descriptive name for the LLM
- `llm_instruction` at spec level — detailed guidance on when to use this API
- `llm_title` at collection level — per-collection naming
- `llm_instruction` at collection level — per-collection guidance
- These fields appear in the tool definitions and available specs list

## Expected behavior

- The LLM sees descriptive titles and instructions in tool definitions
- The `spec_list` tool returns the `llm_title` for each spec
- The `collection_by_spec` tool returns per-collection titles
- The LLM can make better decisions about which API to use
