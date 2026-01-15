package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/kasa021/watabe-lab-app/internal/config"
	"github.com/kasa021/watabe-lab-app/internal/database"
	"github.com/kasa021/watabe-lab-app/internal/domain"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env file not found")
	}

	cfg := config.Load()
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	seedAchievements(db)
}

func seedAchievements(db *gorm.DB) {
	achievements := []domain.Achievement{
		{
			Code:           "check_in_1",
			Name:           "初めての出席",
			Description:    "初めて研究室に出席しました。",
			ConditionType:  "check_in_count",
			ConditionValue: map[string]interface{}{"target": 1},
			PointsReward:   10,
			Category:       "attendance",
			DisplayOrder:   1,
		},
		{
			Code:           "check_in_10",
			Name:           "常連さん",
			Description:    "研究室に10回出席しました。",
			ConditionType:  "check_in_count",
			ConditionValue: map[string]interface{}{"target": 10},
			PointsReward:   50,
			Category:       "attendance",
			DisplayOrder:   2,
		},
		{
			Code:           "check_in_50",
			Name:           "研究室の主",
			Description:    "研究室に50回出席しました。",
			ConditionType:  "check_in_count",
			ConditionValue: map[string]interface{}{"target": 50},
			PointsReward:   200,
			Category:       "attendance",
			DisplayOrder:   3,
		},
		{
			Code:           "duration_100",
			Name:           "研究熱心",
			Description:    "累計滞在時間が100時間を超えました。",
			ConditionType:  "total_duration",
			ConditionValue: map[string]interface{}{"target": 100 * 60}, // 分
			PointsReward:   100,
			Category:       "time",
			DisplayOrder:   4,
		},
	}

	for _, a := range achievements {
		var count int64
		db.Model(&domain.Achievement{}).Where("code = ?", a.Code).Count(&count)
		if count == 0 {
			// JSONB変換（簡易化のため、ここではGORMのMap処理に任せるか、自前でJSON化するか）
			// domain.Achievement の ConditionValue は JSONB (map[string]interface{}) なのでそのまま渡せるはずだが
			// GORMのdriver設定によってはValue()が必要。
			// ここではドメイン定義でValue()を実装しているのでOKなはずだが、Mapリテラル直書きだと型が合わないかも。
			// 修正: domain/achievement.go で JSONB型として定義されているならキャストが必要。
			// しかしここではシンプルに struct初期化で入れている。

			// 念のためキャスト
			// a.ConditionValue は domain.JSONB 型にする必要があるが、
			// 構造体定義が `type JSONB map[string]interface{}` なのでキャスト可能

			if err := db.Create(&a).Error; err != nil {
				log.Printf("Failed to create achievement %s: %v", a.Code, err)
			} else {
				fmt.Printf("Created achievement: %s\n", a.Name)
			}
		}
	}
}

// JSONB対応のために、これをインポートする必要があるが、
// 実行時には main パッケージ内なので適当に再定義するか type alias が見えていればOK
// ここでは domain パッケージを使っているので大丈夫。
