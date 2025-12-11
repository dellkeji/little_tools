package api

import (
	"devops-platform/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MySQLConnectRequest struct {
	Host     string `json:"host" binding:"required"`
	Port     string `json:"port" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
	Database string `json:"database" binding:"required"`
}

func ConnectMySQL(c *gin.Context) {
	var req MySQLConnectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	svc := service.GetMySQLService()
	if err := svc.Connect(req.Host, req.Port, req.Username, req.Password, req.Database); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "连接成功"})
}

type SQLRequest struct {
	Query string `json:"query" binding:"required"`
}

func ValidateSQL(c *gin.Context) {
	var req SQLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	svc := service.GetMySQLService()
	if err := svc.ValidateSQL(req.Query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "valid": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true, "message": "SQL语句合法"})
}

func ExecuteSQL(c *gin.Context) {
	var req SQLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	svc := service.GetMySQLService()
	result, err := svc.ExecuteSQL(req.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

func GetDatabases(c *gin.Context) {
	svc := service.GetMySQLService()
	databases, err := svc.GetDatabases()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"databases": databases})
}

func GetTables(c *gin.Context) {
	database := c.Query("database")
	if database == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "database参数必填"})
		return
	}

	svc := service.GetMySQLService()
	tables, err := svc.GetTables(database)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tables": tables})
}
