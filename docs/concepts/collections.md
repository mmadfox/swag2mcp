# Collections

A collection is a logical group of endpoints within a spec. One spec can have multiple collections.

## How Collections Are Created

Collections are created automatically when parsing a spec:

=== "OpenAPI 3.x"
    Each top-level `tag` becomes a collection.

=== "Swagger 2.0"
    Each tag from the tags list becomes a collection.

=== "Postman"
    Each top-level folder becomes a collection.

## Example

From the Petstore spec:

```yaml
tags:
  - name: pet
    description: Everything about your Pets
  - name: store
    description: Access to Petstore orders
  - name: user
    description: Operations about user
```

Collections created: `pet`, `store`, `user`.

## Overriding Collections

In YAML config:

```yaml
specs:
  - domain: "petstore.swagger.io"
    location: "https://petstore.swagger.io/v2/swagger.json"
    collections:
      - name: "animals"
        tags: ["pet"]
      - name: "orders"
        tags: ["store"]
      - name: "accounts"
        tags: ["user"]
```

## Tag Filtering

Limit endpoints in a collection:

```yaml
collections:
  - name: "pets"
    tags: ["pet", "pets"]
    filter:
      include:
        - method: GET
        - path: "/pet/*"
      exclude:
        - method: DELETE
```

## Management

```bash
# List collections in a spec
swag2mcp ls --spec <spec_id>

# Delete a collection
swag2mcp delete <spec_id>/<collection_id>
```
