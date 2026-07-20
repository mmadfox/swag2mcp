# ls

List specs and collections.

## Syntax

```bash
swag2mcp ls [flags]
```

## Flags

| Flag | Description |
|------|-------------|
| `-s, --spec` | Spec ID to view collections |

## Usage

=== "All specs"
    ```bash
    swag2mcp ls
    ```
    ```
    Specifications:
      petstore (https://petstore.swagger.io/v2/swagger.json)
        pet (3 endpoints)
        store (4 endpoints)
        user (8 endpoints)
      meteo (https://.../forecast.yml)
        forecast (5 endpoints)
    ```

=== "Spec collections"
    ```bash
    swag2mcp ls --spec abc123...
    ```
    ```
    Collections for petstore:
      pet (3 endpoints, 2 tags)
      store (4 endpoints, 1 tag)
      user (8 endpoints, 3 tags)
    ```
