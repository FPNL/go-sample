package middleware

import (
	"github.com/gin-gonic/gin"

	"oltp/pkg/axiom"
)

func NewIpWhitelist() *IpWhitelist {
	return &IpWhitelist{}
}

type IpWhitelist struct {
}

func (m *IpWhitelist) Mid(allowPaths ...string) gin.HandlerFunc {
	return axiom.Mid(allowPaths, func(c *gin.Context) {

	})
}
