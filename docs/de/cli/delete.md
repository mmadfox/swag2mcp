# delete

## Zweck

Entfernt eine **Spec** (API-Dienst) oder **Collection** (Spezifikationsdatei) aus der Konfiguration. Dies ist die Umkehrung von `add`.

## Wann verwenden

- Eine API wird nicht mehr benötigt
- Sie möchten eine bestimmte Spezifikationsdatei aus einer Spec entfernen
- Sie räumen Ihren Arbeitsbereich auf

## Syntax

```bash
swag2mcp delete spec [path]
swag2mcp delete collection [path]
```

## Argumente

| Argument | Position | Erforderlich | Beschreibung |
|----------|----------|-------------|--------------|
| `path` | 1 | Nein | Arbeitsbereichsverzeichnis. Wenn nicht angegeben, wird über die Pfadauflösungsregeln ermittelt. |

## Flags

Keine. Beide Unterbefehle sind rein interaktiv.

## Wie es funktioniert

### Spec löschen

Fordert Sie auf, eine Spec aus einer Liste auszuwählen, und fragt dann nach Bestätigung vor dem Löschen.

```bash
swag2mcp delete spec
```

### Collection löschen

Fordert Sie auf, eine Spec, dann eine Collection innerhalb dieser Spec auszuwählen, und fragt dann nach Bestätigung.

```bash
swag2mcp delete collection
```

## IDs finden

Die interaktiven Eingabeaufforderungen zeigen menschenlesbare Namen, keine IDs. Wenn Sie IDs als Referenz benötigen:

```bash
# Alle Specs mit ihren IDs auflisten
swag2mcp ls

# Collections für eine bestimmte Spec auflisten
swag2mcp ls --tags
```

## Überprüfung nach dem Befehl

```bash
swag2mcp ls [path]
# Die gelöschte Spec oder Collection sollte nicht mehr erscheinen
```

## Nuancen

- **TTY erforderlich:** Beide Befehle erfordern ein interaktives Terminal. Sie funktionieren **nicht** in CI/CD-Pipelines, Cron-Jobs oder nicht-interaktiven Skripten.
- **Kein `--force` oder `--yes`:** Es gibt keine Möglichkeit, die Bestätigungsaufforderung zu überspringen. Dies ist beabsichtigt, um versehentliche Löschungen zu verhindern.
- **Auto-Init:** Wenn keine Konfigurationsdatei existiert, führt `delete` automatisch zuerst den Init-Assistenten aus.
- **Kein YAML-Modus:** Im Gegensatz zu `add` gibt es kein `--yaml`-Flag. Das Löschen erfolgt immer interaktiv.
