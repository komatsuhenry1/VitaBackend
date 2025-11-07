// router/nurse_routes.go
package router

import (
	"medassist/internal/di"
	"medassist/middleware"

	"github.com/gin-gonic/gin"
)

func SetupNurseRoutes(r *gin.RouterGroup, container *di.Container) {
	nurse := r.Group("/nurse")
	{
		nurse.GET("/dashboard", middleware.AuthNurse(), container.NurseHandler.NurseDashboard)         // dados relevenates para dashboard de nurse TODO
		nurse.PATCH("/online", middleware.AuthNurse(), container.NurseHandler.ChangeOnlineNurse)       // ativa online de nurse para receber chamadas de visitas DONE
		nurse.GET("/visits", middleware.AuthNurse(), container.NurseHandler.GetAllVisits)              // retorna todas visitas possiveis / marcadas TODO
		nurse.PATCH("/visit/:id", middleware.AuthNurse(), container.NurseHandler.ConfirmOrCancelVisit) // confirma que uma enfermeira ira para a visita
		nurse.GET("/patient/:id", middleware.AuthUserOrNurse(), container.NurseHandler.GetPatientProfile)
		nurse.PATCH("/update", middleware.AuthNurse(), container.NurseHandler.UpdateNurseProfile)
		nurse.DELETE("/delete", middleware.AuthNurse(), container.NurseHandler.DeleteNurseProfile)
		nurse.GET("/availability", middleware.AuthNurse(), container.NurseHandler.GetAvailabilityInfo)
		nurse.GET("/dashboard_info", middleware.AuthNurse(), container.NurseHandler.NurseDashboardData)
		nurse.GET("/my-profile", middleware.AuthNurse(), container.NurseHandler.GetMyNurseProfile)
		nurse.GET("/visit-info/:id", middleware.AuthNurse(), container.NurseHandler.GetNurseVisitInfo)
		nurse.PATCH("/service-confirmation/:id", middleware.AuthNurse(), container.NurseHandler.VisitServiceConfirmation)
		nurse.PATCH("/offline", middleware.AuthNurse(), container.NurseHandler.TurnOfflineOnLogout)
		nurse.PATCH("/reject-visit/:id", middleware.AuthNurse(), container.NurseHandler.RejectVisit)
		nurse.POST("/review/:id", middleware.AuthNurse(), container.NurseHandler.AddReview)
		nurse.POST("/stripe-onboarding", middleware.AuthNurse(), container.NurseHandler.SetupStripeOnboarding)
	}
}
