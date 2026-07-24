# Установка

## Требования

- **macOS, Linux или Windows** (amd64 / arm64)
- **Go 1.26+** (только для `go install` или сборки из исходников)

## Совместимость

| Метод | macOS | Linux | Windows |
|--------|-------|-------|---------|
| One-liner (curl) | ✅ | ✅ | ❌ |
| Homebrew | ✅ | ✅ | ❌ |
| Scoop | ❌ | ❌ | ✅ |
| APT (deb) | ❌ | ✅ | ❌ |
| RPM | ❌ | ✅ | ❌ |
| Docker | ❌ | ✅ | ❌ |
| GitHub Release | ✅ | ✅ | ✅ |
| go install | ✅ | ✅ | ✅ |
| Сборка из исходников | ✅ | ✅ | ✅ |

---

## macOS

### One-liner (рекомендовано)

```bash
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash
```

Устанавливается в `/usr/local/bin/swag2mcp` (или `~/.local/bin/swag2mcp`, если `/usr/local/bin` недоступен для записи).

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

Убедитесь, что `$GOPATH/bin` находится в `$PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Сборка из исходников

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

---

## Linux

### One-liner (рекомендовано)

```bash
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash
```

Устанавливается в `/usr/local/bin/swag2mcp` (или `~/.local/bin/swag2mcp`, если `/usr/local/bin` недоступен для записи).

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp
```

### APT (Debian / Ubuntu)

```bash
# Скачайте .deb из последнего релиза
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_linux_amd64.deb
sudo dpkg -i swag2mcp_linux_amd64.deb
```

### RPM (Fedora / RHEL)

```bash
# Скачайте .rpm из последнего релиза
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_linux_amd64.rpm
sudo rpm -i swag2mcp_linux_amd64.rpm
```

### Docker

```bash
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

Запуск с stdio-транспортом:

```bash
docker run --rm -i ghcr.io/mmadfox/swag2mcp:latest swag2mcp mcp
```

Запуск с HTTP-транспортом:

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

Убедитесь, что `$GOPATH/bin` находится в `$PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Сборка из исходников

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
# Скачайте последний релиз
curl -LO https://github.com/mmadfox/swag2mcp/releases/latest/download/swag2mcp_windows_amd64.zip
Expand-Archive swag2mcp_windows_amd64.zip -DestinationPath .
move swag2mcp.exe C:\Windows\System32\
```

### go install

```powershell
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

### Сборка из исходников

```powershell
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp.exe ./cmd/swag2mcp
```

---

## Mock-сервер

Бинарный файл `swag2mcp-mock` доступен для отдельной загрузки. Установите его тем же способом, что и основной бинарник:

::: code-group

```bash [macOS / Linux]
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

```powershell [Windows]
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

:::

Или скачайте с [GitHub Releases](https://github.com/mmadfox/swag2mcp/releases) — ищите `swag2mcp-mock_<version>_<os>_<arch>.tar.gz`.

---

## Установка через LLM-агента

Если вы используете AI-помощника в IDE (OpenCode, Cursor, Claude Desktop, VS Code и др.), вы можете установить swag2mcp через агента:

1. Попросите агента добавить навыки swag2mcp:

   ```
   "Создай директорию .agents/skills/swag2mcp-cli и добавь навык из https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md в .agents/skills/swag2mcp-cli/SKILL.md"
   "Создай директорию .agents/skills/swag2mcp-format и добавь навык из https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md в .agents/skills/swag2mcp-format/SKILL.md"
   ```

2. Затем скажите агенту:

   ```
   "Настрой swag2mcp"
   ```

   Агент скачает и установит swag2mcp, затем создаст рабочую область с готовыми к использованию спецификациями.

> Некоторым IDE требуется перезапуск после добавления навыков.

---

## Проверка

```bash
swag2mcp --version
```

Ожидаемый вывод (версия может отличаться):

```
swag2mcp v*.*.*
```

---

## Следующие шаги

- [Быстрый старт](quickstart.md) — запустите за 2 минуты
