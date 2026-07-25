# API Key

Authentication via API key.

## Configuration

=== "Header"
    ```yaml
    auth:
      type: api-key
      api_key:
        name: "X-API-Key"
        in: header
        value: "{{API_KEY}}"
    ```

=== "Query"
    ```yaml
    auth:
      type: api-key
      api_key:
        name: "api_key"
        in: query
        value: "{{API_KEY}}"
    ```

=== "Cookie"
    ```yaml
    auth:
      type: api-key
      api_key:
        name: "session"
        in: cookie
        value: "{{SESSION_ID}}"
    ```

## Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `name` | string | Parameter name |
| `in` | string | Location: header, query, cookie |
| `value` | string | Key value |
