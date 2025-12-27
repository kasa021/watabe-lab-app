package domain

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// JSONB カスタム型（PostgreSQLのJSONB型用）
type JSONB map[string]interface{}

// Value JSONB型をデータベース値に変換
func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan データベース値をJSONB型に変換
func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Achievement 称号マスタ
type Achievement struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	Code           string    `json:"code" gorm:"uniqueIndex;not null"`
	Name           string    `json:"name" gorm:"not null"`
	Description    string    `json:"description"`
	IconURL        string    `json:"icon_url"`
	Category       string    `json:"category"` // attendance, time, streak, special
	ConditionType  string    `json:"condition_type" gorm:"not null"`
	ConditionValue JSONB     `json:"condition_value" gorm:"type:jsonb"`
	PointsReward   int       `json:"points_reward" gorm:"default:0"`
	IsActive       bool      `json:"is_active" gorm:"not null;default:true"`
	DisplayOrder   int       `json:"display_order" gorm:"not null;default:0"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName テーブル名を指定
func (Achievement) TableName() string {
	return "achievements"
}

// UserAchievement ユーザーが獲得した称号
type UserAchievement struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id" gorm:"not null;index"`
	AchievementID uint      `json:"achievement_id" gorm:"not null;index"`
	AchievedAt    time.Time `json:"achieved_at" gorm:"not null"`

	// リレーション
	User        User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Achievement Achievement `json:"achievement,omitempty" gorm:"foreignKey:AchievementID"`
}

// TableName テーブル名を指定
func (UserAchievement) TableName() string {
	return "user_achievements"
}
