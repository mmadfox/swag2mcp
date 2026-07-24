# Umgebungsvariablen

## Übersicht

swag2mcp unterstützt die Substitution von Umgebungsvariablen in der Konfigurationsdatei mit der Syntax `$(VAR_NAME)`. Dadurch können Sie vertrauliche Daten (Tokens, Passwörter, Schlüssel) aus der YAML-Datei heraushalten.

## Wie es funktioniert

Beim Start durchsucht swag2mcp die Konfiguration nach `$(VAR_NAME)`-Mustern und ersetzt sie durch den Wert der entsprechenden Umgebungsvariable.

```yaml
specs:
  - domain: meteo
    auth:
      type: bearer
      config:
        token: "$(API_TOKEN)"
```

Wenn die Umgebungsvariable `API_TOKEN` gesetzt ist, wird sie ersetzt. Wenn sie nicht gesetzt ist, wird der Wert leer.

## Wo `$(VAR)` aufgelöst wird

| Feld | Beispiel |
|------|----------|
| Auth `token` (bearer) | `token: "$(API_TOKEN)"` |
| Auth `username` / `password` (basic, digest) | `password: "$(API_PASSWORD)"` |
| Auth `client_id` / `client_secret` (oauth2-cc, oauth2-pwd) | `client_secret: "$(OAUTH_SECRET)"` |
| Auth `api_key` / `secret_key` (hmac) | `api_key: "$(BINANCE_API_KEY)"` |
| Auth `domain` (script) | `domain: "$(AUTH_DOMAIN)"` |
| MCP-Server-Token | `token: "$(MCP_TOKEN)"` |
| HTTP-Client-Header | `"X-API-Key": "$(API_KEY)"` |
| HTTP-Client-Cookie-Werte | `value: "$(SESSION_TOKEN)"` |

## Wo `$(VAR)` NICHT aufgelöst wird

- Basis-URLs (`base_url`)
- Collection-Speicherorte (`location`)
- Spec-Domain-Namen (`domain`)

## Beispiel

```bash
export API_TOKEN="eyJhbGciOiJIUzI1NiIs..."
export MCP_TOKEN="my-secret-token"

swag2mcp mcp
```

## Sicherheitshinweise

- **Speichern Sie Geheimnisse niemals direkt** in der YAML-Datei
- Verwenden Sie Umgebungsvariablen oder einen externen Geheimnisverwalter
- Fügen Sie die YAML-Datei zu `.gitignore` hinzu, wenn sie hartcodierte Geheimnisse enthält
- Setzen Sie Umgebungsvariablen in Ihrem Shell-Profil, Ihrer IDE-Konfiguration oder Ihrer Bereitstellungspipeline

## Syntax-Details

- `$(VAR_NAME)` — Standardsyntax
- `$( VAR_NAME )` — Leerzeichen innerhalb der Klammern sind erlaubt und werden entfernt
- `$()` — leerer Variablenname gibt die ursprüngliche Zeichenfolge unverändert zurück
- Verschachtelte `$(...)`-Muster werden nicht aufgelöst
