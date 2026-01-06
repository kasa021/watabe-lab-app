package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kasa021/watabe-lab-app/internal/service"
)

type AttendanceHandler struct {
	service service.AttendanceService
}

func NewAttendanceHandler(service service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{service: service}
}

func (h *AttendanceHandler) CheckIn(c *gin.Context) {
	// ミドルウェアでセットされたUserIDを取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req service.CheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CheckIn(c.Request.Context(), userID.(uint), &req); err != nil {
		if errors.Is(err, service.ErrAlreadyCheckedIn) {
			c.JSON(http.StatusConflict, gin.H{"error": "already checked in"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "checked in successfully"})
}

func (h *AttendanceHandler) CheckOut(c *gin.Context) {
	// ミドルウェアでセットされたUserIDを取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.CheckOut(c.Request.Context(), userID.(uint)); err != nil {
		if errors.Is(err, service.ErrNotCheckedIn) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "not checked in"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "checked out successfully"})
}
