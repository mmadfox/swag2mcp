# CLI-Befehle

## Übersicht

Die `swag2mcp`-CLI ist der einzige Einstiegspunkt für alle Operationen — von der Initialisierung eines Arbeitsbereichs und der Verwaltung von API-Spezifikationen bis zum Starten eines MCP-Servers für die LLM-Integration. Sie bietet **13 Befehle**, die den gesamten Lebenszyklus der Arbeit mit OpenAPI/Swagger/Postman-Spezifikationen abdecken.

### Was die CLI löst

- **Arbeitsbereichs-Lebenszyklus** — erstellen (`init`), inspizieren (`info`, `ls`), bereinigen (`clean`), aktualisieren (`update`) und entfernen (`delete`) von Arbeitsbereichen und deren Inhalten
- **Spec- und Collection-Verwaltung** — hinzufügen (`add`), auflisten (`ls`) und löschen (`delete`) von API-Spezifikationen und deren Collections
- **Ausführungsmodi** — MCP-Server für LLM-Tool-Zugriff starten (`mcp`) oder den interaktiven TUI-Explorer starten (`run`)
- **Diagnose** — Konfiguration validieren (`validate`), Version anzeigen (`version`), Laufzeitinfo anzeigen (`info`)
- **Backup & Wiederherstellung** — vollständiger Arbeitsbereichs-Roundtrip über ZIP (`export`, `import`)

### Wichtige Nuancen

- **Pfadauflösung** — Befehle, die `[path]` akzeptieren, erwarten ein **Arbeitsbereichsverzeichnis** (keinen Dateipfad). Auflösungsreihenfolge: expliziter `[path]` → aktuelles Verzeichnis (`./`) → `~/.swag2mcp/`. Die CLI hängt `swag2mcp.yaml` automatisch an. Geben Sie immer einen expliziten Pfad an, wenn Sie es als Dienst oder in der IDE-Konfiguration ausführen, um das Laden des falschen Arbeitsbereichs zu vermeiden.
- **Spec vs Collection** — eine **Spec** repräsentiert einen logischen API-Dienst (z. B. "Open-Meteo API"), während eine **Collection** eine einzelne OpenAPI/Swagger/Postman-Datei ist. Eine Spec kann mehrere Collections haben.
- **`--version`** wird sowohl als Flag (`swag2mcp --version`) als auch als Unterbefehl (`swag2mcp version`) unterstützt.
- **`add spec` / `add collection`** akzeptieren YAML-Eingabe über `--yaml` (Inline-Zeichenfolge oder `-` für stdin). Das Weiterleiten aus einer Datei oder einem Here-Doc vermeidet Shell-Anführungsprobleme mit Sonderzeichen.
- **`delete`** erfordert ein TTY (interaktives Terminal). Es gibt kein `--force`- oder `--yes`-Flag — es fordert immer zur Auswahl und Bestätigung auf.
- **`mcp`** ist der primäre Befehl für die LLM-Integration. Er unterstützt drei Transports: `stdio` (Standard), `sse` und `streamable-http`. Das Flag `--disable-llm-auth` (Standard: `true`) entfernt das `auth`-Tool aus der MCP-Tool-Liste und verhindert, dass der LLM Tokens sieht oder anfordert. Auth funktioniert weiterhin — Tokens werden über den Standard-Konfigurationsmechanismus bezogen, nicht über den LLM. Dieser Modus wird für die **Produktion** empfohlen (LLM hat niemals Zugriff auf Anmeldeinformationen). Für das **Debuggen** oder bei kurzlebigen Tokens setzen Sie `--disable-llm-auth=false`, damit der LLM frische Tokens über das `auth`-Tool anfordern kann.
- **`validate`** prüft YAML-Syntax, Konfigurationsstruktur, Spezifikationsdateiexistenz, URL-Erreichbarkeit, Spezifikationsformat (OpenAPI/Swagger/Postman), Auth-Einstellungen und HTTP-Client-Korrektheit. Es testet **nicht** Authentifizierungsendpunkte oder API-Endpunktverfügbarkeit.
- **`export` / `import`** bieten einen vollständigen Arbeitsbereichs-Roundtrip — Konfigurationsdatei, Spezifikationsdateien, Cache und Auth-Skripte sind alle im ZIP-Archiv enthalten.
- **`clean`** entfernt die Verzeichnisse `cache/` und `responses/`, behält aber `specs/` und `auth_scripts/`. Alte Antworten (>48h) werden auch automatisch beim `mcp`-Start bereinigt.

## Befehle

| Befehl | Beschreibung |
|--------|--------------|
| [`init`](/cli/init) | Arbeitsbereichsverzeichnis mit Standardkonfiguration initialisieren |
| [`add`](/cli/add) | Eine Spec oder Collection zur Konfiguration hinzufügen |
| [`delete`](/cli/delete) | Eine Spec oder Collection interaktiv löschen |
| [`ls`](/cli/ls) | Alle Specs und ihre Collections auflisten |
| [`run`](/cli/run) | Interaktiven TUI-API-Explorer starten |
| [`validate`](/cli/validate) | Konfiguration und Spezifikationsdateien validieren |
| [`clean`](/cli/clean) | Zwischengespeicherte Specs und Aufrufantworten löschen |
| [`update`](/cli/update) | Alle Specs neu validieren, neu cachen und neu indizieren |
| [`mcp`](/cli/mcp) | MCP-Server für LLM-Tool-Zugriff starten |
| [`version`](/cli/version) | swag2mcp-Version anzeigen |
| [`info`](/cli/info) | Detaillierte Konfigurations- und Laufzeitinformationen anzeigen |
| [`import`](/cli/import) | Spezifikationsdateien importieren oder Arbeitsbereich aus ZIP wiederherstellen |
| [`export`](/cli/export) | Arbeitsbereich als portables ZIP-Backup exportieren |

## Globale Flags

| Flag | Beschreibung |
|------|--------------|
| `--version` | Version anzeigen (wie `version`-Unterbefehl) |
| `--help` | Hilfe für jeden Befehl anzeigen |
