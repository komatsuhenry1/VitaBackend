package router

import(
	"github.com/gin-gonic/gin"
	"medassist/internal/di"
	"medassist/internal/chat"
)

func SetupWebsocketRoutes(router *gin.Engine, container *di.Container) {
	hub := container.ChatHub

	ws := router.Group("/ws")
	{
		ws.GET("/chat", func(c *gin.Context) {
			chat.ServeWs(hub, c)	
		})
	}
}
