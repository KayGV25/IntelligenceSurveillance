package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/common/security"
)

func UserContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDHeader := c.GetHeader(security.HeaderUserID)

		var userID *uuid.UUID
		if userIDHeader != "" {
			parsed, err := uuid.Parse(userIDHeader)
			if err == nil {
				userID = &parsed
			}
		}

		userCtx := &security.UserContext{
			UserID:        userID,
			Email:         c.GetHeader(security.HeaderUserEmail),
			Roles:         security.ParseCSVHeader(c.GetHeader(security.HeaderUserRoles)),
			Permissions:   security.ParseCSVHeader(c.GetHeader(security.HeaderUserPerms)),
			TraceID:       c.GetHeader(security.HeaderTraceID),
			RequestID:     c.GetHeader(security.HeaderRequestID),
			CorrelationID: c.GetHeader(security.HeaderCorrelationID),
		}

		c.Set(security.ContextUser, userCtx)

		c.Next()
	}
}
