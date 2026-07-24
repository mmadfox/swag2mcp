# Mock-Server

## Übersicht

Der Mock-Server generiert gefälschte API-Antworten basierend auf Ihren OpenAPI-Schemata. Er ermöglicht es Ihnen, Ihre API-Integration zu testen, ohne echte HTTP-Aufrufe zu tätigen. Dies ist nützlich für Entwicklung, Tests von LLM-Agenten und Demonstrationen.

Der Mock-Server ist eine **separate Binärdatei** — `swag2mcp-mock`. Sie ist nicht in der Hauptbinärdatei `swag2mcp` enthalten und muss separat installiert werden.

## Installation

```bash
# Option 1: Von GitHub Releases herunterladen
# Suchen Sie nach swag2mcp-mock_<version>_<os>_<arch>.tar.gz

# Option 2: Mit Go installieren
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

## Konfiguration

Aktivieren Sie den Mock-Server in Ihrer Konfiguration:

```yaml
mock_enabled: true
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092

specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
        base_mock_url: "127.0.0.1:9090"
```

## Parameter

### mock_enabled

- **Typ:** `bool`
- **Standard:** `false`
- **Wirkung:** Wenn `true`, muss jede aktive Collection `base_mock_url` gesetzt haben. Der Mock-Server startet HTTP-Server für jede Collection.

### mock_auth

Ports für Mock-Authentifizierungsserver. Diese simulieren OAuth2-, Digest- und HMAC-Auth-Endpunkte, damit Sie authentifizierte APIs ohne echte Anmeldeinformationen testen können.

| Feld | Standard | Beschreibung |
|------|----------|--------------|
| `oauth2_port` | `9090` | Port für den Mock-OAuth2-Token-Server |
| `digest_port` | `9091` | Port für den Mock-Digest-Auth-Server |
| `hmac_port` | `9092` | Port für den Mock-HMAC-Auth-Server |

### base_mock_url (pro Collection)

- **Typ:** `string`
- **Erforderlich:** Ja (wenn `mock_enabled: true`)
- **Format:** `host:port` (z. B. `localhost:8080`, `127.0.0.1:9000`)
- **Wirkung:** Jede Collection erhält ihren eigenen HTTP-Server auf dieser Adresse. Der Server antwortet auf alle in der Spezifikation definierten Endpunkte mit zufällig generierten Daten.

## Mock-Server starten

```bash
# Mit Standardkonfiguration starten
swag2mcp-mock mockserver

# Mit TLS starten
swag2mcp-mock mockserver --tls

# Mit benutzerdefiniertem TLS-Zertifikat starten
swag2mcp-mock mockserver --tls --tls-cert cert.pem --tls-key key.pem
```

### TLS-Flags

| Flag | Beschreibung |
|------|--------------|
| `--tls` | TLS mit einem selbstsignierten Zertifikat aktivieren |
| `--tls-cert` | Pfad zur TLS-Zertifikatsdatei |
| `--tls-key` | Pfad zur TLS-Schlüsseldatei |

Wenn `--tls` ohne `--tls-cert` und `--tls-key` gesetzt ist, wird automatisch ein selbstsigniertes Zertifikat für `localhost` generiert.

## Was der Mock-Server tut

Wenn Sie den Mock-Server starten, führt er folgende Schritte aus:

1. **Analysiert alle Spezifikationsdateien** — liest die OpenAPI/Swagger-Spezifikation jeder Collection
2. **Registriert Handler** — erstellt einen HTTP-Handler für jeden Pfad und jede Methode, die in der Spezifikation definiert sind
3. **Generiert gefälschte Daten** — antwortet mit zufällig generierten Daten, die dem Antwortschema entsprechen (korrekte Typen, Formate und Struktur)
4. **Startet Auth-Server** — simuliert OAuth2-, Digest- und HMAC-Auth-Endpunkte zum Testen

### Mock testen

```bash
# In einem Terminal:
swag2mcp-mock mockserver

# In einem anderen Terminal:
curl http://localhost:8080/pets
# → [{"id":1,"name":"Pet_name","status":"available"}]
```

## Wie gefälschte Daten generiert werden

Der Mock-Server generiert realistische gefälschte Daten basierend auf dem OpenAPI-Schema:

- **Zeichenfolgen** — zufällige Wörter, Sätze oder formatspezifische Werte (E-Mail, URL, UUID, Datum, Telefon usw.)
- **Zahlen** — zufällige Ganzzahlen und Gleitkommazahlen innerhalb des angegebenen Bereichs
- **Boolesche Werte** — zufällig wahr/falsch
- **Arrays** — 1 bis 3 zufällige Elemente
- **Objekte** — alle Eigenschaften mit zufälligen Werten gefüllt
- **Aufzählungen** — zufälliger Wert aus der Aufzählungsliste
- **Nullable-Felder** — manchmal wird `null` zurückgegeben (~10% Wahrscheinlichkeit)

## Anwendungsfälle

- **Entwicklung** — Integration ohne echten API-Zugriff testen
- **Testen von LLM-Agenten** — überprüfen, ob der LLM Endpunkte entdecken, inspizieren und aufrufen kann
- **Demonstrationen** — swag2mcp ohne Konfiguration echter APIs zeigen
- **Lasttests** — MCP-Server unter Last testen, ohne echte APIs zu treffen

## Wichtige Hinweise

- **Separate Binärdatei** — `swag2mcp-mock` ist nicht in der Hauptbinärdatei `swag2mcp` enthalten. Installieren Sie sie separat.
- **Jede Collection bekommt ihren eigenen Port** — konfigurieren Sie `base_mock_url` pro Collection
- **Auth-Mock-Server sind global** — OAuth2-, Digest- und HMAC-Server laufen auf den konfigurierten Ports, unabhängig davon, wie viele Collections Sie haben
- **Spezifikations-Parsing-Fehler sind nicht fatal** — wenn die Spezifikation einer Collection nicht analysiert werden kann, wird sie mit einer Warnung übersprungen
- **Selbstsigniertes TLS** — bei Verwendung von `--tls` ohne Zertifikate wird ein selbstsigniertes Zertifikat nur für localhost generiert
