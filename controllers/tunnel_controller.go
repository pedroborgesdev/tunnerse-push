package controllers

import (
	"net/http"
	"tunnerse/config"
	"tunnerse/logger"
	"tunnerse/services"
	"tunnerse/utils"

	"github.com/gin-gonic/gin"
)

type TunnelController struct {
	tunnelService *services.TunnelService
}

func NewTunnelController() *TunnelController {
	return &TunnelController{
		tunnelService: services.NewTunnelService(),
	}
}

func (c *TunnelController) Register(ctx *gin.Context) {
	var request utils.RegisterRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(ctx, gin.H{
			"error": err.Error(),
		})
		return
	}

	tunnelName, err := c.tunnelService.Register(request.Name)

	if err != nil {
		if config.AppConfig.WARNS_ON_HTML {
			if err.Error() == "tunnel not found" {
				c.tunnelService.NotFound(ctx.Writer)
				return
			}
		}

		utils.BadRequest(ctx, gin.H{
			"error": err.Error(),
		})

		logger.Log("ERROR", "Registration failed", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
		})
		return
	}

	utils.Success(ctx, gin.H{
		"message": "tunnel has been registered",
		"tunnel":  tunnelName,
	})

	logger.Log("INFO", "User registered successfully", []logger.LogDetail{
		{Key: "tunnel", Value: tunnelName},
	})
}

func (c *TunnelController) Get(ctx *gin.Context) {
	name := utils.GetTunnelName(ctx)
	if name == "" {
		if config.AppConfig.WARNS_ON_HTML {
			c.tunnelService.Home(ctx.Writer)
			return
		}

		utils.Success(ctx, gin.H{
			"message": "Tunnerse is running :)",
		})
	}

	body, err := c.tunnelService.Get(name)

	if err != nil {
		if config.AppConfig.WARNS_ON_HTML {
			if err.Error() == "tunnel not found" {
				c.tunnelService.NotFound(ctx.Writer)
				return
			}
		}

		utils.BadRequest(ctx, gin.H{
			"error": err.Error(),
		})

		logger.Log("ERROR", "Tunneling failed", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
		})
		return
	}

	ctx.Writer.Write(body)

	logger.Log("INFO", "Messagem has been writed", []logger.LogDetail{
		{Key: "tunnel", Value: name},
	})
}

func (c *TunnelController) Response(ctx *gin.Context) {
	name := utils.GetTunnelName(ctx)
	if name == "" {
		if config.AppConfig.WARNS_ON_HTML {
			c.tunnelService.Home(ctx.Writer)
			return
		}

		utils.Success(ctx, gin.H{
			"message": "Tunnerse is running :)",
		})
	}

	err := c.tunnelService.Response(name, ctx.Request.Body)

	if err != nil {
		if config.AppConfig.WARNS_ON_HTML {
			if err.Error() == "tunnel not found" {
				c.tunnelService.NotFound(ctx.Writer)
				return
			}
		}

		utils.BadRequest(ctx, gin.H{
			"error": err.Error(),
		})

		logger.Log("ERROR", "Tunneling failed", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
		})
		return
	}

	ctx.Writer.WriteHeader(http.StatusOK)

	logger.Log("INFO", "Messagem has been writed", []logger.LogDetail{
		{Key: "tunnel", Value: name},
	})
}

func (c *TunnelController) Tunnel(ctx *gin.Context) {
	name := utils.GetTunnelName(ctx)
	if name == "" {
		if config.AppConfig.WARNS_ON_HTML {
			c.tunnelService.Home(ctx.Writer)
			return
		}

		utils.Success(ctx, gin.H{
			"message": "Tunnerse is running :)",
		})
	}

	err := c.tunnelService.Tunnel(name, ctx.Writer, ctx.Request)

	if err != nil {
		if config.AppConfig.WARNS_ON_HTML {
			if err.Error() == "tunnel not found" {
				c.tunnelService.NotFound(ctx.Writer)
				return
			}
			if err.Error() == "timeout" {
				c.tunnelService.Timeout(ctx.Writer)
				return
			}
		}

		utils.BadRequest(ctx, gin.H{
			"error": err.Error(),
		})

		logger.Log("ERROR", "Tunneling failed", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
		})
		return
	}

	logger.Log("INFO", "Mensagem has been written", []logger.LogDetail{
		{Key: "tunnel", Value: name},
	})
}

func (c *TunnelController) Close(ctx *gin.Context) {
	name := utils.GetTunnelName(ctx)
	if name == "" {
		if config.AppConfig.WARNS_ON_HTML {
			c.tunnelService.Home(ctx.Writer)
			return
		}

		utils.Success(ctx, gin.H{
			"message": "Tunnerse is running :)",
		})
	}

	err := c.tunnelService.Close(name)

	if err != nil {
		if config.AppConfig.WARNS_ON_HTML {
			if err.Error() == "tunnel not found" {
				c.tunnelService.NotFound(ctx.Writer)
				return
			}
		}

		utils.BadRequest(ctx, gin.H{
			"error": err.Error(),
		})

		logger.Log("ERROR", "Failed to delete tunnel", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
		})
		return
	}

	logger.Log("INFO", "Tunnel hass been deleted", []logger.LogDetail{
		{Key: "tunnel", Value: name},
	})
}
