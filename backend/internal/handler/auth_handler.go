package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kasa021/watabe-lab-app/internal/repository"
	"github.com/kasa021/watabe-lab-app/internal/service"
)

// AuthHandler 認証ハンドラー
type AuthHandler struct {
	authService *service.AuthService
	userRepo    repository.UserRepository
}

// NewAuthHandler 認証ハンドラーを作成
func NewAuthHandler(authService *service.AuthService, userRepo repository.UserRepository) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userRepo:    userRepo,
	}
}

// Login ログイン処理
// @Summary ログイン
// @Description LDAPでユーザーを認証し、JWTトークンを発行
// @Tags auth
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "ログイン情報"
// @Success 200 {object} service.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: "リクエストが不正です",
			},
		})
		return
	}

	// LDAPで認証
	ldapUser, err := h.authService.AuthenticateWithLDAP(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: ErrorDetail{
				Code:    "AUTHENTICATION_FAILED",
				Message: "ユーザー名またはパスワードが正しくありません",
			},
		})
		return
	}

	// データベースでユーザーを検索または作成
	user, err := h.userRepo.FindByUsername(req.Username)
	if err != nil {
		// ユーザーが存在しない場合は新規作成
		user = ldapUser
		if err := h.userRepo.Create(user); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: ErrorDetail{
					Code:    "USER_CREATION_FAILED",
					Message: "ユーザーの作成に失敗しました",
				},
			})
			return
		}
	} else {
		// 既存ユーザーの情報を更新
		user.DisplayName = ldapUser.DisplayName
		user.Email = ldapUser.Email
		now := time.Now()
		user.LastLoginAt = &now
		if err := h.userRepo.Update(user); err != nil {
			// 更新失敗してもログインは続行
		}
	}

	// JWTトークンを生成
	token, expiresAt, err := h.authService.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetail{
				Code:    "TOKEN_GENERATION_FAILED",
				Message: "トークンの生成に失敗しました",
			},
		})
		return
	}

	c.JSON(http.StatusOK, service.LoginResponse{
		Token:     token,
		User:      *user,
		ExpiresAt: expiresAt,
	})
}

// Me 現在のユーザー情報を取得
// @Summary 現在のユーザー情報
// @Description JWTトークンから現在のユーザー情報を取得
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.User
// @Failure 401 {object} ErrorResponse
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	// ミドルウェアで設定されたユーザー情報を取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: ErrorDetail{
				Code:    "UNAUTHORIZED",
				Message: "認証が必要です",
			},
		})
		return
	}

	user, err := h.userRepo.FindByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: ErrorDetail{
				Code:    "USER_NOT_FOUND",
				Message: "ユーザーが見つかりません",
			},
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ErrorResponse エラーレスポンス
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail エラー詳細
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

