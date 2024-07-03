package service

import (
	"context"
	"log/slog"

	"github.com/fpnl/go-sample/api"
	"github.com/fpnl/go-sample/biz"
	"github.com/fpnl/go-sample/pkg/logger"
)

func NewGreeterService(greeterUsecase *biz.GreeterUsecase, log *slog.Logger) *Greeter {
	return &Greeter{
		greeterUsecase,
		log,
	}
}

var _ api.GreeterServer = (*Greeter)(nil)

type Greeter struct {
	*biz.GreeterUsecase
	*slog.Logger
}

func (s Greeter) Greet(ctx context.Context, in *api.GreetRequest) (*api.GreetResponse, error) {
	message := s.GreeterUsecase.Hello(ctx, in.Name)

	log := ctx.Value(logger.CtxLogger).(*slog.Logger)
	log.Info("Greet", slog.String("message", message))

	return &api.GreetResponse{Message: message}, nil
}
