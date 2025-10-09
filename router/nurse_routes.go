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
		nurse.GET("/dashboard", middleware.AuthNurse(), container.NurseHandler.NurseDashboard) // dados relevenates para dashboard de nurse TODO
		nurse.PATCH("/online", middleware.AuthNurse(), container.NurseHandler.ChangeOnlineNurse) // ativa online de nurse para receber chamadas de visitas DONE
		nurse.GET("/visits", middleware.AuthNurse(), container.NurseHandler.GetAllVisits) // retorna todas visitas possiveis / marcadas TODO
		nurse.PATCH("/visit/:id", middleware.AuthNurse(), container.NurseHandler.ConfirmOrCancelVisit) // confirma que uma enfermeira ira para a visita
		nurse.GET("/patient/:id", middleware.AuthUserOrNurse(), container.NurseHandler.GetPatientProfile)
		nurse.PATCH("/update", middleware.AuthNurse(), container.NurseHandler.UpdateNurseProfile)
		nurse.DELETE("/delete", middleware.AuthNurse(), container.NurseHandler.DeleteNurseProfile)
	}
}