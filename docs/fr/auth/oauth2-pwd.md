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

## Notes

- `client_secret` est optionnel — les **clients publics** sont pris en charge (par exemple, Keycloak)
- swag2mcp renouvelle automatiquement le jeton à son expiration
- Le jeton est mis en cache jusqu'à l'expiration
- Tous les paramètres peuvent être stockés dans des variables d'environnement
