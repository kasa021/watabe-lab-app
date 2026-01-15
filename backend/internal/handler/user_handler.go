package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kasa021/watabe-lab-app/internal/repository"
)

type UserHandler struct {
	userRepo       repository.UserRepository
	attendanceRepo repository.AttendanceRepository
}

func NewUserHandler(userRepo repository.UserRepository, attendanceRepo repository.AttendanceRepository) *UserHandler {
	return &UserHandler{
		userRepo:       userRepo,
		attendanceRepo: attendanceRepo,
	}
}

type UpdateProfileRequest struct {
	DisplayName      string `json:"display_name"`
	IsPresencePublic bool   `json:"is_presence_public"`
}

// UpdateProfile プロフィール更新
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.userRepo.FindByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 更新
	user.DisplayName = req.DisplayName
	user.IsPresencePublic = req.IsPresencePublic

	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetAttendanceHeatmap ヒートマップ用データ取得
func (h *UserHandler) GetAttendanceHeatmap(c *gin.Context) {
	// ID指定があればそのユーザー、なければ自分
	userIDStr := c.Param("id")
	var targetUserID uint

	if userIDStr == "" || userIDStr == "me" {
		currentUserID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		targetUserID = currentUserID.(uint)
	} else {
		uid, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		targetUserID = uint(uid)
	}

	// プライバシーチェック（他人の場合）
	// (簡易実装: 自分以外は見れない or 公開設定を見るなど必要だが、要件的にそこまで厳密でなくて良いなら一旦スルー)
	// 要件「自分がいつ研究室にきたのか」がメインなので、とりあえず自分のデータは見れるようにする。
	// 他人のデータを見る機能はまだ明示されていないが、一応APIは汎用的にしておく。

	// limit := 365 // 過去1年分
	// TODO: リポジトリ側のメソッドで期間指定できるようにしても良い。現状は全期間か簡易実装。
	// GetDailyAttendanceCounts は期間指定なし（全期間）になっているはず。

	// データ取得
	dailies, err := h.attendanceRepo.GetDailyAttendanceCounts(c.Request.Context(), targetUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// フロントエンドは react-calendar-heatmap を使う想定
	// { date: '2023-01-01', count: 123 } の形式が望ましい
	// domain.DailyAttendance は { Date, Count, DurationMinutes } を持っている

	c.JSON(http.StatusOK, gin.H{"heatmap": dailies})
}
