# Mockサーバーのインストール

`swag2mcp-mock` バイナリは、偽のレスポンスでAPIをテストするための独立したツールです。メインバイナリと同じインストール方法をサポートしています。

## 互換性

| 方法 | macOS | Linux | Windows |
|--------|-------|-------|---------|
| ワンライナー (curl) | ✅ | ✅ | ❌ |
| Homebrew | ✅ | ❌ | ❌ |
| Scoop | ❌ | ❌ | ✅ |
| GitHub Release | ✅ | ✅ | ✅ |
| go install | ✅ | ✅ | ✅ |
| ソースからビルド | ✅ | ✅ | ✅ |
| Docker | ✅ | ✅ | ✅ |

---

## macOS

### ワンライナー（推奨）

2つのインストールモード — 適した方を選択してください：

```bash
# sudo使用（/usr/local/binにインストール — 推奨）
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash

# sudoなし（~/.local/binにインストール）
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash -s -- --local
```

**`--sudo`（デフォルト）：** `sudo` を使用して `/usr/local/bin/swag2mcp-mock` にインストールします。パスワードの入力を求められます。`sudo` が失敗した場合は、`~/.local/bin/swag2mcp-mock` にフォールバックします。

**`--local`：** `sudo` なしで `~/.local/bin/swag2mcp-mock` にインストールします。インストール後、シェル設定に追加してください：

```bash
export PATH="$HOME/.local/bin:$PATH"
```

確認：

```bash
swag2mcp-mock --version
```

### Homebrew

```bash
brew install mmadfox/tap/swag2mcp-mock
```

確認：

```bash
swag2mcp-mock --version
```

### GitHub Release

1. [最新リリースページ](https://github.com/mmadfox/swag2mcp/releases/latest)を開く
2. お使いのMacに合ったアーカイブをダウンロード：
   - **Apple Silicon**：`swag2mcp-mock_darwin_arm64.tar.gz`
   - **Intel**：`swag2mcp-mock_darwin_amd64.tar.gz`
3. ダウンロードフォルダでターミナルを開き、実行：

```bash
tar -xzf swag2mcp-mock_darwin_*.tar.gz
sudo mv swag2mcp-mock /usr/local/bin/
```

確認：

```bash
swag2mcp-mock --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

`$GOPATH/bin` が `$PATH` に含まれていることを確認してください：

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

確認：

```bash
swag2mcp-mock --version
```

### ソースからビルド

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp-mock ./cmd/swag2mcp-mock
sudo mv swag2mcp-mock /usr/local/bin/
```

確認：

```bash
swag2mcp-mock --version
```

### Docker

コンテナは非rootユーザーとして実行され、ワークスペースディレクトリへのアクセスが必要です。
`~/.swag2mcp`（またはカスタムワークスペースパス）を `/home/nonroot/.swag2mcp` にマウントしてください：

> **Apple Silicon (M1/M2/M3/M4) ユーザー：** イメージを実行するには `--platform linux/amd64` を追加してください：
> `docker run --rm --platform linux/amd64 -v ~/.swag2mcp:/home/nonroot/.swag2mcp ...`

```bash
docker pull ghcr.io/mmadfox/swag2mcp-mock:latest
```

Mockサーバーを実行：

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest
```

確認：

```bash
docker run --rm -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest --version
```

---

## Linux

### ワンライナー（推奨）

2つのインストールモード — 適した方を選択してください：

```bash
# sudo使用（/usr/local/binにインストール — 推奨）
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash

# sudoなし（~/.local/binにインストール）
curl -fsSL https://raw.githubusercontent.com/mmadfox/swag2mcp/main/scripts/install-mock.sh | bash -s -- --local
```

**`--sudo`（デフォルト）：** `sudo` を使用して `/usr/local/bin/swag2mcp-mock` にインストールします。パスワードの入力を求められます。`sudo` が失敗した場合は、`~/.local/bin/swag2mcp-mock` にフォールバックします。

**`--local`：** `sudo` なしで `~/.local/bin/swag2mcp-mock` にインストールします。インストール後、シェル設定に追加してください：

```bash
export PATH="$HOME/.local/bin:$PATH"
```

確認：

```bash
swag2mcp-mock --version
```

### GitHub Release

1. [最新リリースページ](https://github.com/mmadfox/swag2mcp/releases/latest)を開く
2. お使いのアーキテクチャに合ったアーカイブをダウンロード：
   - **amd64**：`swag2mcp-mock_linux_amd64.tar.gz`
   - **arm64**：`swag2mcp-mock_linux_arm64.tar.gz`
3. ダウンロードフォルダでターミナルを開き、実行：

```bash
tar -xzf swag2mcp-mock_linux_*.tar.gz
sudo mv swag2mcp-mock /usr/local/bin/
```

確認：

```bash
swag2mcp-mock --version
```

### go install

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

`$GOPATH/bin` が `$PATH` に含まれていることを確認してください：

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

確認：

```bash
swag2mcp-mock --version
```

### ソースからビルド

```bash
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp-mock ./cmd/swag2mcp-mock
sudo mv swag2mcp-mock /usr/local/bin/
```

確認：

```bash
swag2mcp-mock --version
```

### Docker

コンテナは非rootユーザーとして実行され、ワークスペースディレクトリへのアクセスが必要です。
`~/.swag2mcp`（またはカスタムワークスペースパス）を `/home/nonroot/.swag2mcp` にマウントしてください：

```bash
docker pull ghcr.io/mmadfox/swag2mcp-mock:latest
```

Mockサーバーを実行：

```bash
docker run --rm -i -v ~/.swag2mcp:/home/nonroot/.swag2mcp ghcr.io/mmadfox/swag2mcp-mock:latest
```

確認：

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

確認：

```powershell
swag2mcp-mock --version
```

### GitHub Release

1. [最新リリースページ](https://github.com/mmadfox/swag2mcp/releases/latest)を開く
2. `swag2mcp-mock_windows_amd64.zip` をダウンロード
3. ZIPファイルを展開（右クリック → すべて展開、またはPowerShellを使用）
4. `swag2mcp-mock.exe` を `C:\Windows\System32\` に移動

確認：

```powershell
swag2mcp-mock --version
```

### go install

```powershell
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

確認：

```powershell
swag2mcp-mock --version
```

### ソースからビルド

```powershell
git clone https://github.com/mmadfox/swag2mcp.git
cd swag2mcp
go build -o swag2mcp-mock.exe ./cmd/swag2mcp-mock
```

確認：

```powershell
swag2mcp-mock --version
```

---

## 次のステップ

- [クイックスタート](quickstart.md) — 2分で始める
- [Mockサーバー](../advanced/mock-server.md) — Mockサーバーの設定と実行
