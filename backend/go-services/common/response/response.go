package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data"`
	Error     interface{} `json:"error"`
	Timestamp time.Time   `json:"timestamp"`
}

type APIError struct {
	Code    string `json:"code"`
	Key     string `json:"key"`
	Message string `json:"message"`
}

func Success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, APIResponse{
		Success:   true,
		Data:      data,
		Error:     nil,
		Timestamp: time.Now().UTC(),
	})
}

func OK(c *gin.Context, data interface{}) {
	Success(c, http.StatusOK, data)
}

func Created(c *gin.Context, data interface{}) {
	Success(c, http.StatusCreated, data)
}

func Error(c *gin.Context, status int, code string, key string, message string) {
	c.JSON(status, APIResponse{
		Success: false,
		Data:    nil,
		Error: APIError{
			Code:    code,
			Key:     key,
			Message: message,
		},
		Timestamp: time.Now().UTC(),
	})
}
