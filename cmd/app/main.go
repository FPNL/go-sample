package main

import (
	"flag"

	"github.com/fpnl/go-sample/conf"
	"github.com/fpnl/go-sample/pkg/logger"
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

	log, err := logger.NewLogger(config.Log, config.Project)
	if err != nil {
		panic(err)
	}

	app, cleanup, err := initApp(config.Project, config.Server, config.Data, config.Log, log)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err = app.Run(log); err != nil {
		panic(err)
	}
}
