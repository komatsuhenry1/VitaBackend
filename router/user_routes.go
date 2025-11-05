// router/user_routes.go
package router

import (
	"medassist/internal/di"
	"medassist/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.RouterGroup, container *di.Container) {
	user := r.Group("/user")
	{
		user.GET("/all_nurses", middleware.AuthUser(), container.UserHandler.GetAllNurses)  // get all nurses para agendar visita TODO
		user.POST("/visit", middleware.AuthUser(), container.UserHandler.VisitSolicitation) // agendamento de visita TODO
		user.POST("/immediate-visit", middleware.AuthUser(), container.UserHandler.ImmediateVisitSolicitation)
		user.PATCH("/visit/:id", middleware.AuthUser(), container.UserHandler.ConfirmVisitService)
		user.GET("/visits", middleware.AuthUser(), container.UserHandler.GetAllVisits)
		user.GET("/file/:id", container.UserHandler.GetFileByID)
		user.POST("/contact", container.UserHandler.ContactUsMessage)
		user.GET("/nurse/:id", middleware.AuthUserOrNurse(), container.UserHandler.GetNurseProfile)
		user.GET("/my-profile", middleware.AuthUser(), container.UserHandler.GetMyUserProfile)
		user.PATCH("/update", middleware.AuthUser(), container.UserHandler.UpdateUser)
		user.DELETE("/delete", middleware.AuthUser(), container.UserHandler.DeleteUser)
		user.GET("/online_nurses", middleware.AuthUser(), container.UserHandler.GetOnlineNurses)
		user.GET("/visit-info/:id", middleware.AuthUser(), container.UserHandler.GetPatientVisitInfo)
		user.POST("/review/:id", middleware.AuthUser(), container.UserHandler.AddReview)
	}
}
