package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"

	"oltp/conf"
	"oltp/pkg/logger"
)

const (
	version string = "{{VERSION}}"
)

var (
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "conf/env.json", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()

	config, err := conf.InitAPI(flagconf)
	if err != nil {
		panic(err)
	}

	log := logger.NewLogger()

	app, cleanup, err := initApp(config.Project, config.Server, config.Data, log)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err = app.Run(log); err != nil {
		panic(err)
	}
}

func newApp(server *http.Server) *app {
	return &app{
		server,
	}
}

type app struct {
	*http.Server
}

func (receiver *app) New(server *http.Server) app {
	return app{server}
}

func (receiver *app) Run(log *slog.Logger) error {
	log.Info("啟動 APP", "Addr", receiver.Server.Addr)
	return receiver.Server.ListenAndServe()
}

// Stop gracefully stops the application
func (receiver *app) Stop() error {
	return receiver.Server.Shutdown(context.TODO())
}
