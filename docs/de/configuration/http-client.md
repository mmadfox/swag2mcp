# HTTP-Client

swag2mcp verwendet einen konfigurierbaren HTTP-Client für alle API-Aufrufe. Diese Einstellungen werden global definiert und können auf Spec- und Collection-Ebene überschrieben werden.

## Konfiguration

```yaml
http_client:
  timeout: 30s
  max_response_size: 1048576
  user_agent: "swag2mcp-global/1.0"
  follow_redirects: true
  max_redirects: 10
  random: false
  proxy:
    url: ""
    username: ""
    password: ""
    bypass: []
  headers:
    "Accept": "application/json"
  cookies:
    - name: "session"
      value: "abc123"
      domain: ".example.com"
      path: "/"
```

## Timeout

Steuert, wie lange swag2mcp auf eine API-Antwort wartet, bevor es aufgibt.

- **Typ:** Dauer (Go-Format: `30s`, `60s`, `2m`)
- **Standard:** `30s`
- **Bereich:** 1 Sekunde bis 5 Minuten
- **Wirkung:** Wenn die API nicht innerhalb dieser Zeit antwortet, schlägt die Anfrage mit einem Timeout-Fehler fehl.
- **Wann erhöhen:** Langsame APIs, große Nutzlasten, unzuverlässige Netzwerke.
- **Wann verringern:** Interne APIs, Health-Checks, Fast-Fail-Szenarien.

```yaml
http_client:
  timeout: 60s
```

## Maximale Antwortgröße

Begrenzt, wie groß eine Antwort sein kann, bevor swag2mcp sie auf der Festplatte speichert, anstatt sie inline an den LLM zurückzugeben.

- **Typ:** `int` (Bytes)
- **Standard:** `1048576` (1 MB)
- **Bereich:** 256 bis 10.485.760 Bytes (10 MB)
- **Wirkung:** Wenn eine Antwort dieses Limit überschreitet, wird sie als JSON-Datei in `{workspace}/responses/` gespeichert. Der LLM erhält einen Dateiverweis und kann ihn mit den Tools `response_outline`, `response_compress` und `response_slice` erkunden.
- **Wann erhöhen:** APIs, die große Datensätze zurückgeben (Berichte, Protokolle, Analysen).
- **Wann verringern:** Begrenztes LLM-Kontextfenster oder wenn Sie dateibasierten Zugriff für alle Antworten bevorzugen.

```yaml
http_client:
  max_response_size: 4194304  # 4 MB
```

## User-Agent

Der `User-Agent`-Header, der mit jeder Anfrage gesendet wird. Einige APIs erfordern einen bestimmten User-Agent oder blockieren bekannte Bot-User-Agents.

- **Typ:** `string`
- **Standard:** `"swag2mcp-global/1.0"`
- **Wirkung:** Identifiziert Ihre Anwendung gegenüber dem API-Server.
- **Wann ändern:** Die API erfordert einen bestimmten User-Agent, oder Sie möchten Ihre Anwendung für Analysen identifizieren.

```yaml
http_client:
  user_agent: "MyApp/1.0"
```

## Weiterleitungen folgen

Steuert, ob swag2mcp HTTP-Weiterleitungen (3xx-Statuscodes) automatisch folgt.

- **Typ:** `bool`
- **Standard:** `true`
- **Wirkung:** Wenn `true`, folgt swag2mcp Weiterleitungen bis zu `max_redirects`-Mal. Wenn `false`, wird die Weiterleitungsantwort unverändert zurückgegeben.
- **Wann deaktivieren:** APIs, die in einer Schleife weiterleiten, sicherheitsrelevante Endpunkte, bei denen Sie Weiterleitungsziele manuell überprüfen möchten.

```yaml
http_client:
  follow_redirects: false
```

## Maximale Weiterleitungen

Begrenzt, wie vielen Weiterleitungen swag2mcp folgt, bevor es anhält.

- **Typ:** `int`
- **Standard:** `10`
- **Bereich:** 0 bis 50
- **Wirkung:** Wenn die API öfter weiterleitet als dieses Limit, schlägt die Anfrage fehl.
- **Wann ändern:** APIs mit langen Weiterleitungsketten, oder reduzieren für schnelleren Fehlschlag bei Weiterleitungsschleifen.

```yaml
http_client:
  max_redirects: 5
```

## Randomizer

Fügt jeder Anfrage zufällige browserähnliche Header hinzu, um Fingerprinting und Blockierung zu vermeiden.

- **Typ:** `bool`
- **Standard:** `false`
- **Wirkung:** Wenn `true`, generiert swag2mcp zufällige Header für jede Anfrage: `User-Agent` (aus einem Pool echter Browser-Strings), `Accept`, `Accept-Language`, `Accept-Encoding`, `Cache-Control`. Dies überschreibt die `user_agent`-Einstellung.
- **Wann aktivieren:** APIs, die Anfragen basierend auf User-Agent oder Headermustern blockieren, Scraping-Szenarien.

```yaml
http_client:
  random: true
```

## Proxy

Ein Proxy-Server fungiert als Vermittler zwischen swag2mcp und der Ziel-API. Der gesamte HTTP-Verkehr wird durch ihn geleitet.

**Wann Sie einen Proxy benötigen könnten:**
- **Firmennetzwerk** — der gesamte ausgehende Verkehr muss über einen Firmenproxy laufen
- **Geografische Einschränkungen** — einige APIs sind regionsgesperrt, ein Proxy in der richtigen Region umgeht dies
- **Statische IP** — APIs, die eine IP-Whitelist erfordern
- **Anonymität** — die Ursprungs-IP vor der Ziel-API verbergen

### Proxy-URL

- **Typ:** `string`
- **Standard:** `""` (kein Proxy)
- **Unterstützte Schemata:** `http`, `https`, `socks5`, `socks5h`
- **Unterstützt `$(VAR)`:** ✅ zur Laufzeit aufgelöst

| Schema | Beschreibung | Anwendungsfall |
|--------|-------------|----------------|
| `http` | HTTP-Proxy für HTTP-Verkehr | Firmenproxys, einfaches Proxying |
| `https` | HTTPS-Proxy (CONNECT-Tunnel) | Sichere Firmenproxys |
| `socks5` | SOCKS5-Proxy (DNS lokal aufgelöst) | Allgemeiner Zweck, jedes Protokoll |
| `socks5h` | SOCKS5-Proxy (DNS auf Proxy aufgelöst) | Wenn der Proxy eine bessere DNS-Auflösung hat |

### Proxy-Authentifizierung

Wenn der Proxy eine Authentifizierung erfordert, geben Sie `username` und `password` an:

- **Unterstützt `$(VAR)`:** ✅ zur Laufzeit für alle drei Felder (`url`, `username`, `password`) aufgelöst

```yaml
http_client:
  proxy:
    url: "http://proxy.example.com:8080"
    username: "proxyuser"
    password: "$(PROXY_PASSWORD)"
```

### Proxy-Bypass

Eine Liste von Domains, die **nicht** über den Proxy geleitet werden sollen. Nützlich für interne Dienste, localhost oder APIs, die nur direkt erreichbar sind.

```yaml
http_client:
  proxy:
    url: "http://proxy.example.com:8080"
    bypass:
      - "localhost"
      - "127.0.0.1"
      - "*.internal.company.com"
      - "api.local"
```

Bypass unterstützt Platzhaltermuster (`*.example.com` passt auf jede Subdomain).

## Header

Benutzerdefinierte HTTP-Header, die zu jeder Anfrage hinzugefügt werden. Header werden über Kaskadenebenen hinweg zusammengeführt:

```
Globale Header → Spec-Header (zusammengeführt) → Collection-Header (zusammengeführt)
```

Collection-Header überschreiben Spec-Header, die globale Header für denselben Schlüssel überschreiben.

```yaml
http_client:
  headers:
    "Accept": "application/json"
    "Accept-Language": "en-US"
```

Header-Werte unterstützen die Auflösung von `$(ENV_VAR)`.

## Cookies

Cookies, die mit jeder Anfrage gesendet werden. Cookies werden über Kaskadenebenen hinweg zusammengeführt (niedrigere Ebene überschreibt globale für denselben Cookie-Namen).

```yaml
http_client:
  cookies:
    - name: "session"
      value: "abc123"
      domain: ".example.com"
      path: "/"
      secure: false
      http_only: false
```

### Cookie-Felder

| Feld | Erforderlich | Beschreibung |
|------|-------------|--------------|
| `name` | Ja | Cookie-Name |
| `value` | Ja | Cookie-Wert (unterstützt `$(ENV_VAR)`-Auflösung) |
| `domain` | Nein | Domain-Bereich (z. B. `.example.com`) |
| `path` | Nein | Pfad-Bereich (z. B. `/`) |
| `secure` | Nein | Nur über HTTPS senden |
| `http_only` | Nein | Nicht über JavaScript zugänglich |

## Benutzerdefinierte Header auf Spec-Ebene

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    http_client:
      headers:
        "Accept": "application/json"
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Cookies auf Spec-Ebene

```yaml
specs:
  - domain: example
    llm_title: Example API
    base_url: https://api.example.com
    http_client:
      cookies:
        - name: "session"
          value: "abc123"
        - name: "csrf"
          value: "$(CSRF_TOKEN)"
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Kaskade

HTTP-Client-Einstellungen kaskadieren von Global zu Spec zu Collection. Alle Einstellungen können auf jeder Ebene überschrieben werden:

```
Global (http_client)
    ↓ überschreibt (alle Einstellungen)
Spec (specs[].http_client)
    ↓ überschreibt (alle Einstellungen)
Collection (specs[].collections[].http_client)
```

**Alle HTTP-Client-Einstellungen** (Timeout, Proxy, User-Agent, Weiterleitungen, Antwortgröße, Randomizer, Header, Cookies) können sowohl auf Spec- als auch auf Collection-Ebene überschrieben werden.

Siehe [Konfigurationskaskade](./cascade) für Details.
