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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // libera todas as rotas
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

	return router
}