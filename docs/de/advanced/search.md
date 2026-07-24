# Volltextsuche

## Übersicht

swag2mcp enthält eine integrierte Volltext-Suchmaschine (bluge), die alle Endpunkte aller Specs indiziert. Der LLM kann nach Endpunkten nach Methode, Pfad, Zusammenfassung oder Tag suchen — sogar ohne die Endpunkt-ID zu kennen.

## Wie die Indizierung funktioniert

Wenn eine Spec hinzugefügt oder aktualisiert wird, wird jeder Endpunkt indiziert. Die folgenden Felder sind durchsuchbar:

| Feld | Beschreibung | Beispiel |
|------|--------------|----------|
| `method` | HTTP-Methode | `GET`, `POST`, `PUT` |
| `path` | API-Endpunkt-Pfad | `/api/v1/users/{id}` |
| `summary` | OpenAPI-Zusammenfassung | "Find pet by ID" |
| `tag` | Endpunkt-Kategorie | "pets", "users" |
| `_all` | Alle Felder kombiniert | method + path + tag + summary |

Der Index wird bei jedem MCP-Server-Neustart neu aufgebaut. Er wird für schnelle Suchvorgänge im Speicher gehalten.

## Abfragesyntax

Die Suche unterstützt eine umfangreiche Abfragesyntax für präzises Filtern:

| Beispiel | Beschreibung |
|----------|--------------|
| `pet` | Einfache Textsuche über alle Felder |
| `method:GET` | Alle GET-Endpunkte finden |
| `tag:pets` | Endpunkte im Tag "pets" finden |
| `path:"/api/v1/users"` | Exakte Pfadübereinstimmung |
| `+method:POST +tag:pet` | Beide Bedingungen müssen zutreffen |
| `-method:DELETE` | DELETE-Methoden ausschließen |
| `create~` | Unscharfe Suche (toleriert Tippfehler) |
| `cr*` | Platzhaltersuche |
| `"find pet"` | Phrasensuche |
| `+summary:pet -method:DELETE` | "pet" in Zusammenfassung einschließen, DELETE ausschließen |

### Feldspezifische Suche

Sie können innerhalb bestimmter Felder mit der Syntax `field:value` suchen:

```
method:GET
tag:pets
path:"/pet/findByStatus"
summary:"find pet by status"
```

### Boolesche Operatoren

- `+` — der Begriff muss übereinstimmen (AND)
- `-` — der Begriff darf nicht übereinstimmen (NOT)
- Leerzeichen zwischen Begriffen — OR (jeder Begriff kann übereinstimmen)

### Unscharfe Suche und Platzhalter

- `term~` — unscharfe Suche (findet ähnliche Wörter, behandelt Tippfehler)
- `te*` — Platzhalter (passt auf beliebige Zeichen)
- `te?t` — Einzelzeichen-Platzhalter

## Beispiele

```
# Alle GET-Anfragen finden
method:GET

# POST-Anfragen im Tag "pet" finden
+method:POST +tag:pet

# Endpunkte nach exaktem Pfad finden
path:"/pet/findByStatus"

# Nach Beschreibung finden
"find pet by status"

# Alles außer DELETE finden
+summary:pet -method:DELETE

# Unscharfe Suche nach "create" (behandelt Tippfehler)
create~
```

## MCP-Tool

Das `search`-MCP-Tool stellt die Suchmaschine dem LLM zur Verfügung:

```
→ search(query: "find pet by status", limit: 5)
← GET /pet/findByStatus — Findet Haustiere nach Status
   GET /pet/{petId} — Haustier nach ID finden
```

### Parameter

| Parameter | Erforderlich | Beschreibung |
|-----------|-------------|--------------|
| `query` | Ja | Suchabfrage (unterstützt strukturierte Syntax) |
| `limit` | Ja | Maximale Ergebnisse (1-50) |

## Wichtige Hinweise

- **Der Index ist im Speicher** — er wird bei jedem Start des MCP-Servers neu aufgebaut. Es gibt keine persistente Indexdatei.
- **Alle Felder werden kleingeschrieben** — die Suche ist nicht case-sensitive
- **Limit ist auf 50 begrenzt** — Sie können nicht mehr als 50 Ergebnisse anfordern
- **Ungültige Abfragesyntax** gibt eine hilfreiche Fehlermeldung mit Beispielen zurück
- **Das Feld `_all`** kombiniert Methode, Pfad, Tag und Zusammenfassung für einfache Textsuche
