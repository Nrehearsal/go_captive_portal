package webserver

import (
	"github.com/gin-gonic/gin"
)

func GatewaySSLOn() gin.HandlerFunc {
	return func(c *gin.Context) {
		tls := c.Request.TLS
		if tls != nil {
			c.Set("GatewaySSLOn", "yes")
		} else {
			c.Set("GatewaySSLOn", "no")
		}
		c.Header("Cache-Control", "no-cache")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "-1")
		c.Next()
		return
	}
}