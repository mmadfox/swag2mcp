# version

## Zweck

Zeigt die swag2mcp-Version an. Nützlich zum Überprüfen der installierten Version, zum Melden von Fehlern oder zum Prüfen der Kompatibilität.

## Wann verwenden

- Sie möchten überprüfen, welche Version von swag2mcp installiert ist
- Sie melden einen Fehler und müssen die Version angeben
- Sie möchten eine erfolgreiche Installation überprüfen

## Syntax

```bash
swag2mcp version
swag2mcp --version
```

## Argumente

Keine.

## Flags

Keine.

## Wie es funktioniert

```bash
swag2mcp version
# swag2mcp v1.2.0

swag2mcp --version
# swag2mcp v1.2.0
```

## Ausgabeformat

```
swag2mcp &lt;version&gt;
```

Die Version wird zur Build-Zeit über `ldflags` gesetzt. Wenn nicht gesetzt, lautet der Standardwert `"dev"`.

## Nuancen

- **Zwei Formen:** Sowohl `swag2mcp version` (Unterbefehl) als auch `swag2mcp --version` (globales Flag) erzeugen dieselbe Ausgabe.
- **Keine Konfiguration erforderlich:** Dieser Befehl funktioniert ohne Arbeitsbereich oder Konfigurationsdatei.
