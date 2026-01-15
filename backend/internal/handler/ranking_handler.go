package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kasa021/watabe-lab-app/internal/domain"
	"github.com/kasa021/watabe-lab-app/internal/service"
)

type RankingHandler struct {
	service service.RankingService
}

func NewRankingHandler(service service.RankingService) *RankingHandler {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rankings": rankings})
}
