# プロジェクト構造

```
swag2mcp/
├── cmd/
│   ├── swag2mcp/          # メインバイナリ
│   │   └── main.go
│   └── swag2mcp-mock/     # モックサーバー
│       └── main.go
├── internal/
│   ├── auth/              # 9 つの認証方法
│   ├── cache/             # スペックキャッシュ
│   ├── commands/          # 13 の CLI コマンド（cobra）
│   ├── config/            # YAML 設定
│   ├── env/               # 環境変数
│   ├── httpclient/        # HTTP クライアント
│   ├── id/                # MD5 ID 生成
│   ├── index/             # 全文検索（bluge）
│   ├── model/             # データモデル
│   ├── reader/            # 大規模レスポンス読み取り
│   ├── server/
│   │   ├── mcp/           # MCP サーバー（19 ツール）
│   │   └── mockserver/    # モックサーバー
│   ├── service/           # ビジネスロジック
│   ├── spec/              # スペックパーサー
│   ├── tui/               # TUI インターフェース
│   └── workspace/         # ワークスペース管理
├── specs/                 # サンプルスペック
├── tests/                 # 統合テスト
├── docs/                  # ドキュメント
├── examples/              # 設定例
└── playground/            # 開発サンドボックス
```

## 主要パッケージ

| パッケージ | 説明 |
|---------|-------------|
| `auth` | 9 つの認証方法 |
| `cache` | TTL 付きディスクベースのキャッシュ |
| `commands` | Cobra CLI コマンド |
| `config` | カスケード付き YAML 設定 |
| `httpclient` | 設定可能な HTTP クライアント |
| `index` | 全文検索（bluge） |
| `server/mcp` | MCP サーバー（3 つのトランスポート） |
| `service` | ビジネスロジック（コア） |
| `spec` | OpenAPI/Swagger/Postman パーサー |
| `tui` | Bubbletea TUI |
| `workspace` | ファイル管理 |
