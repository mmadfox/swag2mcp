# Installer le serveur Mock

Le binaire `swag2mcp-mock` est un outil séparé pour tester des API avec des réponses simulées. Il prend en charge les mêmes méthodes d'installation que le binaire principal.

## Compatibilité

| Méthode | macOS | Linux | Windows |
|--------|-------|-------|---------|
| En une ligne (curl) | ✅ | ✅ | ❌ |
| Homebrew | ✅ | ❌ | ❌ |
| Scoop | ❌ | ❌ | ✅ |
| GitHub Release | ✅ | ✅ | ✅ |
| go install | ✅ | ✅ | ✅ |
| Compilation depuis les sources | ✅ | ✅ | ✅ |
| Docker | ✅ | ✅ | ✅ |

---

## macOS

### En une ligne (recommandé)

Deux modes d'installation — choisissez celui qui vous convient :

```bash
# Avec sudo (installe dans /usr/local/bin — recommandé)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash

# Sans sudo (installe dans ~/.local/bin)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash -s -- --local
```

**`--sudo` (par défaut) :** Installe dans `/usr/local/bin/swag2mcp-mock` avec `sudo`. Votre mot de passe vous sera demandé. Si `sudo` échoue, bascule vers `~/.local/bin/swag2mcp-mock`.

**`--local` :** Installe dans `~/.local/bin/swag2mcp-mock` sans `sudo`. Après l'installation, ajoutez à votre configuration shell :

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Vérification :

```bash
swag2mcp-mock --version
```

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp-mock
```

Vérification :

```bash
swag2mcp-mock --version
```

### GitHub Release

1. Ouvrez la [page de la dernière version](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Téléchargez l'archive pour votre Mac :
   - **Apple Silicon** : `swag2mcp-mock_darwin_arm64.tar.gz`
   - **Intel** : `swag2mcp-mock_darwin_amd64.tar.gz`
3. Ouvrez le Terminal dans le dossier de téléchargement et exécutez :

```bash
tar -xzf swag2mcp-mock_darwin_*.tar.gz
sudo mv swag2mcp-mock /usr/local/bin/
```

Vérification :

```bash
swag2mcp-mock --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

Assurez-vous que `$GOPATH/bin` est dans votre `$PATH` :

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Vérification :

```bash
swag2mcp-mock --version
```

### Compilation depuis les sources

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp-mock ./cmd/swag2mcp-mock
sudo mv swag2mcp-mock /usr/local/bin/
```

Vérification :

```bash
swag2mcp-mock --version
```

### Docker

Le conteneur s'exécute en tant qu'utilisateur non-root et a besoin d'accéder à votre répertoire de travail.
Montez `~/.swag2mcp` (ou votre chemin de workspace personnalisé) vers `/home/nonroot/.swag2mcp` :

> **Utilisateurs Apple Silicon (M1/M2/M3/M4) :** Ajoutez `--platform linux/amd64` pour exécuter l'image :
> `docker run --rm --platform linux/amd64 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ...`

```bash
docker pull ghcr.io/mmadfox/swag2mcp-mock:latest
```

Exécuter le serveur mock :

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest
```

Vérification :

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest --version
```

---

## Linux

### En une ligne (recommandé)

Deux modes d'installation — choisissez celui qui vous convient :

```bash
# Avec sudo (installe dans /usr/local/bin — recommandé)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash

# Sans sudo (installe dans ~/.local/bin)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash -s -- --local
```

**`--sudo` (par défaut) :** Installe dans `/usr/local/bin/swag2mcp-mock` avec `sudo`. Votre mot de passe vous sera demandé. Si `sudo` échoue, bascule vers `~/.local/bin/swag2mcp-mock`.

**`--local` :** Installe dans `~/.local/bin/swag2mcp-mock` sans `sudo`. Après l'installation, ajoutez à votre configuration shell :

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Vérification :

```bash
swag2mcp-mock --version
```

### GitHub Release

1. Ouvrez la [page de la dernière version](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Téléchargez l'archive pour votre architecture :
   - **amd64** : `swag2mcp-mock_linux_amd64.tar.gz`
   - **arm64** : `swag2mcp-mock_linux_arm64.tar.gz`
3. Ouvrez le Terminal dans le dossier de téléchargement et exécutez :

```bash
tar -xzf swag2mcp-mock_linux_*.tar.gz
sudo mv swag2mcp-mock /usr/local/bin/
```

Vérification :

```bash
swag2mcp-mock --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

Assurez-vous que `$GOPATH/bin` est dans votre `$PATH` :

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Vérification :

```bash
swag2mcp-mock --version
```

### Compilation depuis les sources

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp-mock ./cmd/swag2mcp-mock
sudo mv swag2mcp-mock /usr/local/bin/
```

Vérification :

```bash
swag2mcp-mock --version
```

### Docker

Le conteneur s'exécute en tant qu'utilisateur non-root et a besoin d'accéder à votre répertoire de travail.
Montez `~/.swag2mcp` (ou votre chemin de workspace personnalisé) vers `/home/nonroot/.swag2mcp` :

```bash
docker pull ghcr.io/mmadfox/swag2mcp-mock:latest
```

Exécuter le serveur mock :

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest
```

Vérification :

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

Vérification :

```powershell
swag2mcp-mock --version
```

### GitHub Release

1. Ouvrez la [page de la dernière version](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Téléchargez `swag2mcp-mock_windows_amd64.zip`
3. Extrayez le fichier ZIP (clic droit → Extraire tout, ou avec PowerShell)
4. Déplacez `swag2mcp-mock.exe` vers `C:\Windows\System32\`

Vérification :

```powershell
swag2mcp-mock --version
```

### go install

```powershell
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

Vérification :

```powershell
swag2mcp-mock --version
```

### Compilation depuis les sources

```powershell
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp-mock.exe ./cmd/swag2mcp-mock
```

Vérification :

```powershell
swag2mcp-mock --version
```

---

## Prochaines étapes

- [Démarrage rapide](quickstart.md) — lancez-vous en 2 minutes
- [Serveur Mock](../advanced/mock-server.md) — configurer et exécuter le serveur mock
