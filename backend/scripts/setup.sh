#!/bin/bash

# 研究室出席管理システム セットアップスクリプト

set -e

echo "=========================================="
echo "研究室出席管理システム セットアップ"
echo "=========================================="
echo ""

# 環境変数ファイルの確認
if [ ! -f .env ]; then
    echo "📝 .envファイルが見つかりません。.env.exampleからコピーします..."
    cp .env.example .env
    echo "✅ .envファイルを作成しました"
    echo "⚠️  .envファイルを編集して、適切な設定値を入力してください"
    echo ""
else
    echo "✅ .envファイルが存在します"
    echo ""
fi

# Dockerサービスの起動
echo "🐳 Dockerサービスを起動しています..."
docker-compose up -d postgres
echo ""

# PostgreSQLの起動を待つ
echo "⏳ PostgreSQLの起動を待っています..."
sleep 5

# データベースの接続確認
until docker-compose exec -T postgres pg_isready -U labuser -d lab_attendance > /dev/null 2>&1; do
    echo "   データベースの起動を待っています..."
    sleep 2
done
echo "✅ データベースが起動しました"
echo ""

# マイグレーションの実行
echo "🔄 データベースマイグレーションを実行しています..."
docker-compose --profile migration up migrate
echo "✅ マイグレーションが完了しました"
echo ""

# シードデータの投入確認
read -p "シードデータ（称号・設定）を投入しますか？ (y/N): " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🌱 シードデータを投入しています..."
    docker-compose exec -T postgres psql -U labuser -d lab_attendance < db/seeds/achievements.sql
    docker-compose exec -T postgres psql -U labuser -d lab_attendance < db/seeds/settings.sql
    echo "✅ シードデータの投入が完了しました"
    echo ""
fi

# バックエンドサーバーの起動
echo "🚀 バックエンドサーバーを起動しています..."
docker-compose up -d backend
echo ""

echo "=========================================="
echo "✅ セットアップが完了しました！"
echo "=========================================="
echo ""
echo "📊 サービスの状態を確認:"
docker-compose ps
echo ""
echo "🌐 アクセス先:"
echo "   - API: http://localhost:8080"
echo "   - Health Check: http://localhost:8080/health"
echo ""
echo "📝 次のステップ:"
echo "   1. .envファイルを編集して適切な設定値を入力"
echo "   2. ログを確認: make logs"
echo "   3. 開発を開始！"
echo ""

