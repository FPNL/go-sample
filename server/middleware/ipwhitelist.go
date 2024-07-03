package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/fpnl/go-sample/pkg/tools"
)

func NewIpWhitelist() *IpWhitelist {
	return &IpWhitelist{}
}

type IpWhitelist struct {
}

func (m *IpWhitelist) Mid(allowPaths ...string) gin.HandlerFunc {
	return tools.Mid(allowPaths, func(c *gin.Context) {

	})
}
