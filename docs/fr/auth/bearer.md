# Authentification Bearer

## Objectif

Authentification par jeton Bearer — la méthode la plus courante pour les API REST modernes. Le jeton est envoyé dans l'en-tête `Authorization: Bearer <jeton>`.

## Quand l'utiliser

- API REST modernes
- JWT (JSON Web Tokens)
- Jetons d'accès OAuth2 (lorsque le jeton est déjà obtenu)
- Toute API qui accepte un jeton Bearer

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
      type: bearer
      config:
        token: "eyJhbGciOiJIUzI1NiIs..."
```

## Paramètres

| Paramètre | Requis | Description |
|-----------|--------|-------------|
| `token` | Oui | Jeton Bearer (JWT, jeton OAuth2, etc.) |

## Notes

- Le jeton est statique — s'il expire, vous devez le mettre à jour manuellement dans la configuration
- Pour un renouvellement automatique du jeton, utilisez `oauth2-cc` ou `oauth2-pwd`
- Stockez le jeton dans une variable d'environnement : `token: "$(JETON_API)"`
