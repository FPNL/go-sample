package biz

import (
	"context"
	"log/slog"
)

type Greeter struct {
	Message string
}

type GreeterRepo interface {
	SayHi(ctx context.Context) string
}

func NewGreeterUsecase(repo GreeterRepo, logger slog.Logger) *GreeterUsecase {
	return &GreeterUsecase{repo: repo, logger: logger}
}

type GreeterUsecase struct {
	repo   GreeterRepo
	logger slog.Logger
}

func (s GreeterUsecase) Hello(ctx context.Context, name string) string {
	return s.repo.SayHi(ctx) + name
}
