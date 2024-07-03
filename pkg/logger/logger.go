package logger

import (
	"log/slog"
	"os"
)

const CtxLogger string = "Logger"

func NewLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
