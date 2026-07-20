# Mock Server Examples

## Basic Configuration

```yaml
mock_enabled: true

specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
        base_mock_url: "127.0.0.1:9090"
```

## Multiple Specs with Mock

```yaml
mock_enabled: true

specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
        base_mock_url: "127.0.0.1:9090"

  - domain: pokemon
    llm_title: PokeAPI
    base_url: https://pokeapi.co
    collections:
      - llm_title: Pokemon
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/pokeapi.yaml
        base_mock_url: "127.0.0.1:9091"
```

## Launch

```bash
swag2mcp-mock mockserver
```

With TLS:

```bash
swag2mcp-mock mockserver --tls --tls-cert cert.pem --tls-key key.pem
```
