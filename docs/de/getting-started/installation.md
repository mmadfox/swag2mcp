# Installation

## Anforderungen

- **macOS, Linux oder Windows** (amd64 / arm64)
- **Go 1.26+** (nur für `go install` oder Bauen aus dem Quellcode)

## Kompatibilität

| Methode | macOS | Linux | Windows |
|---------|-------|-------|---------|
| Einzeiler (curl) | ✅ | ✅ | ❌ |
| Homebrew | ✅ | ✅ | ❌ |
| Scoop | ❌ | ❌ | ✅ |
| APT (deb) | ❌ | ✅ | ❌ |
| RPM | ❌ | ✅ | ❌ |
| Docker | ❌ | ✅ | ❌ |
| GitHub Release | ✅ | ✅ | ✅ |
| go install | ✅ | ✅ | ✅ |
| Aus Quellcode bauen | ✅ | ✅ | ✅ |

---

## macOS

### Einzeiler (empfohlen)

```bash
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash
```

Installiert nach `/usr/local/bin/swag2mcp` (oder `~/.local/bin/swag2mcp`, wenn `/usr/local/bin` nicht beschreibbar ist).

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp
```

### GitHub Release

::: code-group

```bash [Apple Silicon]
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_darwin_arm64.tar.gz
tar -xzf swag2mcp_darwin_arm64.tar.gz
sudo mv swag2mcp /usr/local/bin/
```

```bash [Intel]
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_darwin_amd64.tar.gz
tar -xzf swag2mcp_darwin_amd64.tar.gz
sudo mv swag2mcp /usr/local/bin/
```

:::

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

Stellen Sie sicher, dass `$GOPATH/bin` in Ihrem `$PATH` ist:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Aus Quellcode bauen

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

---

## Linux

### Einzeiler (empfohlen)

```bash
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash
```

Installiert nach `/usr/local/bin/swag2mcp` (oder `~/.local/bin/swag2mcp`, wenn `/usr/local/bin` nicht beschreibbar ist).

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp
```

### APT (Debian / Ubuntu)

```bash
# .deb von der neuesten Version herunterladen
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_linux_amd64.deb
sudo dpkg -i swag2mcp_linux_amd64.deb
```

### RPM (Fedora / RHEL)

```bash
# .rpm von der neuesten Version herunterladen
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_linux_amd64.rpm
sudo rpm -i swag2mcp_linux_amd64.rpm
```

### Docker

```bash
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

Mit stdio-Transport ausführen:

```bash
docker run --rm -i ghcr.io/mmadfox/swag2mcp:latest swag2mcp mcp
```

Mit HTTP-Transport ausführen:

```bash
docker run --rm -p 8080:8080 ghcr.io/mmadfox/swag2mcp:latest swag2mcp mcp --transport sse --http-addr :8080
```

### GitHub Release

::: code-group

```bash [amd64]
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_linux_amd64.tar.gz
tar -xzf swag2mcp_linux_amd64.tar.gz
sudo mv swag2mcp /usr/local/bin/
```

```bash [arm64]
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_linux_arm64.tar.gz
tar -xzf swag2mcp_linux_arm64.tar.gz
sudo mv swag2mcp /usr/local/bin/
```

:::

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

Stellen Sie sicher, dass `$GOPATH/bin` in Ihrem `$PATH` ist:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Aus Quellcode bauen

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

---

## Windows

### Scoop

```powershell
scoop bucket add mmadfox https://github.com/mmadfox/scoop-bucket
scoop install mmadfox/swag2mcp
```

### GitHub Release

```powershell
# Die neueste Version herunterladen
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_windows_amd64.zip
Expand-Archive swag2mcp_windows_amd64.zip -DestinationPath .
move swag2mcp.exe C:\Windows\System32\
```

### go install

```powershell
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

### Aus Quellcode bauen

```powershell
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp.exe ./cmd/swag2mcp
```

---

## Mock-Server

Die `swag2mcp-mock`-Binärdatei ist als separater Download verfügbar. Installieren Sie sie mit derselben Methode wie die Hauptbinärdatei:

::: code-group

```bash [macOS / Linux]
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

```powershell [Windows]
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

:::

Oder laden Sie sie von [GitHub Releases](https://github.com/mmadfox/swag2mcp/releases) herunter — suchen Sie nach `swag2mcp-mock_&lt;version&gt;_&lt;os&gt;_&lt;arch&gt;.tar.gz`.

---

## Über LLM-Agenten installieren

Wenn Sie eine KI-gestützte IDE (OpenCode, Cursor, Claude Desktop, VS Code usw.) verwenden, können Sie swag2mcp über Ihren Agenten installieren:

1. Bitten Sie Ihren Agenten, die swag2mcp-Skills hinzuzufügen:

   ```
   "Erstellen Sie das Verzeichnis .agents/skills/swag2mcp-cli und fügen Sie den Skill von https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md zu .agents/skills/swag2mcp-cli/SKILL.md hinzu"
   "Erstellen Sie das Verzeichnis .agents/skills/swag2mcp-format und fügen Sie den Skill von https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md zu .agents/skills/swag2mcp-format/SKILL.md hinzu"
   ```

2. Sagen Sie dann Ihrem Agenten:

   ```
   "Richte swag2mcp ein"
   ```

   Der Agent wird swag2mcp herunterladen und installieren und dann einen Arbeitsbereich mit gebrauchsfertigen Spezifikationen erstellen.

> Einige IDEs erfordern einen Neustart nach dem Hinzufügen von Skills.

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
