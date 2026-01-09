package responce

import "github.com/gin-gonic/gin"

// Response — единая структура ответа для API в 2026 году
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func SendError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, Response{
		Success: false,
		Message: message,
	})
}

func SendSuccess(c *gin.Context, code int, data any) {
	c.JSON(code, Response{
		Success: true,
		Data:    data,
	})
}
