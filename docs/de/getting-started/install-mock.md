# Mock-Server installieren

Die Binärdatei `swag2mcp-mock` ist ein separates Werkzeug zum Testen von APIs mit simulierten Antworten. Es unterstützt dieselben Installationsmethoden wie die Haupt-Binärdatei.

## Kompatibilität

| Methode | macOS | Linux | Windows |
|--------|-------|-------|---------|
| Einzeiler (curl) | ✅ | ✅ | ❌ |
| Homebrew | ✅ | ❌ | ❌ |
| Scoop | ❌ | ❌ | ✅ |
| GitHub Release | ✅ | ✅ | ✅ |
| go install | ✅ | ✅ | ✅ |
| Aus Quellen bauen | ✅ | ✅ | ✅ |
| Docker | ✅ | ✅ | ✅ |

---

## macOS

### Einzeiler (empfohlen)

Zwei Installationsmodi — wählen Sie den passenden:

```bash
# Mit sudo (installiert nach /usr/local/bin — empfohlen)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash

# Ohne sudo (installiert nach ~/.local/bin)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash -s -- --local
```

**`--sudo` (Standard):** Installiert nach `/usr/local/bin/swag2mcp-mock` mit `sudo`. Sie werden nach Ihrem Passwort gefragt. Falls `sudo` fehlschlägt, wird auf `~/.local/bin/swag2mcp-mock` zurückgegriffen.

**`--local`:** Installiert nach `~/.local/bin/swag2mcp-mock` ohne `sudo`. Fügen Sie nach der Installation Folgendes zu Ihrer Shell-Konfiguration hinzu:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Überprüfung:

```bash
swag2mcp-mock --version
```

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp-mock
```

Überprüfung:

```bash
swag2mcp-mock --version
```

### GitHub Release

1. Öffnen Sie die [Seite des neuesten Releases](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Laden Sie das Archiv für Ihren Mac herunter:
   - **Apple Silicon**: `swag2mcp-mock_darwin_arm64.tar.gz`
   - **Intel**: `swag2mcp-mock_darwin_amd64.tar.gz`
3. Öffnen Sie das Terminal im Download-Ordner und führen Sie aus:

```bash
tar -xzf swag2mcp-mock_darwin_*.tar.gz
sudo mv swag2mcp-mock /usr/local/bin/
```

Überprüfung:

```bash
swag2mcp-mock --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

Stellen Sie sicher, dass `$GOPATH/bin` in Ihrem `$PATH` ist:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Überprüfung:

```bash
swag2mcp-mock --version
```

### Aus Quellen bauen

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp-mock ./cmd/swag2mcp-mock
sudo mv swag2mcp-mock /usr/local/bin/
```

Überprüfung:

```bash
swag2mcp-mock --version
```

### Docker

Mounten Sie `~/.swag2mcp` (oder Ihren benutzerdefinierten Arbeitsverzeichnispfad) nach `/home/nonroot/.swag2mcp`.
Der Entrypoint passt die Dateiberechtigungen automatisch an, damit der Container Ihr Arbeitsverzeichnis lesen kann.

```bash
docker pull ghcr.io/mmadfox/swag2mcp-mock:latest
```

Mock-Server ausführen:

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest
```

Überprüfung:

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest --version
```

---

## Linux

### Einzeiler (empfohlen)

Zwei Installationsmodi — wählen Sie den passenden:

```bash
# Mit sudo (installiert nach /usr/local/bin — empfohlen)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash

# Ohne sudo (installiert nach ~/.local/bin)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash -s -- --local
```

**`--sudo` (Standard):** Installiert nach `/usr/local/bin/swag2mcp-mock` mit `sudo`. Sie werden nach Ihrem Passwort gefragt. Falls `sudo` fehlschlägt, wird auf `~/.local/bin/swag2mcp-mock` zurückgegriffen.

**`--local`:** Installiert nach `~/.local/bin/swag2mcp-mock` ohne `sudo`. Fügen Sie nach der Installation Folgendes zu Ihrer Shell-Konfiguration hinzu:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Überprüfung:

```bash
swag2mcp-mock --version
```

### GitHub Release

1. Öffnen Sie die [Seite des neuesten Releases](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Laden Sie das Archiv für Ihre Architektur herunter:
   - **amd64**: `swag2mcp-mock_linux_amd64.tar.gz`
   - **arm64**: `swag2mcp-mock_linux_arm64.tar.gz`
3. Öffnen Sie das Terminal im Download-Ordner und führen Sie aus:

```bash
tar -xzf swag2mcp-mock_linux_*.tar.gz
sudo mv swag2mcp-mock /usr/local/bin/
```

Überprüfung:

```bash
swag2mcp-mock --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

Stellen Sie sicher, dass `$GOPATH/bin` in Ihrem `$PATH` ist:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Überprüfung:

```bash
swag2mcp-mock --version
```

### Aus Quellen bauen

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp-mock ./cmd/swag2mcp-mock
sudo mv swag2mcp-mock /usr/local/bin/
```

Überprüfung:

```bash
swag2mcp-mock --version
```

### Docker

Mounten Sie `~/.swag2mcp` (oder Ihren benutzerdefinierten Arbeitsverzeichnispfad) nach `/home/nonroot/.swag2mcp`.
Der Entrypoint passt die Dateiberechtigungen automatisch an, damit der Container Ihr Arbeitsverzeichnis lesen kann.

```bash
docker pull ghcr.io/mmadfox/swag2mcp-mock:latest
```

Mock-Server ausführen:

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest
```

Überprüfung:

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest --version
```

---

## Windows

### Scoop

```powershell
scoop bucket add mmadfox https://github.com/mmadfox/scoop-bucket
scoop install mmadfox/swag2mcp-mock
```

Überprüfung:

```powershell
swag2mcp-mock --version
```

### GitHub Release

1. Öffnen Sie die [Seite des neuesten Releases](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Laden Sie `swag2mcp-mock_windows_amd64.zip` herunter
3. Entpacken Sie die ZIP-Datei (Rechtsklick → Alle extrahieren, oder mit PowerShell)
4. Verschieben Sie `swag2mcp-mock.exe` nach `C:\Windows\System32\`

Überprüfung:

```powershell
swag2mcp-mock --version
```

### go install

```powershell
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

Überprüfung:

```powershell
swag2mcp-mock --version
```

### Aus Quellen bauen

```powershell
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp-mock.exe ./cmd/swag2mcp-mock
```

Überprüfung:

```powershell
swag2mcp-mock --version
```

---

## Nächste Schritte

- [Schnellstart](quickstart.md) — in 2 Minuten starten
- [Mock-Server](../advanced/mock-server.md) — den Mock-Server konfigurieren und ausführen
