package utils

import (
	"net"
	"strings"
	"tunnerse/config"

	"github.com/gin-gonic/gin"
)

func GetTunnelName(ctx *gin.Context) string {
	name := ""
	host := ctx.Request.Host

	if !config.AppConfig.SUBDOMAIN {
		name = ctx.Param("name")
	} else {
		if h, _, err := net.SplitHostPort(host); err == nil {
			host = h
		}

		parts := strings.Split(host, ".")
		if len(parts) >= 3 {
			name = parts[0]
		}
	}

	return name
}
