# MCP-Server

Der MCP-Server ist der Hauptinteraktionspunkt für LLM-Agenten. Er stellt alle konfigurierten APIs als MCP-Tools bereit, die der LLM aufrufen kann.

## Konfiguration

```yaml
mcp:
  transport: stdio
  addr: ":8080"
  path: "/mcp"
  auth:
    token: ""
```

## Transports

Drei Transporttypen sind verfügbar:

| Transport | Beschreibung | Wann verwenden |
|-----------|-------------|----------------|
| `stdio` | Standard-Ein-/ausgabe | Lokale LLM-Clients (VS Code, Cursor, Claude Desktop) |
| `sse` | Server-Sent Events | Entfernte Clients, HTTP-basierte Kommunikation |
| `streamable-http` | HTTP mit Streaming | Web-Clients, moderne MCP-Clients |

### stdio (Standard)

Der LLM-Client führt swag2mcp als Kindprozess aus. Die Kommunikation erfolgt über Standard-Ein- und -ausgabe. Es wird kein Netzwerkport benötigt.

```yaml
mcp:
  transport: stdio
```

```bash
swag2mcp mcp
```

### SSE

Server-Sent-Events-Transport für HTTP-basierte Kommunikation. Der MCP-Server lauscht auf einem HTTP-Port und der LLM-Client verbindet sich remote.

```yaml
mcp:
  transport: sse
  addr: "127.0.0.1:8080"
  path: "/mcp"
```

```bash
swag2mcp mcp --transport sse --http-addr 127.0.0.1:8080
```

### Streamable HTTP

Moderner HTTP-Transport, der Streaming-Antworten unterstützt. Ähnlich wie SSE, verwendet aber ein anderes Protokoll.

```yaml
mcp:
  transport: streamable-http
  addr: "127.0.0.1:8080"
  path: "/mcp"
```

```bash
swag2mcp mcp --transport streamable-http --http-addr 0.0.0.0:8080
```

## Parameter

### transport

- **Typ:** `string`
- **Standard:** `"stdio"`
- **Optionen:** `stdio`, `sse`, `streamable-http`
- **Wirkung:** Bestimmt, wie der MCP-Server mit dem LLM-Client kommuniziert.

### addr

- **Typ:** `string`
- **Standard:** `":8080"`
- **Beschreibung:** Lauschadresse für SSE- und Streamable-HTTP-Transports. Format: `host:port`.
- **Beispiele:** `":8080"`, `"127.0.0.1:8080"`, `"0.0.0.0:9000"`

### path

- **Typ:** `string`
- **Standard:** `"/mcp"`
- **Beschreibung:** URL-Pfad für den MCP-Endpunkt. Der LLM-Client sendet Anfragen an `http://&lt;addr&gt;&lt;path&gt;`.
- **Beispiele:** `"/mcp"`, `"/api/mcp"`, `"/v1/mcp"`

### auth.token

- **Typ:** `string`
- **Standard:** `""` (kein Auth)
- **Beschreibung:** Bearer-Token für die HTTP-Transport-Authentifizierung. Wenn gesetzt, muss der LLM-Client bei jeder Anfrage `Authorization: Bearer &lt;token&gt;` mitsenden.
- **Hinweis:** Unterstützt die Auflösung von `$(ENV_VAR)`.

## HTTP-Authentifizierung

Schützen Sie den MCP-HTTP-Endpunkt mit einem Bearer-Token:

```yaml
mcp:
  auth:
    token: "my-secret-token"
```

Oder über CLI-Flag:

```bash
swag2mcp mcp --auth-token "my-secret-token"
```

## Health-Check

Der MCP-Server bietet einen Health-Check-Endpunkt, der ohne MCP-Initialisierung funktioniert:

```bash
curl http://127.0.0.1:8080/health
# {"status":"ok","version":"v1.2.0"}
```

## Start-Flags

CLI-Flags überschreiben die YAML-Konfiguration. Wenn ein Flag nicht gesetzt ist, wird der Wert aus dem `mcp`-Abschnitt in YAML als Fallback verwendet.

| Flag | Typ | Standard | Beschreibung |
|------|-----|----------|--------------|
| `--transport` | string | `"stdio"` | Transporttyp: `stdio`, `sse`, `streamable-http` |
| `--http-addr` | string | `":8080"` | HTTP-Serveradresse (für SSE und Streamable HTTP) |
| `--http-path` | string | `"/mcp"` | URL-Pfad für den MCP-Handler |
| `--auth-token` | string | `""` | Bearer-Token für HTTP-Transport-Authentifizierung |
| `--logfile` | string | `""` | Log-Dateipfad (loggt nach stderr, wenn nicht gesetzt) |
| `--disable-llm-auth` | bool | `true` | Entfernt das `auth`-Tool aus der MCP-Tool-Liste |
| `--dump-dir` | string | `""` | Verzeichnis zum Speichern von HTTP-Anfragen zum Debuggen |
| `--tags` | string | `""` | Specs nach Tags filtern (kommagetrennt) |
