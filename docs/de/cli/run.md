# run

## Zweck

Startet den interaktiven **TUI (Terminal User Interface)**-API-Explorer. Dies ist eine Vollbildanwendung zum Suchen, Durchsuchen, Inspizieren und Aufrufen von API-Endpunkten, ohne das Terminal zu verlassen.

## Wann verwenden

- Sie möchten Ihre APIs interaktiv erkunden
- Sie müssen nach einem bestimmten Endpunkt über alle Specs hinweg suchen
- Sie möchten die Spec → Collection → Tag → Endpunkt-Hierarchie durchsuchen
- Sie möchten einen API-Aufruf testen, bevor Sie den MCP-Server konfigurieren

## Syntax

```bash
swag2mcp run [path]
```

## Argumente

| Argument | Position | Erforderlich | Beschreibung |
|----------|----------|-------------|--------------|
| `path` | 1 | Nein | Arbeitsbereichsverzeichnis. Wenn nicht angegeben, wird über die Pfadauflösungsregeln ermittelt. |

## Flags

Keine.

## Modi

### Suchmodus

Volltextsuche über alle Endpunkte aller Specs. Unterstützt Filterung nach HTTP-Methode, Tag und Pfad.

- Geben Sie eine Abfrage ein, um nach Endpunktnamen, Pfaden und Beschreibungen zu suchen
- Filtern Sie Ergebnisse nach Methode (GET, POST, PUT, DELETE usw.)
- Zeigen Sie Endpunktdetails mit einem Tastendruck an

### Durchsuchen-Modus

Baumnavigation durch die Spec-Hierarchie:

```
Spec → Collection → Tag → Endpunkt
```

- Navigieren Sie im Baum nach unten, um bestimmte Endpunkte zu finden
- Zeigen Sie Endpunktdetails an (Parameter, Anforderungstext, Antworten)
- Rufen Sie die API direkt aus der TUI auf

## Navigation

| Taste | Aktion |
|-------|--------|
| `↑` / `↓` | Nach oben/unten navigieren |
| `Enter` | Auswählen oder öffnen |
| `Esc` | Zurückgehen |
| `Tab` | Zwischen Such- und Durchsuchen-Modus wechseln |
| `/` | Sucheingabe fokussieren |
| `q` | Beenden |

## Überprüfung nach dem Befehl

Die TUI lädt alle Specs aus dem Arbeitsbereich. Wenn eine Spec nicht geladen werden kann, wird eine Fehlermeldung in der Oberfläche angezeigt.

## Nuancen

- **Auto-Init:** Wenn keine Konfigurationsdatei existiert, führt `run` automatisch zuerst den Init-Assistenten aus.
- **Keine Flags:** Der Befehl `run` hat keine Flags — die gesamte Konfiguration stammt aus dem Arbeitsbereich.
- **Terminalgröße:** Die TUI erfordert ein Terminal mit mindestens 80×24 Zeichen. In sehr kleinen Terminals wird sie möglicherweise nicht korrekt dargestellt.
- **Abhängigkeiten:** Die TUI verwendet Bubbletea. Sie funktioniert über SSH und in den meisten Terminalemulatoren.
