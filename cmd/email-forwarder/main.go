package main

import (
	"log"
	// "net/http"
	"os"
	"os/signal"
	"syscall"
	"io"
	"time"
	"fmt"

	// "github.com/joegb/email-forwarder/internal/config"
	"github.com/joegb/email-forwarder/internal/database"
	"github.com/joegb/email-forwarder/internal/logger"
	"github.com/joegb/email-forwarder/internal/routes"
	"github.com/joegb/email-forwarder/internal/middleware"
	"github.com/joegb/email-forwarder/internal/cron"

	"github.com/gin-gonic/gin"
)

func main() {
	// 设置日志输出
	setupLogger()

	// 初始化邮件处理日志
	logger.Init()

	// 初始化数据库
	database.Connect()

	// 启动定时任务
	go cron.StartEmailCron()

	// 创建Gin应用
	router := gin.New()

	// 添加速率限制中间件
	router.Use(middleware.RateLimitMiddleware(middleware.RateLimitConfig{
		Period: 30 * time.Second,
		Limit:  30, // 每30秒最多30个请求
	}))
	
	// 添加 Gin 日志中间件
	router.Use(gin.LoggerWithFormatter(logFormatter))
	router.Use(gin.Recovery())
	
	// 设置路由
	routes.SetupTargetRoutes(router)
	
	// 添加健康检查端点
	router.GET("/health", healthCheck)

	// 启动HTTP服务器
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		log.Printf("Server running on port %s", port)
		if err := router.Run(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
}

// 自定义日志格式
func logFormatter(param gin.LogFormatterParams) string {
	return fmt.Sprintf("[%s] %s %s %d %s \"%s\"\n",
		param.TimeStamp.Format(time.RFC3339),
		param.Method,
		param.Path,
		param.StatusCode,
		param.Latency,
		param.ErrorMessage,
	)
}

// 设置日志输出
func setupLogger() {
	// 创建日志目录
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}
	
	// 创建日志文件
	logFile, err := os.OpenFile("logs/email_forwarder.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	
	// 设置日志输出到文件和控制台
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	gin.DefaultWriter = multiWriter
}

// 健康检查端点
func healthCheck(c *gin.Context) {
	// 检查数据库连接
	if database.DB != nil {
		db, err := database.DB.DB()
		if err == nil {
			if err := db.Ping(); err == nil {
				c.JSON(200, gin.H{
					"status": "ok",
					"db": "connected",
					"version": "1.0.0",
					"timestamp": time.Now().Format(time.RFC3339),
				})
				return
			}
		}
	}
	
	c.JSON(500, gin.H{
		"status": "error",
		"db": "disconnected",
		"message": "Database connection failed",
	})
}