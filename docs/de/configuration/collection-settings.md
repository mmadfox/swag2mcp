# Collection-Einstellungen

Collection-Einstellungen definieren eine einzelne OpenAPI/Swagger/Postman-Spezifikationsdatei und überschreiben Spec-Einstellungen für diese bestimmte Datei. Jede Collection gehört zu einer Spec und repräsentiert ein API-Spezifikationsdokument.

## Collection-Abschnitt

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        llm_instruction: "Für aktuelle und Vorhersage-Wetterdaten verwenden"
        disable: false
        base_url: https://forecast-api.open-meteo.com
        base_mock_url: localhost:8081
        http_client:
          timeout: 5s
```

## Parameter

### llm_title

- **Typ:** `string`
- **Erforderlich:** Nein
- **Beschreibung:** Menschenlesbarer Name für diese Collection. Wird in MCP-Tool-Antworten angezeigt.
- **Regeln:** Max. 120 Zeichen. Nur Buchstaben, Ziffern, Leerzeichen und einfache Satzzeichen.
- **Beispiel:** `Forecast`, `Air Quality`, `Market Data`

### llm_instruction

- **Typ:** `string`
- **Standard:** `""`
- **Beschreibung:** Anweisungen für den LLM zu dieser spezifischen Collection. Beschreibt, welche Endpunkte diese Collection bereitstellt.
- **Regeln:** Max. 360 Zeichen. Nur Buchstaben, Ziffern, Leerzeichen und einfache Satzzeichen.
- **Beispiel:** `"Für aktuelle und Vorhersage-Wetterdaten verwenden."`

### title

- **Typ:** `string`
- **Standard:** `""`
- **Beschreibung:** Rohtitel aus der Spezifikationsdatei. Wird zur Laufzeit automatisch befüllt. Normalerweise müssen Sie dies nicht in YAML setzen.

### location

- **Typ:** `string`
- **Erforderlich:** Ja
- **Beschreibung:** URL oder lokaler Dateipfad zur OpenAPI 3.x-, Swagger 2.0- oder Postman-Collection-Spezifikationsdatei.
- **Regeln:** 5-250 Zeichen.
- **Beispiele:**
  - URL: `https://raw.githubusercontent.com/org/repo/main/spec.yaml`
  - Lokal: `./specs/my-api.json`
  - Lokal (absolut): `/home/user/.swag2mcp/specs/my-api.yaml`

### disable

- **Typ:** `bool`
- **Standard:** `false`
- **Beschreibung:** Wenn `true`, wird diese Collection von MCP-Tools ausgeschlossen. Sie wird nicht geladen oder indiziert.
- **Wann verwenden:** Eine Collection vorübergehend deaktivieren, ohne sie aus der Konfiguration zu entfernen. Nützlich, wenn eine Spezifikationsdatei aktualisiert wird oder eine API-Version veraltet ist.

### http_client

- **Typ:** `object`
- **Standard:** erbt von Spec (oder Global)
- **Beschreibung:** HTTP-Client-Einstellungen für diese Collection überschreiben. Alle Einstellungen aus dem globalen `http_client` können überschrieben werden: `timeout`, `max_response_size`, `user_agent`, `follow_redirects`, `max_redirects`, `random`, `proxy`, `headers`, `cookies`.
- **Beispiel:**
  ```yaml
  http_client:
    timeout: 120s
    headers:
      "X-Custom": "value"
    cookies:
      - name: "session"
        value: "abc123"
  ```

### base_url

- **Typ:** `string`
- **Standard:** `""` (erbt von Spec)
- **Beschreibung:** Die `base_url` der Spec-Ebene für diese Collection überschreiben. Verwenden, wenn verschiedene Collections innerhalb derselben Spec unterschiedliche Basis-URLs verwenden.
- **Beispiel:** Wenn die Spec `base_url: https://api.open-meteo.com` hat, aber eine Collection `https://air-quality-api.open-meteo.com` verwendet, setzen Sie `base_url` auf Collection-Ebene.

### base_mock_url

- **Typ:** `string`
- **Standard:** `""`
- **Beschreibung:** Mock-Server-Adresse im Format `host:port`. Erforderlich, wenn `mock_enabled: true` in der globalen Konfiguration.
- **Regeln:** Host muss `localhost`, `127.0.0.1` oder `0.0.0.0` sein. Port muss eine gültige Portnummer sein.
- **Beispiel:** `localhost:8081`, `127.0.0.1:9000`
- **Wann verwenden:** Sie haben `mock_enabled: true` und möchten diese Collection mit gefälschten Antworten testen.

## Mehrere Collections aus einer Spec

Eine Spec kann mehrere Collections haben — zum Beispiel, wenn eine API separate Spezifikationsdateien für verschiedene Dienste hat:

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

## Deaktivieren einer Collection

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
        disable: true
```

## HTTP-Client-Überschreibung

Alle `http_client`-Einstellungen können auf Collection-Ebene überschrieben werden. Collection-Werte haben Vorrang vor Spec- und Global-Werten nur für diese Collection.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          timeout: 120s
          headers:
            "X-Custom": "value"
          cookies:
            - name: "session"
              value: "abc123"
```
