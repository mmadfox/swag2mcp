# TUI-Explorer

## Übersicht

swag2mcp enthält eine integrierte TUI (Terminal User Interface) für die interaktive API-Erkundung. Es ist eine Vollbild-Terminalanwendung, mit der Sie API-Endpunkte durchsuchen, inspizieren und aufrufen können, ohne das Terminal zu verlassen.

## Start

```bash
swag2mcp run
```

Wenn keine Konfigurationsdatei existiert, startet die TUI automatisch zuerst den Initialisierungsassistenten.

## Modi

Die TUI hat drei Modi, die mit der `Tab`-Taste umgeschaltet werden können:

### Suchmodus

Volltextsuche über alle Endpunkte aller Specs. Unterstützt dieselbe Abfragesyntax wie das `search`-MCP-Tool.

- Geben Sie eine Abfrage ein, um nach Endpunktnamen, Pfaden und Beschreibungen zu suchen
- Filtern Sie Ergebnisse nach Methode, Tag oder Pfad
- Zeigen Sie Endpunktdetails mit einem Tastendruck an
- Navigieren Sie durch Ergebnisse mit Seitenumbrüchen (10 Elemente pro Seite)

### Durchsuchen-Modus

Baumnavigation durch die Spec-Hierarchie:

```
Spec → Collection → Tag → Endpunkt
```

- Navigieren Sie im Baum nach unten, um bestimmte Endpunkte zu finden
- Zeigen Sie Endpunktdetails an (Parameter, Anforderungstext, Antworten)
- Rufen Sie die API direkt aus der TUI auf
- Speichern Sie Endpunktdetails als JSON-Datei

### Auth-Modus

Zeigen Sie Authentifizierungstokens und Header für jede Spec an. Nützlich zum Debuggen oder Generieren von curl-Befehlen.

## Steuerung

| Taste | Aktion |
|-------|--------|
| `↑` / `↓` | Nach oben/unten navigieren |
| `Enter` | Auswählen oder öffnen |
| `Esc` | Eine Ebene zurück |
| `Tab` | Zwischen Such-, Durchsuchen- und Auth-Modus wechseln |
| `/` | Sucheingabe fokussieren |
| `N` / `P` | Nächste / vorherige Seite |
| `B` | Zurück zum vorherigen Bildschirm |
| `M` | Zurück zum Hauptmenü |
| `S` | Endpunktdetail als JSON-Datei speichern |
| `q` / `Ctrl+C` | Beenden |

## Zustände

Die TUI durchläuft diese Zustände während der Navigation:

1. **Laden** — Daten aus dem Arbeitsbereich laden
2. **Suche** — Suchmodus mit Abfrageeingabe
3. **Durchsuchen** — Durchsuchen-Modus mit Spec-Liste
4. **Spec-Liste** — Liste aller Specs
5. **Collection-Liste** — Collections innerhalb einer Spec
6. **Tag-Liste** — Tags innerhalb einer Collection
7. **Endpunkt-Liste** — Endpunkte innerhalb eines Tags
8. **Endpunkt-Detail** — vollständige Endpunktinformationen
9. **Aufruf-Ergebnis** — Ergebnis des API-Aufrufs
10. **Fehler** — Fehlerzustand mit Meldung

## Endpunkt-Detailansicht

Wenn Sie einen Endpunkt auswählen, zeigt die TUI:

- HTTP-Methode und Pfad
- Basis-URL und vollständige URL
- Zusammenfassung und Beschreibung
- Alle Parameter (Name, Ort, Typ, erforderlich)
- Anforderungstext-Schema (falls zutreffend)
- Antwortcodes und Schemata
- Veraltungsstatus

## Anforderungen

- **Terminalgröße:** Mindestens 80×24 Zeichen
- **Terminalemulator:** Funktioniert in den meisten modernen Terminals (iTerm2, Terminal.app, GNOME Terminal, Windows Terminal usw.)
- **SSH:** Funktioniert über SSH-Verbindungen

## Wichtige Hinweise

- **Auto-Init** — wenn keine Konfigurationsdatei existiert, startet die TUI automatisch den Initialisierungsassistenten
- **Seitenumbrüche** — Listen werden mit 10 Elementen pro Seite umgebrochen. Verwenden Sie `N` und `P` zum Navigieren
- **Endpunktdetails speichern** — drücken Sie `S` in der Endpunkt-Detailansicht, um das vollständige Detail als JSON-Datei im aktuellen Verzeichnis zu speichern
- **Auth-Modus** — zeigt Tokens und Header zum Debuggen an. In der Produktion kann das Auth-Tool mit `--disable-llm-auth` deaktiviert werden
