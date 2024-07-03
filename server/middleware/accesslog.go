package middleware

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"

	"oltp/pkg/tools"
)

func NewAccessLog(log *slog.Logger) (*AccessLog, func(), error) {
	file, err := os.OpenFile("./log/access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open access.log: %w", err)
	}

	cleanup := func() {
		file.Close()
		log.Info("access log closed")
	}
	return &AccessLog{
		log,
		file,
	}, cleanup, nil
}

type AccessLog struct {
	*slog.Logger
	io.Writer
}

func (m *AccessLog) Mid(allowPaths ...string) gin.HandlerFunc {
	logger := slog.New(slog.NewTextHandler(m.Writer, nil))

	loggerConfig := sloggin.Config{
		WithUserAgent:      false,
		WithRequestID:      true,
		WithRequestBody:    false,
		WithRequestHeader:  false,
		WithResponseBody:   false,
		WithResponseHeader: false,
		WithSpanID:         true,
		WithTraceID:        true,
		Filters:            nil,
	}

	return tools.Mid(allowPaths, sloggin.NewWithConfig(logger, loggerConfig))
}
