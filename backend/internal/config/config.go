package config

import (
	"os"
	"strconv"
)

// Config システム全体の設定
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	LDAP     LDAPConfig
	JWT      JWTConfig
	Location LocationConfig
}

// ServerConfig サーバー設定
type ServerConfig struct {
	Port string
	Env  string
}

// DatabaseConfig データベース設定
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LDAPConfig LDAP認証設定
type LDAPConfig struct {
	Host       string
	Port       string
	BaseDN     string
	BindUser   string
	BindPass   string
	StartTLS   bool
	SkipVerify bool
}

// JWTConfig JWT設定
type JWTConfig struct {
	Secret     string
	ExpireHour int
}

// LocationConfig 位置情報設定
type LocationConfig struct {
	WiFiSSIDs []string
	Latitude  float64
	Longitude float64
	Radius    float64
}

// Load 環境変数から設定を読み込む
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "labuser"),
			Password: getEnv("DB_PASSWORD", "labpass"),
			DBName:   getEnv("DB_NAME", "lab_attendance"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		LDAP: LDAPConfig{
			Host:       getEnv("LDAP_HOST", "localhost"),
			Port:       getEnv("LDAP_PORT", "389"),
			BaseDN:     getEnv("LDAP_BASE_DN", "dc=example,dc=com"),
			BindUser:   getEnv("LDAP_BIND_USER", ""),
			BindPass:   getEnv("LDAP_BIND_PASS", ""),
			StartTLS:   getEnvAsBool("LDAP_START_TLS", true),    // デフォルトは有効
			SkipVerify: getEnvAsBool("LDAP_SKIP_VERIFY", false), // デフォルトは検証する
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key-change-this"),
			ExpireHour: getEnvAsInt("JWT_EXPIRE_HOUR", 24),
		},
		Location: LocationConfig{
			WiFiSSIDs: []string{"WatabeLabWiFi"},
			Latitude:  getEnvAsFloat("LAB_LATITUDE", 35.6812),
			Longitude: getEnvAsFloat("LAB_LONGITUDE", 139.7671),
			Radius:    getEnvAsFloat("LAB_RADIUS_METERS", 100.0),
		},
	}
}

// getEnv 環境変数を取得（デフォルト値あり）
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 環境変数を整数として取得
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsFloat 環境変数を浮動小数点数として取得
func getEnvAsFloat(key string, defaultValue float64) float64 {
	valueStr := os.Getenv(key)
	if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsBool 環境変数を真偽値として取得
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}
