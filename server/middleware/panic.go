package middleware

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"

	"github.com/fpnl/go-sample/conf"
	"github.com/fpnl/go-sample/pkg/logger"
	"github.com/fpnl/go-sample/pkg/tools"
)

func NewRecovery(log *slog.Logger, cl *conf.Log) (*Recovery, func(), error) {
	panicLogFile, err := os.OpenFile(cl.PanicPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open panic.log: %w", err)
	}

	var panicWriter io.Writer = panicLogFile

	if cl.Stdout {
		panicWriter = io.MultiWriter(panicWriter, os.Stderr)
	}

	cleanup := func() {
		panicLogFile.Close()
		log.Info("panic log closed")
	}

	return &Recovery{
		log,
		panicWriter,
	}, cleanup, nil
}

type Recovery struct {
	*slog.Logger
	io.Writer
}

// ErrUnknownRequest is unknown request error.
var ErrUnknownRequest = errors.InternalServer("UNKNOWN", "unknown request error")

func (m *Recovery) Mid(allowPaths ...string) gin.HandlerFunc {
	return tools.Mid(allowPaths, func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				brokenPipe := isBrokenPipe(err)
				requestID := ctx.GetString(RequestIDHeaderKey)

				_, err = fmt.Fprintf(m.Writer,
					"[request: %s] panic recovered:\n%s\n%s",
					requestID,
					err,
					logger.Stack(3),
				)
				if err != nil {
					m.Logger.ErrorContext(ctx, "write panic log fail: %v", err)
				}

				err = ErrUnknownRequest

				if brokenPipe {
					// If the connection is dead, we can't write a status to it.
					_ = ctx.Error(err.(error))
					ctx.Abort()
				} else {
					m.Logger.ErrorContext(ctx, "panic: %v", err)
					ctx.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()

		ctx.Next()
	})
}

func isBrokenPipe(err any) bool {
	if ne, ok := err.(*net.OpError); ok {
		var se *os.SyscallError
		if errors.As(ne.Err, &se) {
			if strings.Contains(strings.ToLower(se.Error()),
				"broken pipe") || strings.Contains(strings.ToLower(se.Error()),
				"connection reset by peer") {
				return true
			}
		}
	}

	return false
}
