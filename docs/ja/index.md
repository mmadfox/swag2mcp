# swag2mcp

OpenAPI/Swagger/Postman の API 仕様を、Model Context Protocol (MCP) を介して LLM エージェントと連携させます。

すべての API が MCP に対応しているわけではありません — プライベートエンドポイント、内部サービス、レガシーシステム、サードパーティ API はほとんど対応していません。swag2mcp はあらゆる REST API を MCP インターフェースでラップし、サーバーコードを一行も変更することなく、LLM エージェントに API サーフェス全体への即時アクセスを提供します。ライブ API 呼び出しを通じて、LLM は現実世界の知識を獲得し、情報に基づいた意思決定、ワークフローの自動化、データに基づく行動を実行します — 単なる推測ではありません。

<a href="https://www.youtube.com/watch?v=1Da4UmE2f9U" target="_blank">
  <img src="https://raw.githubusercontent.com/mmadfox/swag2mcp/main/docs/cover.jpg" alt="プレビュー">
</a>

## あなたの API が LLM と対話

1 行の設定で、あらゆる OpenAPI/Swagger/Postman ファイルを MCP サーバーに変換します。LLM エージェントは API を発見、調査、呼び出し — 統合コードは一切不要です。

<img src="/architecture.svg" width="700" alt="swag2mcp アーキテクチャ">

## ラッパーを書くのはもう終わり

新しい API を LLM に接続するたびに、同じボイラープレート（仕様解析、認証、エラーハンドリング、レート制限）を書いていませんか？ swag2mcp が代わりに行います — 19 の既製 MCP ツールを提供。

## こんな方に

| 役割 | 理由 |
|------|------|
| **AI エージェント開発者** | 2 日ではなく 2 分で API を接続 |
| **MCP エンジニア** | ハンドラーコード不要 — 仕様を指定するだけ |
| **アーキテクト** | 全 LLM 向けの単一 API 統合レイヤー |
| **データアナリスト** | 自然言語で API にアクセス、コーディング不要 |
| **DevOps / SRE** | 追加サービスなしで LLM による監視と自動化 |
| **インテグレーター** | 9 つの認証方式を標準搭載 — Basic から OAuth2、HMAC まで |
| **QA エンジニア** | 実際の API を使わずにモックサーバーで分離テスト |
| **プロダクトマネージャー** | バックエンド作業なしで迅速な AI 機能プロトタイプ |
| **その他多数** | |

---

## ライセンス

**Apache License, Version 2.0** の下でライセンスされています。

全文は [LICENSE](https://github.com/mmadfox/swag2mcp/blob/main/LICENSE) をご覧ください。
