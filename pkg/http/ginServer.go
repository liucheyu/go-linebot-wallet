package http

import (
	"github.com/gin-gonic/gin"
	"github.com/liucheyu/go-linebot-wallet/pkg/http/health"
	"github.com/liucheyu/go-linebot-wallet/pkg/linebot"
)

func GinServer() {
	router := gin.Default()

	router.GET("/health/alive", health.Alive)
	router.POST("/line/bot/callback", linebot.Callback)

	router.Run(":3000")
}
