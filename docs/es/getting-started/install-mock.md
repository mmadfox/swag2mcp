# Instalar el Servidor Mock

El binario `swag2mcp-mock` es una herramienta separada para probar APIs con respuestas simuladas. Soporta los mismos métodos de instalación que el binario principal.

## Compatibilidad

| Método | macOS | Linux | Windows |
|--------|-------|-------|---------|
| Una línea (curl) | ✅ | ✅ | ❌ |
| Homebrew | ✅ | ❌ | ❌ |
| Scoop | ❌ | ❌ | ✅ |
| GitHub Release | ✅ | ✅ | ✅ |
| go install | ✅ | ✅ | ✅ |
| Compilar desde el código fuente | ✅ | ✅ | ✅ |
| Docker | ✅ | ✅ | ✅ |

---

## macOS

### Una línea (recomendado)

Dos modos de instalación — elija el que mejor se adapte:

```bash
# Con sudo (instala en /usr/local/bin — recomendado)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash

# Sin sudo (instala en ~/.local/bin)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash -s -- --local
```

**`--sudo` (predeterminado):** Instala en `/usr/local/bin/swag2mcp-mock` usando `sudo`. Se le solicitará su contraseña. Si `sudo` falla, recurre a `~/.local/bin/swag2mcp-mock`.

**`--local`:** Instala en `~/.local/bin/swag2mcp-mock` sin `sudo`. Después de la instalación, agregue a su configuración de shell:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Verificación:

```bash
swag2mcp-mock --version
```

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp-mock
```

Verificación:

```bash
swag2mcp-mock --version
```

### GitHub Release

1. Abra la [página de la última versión](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Descargue el archivo para su Mac:
   - **Apple Silicon**: `swag2mcp-mock_darwin_arm64.tar.gz`
   - **Intel**: `swag2mcp-mock_darwin_amd64.tar.gz`
3. Abra la Terminal en la carpeta de descargas y ejecute:

```bash
tar -xzf swag2mcp-mock_darwin_*.tar.gz
sudo mv swag2mcp-mock /usr/local/bin/
```

Verificación:

```bash
swag2mcp-mock --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

Asegúrese de que `$GOPATH/bin` esté en su `$PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Verificación:

```bash
swag2mcp-mock --version
```

### Compilar desde el código fuente

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp-mock ./cmd/swag2mcp-mock
sudo mv swag2mcp-mock /usr/local/bin/
```

Verificación:

```bash
swag2mcp-mock --version
```

### Docker

Monte `~/.swag2mcp` (o su ruta de workspace personalizada) en `/home/nonroot/.swag2mcp`.
El entrypoint ajusta automáticamente los permisos de archivos para que el contenedor pueda leer su workspace.

```bash
docker pull ghcr.io/mmadfox/swag2mcp-mock:latest
```

Ejecutar el servidor mock:

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest
```

Verificación:

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest --version
```

---

## Linux

### Una línea (recomendado)

Dos modos de instalación — elija el que mejor se adapte:

```bash
# Con sudo (instala en /usr/local/bin — recomendado)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash

# Sin sudo (instala en ~/.local/bin)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash -s -- --local
```

**`--sudo` (predeterminado):** Instala en `/usr/local/bin/swag2mcp-mock` usando `sudo`. Se le solicitará su contraseña. Si `sudo` falla, recurre a `~/.local/bin/swag2mcp-mock`.

**`--local`:** Instala en `~/.local/bin/swag2mcp-mock` sin `sudo`. Después de la instalación, agregue a su configuración de shell:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Verificación:

```bash
swag2mcp-mock --version
```

### GitHub Release

1. Abra la [página de la última versión](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Descargue el archivo para su arquitectura:
   - **amd64**: `swag2mcp-mock_linux_amd64.tar.gz`
   - **arm64**: `swag2mcp-mock_linux_arm64.tar.gz`
3. Abra la Terminal en la carpeta de descargas y ejecute:

```bash
tar -xzf swag2mcp-mock_linux_*.tar.gz
sudo mv swag2mcp-mock /usr/local/bin/
```

Verificación:

```bash
swag2mcp-mock --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

Asegúrese de que `$GOPATH/bin` esté en su `$PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Verificación:

```bash
swag2mcp-mock --version
```

### Compilar desde el código fuente

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp-mock ./cmd/swag2mcp-mock
sudo mv swag2mcp-mock /usr/local/bin/
```

Verificación:

```bash
swag2mcp-mock --version
```

### Docker

Monte `~/.swag2mcp` (o su ruta de workspace personalizada) en `/home/nonroot/.swag2mcp`.
El entrypoint ajusta automáticamente los permisos de archivos para que el contenedor pueda leer su workspace.

```bash
docker pull ghcr.io/mmadfox/swag2mcp-mock:latest
```

Ejecutar el servidor mock:

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest
```

Verificación:

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

Verificación:

```powershell
swag2mcp-mock --version
```

### GitHub Release

1. Abra la [página de la última versión](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Descargue `swag2mcp-mock_windows_amd64.zip`
3. Extraiga el archivo ZIP (clic derecho → Extraer todo, o use PowerShell)
4. Mueva `swag2mcp-mock.exe` a `C:\Windows\System32\`

Verificación:

```powershell
swag2mcp-mock --version
```

### go install

```powershell
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

Verificación:

```powershell
swag2mcp-mock --version
```

### Compilar desde el código fuente

```powershell
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp-mock.exe ./cmd/swag2mcp-mock
```

Verificación:

```powershell
swag2mcp-mock --version
```

---

## Próximos pasos

- [Inicio rápido](quickstart.md) — póngase en marcha en 2 minutos
- [Servidor Mock](../advanced/mock-server.md) — configurar y ejecutar el servidor mock
