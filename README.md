# 研究室出席管理システム

研究室への登校頻度向上とコミュニケーション促進を目的とした出席管理システムです。

## 📋 プロジェクト概要

### 主な機能

- 📝 **チェックイン/チェックアウト機能**: WiFi/GPS検証による位置確認
- 👥 **リアルタイム在室者表示**: WebSocketによるリアルタイム更新
- 🏆 **ランキング機能**: 週次・月次・年次ランキング
- 🎖️ **称号システム**: ゲーミフィケーション要素
- 📊 **データ可視化**: 個人の滞在履歴とストリーク表示

### 技術スタック

#### バックエンド
- **言語**: Go 1.21+
- **フレームワーク**: Gin (REST API), gorilla/websocket
- **ORM**: GORM
- **データベース**: PostgreSQL 15+
- **認証**: LDAP + JWT

#### フロントエンド
- **フレームワーク**: React 18 + TypeScript
- **ビルドツール**: Vite
- **状態管理**: Zustand
- **UI**: Tailwind CSS
- **HTTPクライアント**: Axios

#### インフラ
- **開発環境**: Docker Compose
- **本番環境**: 研究室サーバー

## 🚀 クイックスタート

### 前提条件

以下のツールがインストールされている必要があります：

- Docker & Docker Compose
- Node.js 18以上
- Go 1.21以上（ローカル開発時）
- Make（オプション）

### 1. リポジトリのクローン

```bash
git clone https://github.com/kasa021/watabe-lab-app.git
cd watabe-lab-app
```

### 2. 環境変数の設定

```bash
# バックエンドの環境変数
cd backend
cp .env.example .env
# .envファイルを編集して適切な値を設定

# フロントエンドの環境変数
cd ../frontend
cp .env.example .env
```

### 3. バックエンドのセットアップと起動

#### 方法1: セットアップスクリプトを使用（推奨）

```bash
cd backend
chmod +x scripts/setup.sh
./scripts/setup.sh
```

#### 方法2: 手動セットアップ

```bash
cd backend

# Dockerサービスの起動
docker-compose up -d postgres

# マイグレーションの実行
make migrate-up

# シードデータの投入
make seed

# バックエンドサーバーの起動
docker-compose up -d backend

# ログの確認
make logs
```

### 4. フロントエンドの起動

```bash
cd frontend

# 依存関係のインストール
npm install

# 開発サーバーの起動
npm run dev
```

### 5. アクセス

- **フロントエンド**: http://localhost:5173
- **バックエンドAPI**: http://localhost:8080
- **ヘルスチェック**: http://localhost:8080/health

## 📁 プロジェクト構造

```
watabe-lab-app/
├── backend/                    # バックエンド（Go）
│   ├── cmd/server/            # アプリケーションエントリーポイント
│   ├── internal/              # 内部パッケージ
│   │   ├── config/           # 設定管理
│   │   ├── domain/           # ドメインモデル
│   │   ├── repository/       # データアクセス層
│   │   ├── service/          # ビジネスロジック
│   │   ├── handler/          # HTTPハンドラー
│   │   ├── middleware/       # ミドルウェア
│   │   └── ws/               # WebSocket管理
│   ├── db/
│   │   ├── migrations/       # データベースマイグレーション
│   │   └── seeds/            # 初期データ
│   ├── scripts/              # セットアップスクリプト
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   └── Makefile
├── frontend/                   # フロントエンド（React + TypeScript）
│   ├── src/
│   │   ├── components/       # 再利用可能なコンポーネント
│   │   ├── pages/            # ページコンポーネント
│   │   ├── api/              # APIクライアント
│   │   ├── types/            # TypeScript型定義
│   │   ├── stores/           # Zustand ストア
│   │   ├── App.tsx
│   │   └── main.tsx
│   ├── index.html
│   ├── package.json
│   ├── tsconfig.json
│   ├── vite.config.ts
│   └── tailwind.config.js
├── docker-compose.yml         # Docker Compose設定
├── .gitignore
└── README.md
```

## 🛠️ 開発コマンド

### バックエンド（backend/）

```bash
make help          # コマンド一覧を表示
make up            # サービスを起動
make down          # サービスを停止
make logs          # ログを表示
make restart       # サービスを再起動
make build         # イメージをビルド
make clean         # すべてをクリーンアップ
make test          # テストを実行
make lint          # リンターを実行
make migrate-up    # マイグレーションを実行
make migrate-down  # マイグレーションをロールバック
make seed          # シードデータを投入
make run-local     # ローカルで実行（Dockerを使わない）
```

### フロントエンド（frontend/）

```bash
npm run dev        # 開発サーバーを起動
npm run build      # プロダクションビルド
npm run preview    # ビルドしたアプリをプレビュー
npm run lint       # リンターを実行
```

## 📚 ドキュメント

詳細なドキュメントは `prompt/` ディレクトリにあります：

- [要件定義書](./prompt/要件定義.md)
- [セットアップ手順書](./prompt/setup.md)
- [開発規約](./.cursorrules)

## 🗄️ データベース設計

### 主要テーブル

- **users**: ユーザー情報
- **check_in_logs**: チェックインログ
- **daily_attendances**: 日次出席記録
- **achievements**: 称号マスタ
- **user_achievements**: ユーザー獲得称号
- **settings**: システム設定

詳細は[要件定義書](./prompt/要件定義.md)を参照してください。

## 🔐 認証

- LDAP認証による既存システムとの連携
- JWTトークンによるセッション管理
- WiFi SSID + GPS座標による位置検証

## 🎮 ゲーミフィケーション要素

### ストリーク機能
- 連続出席日数の記録
- 土日祝日を除外した計算

### ランキング
- 週次・月次・年次ランキング
- 出席日数を基準とした評価

### 称号システム
- 早起き系、夜型系、ストリーク系など多様な称号
- 達成時のボーナスポイント

## 🐛 トラブルシューティング

### ポートが既に使用されている

```bash
# ポート使用状況の確認
lsof -i :8080
lsof -i :5432

# プロセスを終了
kill -9 <PID>
```

### データベース接続エラー

```bash
# PostgreSQLコンテナの状態確認
docker-compose ps

# PostgreSQLのログ確認
docker-compose logs postgres

# コンテナの再起動
docker-compose restart postgres
```

### マイグレーションエラー

```bash
# マイグレーションの状態確認
docker-compose exec postgres psql -U labuser -d lab_attendance -c "SELECT * FROM schema_migrations;"

# マイグレーションを強制的にリセット（開発環境のみ）
make clean
make up
make migrate-up
```

## 📝 開発ワークフロー

### ブランチ戦略

```bash
# 新機能の開発
git checkout -b feature/user-authentication
# ... 開発作業 ...
git add .
git commit -m "feat: Implement LDAP authentication"
git push origin feature/user-authentication
# GitHubでPull Requestを作成

# バグ修正
git checkout -b fix/login-error
# ... 修正作業 ...
git commit -m "fix: Fix login validation error"
```

### コミットメッセージ規約

```
<type>: <subject>

Type:
- feat: 新機能
- fix: バグ修正
- docs: ドキュメント
- style: コードフォーマット
- refactor: リファクタリング
- test: テスト
- chore: その他の変更
```

## 🤝 コントリビューション

1. このリポジトリをフォーク
2. 新しいブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'feat: Add amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. Pull Requestを作成

## 📄 ライセンス

このプロジェクトは研究室内部での使用を目的としています。

## 👥 開発者

- 渡部研究室

## 📞 サポート

問題が発生した場合は、GitHubのIssuesで報告してください。

---

**作成日**: 2024-12-27  
**バージョン**: 1.0.0

