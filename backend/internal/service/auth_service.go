package service

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kasa021/watabe-lab-app/internal/config"
	"github.com/kasa021/watabe-lab-app/internal/domain"
)

// AuthService 認証サービス
type AuthService struct {
	config *config.Config
}

// NewAuthService 認証サービスを作成
func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{
		config: cfg,
	}
}

// LoginRequest ログインリクエスト
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse ログインレスポンス
type LoginResponse struct {
	Token        string      `json:"token"`
	User         domain.User `json:"user"`
	ExpiresAt    time.Time   `json:"expires_at"`
}

// AuthenticateWithLDAP LDAPでユーザーを認証
func (s *AuthService) AuthenticateWithLDAP(username, password string) (*domain.User, error) {
	// LDAPサーバーに接続
	l, err := s.connectLDAP()
	if err != nil {
		return nil, fmt.Errorf("LDAP接続エラー: %w", err)
	}
	defer l.Close()

	// ユーザーのDNを検索
	userDN, attributes, err := s.searchUser(l, username)
	if err != nil {
		return nil, fmt.Errorf("ユーザー検索エラー: %w", err)
	}

	// ユーザーの認証（バインド）
	err = l.Bind(userDN, password)
	if err != nil {
		return nil, fmt.Errorf("認証失敗: %w", err)
	}

	// 認証成功：ユーザー情報を作成
	user := &domain.User{
		Username:    username,
		DisplayName: s.getAttributeValue(attributes, "cn", username),
		Email:       s.getAttributeValue(attributes, "mail", ""),
		Role:        "student", // デフォルトはstudent、必要に応じて変更
		IsActive:    true,
	}

	return user, nil
}

// connectLDAP LDAPサーバーに接続
func (s *AuthService) connectLDAP() (*ldap.Conn, error) {
	ldapURL := fmt.Sprintf("%s:%s", s.config.LDAP.Host, s.config.LDAP.Port)

	// ポート636の場合はLDAPS（TLS）、389の場合は通常のLDAP
	var l *ldap.Conn
	var err error

	if s.config.LDAP.Port == "636" {
		// LDAPS (TLS) 接続
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false, // 本番環境ではfalseに設定
			ServerName:         s.config.LDAP.Host,
		}
		l, err = ldap.DialTLS("tcp", ldapURL, tlsConfig)
	} else {
		// 通常のLDAP接続
		l, err = ldap.Dial("tcp", ldapURL)
		if err != nil {
			return nil, err
		}

		// StartTLSでアップグレード（推奨）
		err = l.StartTLS(&tls.Config{
			InsecureSkipVerify: false,
			ServerName:         s.config.LDAP.Host,
		})
		// StartTLSが失敗しても続行（開発環境用）
		// 本番環境ではエラーを返すべき
	}

	if err != nil {
		return nil, err
	}

	// 管理者アカウントでバインド（ユーザー検索用）
	// BindUserが設定されていない場合は匿名バインド
	if s.config.LDAP.BindUser != "" {
		err = l.Bind(s.config.LDAP.BindUser, s.config.LDAP.BindPass)
		if err != nil {
			l.Close()
			return nil, fmt.Errorf("管理者バインド失敗: %w", err)
		}
	}

	return l, nil
}

// searchUser ユーザーを検索してDNを取得
func (s *AuthService) searchUser(l *ldap.Conn, username string) (string, map[string]string, error) {
	searchRequest := ldap.NewSearchRequest(
		s.config.LDAP.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, // サイズ制限なし
		0, // タイムアウトなし
		false,
		fmt.Sprintf("(uid=%s)", ldap.EscapeFilter(username)), // uidで検索（環境に応じて変更）
		[]string{"dn", "cn", "mail", "displayName"},          // 取得する属性
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return "", nil, err
	}

	if len(sr.Entries) == 0 {
		return "", nil, fmt.Errorf("ユーザーが見つかりません: %s", username)
	}

	if len(sr.Entries) > 1 {
		return "", nil, fmt.Errorf("複数のユーザーが見つかりました: %s", username)
	}

	entry := sr.Entries[0]
	attributes := make(map[string]string)
	for _, attr := range entry.Attributes {
		if len(attr.Values) > 0 {
			attributes[attr.Name] = attr.Values[0]
		}
	}

	return entry.DN, attributes, nil
}

// getAttributeValue 属性値を取得（存在しない場合はデフォルト値）
func (s *AuthService) getAttributeValue(attributes map[string]string, key, defaultValue string) string {
	if val, ok := attributes[key]; ok && val != "" {
		return val
	}
	return defaultValue
}

// GenerateJWT JWTトークンを生成
func (s *AuthService) GenerateJWT(user *domain.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Hour * time.Duration(s.config.JWT.ExpireHour))

	claims := jwt.MapClaims{
		"user_id":      user.ID,
		"username":     user.Username,
		"display_name": user.DisplayName,
		"role":         user.Role,
		"exp":          expiresAt.Unix(),
		"iat":          time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("トークン生成エラー: %w", err)
	}

	return tokenString, expiresAt, nil
}

// ValidateJWT JWTトークンを検証
func (s *AuthService) ValidateJWT(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不正な署名方式: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, fmt.Errorf("無効なトークン")
}

