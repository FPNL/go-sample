package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/fpnl/go-sample/pkg/tools"
)

/*
這裡示範了針對不同的供應商可以有不同的編碼方式
*/

func NewUUCodec() *UUCodec {
	return &UUCodec{}
}

func Mid(allowPaths ...string) gin.HandlerFunc {
	return tools.Mid(allowPaths, func(ctx *gin.Context) {
		ctx.Set(tools.CodecCtx, &UUCodec{})
	})
}

type UUCodec struct{}

func (c *UUCodec) Catch(ctx *gin.Context, err error) {
	// TODO implement me
	panic("implement me")
}

func (c *UUCodec) Bind(ctx *gin.Context, in any) error {
	return ctx.ShouldBind(in)
}

func (c *UUCodec) BindVars(ctx *gin.Context, in any) error {
	return ctx.ShouldBindUri(in)
}

func (c *UUCodec) BindQuery(ctx *gin.Context, in any) error {
	return ctx.BindQuery(in)
}

func (c *UUCodec) BindForm(ctx *gin.Context, in any) error {
	return ctx.ShouldBindWith(in, binding.Form)
}

// Result only supports JSON, XML, and text/plain, default is JSON.
func (c *UUCodec) Result(ctx *gin.Context, out any) error {
	switch ctx.ContentType() {
	case "application/xml":
		ctx.XML(200, out)
	case "application/json":
		ctx.JSON(200, out)
	case "text/plain":
		s, ok := out.(string)
		if !ok {
			return fmt.Errorf("invalid out for text/plain")
		}
		ctx.String(200, s)
	default:
		ctx.JSON(200, out)
	}

	return nil
}
