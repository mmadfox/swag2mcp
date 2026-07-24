# Aucune

## Objectif

Aucune authentification requise. L'API est accessible sans jetons ni clés.

## Quand l'utiliser

- API publiques (Open-Meteo, icanhazdadjoke, PokéAPI)
- Environnements de test et de démonstration
- Quand l'API ne nécessite pas d'autorisation

## Configuration

Définissez `type: none` ou omettez simplement la section `auth` :

```yaml
specs:
  - domain: jokes
    llm_title: API Dad Joke
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Blagues
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: none
```

## Paramètres

Aucun.

## Notes

- Si la section `auth` est complètement absente de la configuration, cela équivaut à `type: none`
- Aucun en-tête d'autorisation n'est ajouté aux requêtes
