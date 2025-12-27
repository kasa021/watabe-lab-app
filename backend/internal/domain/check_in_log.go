package domain

import "time"

// CheckInLog チェックインログ
type CheckInLog struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	UserID          uint       `json:"user_id" gorm:"not null;index"`
	CheckInAt       time.Time  `json:"check_in_at" gorm:"not null;index"`
	CheckOutAt      *time.Time `json:"check_out_at" gorm:"index"`
	DurationMinutes *int       `json:"duration_minutes"`
	CheckInMethod   string     `json:"check_in_method"`
	WiFiSSID        string     `json:"wifi_ssid"`
	GPSLatitude     *float64   `json:"gps_latitude"`
	GPSLongitude    *float64   `json:"gps_longitude"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	// リレーション
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName テーブル名を指定
func (CheckInLog) TableName() string {
	return "check_in_logs"
}
