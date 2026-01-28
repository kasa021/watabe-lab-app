# 研究室出席管理システム

研究室への登校頻度向上とメンバー間のコミュニケーション促進を目的とした出席管理システムです。

## プロジェクト概要

### 主な機能

- **チェックイン/チェックアウト機能**: WiFi SSIDおよびGPS座標を用いた位置情報の検証
- **リアルタイム在室状況**: WebSocketを使用した在室者情報のリアルタイム更新
- **ランキング機能**: 出席日数に基づく週次・月次・年次のランキング表示
- **称号システム**: 出席状況に応じた称号の付与（ゲーミフィケーション要素）
- **データ可視化**: 個人の滞在履歴やストリーク（連続出席日数）の表示

### 技術スタック

#### バックエンド

- **言語**: Go 1.21+
- **フレームワーク**: Gin (REST API), gorilla/websocket
- **ORM**: GORM
- **データベース**: PostgreSQL 15+
- **認証**: LDAP + JWT

#### フロントエンド

- **フレームワーク**: React + TypeScript
- **ビルドツール**: Vite
- **状態管理**: Zustand
- **UI**: Tailwind CSS
- **HTTPクライアント**: Axios

#### インフラ

- **開発環境**: Docker Compose
- **本番環境**: 研究室サーバー

## セットアップ手順

### 前提条件

以下のツールがインストールされていること:

- Docker & Docker Compose
- Node.js 18以上
- Go 1.21以上（ローカル開発を行う場合）
- Make（オプション）

### 1. リポジトリのクローン

```bash
git clone https://github.com/kasa021/watabe-lab-app.git
cd watabe-lab-app
```

### 2. 環境変数の設定

バックエンドとフロントエンドそれぞれの環境変数ファイルを作成します。

```bash
# バックエンド
cd backend
cp .env.example .env
# 必要に応じて.envを編集してください

# フロントエンド
cd ../frontend
cp .env.example .env
```

### 3. バックエンドの起動

#### 推奨: セットアップスクリプトの使用

```bash
cd backend
chmod +x scripts/setup.sh
./scripts/setup.sh
```

#### 手動セットアップの場合

```bash
cd backend

# Dockerサービスの起動
docker-compose up -d postgres

# マイグレーションとシードデータの適用
make migrate-up
make seed

# サーバー起動
docker-compose up -d backend

# ログ確認
make logs
```

### 4. フロントエンドの起動

```bash
cd frontend

# 依存関係のインストール
npm install

# 開発サーバー起動
npm run dev
```

## ライセンスと開発

本システムは渡部研究室内部向けに開発されています。

**開発者**: 渡部研究室
**作成**: 2025-12
