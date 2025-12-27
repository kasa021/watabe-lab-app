package domain

import "time"

// DailyAttendance 日次出席記録
type DailyAttendance struct {
	ID                   uint       `json:"id" gorm:"primaryKey"`
	UserID               uint       `json:"user_id" gorm:"not null;index"`
	AttendanceDate       time.Time  `json:"attendance_date" gorm:"not null;type:date;index"`
	TotalDurationMinutes int        `json:"total_duration_minutes" gorm:"not null;default:0"`
	CheckInCount         int        `json:"check_in_count" gorm:"not null;default:0"`
	FirstCheckInAt       *time.Time `json:"first_check_in_at" gorm:"type:time"`
	LastCheckOutAt       *time.Time `json:"last_check_out_at" gorm:"type:time"`
	Points               int        `json:"points" gorm:"not null;default:0"`
	IsHoliday            bool       `json:"is_holiday" gorm:"not null;default:false"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`

	// リレーション
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName テーブル名を指定
func (DailyAttendance) TableName() string {
	return "daily_attendances"
}
