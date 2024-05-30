package service

import (
	"context"

	"oltp/api"
	"oltp/biz"
)

func NewGreeterService(greeterUsecase *biz.GreeterUsecase) *Greeter {
	return &Greeter{
		greeterUsecase,
	}
}

var _ api.GreeterServer = (*Greeter)(nil)

type Greeter struct {
	*biz.GreeterUsecase
}

func (s Greeter) Greet(ctx context.Context, in *api.GreetRequest) (*api.GreetResponse, error) {
	message := s.GreeterUsecase.Hello(ctx, in.Name)

	return &api.GreetResponse{Message: message}, nil
}
