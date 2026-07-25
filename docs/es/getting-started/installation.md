# Instalación

## Requisitos

- **macOS, Linux o Windows** (amd64 / arm64)
- **Go 1.26+** (solo para `go install` o compilar desde el código fuente)

## Compatibilidad

| Método | macOS | Linux | Windows |
|--------|-------|-------|---------|
| Una línea (curl) | ✅ | ✅ | ❌ |
| Homebrew | ✅ | ❌ | ❌ |
| Scoop | ❌ | ❌ | ✅ |
| APT (deb) | ❌ | ✅ | ❌ |
| RPM | ❌ | ✅ | ❌ |
| Docker | ✅ | ✅ | ✅ |
| GitHub Release | ✅ | ✅ | ✅ |
| go install | ✅ | ✅ | ✅ |
| Compilar desde fuente | ✅ | ✅ | ✅ |

---

## macOS

### Una línea (recomendado)

Dos modos de instalación — elija el que mejor se adapte:

```bash
# Con sudo (instala en /usr/local/bin — recomendado)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash

# Sin sudo (instala en ~/.local/bin)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash -s -- --local
```

**`--sudo` (predeterminado):** Instala en `/usr/local/bin/swag2mcp` usando `sudo`. Se le solicitará su contraseña. Si `sudo` falla, recurre a `~/.local/bin/swag2mcp`.

**`--local`:** Instala en `~/.local/bin/swag2mcp` sin `sudo`. Después de la instalación, agregue a su configuración de shell:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

**Verificación:**

```bash
swag2mcp --version
```

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp
```

**Verificación:**

```bash
swag2mcp --version
```

### GitHub Release

1. Abra la [página de la última versión](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Descargue el archivo para su Mac:
   - **Apple Silicon**: `swag2mcp_darwin_arm64.tar.gz`
   - **Intel**: `swag2mcp_darwin_amd64.tar.gz`
3. Abra la Terminal en la carpeta de descargas y ejecute:

```bash
tar -xzf swag2mcp_darwin_*.tar.gz
sudo mv swag2mcp /usr/local/bin/
```

**Verificación:**

```bash
swag2mcp --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

Asegúrese de que `$GOPATH/bin` esté en su `$PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

**Verificación:**

```bash
swag2mcp --version
```

### Compilar desde el código fuente

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

**Verificación:**

```bash
swag2mcp --version
```

### Docker

El contenedor se ejecuta como un usuario no root y necesita acceso a su directorio de trabajo. Monte `~/.swag2mcp` (o su ruta personalizada) en `/home/nonroot/.swag2mcp`:

> **Usuarios de Apple Silicon (M1/M2/M3/M4):** Agregue `--platform linux/amd64` para ejecutar la imagen:
> `docker run --rm --platform linux/amd64 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ...`

```bash
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

Ejecutar con transporte stdio:

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp
```

Ejecutar con transporte HTTP:

```bash
docker run --rm -p 8080:8080 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp --transport sse --http-addr :8080
```

**Verificación:**

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest --version
```

---

## Linux

### Una línea (recomendado)

Dos modos de instalación — elija el que mejor se adapte:

```bash
# Con sudo (instala en /usr/local/bin — recomendado)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash

# Sin sudo (instala en ~/.local/bin)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash -s -- --local
```

**`--sudo` (predeterminado):** Instala en `/usr/local/bin/swag2mcp` usando `sudo`. Se le solicitará su contraseña. Si `sudo` falla, recurre a `~/.local/bin/swag2mcp`.

**`--local`:** Instala en `~/.local/bin/swag2mcp` sin `sudo`. Después de la instalación, agregue a su configuración de shell:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

**Verificación:**

```bash
swag2mcp --version
```

### APT (Debian / Ubuntu)

1. Abra la [página de la última versión](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Descargue `swag2mcp_linux_amd64.deb`
3. Instale:

```bash
# Asegúrese de estar en la carpeta de descargas
pwd
sudo dpkg -i swag2mcp_linux_amd64.deb
```

**Verificación:**

```bash
swag2mcp --version
```

### RPM (Fedora / RHEL)

1. Abra la [página de la última versión](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Descargue `swag2mcp_linux_amd64.rpm`
3. Instale:

```bash
# Asegúrese de estar en la carpeta de descargas
pwd
sudo rpm -i swag2mcp_linux_amd64.rpm
```

**Verificación:**

```bash
swag2mcp --version
```

### Docker

El contenedor se ejecuta como un usuario no root y necesita acceso a su directorio de trabajo. Monte `~/.swag2mcp` (o su ruta personalizada) en `/home/nonroot/.swag2mcp`:

> **Usuarios de Apple Silicon (M1/M2/M3/M4):** Agregue `--platform linux/amd64` para ejecutar la imagen:
> `docker run --rm --platform linux/amd64 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ...`

```bash
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

Ejecutar con transporte stdio:

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp
```

Ejecutar con transporte HTTP:

```bash
docker run --rm -p 8080:8080 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp --transport sse --http-addr :8080
```

**Verificación:**

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest --version
```

### GitHub Release

1. Abra la [página de la última versión](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Descargue el archivo para su arquitectura:
   - **amd64**: `swag2mcp_linux_amd64.tar.gz`
   - **arm64**: `swag2mcp_linux_arm64.tar.gz`
3. Abra la Terminal en la carpeta de descargas y ejecute:

```bash
tar -xzf swag2mcp_linux_*.tar.gz
sudo mv swag2mcp /usr/local/bin/
```

**Verificación:**

```bash
swag2mcp --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

Asegúrese de que `$GOPATH/bin` esté en su `$PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

**Verificación:**

```bash
swag2mcp --version
```

### Compilar desde el código fuente

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

**Verificación:**

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

**Verificación:**

```powershell
swag2mcp --version
```

### GitHub Release

1. Abra la [página de la última versión](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Descargue `swag2mcp_windows_amd64.zip`
3. Extraiga el archivo ZIP (clic derecho → Extraer todo, o con PowerShell)
4. Mueva `swag2mcp.exe` a `C:\Windows\System32\`

**Verificación:**

```powershell
swag2mcp --version
```

### go install

```powershell
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

**Verificación:**

```powershell
swag2mcp --version
```

### Compilar desde el código fuente

```powershell
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp.exe ./cmd/swag2mcp
```

**Verificación:**

```powershell
swag2mcp --version
```

---

> ¿Necesita el servidor mock? Consulte [Instalar servidor mock](install-mock.md).

## Instalar mediante Agente LLM

Si usa un IDE impulsado por IA (OpenCode, Cursor, Claude Desktop, VS Code, etc.), puede instalar swag2mcp a través de su agente:

1. Pídale a su agente que agregue las habilidades de swag2mcp:

   ```
   "Cree el directorio .agents/skills/swag2mcp-cli y agregue la habilidad desde https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md a .agents/skills/swag2mcp-cli/SKILL.md"
   "Cree el directorio .agents/skills/swag2mcp-format y agregue la habilidad desde https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md a .agents/skills/swag2mcp-format/SKILL.md"
   ```

2. Luego dígale a su agente:

   ```
   "Configura swag2mcp"
   ```

   El agente descargará e instalará swag2mcp, luego creará un espacio de trabajo con especificaciones listas para usar.

> Algunos IDEs requieren un reinicio después de agregar habilidades.

---

## Verificar

```bash
swag2mcp --version
```

Salida esperada (la versión puede variar):

```
swag2mcp v*.*.*
```

---

## Próximos Pasos

- [Inicio Rápido](quickstart.md) — póngase en marcha en 2 minutos
