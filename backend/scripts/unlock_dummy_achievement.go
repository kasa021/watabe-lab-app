package main

import (
	"fmt"
	"log"
	"time"

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

	unlockAchievements(db, "mkobayashi")
}

func unlockAchievements(db *gorm.DB, username string) {
	var user domain.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		log.Fatalf("User %s not found: %v", username, err)
	}

	achievementsInfo := []string{"check_in_1", "duration_100"}

	for _, code := range achievementsInfo {
		var ach domain.Achievement
		if err := db.Where("code = ?", code).First(&ach).Error; err != nil {
			log.Printf("Achievement %s not found, skipping", code)
			continue
		}

		var count int64
		db.Model(&domain.UserAchievement{}).
			Where("user_id = ? AND achievement_id = ?", user.ID, ach.ID).
			Count(&count)

		if count == 0 {
			ua := domain.UserAchievement{
				UserID:        user.ID,
				AchievementID: ach.ID,
				AchievedAt:    time.Now(),
			}
			if err := db.Create(&ua).Error; err != nil {
				log.Printf("Failed to unlock %s for %s: %v", ach.Name, username, err)
			} else {
				fmt.Printf("Unlocked: %s for %s\n", ach.Name, username)
			}
		} else {
			fmt.Printf("Already unlocked: %s for %s\n", ach.Name, username)
		}
	}
}
