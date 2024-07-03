package tools

import (
	"net"
	"net/http"
	"time"
)

func NewTransport() *http.Transport {
	transport := http.DefaultTransport.(*http.Transport).Clone()

	transport.IdleConnTimeout = 120 * time.Second
	transport.DialContext = (&net.Dialer{
		Timeout:   15 * time.Second, // tcp 等待連線的超時
		KeepAlive: 15 * time.Second, // conn 連線探針
	}).DialContext

	return transport
}
