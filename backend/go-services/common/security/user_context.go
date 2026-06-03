package security

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	HeaderUserID        = "X-User-Id"
	HeaderUserEmail     = "X-User-Email"
	HeaderUserRoles     = "X-User-Roles"
	HeaderUserPerms     = "X-User-Permissions"
	HeaderTraceID       = "X-Trace-Id"
	HeaderRequestID     = "X-Request-Id"
	HeaderCorrelationID = "X-Correlation-Id"

	ContextUser = "user_context"
)

type UserContext struct {
	UserID        *uuid.UUID
	Email         string
	Roles         []string
	Permissions   []string
	TraceID       string
	RequestID     string
	CorrelationID string
}

func FromGin(c *gin.Context) *UserContext {
	value, exists := c.Get(ContextUser)
	if !exists {
		return &UserContext{}
	}

	userCtx, ok := value.(*UserContext)
	if !ok {
		return &UserContext{}
	}

	return userCtx
}

func ParseCSVHeader(value string) []string {
	if value == "" {
		return []string{}
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
