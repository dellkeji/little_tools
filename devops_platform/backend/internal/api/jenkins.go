package api

import (
	"devops-platform/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JenkinsConnectRequest struct {
	URL      string `json:"url" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func ConnectJenkins(c *gin.Context) {
	var req JenkinsConnectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	svc := service.GetJenkinsService()
	if err := svc.Connect(req.URL, req.Username, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "连接成功"})
}

func GetJenkinsNodes(c *gin.Context) {
	svc := service.GetJenkinsService()
	nodes, err := svc.GetNodes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"nodes": nodes})
}

func GetJenkinsNodeInfo(c *gin.Context) {
	nodeName := c.Param("name")
	svc := service.GetJenkinsService()
	info, err := svc.GetNodeInfo(c.Request.Context(), nodeName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}

type ToggleNodeRequest struct {
	Offline bool `json:"offline"`
}

func ToggleJenkinsNode(c *gin.Context) {
	nodeName := c.Param("name")
	var req ToggleNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	svc := service.GetJenkinsService()
	if err := svc.ToggleNode(c.Request.Context(), nodeName, req.Offline); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "操作成功"})
}
