package routes

import (
	"net/http"
	"tunnerse/config"
	"tunnerse/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	tunnelController := controllers.NewTunnelController()

	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	tunnel := router.Group("/")
	// tunnel.Use(
	// 	middlewares.RateLimiter(middlewares.RateLimiterConfig{
	// 		RequestLimit:     15,
	// 		WindowMinutes:    60,
	// 		WarningThreshold: 3,
	// 		BlockTimeout:     20,
	// 	}),
	// )

	if config.AppConfig.SUBDOMAIN {
		tunnel.POST("/register", tunnelController.Register)
		tunnel.GET("/tunnel", tunnelController.Get)
		tunnel.POST("/response", tunnelController.Response)
		tunnel.POST("/close", tunnelController.Close)
		tunnel.GET("/", tunnelController.Tunnel)

		router.NoRoute(tunnelController.Tunnel)
	}

	if !config.AppConfig.SUBDOMAIN {
		tunnel.POST("/register", tunnelController.Register)
		tunnel.GET(":name/tunnel", tunnelController.Get)
		tunnel.POST(":name/response", tunnelController.Response)
		tunnel.POST(":name/close", tunnelController.Close)
		tunnel.GET(":name/", tunnelController.Tunnel)

		router.NoRoute(tunnelController.Tunnel)
	}
}
