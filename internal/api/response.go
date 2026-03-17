package api

import "github.com/gin-gonic/gin"

const (
	InvalidParamData = "invalid param id"
)

func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, gin.H{
		"success": false,
		"error":   message,
	})
}

func NewSucceededResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"data":    data,
	})
}

func NewCustomDataSucceededResponse(c *gin.Context, statusCode int, data map[string]interface{}) {
	responseJsonBody := gin.H{
		"succes": true,
	}
	for k, v := range data {
		responseJsonBody[k] = v
	}
	c.JSON(statusCode, responseJsonBody)
}
