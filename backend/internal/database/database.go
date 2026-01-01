package database

import (
	"fmt"
	"log"

	"github.com/kasa021/watabe-lab-app/internal/config"
	"github.com/kasa021/watabe-lab-app/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDatabase データベース接続を作成
func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	// ログレベルの設定
	logLevel := logger.Silent
	if cfg.Server.Env == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("データベース接続エラー: %w", err)
	}

	log.Println("データベース接続成功")
	return db, nil
}

// AutoMigrate テーブルを自動作成（開発用）
// 本番環境ではマイグレーションツールを使用すること
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.CheckInLog{},
		&domain.DailyAttendance{},
		&domain.Achievement{},
		&domain.UserAchievement{},
		&domain.Setting{},
	)
}

