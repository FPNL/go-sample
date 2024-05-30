package axiom

import "github.com/gin-gonic/gin"

const CodecCtx = "context_codec"

type Codec interface {
	Bind(*gin.Context, any) error
	BindVars(*gin.Context, any) error
	BindQuery(*gin.Context, any) error
	BindForm(*gin.Context, any) error
	Result(*gin.Context, any) error
	Catch(ctx *gin.Context, err error)
}
