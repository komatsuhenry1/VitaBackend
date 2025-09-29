// router/user_routes.go
package router

import (
	"medassist/internal/di"
	"github.com/gin-gonic/gin"
	"medassist/middleware"
)
// 
func SetupUserRoutes(r *gin.RouterGroup, container *di.Container) {
	user := r.Group("/user")
	{
		user.GET("/dashboard", middleware.AuthUser(), container.UserHandler.UserDashboard) // retorna dados relevantes para o user
		user.GET("/all_nurses", middleware.AuthUser(), container.UserHandler.GetAllNurses) // get all nurses para agendar visita TODO
		user.POST("/visit", middleware.AuthUser(), container.UserHandler.CreateVisit) // agendamento de visita TODO
		user.GET("/visits", middleware.AuthUser(), container.UserHandler.GetAllVisits)
		user.GET("/file/:id", container.UserHandler.GetFileByID)
		user.POST("/contact", container.UserHandler.ContactUsMessage)
		user.GET("/nurse/:id", middleware.AuthUser(), container.UserHandler.GetNurseProfile)
	}
}