package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"oltp/api"
	"oltp/conf"
	"oltp/server/middleware"
	"oltp/service"
)

func NewHTTPServer(
	cs *conf.Server,
	serviceGreeter *service.Greeter,
	midDefaultCodec *middleware.DefaultCodec,
	midIpWhitelist *middleware.IpWhitelist,
) *http.Server {
	r := gin.Default()
	r.UseH2C = true

	r.Use(midDefaultCodec.Mid())

	// mid 範例 1
	// mid 全域應用
	r.Use(midIpWhitelist.Mid())

	// mid 範例 2
	// 限制在特定 API
	r.Use(midIpWhitelist.Mid(api.GetGreeter))

	// mid 範例 3
	// 限制在特定的 server
	api.RegisterGinGreeterServer(r, serviceGreeter, midIpWhitelist.Mid())

	// mid 範例 4
	// 限制在特定的 server 中的特定 API
	api.RegisterGinGreeterServer(r, serviceGreeter, midIpWhitelist.Mid(api.GetGreeter))

	return &http.Server{
		Addr:         cs.HTTP.Addr,
		Handler:      h2c.NewHandler(r, &http2.Server{}),
		WriteTimeout: 60 * time.Second,
		ErrorLog:     nil,
	}
}
