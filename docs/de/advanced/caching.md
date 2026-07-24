# Caching

## Übersicht

swag2mcp speichert heruntergeladene Spezifikationsdateien zwischen, damit der MCP-Server bei nachfolgenden Starts schneller startet. Anstatt jedes Mal dieselbe Spezifikationsdatei herunterzuladen, wird die zwischengespeicherte Kopie wiederverwendet.

## Wie Caching funktioniert

Wenn Sie eine Spec mit einer entfernten URL hinzufügen, lädt swag2mcp sie herunter und speichert sie im Verzeichnis `cache/`. Beim nächsten Start wird geprüft, ob die zwischengespeicherte Kopie noch aktuell ist. Wenn ja, wird der Download übersprungen.

### Was wird zwischengespeichert

| Quelle | Verhalten |
|--------|-----------|
| **Entfernte URL** (http/https) | Immer zwischengespeichert. Einmal heruntergeladen, bis zum Ablauf des Caches wiederverwendet. |
| **Lokale Datei in `specs/`** | Direkt aus dem Verzeichnis `specs/` verwendet. Nie zwischengespeichert — Änderungen sind sofort sichtbar. |
| **Lokale Datei außerhalb `specs/`** | In den Cache kopiert. Wenn sich die Quelldatei ändert (Änderungszeit), wird der Cache ungültig. |

### Cache-Ablauf (TTL)

Jede zwischengespeicherte Datei erhält eine zufällige Ablaufzeit zwischen **1 Stunde und 48 Stunden**. Die Zufälligkeit verhindert, dass alle zwischengespeicherten Dateien gleichzeitig ablaufen (was zu einer Überlastung durch gleichzeitige Downloads führen würde).

- Die TTL wird bei jedem Start des MCP-Servers zurückgesetzt
- Wenn eine zwischengespeicherte Datei noch innerhalb ihrer TTL liegt, wird sie wiederverwendet
- Wenn die TTL abgelaufen ist, wird die Datei erneut heruntergeladen

### Cache-Struktur

```
~/.swag2mcp/cache/
├── a1b2c3d4e5f6a7b8.spec    # Zwischengespeicherte Spezifikationsdatei
├── a1b2c3d4e5f6a7b8.meta    # Metadaten (Quelle, TTL, Cache-Zeitpunkt)
├── b2c3d4e5f6a7b8c9.spec
├── b2c3d4e5f6a7b8c9.meta
└── ...
```

Der Cache-Schlüssel wird aus der URL oder dem Pfad der Spezifikationsdatei abgeleitet. Jede zwischengespeicherte Datei hat eine begleitende `.meta`-Datei, die speichert, wann sie zwischengespeichert wurde und wann sie abläuft.

## Cache verwalten

### Aktualisierung erzwingen

Führen Sie `swag2mcp update` aus, um den gesamten Cache zu leeren und alle Spezifikationsdateien neu herunterzuladen:

```bash
swag2mcp update
```

Dies validiert die Konfiguration, leert den Cache und lädt alles frisch herunter.

### Cache manuell leeren

```bash
swag2mcp clean
```

Dies entfernt alle zwischengespeicherten Spezifikationsdateien und gespeicherten API-Antworten. Beim nächsten Start des MCP-Servers werden alle Spezifikationen erneut heruntergeladen.

### Automatische Bereinigung

Beim Start des MCP-Servers (`swag2mcp mcp`) werden gespeicherte API-Antworten, die älter als 48 Stunden sind, automatisch entfernt. Dies verhindert, dass das Verzeichnis `responses/` unbegrenzt wächst.

## Wichtige Hinweise

- **Lokale Dateien in `specs/` werden nie zwischengespeichert** — wenn Sie eine Spezifikationsdatei direkt im Verzeichnis `specs/` bearbeiten, sind die Änderungen sofort sichtbar, ohne den Cache zu leeren
- **Entfernte URLs werden immer zwischengespeichert** — es gibt keine Möglichkeit, den Cache für entfernte URLs zu umgehen, außer durch Ausführen von `swag2mcp update` oder `swag2mcp clean`
- **Der Cache ist lokal** — er wird auf der Festplatte gespeichert und nicht zwischen Rechnern synchronisiert. Verwenden Sie `swag2mcp export` und `swag2mcp import`, um Spezifikationen zwischen Rechnern zu übertragen
