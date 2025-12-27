package domain

import "time"

// User ユーザー情報
type User struct {
	ID               uint       `json:"id" gorm:"primaryKey"`
	Username         string     `json:"username" gorm:"uniqueIndex;not null"`
	DisplayName      string     `json:"display_name" gorm:"not null"`
	Email            string     `json:"email"`
	Role             string     `json:"role" gorm:"not null;default:'student'"` // student, teacher, admin
	IsPresencePublic bool       `json:"is_presence_public" gorm:"not null;default:true"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	LastLoginAt      *time.Time `json:"last_login_at"`
	IsActive         bool       `json:"is_active" gorm:"not null;default:true"`
}

// TableName テーブル名を指定
func (User) TableName() string {
	return "users"
}
