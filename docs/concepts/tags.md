# Tags

A tag is a category grouping related endpoints within a collection.

## Hierarchy

```
Spec (domain, e.g. "petstore")
  └── Collection (spec file, e.g. petstore.yaml)
        └── Tag "pets"
              └── GET /pet
              └── POST /pet
        └── Tag "pet_images"
              └── GET /pet/{id}/image
```

## How Tags Are Created

Tags come from the OpenAPI spec:

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

## MCP Tools for Tags

| Tool | Description |
|------|-------------|
| `tag_by_spec` | All tags in a spec |
| `tag_by_collection` | Tags in a collection |
| `tag_by_id` | Tag details |
| `endpoint_by_tag` | Endpoints in a tag |

## Example

```
Query: "Show all tags in the pet collection"
→ tag_by_collection(collectionId: "...")
→ Result: pets (5 methods), pet_images (1 method)
```
