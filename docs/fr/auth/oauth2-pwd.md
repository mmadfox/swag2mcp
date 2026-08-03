# OAuth2 Password Grant

## Objectif

Authentification par OAuth2 Resource Owner Password Grant — utilisant le nom d'utilisateur et le mot de passe d'un utilisateur. Convient aux applications propriétaires où l'utilisateur fait confiance à l'application avec ses identifiants.

## Quand l'utiliser

- Applications propriétaires (mobile, web)
- Intégration avec Keycloak et fournisseurs d'identité similaires
- Quand l'API prend en charge OAuth2 Password Grant

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
      type: oauth2-pwd
      config:
        client_id: "$(ID_CLIENT)"
        client_secret: "$(SECRET_CLIENT)"
        username: "$(NOM_UTILISATEUR)"
        password: "$(MOT_DE_PASSE)"
        token_url: "https://auth.example.com/oauth/token"
        scopes:
          - openid
          - profile
        request_format: form
```

## Paramètres

| Paramètre | Requis | Description |
|-----------|--------|-------------|
| `client_id` | Oui | Identifiant du client |
| `username` | Oui | Nom d'utilisateur |
| `password` | Oui | Mot de passe |
| `token_url` | Oui | URL du point d'accès du jeton |
| `client_secret` | Non | Secret du client (optionnel, pour les clients publics) |
| `scopes` | Non | Liste des permissions (optionnel) |
| `request_format` | Non | Format du corps : `form` (défaut, `application/x-www-form-urlencoded`) ou `json` (`application/json`) |

## Notes

- `client_secret` est optionnel — les **clients publics** sont pris en charge (par exemple, Keycloak)
- swag2mcp renouvelle automatiquement le jeton à son expiration
- Le jeton est mis en cache jusqu'à l'expiration
- `client_id`, `client_secret`, `username` et `password` prennent en charge la syntaxe `$(VAR)` pour les variables d'environnement
- `token_url` et `scopes` sont utilisés tels quels (pas de résolution de variable d'environnement)
- `request_format: json` envoie la demande de jeton en corps JSON — à utiliser lorsque le point d'accès du jeton nécessite `Content-Type: application/json`
