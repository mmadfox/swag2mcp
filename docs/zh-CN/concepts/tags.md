# Tags

A tag is a category that groups related endpoints within a collection. Tags may or may not exist — not all collections have them, and a collection can have any number of tags.

Tags come from the OpenAPI/Swagger/Postman file itself. There are **no YAML config settings** for tags — you cannot create, rename, or delete tags in `swag2mcp.yaml`. The only way to change tags is to edit the original spec file.

## Hierarchy

```
Spec (domain, e.g. "meteo")
  └── Collection (spec file, e.g. forecast.yml)
        └── Tag "weather"
              └── GET /forecast
              └── GET /forecast/hourly
        └── Tag "alerts"
              └── GET /alerts
```

## How Tags Are Created

Tags are extracted from the spec document during parsing:

**OpenAPI 3.x / Swagger 2.0** — each operation's `tags` list becomes tags:

```yaml
paths:
  /pet:
    get:
      tags: ["pets"]
      summary: "Find pet by ID"
    post:
      tags: ["pets"]
      summary: "Add a new pet"
  /pet/{petId}/uploadImage:
    post:
      tags: ["pet_images"]
      summary: "Uploads an image"
```

**Postman** — each top-level folder becomes a tag. Nested folders use the last folder name.

If an endpoint has no tags, it is placed under a `"default"` tag.

## Purpose

Tags help the LLM find groups of related endpoints. Instead of searching through every endpoint in a collection, the LLM can first find the right tag, then list only the endpoints within it.

## MCP Tools for Tags

| Tool | Description |
|------|-------------|
| `tag_by_spec` | All tags across an entire spec |
| `tag_by_collection` | Tags within a specific collection |
| `tag_by_id` | Tag details (title, method count) |
| `endpoint_by_tag` | Endpoints grouped under a tag |

## Example

```
Query: "Show all tags in the pet collection"
→ tag_by_collection(collectionId: "...")
→ Result: pets (5 methods), pet_images (1 method)
```

## Limitations

- Tags are read-only from the config perspective. To add, rename, or remove tags, edit the original OpenAPI/Swagger/Postman file and run `swag2mcp update`.
- Tags cannot be filtered or disabled per-collection in the YAML config.
