# Installation

## Prérequis

- **macOS, Linux ou Windows** (amd64 / arm64)
- **Go 1.26+** (uniquement pour `go install` ou la compilation depuis les sources)

## Compatibilité

| Méthode | macOS | Linux | Windows |
|---------|-------|-------|---------|
| One-liner (curl) | ✅ | ✅ | ❌ |
| Homebrew | ✅ | ❌ | ❌ |
| Scoop | ❌ | ❌ | ✅ |
| APT (deb) | ❌ | ✅ | ❌ |
| RPM | ❌ | ✅ | ❌ |
| Docker | ✅ | ✅ | ✅ |
| GitHub Release | ✅ | ✅ | ✅ |
| go install | ✅ | ✅ | ✅ |
| Compilation depuis les sources | ✅ | ✅ | ✅ |

---

## macOS

### One-liner (recommandé)

Deux modes d'installation — choisissez celui qui vous convient :

```bash
# Avec sudo (installe dans /usr/local/bin — recommandé)
curl -fsSL https://swag2mcp.io/install.sh | bash

# Sans sudo (installe dans ~/.local/bin)
curl -fsSL https://swag2mcp.io/install.sh | bash -s -- --local
```

**`--sudo` (par défaut) :** Installe dans `/usr/local/bin/swag2mcp` via `sudo`. Vous serez invité à saisir votre mot de passe. Si `sudo` échoue, il utilise `~/.local/bin/swag2mcp` comme solution de repli.

**`--local` :** Installe dans `~/.local/bin/swag2mcp` sans `sudo`. Après l'installation, ajoutez à votre configuration de shell :

```bash
export PATH="$HOME/.local/bin:$PATH"
```

**Vérification :**

```bash
swag2mcp --version
```

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp
```

**Vérification :**

```bash
swag2mcp --version
```

### GitHub Release

1. Ouvrez la [page de la dernière version](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Téléchargez l'archive pour votre Mac :
   - **Apple Silicon** : `swag2mcp_darwin_arm64.tar.gz`
   - **Intel** : `swag2mcp_darwin_amd64.tar.gz`
3. Ouvrez le Terminal dans le dossier de téléchargement et exécutez :

```bash
tar -xzf swag2mcp_darwin_*.tar.gz
sudo mv swag2mcp /usr/local/bin/
```

**Vérification :**

```bash
swag2mcp --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

Assurez-vous que `$GOPATH/bin` est dans votre `$PATH` :

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

**Vérification :**

```bash
swag2mcp --version
```

### Compilation depuis les sources

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

**Vérification :**

```bash
swag2mcp --version
```

### Docker

Montez `~/.swag2mcp` (ou votre chemin personnalisé) vers `/home/nonroot/.swag2mcp`.
L'entrypoint ajuste automatiquement les permissions des fichiers pour que le conteneur puisse lire votre espace de travail.

> **Utilisateurs Apple Silicon (M1/M2/M3/M4) :** Ajoutez `--platform linux/amd64` pour exécuter l'image :
> `docker run --rm --platform linux/amd64 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ...`

```bash
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

Exécution avec transport stdio :

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp
```

Exécution avec transport HTTP :

```bash
docker run --rm -p 8080:8080 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp --transport sse --http-addr 0.0.0.0:8080
```

**Vérification :**

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest --version
```

---

## Linux

### One-liner (recommandé)

Deux modes d'installation — choisissez celui qui vous convient :

```bash
# Avec sudo (installe dans /usr/local/bin — recommandé)
curl -fsSL https://swag2mcp.io/install.sh | bash

# Sans sudo (installe dans ~/.local/bin)
curl -fsSL https://swag2mcp.io/install.sh | bash -s -- --local
```

**`--sudo` (par défaut) :** Installe dans `/usr/local/bin/swag2mcp` via `sudo`. Vous serez invité à saisir votre mot de passe. Si `sudo` échoue, il utilise `~/.local/bin/swag2mcp` comme solution de repli.

**`--local` :** Installe dans `~/.local/bin/swag2mcp` sans `sudo`. Après l'installation, ajoutez à votre configuration de shell :

```bash
export PATH="$HOME/.local/bin:$PATH"
```

**Vérification :**

```bash
swag2mcp --version
```

### APT (Debian / Ubuntu)

1. Ouvrez la [page de la dernière version](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Téléchargez `swag2mcp_linux_amd64.deb`
3. Installez :

```bash
# Assurez-vous d'être dans le dossier de téléchargement
pwd
sudo dpkg -i swag2mcp_linux_amd64.deb
```

**Vérification :**

```bash
swag2mcp --version
```

### RPM (Fedora / RHEL)

1. Ouvrez la [page de la dernière version](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Téléchargez `swag2mcp_linux_amd64.rpm`
3. Installez :

```bash
# Assurez-vous d'être dans le dossier de téléchargement
pwd
sudo rpm -i swag2mcp_linux_amd64.rpm
```

**Vérification :**

```bash
swag2mcp --version
```

### Docker

Montez `~/.swag2mcp` (ou votre chemin personnalisé) vers `/home/nonroot/.swag2mcp`.
L'entrypoint ajuste automatiquement les permissions des fichiers pour que le conteneur puisse lire votre espace de travail.

> **Utilisateurs Apple Silicon (M1/M2/M3/M4) :** Ajoutez `--platform linux/amd64` pour exécuter l'image :
> `docker run --rm --platform linux/amd64 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ...`

```bash
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

Exécution avec transport stdio :

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp
```

Exécution avec transport HTTP :

```bash
docker run --rm -p 8080:8080 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp --transport sse --http-addr 0.0.0.0:8080
```

**Vérification :**

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest --version
```

### GitHub Release

1. Ouvrez la [page de la dernière version](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Téléchargez l'archive pour votre architecture :
   - **amd64** : `swag2mcp_linux_amd64.tar.gz`
   - **arm64** : `swag2mcp_linux_arm64.tar.gz`
3. Ouvrez le Terminal dans le dossier de téléchargement et exécutez :

```bash
tar -xzf swag2mcp_linux_*.tar.gz
sudo mv swag2mcp /usr/local/bin/
```

**Vérification :**

```bash
swag2mcp --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

Assurez-vous que `$GOPATH/bin` est dans votre `$PATH` :

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

**Vérification :**

```bash
swag2mcp --version
```

### Compilation depuis les sources

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

**Vérification :**

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

**Vérification :**

```powershell
swag2mcp --version
```

### GitHub Release

1. Ouvrez la [page de la dernière version](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Téléchargez `swag2mcp_windows_amd64.zip`
3. Extrayez le fichier ZIP (clic droit → Extraire tout, ou avec PowerShell)
4. Déplacez `swag2mcp.exe` vers `C:\Windows\System32\`

**Vérification :**

```powershell
swag2mcp --version
```

### go install

```powershell
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

**Vérification :**

```powershell
swag2mcp --version
```

### Compilation depuis les sources

```powershell
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp.exe ./cmd/swag2mcp
```

**Vérification :**

```powershell
swag2mcp --version
```

### Docker

Montez `~/.swag2mcp` (ou votre chemin personnalisé) vers `/home/nonroot/.swag2mcp`. L'entrypoint ajuste automatiquement les permissions des fichiers pour que le conteneur puisse lire votre espace de travail.

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

> Besoin du serveur mock ? Voir [Installer le serveur mock](install-mock.md).

## Installation via Agent LLM

Si vous utilisez un IDE alimenté par l'IA (OpenCode, Cursor, Claude Desktop, VS Code, etc.), consultez [Skills pour LLM](/fr/llm-skills/overview) pour savoir comment installer les skills via votre agent ou manuellement.

> Certains IDE et clients LLM nécessitent un redémarrage après l'ajout de skills.

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
