// router/routes.go
package router

import (
	"medassist/internal/di"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)



func InitializeRoutes() *gin.Engine {
	container := di.NewContainer()
	router := gin.Default()
	go container.ChatHub.Run() //Isso inicia a execução de container.ChatHub.Run() em uma nova goroutine (de forma assíncrona e não bloqueante).


	router.Use(cors.New(cors.Config{
        AllowOrigins: []string{
			"*",
            "http://localhost:3000", // Para acesso local via localhost
            "http://192.168.18.139:3000", // Para acesso via IP na rede local
            "https://vita-frontend-uhje-ghicpj691-komatsuhenry-3753s-projects.vercel.app", // Sua URL de produção/staging
        },
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api/v1")

	SetupAuthRoutes(api, container)
	SetupUserRoutes(api, container)
	SetupNurseRoutes(api, container)
	SetupAdminRoutes(api, container)
	SetupWebsocketRoutes(router, container)
	SetupChatRoutes(api, container)

	return router
}
