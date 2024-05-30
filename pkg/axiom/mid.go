package axiom

import "github.com/gin-gonic/gin"

func Mid(allowPaths []string, m gin.HandlerFunc) gin.HandlerFunc {
	table := make(map[string]struct{})

	for _, path := range allowPaths {
		table[path] = struct{}{}
	}

	return func(c *gin.Context) {
		if _, ok := table[c.FullPath()]; len(table) > 0 && !ok {
			c.Next()
			return
		}
		m(c)
	}
}
