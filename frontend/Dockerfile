# ビルドステージ
FROM node:20-alpine AS builder

WORKDIR /app

# 依存関係のインストール
COPY package*.json ./
RUN npm ci

# ソースコードのコピーとビルド
COPY . .
RUN npm run build

# 実行ステージ
FROM nginx:alpine

# Nginx設定のコピー
COPY nginx.conf /etc/nginx/conf.d/default.conf

# ビルド成果物のコピー
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
