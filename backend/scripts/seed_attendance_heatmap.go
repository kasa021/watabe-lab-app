package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/kasa021/watabe-lab-app/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// ターゲットユーザー: mkobayashi
	var user domain.User
	if err := db.Where("username = ?", "mkobayashi").First(&user).Error; err != nil {
		log.Fatalf("failed to find user mkobayashi: %v", err)
	}

	log.Printf("Seeding attendance data for user: %s (ID: %d)", user.Username, user.ID)

	// 過去1年分
	now := time.Now()
	oneYearAgo := now.AddDate(-1, 0, 0)

	// 生成パラメータ
	attendProbability := 0.7  // 平日の出席確率
	weekendProbability := 0.2 // 休日の出席確率

	rand.Seed(time.Now().UnixNano())

	count := 0
	for d := oneYearAgo; d.Before(now); d = d.AddDate(0, 0, 1) {
		isWeekend := d.Weekday() == time.Saturday || d.Weekday() == time.Sunday
		prob := attendProbability
		if isWeekend {
			prob = weekendProbability
		}

		if rand.Float64() < prob {
			// 出席する

			// 入室時間: 9:00 ~ 11:00 の間
			startHour := 9 + rand.Intn(2)
			startMin := rand.Intn(60)
			checkInAt := time.Date(d.Year(), d.Month(), d.Day(), startHour, startMin, 0, 0, d.Location())

			// 滞在時間: 2時間 ~ 10時間
			durationHours := 2 + rand.Intn(9)
			durationMinutes := durationHours * 60

			checkOutAt := checkInAt.Add(time.Duration(durationMinutes) * time.Minute)

			duration := int(durationMinutes) // int pointer
			lat := 35.6 + rand.Float64()*0.1
			lng := 139.7 + rand.Float64()*0.1
			ssid := "Campus_WiFi"

			logEntry := domain.CheckInLog{
				UserID:          user.ID,
				CheckInAt:       checkInAt,
				CheckOutAt:      &checkOutAt,
				DurationMinutes: &duration,
				CheckInMethod:   "wifi",
				WiFiSSID:        ssid,
				GPSLatitude:     &lat,
				GPSLongitude:    &lng,
				CreatedAt:       time.Now(), // Record creation time is now
				UpdatedAt:       time.Now(),
			}

			if err := db.Create(&logEntry).Error; err != nil {
				log.Printf("Failed to create log for %s: %v", d.Format("2006-01-02"), err)
			} else {
				count++
			}

			// たまに1日2回チェックインする
			if rand.Float64() < 0.1 {
				// 午後また来る
				checkInAt2 := checkOutAt.Add(time.Hour * 1) // 1時間後
				duration2 := 60 + rand.Intn(120)            // 1-3時間
				checkOutAt2 := checkInAt2.Add(time.Duration(duration2) * time.Minute)

				logEntry2 := domain.CheckInLog{
					UserID:          user.ID,
					CheckInAt:       checkInAt2,
					CheckOutAt:      &checkOutAt2,
					DurationMinutes: &duration2,
					CheckInMethod:   "wifi",
					WiFiSSID:        ssid,
					GPSLatitude:     &lat,
					GPSLongitude:    &lng,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				}
				db.Create(&logEntry2)
				count++
			}
		}
	}

	log.Printf("Successfully inserted %d check-in logs.", count)
}
