package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/fpnl/go-sample/api"
	"github.com/fpnl/go-sample/conf"
	"github.com/fpnl/go-sample/server/middleware"
	"github.com/fpnl/go-sample/service"
)

func NewHTTPServer(
	cp *conf.Project,
	cs *conf.Server,
	serviceGreeter *service.Greeter,
	midDefaultCodec *middleware.DefaultCodec,
	midIpWhitelist *middleware.IpWhitelist,
	midRequestUUID *middleware.RequestUUID,
	midAccessLog *middleware.AccessLog,
	midPanicLog *middleware.Recovery,
) *http.Server {
	if !cp.IsDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.UseH2C = true

	r.Use(midRequestUUID.Mid())
	r.Use(midPanicLog.Mid())
	r.Use(midAccessLog.Mid())
	r.Use(midDefaultCodec.Mid())

	// mid 範例 1
	// mid 全域應用
	r.Use(midIpWhitelist.Mid())

	// mid 範例 2
	// 限制在特定 API
	// r.Use(midIpWhitelist.Mid(api.GetGreeter))

	// mid 範例 3
	// 限制在特定的 server
	// api.RegisterGinGreeterServer(r, serviceGreeter, midIpWhitelist.Mid())

	// mid 範例 4
	// 限制在特定的 server 中的特定 API
	// api.RegisterGinGreeterServer(r, serviceGreeter, midIpWhitelist.Mid(api.GetGreeter))

	api.RegisterGinGreeterServer(r, serviceGreeter)

	return &http.Server{
		Addr:         cs.HTTP.Addr,
		Handler:      h2c.NewHandler(r, &http2.Server{}),
		WriteTimeout: 60 * time.Second,
		ErrorLog:     nil,
	}
}
