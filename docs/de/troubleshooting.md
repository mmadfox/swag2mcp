# Fehlerbehebung

## Installationsprobleme

### swag2mcp: Befehl nicht gefunden

Die Binärdatei befindet sich nicht in Ihrem PATH.

```bash
# Prüfen, ob Go installiert ist
go version

# Herausfinden, wo Go Binärdateien installiert
go env GOPATH
# Normalerweise ~/go oder ~/go/bin

# Zum PATH hinzufügen (in ~/.zshrc oder ~/.bashrc einfügen)
export PATH=$PATH:$(go env GOPATH)/bin

# Oder den vollständigen Pfad verwenden
~/go/bin/swag2mcp --version
```

Wenn Sie eine Binärdatei von GitHub Releases heruntergeladen haben, stellen Sie sicher, dass sie sich in einem Verzeichnis befindet, das in Ihrem PATH ist:

```bash
# Nach /usr/local/bin verschieben (macOS/Linux)
sudo mv swag2mcp /usr/local/bin/
```

### Keine Ausführungsberechtigung

Die Binärdatei hat keine Ausführungsberechtigung.

```bash
# Für go install (Besitzer korrigieren)
sudo chown -R $(whoami) $(go env GOPATH)

# Für heruntergeladene Binärdatei
chmod +x /pfad/zu/swag2mcp
```

### Go-Version zu alt

swag2mcp benötigt Go 1.23+.

```bash
go version
# Wenn Version < 1.23, Go aktualisieren:
# https://go.dev/dl/
```

### Mock-Server nicht gefunden

Der Mock-Server ist eine separate Binärdatei. Installieren Sie ihn explizit:

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

## Konfigurationsprobleme

### Konfigurationsdatei nicht gefunden

swag2mcp kann `swag2mcp.yaml` nicht finden.

```bash
# Neue Konfiguration erstellen
swag2mcp init

# Oder den Pfad explizit angeben
swag2mcp mcp /pfad/zu/arbeitsbereich
swag2mcp ls /pfad/zu/arbeitsbereich
```

**Häufige Ursache:** Sie haben `swag2mcp mcp` aus einem beliebigen Verzeichnis ausgeführt, und es hat nach `~/.swag2mcp/` statt nach dem Arbeitsbereich Ihres Projekts gesucht. Geben Sie den Pfad immer explizit an.

### Falscher Arbeitsbereich geladen

swag2mcp hat einen anderen Arbeitsbereich als erwartet geladen.

**Auflösungsreihenfolge:** Expliziter `[path]` → aktuelles Verzeichnis (`./`) → `~/.swag2mcp/`. Wenn Sie `swag2mcp mcp` ohne Pfad aus einem Verzeichnis ausführen, das keine `swag2mcp.yaml` hat, wird auf `~/.swag2mcp/` zurückgegriffen.

**Lösung:** Geben Sie immer den Arbeitsbereichspfad an: `swag2mcp mcp /pfad/zu/ihrem/arbeitsbereich`

### YAML-Parsing-Fehler

Die Konfigurationsdatei hat eine ungültige YAML-Syntax.

```bash
# Konfiguration validieren
swag2mcp validate

# Häufige Fehler:
# - Tabulatoren statt Leerzeichen (YAML benötigt Leerzeichen)
# - Fehlende Einrückung für verschachtelte Felder
# - Nicht in Anführungszeichen gesetzte Zeichenfolgen mit Sonderzeichen (: # & {)
```

**Tipp:** Verwenden Sie einen YAML-Linter oder einen Editor mit YAML-Unterstützung, um Syntaxfehler zu erkennen.

### Validierung fehlgeschlagen: "keine Spezifikationen definiert"

Die Konfigurationsdatei existiert, hat aber keine Specs.

```bash
# Eine Spec hinzufügen
swag2mcp add spec

# Oder swag2mcp.yaml bearbeiten und mindestens eine Spec hinzufügen
```

### Validierung fehlgeschlagen: "doppelte Domain"

Zwei Specs haben denselben `domain`-Wert. Domains müssen eindeutig sein.

```bash
# Aktuelle Specs auflisten
swag2mcp ls

# Auf doppelte Domains in swag2mcp.yaml prüfen
```

### Validierung fehlgeschlagen: "ungültiger Spec-Speicherort"

Die `location`-URL oder der Dateipfad ist nicht erreichbar oder keine gültige Spezifikationsdatei.

```bash
# Prüfen, ob die URL erreichbar ist
curl -I https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml

# Prüfen, ob die lokale Datei existiert
ls -la ./specs/my-api.yaml

# Überprüfen, ob die Datei ein gültiges OpenAPI/Swagger/Postman-Format hat
# (nicht irgendeine JSON- oder HTML-Seite)
```

**Häufige Ursache:** Das `location`-Feld zeigt auf den API-Endpunkt selbst (z. B. `https://api.example.com/v1/users`) anstatt auf die URL der Spezifikationsdatei. Der Speicherort muss auf eine OpenAPI/Swagger/Postman-Datei verweisen.

## MCP-Server-Probleme

### Port bereits belegt

Ein anderer Prozess verwendet den Port.

```bash
# Prozess finden
lsof -i :8080

# Beenden
kill <PID>

# Oder einen anderen Port verwenden
swag2mcp mcp --transport sse --http-addr :9090
```

### Verbindung abgelehnt

Der MCP-Server läuft nicht oder ist nicht erreichbar.

```bash
# Sicherstellen, dass der Server läuft
swag2mcp mcp --transport sse --http-addr 127.0.0.1:8080

# In einem anderen Terminal den Health-Endpunkt prüfen
curl http://127.0.0.1:8080/health

# Bei Verwendung eines benutzerdefinierten Pfads
curl http://127.0.0.1:8080/benutzerdefinierter-pfad/health
```

### MCP-Tools werden im LLM-Client nicht angezeigt

Der LLM-Client kann keine Tools sehen.

```bash
# Prüfen, ob Specs geladen sind
swag2mcp ls

# Prüfen, ob Specs nicht deaktiviert sind
swag2mcp validate

# Server-Logs prüfen
swag2mcp mcp --logfile /tmp/swag2mcp.log
cat /tmp/swag2mcp.log

# Überprüfen, ob der Arbeitsbereichspfad in Ihrer IDE-Konfiguration korrekt ist
# (muss ein absoluter Pfad sein)
```

**Häufige Ursachen:**
- Falscher Arbeitsbereichspfad in der IDE-Konfiguration
- Alle Specs haben `disable: true`
- Specs werden durch `--tags` herausgefiltert
- Konfigurationsdatei existiert nicht am angegebenen Pfad

### MCP-Handshake fehlgeschlagen (HTTP-Transport)

Für SSE- und Streamable-HTTP-Transports erfordert das MCP-Protokoll eine Initialisierung, bevor Tool-Aufrufe funktionieren.

```
Schritt 1: POST /mcp → {"method":"initialize", ...}
Schritt 2: POST /mcp → {"method":"notifications/initialized"}
Schritt 3: POST /mcp → {"method":"tools/list", ...}  ← jetzt funktioniert es
```

Stellen Sie sicher, dass Ihr LLM-Client den Handshake abschließt, bevor er Tools aufruft.

### Health-Check gibt 404 zurück

Der Health-Endpoint-Pfad kann sich vom MCP-Pfad unterscheiden.

```bash
# Standard-Health-Endpoint
curl http://127.0.0.1:8080/health

# Wenn Sie den MCP-Pfad geändert haben, ist Health immer noch unter /health
# (nicht von --http-path betroffen)
```

### Auth-Tool nicht verfügbar

Das `auth`-MCP-Tool wird nicht angezeigt.

Das `auth`-Tool ist **standardmäßig deaktiviert** (`--disable-llm-auth=true`). Dies ist beabsichtigt für die Produktionssicherheit.

```bash
# Auth-Tool aktivieren
swag2mcp mcp --disable-llm-auth=false
```

## Authentifizierungsprobleme

### 401 Nicht autorisiert

Die API hat die Anfrage aufgrund fehlender oder ungültiger Anmeldeinformationen abgelehnt.

```bash
# Prüfen, ob Authentifizierung konfiguriert ist
swag2mcp info

# Konfiguration validieren
swag2mcp validate

# Prüfen, ob Umgebungsvariablen gesetzt sind
echo $MY_TOKEN

# Überprüfen, ob das Token nicht abgelaufen ist (Bearer-Tokens sind statisch)
```

**Häufige Ursachen:**
- Token fehlt oder ist leer
- Umgebungsvariable nicht gesetzt
- Token ist abgelaufen (Bearer-Tokens erneuern sich nicht automatisch)
- Falscher Authentifizierungstyp konfiguriert

### 403 Verboten

Die API hat die Anfrage aufgrund unzureichender Berechtigungen abgelehnt.

- Das Token hat möglicherweise nicht die erforderlichen Bereiche
- Der API-Schlüssel hat möglicherweise keinen Zugriff auf diese Ressource
- Überprüfen Sie die API-Dokumentation auf erforderliche Berechtigungen

### OAuth2-Token-Endpunkt nicht erreichbar

swag2mcp kann die OAuth2-Token-URL nicht erreichen.

```bash
# token_url in Ihrer Konfiguration prüfen
# Überprüfen, ob die URL korrekt und erreichbar ist
curl -X POST https://auth.example.com/oauth/token \
  -d "grant_type=client_credentials" \
  -d "client_id=test" \
  -d "client_secret=test"

# Netzwerkkonnektivität prüfen
# Proxy-Einstellungen prüfen, wenn hinter einem Firmenproxy
```

### Digest-Authentifizierung schlägt fehl

swag2mcp kann den Digest-Authentifizierungs-Handshake nicht abschließen.

- Der Server muss einen `WWW-Authenticate: Digest ...`-Header mit einer 401-Antwort zurückgeben
- Die Challenge wird 5 Minuten lang zwischengespeichert — wenn der Server sein Nonce ändert, warten Sie, bis der Cache abläuft
- Prüfen Sie, ob Benutzername und Passwort korrekt sind

### HMAC-Signatur stimmt nicht überein

Die API hat die HMAC-signierte Anfrage abgelehnt.

- Überprüfen Sie, ob `api_key` und `secret_key` korrekt sind
- Prüfen Sie, ob die API Binance-kompatibles HMAC-SHA256-Signing verwendet
- Einige Börsen verwenden andere Signiermethoden — HMAC-Auth ist speziell für Binance-kompatible APIs

### Skript-Authentifizierung schlägt fehl

Das externe Authentifizierungsskript ist fehlgeschlagen.

```bash
# Prüfen, ob das Skript existiert
ls -la ~/.swag2mcp/auth_scripts/my-domain.sh

# Skript manuell testen
sh ~/.swag2mcp/auth_scripts/my-domain.sh

# Ausgabeformat des Skripts prüfen (muss JSON sein: {"token": "...", "expires_in": 3600})
# Prüfen, ob das Skript innerhalb von 30 Sekunden abschließt
# Prüfen, ob das Skript Ausführungsberechtigung hat
chmod +x ~/.swag2mcp/auth_scripts/my-domain.sh
```

## Suchprobleme

### Keine Suchergebnisse

Die Suche hat keine Endpunkte zurückgegeben.

```bash
# Prüfen, ob Specs geladen sind
swag2mcp ls

# Prüfen, ob Specs nicht deaktiviert sind
swag2mcp validate

# Einfachere Abfrage versuchen
# Nach Methode suchen: method:GET
# Nach Tag suchen: tag:pets

# Der Index wird bei jedem MCP-Server-Neustart neu aufgebaut
# Wenn Sie gerade eine Spec hinzugefügt haben, starten Sie den Server neu
```

### Suche liefert irrelevante Ergebnisse

Die Abfrage ist zu breit oder mehrdeutig.

- Verwenden Sie Feldfilter zur Eingrenzung: `method:GET +tag:pets`
- Verwenden Sie genaue Phrasen: `"find pet by status"`
- Verwenden Sie den `limit`-Parameter für fokussiertere Ergebnisse

## API-Aufruf-Probleme

### invoke gibt einen Fehler zurück

Der API-Aufruf ist fehlgeschlagen.

```bash
# Fehlermeldung prüfen — sie enthält den HTTP-Statuscode
# 4xx-Fehler: Parameter, Authentifizierung oder Berechtigungen prüfen
# 5xx-Fehler: Der API-Server hat ein Problem

# Endpunkt vor dem Aufruf immer inspizieren
inspect(endpointId: "...")

# Prüfen, ob alle erforderlichen Parameter angegeben sind
# Parametertypen prüfen (Zeichenfolge, Zahl, boolesch)
```

### Ratenbegrenzungsfehler

Der LLM hat denselben Endpunkt zu schnell aufgerufen.

Jeder Endpunkt hat eine 10-Sekunden-Abklingzeit. Warten Sie vor dem erneuten Aufruf oder deaktivieren Sie den Ratenbegrenzer:

```yaml
disable_ratelimiter: true
```

### Antwort zu groß (fileRef zurückgegeben)

Die Antwort hat `max_response_size` überschritten.

Dies ist normal. Verwenden Sie die Antwort-Tools, um die Daten zu erkunden:

```
1. response_outline(path) → Struktur verstehen
2. response_compress(path, mode: "first_of_array") → Beispiel abrufen
3. response_slice(path, jsonPath: "data.0") → Bestimmte Daten abrufen
```

Oder erhöhen Sie das Limit:

```yaml
http_client:
  max_response_size: 4194304  # 4 MB
```

### Langsame API-Antworten

Die API braucht zu lange zum Antworten.

```yaml
http_client:
  timeout: 120s  # Von Standard 30s erhöhen
```

## Arbeitsbereichsprobleme

### swag2mcp init fehlgeschlagen: "Verzeichnis ist nicht leer"

Das Zielverzeichnis enthält bereits Dateien.

```bash
# --force zum Überschreiben verwenden
swag2mcp init --force

# Oder ein anderes Verzeichnis verwenden
swag2mcp init ./neuer-arbeitsbereich
```

### swag2mcp update fehlgeschlagen

Eine oder mehrere Spezifikationsdateien konnten nicht heruntergeladen werden.

```bash
# Fehlermeldung prüfen, welche URL fehlgeschlagen ist
# Überprüfen, ob die URL erreichbar ist
curl -I <fehlgeschlagene-url>

# Netzwerkkonnektivität prüfen
# Proxy-Einstellungen prüfen
```

### Export erstellt kein ZIP

Das Argument `[output]` muss ein Dateipfad sein, der auf `.zip` endet, kein Verzeichnis.

```bash
# Richtig
swag2mcp export /pfad/zu/arbeitsbereich /pfad/zu/sicherung.zip

# Falsch (es wird kein ZIP erstellt)
swag2mcp export /pfad/zu/arbeitsbereich /ein/verzeichnis
```

### Import fehlgeschlagen: "keine gültige swag2mcp-Sicherung"

Die ZIP-Datei wurde nicht von `swag2mcp export` erstellt.

Nur ZIP-Archive, die von `swag2mcp export` erstellt wurden, können importiert werden. Das Archiv hat eine spezifische interne Struktur (`swag2mcp.yaml`, `specs/`, `auth_scripts/`).

## TUI-Probleme

### TUI wird nicht korrekt dargestellt

Das Terminal ist zu klein oder unterstützt die erforderlichen Funktionen nicht.

- Minimale Terminalgröße: 80×24 Zeichen
- Die TUI verwendet Bubbletea und funktioniert in den meisten modernen Terminals
- Versuchen Sie, Ihr Terminalfenster zu vergrößern
- Versuchen Sie einen anderen Terminalemulator

### TUI zeigt "keine Specs gefunden"

Der Arbeitsbereich hat keine konfigurierten Specs.

```bash
# Specs prüfen
swag2mcp ls

# Eine Spec hinzufügen
swag2mcp add spec
```

## Mock-Server-Probleme

### Mock-Server startet nicht

```bash
# Prüfen, ob mock_enabled: true in der Konfiguration
# Prüfen, ob jede Collection base_mock_url gesetzt hat
# Prüfen, ob Ports nicht belegt sind
lsof -i :9090

# Mock-Server-Logs prüfen
swag2mcp-mock mockserver
```

### Mock-Server gibt leere Antworten zurück

Die Spezifikationsdatei hat möglicherweise keine Antwortschemata definiert.

- Der Mock-Server generiert Daten aus Antwortschemata
- Wenn kein Schema gefunden wird, gibt er `{}` zurück
- Prüfen Sie, ob Ihre OpenAPI-Spezifikation `responses` mit `schema` definiert hat

## Netzwerkprobleme

### Proxy-Verbindung fehlgeschlagen

swag2mcp kann keine Verbindung über den konfigurierten Proxy herstellen.

```bash
# Proxy-URL-Format prüfen (muss Schema enthalten: http://, https://, socks5://)
# Proxy-Anmeldeinformationen prüfen
# Bypass-Liste prüfen — das Ziel könnte in der Bypass-Liste sein
# Proxy mit curl testen
curl -x http://proxy.company.com:8080 https://api.example.com
```

### TLS/SSL-Fehler

Die Zertifikatsprüfung ist fehlgeschlagen.

- Bei Verwendung eines selbstsignierten Zertifikats für den MCP-Server muss der Client ihm vertrauen
- Für den Mock-Server mit `--tls` wird automatisch ein selbstsigniertes Zertifikat generiert
- Für API-Aufrufe verwendet swag2mcp den Systemzertifikatsspeicher

## Sonstige Probleme

### Hohe Festplattennutzung

Die Cache- und Antwortverzeichnisse können mit der Zeit wachsen.

```bash
# Alles bereinigen
swag2mcp clean

# Alte Antworten (>48h) werden automatisch beim MCP-Server-Start bereinigt
# Cache-Dateien laufen zufällig zwischen 1-48 Stunden ab
```

### "Befehl nicht gefunden" nach go install

Das `go install`-Verzeichnis ist nicht in Ihrem PATH.

```bash
# Herausfinden, wo Go Binärdateien installiert
go env GOPATH
# Zum PATH hinzufügen
export PATH=$PATH:$(go env GOPATH)/bin
```

### LLM verwendet die Tools nicht korrekt

Der LLM benötigt möglicherweise bessere Anweisungen oder eine Formatierungs-Skill.

- Verwenden Sie `llm_instruction` in Ihrer Spec-Konfiguration, um zu beschreiben, was die API tut
- Erwägen Sie die Verwendung der [swag2mcp-format-Skill](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md) für eine konsistente Ausgabeformatierung
- Die Qualität der LLM-Antworten hängt vom Modell und den Anweisungen ab, die es erhält

### Wie melde ich einen Fehler?

Eröffnen Sie ein Issue auf [GitHub](https://github.com/mmadfox/swag2mcp/issues) mit:
- swag2mcp-Version (`swag2mcp --version`)
- Ihrem Betriebssystem und Ihrer Architektur
- Dem genauen Befehl, den Sie ausgeführt haben
- Der vollständigen Fehlermeldung
- Ihrer Konfigurationsdatei (ohne Geheimnisse)
