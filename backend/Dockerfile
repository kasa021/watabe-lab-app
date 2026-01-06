# ビルドステージ
FROM golang:1.24-alpine AS builder

WORKDIR /app

# 依存関係のコピーとダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .

# バイナリのビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# 実行ステージ
FROM alpine:latest

WORKDIR /app

# タイムゾーンとCA証明書のインストール
RUN apk --no-cache add ca-certificates tzdata

# ビルドしたバイナリのコピー
COPY --from=builder /app/server .

# ポートの公開
EXPOSE 8080

# 実行
CMD ["./server"]

