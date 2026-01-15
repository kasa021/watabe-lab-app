package service

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kasa021/watabe-lab-app/internal/domain"
	"github.com/kasa021/watabe-lab-app/internal/repository"
)

type RankingService interface {
	GetWeeklyRanking(ctx context.Context) ([]domain.UserRanking, error)
	GetMonthlyRanking(ctx context.Context) ([]domain.UserRanking, error)
	GetTotalRanking(ctx context.Context) ([]domain.UserRanking, error)
}

type rankingService struct {
	repo repository.AttendanceRepository
}

func NewRankingService(repo repository.AttendanceRepository) RankingService {
	return &rankingService{repo: repo}
}

func (s *rankingService) GetWeeklyRanking(ctx context.Context) ([]domain.UserRanking, error) {
	now := time.Now()
	// 週の開始（月曜日）を取得
	offset := int(now.Weekday())
	if offset == 0 { // 日曜日の場合
		offset = 7
	}
	startOfWeek := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -offset+1)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	return s.repo.GetUserRanking(ctx, startOfWeek, endOfWeek)
}

func (s *rankingService) GetMonthlyRanking(ctx context.Context) ([]domain.UserRanking, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	return s.repo.GetUserRanking(ctx, startOfMonth, endOfMonth)
}

func (s *rankingService) GetTotalRanking(ctx context.Context) ([]domain.UserRanking, error) {
	// 全期間なので、十分に古い日付から未来まで
	start := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Now().AddDate(100, 0, 0)

	return s.repo.GetUserRanking(ctx, start, end)
}

// Handler integration
type RankingHandler struct {
	service RankingService
}

func NewRankingHandler(service RankingService) *RankingHandler {
	return &RankingHandler{service: service}
}

func (h *RankingHandler) GetRankings(c *gin.Context) {
	rankingType := c.DefaultQuery("type", "weekly")

	var rankings []domain.UserRanking
	var err error

	switch rankingType {
	case "weekly":
		rankings, err = h.service.GetWeeklyRanking(c.Request.Context())
	case "monthly":
		rankings, err = h.service.GetMonthlyRanking(c.Request.Context())
	case "total":
		rankings, err = h.service.GetTotalRanking(c.Request.Context())
	default:
		// default to weekly
		rankings, err = h.service.GetWeeklyRanking(c.Request.Context())
	}

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"rankings": rankings})
}
