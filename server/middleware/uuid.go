package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/fpnl/go-sample/pkg/logger"
	"github.com/fpnl/go-sample/pkg/tools"
)

const RequestIDHeaderKey = "X-Request-Id"

func NewRequestUUID(logger *slog.Logger) *RequestUUID {
	return &RequestUUID{
		logger: logger,
	}
}

type RequestUUID struct {
	logger *slog.Logger
}

func (m *RequestUUID) Mid(allowPaths ...string) gin.HandlerFunc {
	return tools.Mid(allowPaths, func(c *gin.Context) {
		if c.GetHeader(RequestIDHeaderKey) != "" {
			return
		}

		obj := struct {
			RequestUUID string `json:"request_uuid"`
		}{}

		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			if err := c.ShouldBindBodyWithJSON(&obj); err != nil {
				c.String(500, err.Error())
				c.Abort()
				return
			}
		} else {
			if err := c.ShouldBindQuery(&obj); err != nil {
				c.String(500, err.Error())
				c.Abort()
				return
			}
		}

		requestUUID := obj.RequestUUID

		if requestUUID == "" {
			requestUUID = uuid.New().String()
		}

		c.Header(RequestIDHeaderKey, requestUUID)
		c.Set(logger.CtxLogger, m.logger.With("id", requestUUID))
	})
}
