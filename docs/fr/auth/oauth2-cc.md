# OAuth2 Client Credentials

## Objectif

Authentification par OAuth2 Client Credentials Grant — pour la communication serveur à serveur. L'application obtient un jeton en utilisant son client_id et client_secret, sans intervention de l'utilisateur.

## Quand l'utiliser

- Microservices et intégrations serveur à serveur
- Communication machine à machine
- Quand l'API utilise OAuth2 et que vous avez un client_id + client_secret

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
      type: oauth2-cc
      config:
        client_id: "$(ID_CLIENT)"
        client_secret: "$(SECRET_CLIENT)"
        token_url: "https://auth.example.com/oauth/token"
        scopes:
          - read
          - write
        request_format: form
```

## Paramètres

| Paramètre | Requis | Description |
|-----------|--------|-------------|
| `client_id` | Oui | Identifiant du client |
| `client_secret` | Oui | Secret du client |
| `token_url` | Oui | URL du point d'accès du jeton |
| `scopes` | Non | Liste des permissions (optionnel) |
| `request_format` | Non | Format du corps : `form` (défaut, `application/x-www-form-urlencoded`) ou `json` (`application/json`) |

## Notes

- swag2mcp demande automatiquement un nouveau jeton lorsque le jeton actuel expire
- Le jeton est mis en cache jusqu'à sa date d'expiration (`expires_in`)
- Si le serveur ne fournit pas `expires_in`, le jeton est considéré comme valide pendant 1 heure
- `client_id` et `client_secret` prennent en charge la syntaxe `$(VAR)` pour les variables d'environnement
- `token_url` et `scopes` sont utilisés tels quels (pas de résolution de variable d'environnement)
- `request_format: json` envoie la demande de jeton en corps JSON — à utiliser lorsque le point d'accès du jeton nécessite `Content-Type: application/json`
