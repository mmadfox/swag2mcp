# Konzepte

## Architektur

swag2mcp fungiert als Brücke zwischen API-Spezifikationen und LLM-Agenten:

<img src="/architecture.svg" width="800" alt="swag2mcp-Architektur">

## Kernkonzepte

**Spec** — ein logischer Container, der eine API-Domain oder einen API-Dienst repräsentiert (z. B. YouTube, Binance, Open-Meteo). Jede Spec hat eine eindeutige `domain`, eine `base_url`, optional `auth` und enthält eine oder mehrere Collections. Sie können auch `llm_instruction` setzen — einen kurzen Hinweis, der in den swag2mcp-System-Prompt eingefügt wird und dem LLM sagt, wofür diese Spec ist und wann sie verwendet werden soll. Mehr erfahren: [Specs](./specs).

**Collection** — eine einzelne OpenAPI/Swagger/Postman-Datei, die eine bestimmte API beschreibt. Sie verweist auf einen `location` (URL oder lokalen Dateipfad). Eine Spec kann mehrere Collections haben — zum Beispiel könnte die "meteo"-Spec die Collections "Forecast", "Air Quality" und "Marine" haben, die jeweils auf eine andere Spezifikationsdatei verweisen. Mehr erfahren: [Collections](./collections).

**Tag** — eine Kategorie von Endpunkten innerhalb einer Collection. Hilft dem LLM, die richtigen Operationen genauer zu finden. Mehr erfahren: [Tags](./tags).

**Endpunkt** — eine bestimmte HTTP-Methode + Pfad (z. B. `GET /api/users`). Der LLM kann einen Endpunkt anhand der Beschreibung finden, seine Parameter und Schemata inspizieren und ihn dann aufrufen. Mehr erfahren: [Endpunkte](./endpoints).

**Arbeitsbereich** — das Verzeichnis, in dem swag2mcp Konfiguration, Spec-Cache, gespeicherte Antworten und Auth-Skripte speichert. Mehr erfahren: [Arbeitsbereich](./workspace).

## Wie es funktioniert

1. **Eine Spec oder Collection hinzufügen** — in der YAML-Konfiguration definieren (`~/.swag2mcp/swag2mcp.yaml`). Zum Beispiel:

   ```yaml
   specs:
     - domain: jokes
       llm_title: Dad Joke API
       base_url: https://icanhazdadjoke.com
       collections:
         - llm_title: Jokes
           location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
   ```
2. **swag2mcp parst jede Collection** — erstellt Tags und Endpunkte, indiziert sie für die Suche.
3. **LLM findet den richtigen Endpunkt** — über MCP-Tools (`search`, `endpoint_by_tag`, `inspect`) sucht der LLM nach einem passenden Endpunkt anhand der Beschreibung, prüft seine Parameter und das Anforderungsschema.
4. **LLM ruft den Endpunkt auf** — über das MCP-Tool `invoke` sendet der LLM die Anfrage. swag2mcp validiert jeden Eingabeparameter gegen das OpenAPI-Schema des Endpunkts (Pfadparameter, Abfrageparameter, Header, Anforderungstext), bevor der Aufruf erfolgt. Wenn etwas nicht zum Schema passt, erhält der LLM eine klare Fehlermeldung, die erklärt, was falsch ist. Nach der Validierung führt swag2mcp den echten HTTP-Aufruf aus und gibt das Ergebnis zurück.
5. **Ergebnis geht zurück an den LLM** — die API-Antwort wird an den Agenten zurückgegeben. Große Antworten werden im Arbeitsbereich gespeichert und können mit drei speziellen MCP-Tools erkundet werden: `response_outline` (Struktur anzeigen), `response_compress` (auf eine repräsentative Stichprobe verkleinern) und `response_slice` (bestimmte Fragmente extrahieren).

swag2mcp ist eine Brücke zwischen LLMs und der Welt der APIs. Sie fügen API-Spezifikationen hinzu, und der LLM — über das MCP-Protokoll — findet die richtigen Endpunkte, inspiziert ihre Dokumentation und ruft sie auf. Alles, was Sie tun müssen, ist eine Spec hinzuzufügen und den MCP-Server zu starten.

> **Die Konfiguration kann jederzeit bearbeitet werden.** Die YAML-Konfigurationsdatei (`~/.swag2mcp/swag2mcp.yaml`) kann von Hand bearbeitet werden — Specs hinzufügen, Auth ändern, Einstellungen anpassen. Nach jeder Bearbeitung starten Sie den MCP-Server (`swag2mcp mcp`) neu, damit die Änderungen wirksam werden.

## Hierarchie

```
Spec (domain, z. B. "meteo")
  └── Collection 1 (Spezifikationsdatei, z. B. forecast.yml)
        └── Tag 1 (Kategorie)
              └── Endpunkt (GET /api/forecast)
              └── Endpunkt (POST /api/forecast)
        └── Tag 2
              └── Endpunkt (GET /api/forecast/{id})
  └── Collection 2 (Spezifikationsdatei, z. B. air-quality.yml)
        └── Tag 3
              └── Endpunkt (GET /api/air-quality)
```
