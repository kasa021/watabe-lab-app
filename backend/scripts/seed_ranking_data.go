package main

import (
	"fmt"
	"log"
	"math/rand"
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

	seedUsers(db)
	seedCheckInLogs(db)
}

func seedUsers(db *gorm.DB) {
	users := []domain.User{
		{Username: "student1", DisplayName: "山田 太郎", Email: "taro@example.com", Role: "student"},
		{Username: "student2", DisplayName: "佐藤 花子", Email: "hanako@example.com", Role: "student"},
		{Username: "student3", DisplayName: "鈴木 一郎", Email: "ichiro@example.com", Role: "student"},
		{Username: "student4", DisplayName: "田中 次郎", Email: "jiro@example.com", Role: "student"},
		{Username: "student5", DisplayName: "高橋 三郎", Email: "saburo@example.com", Role: "student"},
	}

	for _, u := range users {
		var count int64
		db.Model(&domain.User{}).Where("username = ?", u.Username).Count(&count)
		if count == 0 {
			if err := db.Create(&u).Error; err != nil {
				log.Printf("Failed to create user %s: %v", u.Username, err)
			} else {
				fmt.Printf("Created user: %s\n", u.DisplayName)
			}
		}
	}
}

func seedCheckInLogs(db *gorm.DB) {
	var users []domain.User
	db.Find(&users)

	// Past 30 days
	now := time.Now()
	for i := 0; i < 30; i++ {
		date := now.AddDate(0, 0, -i)
		// Skip weekends randomly (students often come on weekends too, but let's make it varied)
		if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
			if rand.Float32() < 0.5 {
				continue
			}
		}

		for _, u := range users {
			// Randomly decide if user came on this day
			if rand.Float32() < 0.2 { // 80% attendance rate is too high, let's say 80% chance
				continue
			}

			// Random duration between 2 and 10 hours
			durationHours := 2 + rand.Intn(8)
			durationMinutes := durationHours * 60

			checkIn := time.Date(date.Year(), date.Month(), date.Day(), 9+rand.Intn(3), rand.Intn(60), 0, 0, date.Location())
			checkOut := checkIn.Add(time.Duration(durationMinutes) * time.Minute)

			wifiSSID := "WatabeLabWiFi"

			logEntry := domain.CheckInLog{
				UserID:          u.ID,
				CheckInAt:       checkIn,
				CheckOutAt:      &checkOut,
				DurationMinutes: &durationMinutes,
				CheckInMethod:   "web_manual",
				WiFiSSID:        wifiSSID,
			}

			if err := db.Create(&logEntry).Error; err != nil {
				log.Printf("Failed to create log for user %d: %v", u.ID, err)
			}
		}
	}
	fmt.Println("Dummy check-in logs inserted.")
}
