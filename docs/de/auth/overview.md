# Authentifizierung

## Übersicht

swag2mcp unterstützt **9 Authentifizierungsmethoden** für die Arbeit mit APIs, die eine Autorisierung erfordern. Sie konfigurieren es einmal in der Konfigurationsdatei — danach fügt jeder API-Aufruf über `invoke` automatisch die richtigen Tokens und Header hinzu.

### Wo konfigurieren

Die Authentifizierung wird auf **Spec**-Ebene in `swag2mcp.yaml` festgelegt:

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: bearer
      config:
        token: "my-token"
```

### Wie es funktioniert

- Sie geben den Authentifizierungstyp und die Parameter in der Konfiguration an
- swag2mcp wendet sie automatisch auf jede Anfrage an, wenn Sie `invoke` aufrufen
- Sie **müssen kein Token anfordern**, bevor Sie eine API aufrufen — es passiert automatisch
- Wenn ein Token abläuft (OAuth2, Script), erneuert swag2mcp es selbstständig

### Umgebungsvariablen

Vertrauliche Daten (Tokens, Passwörter, Schlüssel) können mit der Syntax `$(VAR_NAME)` in Umgebungsvariablen gespeichert werden:

```yaml
auth:
  type: bearer
  config:
    token: "$(MY_API_TOKEN)"
```

swag2mcp setzt den Wert von `MY_API_TOKEN` beim Start ein.

### MCP-Auth-Tool

Der LLM-Agent kann ein Token oder Header über das `auth`-MCP-Tool abrufen — zum Beispiel, um einen curl-Befehl zu erstellen oder dem Benutzer anzuzeigen.

In der **Produktion** sollte dieses Tool mit `--disable-llm-auth` deaktiviert werden (standardmäßig aktiviert), damit der LLM niemals Zugriff auf Tokens hat.

### Methoden

| Methode | Beschreibung | Am besten für |
|---------|--------------|---------------|
| [`none`](/auth/none) | Keine Authentifizierung | Öffentliche APIs |
| [`basic`](/auth/basic) | HTTP Basic (Benutzername + Passwort) | Legacy-APIs, einfache Authentifizierung |
| [`bearer`](/auth/bearer) | Bearer-Token (JWT, Token) | Moderne REST-APIs |
| [`api-key`](/auth/api-key) | API-Schlüssel in Header oder Abfrageparameter | Dienste mit API-Schlüsseln |
| [`digest`](/auth/digest) | HTTP Digest (Benutzername + Passwort) | Legacy-APIs, sicherer als Basic |
| [`hmac`](/auth/hmac) | HMAC-SHA256-Signatur (Binance-Stil) | Kryptowährungsbörsen |
| [`oauth2-cc`](/auth/oauth2-cc) | OAuth2 Client Credentials | Server-zu-Server, Microservices |
| [`oauth2-pwd`](/auth/oauth2-pwd) | OAuth2 Password Grant | Apps mit Benutzeranmeldung |
| [`script`](/auth/script) | Externes Skript zum Abrufen eines Tokens | Jedes benutzerdefinierte Auth-Schema |
