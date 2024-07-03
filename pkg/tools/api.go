package tools

import (
	"github.com/gin-gonic/gin"
)

func Handle(handler func(ctx *gin.Context, codec Codec) error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		codec := ctx.MustGet(CodecCtx).(Codec)
		if err := handler(ctx, codec); err != nil {
			codec.Catch(ctx, err)
		}
	}
}
