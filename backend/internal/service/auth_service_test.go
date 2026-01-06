package service_test

import (
	"log"
	"testing"

	"github.com/kasa021/watabe-lab-app/internal/config"
	"github.com/kasa021/watabe-lab-app/internal/service"
)

// TestAuthService_AuthenticateWithLDAP_Integration は実際のLDAPサーバーを使った統合テストです。
// 事前に docker-compose up -d openldap を実行しておく必要があります。
func TestAuthService_AuthenticateWithLDAP_Integration(t *testing.T) {
	// 設定オブジェクトの手動構築
	// 注意: ポート389はdocker-composeでlocalhost:389にマップされている前提
	cfg := &config.Config{
		LDAP: config.LDAPConfig{
			Host:     "localhost",
			Port:     "389",
			BaseDN:   "ou=People,dc=example,dc=com",
			BindUser: "cn=admin,dc=example,dc=com", // osixia/openldapのデフォルト管理者DN
			BindPass: "adminpassword",              // docker-compose.ymlで設定したパスワード
			StartTLS: false,                        // ローカル環境ではTLSなし
		},
	}

	authService := service.NewAuthService(cfg)

	// ケース1: 正しいユーザー情報で認証成功
	t.Run("ValidCredentials", func(t *testing.T) {
		user, err := authService.AuthenticateWithLDAP("student1", "password")
		if err != nil {
			t.Fatalf("LDAP認証に失敗しました: %v", err)
		}

		if user.Username != "student1" {
			t.Errorf("ユーザー名が期待値と異なります: got %v, want %v", user.Username, "student1")
		}

		log.Printf("認証成功: %+v", user)
	})

	// ケース2: 間違ったパスワードで認証失敗
	t.Run("InvalidPassword", func(t *testing.T) {
		_, err := authService.AuthenticateWithLDAP("student1", "wrongpass")
		if err == nil {
			t.Fatal("パスワード間違いでエラーが発生しませんでした")
		}
	})

	// ケース3: 存在しないユーザー
	t.Run("UserNotFound", func(t *testing.T) {
		_, err := authService.AuthenticateWithLDAP("unknown_user", "password")
		if err == nil {
			t.Fatal("存在しないユーザーでエラーが発生しませんでした")
		}
	})
}
