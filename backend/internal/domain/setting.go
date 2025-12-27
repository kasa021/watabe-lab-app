package domain

import "time"

// Setting システム設定
type Setting struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Key         string    `json:"key" gorm:"uniqueIndex;not null"`
	Value       JSONB     `json:"value" gorm:"type:jsonb;not null"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   *uint     `json:"updated_by"`
}

// TableName テーブル名を指定
func (Setting) TableName() string {
	return "settings"
}
