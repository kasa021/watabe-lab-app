package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kasa021/watabe-lab-app/internal/service"
)

// AuthMiddleware 認証ミドルウェア
func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorizationヘッダーを取得
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "MISSING_TOKEN",
					"message": "認証トークンが必要です",
				},
			})
			c.Abort()
			return
		}

		// Bearerトークンを抽出
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "INVALID_TOKEN_FORMAT",
					"message": "トークンの形式が不正です",
				},
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// トークンを検証
		claims, err := authService.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "INVALID_TOKEN",
					"message": "無効なトークンです",
				},
			})
			c.Abort()
			return
		}

		// ユーザー情報をコンテキストに設定
		c.Set("user_id", uint((*claims)["user_id"].(float64)))
		c.Set("username", (*claims)["username"].(string))
		c.Set("role", (*claims)["role"].(string))

		c.Next()
	}
}

// RoleMiddleware ロール（権限）チェックミドルウェア
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "認証が必要です",
				},
			})
			c.Abort()
			return
		}

		userRole := role.(string)
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{
				"code":    "FORBIDDEN",
				"message": "この操作を実行する権限がありません",
			},
		})
		c.Abort()
	}
}

