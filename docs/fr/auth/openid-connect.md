# OpenID Connect (OIDC Discovery)

## Objectif

Authentification basée sur OpenID Connect Discovery. swag2mcp récupère le document de découverte OIDC auprès de l'émetteur, découvre le point de terminaison du jeton et obtient un jeton Bearer via le flux **client credentials**. Le jeton est mis en cache jusqu'à son expiration.

## Quand l'utiliser

- API qui exposent un document de découverte OIDC (`.well-known/openid-configuration`)
- Intégrations serveur à serveur avec un `client_id` + `client_secret`
- Lorsque l'API utilise OIDC et que vous souhaitez une acquisition automatique du jeton

## Configuration

```yaml
specs:
  - domain: my-api
    llm_title: My API
    base_url: https://api.example.com
    collections:
      - llm_title: Main
        location: https://example.com/spec.yaml
    auth:
      type: openid-connect
      config:
        issuer: "https://auth.example.com"
        client_id: "$(CLIENT_ID)"
        client_secret: "$(CLIENT_SECRET)"
        scopes:
          - openid
          - profile
```

## Paramètres

| Paramètre | Requis | Description |
|-----------|----------|-------------|
| `issuer` | Oui | URL de l'émetteur OIDC. swag2mcp ajoute `/.well-known/openid-configuration` pour découvrir le point de terminaison du jeton |
| `client_id` | Oui | Identifiant du client |
| `client_secret` | Oui | Secret du client |
| `scopes` | Non | Liste des autorisations (optionnel) |

## Remarques

- swag2mcp demande automatiquement un nouveau jeton lorsque le courant expire
- Le jeton est mis en cache jusqu'à son expiration (`expires_in`)
- Si le serveur ne fournit pas `expires_in`, le jeton est valide 1 heure
- `issuer`, `client_id` et `client_secret` prennent en charge `$(VAR)` pour les variables d'environnement
- Le jeton est obtenu via le flux **client credentials** (serveur à serveur, sans utilisateur)
- Le document de découverte doit contenir un champ `token_endpoint`
