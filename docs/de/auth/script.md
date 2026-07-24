# Skript-Authentifizierung

## Zweck

Authentifizierung über ein externes Skript — die flexibelste Methode. Sie können ein Skript in jeder Sprache (bash, Python usw.) schreiben, das ein Token nach Ihren Wünschen abruft und an swag2mcp zurückgibt.

## Wann verwenden

- Benutzerdefinierte oder nicht standardmäßige Authentifizierungsschemata
- Komplexe Token-Beschaffungslogik (mehrstufig, mit zusätzlichen Prüfungen)
- Wenn keine der Standardmethoden Ihren Anforderungen entspricht

## Konfiguration

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: script
      config:
        domain: "my-auth"
```

## Parameter

| Parameter | Erforderlich | Beschreibung |
|-----------|-------------|--------------|
| `domain` | Ja | Skriptdateiname (ohne Erweiterung) |

## Skript-Speicherort

Das Skript muss im Verzeichnis `auth_scripts` Ihres Arbeitsbereichs platziert werden:

- **Linux / macOS:** `{workspace}/auth_scripts/{domain}.sh`
- **Windows:** `{workspace}/auth_scripts/{domain}.bat`

## Skript-Ausgabeformat

Das Skript muss JSON an die Standardausgabe ausgeben, mit dem Token und seiner Ablaufzeit:

```bash
#!/bin/bash
# auth_scripts/my-auth.sh

TOKEN=$(curl -s -X POST https://auth.example.com/token \
  -d "grant_type=client_credentials" \
  -d "client_id=$CLIENT_ID" \
  -d "client_secret=$CLIENT_SECRET" | jq -r '.access_token')

echo "{\"token\": \"$TOKEN\", \"expires_in\": 3600}"
```

### JSON-Felder

| Feld | Erforderlich | Beschreibung |
|------|-------------|--------------|
| `token` | Ja | Authentifizierungstoken |
| `expires_in` | Nein | Token-Lebensdauer in Sekunden (Standard: 3600) |

## Hinweise

- swag2mcp führt das Skript bei jeder Anfrage aus, wenn das zwischengespeicherte Token abgelaufen ist
- Das Skript muss innerhalb von 30 Sekunden abgeschlossen sein
- Das Token wird bis zu seiner Ablaufzeit zwischengespeichert
- Skriptdateiname = `{domain}.sh` (Unix) oder `{domain}.bat` (Windows)
- `domain` darf keine `/` oder `\` enthalten
