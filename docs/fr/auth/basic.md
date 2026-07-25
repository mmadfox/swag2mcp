# Authentification Basic

## Objectif

Authentification HTTP Basic — la façon la plus simple de s'authentifier avec un nom d'utilisateur et un mot de passe.

## Quand l'utiliser

- API legacy qui ne prennent en charge que Basic Auth
- Authentification simple sans jetons complexes
- Services internes

## Configuration

```yaml
specs:
  - domain: jokes
    llm_title: API Dad Joke
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Blagues
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: basic
      config:
        username: "admin"
        password: "$(MOT_DE_PASSE)"
```

## Paramètres

| Paramètre | Requis | Description |
|-----------|--------|-------------|
| `username` | Oui | Nom d'utilisateur |
| `password` | Oui | Mot de passe |

## Notes

- Le mot de passe est envoyé dans l'en-tête `Authorization: Basic ...` encodé en Base64 — ce n'est **pas du chiffrement**. Utilisez toujours HTTPS.
- Stockez le mot de passe dans une variable d'environnement : `password: "$(MON_MOT_DE_PASSE)"`
