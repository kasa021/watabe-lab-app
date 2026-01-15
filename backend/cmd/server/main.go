package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kasa021/watabe-lab-app/internal/config"
	"github.com/kasa021/watabe-lab-app/internal/database"
	"github.com/kasa021/watabe-lab-app/internal/handler"
	"github.com/kasa021/watabe-lab-app/internal/middleware"
	"github.com/kasa021/watabe-lab-app/internal/repository"
	"github.com/kasa021/watabe-lab-app/internal/service"
	"github.com/kasa021/watabe-lab-app/internal/ws"
)

func main() {
	// 環境変数の読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// 設定の読み込み
	cfg := config.Load()
	log.Printf("Server starting in %s mode", cfg.Server.Env)

	// データベース接続
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// リポジトリの初期化
	userRepo := repository.NewUserRepository(db)

	// サービスの初期化
	authService := service.NewAuthService(cfg)

	// ハンドラーの初期化
	authHandler := handler.NewAuthHandler(authService, userRepo)

	// WebSocket Hubの初期化と起動
	hub := ws.NewHub()
	go hub.Run()

	// 出席管理機能の初期化
	attendanceRepo := repository.NewAttendanceRepository(db)
	attendanceService := service.NewAttendanceService(attendanceRepo, hub)
	attendanceHandler := handler.NewAttendanceHandler(attendanceService)

	// Ginエンジンの作成
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// CORS設定
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:5173", "http://localhost:3000"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	// ヘルスチェックエンドポイント
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Lab Attendance System is running",
		})
	})

	// APIルーティング
	api := r.Group("/api/v1")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		// WebSocket エンドポイント
		api.GET("/ws", func(c *gin.Context) {
			ws.ServeWs(hub, c)
		})

		// 認証エンドポイント（認証不要）
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}

		// 認証が必要なエンドポイント
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			protected.GET("/auth/me", authHandler.Me)

			// 出席管理エンドポイント
			attendance := protected.Group("/attendance")
			{
				attendance.POST("/checkin", attendanceHandler.CheckIn)
				attendance.POST("/checkout", attendanceHandler.CheckOut)
				attendance.GET("/active", attendanceHandler.GetActiveUsers)
			}

			// 管理者のみアクセス可能なエンドポイント
			admin := protected.Group("")
			admin.Use(middleware.RoleMiddleware("admin"))
			{
				// TODO: 管理者用エンドポイント
			}
		}
	}

	// サーバー起動
	log.Printf("Starting server on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
