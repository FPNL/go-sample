package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"

	"oltp/conf"
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

	// TODO: logger

	app, cleanup, err := initApp(config.Server, config.Data, slog.Logger{})
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err = app.Run(); err != nil {
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

func (receiver *app) Run() error {
	return receiver.Server.ListenAndServe()
}

// Stop gracefully stops the application
func (receiver *app) Stop() error {
	return receiver.Server.Shutdown(context.TODO())
}
