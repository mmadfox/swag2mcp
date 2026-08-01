# Установка

## Требования

- **macOS, Linux или Windows** (amd64 / arm64)
- **Go 1.26+** (только для `go install` или сборки из исходников)

## Совместимость

| Метод | macOS | Linux | Windows |
|--------|-------|-------|---------|
| One-liner (curl) | ✅ | ✅ | ❌ |
| Homebrew | ✅ | ❌ | ❌ |
| Scoop | ❌ | ❌ | ✅ |
| APT (deb) | ❌ | ✅ | ❌ |
| RPM | ❌ | ✅ | ❌ |
| Docker | ✅ | ✅ | ✅ |
| GitHub Release | ✅ | ✅ | ✅ |
| go install | ✅ | ✅ | ✅ |
| Сборка из исходников | ✅ | ✅ | ✅ |

---

## macOS

### One-liner (рекомендовано)

Два режима установки — выберите подходящий:

```bash
# С sudo (устанавливается в /usr/local/bin — рекомендуется)
curl -fsSL https://swag2mcp.io/install.sh | bash

# Без sudo (устанавливается в ~/.local/bin)
curl -fsSL https://swag2mcp.io/install.sh | bash -s -- --local
```

**`--sudo` (по умолчанию):** Устанавливает в `/usr/local/bin/swag2mcp` с помощью `sudo`. Вам будет предложено ввести пароль. Если `sudo` не сработает, используется резервный путь `~/.local/bin/swag2mcp`.

**`--local`:** Устанавливает в `~/.local/bin/swag2mcp` без `sudo`. После установки добавьте в конфигурацию вашей оболочки:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

**Проверка:**

```bash
swag2mcp --version
```

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp
```

**Проверка:**

```bash
swag2mcp --version
```

### GitHub Release

1. Откройте [страницу последнего релиза](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Скачайте архив для вашего Mac:
   - **Apple Silicon**: `swag2mcp_darwin_arm64.tar.gz`
   - **Intel**: `swag2mcp_darwin_amd64.tar.gz`
3. Откройте Терминал в папке загрузок и выполните:

```bash
tar -xzf swag2mcp_darwin_*.tar.gz
sudo mv swag2mcp /usr/local/bin/
```

**Проверка:**

```bash
swag2mcp --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

Убедитесь, что `$GOPATH/bin` находится в `$PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

**Проверка:**

```bash
swag2mcp --version
```

### Сборка из исходников

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

**Проверка:**

```bash
swag2mcp --version
```

### Docker

Смонтируйте `~/.swag2mcp` (или ваш кастомный путь) в `/home/nonroot/.swag2mcp`.
Entrypoint автоматически настраивает права доступа, чтобы контейнер мог читать вашу рабочую директорию.

> **Пользователи Apple Silicon (M1/M2/M3/M4):** Добавьте `--platform linux/amd64` для запуска образа:
> `docker run --rm --platform linux/amd64 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ...`

```bash
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

Запуск с stdio-транспортом:

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp
```

Запуск с HTTP-транспортом:

```bash
docker run --rm -p 8080:8080 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp --transport sse --http-addr :8080
```

**Проверка:**

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest --version
```

---

## Linux

### One-liner (рекомендовано)

Два режима установки — выберите подходящий:

```bash
# С sudo (устанавливается в /usr/local/bin — рекомендуется)
curl -fsSL https://swag2mcp.io/install.sh | bash

# Без sudo (устанавливается в ~/.local/bin)
curl -fsSL https://swag2mcp.io/install.sh | bash -s -- --local
```

**`--sudo` (по умолчанию):** Устанавливает в `/usr/local/bin/swag2mcp` с помощью `sudo`. Вам будет предложено ввести пароль. Если `sudo` не сработает, используется резервный путь `~/.local/bin/swag2mcp`.

**`--local`:** Устанавливает в `~/.local/bin/swag2mcp` без `sudo`. После установки добавьте в конфигурацию вашей оболочки:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

**Проверка:**

```bash
swag2mcp --version
```

### APT (Debian / Ubuntu)

1. Откройте [страницу последнего релиза](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Скачайте `swag2mcp_linux_amd64.deb`
3. Установите:

```bash
# Убедитесь, что вы в папке загрузок
pwd
sudo dpkg -i swag2mcp_linux_amd64.deb
```

**Проверка:**

```bash
swag2mcp --version
```

### RPM (Fedora / RHEL)

1. Откройте [страницу последнего релиза](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Скачайте `swag2mcp_linux_amd64.rpm`
3. Установите:

```bash
# Убедитесь, что вы в папке загрузок
pwd
sudo rpm -i swag2mcp_linux_amd64.rpm
```

**Проверка:**

```bash
swag2mcp --version
```

### Docker

Смонтируйте `~/.swag2mcp` (или ваш кастомный путь) в `/home/nonroot/.swag2mcp`.
Entrypoint автоматически настраивает права доступа, чтобы контейнер мог читать вашу рабочую директорию.

> **Пользователи Apple Silicon (M1/M2/M3/M4):** Добавьте `--platform linux/amd64` для запуска образа:
> `docker run --rm --platform linux/amd64 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ...`

```bash
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

Запуск с stdio-транспортом:

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp
```

Запуск с HTTP-транспортом:

```bash
docker run --rm -p 8080:8080 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp --transport sse --http-addr :8080
```

**Проверка:**

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest --version
```

### GitHub Release

1. Откройте [страницу последнего релиза](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Скачайте архив для вашей архитектуры:
   - **amd64**: `swag2mcp_linux_amd64.tar.gz`
   - **arm64**: `swag2mcp_linux_arm64.tar.gz`
3. Откройте Терминал в папке загрузок и выполните:

```bash
tar -xzf swag2mcp_linux_*.tar.gz
sudo mv swag2mcp /usr/local/bin/
```

**Проверка:**

```bash
swag2mcp --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

Убедитесь, что `$GOPATH/bin` находится в `$PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

**Проверка:**

```bash
swag2mcp --version
```

### Сборка из исходников

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

**Проверка:**

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

**Проверка:**

```powershell
swag2mcp --version
```

### GitHub Release

1. Откройте [страницу последнего релиза](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Скачайте `swag2mcp_windows_amd64.zip`
3. Распакуйте ZIP-файл (правой кнопкой мыши → Извлечь все, или через PowerShell)
4. Переместите `swag2mcp.exe` в `C:\Windows\System32\`

**Проверка:**

```powershell
swag2mcp --version
```

### go install

```powershell
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

**Проверка:**

```powershell
swag2mcp --version
```

### Сборка из исходников

```powershell
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp.exe ./cmd/swag2mcp
```

**Проверка:**

```powershell
swag2mcp --version
```

### Docker

Смонтируйте `~/.swag2mcp` (или ваш кастомный путь) в `/home/nonroot/.swag2mcp`. Entrypoint автоматически настраивает права доступа, чтобы контейнер мог читать вашу рабочую директорию.

```powershell
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

Run with stdio transport:

```powershell
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp
```

Run with HTTP transport:

```powershell
docker run --rm -p 8080:8080 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp --transport sse --http-addr :8080
```

Verify:

```powershell
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest --version
```

---

> Нужен mock-сервер? Смотрите [Установка Mock Server](install-mock.md).

## Установка через LLM-агента

Если вы используете IDE с ИИ (OpenCode, Cursor, Claude Desktop, VS Code и т.д.), перейдите в раздел [Скиллы для LLM](/ru/llm-skills/overview) для инструкций по установке скиллов через агента или вручную.

> Некоторые IDE и LLM-клиенты требуют перезагрузки после добавления скиллов.

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
