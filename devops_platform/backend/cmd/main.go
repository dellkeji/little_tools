package main

import (
	"devops-platform/internal/api"
	"devops-platform/internal/config"
	"devops-platform/internal/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 创建Gin引擎
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORS())
	r.Use(middleware.ErrorHandler())

	// API路由
	apiGroup := r.Group("/api/v1")
	{
		// Jenkins相关
		jenkins := apiGroup.Group("/jenkins")
		{
			jenkins.POST("/connect", api.ConnectJenkins)
			jenkins.GET("/nodes", api.GetJenkinsNodes)
			jenkins.POST("/nodes/:name/toggle", api.ToggleJenkinsNode)
			jenkins.GET("/nodes/:name", api.GetJenkinsNodeInfo)
		}

		// MySQL相关
		mysql := apiGroup.Group("/mysql")
		{
			mysql.POST("/connect", api.ConnectMySQL)
			mysql.POST("/execute", api.ExecuteSQL)
			mysql.POST("/validate", api.ValidateSQL)
			mysql.GET("/databases", api.GetDatabases)
			mysql.GET("/tables", api.GetTables)
		}

		// Redis相关
		redis := apiGroup.Group("/redis")
		{
			redis.POST("/connect", api.ConnectRedis)
			redis.GET("/keys", api.GetRedisKeys)
			redis.GET("/key/:key", api.GetRedisValue)
			redis.POST("/key", api.SetRedisValue)
			redis.DELETE("/key/:key", api.DeleteRedisKey)
			redis.GET("/info", api.GetRedisInfo)
		}
	}

	// 启动服务器
	log.Printf("服务器启动在端口 %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
