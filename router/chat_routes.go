package router

import (
	"medassist/internal/di"
	"medassist/middleware"

	"github.com/gin-gonic/gin"
)

func SetupChatRoutes(api *gin.RouterGroup, container *di.Container) {
	handler := container.ChatHandler
	chatGroup := api.Group("/chat")

	chatGroup.Use(middleware.AuthUserOrNurse())

	{
		chatGroup.GET("/messages/:nurseId", handler.GetMessagesHistory)
		chatGroup.GET("/conversations", middleware.AuthNurse(), handler.GetConversations)
	}
}
