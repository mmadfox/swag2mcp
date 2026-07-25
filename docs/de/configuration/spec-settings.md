# Spec-Einstellungen

Spec-Einstellungen definieren einen API-Dienst und überschreiben globale Einstellungen für diese bestimmte API. Jede Spec repräsentiert eine logische API (z. B. "Open-Meteo Weather APIs") und kann mehrere Collections (Spezifikationsdateien) enthalten.

## Spec-Abschnitt

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    llm_instruction: "Verwenden Sie diese API für Wettervorhersagen und Klimadaten"
    base_url: https://api.open-meteo.com
    disable: false
    tags: ["weather", "climate"]
    http_client:
      timeout: 10s
      max_response_size: 1024
    auth:
      type: bearer
      config:
        token: "$(TOKEN)"
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

## Parameter

### domain

- **Typ:** `string`
- **Erforderlich:** Ja
- **Beschreibung:** Eindeutiger Identifikator für diese API-Spec. Wird intern verwendet, um auf die Spec zu verweisen.
- **Regeln:** 1-60 Zeichen. Nur Kleinbuchstaben (`a-z`), Ziffern (`0-9`), Bindestriche (`-`) und Unterstriche (`_`).
- **Beispiel:** `meteo`, `binance`, `my-api`

### llm_title

- **Typ:** `string`
- **Erforderlich:** Ja
- **Beschreibung:** Menschenlesbarer Name, den der LLM zur Referenzierung dieser API verwendet. Wird in MCP-Tool-Antworten angezeigt.
- **Regeln:** 5-120 Zeichen. Nur Buchstaben, Ziffern, Leerzeichen und einfache Satzzeichen.
- **Beispiel:** `Open-Meteo Weather APIs`, `Binance Market Data`

### llm_instruction

- **Typ:** `string`
- **Standard:** `""`
- **Beschreibung:** Anweisungen für den LLM zur Verwendung dieser API. Beschreibt, was die API tut und wann sie verwendet werden soll.
- **Regeln:** Max. 500 Zeichen. Nur Buchstaben, Ziffern, Leerzeichen und einfache Satzzeichen.
- **Beispiel:** `"Verwenden Sie diese API für Wettervorhersagen, aktuelle Bedingungen und Klimadaten."`

### base_url

- **Typ:** `string`
- **Erforderlich:** Ja
- **Beschreibung:** Basis-URL für alle API-Anfragen in dieser Spec. Die Endpunkt-Pfade aus der OpenAPI-Spezifikation werden an diese URL angehängt.
- **Beispiel:** `https://api.open-meteo.com`, `https://api.binance.com`
- **Hinweis:** Kann auf Collection-Ebene überschrieben werden, wenn verschiedene Collections unterschiedliche Basis-URLs verwenden.

### disable

- **Typ:** `bool`
- **Standard:** `false`
- **Beschreibung:** Wenn `true`, wird diese Spec von MCP-Tools ausgeschlossen. Sie wird nicht geladen, indiziert oder dem LLM zur Verfügung gestellt.
- **Wann verwenden:** Eine API vorübergehend deaktivieren, ohne sie aus der Konfiguration zu entfernen. Nützlich für APIs, die nicht verfügbar, veraltet oder in Wartung sind.

### tags

- **Typ:** `[]string` (Array von Zeichenfolgen)
- **Standard:** `[]`
- **Beschreibung:** Tags zum Filtern von Specs. Wird mit dem Flag `--tags` in CLI-Befehlen verwendet (`ls`, `validate`, `mcp`, `update`).
- **Beispiel:** `["public", "weather"]`, `["internal", "production"]`
- **Wirkung:** Wenn Sie `swag2mcp mcp --tags=public` ausführen, werden nur Specs mit dem Tag `public` geladen.

### http_client

- **Typ:** `object`
- **Standard:** erbt von Global
- **Beschreibung:** Globale HTTP-Client-Einstellungen für diese Spec überschreiben. Alle Einstellungen aus dem globalen `http_client` können überschrieben werden: `timeout`, `max_response_size`, `user_agent`, `follow_redirects`, `max_redirects`, `random`, `proxy`, `headers`, `cookies`.
- **Beispiel:**
  ```yaml
  http_client:
    timeout: 60s
    max_response_size: 4194304
    headers:
      "X-DC": "us-east-1"
  ```

### auth

- **Typ:** `object`
- **Standard:** `none` (keine Authentifizierung)
- **Beschreibung:** Authentifizierungskonfiguration für diese Spec. Siehe den Abschnitt [Authentifizierung](/auth/overview) für alle 9 Methoden und ihre Parameter.
- **Beispiel:**
  ```yaml
  auth:
    type: bearer
    config:
      token: "$(API_TOKEN)"
  ```

### collections

- **Typ:** `[]object` (Array von Collections)
- **Erforderlich:** Ja (mindestens 1)
- **Beschreibung:** Liste der OpenAPI/Swagger/Postman-Spezifikationsdateien, die zu dieser Spec gehören. Jede Collection ist eine Spezifikationsdatei.
- **Regeln:** 1-30 Collections pro Spec.
- **Siehe:** [Collection-Einstellungen](./collection-settings) für alle Collection-Parameter.

## Deaktivieren einer Spec

Deaktivierte Specs werden nicht geladen oder indiziert. Der LLM kann sie weder sehen noch verwenden.

```yaml
specs:
  - domain: old-api
    llm_title: Old API
    base_url: https://old-api.example.com
    disable: true
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## HTTP-Client-Überschreibung

Alle `http_client`-Einstellungen der globalen Ebene können auf Spec-Ebene überschrieben werden. Die Spec-Werte haben Vorrang vor globalen Werten nur für diese Spec.

```yaml
specs:
  - domain: slow-api
    llm_title: Slow API
    base_url: https://slow-api.example.com
    http_client:
      timeout: 120s
      max_response_size: 8388608
      headers:
        "X-DC": "us-east-1"
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Proxy-Überschreibung

Wenn diese Spec einen anderen Proxy als den globalen benötigt, konfigurieren Sie ihn auf Spec-Ebene:

```yaml
specs:
  - domain: proxied-api
    llm_title: Proxied API
    base_url: https://api.example.com
    http_client:
      proxy:
        url: http://proxy.company.com:8080
        username: $(PROXY_USER)
        password: $(PROXY_PASS)
        bypass:
          - "*.local"
          - "10.0.0.0/8"
    collections:
      - llm_title: Main
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo.json
```
