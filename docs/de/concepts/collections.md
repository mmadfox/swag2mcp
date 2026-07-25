# Collections

Eine Collection ist eine einzelne OpenAPI/Swagger/Postman-Datei, die eine bestimmte API beschreibt. Sie verweist auf einen `location` (URL oder lokalen Dateipfad) und gehört zu einer Spec (Domain).

Eine Spec kann mehrere Collections haben — zum Beispiel könnte die "meteo"-Spec die Collections "Forecast", "Air Quality" und "Marine" haben, die jeweils auf eine andere Spezifikationsdatei verweisen.

## Collection-Felder

| Feld | YAML-Schlüssel | Erforderlich | Beschreibung |
|------|---------------|-------------|--------------|
| [LLM-Titel](#llm-instruction) | `llm_title` | ❌ | Collection-Anzeigename für den LLM (max. 120 Zeichen). Wird automatisch aus dem Spec-Dokument befüllt, wenn nicht gesetzt |
| [LLM-Instruktion](#llm-instruction) | `llm_instruction` | ❌ | Kurzer Hinweis für den LLM (max. 360 Zeichen). Wird automatisch aus dem Spec-Dokument befüllt, wenn nicht gesetzt |
| Titel | `title` | ❌ | Überschreibung des ursprünglichen Spec-Titels (wird automatisch aus dem geparsten Dokument befüllt) |
| [Speicherort](#location--wie-spezifikationsdateien-aufgelöst-werden) | `location` | ✅ | URL oder Pfad zur Spezifikationsdatei (5–250 Zeichen) |
| [Deaktivieren](#disable) | `disable` | ❌ | Diese Collection beim Laden überspringen |
| [HTTP-Client](#http-client-überschreibung) | `http_client` | ❌ | Pro-Collection-HTTP-Einstellungen (Header, Cookies) |
| [Basis-URL](#basis-url-überschreibung) | `base_url` | ❌ | Die Basis-URL der Spec für diese Collection überschreiben |
| [Mock-Server](#mock-server) | `base_mock_url` | ❌ | Mock-Server-Adresse im Format `host:port`. Erforderlich, wenn `mock_enabled: true` |

## Speicherort — Wie Spezifikationsdateien aufgelöst werden

Das Feld `location` teilt swag2mcp mit, wo die OpenAPI/Swagger/Postman-Datei zu finden ist. Es unterstützt mehrere Quelltypen:

| Quelle | Beispiel | Beschreibung |
|--------|---------|--------------|
| **Entfernte URL** | `https://raw.githubusercontent.com/.../spec.yaml` | Heruntergeladen und zwischengespeichert |
| **Lokale Datei (absolut)** | `/home/user/my-api.yaml` | Vom Dateisystem gelesen, zwischengespeichert |
| **Lokale Datei (relativ)** | `./my-api.yaml` | Zu absolutem Pfad aufgelöst, zwischengespeichert |
| **Lokale Datei im Arbeitsbereich** | `specs/my-api.yaml` | In `~/.swag2mcp/specs/` gespeichert, direkt verwendet (nicht zwischengespeichert) |
| **file://-URI** | `file:///home/user/spec.yaml` | In lokalen Pfad umgewandelt, zwischengespeichert |

swag2mcp erkennt den Quelltyp automatisch:

- `https://` oder `http://` → entfernte URL (zwischengespeichert)
- `file://` → lokale Datei (in Dateisystempfad umgewandelt)
- Alles andere → lokale Datei (mit `~`-Erweiterung für das Home-Verzeichnis)

### Entfernte URLs

Wenn Sie eine entfernte URL verwenden, lädt swag2mcp die Datei herunter und speichert sie lokal zwischen. Der Cache wird bei nachfolgenden Starts wiederverwendet, um wiederholte Downloads zu vermeiden.

### Lokale Dateien

Lokale Dateien werden direkt vom Dateisystem gelesen. Wenn die Datei außerhalb des Arbeitsbereichsverzeichnisses `specs/` liegt, wird sie zur Konsistenz in den Cache kopiert.

### Lokale Dateien im Arbeitsbereich

Das Verzeichnis `specs/` innerhalb des Arbeitsbereichs (`~/.swag2mcp/specs/`) ist der empfohlene Ort für lokale Spezifikationsdateien. Hier gespeicherte Dateien werden direkt ohne Zwischenspeicherung verwendet. Verwenden Sie einen relativen Pfad, der mit `specs/` beginnt, um darauf zu verweisen.

> **Hinweis:** `specs/` ist nur ein Verzeichnisname (wie `cache/` oder `responses/`), nicht das Konzept "Spec". Es speichert die eigentlichen OpenAPI/Swagger/Postman-Dateien, auf die Collections verweisen.

```bash
# Eine Spezifikationsdatei in den Arbeitsbereich importieren
swag2mcp import https://example.com/api.yaml myspec

# Nach dem Import wird der Speicherort:
# specs/myspec.yaml
```

## Cache-System

swag2mcp speichert entfernte Spezifikationsdateien zwischen, um sie nicht bei jedem Start herunterladen zu müssen.

### Wie es funktioniert

1. Wenn eine Collection mit einer entfernten URL geladen wird, überprüft swag2mcp den Cache
2. Wenn ein gültiger (nicht abgelaufener) Cache-Eintrag existiert, wird er direkt verwendet
3. Wenn nicht, wird die Datei heruntergeladen, geparst und im Cache gespeichert

### Cache-Struktur

```
~/.swag2mcp/
  cache/
    {sha256_hash}.spec    # Zwischengespeicherter Spezifikationsdateiinhalt
    {sha256_hash}.meta    # Cache-Metadaten (JSON)
```

Jede zwischengespeicherte Datei hat eine Metadatendatei mit:

```json
{
  "source": "https://example.com/api.yaml",
  "source_type": "url",
  "cached_at": "2024-01-01T00:00:00Z",
  "mod_time": "2024-01-01T00:00:00Z",
  "ttl_sec": 3600
}
```

### Cache-TTL

Jede zwischengespeicherte Datei erhält eine **zufällige TTL** zwischen 1 Stunde und 48 Stunden. Dies verhindert, dass alle zwischengespeicherten Dateien gleichzeitig ablaufen (Thundering-Herd-Problem).

### Cache-Schlüssel

Der Cache-Schlüssel ist ein SHA-256-Hash der rohen Speicherort-Zeichenfolge (erste 16 Bytes = 32 Hex-Zeichen).

### Cache verwalten

```bash
# Cache und Antworten leeren, alle Spezifikationsdateien neu herunterladen
swag2mcp update

# Nur Cache und Antworten leeren
swag2mcp clean
```

- `swag2mcp update` — validiert Konfiguration, leert `cache/` und `responses/`, dann werden alle Collection-Speicherorte neu zwischengespeichert
- `swag2mcp clean` — entfernt alle Inhalte von `cache/` und `responses/`, plus verwaiste Auth-Skripte
- Alte Antworten werden automatisch nach 48 Stunden beim MCP-Server-Start bereinigt

## Validierung

Jede Collection wird beim Laden der Konfiguration validiert. Die Validierung läuft bei jedem `swag2mcp mcp`-Start. Wenn sie fehlschlägt, startet der MCP-Server nicht — in manchen IDEs bedeutet dies, dass der Server einfach keine Verbindung herstellt, und der LLM erhält eine klare Fehlermeldung, die erklärt, was zu beheben ist.

| Prüfung | Regel |
|---------|-------|
| **Speicherort** | Erforderlich, 5–250 Zeichen |
| **Speicherort-Erreichbarkeit** | Muss eine erreichbare URL oder vorhandene Datei sein |
| **Speicherort-Gültigkeit** | Muss eine gültige OpenAPI 3.x-, Swagger 2.0- oder Postman-Datei sein |
| **LLM-Titel** | Max. 120 Zeichen, Buchstaben/Ziffern/einfache Satzzeichen |
| **LLM-Instruktion** | Max. 360 Zeichen, gleicher Zeichensatz wie Titel |
| **Basis-URL** | Muss eine gültige URL sein, wenn gesetzt |
| **Basis-Mock-URL** | Muss `host:port` oder `host:port/pfad` sein, wobei host `localhost`, `127.0.0.1` oder `0.0.0.0` ist |
| **Mock erforderlich** | Wenn `mock_enabled: true`, muss jede Collection `base_mock_url` haben |
| **Doppelte Mock-Ports** | Keine zwei Collections dürfen denselben Mock-Port teilen |

Verwenden Sie den Befehl [`validate`](../cli/validate.md), um Probleme vor dem Start des Servers zu diagnostizieren:

```bash
# Standard-Arbeitsbereich validieren (~/.swag2mcp)
swag2mcp validate

# Benutzerdefinierten Projektarbeitsbereich validieren
swag2mcp validate ./my-project
```

## Collections hinzufügen

### Über YAML-Konfiguration

Bearbeiten Sie `~/.swag2mcp/swag2mcp.yaml` direkt:

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

Nach der Bearbeitung starten Sie den MCP-Server (`swag2mcp mcp`) neu, damit die Änderungen wirksam werden.

### Über CLI

```bash
# Interaktiver Modus
swag2mcp add collection

# Nicht-interaktiv mit YAML
swag2mcp add collection --yaml 'spec_domain: meteo
llm_title: Forecast
location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml'

# Von stdin weiterleiten
cat collection.yaml | swag2mcp add collection --yaml -

# YAML-Beispiel anzeigen
swag2mcp add collection --example
```

### Über Import

```bash
# Eine Spezifikationsdatei in den Arbeitsbereich importieren
swag2mcp import https://example.com/api.yaml
```

## LLM-Instruktion

Collections können eine eigene `llm_instruction` (bis zu 360 Zeichen) für spezifischere Anleitungen haben. Diese wird zusammen mit der Spec-Ebene-Instruktion in den swag2mcp-System-Prompt eingefügt.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        llm_instruction: "Verwenden Sie diese Collection für aktuelles Wetter und tägliche Vorhersagen."
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        llm_instruction: "Verwenden Sie diese Collection für Luftqualitätsindex und Verschmutzungsdaten."
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
```

Wenn `llm_title` nicht gesetzt ist, wird es automatisch aus dem `title`-Feld des Spec-Dokuments befüllt. Wenn `llm_instruction` nicht gesetzt ist, wird es aus dem `description`-Feld des Spec-Dokuments befüllt.

## Deaktivieren

Setzen Sie `disable: true`, um eine Collection zu überspringen. Sie wird nicht geladen, indiziert oder dem LLM zur Verfügung gestellt.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
        disable: true
```

## Basis-URL-Überschreibung

Jede Collection kann die `base_url` der Spec überschreiben. Dies ist nützlich, wenn verschiedene Collections innerhalb derselben Spec unterschiedliche API-Endpunkte verwenden.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
      - llm_title: Marine
        base_url: https://marine-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

## HTTP-Client-Überschreibung

Collections können HTTP-Einstellungen (Header, Cookies) von der Spec- und globalen Ebene überschreiben.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          headers:
            X-API-Version: "2"
          cookies:
            - name: session
              value: abc123
```

Einstellungen kaskadieren: global → spec → collection. Siehe [Konfigurationskaskade](../configuration/cascade.md) für Details.

## Mock-Server

Wenn `mock_enabled: true` auf Konfigurationsebene gesetzt ist, muss jede Collection `base_mock_url` gesetzt haben. Dies teilt swag2mcp mit, wo der Mock-Server für diese Collection läuft.

```yaml
mock_enabled: true
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        base_mock_url: localhost:8080
```

Siehe [Mock-Server](../advanced/mock-server.md) für vollständige Details.

## Beispiele

### Minimale Collection

```yaml
specs:
  - domain: dadjokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

### Vollständige Collection mit allen Feldern

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        llm_instruction: "Für aktuelles Wetter und tägliche Vorhersagen verwenden."
        title: "Benutzerdefinierter Titel"
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        disable: false
        base_url: https://forecast-api.open-meteo.com
        base_mock_url: localhost:8080
        http_client:
          headers:
            X-Custom: value
```

### Mehrere Collections pro Spec

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
      - llm_title: Marine
        base_url: https://marine-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

### Lokale Datei im Arbeitsbereich (Verzeichnis specs/)

```yaml
specs:
  - domain: myapi
    llm_title: My Internal API
    base_url: https://api.mycompany.com
    collections:
      - llm_title: Users
        location: specs/users.openapi.json
      - llm_title: Orders
        location: specs/orders.openapi.json
```

### Deaktivierte Collection

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
        disable: true
```

## Verwandte Themen

- [Collection-Einstellungen (Konfiguration)](../configuration/collection-settings.md) — vollständige YAML-Referenz
- [Konfigurationskaskade](../configuration/cascade.md) — wie Einstellungen einander überschreiben
- [Specs](./specs) — logische Container für Collections
- [HTTP-Client](../configuration/http-client.md) — HTTP-Client-Konfiguration
- [Mock-Server](../advanced/mock-server.md) — Mock-Server-Einrichtung
- [CLI: validate](../cli/validate.md) — validate-Befehlsreferenz
- [CLI: update](../cli/update.md) — update-Befehlsreferenz
- [CLI: clean](../cli/clean.md) — clean-Befehlsreferenz
