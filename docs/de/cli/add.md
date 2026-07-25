# add

## Zweck

Fügt eine neue **Spec** (API-Dienst) oder **Collection** (OpenAPI/Swagger/Postman-Datei) zu einer bestehenden Konfiguration hinzu. Dies ist der primäre Weg, um Ihren Arbeitsbereich mit neuen APIs zu erweitern.

## Wann verwenden

- Sie haben eine neue API, die Sie mit Ihrem LLM-Agenten verbinden möchten
- Sie haben eine OpenAPI-Spec-URL gefunden und möchten sie hinzufügen
- Sie möchten eine zusätzliche Spezifikationsdatei (Collection) zu einer bestehenden Spec hinzufügen
- Sie bevorzugen das direkte Schreiben von YAML anstelle des interaktiven Assistenten

## Syntax

```bash
swag2mcp add spec [path] [flags]
swag2mcp add collection [path] [flags]
```

## Argumente

| Argument | Position | Erforderlich | Beschreibung |
|----------|----------|-------------|--------------|
| `path` | 1 | Nein | Arbeitsbereichsverzeichnis. Wenn nicht angegeben, wird über die Pfadauflösungsregeln ermittelt. |

## Flags

### `add spec`

| Flag | Kurzform | Typ | Standard | Beschreibung |
|------|----------|-----|----------|--------------|
| `--yaml` | `-y` | `string` | `""` | YAML-Eingabe inline oder `-` für stdin |
| `--example` | `-e` | `bool` | `false` | YAML-Vorlage ausgeben und beenden |

### `add collection`

| Flag | Kurzform | Typ | Standard | Beschreibung |
|------|----------|-----|----------|--------------|
| `--yaml` | `-y` | `string` | `""` | YAML-Eingabe inline oder `-` für stdin |
| `--example` | `-e` | `bool` | `false` | YAML-Vorlage ausgeben und beenden |

## Wie es funktioniert

### Interaktiver Modus (Standard)

Startet einen TUI-Assistenten, der Sie Schritt für Schritt durch die Eingabe der Spec- oder Collection-Felder führt.

```bash
swag2mcp add spec
swag2mcp add collection
```

### YAML-Inline-Modus

Übergeben Sie das YAML direkt als Zeichenfolge. **Achten Sie auf Shell-Anführungszeichen** — Sonderzeichen wie `:`, `#`, `&`, `{` können den Befehl stören.

```bash
swag2mcp add spec --yaml 'domain: meteo
llm_title: Open-Meteo API
base_url: https://meteo.swagger.io/v2
collections:
  - llm_title: Main
    location: https://example.com/spec.json'
```

### YAML von stdin (für komplexes YAML empfohlen)

Leiten Sie aus einer Datei weiter oder verwenden Sie einen Here-Doc, um Shell-Anführungsprobleme vollständig zu vermeiden:

```bash
# Aus Datei weiterleiten
cat spec.yaml | swag2mcp add spec --yaml -

# Here-Doc
swag2mcp add spec --yaml - <<EOF
domain: my-api
llm_title: My API
llm_instruction: "Use this API for X & Y # important"
base_url: https://api.example.com/v1
collections:
  - llm_title: Main
    location: https://raw.githubusercontent.com/org/repo/main/spec.yaml
EOF
```

### YAML-Vorlage

Die erwartete YAML-Struktur ausgeben und beenden:

```bash
swag2mcp add spec --example
swag2mcp add collection --example
```

## YAML-Format

### Spec

```yaml
domain: meteo
llm_title: Open-Meteo API
llm_instruction: Use this API to manage pets.
base_url: https://meteo.swagger.io/v2
tags: [public, demo]
auth:
  type: bearer
  config:
    token: $(TOKEN)
collections:
  - llm_title: Open-Meteo Swagger
    location: https://example.com/spec.json
```

### Collection

```yaml
spec_domain: meteo
llm_title: Orders Collection
location: https://example.com/orders.json
```

## Überprüfung nach dem Befehl

```bash
swag2mcp ls [path]
# Die neue Spec oder Collection sollte in der Liste erscheinen
```

## Nuancen

- **Auto-Init:** Wenn keine Konfigurationsdatei existiert, führt `add` automatisch zuerst den Init-Assistenten aus. Sie müssen `init` nicht separat ausführen.
- **Shell-Anführungszeichen:** Inline-YAML (`--yaml '...'`) ist bei Sonderzeichen anfällig. Bevorzugen Sie `--yaml -` mit einem Here-Doc oder einer Pipe für alles, was über einfache Werte hinausgeht.
- **`--example` beendet sofort** ohne auf eine bestehende Konfiguration zu prüfen oder etwas zu ändern.
- **`add spec` vs `add collection`:** Verwenden Sie `add spec` für einen neuen API-Dienst (neue Domain). Verwenden Sie `add collection`, um eine weitere Spezifikationsdatei zu einer bestehenden Spec hinzuzufügen.
