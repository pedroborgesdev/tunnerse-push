package routes

import (
	"net/http"
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

	tunnel.POST("/register", tunnelController.Register)
	tunnel.GET("/tunnel", tunnelController.Get)
	tunnel.POST("/response", tunnelController.Response)
	tunnel.POST("/close", tunnelController.Close)
	tunnel.GET("/", tunnelController.Tunnel)
}
