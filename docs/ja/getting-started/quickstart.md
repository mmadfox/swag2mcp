# クイックスタート

swag2mcp を 2 分で起動する。

## 1. 初期化

### ホームディレクトリ（推奨）

システム全体で一度だけセットアップします。設定はホームフォルダに保存されます。

::: code-group

```bash [macOS / Linux]
swag2mcp init
# ~/.swag2mcp/swag2mcp.yaml が作成される
```

```powershell [Windows]
swag2mcp.exe init
# %USERPROFILE%\.swag2mcp\swag2mcp.yaml が作成される
```

:::

### プロジェクトディレクトリ

プロジェクト内で独立したワークスペースを使用する場合。

::: code-group

```bash [macOS / Linux]
mkdir -p ./swag2mcp && swag2mcp init ./swag2mcp
```

```powershell [Windows]
mkdir ./swag2mcp; swag2mcp.exe init ./swag2mcp
```

:::

### ZIP から

既製のワークスペースがある場合（例：同僚から受け取った場合）：

```bash
swag2mcp import --from-zip workspace.zip
```

## 2. エージェントスキルのインストール（推奨）

swag2mcp のスキルをインストールして、AI エージェントにすべてのコマンド、フラグ、設定形式、実際の使用例を教えます。

エージェントに以下のように依頼してください：

```bash
"swag2mcp-cli スキルを https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md から追加して"
"swag2mcp-format スキルを https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md から追加して"
```

> 一部の IDE ではスキル追加後に再起動が必要です。

## 3. LLM クライアント / IDE の設定

IDE を swag2mcp に接続するように設定します。IDE は必要に応じて MCP サーバーを自動的に起動します。

::: code-group

```json [OpenCode]
{
  "mcp": {
    "swag2mcp": {
      "type": "local",
      "command": ["swag2mcp", "mcp"],
      "enabled": true
    }
  }
}
```

```json [Claude Desktop]
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp"]
    }
  }
}
```

```json [Crush]
{
  "mcp": {
    "swag2mcp": {
      "type": "stdio",
      "command": "swag2mcp",
      "args": ["mcp"]
    }
  }
}
```

:::

他の IDE（Cursor、VS Code、JetBrains）については[統合ガイド](../integration/opencode.md)を参照してください。

> カスタムパス（例：`./swag2mcp`）でワークスペースを初期化した場合は、コマンドにフルパスを指定してください：
> `"command": ["swag2mcp", "mcp", "/absolute/path/to/swag2mcp"]`

> **設定変更後は MCP サーバーを再起動**してください。変更が反映されます。

## 4. MCP サーバーの起動

### stdio（デフォルト）— ローカル IDE 用

設定は不要です。上記の設定により、IDE が自動的に swag2mcp を起動します。

```bash
swag2mcp mcp
```

### SSE / Streamable HTTP — リモートアクセス用

```bash
swag2mcp mcp --transport sse --http-addr :8080
```

または `swag2mcp.yaml` で設定：

```yaml
mcp:
  transport: sse
  addr: ":8080"
  path: "/mcp"
```

すべてのフラグについては [MCP サーバーリファレンス](../configuration/mcp-server.md) を参照してください。

### タグでスペックをフィルタリング

```bash
swag2mcp mcp --tags weather,public
```

一致するタグを持つスペックのみが LLM から利用可能になります。

### 動作確認

接続後、LLM エージェントに以下のように尋ねてください：

```bash
"どの MCP ツールをサポートしていますか？"
```

エージェントが swag2mcp のツール（`spec_list`、`search`、`invoke` など）をリスト表示すれば、正常に動作しています。

### 試せるクエリ例

| エージェントへの依頼 | 動作 |
|-------|-------------|
| "ニューヨークの天気は？" | `invoke` — Open-Meteo 予報 API を呼び出し |
| "現在の BTC 価格は？" | `invoke` — Binance ティッカー API を呼び出し |
| "ダジャレを教えて" | `invoke` — icanhazdadjoke API を呼び出し |
| "ピカチュウを見せて" | `invoke` — PokéAPI を名前で呼び出し |
| "Rick Sanchez は誰？" | `invoke` — Rick and Morty キャラクター API を呼び出し |
| "北京の大気質は？" | `invoke` — Open-Meteo 大気質 API を呼び出し |
| "ポルトガル近くの波の高さは？" | `invoke` — Open-Meteo 海洋 API を呼び出し |
| "犬に関するジョークを検索" | `invoke` — dadjoke 検索エンドポイントを呼び出し |
| "すべてのポケモンをリスト表示" | `invoke` — PokéAPI リストエンドポイントを呼び出し |
| "エベレストの標高は？" | `invoke` — Open-Meteo 標高 API を呼び出し |

## 5. 次のステップ

- [コンセプト](../concepts/overview.md) — アーキテクチャを理解する
- [設定](../configuration/config-file.md) — 設定をカスタマイズする
- [CLI コマンド](../cli/overview.md) — 全コマンドリファレンス
