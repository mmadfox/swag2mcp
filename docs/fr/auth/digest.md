# Authentification Digest

## Objectif

Authentification HTTP Digest Access — une alternative plus sécurisée à Basic Auth. Le mot de passe n'est pas envoyé en texte clair ; des hachages MD5 sont utilisés à la place.

## Quand l'utiliser

- API legacy qui ne prennent en charge que Digest
- Quand vous avez besoin d'authentification sans envoyer le mot de passe en texte clair
- Systèmes d'entreprise internes

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
      type: digest
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

- swag2mcp envoie d'abord une requête sans authentification, reçoit un défi du serveur (HTTP 401), calcule la réponse et réessaie avec l'en-tête `Authorization: Digest ...`
- Le défi est mis en cache pendant 5 minutes — les requêtes suivantes n'ont pas besoin d'un aller-retour supplémentaire
- Stockez le mot de passe dans une variable d'environnement : `password: "$(MOT_DE_PASSE_API)"`
