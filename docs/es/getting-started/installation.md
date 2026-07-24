# Instalación

## Requisitos

- **macOS, Linux o Windows** (amd64 / arm64)
- **Go 1.26+** (solo para `go install` o compilar desde el código fuente)

## Compatibilidad

| Método | macOS | Linux | Windows |
|--------|-------|-------|---------|
| Una línea (curl) | ✅ | ✅ | ❌ |
| Homebrew | ✅ | ✅ | ❌ |
| Scoop | ❌ | ❌ | ✅ |
| APT (deb) | ❌ | ✅ | ❌ |
| RPM | ❌ | ✅ | ❌ |
| Docker | ❌ | ✅ | ❌ |
| GitHub Release | ✅ | ✅ | ✅ |
| go install | ✅ | ✅ | ✅ |
| Compilar desde fuente | ✅ | ✅ | ✅ |

---

## macOS

### Una línea (recomendado)

```bash
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash
```

Instala en `/usr/local/bin/swag2mcp` (o `~/.local/bin/swag2mcp` si `/usr/local/bin` no es escribible).

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

Asegúrese de que `$GOPATH/bin` esté en su `$PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Compilar desde el código fuente

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

---

## Linux

### Una línea (recomendado)

```bash
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash
```

Instala en `/usr/local/bin/swag2mcp` (o `~/.local/bin/swag2mcp` si `/usr/local/bin` no es escribible).

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp
```

### APT (Debian / Ubuntu)

```bash
# Descargue el .deb desde la última versión
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_linux_amd64.deb
sudo dpkg -i swag2mcp_linux_amd64.deb
```

### RPM (Fedora / RHEL)

```bash
# Descargue el .rpm desde la última versión
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_linux_amd64.rpm
sudo rpm -i swag2mcp_linux_amd64.rpm
```

### Docker

```bash
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

Ejecutar con transporte stdio:

```bash
docker run --rm -i ghcr.io/mmadfox/swag2mcp:latest swag2mcp mcp
```

Ejecutar con transporte HTTP:

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

Asegúrese de que `$GOPATH/bin` esté en su `$PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Compilar desde el código fuente

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
# Descargue la última versión
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_windows_amd64.zip
Expand-Archive swag2mcp_windows_amd64.zip -DestinationPath .
move swag2mcp.exe C:\Windows\System32\
```

### go install

```powershell
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

### Compilar desde el código fuente

```powershell
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp.exe ./cmd/swag2mcp
```

---

## Servidor Simulado

El binario `swag2mcp-mock` está disponible como una descarga separada. Instálelo usando el mismo método que el binario principal:

::: code-group

```bash [macOS / Linux]
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

```powershell [Windows]
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

:::

O descárguelo desde [GitHub Releases](https://github.com/mmadfox/swag2mcp/releases) — busque `swag2mcp-mock_<version>_<os>_<arch>.tar.gz`.

---

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
