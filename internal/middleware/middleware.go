package middleware

import (
	"encoding/base64"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	ginLimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从环境变量获取 API 密钥
		apiKey := os.Getenv("API_KEY")
		if apiKey == "" {
			c.Next() // 没有设置 API 密钥则跳过认证
			return
		}
		
		// 获取 Authorization 头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header required"})
			return
		}
		
		// 解析 Basic Auth
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Basic" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid authorization format"})
			return
		}
		
		// 解码凭证
		payload, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid authorization encoding"})
			return
		}
		
		// 验证凭证
		if string(payload) != "api:"+apiKey {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid API key"})
			return
		}
		
		c.Next()
	}
}

// RateLimitConfig 速率限制配置
type RateLimitConfig struct {
	Period time.Duration
	Limit  int64
}

// DefaultRateLimitConfig 默认速率限制配置
var DefaultRateLimitConfig = RateLimitConfig{
	Period: 1 * time.Minute,
	Limit:  60,
}

// RateLimitMiddleware 创建速率限制中间件
func RateLimitMiddleware(config ...RateLimitConfig) gin.HandlerFunc {
	// 使用提供的配置或默认配置
	cfg := DefaultRateLimitConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	// 创建速率限制器配置
	rate := limiter.Rate{
		Period: cfg.Period,
		Limit:  cfg.Limit,
	}

	// 使用内存存储
	store := memory.NewStore()

	// 创建限制器实例
	limiterInstance := limiter.New(store, rate)

	// 创建 Gin 中间件
	return ginLimiter.NewMiddleware(limiterInstance)
}
