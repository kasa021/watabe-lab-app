package handler

import (
	"errors"
	"net/http"
	"strings"

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

	// Set Client IP
	req.ClientIP = c.ClientIP()

	if err := h.service.CheckIn(c.Request.Context(), userID.(uint), &req); err != nil {
		if errors.Is(err, service.ErrAlreadyCheckedIn) {
			c.JSON(http.StatusConflict, gin.H{"error": "already checked in"})
			return
		}
		if errors.Is(err, service.ErrRestrictionViolation) {
			// Clean up message for frontend display by removing the technical suffix
			// fmt.Errorf("msg: %w", err) produces "msg: error_string"
			msg := strings.TrimSuffix(err.Error(), ": "+service.ErrRestrictionViolation.Error())
			c.JSON(http.StatusForbidden, gin.H{"error": msg})
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

func (h *AttendanceHandler) GetActiveUsers(c *gin.Context) {
	logs, err := h.service.GetActiveUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"active_users": logs})
}
