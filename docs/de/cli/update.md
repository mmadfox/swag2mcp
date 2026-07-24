# update

## Zweck

Validiert die Konfiguration erneut, leert den Cache und lädt alle Spezifikationsdateien neu herunter. Dies ist eine **vollständige Aktualisierung** des Arbeitsbereichs — sie stellt sicher, dass alle zwischengespeicherten Spezifikationen auf dem neuesten Stand sind und der Index neu aufgebaut wird.

## Wann verwenden

- Entfernte Spezifikationsdateien haben sich geändert und Sie möchten die neueste Version
- Nach dem Bearbeiten von `swag2mcp.yaml`, um Spezifikationsspeicherorte hinzuzufügen oder zu ändern
- Bei der Fehlerbehebung von veraltetem oder beschädigtem Cache
- Vor dem Ausführen von `mcp`, um sicherzustellen, dass alles aktuell ist

## Syntax

```bash
swag2mcp update [path]
```

## Argumente

| Argument | Position | Erforderlich | Beschreibung |
|----------|----------|-------------|--------------|
| `path` | 1 | Nein | Arbeitsbereichsverzeichnis. Wenn nicht angegeben, wird über die Pfadauflösungsregeln ermittelt. |

## Flags

Keine.

## Wie es funktioniert

Der Befehl `update` führt eine Pipeline von Operationen aus:

1. **Konfiguration laden** — liest `swag2mcp.yaml` aus dem Arbeitsbereich
2. **Validieren** — führt dieselben Prüfungen wie `validate` durch (YAML-Syntax, Struktur, Spezifikationsdatei-Erreichbarkeit, Format, Auth, HTTP-Client)
3. **Bereinigen** — entfernt alle Inhalte von `cache/` und `responses/`
4. **Neu cachen** — lädt alle entfernten Spezifikationsdateien herunter und kopiert lokale Spezifikationsdateien in den Cache
5. **Neu indizieren** — baut den Volltext-Suchindex für alle Endpunkte neu auf
6. **Auth-Skripte** — erstellt Stub-Auth-Skripte für Specs, die `ScriptAuth` verwenden
7. **Bereinigung verwaister Einträge** — entfernt Auth-Skripte für Specs, die nicht mehr existieren

```bash
swag2mcp update
swag2mcp update ./my-workspace
```

## Was mit deaktivierten Collections passiert

Collections mit `disable: true` werden vollständig übersprungen — sie werden nicht zwischengespeichert oder indiziert.

## Überprüfung nach dem Befehl

```bash
swag2mcp ls [path]
# Alle Specs sollten weiterhin aufgelistet und erreichbar sein
```

## Nuancen

- **Kein Auto-Init:** Wenn die Konfigurationsdatei nicht existiert, gibt `update` einen Fehler zurück: `"Konfiguration nicht gefunden unter &lt;path&gt;"`. Führen Sie zuerst `init` aus.
- **Netzwerkabhängigkeit:** Alle entfernten Spec-URLs müssen erreichbar sein. Wenn ein Download fehlschlägt, schlägt das gesamte Update mit einer klaren Fehlermeldung fehl.
- **Auth-Skript-Erstellung:** Wenn eine Spec `ScriptAuth` verwendet und das Stub-Skript nicht existiert, erstellt `update` es. Wenn die Erstellung fehlschlägt, schlägt das Update fehl.
- **`update` vs `clean`:** `clean` entfernt nur den Cache. `update` entfernt den Cache **und** lädt alles neu herunter. Verwenden Sie `clean`, wenn Sie nur Speicherplatz freigeben möchten; verwenden Sie `update`, wenn Sie eine Aktualisierung wünschen.
