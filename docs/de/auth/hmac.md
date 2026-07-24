# HMAC Auth

## Zweck

HMAC-SHA256-Anfragesignierung — die Authentifizierungsmethode, die von Kryptowährungsbörsen (Binance, Bybit und anderen) verwendet wird. Jede Anfrage wird mit einem geheimen Schlüssel signiert.

## Wann verwenden

- Binance-API und Binance-kompatible Börsen
- Kryptowährungs-Handelsplattformen
- APIs, die eine Anfragesignierung erfordern

## Konfiguration

```yaml
specs:
  - domain: binance
    llm_title: Binance Market Data
    base_url: https://api.binance.com
    collections:
      - llm_title: Market Data
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/binance.yaml
    auth:
      type: hmac
      config:
        api_key: "$(BINANCE_API_KEY)"
        secret_key: "$(BINANCE_SECRET_KEY)"
```

## Parameter

| Parameter | Erforderlich | Beschreibung |
|-----------|-------------|--------------|
| `api_key` | Ja | Öffentlicher API-Schlüssel |
| `secret_key` | Ja | Geheimer Schlüssel zum Signieren |

## Hinweise

- swag2mcp fügt automatisch einen Zeitstempel (Unix in Millisekunden) zu jeder Anfrage hinzu
- Die Signatur wird aus allen Anfrageparametern berechnet
- Speichern Sie Schlüssel in Umgebungsvariablen: `api_key: "$(BINANCE_API_KEY)"`
- Diese Methode ist kompatibel mit der Binance-API und ähnlichen Börsen
