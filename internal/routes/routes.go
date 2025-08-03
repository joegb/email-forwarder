package routes

import (
	"github.com/joegb/email-forwarder/internal/controllers"
	"github.com/joegb/email-forwarder/internal/middleware"
	
	"github.com/gin-gonic/gin"
)

func SetupTargetRoutes(router *gin.Engine) {
	// 添加认证中间件
	authMiddleware := middleware.AuthMiddleware()
	
	targetGroup := router.Group("/api/targets")
	targetGroup.Use(authMiddleware) // 应用认证中间件
	{
		targetGroup.POST("", controllers.CreateTarget)
		targetGroup.GET("", controllers.ListTargets)
		targetGroup.GET("/:id", controllers.GetTarget)
		targetGroup.PUT("/:id", controllers.UpdateTarget)
		targetGroup.DELETE("/:id", controllers.DeleteTarget)
	}
}