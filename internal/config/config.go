package config

import (
	"os"
	"strconv"
	"time"
)

// 获取配置值，支持默认值
func GetString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func GetInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func GetDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// 应用配置
type AppConfig struct {
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	GmailEmail       string
	GmailAppPassword string
	Port             string
	APIKey           string
	CronSchedule     string
}

// 获取应用配置
func GetConfig() *AppConfig {
	return &AppConfig{
		DBHost:           GetString("DB_HOST", "db"),
		DBPort:           GetString("DB_PORT", "3306"),
		DBUser:           GetString("DB_USER", "root"),
		DBPassword:       GetString("DB_PASSWORD", "mysecretpassword"),
		DBName:           GetString("DB_NAME", "email_forwarder"),
		GmailEmail:       GetString("GMAIL_EMAIL", ""),
		GmailAppPassword: GetString("GMAIL_APP_PASSWORD", ""),
		Port:             GetString("PORT", "8080"),
		APIKey:           GetString("API_KEY", ""),
		CronSchedule:     GetString("CRON_SCHEDULE", "*/5 * * * *"),
	}
}