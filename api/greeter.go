package api

import (
	"context"

	"github.com/gin-gonic/gin"

	"oltp/pkg/axiom"
)

const (
	GetGreeter = "/api/v1/greeter"
)

// GreeterServer 有幾個 API 就應該要有幾個 method
type GreeterServer interface {
	Greet(ctx context.Context, in *GreetRequest) (*GreetResponse, error)
}

// RegisterGinGreeterServer 採用 Gin 當作 router
func RegisterGinGreeterServer(router gin.IRouter, service GreeterServer, middlewares ...gin.HandlerFunc) {
	group := router.Group("")
	group.Use(middlewares...)
	// 按需求新增 router
	group.GET(GetGreeter, greeterHandler(service))
}

// GreetRequest 按需求修改
type GreetRequest struct {
	Name string `json:"name"`
}

// GreetResponse 按需求修改，如果是 text/plain 則改成 string
type GreetResponse struct {
	Message string `json:"message"`
}

// greeterHandler
func greeterHandler(server GreeterServer) gin.HandlerFunc {
	return axiom.Handle(func(ctx *gin.Context, codec axiom.Codec) (err error) {
		// 這裏應為固定
		var in *GreetRequest
		var out *GreetResponse

		// 這行按需求修改
		if err = codec.Bind(ctx, &in); err != nil {
			return err
		}

		// 這行按需求修改
		out, err = server.Greet(ctx, in)
		if err != nil {
			return err
		}

		return codec.Result(ctx, out)
	})
}
