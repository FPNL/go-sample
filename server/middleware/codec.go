package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/fpnl/go-sample/pkg/tools"
)

func NewDefaultCodec() *DefaultCodec {
	return &DefaultCodec{}
}

// DefaultCodec 將編碼轉換成 service 需要的 input/output
type DefaultCodec struct{}

func (c *DefaultCodec) Mid(allowPaths ...string) gin.HandlerFunc {
	return tools.Mid(allowPaths, func(ctx *gin.Context) {
		ctx.Set(tools.CodecCtx, c)
	})
}

func (c *DefaultCodec) Bind(ctx *gin.Context, in any) error {
	return ctx.ShouldBind(in)
}

func (c *DefaultCodec) BindVars(ctx *gin.Context, in any) error {
	return ctx.ShouldBindUri(in)
}

func (c *DefaultCodec) BindQuery(ctx *gin.Context, in any) error {
	return ctx.BindQuery(in)
}

func (c *DefaultCodec) BindForm(ctx *gin.Context, in any) error {
	return ctx.ShouldBindWith(in, binding.Form)
}

// Result only supports JSON, XML, and text/plain, default is JSON.
func (c *DefaultCodec) Result(ctx *gin.Context, out any) error {
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

func (c *DefaultCodec) Catch(ctx *gin.Context, err error) {
	// TODO implement me
	panic("implement me")
}
