# Installation

## Prérequis

- **macOS, Linux ou Windows** (amd64 / arm64)
- **Go 1.26+** (uniquement pour `go install` ou la compilation depuis les sources)

## Compatibilité

| Méthode | macOS | Linux | Windows |
|---------|-------|-------|---------|
| One-liner (curl) | ✅ | ✅ | ❌ |
| Homebrew | ✅ | ✅ | ❌ |
| Scoop | ❌ | ❌ | ✅ |
| APT (deb) | ❌ | ✅ | ❌ |
| RPM | ❌ | ✅ | ❌ |
| Docker | ❌ | ✅ | ❌ |
| GitHub Release | ✅ | ✅ | ✅ |
| go install | ✅ | ✅ | ✅ |
| Compilation depuis les sources | ✅ | ✅ | ✅ |

---

## macOS

### One-liner (recommandé)

```bash
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash
```

Installe dans `/usr/local/bin/swag2mcp` (ou `~/.local/bin/swag2mcp` si `/usr/local/bin` n'est pas accessible en écriture).

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

Assurez-vous que `$GOPATH/bin` est dans votre `$PATH` :

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Compilation depuis les sources

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

---

## Linux

### One-liner (recommandé)

```bash
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash
```

Installe dans `/usr/local/bin/swag2mcp` (ou `~/.local/bin/swag2mcp` si `/usr/local/bin` n'est pas accessible en écriture).

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp
```

### APT (Debian / Ubuntu)

```bash
# Téléchargez le .deb depuis la dernière version
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_linux_amd64.deb
sudo dpkg -i swag2mcp_linux_amd64.deb
```

### RPM (Fedora / RHEL)

```bash
# Téléchargez le .rpm depuis la dernière version
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_linux_amd64.rpm
sudo rpm -i swag2mcp_linux_amd64.rpm
```

### Docker

```bash
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

Exécution avec transport stdio :

```bash
docker run --rm -i ghcr.io/mmadfox/swag2mcp:latest swag2mcp mcp
```

Exécution avec transport HTTP :

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

Assurez-vous que `$GOPATH/bin` est dans votre `$PATH` :

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Compilation depuis les sources

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
# Téléchargez la dernière version
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_windows_amd64.zip
Expand-Archive swag2mcp_windows_amd64.zip -DestinationPath .
move swag2mcp.exe C:\Windows\System32\
```

### go install

```powershell
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

### Compilation depuis les sources

```powershell
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp.exe ./cmd/swag2mcp
```

---

## Serveur Mock

Le binaire `swag2mcp-mock` est disponible en téléchargement séparé. Installez-le en utilisant la même méthode que le binaire principal :

::: code-group

```bash [macOS / Linux]
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

```powershell [Windows]
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

:::

Ou téléchargez-le depuis [GitHub Releases](https://github.com/mmadfox/swag2mcp/releases) — cherchez `swag2mcp-mock_&lt;version&gt;_&lt;os&gt;_&lt;arch&gt;.tar.gz`.

---

## Installation via Agent LLM

Si vous utilisez un IDE avec IA intégrée (OpenCode, Cursor, Claude Desktop, VS Code, etc.), vous pouvez installer swag2mcp via votre agent :

1. Demandez à votre agent d'ajouter les compétences swag2mcp :

   ```
   "Créez le répertoire .agents/skills/swag2mcp-cli et ajoutez la compétence depuis https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md vers .agents/skills/swag2mcp-cli/SKILL.md"
   "Créez le répertoire .agents/skills/swag2mcp-format et ajoutez la compétence depuis https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md vers .agents/skills/swag2mcp-format/SKILL.md"
   ```

2. Ensuite, dites à votre agent :

   ```
   "Configurez swag2mcp"
   ```

   L'agent téléchargera et installera swag2mcp, puis créera un espace de travail avec des spécifications prêtes à l'emploi.

> Certains IDE nécessitent un redémarrage après l'ajout de compétences.

---

## Vérification

```bash
swag2mcp --version
```

Sortie attendue (la version peut varier) :

```
swag2mcp v*.*.*
```

---

## Prochaines étapes

- [Démarrage rapide](quickstart.md) — opérationnel en 2 minutes
