// router/auth_routes.go
package router

import (
	"medassist/internal/di"
	"medassist/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.RouterGroup, container *di.Container) {
	auth := r.Group("/auth")
	{
		auth.POST("/adm", container.AuthHandler.FirstLoginAdmin)
		auth.POST("/user", container.AuthHandler.UserRegister)
		auth.POST("/nurse", container.AuthHandler.NurseRegister)
		auth.POST("/email", container.AuthHandler.SendEmailForgotPassword)
		auth.PATCH("/code", container.AuthHandler.SendCode) //comentario pra subir certo
		auth.POST("/validate", container.AuthHandler.ValidateCode)
		//auth.POST("/reset-password", container.AuthHandler.ChangePasswordUnlogged)
		auth.POST("/reset-password", container.AuthHandler.ResetPassword)
		auth.POST("/login", container.AuthHandler.LoginUser)
		auth.PATCH("/logged/password", middleware.AuthUserOrNurse(), container.AuthHandler.ChangePasswordLogged)
		auth.POST("/validate-token", container.AuthHandler.ValidateResetToken)
	}
}
