package domain

type UserRanking struct {
	UserID        uint   `json:"user_id"`
	DisplayName   string `json:"display_name"`
	Username      string `json:"username"`
	TotalDuration int    `json:"total_duration"` // 分単位
}
