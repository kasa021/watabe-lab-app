.PHONY: help up down logs restart build clean test lint migrate-up migrate-down run-local

help: ## ヘルプを表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

up: ## サービスを起動
	docker-compose up -d

down: ## サービスを停止
	docker-compose down

logs: ## ログを表示
	docker-compose logs -f

restart: ## サービスを再起動
	docker-compose restart

build: ## イメージをビルド
	docker-compose build

clean: ## すべてをクリーンアップ（データも削除）
	docker-compose down -v
	rm -rf tmp/

test: ## テストを実行
	go test -v ./...

lint: ## リンターを実行
	golangci-lint run

migrate-up: ## マイグレーションを実行
	docker-compose --profile migration up migrate

migrate-down: ## マイグレーションをロールバック
	docker-compose run --rm migrate -path /migrations -database "postgres://labuser:labpass@postgres:5432/lab_attendance?sslmode=disable" down

run-local: ## ローカルで実行（Dockerを使わない）
	go run cmd/server/main.go

seed: ## シードデータを投入
	docker-compose exec -T postgres psql -U labuser -d lab_attendance < db/seeds/achievements.sql
	docker-compose exec -T postgres psql -U labuser -d lab_attendance < db/seeds/settings.sql

