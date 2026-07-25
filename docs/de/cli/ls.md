# ls

## Zweck

Listet alle konfigurierten **Specs** und ihre **Collections** in einem menschenlesbaren Format auf. Dies ist der primäre Weg, um zu überprüfen, welche APIs in Ihrem Arbeitsbereich verfügbar sind.

## Wann verwenden

- Sie möchten sehen, welche APIs konfiguriert sind
- Sie müssen eine Spec- oder Collection-ID finden
- Sie möchten überprüfen, wie viele Endpunkte jede Collection hat
- Sie möchten Specs nach Tags filtern

## Syntax

```bash
swag2mcp ls [path] [flags]
```

## Argumente

| Argument | Position | Erforderlich | Beschreibung |
|----------|----------|-------------|--------------|
| `path` | 1 | Nein | Arbeitsbereichsverzeichnis. Wenn nicht angegeben, wird über die Pfadauflösungsregeln ermittelt. |

## Flags

| Flag | Kurzform | Typ | Standard | Beschreibung |
|------|----------|-----|----------|--------------|
| `--tags` | `-t` | `string` | `""` | Specs nach Tags filtern (kommagetrennt) |

## Wie es funktioniert

### Alle Specs auflisten

Zeigt jede Spec mit ihrer Domain, Collections und Endpunktanzahl:

```bash
swag2mcp ls
```

Beispielausgabe:

```
Specifications:
  dadjoke (https://icanhazdadjoke.com)
    jokes (3 Endpunkte)
  meteo (https://meteo.swagger.io/v2)
    forecast (5 Endpunkte)
    current (8 Endpunkte)
  binance (https://api.binance.com)
    market-data (12 Endpunkte)
```

### Nach Tags filtern

Nur Specs anzeigen, die die angegebenen Tags haben:

```bash
swag2mcp ls --tags=public
swag2mcp ls --tags=public,internal
```

## Überprüfung nach dem Befehl

Verwenden Sie `ls` nach `add`, `delete`, `update` oder `import`, um zu bestätigen, dass der Arbeitsbereichszustand Ihren Erwartungen entspricht.

## Nuancen

- **Auto-Init:** Wenn keine Konfigurationsdatei existiert, führt `ls` automatisch zuerst den Init-Assistenten aus.
- **Tag-Filterung:** Tags sind kommagetrennt. Specs, die **irgendeinen** der angegebenen Tags haben, werden angezeigt (ODER-Logik).
- **Ausgabeformat:** Die Ausgabe ist reiner Text, kein JSON. Für maschinenlesbare Ausgabe verwenden Sie `info`.
