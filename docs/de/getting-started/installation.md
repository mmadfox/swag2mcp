# Installation

## Anforderungen

- **macOS, Linux oder Windows** (amd64 / arm64)
- **Go 1.26+** (nur für `go install` oder Bauen aus dem Quellcode)

## Kompatibilität

| Methode | macOS | Linux | Windows |
|---------|-------|-------|---------|
| Einzeiler (curl) | ✅ | ✅ | ❌ |
| Homebrew | ✅ | ❌ | ❌ |
| Scoop | ❌ | ❌ | ✅ |
| APT (deb) | ❌ | ✅ | ❌ |
| RPM | ❌ | ✅ | ❌ |
| Docker | ✅ | ✅ | ✅ |
| GitHub Release | ✅ | ✅ | ✅ |
| go install | ✅ | ✅ | ✅ |
| Aus Quellcode bauen | ✅ | ✅ | ✅ |

---

## macOS

### Einzeiler (empfohlen)

Zwei Installationsmodi — wählen Sie den passenden:

```bash
# Mit sudo (installiert nach /usr/local/bin — empfohlen)
curl -fsSL https://swag2mcp.io/install.sh | bash

# Ohne sudo (installiert nach ~/.local/bin)
curl -fsSL https://swag2mcp.io/install.sh | bash -s -- --local
```

**`--sudo` (Standard):** Installiert nach `/usr/local/bin/swag2mcp` mit `sudo`. Sie werden nach Ihrem Passwort gefragt. Wenn `sudo` fehlschlägt, wird auf `~/.local/bin/swag2mcp` zurückgegriffen.

**`--local`:** Installiert nach `~/.local/bin/swag2mcp` ohne `sudo`. Fügen Sie nach der Installation Folgendes zu Ihrer Shell-Konfiguration hinzu:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

**Überprüfung:**

```bash
swag2mcp --version
```

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp
```

**Überprüfung:**

```bash
swag2mcp --version
```

### GitHub Release

1. Öffnen Sie die [Seite der neuesten Version](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Laden Sie das Archiv für Ihren Mac herunter:
   - **Apple Silicon**: `swag2mcp_darwin_arm64.tar.gz`
   - **Intel**: `swag2mcp_darwin_amd64.tar.gz`
3. Öffnen Sie das Terminal im Download-Ordner und führen Sie aus:

```bash
tar -xzf swag2mcp_darwin_*.tar.gz
sudo mv swag2mcp /usr/local/bin/
```

**Überprüfung:**

```bash
swag2mcp --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

Stellen Sie sicher, dass `$GOPATH/bin` in Ihrem `$PATH` ist:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

**Überprüfung:**

```bash
swag2mcp --version
```

### Aus Quellcode bauen

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

**Überprüfung:**

```bash
swag2mcp --version
```

### Docker

Mounten Sie `~/.swag2mcp` (oder Ihren benutzerdefinierten Pfad) nach `/home/nonroot/.swag2mcp`.
Der Entrypoint passt Dateiberechtigungen automatisch an, damit der Container Ihr Arbeitsverzeichnis lesen kann.

> **Apple Silicon (M1/M2/M3/M4) Benutzer:** Fügen Sie `--platform linux/amd64` hinzu:
> `docker run --rm --platform linux/amd64 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ...`

```bash
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

Mit stdio-Transport ausführen:

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp
```

Mit HTTP-Transport ausführen:

```bash
docker run --rm -p 8080:8080 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp --transport sse --http-addr 0.0.0.0:8080
```

**Überprüfung:**

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest --version
```

---

## Linux

### Einzeiler (empfohlen)

Zwei Installationsmodi — wählen Sie den passenden:

```bash
# Mit sudo (installiert nach /usr/local/bin — empfohlen)
curl -fsSL https://swag2mcp.io/install.sh | bash

# Ohne sudo (installiert nach ~/.local/bin)
curl -fsSL https://swag2mcp.io/install.sh | bash -s -- --local
```

**`--sudo` (Standard):** Installiert nach `/usr/local/bin/swag2mcp` mit `sudo`. Sie werden nach Ihrem Passwort gefragt. Wenn `sudo` fehlschlägt, wird auf `~/.local/bin/swag2mcp` zurückgegriffen.

**`--local`:** Installiert nach `~/.local/bin/swag2mcp` ohne `sudo`. Fügen Sie nach der Installation Folgendes zu Ihrer Shell-Konfiguration hinzu:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

**Überprüfung:**

```bash
swag2mcp --version
```

### APT (Debian / Ubuntu)

1. Öffnen Sie die [Seite der neuesten Version](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Laden Sie `swag2mcp_linux_amd64.deb` herunter
3. Installieren:

```bash
# Stellen Sie sicher, dass Sie sich im Download-Ordner befinden
pwd
sudo dpkg -i swag2mcp_linux_amd64.deb
```

**Überprüfung:**

```bash
swag2mcp --version
```

### RPM (Fedora / RHEL)

1. Öffnen Sie die [Seite der neuesten Version](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Laden Sie `swag2mcp_linux_amd64.rpm` herunter
3. Installieren:

```bash
# Stellen Sie sicher, dass Sie sich im Download-Ordner befinden
pwd
sudo rpm -i swag2mcp_linux_amd64.rpm
```

**Überprüfung:**

```bash
swag2mcp --version
```

### Docker

Mounten Sie `~/.swag2mcp` (oder Ihren benutzerdefinierten Pfad) nach `/home/nonroot/.swag2mcp`.
Der Entrypoint passt Dateiberechtigungen automatisch an, damit der Container Ihr Arbeitsverzeichnis lesen kann.

> **Apple Silicon (M1/M2/M3/M4) Benutzer:** Fügen Sie `--platform linux/amd64` hinzu:
> `docker run --rm --platform linux/amd64 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ...`

```bash
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

Mit stdio-Transport ausführen:

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp
```

Mit HTTP-Transport ausführen:

```bash
docker run --rm -p 8080:8080 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp --transport sse --http-addr 0.0.0.0:8080
```

**Überprüfung:**

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest --version
```

### GitHub Release

1. Öffnen Sie die [Seite der neuesten Version](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Laden Sie das Archiv für Ihre Architektur herunter:
   - **amd64**: `swag2mcp_linux_amd64.tar.gz`
   - **arm64**: `swag2mcp_linux_arm64.tar.gz`
3. Öffnen Sie das Terminal im Download-Ordner und führen Sie aus:

```bash
tar -xzf swag2mcp_linux_*.tar.gz
sudo mv swag2mcp /usr/local/bin/
```

**Überprüfung:**

```bash
swag2mcp --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

Stellen Sie sicher, dass `$GOPATH/bin` in Ihrem `$PATH` ist:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

**Überprüfung:**

```bash
swag2mcp --version
```

### Aus Quellcode bauen

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

**Überprüfung:**

```bash
swag2mcp --version
```

---

## Windows

### Scoop

```powershell
scoop bucket add mmadfox https://github.com/mmadfox/scoop-bucket
scoop install mmadfox/swag2mcp
```

**Überprüfung:**

```powershell
swag2mcp --version
```

### GitHub Release

1. Öffnen Sie die [Seite der neuesten Version](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Laden Sie `swag2mcp_windows_amd64.zip` herunter
3. Entpacken Sie die ZIP-Datei (rechtsklick → Alle extrahieren, oder mit PowerShell)
4. Verschieben Sie `swag2mcp.exe` nach `C:\Windows\System32\`

**Überprüfung:**

```powershell
swag2mcp --version
```

### go install

```powershell
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

**Überprüfung:**

```powershell
swag2mcp --version
```

### Aus Quellcode bauen

```powershell
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp.exe ./cmd/swag2mcp
```

**Überprüfung:**

```powershell
swag2mcp --version
```

### Docker

Mounten Sie `~/.swag2mcp` (oder Ihren benutzerdefinierten Pfad) nach `/home/nonroot/.swag2mcp`. Der Entrypoint passt Dateiberechtigungen automatisch an, damit der Container Ihr Arbeitsverzeichnis lesen kann.

```powershell
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

Run with stdio transport:

```powershell
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp
```

Run with HTTP transport:

```powershell
docker run --rm -p 8080:8080 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp --transport sse --http-addr 0.0.0.0:8080
```

Verify:

```powershell
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest --version
```

---

> Benötigen Sie den Mock-Server? Siehe [Mock Server installieren](install-mock.md).

## Über LLM-Agenten installieren

Wenn Sie eine KI-gestützte IDE (OpenCode, Cursor, Claude Desktop, VS Code usw.) verwenden, lesen Sie [Skills für LLM](/de/llm-skills/overview) für Anweisungen zur Installation von Skills über den Agenten oder manuell.

> Einige IDEs und LLM-Clients erfordern einen Neustart nach dem Hinzufügen von Skills.

---

## Überprüfung

```bash
swag2mcp --version
```

Erwartete Ausgabe (Version kann variieren):

```
swag2mcp v*.*.*
```

---

## Nächste Schritte

- [Schnellstart](quickstart.md) — in 2 Minuten startklar
