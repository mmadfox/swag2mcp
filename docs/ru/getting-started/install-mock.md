# Установка Mock-сервера

Бинарный файл `swag2mcp-mock` — это отдельный инструмент для тестирования API с фиктивными ответами. Он поддерживает те же методы установки, что и основной бинарный файл.

## Совместимость

| Метод | macOS | Linux | Windows |
|--------|-------|-------|---------|
| Одной строкой (curl) | ✅ | ✅ | ❌ |
| Homebrew | ✅ | ❌ | ❌ |
| Scoop | ❌ | ❌ | ✅ |
| GitHub Release | ✅ | ✅ | ✅ |
| go install | ✅ | ✅ | ✅ |
| Сборка из исходников | ✅ | ✅ | ✅ |
| Docker | ✅ | ✅ | ✅ |

---

## macOS

### Одной строкой (рекомендуется)

Два режима установки — выберите подходящий:

```bash
# С sudo (устанавливается в /usr/local/bin — рекомендуется)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash

# Без sudo (устанавливается в ~/.local/bin)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash -s -- --local
```

**`--sudo` (по умолчанию):** Устанавливает в `/usr/local/bin/swag2mcp-mock` с помощью `sudo`. Будет запрошен пароль. Если `sudo` не сработает, выполняется откат к `~/.local/bin/swag2mcp-mock`.

**`--local`:** Устанавливает в `~/.local/bin/swag2mcp-mock` без `sudo`. После установки добавьте в конфигурацию вашей оболочки:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Проверка:

```bash
swag2mcp-mock --version
```

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp-mock
```

Проверка:

```bash
swag2mcp-mock --version
```

### GitHub Release

1. Откройте [страницу последнего релиза](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Скачайте архив для вашего Mac:
   - **Apple Silicon**: `swag2mcp-mock_darwin_arm64.tar.gz`
   - **Intel**: `swag2mcp-mock_darwin_amd64.tar.gz`
3. Откройте Терминал в папке загрузок и выполните:

```bash
tar -xzf swag2mcp-mock_darwin_*.tar.gz
sudo mv swag2mcp-mock /usr/local/bin/
```

Проверка:

```bash
swag2mcp-mock --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

Убедитесь, что `$GOPATH/bin` находится в вашем `$PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Проверка:

```bash
swag2mcp-mock --version
```

### Сборка из исходников

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp-mock ./cmd/swag2mcp-mock
sudo mv swag2mcp-mock /usr/local/bin/
```

Проверка:

```bash
swag2mcp-mock --version
```

### Docker

Подключите `~/.swag2mcp` (или ваш собственный путь к рабочему каталогу) к `/home/nonroot/.swag2mcp`.
Точка входа автоматически настраивает права доступа к файлам, чтобы контейнер мог читать ваш рабочий каталог.

```bash
docker pull ghcr.io/mmadfox/swag2mcp-mock:latest
```

Запуск mock-сервера:

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest
```

Проверка:

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest --version
```

---

## Linux

### Одной строкой (рекомендуется)

Два режима установки — выберите подходящий:

```bash
# С sudo (устанавливается в /usr/local/bin — рекомендуется)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash

# Без sudo (устанавливается в ~/.local/bin)
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash -s -- --local
```

**`--sudo` (по умолчанию):** Устанавливает в `/usr/local/bin/swag2mcp-mock` с помощью `sudo`. Будет запрошен пароль. Если `sudo` не сработает, выполняется откат к `~/.local/bin/swag2mcp-mock`.

**`--local`:** Устанавливает в `~/.local/bin/swag2mcp-mock` без `sudo`. После установки добавьте в конфигурацию вашей оболочки:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Проверка:

```bash
swag2mcp-mock --version
```

### GitHub Release

1. Откройте [страницу последнего релиза](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Скачайте архив для вашей архитектуры:
   - **amd64**: `swag2mcp-mock_linux_amd64.tar.gz`
   - **arm64**: `swag2mcp-mock_linux_arm64.tar.gz`
3. Откройте Терминал в папке загрузок и выполните:

```bash
tar -xzf swag2mcp-mock_linux_*.tar.gz
sudo mv swag2mcp-mock /usr/local/bin/
```

Проверка:

```bash
swag2mcp-mock --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

Убедитесь, что `$GOPATH/bin` находится в вашем `$PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Проверка:

```bash
swag2mcp-mock --version
```

### Сборка из исходников

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp-mock ./cmd/swag2mcp-mock
sudo mv swag2mcp-mock /usr/local/bin/
```

Проверка:

```bash
swag2mcp-mock --version
```

### Docker

Подключите `~/.swag2mcp` (или ваш собственный путь к рабочему каталогу) к `/home/nonroot/.swag2mcp`.
Точка входа автоматически настраивает права доступа к файлам, чтобы контейнер мог читать ваш рабочий каталог.

```bash
docker pull ghcr.io/mmadfox/swag2mcp-mock:latest
```

Запуск mock-сервера:

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest
```

Проверка:

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

Проверка:

```powershell
swag2mcp-mock --version
```

### GitHub Release

1. Откройте [страницу последнего релиза](https://github.com/mmadfox/swag2mcp/releases/latest)
2. Скачайте `swag2mcp-mock_windows_amd64.zip`
3. Распакуйте ZIP-файл (правой кнопкой мыши → Извлечь всё, или через PowerShell)
4. Переместите `swag2mcp-mock.exe` в `C:\Windows\System32\`

Проверка:

```powershell
swag2mcp-mock --version
```

### go install

```powershell
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

Проверка:

```powershell
swag2mcp-mock --version
```

### Сборка из исходников

```powershell
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp-mock.exe ./cmd/swag2mcp-mock
```

Проверка:

```powershell
swag2mcp-mock --version
```

---

## Следующие шаги

- [Быстрый старт](quickstart.md) — начните работу за 2 минуты
- [Mock-сервер](../advanced/mock-server.md) — настройка и запуск mock-сервера
