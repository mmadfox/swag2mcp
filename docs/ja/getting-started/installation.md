# インストール

## システム要件

- **macOS、Linux、Windows**（amd64 / arm64）
- **Go 1.26+**（`go install` またはソースからのビルドの場合のみ）

## 互換性

| 方法 | macOS | Linux | Windows |
|------|-------|-------|---------|
| One-liner (curl) | ✅ | ✅ | ❌ |
| Homebrew | ✅ | ❌ | ❌ |
| Scoop | ❌ | ❌ | ✅ |
| APT (deb) | ❌ | ✅ | ❌ |
| RPM | ❌ | ✅ | ❌ |
| Docker | ✅ | ✅ | ✅ |
| GitHub Release | ✅ | ✅ | ✅ |
| go install | ✅ | ✅ | ✅ |
| ソースからビルド | ✅ | ✅ | ✅ |

---

## macOS

### One-liner（推奨）

2つのインストールモード — 適した方を選んでください：

```bash
# sudo を使用（/usr/local/bin にインストール — 推奨）
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash

# sudo なし（~/.local/bin にインストール）
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash -s -- --local
```

**`--sudo`（デフォルト）：** `sudo` を使用して `/usr/local/bin/swag2mcp` にインストールします。パスワードの入力を求められます。`sudo` が失敗した場合は `~/.local/bin/swag2mcp` にフォールバックします。

**`--local`：** `sudo` なしで `~/.local/bin/swag2mcp` にインストールします。インストール後、シェル設定に追加してください：

```bash
export PATH="$HOME/.local/bin:$PATH"
```

**確認：**

```bash
swag2mcp --version
```

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp
```

**確認：**

```bash
swag2mcp --version
```

### GitHub Release

1. [最新リリースページ](https://github.com/mmadfox/swag2mcp/releases/latest)を開く
2. お使いのMacに合ったアーカイブをダウンロード：
   - **Apple Silicon**：`swag2mcp_darwin_arm64.tar.gz`
   - **Intel**：`swag2mcp_darwin_amd64.tar.gz`
3. ダウンロードフォルダでターミナルを開き、以下を実行：

```bash
tar -xzf swag2mcp_darwin_*.tar.gz
sudo mv swag2mcp /usr/local/bin/
```

**確認：**

```bash
swag2mcp --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

`$GOPATH/bin` が `$PATH` に含まれていることを確認してください：

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

**確認：**

```bash
swag2mcp --version
```

### ソースからビルド

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

**確認：**

```bash
swag2mcp --version
```

### Docker

コンテナは非 root ユーザーとして実行され、ワークスペースディレクトリへのアクセスが必要です。`~/.swag2mcp`（またはカスタムパス）を `/home/nonroot/.swag2mcp` にマウントしてください：

> **Apple Silicon（M1/M2/M3/M4）ユーザー：** イメージを実行するには `--platform linux/amd64` を追加してください：
> `docker run --rm --platform linux/amd64 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ...`

```bash
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

stdio トランスポートで実行：

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp
```

HTTP トランスポートで実行：

```bash
docker run --rm -p 8080:8080 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp --transport sse --http-addr :8080
```

**確認：**

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest --version
```

---

## Linux

### One-liner（推奨）

2つのインストールモード — 適した方を選んでください：

```bash
# sudo を使用（/usr/local/bin にインストール — 推奨）
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash

# sudo なし（~/.local/bin にインストール）
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install.sh | bash -s -- --local
```

**`--sudo`（デフォルト）：** `sudo` を使用して `/usr/local/bin/swag2mcp` にインストールします。パスワードの入力を求められます。`sudo` が失敗した場合は `~/.local/bin/swag2mcp` にフォールバックします。

**`--local`：** `sudo` なしで `~/.local/bin/swag2mcp` にインストールします。インストール後、シェル設定に追加してください：

```bash
export PATH="$HOME/.local/bin:$PATH"
```

**確認：**

```bash
swag2mcp --version
```

### APT（Debian / Ubuntu）

1. [最新リリースページ](https://github.com/mmadfox/swag2mcp/releases/latest)を開く
2. `swag2mcp_linux_amd64.deb` をダウンロード
3. インストール：

```bash
# ダウンロードフォルダにいることを確認してください
pwd
sudo dpkg -i swag2mcp_linux_amd64.deb
```

**確認：**

```bash
swag2mcp --version
```

### RPM（Fedora / RHEL）

1. [最新リリースページ](https://github.com/mmadfox/swag2mcp/releases/latest)を開く
2. `swag2mcp_linux_amd64.rpm` をダウンロード
3. インストール：

```bash
# ダウンロードフォルダにいることを確認してください
pwd
sudo rpm -i swag2mcp_linux_amd64.rpm
```

**確認：**

```bash
swag2mcp --version
```

### Docker

コンテナは非 root ユーザーとして実行され、ワークスペースディレクトリへのアクセスが必要です。`~/.swag2mcp`（またはカスタムパス）を `/home/nonroot/.swag2mcp` にマウントしてください：

> **Apple Silicon（M1/M2/M3/M4）ユーザー：** イメージを実行するには `--platform linux/amd64` を追加してください：
> `docker run --rm --platform linux/amd64 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ...`

```bash
docker pull ghcr.io/mmadfox/swag2mcp:latest
```

stdio トランスポートで実行：

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp
```

HTTP トランスポートで実行：

```bash
docker run --rm -p 8080:8080 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest mcp --transport sse --http-addr :8080
```

**確認：**

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp:latest --version
```

### GitHub Release

1. [最新リリースページ](https://github.com/mmadfox/swag2mcp/releases/latest)を開く
2. お使いのアーキテクチャに合ったアーカイブをダウンロード：
   - **amd64**：`swag2mcp_linux_amd64.tar.gz`
   - **arm64**：`swag2mcp_linux_arm64.tar.gz`
3. ダウンロードフォルダでターミナルを開き、以下を実行：

```bash
tar -xzf swag2mcp_linux_*.tar.gz
sudo mv swag2mcp /usr/local/bin/
```

**確認：**

```bash
swag2mcp --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

`$GOPATH/bin` が `$PATH` に含まれていることを確認してください：

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

**確認：**

```bash
swag2mcp --version
```

### ソースからビルド

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
make build
sudo mv swag2mcp /usr/local/bin/
```

**確認：**

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

**確認：**

```powershell
swag2mcp --version
```

### GitHub Release

1. [最新リリースページ](https://github.com/mmadfox/swag2mcp/releases/latest)を開く
2. `swag2mcp_windows_amd64.zip` をダウンロード
3. ZIPファイルを展開（右クリック → すべて展開、またはPowerShellを使用）
4. `swag2mcp.exe` を `C:\Windows\System32\` に移動

**確認：**

```powershell
swag2mcp --version
```

### go install

```powershell
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

**確認：**

```powershell
swag2mcp --version
```

### ソースからビルド

```powershell
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp.exe ./cmd/swag2mcp
```

**確認：**

```powershell
swag2mcp --version
```

---

> モックサーバーが必要ですか？[Mock Server のインストール](install-mock.md) を参照してください。

## LLM エージェント経由のインストール

AI 搭載 IDE（OpenCode、Cursor、Claude Desktop、VS Code など）を使用している場合、エージェントを通じて swag2mcp をインストールできます：

1. エージェントに swag2mcp スキルを追加するよう依頼します：

   ```
   ".agents/skills/swag2mcp-cli ディレクトリを作成し、https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md のスキルを .agents/skills/swag2mcp-cli/SKILL.md に追加してください"
   ".agents/skills/swag2mcp-format ディレクトリを作成し、https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md のスキルを .agents/skills/swag2mcp-format/SKILL.md に追加してください"
   ```

2. 次にエージェントに指示します：

   ```
   "swag2mcp をセットアップしてください"
   ```

   エージェントが swag2mcp をダウンロードしてインストールし、使用可能なスペックを含むワークスペースを作成します。

> 一部の IDE ではスキル追加後に再起動が必要です。

---

## 確認

```bash
swag2mcp --version
```

期待される出力（バージョンは異なる場合があります）：

```
swag2mcp v*.*.*
```

---

## 次のステップ

- [クイックスタート](quickstart.md) — 2分で使い始める
