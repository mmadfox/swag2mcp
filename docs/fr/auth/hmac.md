# Authentification HMAC

## Objectif

Signature de requête HMAC-SHA256 — la méthode d'authentification utilisée par les échanges de cryptomonnaies (Binance, Bybit et autres). Chaque requête est signée avec une clé secrète.

## Quand l'utiliser

- API Binance et échanges compatibles Binance
- Plateformes de trading de cryptomonnaies
- API qui nécessitent une signature de requête

## Configuration

```yaml
specs:
  - domain: binance
    llm_title: Données de marché Binance
    base_url: https://api.binance.com
    collections:
      - llm_title: Données de marché
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/binance.yaml
    auth:
      type: hmac
      config:
        api_key: "$(CLE_API_BINANCE)"
        secret_key: "$(CLE_SECRETE_BINANCE)"
```

## Paramètres

| Paramètre | Requis | Description |
|-----------|--------|-------------|
| `api_key` | Oui | Clé API publique |
| `secret_key` | Oui | Clé secrète pour la signature |

## Notes

- swag2mcp ajoute automatiquement un horodatage (Unix en millisecondes) à chaque requête
- La signature est calculée à partir de tous les paramètres de la requête
- Stockez les clés dans des variables d'environnement : `api_key: "$(CLE_API_BINANCE)"`
- Cette méthode est compatible avec l'API Binance et les échanges similaires
