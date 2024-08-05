package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"

	sloglogrus "github.com/samber/slog-logrus/v2"
	"github.com/sirupsen/logrus"

	"github.com/fpnl/go-sample/conf"
)

const CtxLogger string = "Logger"

func NewLogger(cl *conf.Log, cp *conf.Project) (*slog.Logger, error) {
	logrusLogger, err := NewLogrus(cl)
	if err != nil {
		return nil, fmt.Errorf("new logrus logger fail: %w", err)
	}

	option := sloglogrus.Option{
		Level:     slog.LevelDebug,
		Logger:    logrusLogger,
		AddSource: true,
	}

	if cp.IsDebug {
		option.Level = slog.LevelDebug
	} else {
		option.Level = slog.LevelInfo
	}

	logger := slog.New(
		option.NewLogrusHandler(),
	)

	// slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return logger, nil
}

func NewLogrus(cl *conf.Log) (logrusLogger *logrus.Logger, err error) {
	logrusLogger = logrus.New()
	logrusLogger.SetFormatter(NewJSONFormatter())

	if err = EnsureDir(cl.OutPath); err != nil {
		return nil, fmt.Errorf("ensure log dir fail: %w", err)
	}

	var outLog io.Writer

	if outLog, err = os.OpenFile(cl.OutPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		return nil, fmt.Errorf("initial log file fail: %w", err)
	}

	if cl.Stdout {
		outLog = io.MultiWriter(outLog, os.Stdout)
	}

	logrusLogger.SetOutput(outLog)

	return logrusLogger, nil
}

func NewJSONFormatter() *logrus.JSONFormatter {
	return &logrus.JSONFormatter{
		TimestampFormat:   time.DateTime,
		DisableHTMLEscape: true,
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			_, err := os.Getwd()
			if err == nil {
				pcs := make([]uintptr, 50)
				_ = runtime.Callers(7, pcs)
				frames := runtime.CallersFrames(pcs)

			outter:
				for next, again := frames.Next(); again; next, again = frames.Next() {
					if strings.HasSuffix(next.File, "logger_do_not_rename.go") {
						for find, more := frames.Next(); more; find, more = frames.Next() {
							if !strings.HasSuffix(find.File, "logger_do_not_rename.go") {
								return find.Function, fmt.Sprintf("%s:%d", find.File, find.Line)
							}
						}

						break outter
					}
				}
			}

			return f.Function, fmt.Sprintf("%s:%d", f.File, f.Line)
		},
		PrettyPrint: false,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "Time",
			logrus.FieldKeyLevel: "Level",
			logrus.FieldKeyMsg:   "Message",
			logrus.FieldKeyFunc:  "Caller",
		},
	}
}
