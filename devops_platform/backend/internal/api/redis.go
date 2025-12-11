package api

import (
	"devops-platform/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RedisConnectRequest struct {
	Host     string `json:"host" binding:"required"`
	Port     string `json:"port" binding:"required"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

func ConnectRedis(c *gin.Context) {
	var req RedisConnectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	svc := service.GetRedisService()
	if err := svc.Connect(req.Host, req.Port, req.Password, req.DB); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "连接成功"})
}

func GetRedisKeys(c *gin.Context) {
	pattern := c.Query("pattern")
	svc := service.GetRedisService()
	keys, err := svc.GetKeys(c.Request.Context(), pattern)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"keys": keys})
}

func GetRedisValue(c *gin.Context) {
	key := c.Param("key")
	svc := service.GetRedisService()
	value, err := svc.GetValue(c.Request.Context(), key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, value)
}

type SetRedisValueRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
	TTL   int    `json:"ttl"`
}

func SetRedisValue(c *gin.Context) {
	var req SetRedisValueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	svc := service.GetRedisService()
	if err := svc.SetValue(c.Request.Context(), req.Key, req.Value, req.TTL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "设置成功"})
}

func DeleteRedisKey(c *gin.Context) {
	key := c.Param("key")
	svc := service.GetRedisService()
	if err := svc.DeleteKey(c.Request.Context(), key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func GetRedisInfo(c *gin.Context) {
	svc := service.GetRedisService()
	info, err := svc.GetInfo(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}
