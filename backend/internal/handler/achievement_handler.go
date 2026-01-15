package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kasa021/watabe-lab-app/internal/service"
)

type AchievementHandler struct {
	service service.AchievementService
}

func NewAchievementHandler(service service.AchievementService) *AchievementHandler {
	return &AchievementHandler{service: service}
}

func (h *AchievementHandler) GetAchievements(c *gin.Context) {
	achievements, err := h.service.GetAchievements(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"achievements": achievements})
}

func (h *AchievementHandler) GetUserAchievements(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	achievements, err := h.service.GetUserAchievements(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_achievements": achievements})
}

func (h *AchievementHandler) GetMyAchievements(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	achievements, err := h.service.GetUserAchievements(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_achievements": achievements})
}
