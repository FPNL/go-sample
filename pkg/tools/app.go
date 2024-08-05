package tools

import (
	"context"
	"log/slog"
	"net/http"
)

func NewApp(server *http.Server) *App {
	return &App{
		server,
	}
}

type App struct {
	*http.Server
}

func (receiver *App) New(server *http.Server) App {
	return App{server}
}

func (receiver *App) Run(log *slog.Logger) error {
	log.Info("啟動 APP", "Addr", receiver.Server.Addr)
	return receiver.Server.ListenAndServe()
}

// Stop gracefully stops the application
func (receiver *App) Stop() error {
	return receiver.Server.Shutdown(context.TODO())
}
