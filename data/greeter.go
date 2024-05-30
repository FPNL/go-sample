package data

import (
	"context"
	"log/slog"

	"oltp/biz"
)

func NewGreeterRepo(data *Data, log slog.Logger) biz.GreeterRepo {
	return &greeterRepo{
		data,
		log,
	}
}

var _ biz.GreeterRepo = (*greeterRepo)(nil)

type greeterRepo struct {
	data *Data
	log  slog.Logger
}

func (g greeterRepo) SayHi(_ context.Context) string {
	return "hi, "
}
